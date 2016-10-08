package repo

import (
	"io"
	"fmt"
	"os"
	"log"
)

const FILEMETADATA = "file_metadata"

type (
	PipeItem interface {
		Write(*Context, io.ReadCloser) error
		Exists(*Context) bool
		Read(*Context) ([]byte, error)
	}

	Chain []PipeItem

	Repository struct {
		chain Chain
	}
)

func selectStorageAdapter(flag string) PipeItem {
	if flag == "filesystem" {
		path := fromEnv("FS_PATH", "files")
		return NewFileSystemAdapter(path)
	} else {
		// TODO add more storage adapters.
		panic(fmt.Sprintf("No storage adapter matching %s found.", flag))
	}
}

func NewRepository(storageDriver string) *Repository {
	storageAdapter := selectStorageAdapter(storageDriver)
	postgresAdapter := NewPostgresAdapter()

	repo := &Repository{
		chain:make([]PipeItem, 0, 5),
	}
	repo.use(postgresAdapter, storageAdapter)

	return repo
}

func (repo *Repository) use(pipes ...PipeItem) {
	repo.chain = append(repo.chain, pipes...)
}

func (repo *Repository) context(meta *FileMetadata) *Context {
	c := NewContext(&repo.chain)
	c.Set(FILEMETADATA, meta)

	return c
}

func (repo *Repository) Write(meta *FileMetadata, bytes io.ReadCloser) error {
	return repo.chain[0].Write(repo.context(meta), bytes)
}

func (repo *Repository) Exists(meta *FileMetadata) bool {
	return repo.chain[0].Exists(repo.context(meta))
}

func (repo *Repository) Read(meta *FileMetadata) ([]byte, error) {
	return repo.chain[0].Read(repo.context(meta))
}

func fromEnv(param string, defaultVal string) string {
	env := os.Getenv(param)

	if env == "" {
		env = defaultVal
	}

	return env
}