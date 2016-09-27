package twitter

import (
	"github.com/AskUbuntu/tbot/config"
	"github.com/AskUbuntu/tbot/scraper"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

// Client sends tweets as soon as they are ready.
type Twitter struct {
	client    *twitter.Client
	closeChan chan bool
}

// run receives messages on the specified channel and tweets them.
func (t *Twitter) run(ch <-chan *scraper.Message) {
	for {
		quit := false
		select {
		case m := <-ch:
			if _, _, err := t.client.Statuses.Update(m.String(), nil); err != nil {
				// TODO: some sort of error should be shown
			}
		case <-t.closeChan:
			quit = true
		}
		if quit {
			break
		}
	}
	close(t.closeChan)
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
		client:    twitter.NewClient(httpClient),
		closeChan: make(chan bool),
	}
	params := &twitter.AccountVerifyParams{
		SkipStatus: twitter.Bool(true),
	}
	if _, _, err := t.client.Accounts.VerifyCredentials(params); err != nil {
		return nil, err
	}
	go t.run(ch)
	return t, nil
}

// Waits for the client to shutdown. The channel passed to New *must* be
// closed first.
func (t *Twitter) Close() {
	t.closeChan <- true
	<-t.closeChan
}
