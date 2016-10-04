package server

import (
	"github.com/AskUbuntu/tbot/util"
	"github.com/flosch/pongo2"

	"net/http"
)

func (s *Server) queueHandler(w http.ResponseWriter, r *http.Request) {
	s.render(w, r, "queue.html", pongo2.Context{
		"title":    "Queue",
		"messages": s.queue.Messages(),
	})
}

func (s *Server) queueDeleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		id := util.Atoi(r.Form.Get("id"))
		if err := s.queue.Delete(id); err != nil {
			s.addAlert(w, r, dangerType, err.Error())
		} else {
			s.addAlert(w, r, infoType, "message removed from queue")
		}
	}
	http.Redirect(w, r, "/queue", http.StatusFound)
}
