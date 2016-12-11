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
		Admin bool
	}

	Users struct {
	}
)

func NewUsers() *Users {
	return &Users{}
}

func (u *Users) Store(user *User) error {
	conn, err := sql.Open("postgres", "")

	if err != nil {
		return err
	}

	defer conn.Close()

	_, err = conn.Exec("insert into accounts (username, password) values ($1, $2)", user.Username, hash(user.Password))

	if err != nil {
		return err
	}

	return nil
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

	row := conn.QueryRow("select id, username, password, admin from accounts where username = $1 and password = $2", username, hash(password))

	user := &User{}
	err = row.Scan(&user.Id, &user.Username, &user.Password, &user.Admin)

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

	row := conn.QueryRow("select id, username, password, admin from accounts where id = $1", id)

	user := &User{}
	err = row.Scan(&user.Id, &user.Username, &user.Password, &user.Admin)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *Users) ListUsers(skip, limit int) ([]*User, error) {
	conn, err := sql.Open("postgres", "")

	if err != nil {
		return nil, err
	}

	defer conn.Close()

	rows, err := conn.Query("select * from accounts limit $2 offset $1", skip, limit)

	if err != nil {
		return nil, err
	}

	data := make([]*User, 0)
	for rows.Next() {
		row := &User{}
		err = rows.Scan(&row.Id, &row.Username, &row.Password, &row.Admin)

		if err != nil {
			return nil, err
		}

		data = append(data, row)
	}

	return data, nil
}

func hash(value string) string {
	hash := sha256.New()
	hash.Write([]byte(value))

	return hex.EncodeToString(hash.Sum(nil))
}
