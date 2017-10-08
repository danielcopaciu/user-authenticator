package main

import (
	"database/sql"
	"errors"
	"fmt"
)

const driver = "postgress"

var (
	ErrMoreThanOneUserFound = errors.New("More than one user found with the same username")
	ErrNoUsersFound         = errors.New("No users found")
)

type readerWriter interface {
	read(username string) readResult
	write(user User) error
	delete(username string) error
	close()
}

type readResult struct {
	user User
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
		return readResult{user: User{}, err: err}
	}

	for rows.Next() {
		var username string
		var password string

		err := rows.Scan(&username, &password)
		if err != nil {
			return readResult{user: User{}, err: err}
		}

		if rows.Next() {
			return readResult{user: User{}, err: ErrMoreThanOneUserFound}
		}

		return readResult{user: User{username: username, password: password}, err: nil}
	}

	return readResult{user: User{}, err: ErrNoUsersFound}
}

func (rw userReaderWriter) write(user User) error {
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
