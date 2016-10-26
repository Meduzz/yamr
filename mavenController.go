package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"strings"
)

func Upload(g *gin.Context) {
	context := repository.Context(g)
	err := repository.Write(context, g.Request.Body)

	if err != nil {
		log.Printf("There was an error. %s", err)
		if isAccessDenied(err) {
			// TODO figure out a good realm.
			g.Header("WWW-Authenticate", "Basic realm=\"test\"")
			g.AbortWithError(401, err)
		} else {
			g.AbortWithError(500, err)
		}
	} else {
		g.Next()
	}
}

func Exists(g *gin.Context) {
	context := repository.Context(g)
	exists, err := repository.Exists(context)

	if err != nil {
		if isAccessDenied(err) {
			// TODO figure out a good realm.
			g.Header("WWW-Authenticate", "Basic realm=\"test\"")
			g.AbortWithError(401, err)
		} else {
			g.AbortWithError(500, err)
		}
	} else if exists {
		g.String(200, "")
	} else {
		// 404?
		g.String(404, "")
	}
}

func Download(g *gin.Context) {
	context := repository.Context(g)
	bytes, err := repository.Read(context)

	if err != nil {
		log.Printf("There was an error. %s", err)
		if isAccessDenied(err) {
			// TODO figure out a good realm.
			g.Header("WWW-Authenticate", "Basic realm=\"test\"")
			g.AbortWithError(401, err)
		} else {
			g.AbortWithError(500, err)
		}
	} else {
		g.Status(200)
		g.Writer.Write(bytes)
	}
}

func isAccessDenied(err error) bool {
	return strings.Contains(err.Error(), "Access denied")
}