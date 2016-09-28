package main

import (
	"github.com/gin-gonic/gin"
	"strings"
)

// TODO turn storage and path into flags and env-var.
var storage = "filesystem"
var path = "files"

var storageAdapter StorageAdapter = selectStorageAdapter(storage)

func extract(path string) *FileMetadata {
	split := strings.Split(path[1:], "/")

	binary := split[len(split) - 1]
	version := split[len(split) - 2]
	artifact := split[len(split) - 3]
	group := split[0:len(split) - 3]

	return NewFileMetadata(group, artifact, version, binary)
}

func selectStorageAdapter(flag string) StorageAdapter {
	if flag == "filesystem" {
		return NewFileSystemAdapter(path)
	} else {
		// TODO add more adapters.
		return nil
	}
}

func main() {

	webserver := gin.Default()

	webserver.PUT("/maven/*any", func(g *gin.Context) {
		any := g.Param("any")

		fileMeta := extract(any)
		err := storageAdapter.Write(g.Request.Body, fileMeta)

		if err != nil {
			g.Error(err)
		}

		g.Next()
	})

	// TODO implement
	webserver.HEAD("/maven/*any", func(g *gin.Context) {
		g.String(404, "")
	})

	// At publish
	// 1. HEAD .pom
	// 2. PUT .pom, .pom.sha1, .pom.md5
	// 3. HEAD .jar
	// 4. PUT .jar, .jar.sha1, .jar.md5
	// 5. HEAD .sources.jar
	// 6. PUT .sources.jar etc...
	// 7. HEAD .javadoc.jar
	// 8. PUT .javadoc.jar etc...

	// At read
	// 1. HEAD .pom
	// 2. GET .pom
	// 3. HEAD .pom.sha1 (twice?)
	// 4. GET .pom.sha1

	webserver.Run(":4040")
}
