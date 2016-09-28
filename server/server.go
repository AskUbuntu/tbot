package server

import (
	"github.com/AskUbuntu/tbot/auth"
	"github.com/AskUbuntu/tbot/config"
	"github.com/AskUbuntu/tbot/queue"
	"github.com/AskUbuntu/tbot/scraper"
	"github.com/AskUbuntu/tbot/twitter"
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
	messages         chan *scraper.Message
	auth             *auth.Auth
	queue            *queue.Queue
	scraper          *scraper.Scraper
	twitter          *twitter.Twitter
	queueTemplate    *template.Template
	settingsTemplate *template.Template
	usersTemplate    *template.Template
}

type message struct {
	Type string
	Body string
}

// New creates a new server bound to the address specified in the config. The
// server also acts as the central coordination point for the other objects,
// such as the scraper and twitter client.
func New(config *config.Config) (*Server, error) {
	var (
		messagesIn  = make(chan *scraper.Message)
		messagesOut = make(chan *scraper.Message)
	)
	a, err := auth.New(config)
	if err != nil {
		return nil, err
	}
	q, err := queue.New(config, messagesIn, messagesOut)
	if err != nil {
		return nil, err
	}
	s, err := scraper.New(config)
	if err != nil {
		return nil, err
	}
	t, err := twitter.New(config, messagesOut)
	if err != nil {
		return nil, err
	}
	srv := &Server{
		server:   server.New(config.Addr),
		mux:      mux.NewRouter(),
		messages: messagesIn,
		auth:     a,
		queue:    q,
		scraper:  s,
		twitter:  t,
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
	srv.server.Handler = srv
	srv.mux.HandleFunc("/", srv.indexHandler)
	srv.mux.HandleFunc("/queue", srv.queueHandler)
	srv.mux.HandleFunc("/settings", srv.settingsHandler)
	srv.mux.HandleFunc("/users", srv.usersHandler)
	srv.mux.PathPrefix("/").Handler(http.FileServer(http.Dir(config.RootPath)))
	if err := srv.server.Start(); err != nil {
		return nil, err
	}
	return srv, nil
}

// ServeHTTP ensures that a valid username and password are provided before
// passing the request along to the mux. The administrator's credentials are
// always accepted.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	username, password, _ := r.BasicAuth()
	if username != "" {
		if _, err := s.auth.Authenticate(username, password); err == nil {
			s.mux.ServeHTTP(w, r)
			return
		}
	}
	w.Header().Set("WWW-Authenticate", "Basic realm=tbot")
	http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
}

// Close shuts down the server.
func (s *Server) Close() {
	s.server.Stop()
	s.queue.Close()
	close(s.messages)
	s.twitter.Close()
	s.scraper.Close()
}
