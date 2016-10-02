package scraper

import (
	"fmt"
	"strings"
)

// Message is an individual message found by the scraper.
type Message struct {
	ID     int    `json:"id"`
	URL    string `json:"url"`
	Body   string `json:"body"`
	Author string `json:"author"`
	Stars  int    `json:"stars"`
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
