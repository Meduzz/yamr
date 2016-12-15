package domains

import (
	"database/sql"
)

type (
	Domain struct {
		Id int64 `json:",omitempty"`
		Name string
		Active bool `json:",omitempty"`
	}

	Domains struct {
	}
)

func NewDomains() *Domains {
	return &Domains{}
}

func (d *Domains) Create(domain *Domain, userId int64) error {
	conn, err := sql.Open("postgres", "")

	if err != nil {
		return err
	}

	defer conn.Close()

	_, err = conn.Exec("insert into domains (name, active, userid) values ($1, $2, $3)", domain.Name, false, userId)

	if err != nil {
		return err
	}

	return nil
}

func (d *Domains) ListDomainsForUser(userId int64, skip, limit int) ([]*Domain, error) {
	conn, err := sql.Open("postgres", "")

	if err != nil {
		return nil, err
	}

	defer conn.Close()

	rows, err := conn.Query("select id, name, active from domains where userId = $1 limit $3 offset $2", userId, skip, limit)

	if err != nil {
		return nil, err
	}

	data := make([]*Domain, 0)

	for rows.Next() {
		domain := &Domain{}
		err = rows.Scan(&domain.Id, &domain.Name, &domain.Active)

		if err!= nil {
			return nil, err
		}

		data = append(data, domain)
	}

	return data, nil
}

func (d *Domains) FindInactiveDomains(skip, limit int) ([]*Domain, error) {
	conn, err := sql.Open("postgres", "")

	if err != nil {
		return nil, err
	}

	defer conn.Close()

	rows, err := conn.Query("select id, name, active from domains where active = $1 limit $3 offset $2", false, skip, limit)

	if err != nil {
		return nil, err
	}

	data := make([]*Domain, 0)

	for rows.Next() {
		domain := &Domain{}
		err = rows.Scan(&domain.Id, &domain.Name, &domain.Active)

		if err!= nil {
			return nil, err
		}

		data = append(data, domain)
	}

	return data, nil
}

func (d *Domains) GetById(id int64) (*Domain, error) {
	conn, err := sql.Open("postgres", "")

	if err != nil {
		return nil, err
	}

	defer conn.Close()

	row := conn.QueryRow("select id, name, active from domains where id = $1", id)
	domain := &Domain{}

	err = row.Scan(domain.Id, domain.Name, domain.Active)

	if err != nil {
		return nil, err
	}

	return domain, nil
}

func (d *Domains) Activate(id int64) error {
	return setActive(id, true)
}

func (d *Domains) DeActivate(id int64) error {
	return setActive(id, false)
}

func setActive(id int64, active bool) error {
	conn, err := sql.Open("postgres", "")

	if err != nil {
		return err
	}

	defer conn.Close()

	_, err = conn.Exec("update domains set active = $1 where id = $2", active, id)

	if err != nil {
		return err
	}

	return nil
}

func (d *Domains) OwnedBy(domainId, userId int64) bool {
	conn, err := sql.Open("postgres", "")

	if err != nil {
		return false
	}

	defer conn.Close()

	row := conn.QueryRow("select count(id) from domains where id = $1 and userId = $2", domainId, userId)

	var count int64
	err = row.Scan(&count)

	if err != nil {
		return false
	}

	return count == 1
}

func (d *Domains) DomainByPackage(id int64) (*Domain, error) {
	conn, err := sql.Open("postgres", "")

	if err != nil {
		return nil, err
	}

	defer conn.Close()

	row := conn.QueryRow("select d.id, d.name, d.active from domains d left join packages p on (d.id = p.domainId) where p.id = $1", id)

	domain := &Domain{}
	err = row.Scan(&domain.Id, &domain.Name, &domain.Active)

	if err != nil {
		return nil, err
	}

	return domain, nil
}