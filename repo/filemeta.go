package repo

import (
	"strings"
)

type FileMetadata struct {
	Group []string
	Artifact string
	Version string
	File string
}

func NewFileMetadata(group []string, artifact string, version string, file string) *FileMetadata {
	m := &FileMetadata{}

	m.Group = group
	m.Artifact = artifact
	m.Version = version
	m.File = file

	return m
}

func (m *FileMetadata) GroupAsPackage() string {
	return strings.Join(m.Group, ".")
}

func (m *FileMetadata) GroupAsPath() string {
	return strings.Join(m.Group, "/")
}

func (m *FileMetadata) Path() string {
	return m.GroupAsPath() + "/" + m.Version + "/" + m.File
}