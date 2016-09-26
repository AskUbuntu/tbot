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
	// Tweet format:
	//
	//     "[body]" - [author]\n\n[url]
	//
	// We start with 140 characters and then subtract:
	//  - 7 for punctuation
	//  - 23 for the URL (per current shortened URL length)
	//
	// If the author and body fit, use them as-is. Otherwise truncate the
	// author to 12 characters and/or the body to whatever length is needed
	var (
		charsRemaining = 110
		body           = m.Body
		author         = m.Author
	)
	if len(body)+len(author) > charsRemaining {
		if len(author) > 12 {
			author = strings.TrimSpace(author[:9])
			author += "…"
		}
		charsRemaining -= len(author)
		if len(body) > charsRemaining {
			body = strings.TrimSpace(body[:charsRemaining-1])
			body += "…"
		}
	}
	return fmt.Sprintf("“%s” — %s\n\n%s", body, author, m.URL)
}
