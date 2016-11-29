package main

import (
	"github.com/gin-gonic/gin"
	"os"
	"fmt"
	"flag"
	"github.com/Meduzz/yamr/maven"
)

var storage = flag.String("storage", "filesystem", "The storage engine to use (defaults to \"filesystem\").")
var repository *maven.Repository

func main() {
	flag.Parse()

	port := fromEnv("PORT", "4040")

	repository = maven.NewRepository(*storage)

	webserver := gin.Default()

	// maven repo urls
	webserver.PUT("/maven/*any", Upload)
	webserver.HEAD("/maven/*any", Exists)
	webserver.GET("/maven/*any", Download)

	// api urls
	webserver.POST("/api/register", Register)
	webserver.POST("/api/login", Login)
	webserver.GET("/api/username/:username", UsernameExists)
	webserver.GET("/api/domain/:domain", DomainExists)
	webserver.GET("/api/packages", Packages)
	webserver.POST("/api/packages", UpdatePackage)
	webserver.GET("/api/search", Search)
	// logout will simply be a reload of the single page.

	webserver.StaticFile("/", "./static/index.html")
	webserver.Static("/static", "./static")

	webserver.Run(fmt.Sprintf(":%s", port))
}

func fromEnv(param string, defaultVal string) string {
	env := os.Getenv(param)

	if env == "" {
		env = defaultVal
	}

	return env
}