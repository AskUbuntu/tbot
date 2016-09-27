package auth

import (
	"github.com/AskUbuntu/tbot/config"

	"errors"
	"path"
	"time"
)

// Auth manages access to users. This includes adding and removing users. Only
// staff can change settings and only the administrator can perform user
// actions.
type Auth struct {
	data          *data
	adminUser     *User
	adminPassword string
}

// New creates a new authenticator for registered users. A special entry is
// created for the admin user.
func New(config *Config) (*Auth, error) {
	a := &Auth{
		data: &data{name: path.Join(config.DataPath, "auth_data.json")},
		adminUser: &User{
			Type: AdminUser,
		},
		adminPassword: config.AdminPassword,
	}
	if err := a.data.load(); err != nil {
		return nil, err
	}
	return a, nil
}

// Users returns a map of usernames to their account information.
func (a *Auth) Users() map[string]*User {
	a.data.Lock()
	defer a.data.Unlock()
	return a.data.Users
}

// CreateUser creates a new user. Their randomly-generated password is
// returned if the process completes without error.
func (a *Auth) CreateUser(username string) (string, error) {
	u := &User{
		Type:    StandardUser,
		Created: time.Now(),
	}
	p, err := u.resetPassword()
	if err != nil {
		return "", err
	}
	a.data.Lock()
	defer a.data.Unlock()
	a.data.Users[username] = u
	if err := a.data.save(); err != nil {
		return "", err
	}
	return p, nil
}

// Authenticate attempts to authenticate the specified using their username
// and password. To make things harder for malicious users, there is no
// distinguishing between invalid usernames and invalid passwords.
func (a *Auth) Authenticate(username, password string) (User, error) {
	if username == "admin" && password == a.adminPassword {
		return *a.adminUser, nil
	}
	a.data.Lock()
	defer a.data.Unlock()
	u, ok := a.data.Users[username]
	if ok {
		if u.authenticate(password) {
			return u, nil
		}
	}
	return nil, errors.New("invalid username or password supplied")
}
