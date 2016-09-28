package main

import (
	"io"
	"os"
	"log"
)

type FilesystemAdapter struct {
	url string
}

func NewFileSystemAdapter(url string) FilesystemAdapter {
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

	return FilesystemAdapter{baseUrl}
}

func (fs FilesystemAdapter) Write(bytes io.ReadCloser, meta *FileMetadata) error {
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

	log.Printf("%s bytes stored into %s", read, meta.Path())

	return nil
}

// TODO implement
func (fs FilesystemAdapter) Read(meta *FileMetadata) (string, error) {
	return "", nil
}
