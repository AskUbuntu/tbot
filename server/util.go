package server

import (
	"github.com/AskUbuntu/tbot/auth"
	"github.com/flosch/pongo2"
	"github.com/gorilla/context"

	"net/http"
	"path"
)

const (
	sessionName     = "auth"
	sessionUser     = "user"
	contextUser     = "user"
	contextMessages = "messages"
)

// initRequest initializes the request with context variables, such as the
// current user and flashes.
func (s *Server) initRequest(r *http.Request) {
	var user *auth.User
	session, _ := s.sessions.Get(r, sessionName)
	if v, ok := session.Values[sessionUser]; ok {
		if u, ok := v.(*auth.User); ok {
			user = u
		}
	}
	context.Set(r, contextUser, user)
	context.Set(r, contextMessages, session.Flashes())
}

// getUser retrieves the user for the request.
func (s *Server) getUser(r *http.Request) *auth.User {
	return context.Get(r, contextUser).(*auth.User)
}

// setUser sets the user for the current session to the provided user.
func (s *Server) setUser(w http.ResponseWriter, r *http.Request, u *auth.User) {
	session, _ := s.sessions.Get(r, sessionName)
	session.Values[sessionUser] = u
	session.Save(r, w)
}

const (
	infoType    = "info"
	warningType = "warning"
	dangerType  = "danger"
)

// addMessage registers the provided message for display on the next page the
// user displays.
func (s *Server) addMessage(w http.ResponseWriter, r *http.Request, flashType, body string) {
	session, _ := s.sessions.Get(r, sessionName)
	session.AddFlash(map[string]string{
		"type": flashType,
		"body": body,
	})
	session.Save(r, w)
}

// r prevents users from accessing pages for which they do not have the correct
// permissions. The first argument is the minimum required permission for
// accessing the page. The second argument is the handler which will be invoked
// upon confirmation of authorization.
func (s *Server) r(userType string, fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if u := s.getUser(r); u != nil {
			if userType == auth.StandardUser ||
				userType == auth.StaffUser && u.Type != auth.StandardUser ||
				userType == auth.AdminUser && u.Type == auth.AdminUser {
				fn(w, r)
				return
			}
		}
		s.addMessage(w, r, dangerType, "page requires authentication")
		http.Redirect(w, r, "/users/login", http.StatusTemporaryRedirect)
	}
}

// render loads the specified template, injects the provided context into it,
// and renders it.
func (s *Server) render(w http.ResponseWriter, r *http.Request, templateName string, ctx pongo2.Context) {
	t, err := pongo2.FromFile(path.Join(s.templatePath, templateName))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ctx[contextUser] = context.Get(r, contextUser)
	ctx[contextMessages] = context.Get(r, contextMessages)
	b, err := t.ExecuteBytes(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(b)
}
