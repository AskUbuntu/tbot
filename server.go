package main

import (
	"github.com/gorilla/mux"
	"github.com/hectane/go-asyncserver"

	"html/template"
	"net/http"
	"path"
)

// Server acts as a front end to the application.
type Server struct {
	server           *server.AsyncServer
	mux              *mux.Router
	settings         *Settings
	auth             *Auth
	queue            *Queue
	queueTemplate    *template.Template
	settingsTemplate *template.Template
	usersTemplate    *template.Template
}

// indexHandler redirects the client to the queue tab.
func (s *Server) indexHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/queue", http.StatusTemporaryRedirect)
}

// queueHandler manages the queuing of items and custom tweets.
func (s *Server) queueHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	s.queueTemplate.Execute(w, map[string]interface{}{
		"messages": []string{},
		"tweets":   s.queue.QueuedTweets,
	})
}

// settingsHandler manages access to settings that control tweets.
func (s *Server) settingsHandler(w http.ResponseWriter, r *http.Request) {
}

// usersHandler manages registered users and their permissions.
func (s *Server) usersHandler(w http.ResponseWriter, r *http.Request) {
}

// NewServer creates a new server bound to the address specified in the config.
func NewServer(config *Config, settings *Settings, auth *Auth, queue *Queue) (*Server, error) {
	s := &Server{
		server:   server.New(config.Addr),
		mux:      mux.NewRouter(),
		settings: settings,
		auth:     auth,
		queue:    queue,
		queueTemplate: template.Must(template.ParseFiles(
			path.Join(config.RootPath, "base.tmpl"),
			path.Join(config.RootPath, "queue.tmpl"),
		)),
		settingsTemplate: template.Must(template.ParseFiles(
			path.Join(config.RootPath, "base.tmpl"),
			path.Join(config.RootPath, "settings.tmpl"),
		)),
		usersTemplate: template.Must(template.ParseFiles(
			path.Join(config.RootPath, "base.tmpl"),
			path.Join(config.RootPath, "users.tmpl"),
		)),
	}
	s.server.Handler = s
	s.mux.HandleFunc("/", s.indexHandler)
	s.mux.HandleFunc("/queue", s.queueHandler)
	s.mux.HandleFunc("/settings", s.settingsHandler)
	s.mux.HandleFunc("/users", s.usersHandler)
	s.mux.PathPrefix("/").Handler(http.FileServer(http.Dir(config.RootPath)))
	if err := s.server.Start(); err != nil {
		return nil, err
	}
	return s, nil
}

// ServeHTTP ensures that a valid username and password are provided before
// passing the request along to the mux. The administrator's credentials are
// always accepted.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	username, password, _ := r.BasicAuth()
	if username != "" {
		if user, ok := s.auth.Users[username]; ok {
			if user.Authenticate(password) {
				s.mux.ServeHTTP(w, r)
				return
			}
		}
	}
	w.Header().Set("WWW-Authenticate", "Basic realm=tbot")
	http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
}

// Close shuts down the server.
func (s *Server) Close() {
	s.server.Stop()
}
