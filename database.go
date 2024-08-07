package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type Db struct {
	Host string
	Port int

	Name string

	Database *sql.DB
}

var db *Db

func NewDb(host string, port int, name string) (*Db, error) {
	var database *Db
	if db, err :=
		sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", "root", "root", host, port, name)); err != nil {
		return nil, err
	} else {
		database = &Db{
			Host:     host,
			Port:     port,
			Name:     name,
			Database: db,
		}
	}

	return database, nil
}

func (db *Db) authenticate(username, password string) bool {
	var dbPassword string
	err := db.Database.QueryRow("SELECT password FROM users WHERE username = ?", username).Scan(&dbPassword)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Printf("Error authenticating user: %v", err)
		}
		return false
	}

	return dbPassword == password
}

func (db *Db) existsUser(username string) bool {
	var user string
	err := db.Database.QueryRow("SELECT username FROM users WHERE username = ?", username).Scan(&user)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("Error checking user existence: %v", err)
	}
	return err != sql.ErrNoRows
}

func (db *Db) createUser(username, password, email string) error {
	result, err := db.Database.Exec("INSERT INTO users (username, password, email) VALUES (?,?,?)", username, password, email)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no rows were inserted for user %s", username)
	}

	return nil
}

func (db *Db) updateUser(username, password, email string) error {
	if !db.existsUser(username) {
		return fmt.Errorf("user %s does not exist", username)
	}

	result, err := db.Database.Exec("UPDATE users SET password = ?, email = ? WHERE username = ?", password, email, username)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no rows were updated for user %s", username)
	}

	return nil
}

func (db *Db) deleteUserByEmail(email string) error {
	result, err := db.Database.Exec("DELETE FROM users WHERE email = ?", email)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no rows were deleted for email %s", email)
	}

	return nil
}
