package main

import (
	"errors"
	"fmt"
)

type manager interface {
	authenticate(username string, password string) (User, error)
}

type userManager struct {
	readerWriter
}

const hiddenPassword = "xxxx"

var (
	ErrInvalidUsername = errors.New("Invalid username")
	ErrInvalidPassword = errors.New("Invalid password")
)

func newUserManager(rw readerWriter) userManager {
	return userManager{rw}
}

func (u userManager) authenticate(username string, password string) (User, error) {
	rr := u.read(username)
	if rr.err == ErrNoUsersFound || rr.err == ErrMoreThanOneUserFound {
		return User{}, ErrInvalidUsername
	}

	if rr.err != nil {
		return User{}, fmt.Errorf("Error while retrieving user %s: %v", username, rr.err.Error())
	}

	if rr.user.username == username && rr.user.password == password {
		rr.user.password = hiddenPassword
		return rr.user, nil
	}

	return User{}, ErrInvalidPassword
}
