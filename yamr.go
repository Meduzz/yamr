package main

import (
	"github.com/gin-gonic/gin"
	"os"
	"fmt"
	"flag"
	"./maven"
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
	webserver.GET("/api/profile", Profile)
	webserver.POST("/api/profile", UpdateProfile)
	// TODO logout will simply be a reload of the single page.

	// TODO add a static url.

	webserver.Run(fmt.Sprintf(":%s", port))
}

func fromEnv(param string, defaultVal string) string {
	env := os.Getenv(param)

	if env == "" {
		env = defaultVal
	}

	return env
}
