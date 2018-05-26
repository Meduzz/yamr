package main

import (
	"strings"

	"github.com/Meduzz/yamr/artifacts"
	"github.com/gin-gonic/gin"
)

const (
	USER     = "user"
	SESSION  = "session"
	FILEMETA = "filemeta"
)

func authenticate(ctx *gin.Context) {
	username, password, ok := ctx.Request.BasicAuth()

	if ok {
		user, err := userManager.LoadByUsernameAndPassword(username, password)

		if err != nil {
			// TODO WWW-Authenticate?
			ctx.AbortWithError(401, err)
		} else {
			ctx.Set(USER, user)
			ctx.Next()
		}
	} else {
		sessionId := ctx.GetHeader("Session")
		ip := cleanIp(ctx.Request)

		if sessionId == "" {
			ctx.AbortWithStatus(401)
		} else {
			session, err := sessionManager.LoadById(sessionId)

			if err != nil {
				ctx.AbortWithError(500, err)
			} else {
				if !valid(session, ip) {
					ctx.AbortWithStatus(403)
				} else {
					ctx.Set(SESSION, session)

					user, err := userManager.LoadById(session.UserId)

					if err != nil {
						ctx.AbortWithError(500, err)
					} else {
						ctx.Set(USER, user)
						ctx.Next()
					}
				}
			}
		}
	}
}

func loadPackage(ctx *gin.Context) {
	any := ctx.Param("any")

	ctx.Set(FILEMETA, extract(any))
	ctx.Next()
}

func extract(path string) *artifacts.FileMetadata {
	split := strings.Split(path[1:], "/")

	binary := split[len(split)-1]
	version := split[len(split)-2]
	artifact := split[len(split)-3]
	group := split[0 : len(split)-3]

	return artifacts.NewFileMetadata(group, artifact, version, binary)
}
