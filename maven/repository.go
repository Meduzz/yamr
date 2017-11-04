package maven

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/Meduzz/yamr/artifacts"
	"github.com/gin-gonic/gin"
)

const FILEMETADATA = "file_metadata"

type (
	PipeItem interface {
		Write(*Context, io.ReadCloser) error
		Exists(*Context) (bool, error)
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
		log.Printf("Using filesystem (%s).", path)
		return NewFileSystemAdapter(path)
	} else if flag == "gce" {
		log.Println("Using Google Storage.")
		return NewGceStorage()
	} else {
		// TODO add more storage adapters.
		panic(fmt.Sprintf("No storage adapter matching %s found.", flag))
	}
}

func NewRepository(storageDriver string) *Repository {
	storageAdapter := selectStorageAdapter(storageDriver)
	postgresAdapter := NewPostgresAdapter()
	authorizeAdapter := NewAuthorizeAdapter()

	repo := &Repository{
		chain: make([]PipeItem, 0, 5),
	}
	repo.use(authorizeAdapter, postgresAdapter, storageAdapter)

	return repo
}

func (repo *Repository) use(pipes ...PipeItem) {
	repo.chain = append(repo.chain, pipes...)
}

func (repo *Repository) Context(g *gin.Context) *Context {
	any := g.Param("any")
	username, password, ok := g.Request.BasicAuth()

	c := NewContext(&repo.chain)
	c.Set(FILEMETADATA, extract(any))
	if ok {
		c.Set(AUTHORIZATIONS, &Credential{username, password})
	}

	return c
}

func (repo *Repository) Write(context *Context, bytes io.ReadCloser) error {
	return repo.chain[0].Write(context, bytes)
}

func (repo *Repository) Exists(context *Context) (bool, error) {
	return repo.chain[0].Exists(context)
}

func (repo *Repository) Read(context *Context) ([]byte, error) {
	return repo.chain[0].Read(context)
}

func fromEnv(param string, defaultVal string) string {
	env := os.Getenv(param)

	if env == "" {
		env = defaultVal
	}

	return env
}

func extract(path string) *artifacts.FileMetadata {
	split := strings.Split(path[1:], "/")

	binary := split[len(split)-1]
	version := split[len(split)-2]
	artifact := split[len(split)-3]
	group := split[0 : len(split)-3]

	return artifacts.NewFileMetadata(group, artifact, version, binary)
}
