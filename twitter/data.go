package twitter

import (
	"github.com/AskUbuntu/tbot/scraper"
	"github.com/AskUbuntu/tbot/util"

	"sync"
	"time"
)

// Tweet represents an individual tweet and its accompanying message.
type Tweet struct {
	Message   *scraper.Message `json:"message"`
	TweetID   int64            `json:"tweet_id"`
	TweetTime time.Time        `json:"tweet_time"`
}

type data struct {
	sync.Mutex
	name   string
	Tweets []*Tweet `json:"tweets"`
}

func (d *data) load() error {
	_, err := util.LoadJSON(d.name, d)
	return err
}

func (d *data) save() error {
	return util.SaveJSON(d.name, d)
}
