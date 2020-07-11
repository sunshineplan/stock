package main

import (
	"crypto/rand"
	"io"
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
		log.Fatalf("Failed to glob stock templates: %v", err)
	}

	for _, include := range includes {
		r.AddFromFiles(filepath.Base(include), joinPath(dir(self), "templates/base.html"), include)
	}
	return r
}

func run() {
	f, err := os.OpenFile(*logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0640)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	gin.DefaultWriter = io.MultiWriter(f)
	log.SetOutput(gin.DefaultWriter)

	secret := make([]byte, 16)
	_, err = rand.Read(secret)
	if err != nil {
		log.Fatalf("Failed to get secret: %v", err)
	}

	router := gin.Default()
	router.Use(sessions.Sessions("session", sessions.NewCookieStore(secret)))
	router.StaticFS("/static", http.Dir(joinPath(dir(self), "static")))
	router.HTMLRender = loadTemplates()
	router.GET("/", func(c *gin.Context) {
		session := sessions.Default(c)
		username := session.Get("username")
		c.HTML(200, "mystocks.html", gin.H{"user": username})
	})

	auth := router.Group("/auth")
	auth.GET("/login", func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get("user_id")
		if user != nil {
			c.Redirect(302, "/")
			return
		}
		c.HTML(200, "login.html", gin.H{"error": ""})
	})
	auth.POST("/login", login)
	auth.GET("/logout", authRequired, func(c *gin.Context) {
		session := sessions.Default(c)
		session.Clear()
		session.Save()
		c.Redirect(302, "/auth/login")
	})
	auth.GET("/setting", authRequired, func(c *gin.Context) {
		c.HTML(200, "setting.html", nil)
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
				log.Fatalf("Failed to remove socket file: %v", err)
			}
		}

		listener, err := net.Listen("unix", *unix)
		if err != nil {
			log.Fatalf("Failed to listen socket file: %v", err)
		}

		idleConnsClosed := make(chan struct{})
		go func() {
			quit := make(chan os.Signal, 1)
			signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
			<-quit

			if err := listener.Close(); err != nil {
				log.Printf("Failed to close listener: %v", err)
			}
			if _, err := os.Stat(*unix); err == nil {
				err = os.Remove(*unix)
				if err != nil {
					log.Printf("Failed to remove socket file: %v", err)
				}
			}
			close(idleConnsClosed)
		}()

		if err = os.Chmod(*unix, 0666); err != nil {
			log.Fatalf("Failed to chmod socket file: %v", err)
		}

		http.Serve(listener, router)
		<-idleConnsClosed
	} else {
		router.Run(*host + ":" + *port)
	}
}
