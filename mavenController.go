package main

import (
	"github.com/gin-gonic/gin"
	"strings"
	"github.com/Meduzz/yamr/maven"
)

func Upload(g *gin.Context) {
	context := repository.Context(g)
	err := repository.Write(context, g.Request.Body)

	if err != nil {
		if isAccessDenied(err) {
			meta := context.Get(maven.FILEMETADATA).(*maven.FileMetadata)
			g.Header("WWW-Authenticate", "Basic realm=\"" + meta.GroupAsPackage() + "\"")
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
			meta := context.Get(maven.FILEMETADATA).(*maven.FileMetadata)
			g.Header("WWW-Authenticate", "Basic realm=\"" + meta.GroupAsPackage() + "\"")
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
		if isAccessDenied(err) {
			meta := context.Get(maven.FILEMETADATA).(*maven.FileMetadata)
			g.Header("WWW-Authenticate", "Basic realm=\"" + meta.GroupAsPackage() + "\"")
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
	return strings.Contains(err.Error(), "denied") ||
			strings.Contains(err.Error(), "credential")
}