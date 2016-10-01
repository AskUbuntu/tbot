package server

import (
	"github.com/AskUbuntu/tbot/auth"
	"github.com/flosch/pongo2"
	"github.com/gorilla/context"

	"net/http"
	"path"
)

const (
	sessionName     = "session"
	sessionUser     = "user"
	contextRequest  = "request"
	contextUser     = "user"
	contextMessages = "messages"
)

// initRequest initializes the request with context variables, such as the
// current user.
func (s *Server) initRequest(r *http.Request) {
	var user *auth.User
	session, _ := s.sessions.Get(r, sessionName)
	if v, ok := session.Values[sessionUser]; ok {
		if u, ok := v.(*auth.User); ok {
			user = u
		}
	}
	context.Set(r, contextUser, user)
}

// getUser retrieves the user for the request.
func (s *Server) getUser(r *http.Request) *auth.User {
	return context.Get(r, contextUser).(*auth.User)
}

// setUser sets the user for the current session to the provided user.
func (s *Server) setUser(w http.ResponseWriter, r *http.Request, u *auth.User) {
	session, _ := s.sessions.Get(r, sessionName)
	defer session.Save(r, w)
	session.Values[sessionUser] = u
}

// deleteUser removes the user from the current session.
func (s *Server) deleteUser(w http.ResponseWriter, r *http.Request) {
	session, _ := s.sessions.Get(r, sessionName)
	defer session.Save(r, w)
	delete(session.Values, sessionUser)
}

const (
	infoType    = "info"
	warningType = "warning"
	dangerType  = "danger"
)

type message struct {
	Type string
	Body string
}

// addMessage registers the provided message for display on the next page the
// user displays.
func (s *Server) addMessage(w http.ResponseWriter, r *http.Request, flashType, body string) {
	session, _ := s.sessions.Get(r, sessionName)
	defer session.Save(r, w)
	session.AddFlash(&message{
		Type: flashType,
		Body: body,
	})
}

// getMessages retrieves the messages from the current session.
func (s *Server) getMessages(w http.ResponseWriter, r *http.Request) interface{} {
	session, _ := s.sessions.Get(r, sessionName)
	defer session.Save(r, w)
	return session.Flashes()
}

// r prevents users from accessing pages for which they do not have the correct
// permissions. The first argument is the handler which will be invoked upon
// confirmation of authorization. The second argument is the minimum required
// permission for accessing the page.
func (s *Server) r(userType string, fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.initRequest(r)
		if u := s.getUser(r); userType == auth.NoUser ||
			(u != nil && userType == auth.StandardUser ||
				userType == auth.StaffUser && u.IsStaff() ||
				userType == auth.AdminUser && u.IsAdmin()) {
			fn(w, r)
			return
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
	ctx[contextRequest] = r
	ctx[contextUser] = s.getUser(r)
	ctx[contextMessages] = s.getMessages(w, r)
	b, err := t.ExecuteBytes(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}
