package main

import (
	"github.com/gin-gonic/gin"
	"fmt"
)

func main() {

	webserver := gin.Default()

/*
	webserver.GET("/maven/*binary", func(c *gin.Context) {
		all := c.Param("binary")

		split := strings.Split(all[1:], "/")

		binary := split[len(split) - 1]
		version := split[len(split) - 2]
		artifact := split[len(split) - 3]
		group := strings.Join(split[0:len(split) - 3], "/")

		c.String(200, "%s %s %s %s", group, artifact, version, binary)
	})
*/

	webserver.Any("/*any", func(g *gin.Context) {
		url := g.Param("any")
		method := g.Request.Method
		fmt.Printf("%s %s", method, url)

		if method == "HEAD" {
			g.String(200, "")
		} else {
			g.Next()
		}
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
