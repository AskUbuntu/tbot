package server

import (
	"github.com/AskUbuntu/tbot/util"
	"github.com/flosch/pongo2"

	"net/http"
)

func (s *Server) messagesHandler(w http.ResponseWriter, r *http.Request) {
	s.render(w, r, "messages.html", pongo2.Context{
		"title":    "Messages",
		"messages": s.scraper.Messages(),
	})
}

func (s *Server) messagesQueueHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		id := util.Atoi(r.Form.Get("id"))
		m, err := s.scraper.Get(id)
		if err != nil {
			s.addAlert(w, r, dangerType, err.Error())
		} else {
			s.messages <- m
			s.addAlert(w, r, infoType, "message added to queue")
		}
	}
	http.Redirect(w, r, "/messages", http.StatusFound)
}

func (s *Server) messagesDeleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		id := util.Atoi(r.Form.Get("id"))
		_, err := s.scraper.Get(id)
		if err != nil {
			s.addAlert(w, r, dangerType, err.Error())
		} else {
			s.addAlert(w, r, infoType, "message removed")
		}
	}
	http.Redirect(w, r, "/messages", http.StatusFound)
}
