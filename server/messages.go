package server

import (
	"github.com/AskUbuntu/tbot/util"
	"github.com/flosch/pongo2"
	"github.com/gorilla/mux"

	"net/http"
	"time"
)

func (s *Server) messagesHandler(w http.ResponseWriter, r *http.Request) {
	s.render(w, r, "messages.html", pongo2.Context{
		"title":    "Messages",
		"messages": s.scraper.Messages(),
		"next_scrape": s.scraper.LastScrape().Add(
			time.Duration(s.scraper.Settings().PollFrequency)*time.Minute,
		).Sub(time.Now()) / time.Minute,
	})
}

func (s *Server) messagesScrapeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		s.scraper.Scrape()
		s.addAlert(w, r, infoType, "scrape has begun and may take a few seconds")
	}
	http.Redirect(w, r, "/messages", http.StatusFound)
}

func (s *Server) messagesByIdEditHandler(w http.ResponseWriter, r *http.Request) {
	id := util.Atoi(mux.Vars(r)["id"])
	m, err := s.scraper.Get(id, false)
	if err != nil {
		s.addAlert(w, r, dangerType, err.Error())
		http.Redirect(w, r, "/messages", http.StatusFound)
		return
	}
	if r.Method == http.MethodPost {
		message := r.Form.Get("message")
		if message == "" {
			s.addAlert(w, r, dangerType, "message is empty")
		} else {
			m, err = s.scraper.Get(id, true)
			if err != nil {
				s.addAlert(w, r, dangerType, err.Error())
			} else {
				m.Body = message
				if err := s.queue.Add(m); err != nil {
					s.addAlert(w, r, dangerType, err.Error())
				} else {
					s.addAlert(w, r, infoType, "message edited and added to queue")
					http.Redirect(w, r, "/messages", http.StatusFound)
					return
				}
			}
		}
	}
	s.render(w, r, "messages_id_edit.html", pongo2.Context{
		"title": "Edit Message",
		"m":     m,
	})
}

func (s *Server) messagesQueueHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		id := util.Atoi(r.Form.Get("id"))
		m, err := s.scraper.Get(id, true)
		if err != nil {
			s.addAlert(w, r, dangerType, err.Error())
		} else {
			if err := s.queue.Add(m); err != nil {
				s.addAlert(w, r, dangerType, err.Error())
			} else {
				s.addAlert(w, r, infoType, "message added to queue")
			}
		}
	}
	http.Redirect(w, r, "/messages", http.StatusFound)
}

func (s *Server) messagesDeleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		id := util.Atoi(r.Form.Get("id"))
		_, err := s.scraper.Get(id, true)
		if err != nil {
			s.addAlert(w, r, dangerType, err.Error())
		} else {
			s.addAlert(w, r, infoType, "message removed")
		}
	}
	http.Redirect(w, r, "/messages", http.StatusFound)
}
