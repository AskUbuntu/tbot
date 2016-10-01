package server

import (
	"github.com/flosch/pongo2"

	"net/http"
)

func (s *Server) usersHandler(w http.ResponseWriter, r *http.Request) {
}

func (s *Server) usersLoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			s.addMessage(w, r, dangerType, err.Error())
		} else {
			var (
				username = r.Form.Get("username")
				password = r.Form.Get("password")
			)
			u, err := s.auth.Authenticate(username, password)
			if err != nil {
				s.addMessage(w, r, dangerType, err.Error())
			} else {
				s.setUser(w, r, u)
				s.addMessage(w, r, infoType, "you have been logged in")
				http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
				return
			}
		}
	}
	s.render(w, r, "users_login.html", pongo2.Context{
		"title": "Login",
	})
}

func (s *Server) usersLogoutHandler(w http.ResponseWriter, r *http.Request) {
	s.deleteUser(w, r)
	s.addMessage(w, r, infoType, "you have been logged out")
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func (s *Server) usersResetHandler(w http.ResponseWriter, r *http.Request) {
}

func (s *Server) usersCreateHandler(w http.ResponseWriter, r *http.Request) {
}

func (s *Server) usersDeleteHandler(w http.ResponseWriter, r *http.Request) {
}
