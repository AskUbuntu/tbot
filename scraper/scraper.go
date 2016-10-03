package scraper

import (
	"github.com/AskUbuntu/tbot/config"
	"github.com/AskUbuntu/tbot/util"
	"github.com/PuerkitoBio/goquery"

	"errors"
	"fmt"
	"log"
	"path"
	"reflect"
	"strings"
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

func (s *Scraper) scrapePage(document *goquery.Document) (earliestID int, messages []*Message) {
	s.data.Lock()
	messagesUsed := s.data.MessagesUsed
	s.data.Unlock()
	s.settings.Lock()
	var (
		pollURL       = s.settings.PollURL
		minStars      = s.settings.MinStars
		matchingWords = s.settings.MatchingWords
	)
	s.settings.Unlock()
	d, _ := time.Parse("2006-01-02", document.Find("#info .icon").AttrOr("title", ""))
	document.Find(".message").Each(func(i int, selection *goquery.Selection) {
		var (
			link = selection.Find("a[name]")
			id   = util.Atoi(link.AttrOr("name", ""))
		)
		if id == 0 || util.ContainsInt(messagesUsed, id) {
			return
		}
		if earliestID == 0 {
			earliestID = id
		}
		var (
			content = selection.Find(".content")
			body    = strings.TrimSpace(content.Text())
			onebox  = content.Find(".onebox").Length() != 0
			stars   = util.Atoi(selection.Find(".stars .times").Text())
		)
		if onebox {
			body = content.Find(".ob-image img").AttrOr("src", "")
		}
		if body != "" &&
			(util.ContainsString(body, matchingWords, false) ||
				stars >= minStars) {
			var (
				signature = selection.Parent().Parent().Find(".signature")
				m         = &Message{
					ID:     id,
					URL:    fmt.Sprintf("%s%s", pollURL, link.AttrOr("href", "")),
					Body:   body,
					Stars:  stars,
					Author: signature.Find(".username").Text(),
					AuthorImage: strings.Replace(
						signature.Find(".avatar img").AttrOr("src", ""),
						"?s=16", "?s=48", -1,
					),
					Created: d,
				}
			)
			messages = append(messages, m)
		}
	})
	return
}

func (s *Scraper) scrape() error {
	s.settings.Lock()
	var (
		pollURL     = s.settings.PollURL
		pollRoomID  = s.settings.PollRoomID
		pollHistory = s.settings.PollHistory
	)
	s.settings.Unlock()
	var (
		start = time.Now().Add(time.Duration(pollHistory) * -24 * time.Hour)
		path  = fmt.Sprintf(
			"/transcript/%d/%d/%d/%d",
			pollRoomID,
			start.Year(),
			start.Month(),
			start.Day(),
		)
		earliestID = 0
		messages   = []*Message{}
	)
	for path != "" {
		log.Printf("Scraping '%s'...\n", path)
		document, err := goquery.NewDocument(
			fmt.Sprintf("%s%s", pollURL, path),
		)
		if err != nil {
			return err
		}
		newEarliestID, newMessages := s.scrapePage(document)
		if earliestID == 0 {
			earliestID = newEarliestID
		}
		messages = append(messages, newMessages...)
		selection := document.Find(".pager .current").Next()
		if selection.Length() == 0 {
			selection = document.Find("a[rel=prev]").NextAllFiltered("a")
		}
		path = selection.AttrOr("href", "")
	}
	s.data.Lock()
	s.data.LastScrape = time.Now()
	s.data.EarliestID = earliestID
	s.data.Messages = messages
	s.data.MessagesUsed = util.FilterInt(s.data.MessagesUsed, func(i int) bool {
		return i >= earliestID
	})
	if err := s.data.save(); err != nil {
		return err
	}
	s.data.Unlock()
	return nil
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

// Get removes the message from the list in preparation for use. This will also
// cause the message to be ignored in future scrapes.
func (s *Scraper) Get(id int) (*Message, error) {
	s.data.Lock()
	defer s.data.Unlock()
	var message *Message
	for i := len(s.data.Messages) - 1; i >= 0; i-- {
		m := s.data.Messages[i]
		if m.ID == id {
			message = m
			s.data.Messages = append(
				s.data.Messages[:i],
				s.data.Messages[i+1:]...,
			)
		}
	}
	if message == nil {
		return nil, errors.New("Invalid message index")
	}
	s.data.MessagesUsed = append(s.data.MessagesUsed, message.ID)
	s.data.save()
	return message, nil
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
