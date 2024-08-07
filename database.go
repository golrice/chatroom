package main

import (
	_ "github.com/go-sql-driver/mysql"
)

type Db struct {
	Host string
	Port int

	Name string
}

var db *Db

func NewDb(host string, port int, name string) *Db {
	return &Db{
		Host: host,
		Port: port,
		Name: name,
	}
}

func (db *Db) connect() {
	// connect to mysql
}

func (db *Db) authenticate(username, password string) bool {
	// TODO: Implement authentication logic here
	return true
}

func (db *Db) existsUser(username string) bool {
	// TODO: Implement user existence check logic here
	return false
}

func (db *Db) createUser(username, password, email string) error {
	// TODO: Implement user creation logic here
	return nil
}

func (db *Db) updateUser(username, password, email string) error {
	// TODO: Implement user update logic here
	return nil
}

func (db *Db) deleteUser(username string) error {
	// TODO: Implement user deletion logic here
	return nil
}
