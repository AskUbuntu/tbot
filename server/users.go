package server

import (
	"fmt"
	"net/http"
)

// usersHandler manages registered users and their permissions.
func (s *Server) usersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	messages := []message{}
	defer func() {
		s.usersTemplate.Execute(w, map[string]interface{}{
			"Messages": messages,
			"Users":    s.auth.Users(),
		})
	}()
	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			messages = append(messages, message{"danger", err.Error()})
			return
		}
		username := r.Form.Get("username")
		if username == "" {
			messages = append(messages, message{
				Type: "info",
				Body: "'username' missing from form",
			})
			return
		}
		p, err := s.auth.CreateUser(username)
		if err != nil {
			messages = append(messages, message{"danger", err.Error()})
			return
		}
		messages = append(messages, message{
			Type: "info",
			Body: fmt.Sprintf(
				"user %s created with password '%s'",
				username,
				p,
			),
		})
	}
}
