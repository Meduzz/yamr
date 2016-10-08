package repo

import (
	"io"
	"errors"
	"database/sql"
	_ "github.com/lib/pq"
)

// This is the PipeItem, recording files.
type PostgresAdapter struct {
}

func NewPostgresAdapter() PipeItem {
	err := SetupDatabase()
	if err != nil {
		panic(err)
	}
	return &PostgresAdapter{}
}

func (p *PostgresAdapter) Write(context *Context, bytes io.ReadCloser) error {
	meta := context.Get(FILEMETADATA).(*FileMetadata)
	// let the items further down in the pipe, handle the actual write.
	err := context.Write(bytes)

	if err != nil {
		return err
	}

	err = insert(meta.GroupAsPackage(), meta.Artifact, meta.Version, meta.File)

	if err != nil {
		return err
	}

	return nil
}

func (p *PostgresAdapter) Exists(context *Context) bool {
	meta := context.Get(FILEMETADATA).(*FileMetadata)
	torf, err := exists(meta.GroupAsPackage(), meta.Artifact, meta.Version, meta.File)

	if err != nil {
		return false
	}

	return torf
}

func (p *PostgresAdapter) Read(context *Context) ([]byte, error) {
	meta := context.Get(FILEMETADATA).(*FileMetadata)
	torf, err := exists(meta.GroupAsPackage(), meta.Artifact, meta.Version, meta.File)

	if err != nil {
		return nil, err
	}

	// if the file exists in db, delegate to lower pipeItems to fetch it.
	if torf {
		return context.Read()
	}

	return nil, errors.New("File not found.")
}

func insert(group string, artifact string, version string, file string) error {
	conn, err := sql.Open("postgres", "")

	if err != nil {
		return err
	}

	defer conn.Close()

	// insert into
	_, err = conn.Exec("insert into artifacts (groupname, artifactname, version, filename) values ($1, $2, $3, $4)", group, artifact, version, file)

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