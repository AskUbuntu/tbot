package server

import (
	"net/http"
)

// indexHandler redirects the client to the queue tab.
func (s *Server) indexHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/queue", http.StatusTemporaryRedirect)
}
