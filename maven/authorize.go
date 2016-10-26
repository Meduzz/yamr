package maven

import (
	"io"
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"errors"
)

type (
	AuthorizePipeItem struct {
	}

	Credential struct {
		Username string
		Password string
	}

	Package struct {
		Name string
		Read bool
		Username string
		Password string
	}
)

const AUTHORIZATIONS = "authorizations"

func NewAuthorizeAdapter() *AuthorizePipeItem {
	return &AuthorizePipeItem{}
}

func (a *AuthorizePipeItem) Write(context *Context, bytes io.ReadCloser) error {
	meta := context.Get(FILEMETADATA).(*FileMetadata)
	credentials := credentialsOrNull(context)
	auths := authorizationsStartingWith(meta.TopDomain("%"))

	if len(auths) == 0 {
		return context.Write(bytes)
	} else if credentials != nil {
		for _, auth := range(auths) {
			if auth.Name == meta.TopDomain("") {
				if auth.Username == credentials.Username && auth.Password == credentials.Password {
					return context.Write(bytes)
				}
			} else if auth.Name == meta.GroupAsPackage() {
				if auth.Username == credentials.Username && auth.Password == credentials.Password {
					return context.Write(bytes)
				}
			}
		}
	}

	return errors.New("Access denied")
}

func (a *AuthorizePipeItem) Read(context *Context) ([]byte, error) {
	meta := context.Get(FILEMETADATA).(*FileMetadata)
	credentials := credentialsOrNull(context)
	auths := authorizationsStartingWith(meta.TopDomain("%"))

	if len(auths) == 0 {
		return context.Read()
	} else {
		for _, auth := range(auths) {
			if auth.Read {
				return context.Read()
			} else if credentials == nil {
				return nil, errors.New("Access denied")
			} else {
				if auth.Name == meta.TopDomain("") {
					if auth.Username == credentials.Username && auth.Password == credentials.Password {
						return context.Read()
					}
				} else if auth.Name == meta.GroupAsPackage() {
					if auth.Username == credentials.Username && auth.Password == credentials.Password {
						return context.Read()
					}
				}
			}
		}
	}

	return nil, errors.New("Access denied")
}

func (a *AuthorizePipeItem) Exists(context *Context) (bool, error) {
	meta := context.Get(FILEMETADATA).(*FileMetadata)
	credentials := credentialsOrNull(context)
	auths := authorizationsStartingWith(meta.TopDomain("%"))

	if len(auths) == 0 {
		return context.Exists()
	} else {
		for _, auth := range(auths) {
			if auth.Read {
				return context.Exists()
			} else if credentials == nil {
				return false, errors.New("Access denied")
			} else {
				if auth.Name == meta.TopDomain("") {
					if auth.Username == credentials.Username && auth.Password == credentials.Password {
						return context.Exists()
					}
				} else if auth.Name == meta.GroupAsPackage() {
					if auth.Username == credentials.Username && auth.Password == credentials.Password {
						return context.Exists()
					}
				}
			}
		}
	}

	return false, errors.New("Access denied")
}

func authorizationsStartingWith(group string) []Package {
	conn, err := sql.Open("postgres", "")

	matches := make([]Package, 0)

	if err != nil {
		log.Printf("There was an error connecting to db: %s.", err)
		return matches
	}

	defer conn.Close()

	rows, err := conn.Query("select name, read, username, password from packages where name like $1", group)

	if  err != nil {
		log.Printf("Error executing query: %s.", err)
		return matches
	}

	defer rows.Close()

	for rows.Next() {
		p := new(Package)

		err = rows.Scan(&p.Name, &p.Read, &p.Username, &p.Password)
		if err != nil {
			log.Printf("There was an error fetching data from db. (%s)", err)
		} else {
			matches = append(matches, *p)
		}
	}

	return matches
}

func credentialsOrNull(c *Context) *Credential {
	credentials := c.Get(AUTHORIZATIONS)

	if credentials != nil {
		return credentials.(*Credential)
	} else {
		return nil
	}
}

