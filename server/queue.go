package server

import (
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
}
