package main

import (
	"flag"
	"fmt"
	"os"

	"./maven"
	"github.com/gin-gonic/gin"
)

var storage = flag.String("storage", "filesystem", "The storage engine to use (defaults to \"filesystem\").")
var repository *maven.Repository

func main() {
	flag.Parse()

	port := fromEnv("PORT", "4040")

	repository = maven.NewRepository(*storage)

	webserver := gin.Default()

	// maven repo urls
	webserver.PUT("/maven/*any", upload)
	webserver.HEAD("/maven/*any", exists)
	webserver.GET("/maven/*any", download)

	// api urls
	webserver.POST("/api/register", register)
	webserver.POST("/api/login", login)
	webserver.GET("/api/username/:username", usernameExists)
	webserver.PUT("/api/domain/apply", applyForDomain)
	webserver.GET("/api/domains", loadDomains)
	webserver.GET("/api/search", search)
	webserver.GET("/api/packages", loadPackages)
	webserver.POST("/api/packages", updatePackage)
	// logout will simply be a reload of the single page.

	// admin api urls
	webserver.GET("/admin/domains", findDomains)
	webserver.GET("/admin/activate/:domain", activateDomain)
	webserver.GET("/admin/users", listUsers)

	// static files
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
