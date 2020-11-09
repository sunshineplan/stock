package main

import (
	"crypto/rand"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

func loadTemplates() multitemplate.Renderer {
	r := multitemplate.NewRenderer()

	includes, err := filepath.Glob(joinPath(dir(self), "templates/*/*"))
	if err != nil {
		log.Fatalln("Failed to glob stock templates:", err)
	}

	for _, include := range includes {
		r.AddFromFiles(filepath.Base(include), joinPath(dir(self), "templates/base.html"), include)
	}
	return r
}

func run() {
	if *logPath != "" {
		f, err := os.OpenFile(*logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0640)
		if err != nil {
			log.Fatalln("Failed to open log file:", err)
		}
		gin.DefaultWriter = f
		gin.DefaultErrorWriter = f
		log.SetOutput(f)
	}

	secret := make([]byte, 16)
	if _, err := rand.Read(secret); err != nil {
		log.Fatalln("Failed to get secret:", err)
	}

	router := gin.Default()
	server.Handler = router
	router.Use(sessions.Sessions("session", sessions.NewCookieStore(secret)))
	router.StaticFS("/static", http.Dir(joinPath(dir(self), "static")))
	router.HTMLRender = loadTemplates()
	router.GET("/", func(c *gin.Context) {
		session := sessions.Default(c)
		username := session.Get("username")
		c.HTML(200, "mystocks.html", gin.H{"user": username, "refresh": refresh})
	})

	auth := router.Group("/auth")
	auth.GET("/login", func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get("user_id")
		if user != nil && user != 0 {
			c.Redirect(302, "/")
			return
		}
		c.HTML(200, "login.html", gin.H{"error": ""})
	})
	auth.GET("/logout", authRequired, func(c *gin.Context) {
		session := sessions.Default(c)
		session.Clear()
		session.Save()
		c.Redirect(302, "/")
	})
	auth.POST("/login", login)
	auth.GET("/setting", authRequired, func(c *gin.Context) {
		session := sessions.Default(c)
		username := session.Get("username")
		c.HTML(200, "setting.html", gin.H{"user": username, "error": ""})
	})
	auth.POST("/setting", authRequired, setting)

	base := router.Group("/")
	base.GET("/stock/:index/:code", showStock)
	base.GET("/mystocks", myStocks)
	base.GET("/indices", indices)
	base.GET("/get", getStock)
	base.GET("/suggest", getSuggest)
	base.GET("/star", star)
	base.POST("/star", doStar)
	base.POST("/reorder", reorder)

	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
