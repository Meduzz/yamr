package packages

import (
	"database/sql"
)

type (

	Package struct {
		Id int64 `json:",omitempty"`
		Name string
		Password string
		Public bool
	}

	Packages struct {
	}

)

func NewPackages() *Packages {
	return &Packages{}
}

func (p *Packages) UpdateOrCreate(userId int64, pac *Package) error {
	conn, err := sql.Open("postgres", "")

	if err != nil {
		return err
	}

	defer conn.Close()

	if pac.Id > 0 {
		_, err = conn.Exec("update packages set password = $1, public = $2 where id = $3", pac.Password, pac.Public, pac.Id)
	} else {
		_, err = conn.Exec("insert into packages (groupName, password, public, userId) values ($1, $2, $3, $4)", pac.Name, pac.Password, pac.Public, userId)
	}

	if err != nil {
		return err
	}

	return nil
}

func (p *Packages) List(userId int64, page, limit int) ([]*Package, error) {
	conn, err := sql.Open("postgres", "")

	if err != nil {
		return nil, err
	}

	defer conn.Close()

	page = page * limit

	rows, err := conn.Query("select id, groupname, password, public from packages where userId = $1 limit $2 offset $3", userId, limit, page)

	if err != nil {
		return nil, err
	}

	packages := make([]*Package, 0)
	defer rows.Close()

	for rows.Next() {
		pac := &Package{}
		err = rows.Scan(&pac.Id, &pac.Name, &pac.Password, &pac.Public)

		if err != nil {
			return nil, err
		} else {
			packages = append(packages, pac)
		}
	}

	return packages, nil
}

func (p *Packages) Load(id int64) (*Package, error) {
	conn, err := sql.Open("postgres", "")

	if err != nil {
		return nil, err
	}

	defer conn.Close()

	row := conn.QueryRow("select id, groupname, password, public from packages where id = $1", id)

	pac := &Package{}
	err = row.Scan(&pac.Id, &pac.Name, &pac.Password, &pac.Public)

	if err != nil {
		return nil, err
	}

	return pac, nil
}