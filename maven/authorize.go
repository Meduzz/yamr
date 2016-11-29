package maven

import (
	"io"
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"errors"
	"github.com/Meduzz/yamr/artifacts"
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
		Id int64
		Name string
		Password string
		Public bool
	}
)

const AUTHORIZATIONS = "authorizations"
const PACKAGE = "package"

func NewAuthorizeAdapter() *AuthorizePipeItem {
	return &AuthorizePipeItem{}
}

// TODO check that credential.username owns this package.
func (a *AuthorizePipeItem) Write(context *Context, bytes io.ReadCloser) error {
	meta := context.Get(FILEMETADATA).(*artifacts.FileMetadata)
	credentials := credentialsOrNull(context)
	packageDetails, err := authorizationForGroup(meta.GroupAsPackage())

	if err != nil {
		return err
	} else if credentials == nil {
		return errors.New("Access denied.")
	} else {
		if packageDetails.Password == credentials.Password {
			context.Set(PACKAGE, packageDetails)
			return context.Write(bytes)
		}
		return errors.New("Invalid credentials.")
	}
}

func (a *AuthorizePipeItem) Read(context *Context) ([]byte, error) {
	meta := context.Get(FILEMETADATA).(*artifacts.FileMetadata)
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
	meta := context.Get(FILEMETADATA).(*artifacts.FileMetadata)
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

	row := conn.QueryRow("select id, groupName, public, password from packages where groupName = $1", group)

	p := &Package{}
	err = row.Scan(&p.Id, &p.Name, &p.Public, &p.Password)

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