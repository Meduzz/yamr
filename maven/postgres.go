package maven

import (
	"io"
	"errors"
	_ "github.com/lib/pq"
	"github.com/Meduzz/yamr/artifacts"
)

// This is the PipeItem, recording files.
type PostgresPipeItem struct {
}

var artifactManager = artifacts.NewArtifacts()

func NewPostgresAdapter() PipeItem {
	err := SetupDatabase()
	if err != nil {
		panic(err)
	}
	return &PostgresPipeItem{}
}

func (p *PostgresPipeItem) Write(context *Context, bytes io.ReadCloser) error {
	meta := context.Get(FILEMETADATA).(*artifacts.FileMetadata)
	packages := context.Get(PACKAGE).(*Package)

	// let the items further down in the pipe, handle the actual write.
	err := context.Write(bytes)

	if err != nil {
		return err
	}

	err = artifactManager.Store(meta, packages.Id)

	if err != nil {
		return err
	}

	return nil
}

func (p *PostgresPipeItem) Exists(context *Context) (bool, error) {
	meta := context.Get(FILEMETADATA).(*artifacts.FileMetadata)
	torf, err := artifactManager.Exists(meta)

	if err != nil {
		return false, err
	}

	return torf, nil
}

func (p *PostgresPipeItem) Read(context *Context) ([]byte, error) {
	meta := context.Get(FILEMETADATA).(*artifacts.FileMetadata)
	torf, err := artifactManager.Exists(meta)

	if err != nil {
		return nil, err
	}

	// if the file exists in db, delegate to lower pipeItems to fetch it.
	if torf {
		return context.Read()
	}

	return nil, errors.New("File not found.")
}