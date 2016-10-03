package twitter

import (
	"github.com/AskUbuntu/tbot/config"
	"github.com/AskUbuntu/tbot/scraper"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"

	"log"
	"path"
	"time"
)

// Client sends tweets as soon as they are ready.
type Twitter struct {
	client  *twitter.Client
	data    *data
	trigger chan bool
}

// run receives messages on the specified channel and tweets them.
func (t *Twitter) run(ch <-chan *scraper.Message) {
	for {
		quit := false
		select {
		case m := <-ch:
			log.Printf("tweeting '%s'", m.String())
			tweet, _, err := t.client.Statuses.Update(m.String(), nil)
			if err != nil {
				log.Printf("twitter API error '%s'", err.Error())
			} else {
				t.data.Lock()
				if len(t.data.Tweets) > 9 {
					t.data.Tweets = t.data.Tweets[:9]
				}
				t.data.Tweets = append(t.data.Tweets, &Tweet{
					Message:   m,
					TweetID:   tweet.ID,
					TweetTime: time.Now(),
				})
				if err := t.data.save(); err != nil {
					log.Printf("twitter serialization error '%s'", err.Error())
				}
				t.data.Unlock()
			}
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
	twitterConfig := oauth1.NewConfig(
		config.TwitterConsumerKey,
		config.TwitterConsumerSecret,
	)
	token := oauth1.NewToken(
		config.TwitterAccessToken,
		config.TwitterAccessSecret,
	)
	httpClient := twitterConfig.Client(oauth1.NoContext, token)
	t := &Twitter{
		client:  twitter.NewClient(httpClient),
		data:    &data{name: path.Join(config.DataPath, "twitter_data.json")},
		trigger: make(chan bool),
	}
	params := &twitter.AccountVerifyParams{
		SkipStatus: twitter.Bool(true),
	}
	_, _, err := t.client.Accounts.VerifyCredentials(params)
	if err != nil {
		return nil, err
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
	_, _, err := t.client.Statuses.Update(status, nil)
	return err
}

// Waits for the client to shutdown. The channel passed to New *must* be
// closed first.
func (t *Twitter) Close() {
	t.trigger <- true
	<-t.trigger
}
