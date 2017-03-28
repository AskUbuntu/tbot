package twitter

import (
	"bytes"
	"encoding/base64"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/AskUbuntu/tbot/config"
	"github.com/AskUbuntu/tbot/scraper"
	"github.com/ChimeraCoder/anaconda"
)

// Client sends tweets as soon as they are ready.
type Twitter struct {
	api     *anaconda.TwitterApi
	data    *data
	trigger chan bool
}

// retrieveImage fetches a remote image and returns its base64 content.
func retrieveImage(resource string) (string, error) {
	log.Printf("retrieving '%s'...", resource)
	req, err := http.NewRequest(http.MethodGet, resource, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "Ask Ubuntu Twitter Bot")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	var (
		buffer = bytes.Buffer{}
		w      = base64.NewEncoder(base64.StdEncoding, &buffer)
	)
	if _, err := io.Copy(w, res.Body); err != nil {
		return "", err
	}
	if err := w.Close(); err != nil {
		return "", err
	}
	return buffer.String(), nil
}

// run receives messages on the specified channel and tweets them.
func (t *Twitter) run(ch <-chan *scraper.Message) {
	for {
		var quit = false
		select {
		case m := <-ch:
			v := url.Values{}
			if m.Onebox == scraper.OneboxImage {
				// TODO: do something else if this fails
				img, err := retrieveImage(m.Body)
				if err != nil {
					log.Printf("HTTP error '%s'", err.Error())
					break
				}
				media, err := t.api.UploadMedia(img)
				if err != nil {
					log.Printf("twitter API error '%s'", err.Error())
					break
				}
				v.Set("media_ids", media.MediaIDString)
			}
			tweet, err := t.api.PostTweet(m.String(), v)
			if err != nil {
				log.Printf("twitter API error '%s'", err.Error())
				break
			}
			t.data.Lock()
			if len(t.data.Tweets) > 9 {
				t.data.Tweets = t.data.Tweets[:9]
			}
			t.data.Tweets = append([]*Tweet{
				&Tweet{
					Message:   m,
					TweetID:   tweet.Id,
					TweetTime: time.Now(),
				},
			},
				t.data.Tweets...,
			)
			if err := t.data.save(); err != nil {
				log.Printf("twitter serialization error '%s'", err.Error())
			}
			t.data.Unlock()
		case quit = <-t.trigger:
		}
		if quit {
			break
		}
	}
	close(t.trigger)
}

// New creates a new Twitter client. The credentials are checked to ensure that
// they are valid.
func New(config *config.Config, ch <-chan *scraper.Message) (*Twitter, error) {
	anaconda.SetConsumerKey(config.TwitterConsumerKey)
	anaconda.SetConsumerSecret(config.TwitterConsumerSecret)
	t := &Twitter{
		api: anaconda.NewTwitterApi(
			config.TwitterAccessToken,
			config.TwitterAccessSecret,
		),
		data:    &data{name: path.Join(config.DataPath, "twitter_data.json")},
		trigger: make(chan bool),
	}
	if ok, err := t.api.VerifyCredentials(); err != nil {
		return nil, err
	} else if !ok {
		return nil, errors.New("invalid Twitter credentials")
	}
	if err := t.data.load(); err != nil {
		return nil, err
	}
	go t.run(ch)
	return t, nil
}

// Tweets retrieves recently tweeted messages.
func (t *Twitter) Tweets() []*Tweet {
	t.data.Lock()
	defer t.data.Unlock()
	return t.data.Tweets
}

// Send posts a single status update.
func (t *Twitter) Send(status string) error {
	_, err := t.api.PostTweet(status, nil)
	return err
}

// Delete the tweet with the given ID.
func (t *Twitter) Delete(tweetID int64) error {
	_, err := t.api.DeleteTweet(tweetID, true)
	if err != nil {
		return err
	}
	t.data.Lock()
	defer t.data.Unlock()
	for i, tweet := range t.data.Tweets {
		if tweet.TweetID == tweetID {
			t.data.Tweets = append(
				t.data.Tweets[:i],
				t.data.Tweets[i+1:]...,
			)
			if err := t.data.save(); err != nil {
				return err
			}
			break
		}
	}
	return nil
}

// Waits for the client to shutdown. The channel passed to New *must* be
// closed first.
func (t *Twitter) Close() {
	t.trigger <- true
	<-t.trigger
}
