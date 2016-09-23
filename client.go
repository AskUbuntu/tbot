package main

import (
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

// Client sends tweets as soon as they are ready.
type Client struct {
	client    *twitter.Client
	closeChan chan bool
}

// run receives messages on the specified channel and tweets them.
func (c *Client) run(ch <-chan string) {
	for t := range ch {
		if _, _, err := c.client.Statuses.Update(t, nil); err != nil {
			// TODO: some sort of error should be shown
		}
	}
	close(c.closeChan)
}

// NewClient creates a new Twitter client. The credentials are checked to
// ensure that they are valid.
func NewClient(config *Config, ch <-chan string) (*Client, error) {
	twitterConfig := oauth1.NewConfig(
		config.TwitterConsumerKey,
		config.TwitterConsumerSecret,
	)
	token := oauth1.NewToken(
		config.TwitterAccessToken,
		config.TwitterAccessSecret,
	)
	httpClient := twitterConfig.Client(oauth1.NoContext, token)
	c := &Client{
		client:    twitter.NewClient(httpClient),
		closeChan: make(chan bool),
	}
	params := &twitter.AccountVerifyParams{
		SkipStatus: twitter.Bool(true),
	}
	if _, _, err := c.client.Accounts.VerifyCredentials(params); err != nil {
		return nil, err
	}
	go c.run(ch)
	return c, nil
}

// Waits for the client to shutdown. The channel passed to NewClient *must* be
// closed first.
func (c *Client) Close() {
	<-c.closeChan
}
