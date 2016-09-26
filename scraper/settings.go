package scraper

import (
	"github.com/AskUbuntu/tbot/util"

	"sync"
)

// Settings controls the behavior of the scraper. This includes such things as
// the room, polling rate, and matching criteria.
type Settings struct {
	PollURL       string   `json:"poll_url"`       // Base URL for chat
	PollRoomID    int      `json:"poll_room_id"`   // Room to
	PollFrequency int      `json:"poll_frequency"` // Time in minutes between polling
	MinStars      int      `json:"min_stars"`      // Select messages with these many stars
	MatchingWords []string `json:"matching_words"` // Select messages with these words
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
		s.PollURL = "http://chat.stackexchange.com"
		s.PollRoomID = 201
		s.PollFrequency = 60
		s.MinStars = 4
		s.MatchingWords = []string{"ಠ_ಠ", "bacon"}
	}
	return nil
}

func (s *settings) save() error {
	return util.SaveJSON(s.name, s)
}
