package main

type Db struct {
	// TODO: Implement database logic here
}

func NewDb() *Db {
	return &Db{}
}

func (db *Db) authenticate(username, password string) bool {
	// TODO: Implement authentication logic here
	return true
}

func (db *Db) createUser(username, password string) error {
	// TODO: Implement user creation logic here
	return nil
}

func (db *Db) updateUser(username, password string) error {
	// TODO: Implement user update logic here
	return nil
}

func (db *Db) deleteUser(username string) error {
	// TODO: Implement user deletion logic here
	return nil
}
