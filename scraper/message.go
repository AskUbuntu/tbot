package scraper

import (
	"fmt"
	"strings"
	"time"
)

// Message is an individual message found by the scraper.
type Message struct {
	ID          int       `json:"id"`
	URL         string    `json:"url"`
	Body        string    `json:"body"`
	Onebox      bool      `json:"onebox"`
	Stars       int       `json:"stars"`
	Author      string    `json:"author"`
	AuthorImage string    `json:"author_image"`
	Created     time.Time `json:"created"`
}

// String converts the message into a properly formatted string ready for
// tweeting. Tweets that exceed the 140-character limit are truncated.
func (m *Message) String() string {
	var (
		charsRemaining = 140
		body           = m.Body
	)
	switch {
	case m.Onebox:
	case len(body) > charsRemaining:
		charsRemaining -= 2
		body = strings.TrimSpace(body[:charsRemaining-1])
		body += "…"
		fallthrough
	default:
		body = fmt.Sprintf("“%s”", body)
	}
	return body
}
