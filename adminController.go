package main

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

func FindDomains(g *gin.Context) {
	sessionId := g.Request.Header.Get("Session")
	ip := cleanIp(g.Request)
	sPage := g.Query("skip")
	sLimit := g.Query("limit")

	page := 0
	limit := 20

	if sPage != "" {
		page, _ = strconv.Atoi(sPage)
	}

	if sLimit != "" {
		limit, _ = strconv.Atoi(sLimit)
	}

	session, err := sessionManager.LoadById(sessionId)

	if err != nil {
		g.AbortWithError(500, err)
	}

	u, err := userManager.LoadById(session.UserId)

	if err != nil {
		g.AbortWithError(500, err)
	} else if !valid(session, ip) || !u.Admin {
		g.AbortWithStatus(403)
	} else {
		sessionManager.Extend(session)

		inactives, err := domainManager.FindInactiveDomains(page, limit)

		if err != nil {
			g.AbortWithError(500, err)
		} else {
			g.JSON(200, inactives)
		}
	}
}

func ActivateDomain(g *gin.Context) {
	sessionId := g.Request.Header.Get("Session")
	ip := cleanIp(g.Request)
	sDomainId := g.Param("domain")

	iDomainId, err := strconv.Atoi(sDomainId)

	session, err := sessionManager.LoadById(sessionId)

	if err != nil {
		g.AbortWithError(500, err)
	}

	u, err := userManager.LoadById(session.UserId)

	if err != nil {
		g.AbortWithError(500, err)
	} else if !valid(session, ip) || !u.Admin {
		g.AbortWithStatus(403)
	} else {
		sessionManager.Extend(session)

		err = domainManager.Activate(int64(iDomainId))

		if err != nil {
			g.AbortWithError(500, err)
		} else {
			g.JSON(200, gin.H{})
		}
	}
}

func ListUsers(g *gin.Context) {
	sessionId := g.Request.Header.Get("Session")
	ip := cleanIp(g.Request)
	sPage := g.Query("skip")
	sLimit := g.Query("limit")

	page := 0
	limit := 20

	if sPage != "" {
		page, _ = strconv.Atoi(sPage)
	}

	if sLimit != "" {
		limit, _ = strconv.Atoi(sLimit)
	}

	session, err := sessionManager.LoadById(sessionId)

	if err != nil {
		g.AbortWithError(500, err)
	}

	u, err := userManager.LoadById(session.UserId)

	if err != nil {
		g.AbortWithError(500, err)
	} else if !valid(session, ip) || !u.Admin {
		g.AbortWithStatus(403)
	} else {
		sessionManager.Extend(session)

		data, err := userManager.ListUsers(page, limit)

		if err != nil {
			g.AbortWithError(500, err)
		} else {
			g.JSON(200, data)
		}
	}
}