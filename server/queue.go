package server

import (
	"github.com/AskUbuntu/tbot/util"

	"fmt"
	"net/http"
)

func (s *Server) queueHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	messages := []message{}
	defer func() {
		s.queueTemplate.Execute(w, map[string]interface{}{
			"Messages":        messages,
			"ScrapedMessages": s.scraper.Messages(),
			"QueuedMessages":  s.queue.Messages(),
		})
	}()
	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			messages = append(messages, message{"danger", err.Error()})
			return
		}
		id := util.Atoi(r.Form.Get("id"))
		if id == 0 {
			messages = append(messages, message{
				Type: "danger",
				Body: "'id' missing from form",
			})
			return
		}
		m, err := s.scraper.Use(id)
		if err != nil {
			messages = append(messages, message{"danger", err.Error()})
			return
		}
		switch r.Form.Get("action") {
		case "queue":
			s.messages <- m
			messages = append(messages, message{
				Type: "info",
				Body: fmt.Sprintf("message #%d added to queue", id),
			})
		case "delete":
			s.scraper.Blacklist(id)
			messages = append(messages, message{
				Type: "info",
				Body: fmt.Sprintf("message #%d removed from queue", id),
			})
		default:
			messages = append(messages, message{
				Type: "danger",
				Body: "invalid form action",
			})
		}
	}
}
