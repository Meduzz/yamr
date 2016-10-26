package maven

import (
	"io"
	"os"
	"log"
	"io/ioutil"
)

type FilesystemPipeItem struct {
	url string
}

func NewFileSystemAdapter(url string) *FilesystemPipeItem {
	var baseUrl string

	if url[0] != '/' {
		dir, err := os.Getwd()

		if err != nil {
			// TODO room for improvements
			log.Fatalf("Uh oh! (%s)", url)
		}

		baseUrl = dir + "/" + url
	} else {
		baseUrl = url
	}

	log.Printf("Using %s as base dir.", baseUrl)

	return &FilesystemPipeItem{baseUrl}
}

func (fs *FilesystemPipeItem) Write(context *Context, bytes io.ReadCloser) error {
	// I expect this will be ran in it's own go-routine sooner or later.
	meta := context.Get(FILEMETADATA).(*FileMetadata)
	err := os.MkdirAll(fs.url + "/" + meta.GroupAsPath() + "/" + meta.Version, 0755)
	if err != nil && !os.IsExist(err) {
		return err
	}

	tmp, err := os.Create(fs.url + "/" + meta.Path())
	if (err != nil) {
		return err
	}

	read, err := io.Copy(tmp, bytes)
	if err != nil {
		return err
	}

	log.Printf("%d bytes stored into %s", read, meta.Path())

	return nil
}

func (fs *FilesystemPipeItem) Read(context *Context) ([]byte, error) {
	meta := context.Get(FILEMETADATA).(*FileMetadata)
	file, err := os.Open(fs.url + "/" + meta.Path())

	if err != nil {
		return nil, err
	} else {
		defer file.Close()
		return ioutil.ReadAll(file)
	}
}

func (fs *FilesystemPipeItem) Exists(context *Context) (bool, error) {
	meta := context.Get(FILEMETADATA).(*FileMetadata)
	file, err := os.Open(meta.Path())

	if err != nil {
		return false, err
	} else {
		defer file.Close()
		return true, nil
	}
}