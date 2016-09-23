package main

import (
	"github.com/gorilla/mux"
	"github.com/hectane/go-asyncserver"

	"net/http"
)

// Server acts as a front end to the application.
type Server struct {
	server   *server.AsyncServer
	settings *Settings
	auth     *Auth
}

// queueHandler manages the queuing of items and custom tweets.
func (s *Server) queueHandler(w http.ResponseWriter, r *http.Request) {
}

// settingsHandler manages access to settings that control tweets.
func (s *Server) settingsHandler(w http.ResponseWriter, r *http.Request) {
}

// usersHandler manages registered users and their permissions.
func (s *Server) usersHandler(w http.ResponseWriter, r *http.Request) {
}

// NewServer creates a new server bound to the address specified in the config.
func NewServer(config *Config, settings *Settings, auth *Auth) (*Server, error) {
	var (
		r = mux.NewRouter()
		s = &Server{
			server:   server.New(config.Addr),
			settings: settings,
			auth:     auth,
		}
	)
	s.server.Handler = r
	return s, nil
}

// Close shuts down the server.
func (s *Server) Close() {
	s.server.Stop()
}
