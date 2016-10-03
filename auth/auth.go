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
	data      *data
	adminUser *User
}

func (a *Auth) get(username string) (*User, error) {
	if username == "admin" {
		return a.adminUser, nil
	}
	u, ok := a.data.Users[username]
	if ok {
		return u, nil
	}
	return nil, errors.New("user does not exist")
}

// New creates a new authenticator for registered users. A special entry is
// created for the admin user.
func New(config *config.Config) (*Auth, error) {
	a := &Auth{
		data: &data{
			name:  path.Join(config.DataPath, "auth_data.json"),
			Users: make(map[string]*User),
		},
		adminUser: &User{
			Type: AdminUser,
		},
	}
	if err := a.adminUser.setPassword(config.AdminPassword); err != nil {
		return nil, err
	}
	if err := a.data.load(); err != nil {
		return nil, err
	}
	return a, nil
}

// Get retrieves a specific user.
func (a *Auth) Get(username string) (*User, error) {
	a.data.Lock()
	defer a.data.Unlock()
	return a.get(username)
}

// Users returns a map of usernames to their account information.
func (a *Auth) Users() map[string]*User {
	a.data.Lock()
	defer a.data.Unlock()
	return a.data.Users
}

// CreateUser creates a new user. Their randomly-generated password is
// returned if the process completes without error.
func (a *Auth) CreateUser(username, userType string) (string, error) {
	u := &User{
		Type:    userType,
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
func (a *Auth) Authenticate(username, password string) (*User, error) {
	a.data.Lock()
	defer a.data.Unlock()
	u, err := a.get(username)
	if err == nil {
		if u.authenticate(password) {
			return u, nil
		}
	}
	return nil, errors.New("invalid username or password supplied")
}

// SetPassword attempts to set a new password for a user.
func (a *Auth) SetPassword(username, password string) error {
	a.data.Lock()
	defer a.data.Unlock()
	u, err := a.get(username)
	if err != nil {
		return err
	}
	if err := u.setPassword(password); err != nil {
		return err
	}
	if err := a.data.save(); err != nil {
		return err
	}
	return nil
}

// ResetPassword attempts to reset a user's password.
func (a *Auth) ResetPassword(username string) (string, error) {
	a.data.Lock()
	defer a.data.Unlock()
	u, err := a.get(username)
	if err != nil {
		return "", err
	}
	p, err := u.resetPassword()
	if err != nil {
		return "", err
	}
	if err := a.data.save(); err != nil {
		return "", err
	}
	return p, nil
}

// Delete removes the specified user account.
func (a *Auth) Delete(username string) error {
	a.data.Lock()
	defer a.data.Unlock()
	if _, err := a.get(username); err != nil {
		return err
	}
	delete(a.data.Users, username)
	return a.data.save()
}
