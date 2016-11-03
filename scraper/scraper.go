package scraper

import (
	"github.com/AskUbuntu/tbot/config"

	"errors"
	"path"
	"reflect"
	"time"
)

// Scraper regularly scrapes the chat transcript for the current and previous
// day for messages that match the criteria in use. Once a message matches, it
// is added to the list of candidates for tweeting. The IDs of messages that
// are used is kept to prevent duplicates.
type Scraper struct {
	data     *data
	settings *settings
	trigger  chan bool
}

func (s *Scraper) run() {
	for {
		s.data.Lock()
		lastScrape := s.data.LastScrape
		s.data.Unlock()
		s.settings.Lock()
		pollFrequency := s.settings.PollFrequency
		s.settings.Unlock()
		var (
			now      = time.Now()
			duration = time.Duration(pollFrequency)
			diff     = lastScrape.Add(duration * time.Minute).Sub(now)
		)
		if diff <= 0 {
			// TODO: log scrape error
			s.scrape()
			diff = duration
		}
		var (
			timer = time.NewTimer(diff)
			quit  = false
		)
		select {
		case <-timer.C:
		case quit = <-s.trigger:
		}
		timer.Stop()
		if quit {
			break
		}
	}
	close(s.trigger)
}

// New creates a new scraper.
func New(c *config.Config) (*Scraper, error) {
	s := &Scraper{
		data:     &data{name: path.Join(c.DataPath, "scraper_data.json")},
		settings: &settings{name: path.Join(c.DataPath, "scraper_settings.json")},
		trigger:  make(chan bool),
	}
	if err := s.data.load(); err != nil {
		return nil, err
	}
	if err := s.settings.load(); err != nil {
		return nil, err
	}
	go s.run()
	return s, nil
}

// Messages retrieves the current list of matching messages.
func (s *Scraper) Messages() []*Message {
	s.data.Lock()
	defer s.data.Unlock()
	return s.data.Messages
}

// Get retrieves a message by its ID. It can also optionally remove the message
// from the scraper and prevent it from showing up in future scrapes.
func (s *Scraper) Get(id int, remove bool) (*Message, error) {
	s.data.Lock()
	defer s.data.Unlock()
	for i, m := range s.data.Messages {
		if m.ID == id {
			message := m
			if remove {
				s.data.Messages = append(
					s.data.Messages[:i],
					s.data.Messages[i+1:]...,
				)
				s.data.MessagesUsed = append(s.data.MessagesUsed, message.ID)
				if err := s.data.save(); err != nil {
					return nil, err
				}
			}
			return message, nil
		}
	}
	return nil, errors.New("invalid message index")
}

// Settings retrieves the current settings for the scraper.
func (s *Scraper) Settings() Settings {
	s.settings.Lock()
	defer s.settings.Unlock()
	return s.settings.Settings
}

// SetSettings stores the current settings for the scraper.
func (s *Scraper) SetSettings(settings Settings) error {
	s.settings.Lock()
	defer s.settings.Unlock()
	if !reflect.DeepEqual(settings, s.settings.Settings) {
		s.settings.Settings = settings
		s.trigger <- false
		return s.settings.save()
	}
	return nil
}

// Close shuts down the scraper and waits for it to exit.
func (s *Scraper) Close() {
	s.trigger <- true
	<-s.trigger
}
