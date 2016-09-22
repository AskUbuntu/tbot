package main

import (
	"encoding/json"
	"os"
	"path"
)

// Settings manages the storage, retrieval, and updating of application
// settings. Unlike configuration, settings can be modified without restarting
// the application. For the MinStars and MatchingWords members, messages will
// be selected if *either* one triggers.
type Settings struct {
	name           string
	channels       []chan<- bool
	PollFrequency  int      `json:"poll_frequency"`  // Wait this many minutes between polling attempts
	MinStars       int      `json:"min_stars"`       // Minimum stars for selecting messages
	MatchingWords  []string `json:"matching_words"`  // Select messages with any of these words
	QueueFrequency int      `json:"queue_frequency"` // Time (in minutes) between tweets
}

// NewSettings creates a settings manager using the specified configuration. If
// the file exists, it is loaded from disk. Otherwise, defaults are used.
func NewSettings(config *Config) (*Settings, error) {
	s := &Settings{
		name: path.Join(config.DataPath, "settings.json"),
	}
	r, err := os.Open(s.name)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		} else {
			s.PollFrequency = 10
			s.MinStars = 6
			s.MatchingWords = []string{}
			s.QueueFrequency = 180 // 3 hours
		}
	} else {
		defer r.Close()
		if err := json.NewDecoder(r).Decode(s); err != nil {
			return nil, err
		}
	}
	return s, nil
}

// Save the settings to disk.
func (s *Settings) Save() error {
	w, err := os.Create(s.name)
	if err != nil {
		return err
	}
	defer w.Close()
	if err := json.NewEncoder(w).Encode(s); err != nil {
		return err
	}
	return nil
}

// Register adds the channel to a list of ones to be notified when settings are
// changed.
func (s *Settings) Register(ch chan<- bool) {
	s.channels = append(s.channels, ch)
}

// Notify lets all registered channels know that settings have been changed.
func (s *Settings) Notify() error {
	for _, ch := range s.channels {
		go func() {
			ch <- true
		}()
	}
}

// Close closes all of the registered channels in preparation for shutdown.
func (s *Settings) Close() {
	for _, ch := range s.channels {
		close(ch)
	}
}
