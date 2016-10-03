package server

import (
	"github.com/AskUbuntu/tbot/util"
	"github.com/flosch/pongo2"

	"net/http"
)

func (s *Server) twitterHandler(w http.ResponseWriter, r *http.Request) {
	s.render(w, r, "twitter.html", pongo2.Context{
		"tweets": s.twitter.Tweets(),
	})
}

func (s *Server) twitterCustomHandler(w http.ResponseWriter, r *http.Request) {
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
