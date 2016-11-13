package user

import (
	"database/sql"
	"crypto/sha256"
	"encoding/hex"
	"time"
	"strings"
)

type (
	User struct {
		Username string
		Password string
		Package string
	}

	Users struct {
	}

	Session struct {
		Id string
		Package string
		Created time.Time
		Expires time.Time
		Ip string
	}

	Sessions struct {
	}

	Package struct {
		Name string
		Password string
		Public bool
	}

	Packages struct {
	}
)

const SESSION_LIFE = 30

// TODO turn these NewX functions into factories...
func NewUsers() *Users {
	return &Users{}
}

func NewSessions() *Sessions {
	return &Sessions{}
}

func NewPackages() *Packages {
	return &Packages{}
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

	row := conn.QueryRow("select username, password, name from accounts where username = $1 and password = $2", username, hash(password))

	user := &User{}
	err = row.Scan(&user.Username, &user.Password, &user.Package)

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

func (s *Sessions) LoadById(id string) (*Session, error) {
	conn, err := sql.Open("postgres", "")

	if err != nil {
		return nil, err
	}

	defer conn.Close()

	row := conn.QueryRow("select id, package, created, expires, ip from sessions where id = $1", id)

	session := &Session{}
	var createdRead int64
	var expiresRead int64
	err = row.Scan(&session.Id, &session.Package, &createdRead, &expiresRead, &session.Ip)

	session.Created = time.Unix(createdRead, 0)
	session.Expires = time.Unix(expiresRead, 0)

	if err != nil {
		return nil, err
	}

	return session, nil
}

func (s *Sessions) CreateForUser(user string, ip string) (*Session, error) {
	conn, err := sql.Open("postgres", "")

	if err != nil {
		return nil, err
	}

	defer conn.Close()

	row := conn.QueryRow("select nextval('session_seq')")

	id := int64(0)
	err = row.Scan(&id)

	if err != nil {
		return nil, err
	}

	now := time.Now()
	expires := now.Add(time.Minute * SESSION_LIFE)

	_, err = conn.Exec("insert into sessions (id, package, created, expires, ip) values ($1, $2, $3, $4, $5)", hash(string(id)), user, now.Unix(), expires.Unix(), ip)

	if err != nil {
		return nil, err
	}

	row = conn.QueryRow("select id, package, created, expires, ip from sessions where id = $1", hash(string(id)))

	session := &Session{}
	var createdRead int64
	var expiresRead int64
	err = row.Scan(&session.Id, &session.Package, &createdRead, &expiresRead, &session.Ip)

	session.Created = time.Unix(createdRead, 0)
	session.Expires = time.Unix(expiresRead, 0)

	if err != nil {
		return nil, err
	}

	return session, nil
}

func (s *Sessions) Extend(session *Session) error {
	conn, err := sql.Open("postgres", "")

	if err != nil {
		return err
	}

	defer conn.Close()

	now := time.Now()
	expires := now.Add(time.Minute * SESSION_LIFE)

	_, err = conn.Exec("update sessions set expires = $1 where id = $2", expires.Unix(), session.Id)

	if err != nil {
		return err
	}

	return nil
}

func (p *Packages) UpdateOrCreate(pac *Package) error {
	conn, err := sql.Open("postgres", "")

	if err != nil {
		return err
	}

	defer conn.Close()

	row := conn.QueryRow("select count(groupName) from packages where groupName = $1", pac.Name)
	var count int64

	err = row.Scan(&count)

	if err != nil {
		return err
	}

	if count > 0 {
		_, err = conn.Exec("update packages set password = $1, public = $2 where groupName = $3", pac.Password, pac.Public, pac.Name)
	} else {
		_, err = conn.Exec("insert into packages (groupName, password, public) values ($1, $2, $3)", pac.Name, pac.Password, pac.Public)
	}

	if err != nil {
		return err
	}

	return nil
}

// TODO add pagination support.
func (p *Packages) List(topDomain string) ([]*Package, error) {
	conn, err := sql.Open("postgres", "")

	if err != nil {
		return nil, err
	}

	defer conn.Close()

	if !strings.HasSuffix(topDomain, "%") {
		topDomain += "%"
	}

	rows, err := conn.Query("select * from packages where groupName like $1", topDomain)

	if err != nil {
		return nil, err
	}

	packages := make([]*Package, 0)
	defer rows.Close()

	for rows.Next() {
		pac := &Package{}
		err = rows.Scan(&pac.Name, &pac.Password, &pac.Public)

		if err != nil {
			return nil, err
		} else {
			packages = append(packages, pac)
		}
	}

	return packages, nil
}

func (p *Packages) Load(groupName string) (*Package, error) {
	conn, err := sql.Open("postgres", "")

	if err != nil {
		return nil, err
	}

	defer conn.Close()

	row := conn.QueryRow("select * from packages where groupName = $1", groupName)

	pac := &Package{}
	err = row.Scan(&pac.Name, &pac.Password, &pac.Public)

	if err != nil {
		return nil, err
	}

	return pac, nil
}