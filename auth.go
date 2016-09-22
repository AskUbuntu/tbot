package main

import (
	"path"
)

// Auth manages access to users. This includes adding and removing users. Only
// staff can change settings and only the administrator can perform user
// actions.
type Auth struct {
	name  string
	Users map[string]*User `json:"users"`
}

// NewAuth creates a new authenticator for registered users.
func NewAuth(config *Config) (*Auth, error) {
	a := &Auth{
		name: path.Join(config.DataPath, "auth.json"),
	}
	_, err := LoadJSON(a.name, a)
	if err != nil {
		return nil, err
	}
	return a, nil
}

// Save the user list to disk.
func (a *Auth) Save() error {
	return SaveJSON(a.name, a)
}
