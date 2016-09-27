package queue

import (
	"github.com/AskUbuntu/tbot/util"

	"sync"
	"time"
)

type data struct {
	sync.Mutex
	name         string
	LastMessage  time.Time `json:"last_message"` // Time of last message
	QueuedTweets []string  `json:"tweets"`       // List of tweets in queue
}

func (d *data) load() error {
	_, err := util.LoadJSON(d.name, d)
	return err
}

func (d *data) save() error {
	return util.SaveJSON(d.name, d)
}
