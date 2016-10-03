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
	Stars       int       `json:"stars"`
	Author      string    `json:"author"`
	AuthorImage string    `json:"author_image"`
	Created     time.Time `json:"created"`
}

// String converts the message into a properly formatted string ready for
// tweeting. Tweets that exceed the 140-character limit are truncated.
func (m *Message) String() string {
	var (
		charsRemaining = 138
		body           = m.Body
	)
	if len(body) > charsRemaining {
		body = strings.TrimSpace(body[:charsRemaining-1])
		body += "…"
	}
	return fmt.Sprintf("“%s”", body)
}
