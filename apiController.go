package main

import (
	"github.com/gin-gonic/gin"
	"time"
	"strings"
	"github.com/Meduzz/yamr/user"
	"github.com/Meduzz/yamr/maven"
)

// register a user (incl top domain (se.kodiak)).
func Register(g *gin.Context) {
	u := &user.User{}
	err := g.BindJSON(u)

	if err != nil {
		g.AbortWithError(500, err)
	} else {
		err := user.NewUsers().Store(u)

		if err != nil {
			g.AbortWithError(500, err)
		} else {
			g.JSON(200, "")
		}
	}
}

// login a user.
func Login(g *gin.Context) {
	credential := &maven.Credential{}
	err := g.BindJSON(credential)

	if err != nil {
		g.AbortWithError(500, err)
	}

	ip := cleanIp(g.Request.RemoteAddr)
	u, err := user.NewUsers().LoadByUsernameAndPassword(credential.Username, credential.Password)

	if err != nil {
		g.AbortWithError(500, err)
	} else {
		session, err := user.NewSessions().CreateForUser(u.Package, ip)

		if err != nil {
			g.AbortWithError(500, err)
		} else {
			g.JSON(200, gin.H{"Id":session.Id})
		}
	}
}

// does the username already exists?
func UsernameExists(g *gin.Context) {
	username := g.Param("username")

	exists, err := user.NewUsers().UserExists(username)

	if err != nil {
		g.AbortWithError(500, err)
	} else {
		if exists {
			g.JSON(400, "")
		} else {
			g.JSON(200, "")
		}
	}
}

// are the top domain already registered?
func DomainExists(g *gin.Context) {
	domain := g.Param("domain")

	exists, err := user.NewUsers().DomainExists(domain)

	if err != nil {
		g.AbortWithError(500, err)
	} else {
		if exists {
			g.JSON(400, "")
		} else {
			g.JSON(200, "")
		}
	}
}

// list the users packages.
func Packages(g *gin.Context) {
	sessionId := g.Request.Header.Get("Session")
	ip := cleanIp(g.Request.RemoteAddr)

	session, err := user.NewSessions().LoadById(sessionId)

	if err != nil {
		g.AbortWithError(500, err)
	} else if !valid(session, ip) {
		g.AbortWithStatus(403)
	} else {
		user.NewSessions().Extend(session)
		ps, err := user.NewPackages().List(session.Package)

		if err != nil {
			g.AbortWithError(500, err)
		} else  {
			g.JSON(200, ps)
		}
	}
}

// update/create a package.
func UpdatePackage(g *gin.Context) {
	sessionId := g.Request.Header.Get("Session")
	ip := cleanIp(g.Request.RemoteAddr)

	session, err := user.NewSessions().LoadById(sessionId)

	if err != nil {
		g.AbortWithError(500, err)
	} else if !valid(session, ip) {
		g.AbortWithStatus(403)
	} else {
		user.NewSessions().Extend(session)
		p := &user.Package{}
		err = g.BindJSON(p)

		if err != nil {
			g.AbortWithStatus(400)
		} else {
			err := user.NewPackages().UpdateOrCreate(p)

			if err != nil {
				g.AbortWithStatus(400)
			} else {
				g.JSON(200, gin.H{})
			}
		}
	}
}

func valid(session *user.Session, ip string) bool {
	now := time.Now()
	return session.Expires.After(now) && session.Ip == ip
}

func cleanIp(ip string) string {
	return strings.Split(ip, ":")[0]
}