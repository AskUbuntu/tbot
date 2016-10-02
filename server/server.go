package server

import (
	"github.com/AskUbuntu/tbot/auth"
	"github.com/AskUbuntu/tbot/config"
	"github.com/AskUbuntu/tbot/queue"
	"github.com/AskUbuntu/tbot/scraper"
	"github.com/AskUbuntu/tbot/twitter"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/hectane/go-asyncserver"

	"encoding/gob"
	"net/http"
	"path"
)

// Server acts as a front end to the application, allowing the entire
// application to be controlled directly from the web.
type Server struct {
	server       *server.AsyncServer
	mux          *mux.Router
	sessions     *sessions.CookieStore
	templatePath string
	messages     chan *scraper.Message
	auth         *auth.Auth
	queue        *queue.Queue
	twitter      *twitter.Twitter
	scraper      *scraper.Scraper
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
	t, err := twitter.New(config, messagesOut)
	if err != nil {
		return nil, err
	}
	s, err := scraper.New(config)
	if err != nil {
		return nil, err
	}
	srv := &Server{
		server:       server.New(config.Addr),
		mux:          mux.NewRouter(),
		sessions:     sessions.NewCookieStore([]byte(config.SecretKey)),
		templatePath: path.Join(config.RootPath, "templates"),
		messages:     messagesIn,
		auth:         a,
		queue:        q,
		twitter:      t,
		scraper:      s,
	}
	srv.server.Handler = srv
	srv.mux.HandleFunc("/", srv.r(auth.NoUser, srv.indexHandler))
	srv.mux.HandleFunc("/users", srv.r(auth.AdminUser, srv.usersHandler))
	srv.mux.HandleFunc("/users/login", srv.r(auth.NoUser, srv.usersLoginHandler))
	srv.mux.HandleFunc("/users/password", srv.r(auth.StandardUser, srv.usersPasswordHandler))
	srv.mux.HandleFunc("/users/logout", srv.r(auth.StandardUser, srv.usersLogoutHandler))
	srv.mux.HandleFunc("/users/reset", srv.r(auth.AdminUser, srv.usersResetHandler))
	srv.mux.HandleFunc("/users/create", srv.r(auth.AdminUser, srv.usersCreateHandler))
	srv.mux.HandleFunc("/users/delete", srv.r(auth.AdminUser, srv.usersDeleteHandler))
	srv.mux.HandleFunc("/messages", srv.r(auth.StandardUser, srv.messagesHandler))
	srv.mux.HandleFunc("/messages/queue", srv.r(auth.StandardUser, srv.messagesQueueHandler))
	srv.mux.HandleFunc("/messages/delete", srv.r(auth.StandardUser, srv.messagesDeleteHandler))
	srv.mux.HandleFunc("/queue", srv.r(auth.StandardUser, srv.queueHandler))
	srv.mux.HandleFunc("/queue/delete", srv.r(auth.StandardUser, srv.queueDeleteHandler))
	srv.mux.HandleFunc("/twitter", srv.r(auth.StandardUser, srv.twitterHandler))
	srv.mux.HandleFunc("/twitter/custom", srv.r(auth.StandardUser, srv.twitterCustomHandler))
	srv.mux.HandleFunc("/twitter/delete", srv.r(auth.StandardUser, srv.twitterDeleteHandler))
	srv.mux.HandleFunc("/settings", srv.r(auth.StaffUser, srv.settingsHandler))
	srv.mux.PathPrefix("/static").Handler(
		http.FileServer(http.Dir(config.RootPath)),
	)
	gob.Register(&auth.User{})
	gob.Register(&alert{})
	if err := srv.server.Start(); err != nil {
		return nil, err
	}
	return srv, nil
}

// ServeHTTP initializes the request (with context variables) and dispatches it
// to the appropriate handler.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

// Close shuts down the server.
func (s *Server) Close() {
	s.server.Stop()
	s.queue.Close()
	close(s.messages)
	s.twitter.Close()
	s.scraper.Close()
}
