package scraper

import (
	"github.com/AskUbuntu/tbot/util"
	"github.com/PuerkitoBio/goquery"

	"fmt"
	"log"
	"regexp"
	"strings"
	"time"
)

var mentionRegexp = regexp.MustCompile("@\\w+[,:]?")

func (s *Scraper) cleanBody(body string) string {
	if strings.HasPrefix(body, "//i.stack.imgur.com") {
		body = "http:" + body
	}
	body = mentionRegexp.ReplaceAllString(body, "")
	body = strings.TrimSpace(body)
	return body
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
			body    = content.Text()
			onebox  = content.Find(".onebox").Length() != 0
			stars   = util.Atoi(selection.Find(".stars .times").Text())
		)
		if onebox {
			body = content.Find(".ob-image img").AttrOr("src", "")
		}
		body = s.cleanBody(body)
		if body != "" &&
			(util.ContainsString(body, matchingWords, false) ||
				stars >= minStars) {
			var (
				signature = selection.Parent().Parent().Find(".signature")
				m         = &Message{
					ID:     id,
					URL:    fmt.Sprintf("%s%s", pollURL, link.AttrOr("href", "")),
					Body:   body,
					Onebox: onebox,
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
