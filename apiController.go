package main

import (
	"github.com/gin-gonic/gin"
	"time"
	"strings"
	"github.com/Meduzz/yamr/users"
	"github.com/Meduzz/yamr/sessions"
	"github.com/Meduzz/yamr/packages"
	"github.com/Meduzz/yamr/maven"
	"github.com/Meduzz/yamr/artifacts"
	"net/http"
	"strconv"
	"github.com/Meduzz/yamr/domains"
)

var (
	sessionManager = sessions.NewSessions()
	userManager = users.NewUsers()
	packageManager = packages.NewPackages()
	artifactManager = artifacts.NewArtifacts()
	domainManager = domains.NewDomains()
)

// register a user (incl top domain (se.kodiak)).
func Register(g *gin.Context) {
	u := &users.User{}
	err := g.BindJSON(u)

	if err != nil {
		g.AbortWithError(500, err)
	} else {
		err := userManager.Store(u)

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

	ip := cleanIp(g.Request)
	u, err := userManager.LoadByUsernameAndPassword(credential.Username, credential.Password)

	if err != nil {
		g.AbortWithError(404, err)
	} else {
		session, err := sessionManager.CreateForUser(u.Id, ip)

		if err != nil {
			g.AbortWithError(500, err)
		} else {
			g.JSON(200, gin.H{"Id":session.Id, "Admin":u.Admin})
		}
	}
}

// does the username already exists?
func UsernameExists(g *gin.Context) {
	username := g.Param("username")

	exists, err := userManager.UserExists(username)

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

func ApplyForDomain(g *gin.Context) {
	sessionId := g.Request.Header.Get("Session")
	ip := cleanIp(g.Request)
	domain := &domains.Domain{}

	session, err := sessionManager.LoadById(sessionId)

	if err != nil {
		g.AbortWithError(500, err)
	} else if !valid(session, ip) {
		g.AbortWithStatus(403)
	} else {
		sessionManager.Extend(session)
		err = g.BindJSON(domain)
		if err != nil {
			g.AbortWithError(500, err)
		} else {
			err := domainManager.Create(domain, session.UserId)

			if err != nil {
				g.AbortWithError(500, err)
			} else {
				g.JSON(201, "")
			}
		}
	}
}

func Domains(g *gin.Context) {
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
	} else if !valid(session, ip) {
		g.AbortWithStatus(403)
	} else {
		sessionManager.Extend(session)
		usersDomains, err := domainManager.ListDomainsForUser(session.UserId, page, limit)

		if err != nil {
			g.AbortWithError(500, err)
		} else {
			g.JSON(200, usersDomains)
		}
	}
}

// list the users packages.
func Packages(g *gin.Context) {
	sessionId := g.Request.Header.Get("Session")
	ip := cleanIp(g.Request)
	sPage := g.Query("skip")
	sLimit := g.Query("limit")
	sDomainId := g.Query("id")

	page := 0
	limit := 20

	if sPage != "" {
		page, _ = strconv.Atoi(sPage)
	}

	if sLimit != "" {
		limit, _ = strconv.Atoi(sLimit)
	}

	domainId, _ := strconv.Atoi(sDomainId)

	session, err := sessionManager.LoadById(sessionId)

	if err != nil {
		g.AbortWithError(500, err)
	} else if !valid(session, ip) {
		g.AbortWithStatus(403)
	} else {
		if !domainManager.OwnedBy(int64(domainId), session.UserId) {
			g.JSON(404, gin.H{})
		} else {

			sessionManager.Extend(session)
			ps, err := packageManager.List(int64(domainId), page, limit)

			if err != nil {
				g.AbortWithError(500, err)
			} else {
				g.JSON(200, ps)
			}
		}
	}
}

// update/create a package.
func UpdatePackage(g *gin.Context) {
	sessionId := g.Request.Header.Get("Session")
	ip := cleanIp(g.Request)
	sDomainId := g.Query("id")

	domainId, _ := strconv.Atoi(sDomainId)

	session, err := sessionManager.LoadById(sessionId)

	if err != nil {
		g.AbortWithError(500, err)
	} else if !valid(session, ip) {
		g.AbortWithStatus(403)
	} else {
		sessionManager.Extend(session)

		p := &packages.Package{}
		err = g.BindJSON(p)

		if err != nil {
			g.AbortWithError(400, err)
		} else {
			if !domainManager.OwnedBy(int64(domainId), session.UserId) {
				g.JSON(401, gin.H{})
			} else {
				err = packageManager.UpdateOrCreate(int64(domainId), p)

				if err != nil {
					g.AbortWithError(400, err)
				} else {
					g.JSON(200, gin.H{})
				}
			}
		}
	}
}

// handle queries for packages.
func Search(g *gin.Context) {
	sessionId := g.Request.Header.Get("Session")
	query := g.Query("q")
	sPage := g.Query("page")
	sLimit := g.Query("limit")

	page := 0
	limit := 20

	if sPage != "" {
		page, _ = strconv.Atoi(sPage)
	}

	if sLimit != "" {
		limit, _ = strconv.Atoi(sLimit)
	}

	if len(sessionId) > 0 {
		ip := cleanIp(g.Request)

		session, err := sessionManager.LoadById(sessionId)

		if err != nil {
			// search without user.
			result, err := artifactManager.Search(query, 0, page, limit)
			if err != nil {
				g.AbortWithError(500, err)
			} else {
				g.JSON(200, result)
			}
		} else if !valid(session, ip) {
			// search without user.
			result, err := artifactManager.Search(query, 0, page, limit)
			if err != nil {
				g.AbortWithError(500, err)
			} else {
				g.JSON(200, result)
			}
		} else {
			sessionManager.Extend(session)
			// search with user.
			// session.Package
			result, err := artifactManager.Search(query, session.UserId, page, limit)
			if err != nil {
				g.AbortWithError(500, err)
			} else {
				g.JSON(200, result)
			}
		}
	} else {
		// search without user.
		result, err := artifactManager.Search(query, 0, page, limit)
		if err != nil {
			g.AbortWithError(500, err)
		} else {
			g.JSON(200, result)
		}
	}
}

func valid(session *sessions.Session, ip string) bool {
	now := time.Now()
	return session.Expires.After(now) && session.Ip == ip
}

func cleanIp(req *http.Request) string {
	proxied := req.Header.Get("X-Forwarded-For")

	if len(proxied) == 0 {
		ip := strings.Replace(req.RemoteAddr, "[::1]", "127.0.0.1", -1)
		return strings.Split(ip, ":")[0]
	} else  {
		return strings.Split(proxied, ":")[0]
	}
}