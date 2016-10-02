package queue

import (
	"github.com/AskUbuntu/tbot/util"

	"sync"
)

// Settings controls the behavior of the queue. This includes the interval
// between consecutive tweets for the account.
type Settings struct {
	QueueFrequency int `json:"queue_frequency"` // Time in minutes between tweets
}

type settings struct {
	Settings
	sync.Mutex
	name string
}

func (s *settings) load() error {
	e, err := util.LoadJSON(s.name, s)
	if err != nil {
		return err
	}
	if !e {
		s.QueueFrequency = 480 // 8 hours
	}
	return nil
}

func (s *settings) save() error {
	return util.SaveJSON(s.name, s)
}
