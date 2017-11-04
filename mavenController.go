package main

import (
	"strings"

	"github.com/Meduzz/yamr/artifacts"
	"github.com/Meduzz/yamr/maven"
	"github.com/gin-gonic/gin"
)

func upload(g *gin.Context) {
	context := repository.Context(g)
	err := repository.Write(context, g.Request.Body)

	if err != nil {
		if isAccessDenied(err) {
			meta := context.Get(maven.FILEMETADATA).(*artifacts.FileMetadata)
			g.Header("WWW-Authenticate", "Basic realm=\""+meta.GroupAsPackage()+"\"")
			g.AbortWithError(401, err)
		} else {
			g.AbortWithError(500, err)
		}
	} else {
		g.Next()
	}
}

func exists(g *gin.Context) {
	context := repository.Context(g)
	exists, err := repository.Exists(context)

	if err != nil {
		if isAccessDenied(err) {
			meta := context.Get(maven.FILEMETADATA).(*artifacts.FileMetadata)
			g.Header("WWW-Authenticate", "Basic realm=\""+meta.GroupAsPackage()+"\"")
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

func download(g *gin.Context) {
	context := repository.Context(g)
	bytes, err := repository.Read(context)

	if err != nil {
		if isAccessDenied(err) {
			meta := context.Get(maven.FILEMETADATA).(*artifacts.FileMetadata)
			g.Header("WWW-Authenticate", "Basic realm=\""+meta.GroupAsPackage()+"\"")
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
