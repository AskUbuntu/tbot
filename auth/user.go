package auth

import (
	"golang.org/x/crypto/bcrypt"

	"crypto/rand"
	"encoding/base64"
	"time"
)

const (
	NoUser       = ""
	StandardUser = "standard"
	StaffUser    = "staff"
	AdminUser    = "admin"
)

// User represents information for a registered user.
type User struct {
	PasswordHash   []byte    `json:"password_hash"`
	ChangePassword bool      `json:"change_password"`
	Type           string    `json:"type"`
	Created        time.Time `json:"created"`
}

// Authenticate will check the specified password against its stored hash.
func (u *User) authenticate(password string) bool {
	if bcrypt.CompareHashAndPassword(u.PasswordHash, []byte(password)) == nil {
		return true
	}
	return false
}

// resetPassword generates a password for the user and forces it to be changed
// immediately after login.
func (u *User) resetPassword() (string, error) {
	b := make([]byte, 10)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	password := base64.StdEncoding.EncodeToString(b)
	if err := u.setPassword(password); err != nil {
		return "", err
	}
	u.ChangePassword = true
	return password, nil
}

// setPassword changes the password set on the account.
func (u *User) setPassword(password string) error {
	h, err := bcrypt.GenerateFromPassword([]byte(password), 0)
	if err != nil {
		return err
	}
	u.PasswordHash = h
	return nil
}

// IsStaff determines if the user is a staff member.
func (u *User) IsStaff() bool {
	return u.Type == StaffUser || u.Type == AdminUser
}

// IsAdmin determines if the user is an admin.
func (u *User) IsAdmin() bool {
	return u.Type == AdminUser
}
