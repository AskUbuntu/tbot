package queue

import (
	"github.com/AskUbuntu/tbot/scraper"
	"github.com/AskUbuntu/tbot/util"

	"sync"
	"time"
)

type data struct {
	sync.Mutex
	name           string
	LastMessage    time.Time          `json:"last_message"`    // Time of last message
	QueuedMessages []*scraper.Message `json:"queued_messages"` // List of messages in queue
}

func (d *data) load() error {
	_, err := util.LoadJSON(d.name, d)
	return err
}

func (d *data) save() error {
	return util.SaveJSON(d.name, d)
}
