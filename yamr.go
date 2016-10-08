package main

import (
	"github.com/gin-gonic/gin"
	"strings"
	"os"
	"fmt"
	"flag"
	"./repo"
)

var storage = flag.String("storage", "filesystem", "The storage engine to use (defaults to \"filesystem\").")

func main() {
	flag.Parse()

	port := fromEnv("PORT", "4040")

	repository := repo.NewRepository(*storage)

	webserver := gin.Default()

	webserver.PUT("/maven/*any", func(g *gin.Context) {
		any := g.Param("any")

		err := repository.Write(extract(any), g.Request.Body)

		if err != nil {
			g.Error(err)
		}

		g.Next()
	})

	webserver.HEAD("/maven/*any", func(g *gin.Context) {
		any := g.Param("any")

		exists := repository.Exists(extract(any))

		if exists {
			g.String(200, "")
		} else {
			g.String(404, "")
		}
	})

	webserver.GET("/maven/*any", func(g *gin.Context) {
		any := g.Param("any")

		bytes, err := repository.Read(extract(any))

		if err != nil {
			g.Error(err)
		} else {
			g.Status(200)
			g.Writer.Write(bytes)
		}
	})

	webserver.Run(fmt.Sprintf(":%s", port))
}

func fromEnv(param string, defaultVal string) string {
	env := os.Getenv(param)

	if env == "" {
		env = defaultVal
	}

	return env
}

func extract(path string) *repo.FileMetadata {
	split := strings.Split(path[1:], "/")

	binary := split[len(split) - 1]
	version := split[len(split) - 2]
	artifact := split[len(split) - 3]
	group := split[0:len(split) - 3]

	return repo.NewFileMetadata(group, artifact, version, binary)
}
