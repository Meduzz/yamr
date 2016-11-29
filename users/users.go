package users

import (
	"database/sql"
	"crypto/sha256"
	"encoding/hex"
)

type (
	User struct {
		Id int64
		Username string
		Password string
		Package string
	}

	Users struct {
	}
)

// TODO turn these NewX functions into factories...
func NewUsers() *Users {
	return &Users{}
}

func (u *Users) Store(user *User) error {
	conn, err := sql.Open("postgres", "")

	if err != nil {
		return err
	}

	defer conn.Close()

	_, err = conn.Exec("insert into accounts (name, username, password) values ($1, $2, $3)", user.Package, user.Username, hash(user.Password))

	if err != nil {
		return err
	}

	return nil
}

func (u *Users) DomainExists(name string) (bool, error) {
	conn, err := sql.Open("postgres", "")

	if err != nil {
		return true, err
	}

	defer conn.Close()

	row := conn.QueryRow("select count(name) from accounts where name = $1", name)

	var count int64
	err = row.Scan(&count)

	if err != nil {
		return true, err
	}

	return count == 1, nil
}

func (u *Users) UserExists(name string) (bool, error) {
	conn, err := sql.Open("postgres", "")

	if err != nil {
		return true, err
	}

	defer conn.Close()

	row := conn.QueryRow("select count(username) from accounts where username = $1", name)

	var count int64
	err = row.Scan(&count)

	if err != nil {
		return true, err
	}

	return count == 1, nil
}

func (u *Users) LoadByUsernameAndPassword(username string, password string) (*User, error) {
	conn, err := sql.Open("postgres", "")

	if err != nil {
		return nil, err
	}

	defer conn.Close()

	row := conn.QueryRow("select id, username, password, name from accounts where username = $1 and password = $2", username, hash(password))

	user := &User{}
	err = row.Scan(&user.Id, &user.Username, &user.Password, &user.Package)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *Users) LoadById(id int64) (*User, error) {
	conn, err := sql.Open("postgres", "")

	if err != nil {
		return nil, err
	}

	defer conn.Close()

	row := conn.QueryRow("select id, username, password, name from accounts where id = $1", id)

	user := &User{}
	err = row.Scan(&user.Id, &user.Username, &user.Password, &user.Package)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func hash(value string) string {
	hash := sha256.New()
	hash.Write([]byte(value))

	return hex.EncodeToString(hash.Sum(nil))
}
