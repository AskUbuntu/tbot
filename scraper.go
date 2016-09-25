package main

import (
	"github.com/PuerkitoBio/goquery"

	"path"
	"time"
)

// Scraper regularly scrapes the chat transcript for the current and previous
// day for messages that match the criteria in use. Once a message matches, it
// is sent on the provided channel. The IDs of messages that matched are kept
// to ensure a message isn't matched twice.
type Scraper struct {
	settings       *Settings
	name           string
	messageOutChan chan<- *Message
	notifyChan     chan bool
	closeChan      chan bool
	LastScrape     time.Time  `json:"last_scrape"`
	EarliestID     int        `json:"earliest_id"`
	MessagesSent   []int      `json:"messages_sent"`
	Messages       []*Message `json:"messages"`
}

// Save writes the scraper state to disk.
func (s *Scraper) save() error {
	return SaveJSON(s.name, s)
}

// scrapePage scrapes an individual page for messages.
func (s *Scraper) scrapePage() {
	//...
}

// run periodically scrapes the page for matching messages.
func (s *Scraper) run() {
	//...
}

// NewScraper creates a new scraper.
func NewScraper(config *Config, settings *Settings, ch chan<- *Message) (*Scraper, error) {
	s := &Scraper{
		settings:       settings,
		name:           path.Join(config.DataPath, "scraper.json"),
		messageOutChan: ch,
		notifyChan:     make(chan bool),
		closeChan:      make(chan bool),
	}
	_, err := LoadJSON(s.name, s)
	if err != nil {
		return nil, err
	}
	settings.Register(s.notifyChan)
	go s.run()
	return s, nil
}

// Close shuts down the scraper and waits for it to exit.
func (s *Scraper) Close() {
	s.closeChan <- true
	<-s.closeChan
}
