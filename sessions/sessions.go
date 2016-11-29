package sessions

import (
	"time"
	"database/sql"
	"crypto/sha256"
	"encoding/hex"
)

type (

	Session struct {
		Id string
		Created time.Time
		Expires time.Time
		Ip string
		UserId int64
	}

	Sessions struct {
	}
)

const SESSION_LIFE = 30

func NewSessions() *Sessions {
	return &Sessions{}
}

func (s *Sessions) LoadById(id string) (*Session, error) {
	conn, err := sql.Open("postgres", "")

	if err != nil {
		return nil, err
	}

	defer conn.Close()

	row := conn.QueryRow("select id, created, expires, ip, userId from sessions where id = $1", id)

	session := &Session{}
	var createdRead int64
	var expiresRead int64
	err = row.Scan(&session.Id, &createdRead, &expiresRead, &session.Ip, &session.UserId)

	session.Created = time.Unix(createdRead, 0)
	session.Expires = time.Unix(expiresRead, 0)

	if err != nil {
		return nil, err
	}

	return session, nil
}

func (s *Sessions) CreateForUser(userId int64, ip string) (*Session, error) {
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

	_, err = conn.Exec("insert into sessions (id, created, expires, ip, userId) values ($1, $2, $3, $4, $5)", hash(string(id)), now.Unix(), expires.Unix(), ip, userId)

	if err != nil {
		return nil, err
	}

	return s.LoadById(hash(string(id)))
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

func hash(value string) string {
	hash := sha256.New()
	hash.Write([]byte(value))

	return hex.EncodeToString(hash.Sum(nil))
}
