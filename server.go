package main

import (
	"crypto/rand"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

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

	if *unix != "" && OS == "linux" {
		if _, err := os.Stat(*unix); err == nil {
			err = os.Remove(*unix)
			if err != nil {
				log.Fatalln("Failed to remove socket file:", err)
			}
		}

		listener, err := net.Listen("unix", *unix)
		if err != nil {
			log.Fatalln("Failed to listen socket file:", err)
		}

		idleConnsClosed := make(chan struct{})
		go func() {
			quit := make(chan os.Signal, 1)
			signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
			<-quit

			if err := listener.Close(); err != nil {
				log.Println("Failed to close listener:", err)
			}
			if _, err := os.Stat(*unix); err == nil {
				if err := os.Remove(*unix); err != nil {
					log.Println("Failed to remove socket file:", err)
				}
			}
			close(idleConnsClosed)
		}()

		if err := os.Chmod(*unix, 0666); err != nil {
			log.Fatalln("Failed to chmod socket file:", err)
		}

		http.Serve(listener, router)
		<-idleConnsClosed
	} else {
		router.Run(*host + ":" + *port)
	}
}
