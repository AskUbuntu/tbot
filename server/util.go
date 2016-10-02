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
	sessionUsername = "username"
	contextRequest  = "request"
	contextUsername = "username"
	contextUser     = "user"
	contextAlerts   = "alerts"
)

// setUsername sets the user for the current session to the provided user.
func (s *Server) setUsername(w http.ResponseWriter, r *http.Request, username string) {
	session, _ := s.sessions.Get(r, sessionName)
	defer session.Save(r, w)
	session.Values[sessionUsername] = username
}

const (
	infoType    = "info"
	warningType = "warning"
	dangerType  = "danger"
)

type alert struct {
	Type string
	Body string
}

// addAlert registers the provided alert for display on the next page the user
// displays.
func (s *Server) addAlert(w http.ResponseWriter, r *http.Request, alertType, body string) {
	session, _ := s.sessions.Get(r, sessionName)
	defer session.Save(r, w)
	session.AddFlash(&alert{
		Type: alertType,
		Body: body,
	})
}

// getAlerts retrieves the alerts from the current session.
func (s *Server) getAlerts(w http.ResponseWriter, r *http.Request) interface{} {
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
		var (
			username string
			user     *auth.User
		)
		session, _ := s.sessions.Get(r, sessionName)
		if v, ok := session.Values[sessionUsername]; ok {
			u, ok := v.(string)
			if ok {
				username = u
				user, _ = s.auth.Get(u)
			}
		}
		context.Set(r, contextUsername, username)
		context.Set(r, contextUser, user)
		if userType != auth.NoUser && (user == nil ||
			userType == auth.StaffUser && !user.IsStaff() ||
			userType == auth.AdminUser && !user.IsAdmin()) {
			s.addAlert(w, r, dangerType, "page requires authentication")
			http.Redirect(w, r, "/users/login", http.StatusFound)
			return
		}
		if r.Method == http.MethodPost {
			if err := r.ParseForm(); err != nil {
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}
		}
		fn(w, r)
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
	ctx[contextUsername] = context.Get(r, contextUsername).(string)
	ctx[contextUser] = context.Get(r, contextUser).(*auth.User)
	ctx[contextAlerts] = s.getAlerts(w, r)
	b, err := t.ExecuteBytes(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}
