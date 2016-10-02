package server

import (
	"github.com/AskUbuntu/tbot/auth"
	"github.com/flosch/pongo2"
	"github.com/gorilla/context"

	"fmt"
	"net/http"
)

func (s *Server) usersHandler(w http.ResponseWriter, r *http.Request) {
	s.render(w, r, "users.html", pongo2.Context{
		"title": "Users",
		"users": s.auth.Users(),
	})
}

func (s *Server) usersLoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var (
			username = r.Form.Get("username")
			password = r.Form.Get("password")
		)
		u, err := s.auth.Authenticate(username, password)
		if err != nil {
			s.addAlert(w, r, dangerType, err.Error())
		} else {
			s.setUsername(w, r, username)
			if u.ChangePassword {
				s.addAlert(w, r, infoType, "password reset automatically triggered")
				http.Redirect(w, r, "/users/password", http.StatusFound)
			} else {
				s.addAlert(
					w, r, infoType,
					fmt.Sprintf("welcome back, %s", username),
				)
				http.Redirect(w, r, "/messages", http.StatusFound)
			}
			return
		}
	}
	s.render(w, r, "users_login.html", pongo2.Context{
		"title": "Login",
	})
}

func (s *Server) usersPasswordHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var (
			password  = r.Form.Get("password")
			password2 = r.Form.Get("password2")
		)
		if password == "" || password != password2 {
			s.addAlert(w, r, dangerType, "both passwords must be set to identical values")
		} else {
			username := context.Get(r, contextUsername).(string)
			if err := s.auth.SetPassword(username, password); err != nil {
				s.addAlert(w, r, dangerType, err.Error())
			} else {
				s.addAlert(w, r, infoType, "password successfully changed")
				http.Redirect(w, r, "/messages", http.StatusFound)
				return
			}
		}
	}
	s.render(w, r, "users_password.html", pongo2.Context{
		"title": "Set Password",
	})
}

func (s *Server) usersLogoutHandler(w http.ResponseWriter, r *http.Request) {
	s.setUsername(w, r, "")
	s.addAlert(w, r, infoType, "you have been logged out")
	http.Redirect(w, r, "/", http.StatusFound)
}

func (s *Server) usersResetHandler(w http.ResponseWriter, r *http.Request) {
}

func (s *Server) usersDeleteHandler(w http.ResponseWriter, r *http.Request) {
}

func (s *Server) usersCreateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var (
			username = r.Form.Get("username")
			userType = r.Form.Get("type")
		)
		if username == "" || (userType != auth.StandardUser &&
			userType != auth.StaffUser &&
			userType != auth.AdminUser) {
			s.addAlert(w, r, dangerType, "invalid username or type")
		} else {
			// TODO: no checks for existing users
			p, err := s.auth.CreateUser(username, userType)
			if err != nil {
				s.addAlert(w, r, dangerType, err.Error())
			} else {
				s.addAlert(
					w, r, infoType,
					fmt.Sprintf("user created with password '%s'", p),
				)
				http.Redirect(w, r, "/users", http.StatusFound)
				return
			}
		}
	}
	s.render(w, r, "users_create.html", pongo2.Context{
		"title": "Create User",
	})
}
