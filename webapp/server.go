package main

import (
	"crypto/rand"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

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
	router.LoadHTMLFiles(joinPath(dir(self), "templates/index.html"))
	router.GET("/", func(c *gin.Context) {
		username := sessions.Default(c).Get("username")
		c.HTML(200, "index.html", gin.H{"user": username, "refresh": refresh})
	})

	auth := router.Group("/auth")
	auth.POST("/login", login)
	auth.GET("/logout", authRequired, func(c *gin.Context) {
		session := sessions.Default(c)
		session.Clear()
		session.Save()
		c.Redirect(302, "/")
	})
	auth.POST("/setting", authRequired, setting)

	base := router.Group("/")
	base.GET("/stock/:index/:code", func(c *gin.Context) {
		username := sessions.Default(c).Get("username")
		c.HTML(200, "index.html", gin.H{"user": username, "refresh": refresh})
	})
	base.GET("/mystocks", myStocks)
	base.GET("/indices", indices)
	base.POST("/get", getStock)
	base.POST("/suggest", getSuggest)
	base.GET("/star", star)
	base.POST("/star", doStar)
	base.POST("/reorder", reorder)

	router.NoRoute(func(c *gin.Context) {
		c.Redirect(302, "/")
	})

	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
