package main

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
)

type manager interface {
	authenticate(username string, password string) (user, error)
}

type userManager struct {
	readerWriter
}

const hiddenPassword = "xxxx"

var (
	//ErrInvalidUsername is a generic error for inavild username
	ErrInvalidUsername = errors.New("Invalid username")
	//ErrInvalidPassword is a generic error for inavild passwrod
	ErrInvalidPassword = errors.New("Invalid password")
)

func newUserManager(rw readerWriter) userManager {
	return userManager{rw}
}

func (u userManager) authenticate(username string, password string) (user, error) {
	rr := u.read(username)
	if rr.err == ErrNoUsersFound || rr.err == ErrMoreThanOneUserFound {
		return user{}, ErrInvalidUsername
	}

	if rr.err != nil {
		return user{}, fmt.Errorf("Error while retrieving user %s: %v", username, rr.err.Error())
	}

	if rr.user.username == username && rr.user.password == password {
		rr.user.password = hiddenPassword
		return rr.user, nil
	}

	return user{}, ErrInvalidPassword
}

func sessionID() string {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}
