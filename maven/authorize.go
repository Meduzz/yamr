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

	// TODO improvement when it comes to dual declare structs...
	// Or simply package stuff better...
	Package struct {
		Name string
		Password string
		Public bool
	}
)

const AUTHORIZATIONS = "authorizations"

func NewAuthorizeAdapter() *AuthorizePipeItem {
	return &AuthorizePipeItem{}
}

// TODO check that credentail.username owns this package.
func (a *AuthorizePipeItem) Write(context *Context, bytes io.ReadCloser) error {
	meta := context.Get(FILEMETADATA).(*FileMetadata)
	credentials := credentialsOrNull(context)
	packageDetails, err := authorizationForGroup(meta.GroupAsPackage())

	if err != nil {
		return err
	} else if credentials == nil {
		return errors.New("Access denied.")
	} else {
		if packageDetails.Password == credentials.Password {
			return context.Write(bytes)
		}
		return errors.New("Invalid credentials.")
	}
}

func (a *AuthorizePipeItem) Read(context *Context) ([]byte, error) {
	meta := context.Get(FILEMETADATA).(*FileMetadata)
	credentials := credentialsOrNull(context)
	packageDetails, err := authorizationForGroup(meta.GroupAsPackage())

	if err != nil {
		return nil, err
	} else {
		if packageDetails.Public {
			return context.Read()
		} else if credentials == nil {
			return nil, errors.New("Access denied.")
		} else if packageDetails.Password == credentials.Password {
			return context.Read()
		} else {
			return nil, errors.New("Invalid credentials.")
		}
	}
}

func (a *AuthorizePipeItem) Exists(context *Context) (bool, error) {
	meta := context.Get(FILEMETADATA).(*FileMetadata)
	credentials := credentialsOrNull(context)
	packageDetails, err := authorizationForGroup(meta.GroupAsPackage())

	if err != nil {
		return false, err
	} else {
		if packageDetails.Public {
			return context.Exists()
		} else if credentials == nil {
			return false, errors.New("Access denied.")
		} else if packageDetails.Password == credentials.Password {
			return context.Exists()
		} else {
			return false, errors.New("Invalid credentials.")
		}
	}
}

func authorizationForGroup(group string) (*Package, error) {
	conn, err := sql.Open("postgres", "")

	if err != nil {
		log.Printf("There was an error connecting to db: %s.", err)
		return nil, err
	}

	defer conn.Close()

	row := conn.QueryRow("select groupName, public, password from packages where groupName = $1", group)

	p := &Package{}
	err = row.Scan(&p.Name, &p.Public, &p.Password)

	if err != nil {
		log.Printf("There was an error fetching data from db. (%s)", err)
		return nil, err
	}

	return p, nil
}

func credentialsOrNull(c *Context) *Credential {
	credentials := c.Get(AUTHORIZATIONS)

	if credentials != nil {
		return credentials.(*Credential)
	} else {
		return nil
	}
}

/*
	Simplifications:
	1. Disallow writes without credentials.
	2. Only allow writes to packages specified in UI (with that package credential)
	3. Unless package.read is true, disallow reads to that package without credentials, and only allow with that package credential.
 */