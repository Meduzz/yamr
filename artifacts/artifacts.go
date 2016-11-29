package artifacts

import (
	"database/sql"
	"strings"
)

type (
	Artifacts struct { }

	Artifact struct {
		Id int64
		Group string
		Name string
		Version string
		File string
	}
)

func NewArtifacts() *Artifacts {
	return &Artifacts{}
}

func (a *Artifacts) Store(meta *FileMetadata, packageId int64) error {
	return insert(meta.GroupAsPackage(), meta.Artifact, meta.Version, meta.File, packageId)
}

func (a *Artifacts) Exists(meta *FileMetadata) (bool, error) {
	return exists(meta.GroupAsPackage(), meta.Artifact, meta.Version, meta.File)
}

func (a *Artifacts) Search(query string, userId int64, page, limit int) ([]*Artifact, error) {
	conn, err := sql.Open("postgres", "")

	if err != nil {
		return nil, err
	}

	defer conn.Close()

	page = page * limit

	query = strings.Replace(query, "*", "%", -1)
	if !strings.Contains(query, "%") {
		query = strings.Join([]string{"%", query, "%"}, "")
	}

	rows, err := conn.Query("select a.id, a.groupName, a.artifactName, a.version, a.filename from artifacts a left join packages p on (a.package_id = p.id) where (p.public = true or p.userId = $2) and (a.groupname like ($1) or a.artifactname like ($1)) limit $3 offset $4", query, userId, limit, page)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	result := make([]*Artifact, 0)

	for rows.Next() {
		artifact := &Artifact{}
		err = rows.Scan(&artifact.Id, &artifact.Group, &artifact.Name, &artifact.Version, &artifact.File)

		if err != nil {
			return nil, err
		}

		result = append(result, artifact)
	}

	return result, nil
}

func insert(group string, artifact string, version string, file string, packageId int64) error {
	conn, err := sql.Open("postgres", "")

	if err != nil {
		return err
	}

	defer conn.Close()

	// insert into
	_, err = conn.Exec("insert into artifacts (groupname, artifactname, version, filename, package_id) values ($1, $2, $3, $4, $5)", group, artifact, version, file, packageId)

	if err != nil {
		return err
	}

	return nil
}

func exists(group string, artifact string, version string, file string) (bool, error) {
	conn, err := sql.Open("postgres", "")

	if err != nil {
		return false, err
	}

	defer conn.Close()

	// select count(*)
	row := conn.QueryRow("select count(id) from artifacts where groupname=$1 and artifactname=$2 and version=$3 and filename=$4", group, artifact, version, file)

	var count int = 0
	err = row.Scan(&count)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}