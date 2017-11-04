package maven

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"

	"cloud.google.com/go/storage"
	"github.com/Meduzz/yamr/artifacts"
	"golang.org/x/net/context"
)

var BUCKET string

type (
	GceStorage struct{}
)

func NewGceStorage() *GceStorage {
	BUCKET = fromEnv("BUCKET", "")

	if BUCKET == "" {
		panic("A bucket environment parameter must be set (BUCKET=<bucket name>).")
	}

	return &GceStorage{}
}

func (g *GceStorage) Write(ctx *Context, data io.ReadCloser) error {
	meta := ctx.Get(FILEMETADATA).(*artifacts.FileMetadata)
	key := mkKey(meta)

	gctx := context.Background()
	gce, err := storage.NewClient(gctx)

	if err != nil {
		return err
	}

	defer gce.Close()

	writer := gce.Bucket(BUCKET).Object(key).NewWriter(gctx)

	written, err := io.Copy(writer, data)

	log.Printf("Wrote %d bytes to %s.", written, key)

	if err != nil {
		return err
	}

	return writer.Close()
}

func (g *GceStorage) Exists(ctx *Context) (bool, error) {
	meta := ctx.Get(FILEMETADATA).(*artifacts.FileMetadata)
	key := mkKey(meta)

	gctx := context.Background()
	gce, err := storage.NewClient(gctx)

	if err != nil {
		return false, err
	}

	defer gce.Close()

	reader, err := gce.Bucket(BUCKET).Object(key).NewRangeReader(gctx, 0, 1)

	if err != nil {
		return false, err
	} else {
		reader.Close()
		return true, nil
	}
}

func (g *GceStorage) Read(ctx *Context) ([]byte, error) {
	meta := ctx.Get(FILEMETADATA).(*artifacts.FileMetadata)
	key := mkKey(meta)

	gctx := context.Background()
	gce, err := storage.NewClient(gctx)

	if err != nil {
		return nil, err
	}

	defer gce.Close()

	reader, err := gce.Bucket(BUCKET).Object(key).NewReader(gctx)

	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(reader)
}

func mkKey(meta *artifacts.FileMetadata) string {
	root := meta.GroupAsPackage()
	file := meta.File

	return fmt.Sprintf("%s.%s", root, file)
}
