package scraper

import (
	"github.com/AskUbuntu/tbot/util"

	"sync"
	"time"
)

type data struct {
	sync.Mutex
	name         string
	LastScrape   time.Time  `json:"last_scrape"`   // Time of last scrape
	EarliestID   int        `json:"earliest_id"`   // ID of first message in last scrape
	Messages     []*Message `json:"messages"`      // Messages matching the criteria
	MessagesUsed []int      `json:"messages_used"` // Messages used from last period
}

func (s data) load() error {
	_, err := util.LoadJSON(s.name, s)
	return err
}

func (s *data) save() error {
	return util.SaveJSON(s.name, s)
}
