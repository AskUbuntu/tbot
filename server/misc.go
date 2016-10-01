package server

import (
	"github.com/flosch/pongo2"

	"net/http"
)

func (s *Server) indexHandler(w http.ResponseWriter, r *http.Request) {
	s.render(w, r, "index.html", pongo2.Context{})
}

func (s *Server) settingsHandler(w http.ResponseWriter, r *http.Request) {
}
