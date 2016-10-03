package server

import (
	"github.com/AskUbuntu/tbot/util"
	"github.com/flosch/pongo2"

	"net/http"
)

func (s *Server) twitterHandler(w http.ResponseWriter, r *http.Request) {
	s.render(w, r, "twitter.html", pongo2.Context{
		"title":  "Twitter",
		"tweets": s.twitter.Tweets(),
	})
}

func (s *Server) twitterCustomHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		message := r.Form.Get("message")
		if message == "" {
			s.addAlert(w, r, dangerType, "message is empty")
		} else {
			if err := s.twitter.Send(message); err != nil {
				s.addAlert(w, r, dangerType, err.Error())
			} else {
				s.addAlert(w, r, infoType, "tweet sent successfully")
				http.Redirect(w, r, "/twitter", http.StatusFound)
				return
			}
		}
	}
	s.render(w, r, "twitter_custom.html", pongo2.Context{
		"title":  "Custom Tweet",
		"tweets": s.twitter.Tweets(),
	})
}

func (s *Server) twitterDeleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		id := util.Atoi64(r.Form.Get("id"))
		if err := s.twitter.Delete(id); err != nil {
			s.addAlert(w, r, dangerType, err.Error())
		} else {
			s.addAlert(w, r, infoType, "tweet deleted")
		}
	}
	http.Redirect(w, r, "/twitter", http.StatusFound)
}
