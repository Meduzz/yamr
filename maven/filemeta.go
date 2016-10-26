package maven

import (
	"strings"
	"fmt"
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

func (m *FileMetadata) TopDomain(append string) string {
	return fmt.Sprintf("%s.%s%s", m.Group[0], m.Group[1], append)
}