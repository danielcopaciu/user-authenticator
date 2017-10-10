package main

import (
	"database/sql"
	"errors"
	"fmt"
)

const driver = "postgress"

var (
	// ErrMoreThanOneUserFound is an error for specifing that more than one user with the same username was found
	ErrMoreThanOneUserFound = errors.New("More than one user found with the same username")
	//ErrNoUsersFound that no user was found in the database
	ErrNoUsersFound = errors.New("No users found")
)

type user struct {
	username string
	password string
}

type readerWriter interface {
	read(username string) readResult
	write(user user) error
	delete(username string) error
	close()
}

type readResult struct {
	user user
	err  error
}

type userReaderWriter struct {
	*sql.DB
}

func newUserReaderWriter(dbUser string, dbPassword string, dbName string) (userReaderWriter, error) {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", dbUser, dbPassword, dbName)
	db, err := sql.Open(driver, dbinfo)
	if err != nil {
		return userReaderWriter{}, err
	}

	return userReaderWriter{db}, nil
}

func (rw userReaderWriter) read(username string) readResult {

	rows, err := rw.Query("SELECT * FROM users WHERE username=$1", username)
	if err != nil {
		return readResult{user: user{}, err: err}
	}

	for rows.Next() {
		var username string
		var password string

		err := rows.Scan(&username, &password)
		if err != nil {
			return readResult{user: user{}, err: err}
		}

		if rows.Next() {
			return readResult{user: user{}, err: ErrMoreThanOneUserFound}
		}

		return readResult{user: user{username: username, password: password}, err: nil}
	}

	return readResult{user: user{}, err: ErrNoUsersFound}
}

func (rw userReaderWriter) write(user user) error {
	return rw.QueryRow("INSERT INTO users(usrname,password) VALUES ($1,$2)", user.username, user.password).Scan()
}

func (rw userReaderWriter) delete(username string) error {
	stmt, err := rw.Prepare("DELETE FROM users WHERE username=$1")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(username)
	return err
}

func (rw userReaderWriter) close() {
	rw.close()
}
