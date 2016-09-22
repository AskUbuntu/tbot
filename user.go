package main

import (
	"golang.org/x/crypto/bcrypt"

	"crypto/rand"
	"encoding/base64"
	"path"
	"time"
)

// User represents information for a registered user.
type User struct {
	PasswordHash   []byte    `json:"password_hash"`
	ChangePassword bool      `json:"change_password"`
	IsStaff        bool      `json:"is_staff"`
	Created        time.Time `json:"created"`
}

// ResetPassword generates a password for the user and forces it to be changed
// immediately after login.
func (u *User) ResetPassword() (string, error) {
	b := make([]byte, 6)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	if err := u.SetPassword(base64.StdEncoding.EncodeToString(b)); err != nil {
		return "", err
	}
	return
}

// SetPassword changes the password set on the account.
func (u *User) SetPassword(password string) error {
	h, err = bcrypt.GenerateFromPassword([]byte(password), 0)
	if err != nil {
		return err
	}
	u.PasswordHash = h
	return nil
}

// Authenticate will check the specified password against its stored hash.
func (u *User) Authenticate(password string) bool {
	if bcrypt.CompareHashAndPassword(u.PasswordHash, []byte(password)) == nil {
		return true
	}
	return false
}
