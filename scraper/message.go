package scraper

import (
	"fmt"
	"time"

	"github.com/AskUbuntu/tbot/util"
	"github.com/flosch/pongo2"
)

const (
	OneboxImage = "image"
)

// Message is an individual message found by the scraper.
type Message struct {
	ID          int       `json:"id"`
	URL         string    `json:"url"`
	Body        string    `json:"body"`
	Onebox      string    `json:"onebox"`
	Stars       int       `json:"stars"`
	Author      string    `json:"author"`
	AuthorImage string    `json:"author_image"`
	Created     time.Time `json:"created"`
}

// String returns the value that should be passed to Twitter to update the
// status.
func (m *Message) String() string {
	switch m.Onebox {
	case OneboxImage:
		return ""
	default:
		return fmt.Sprintf("“%s”", util.Truncate(m.Body, 138))
	}
}

// HTML returns an object that can be rendered in a template, either as a
// single string or as pre-escaped HTML.
func (m *Message) HTML() interface{} {
	switch m.Onebox {
	case OneboxImage:
		return pongo2.AsSafeValue(
			fmt.Sprintf("<img src=\"%s\" class=\"img-responsive\">", m.Body),
		)
	default:
		return fmt.Sprintf("“%s”", util.Truncate(m.Body, 138))
	}
}
