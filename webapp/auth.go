package main

import (
	"log"
	"strings"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type user struct {
	ID       int
	Username string
	Password string
}

func authRequired(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user_id")
	if userID == nil || userID == 0 {
		c.Abort()
		c.Redirect(302, "/")
	}
}

func login(c *gin.Context) {
	db, err := getDB()
	if err != nil {
		log.Println("Failed to connect to database:", err)
		c.HTML(200, "login.html", gin.H{"error": "Failed to connect to database."})
		return
	}
	defer db.Close()

	session := sessions.Default(c)
	username := strings.TrimSpace(strings.ToLower(c.PostForm("username")))
	password := c.PostForm("password")

	var user user
	var message string
	if err := db.QueryRow(
		"SELECT id, username, password FROM user WHERE username = ?", username).Scan(&user.ID, &user.Username, &user.Password); err != nil {
		if strings.Contains(err.Error(), "doesn't exist") {
			restore("")
			c.HTML(200, "login.html", gin.H{"error": "Detected first time running. Initialized the database."})
			return
		}
		if strings.Contains(err.Error(), "no rows") {
			message = "Incorrect username"
		} else {
			log.Println(err)
			message = "Critical Error! Please contact your system administrator."
		}
	} else {
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
			if (strings.Contains(err.Error(), "too short") && user.Password != password) || strings.Contains(err.Error(), "is not") {
				message = "Incorrect password"
			} else if user.Password != password {
				log.Println(err)
				message = "Critical Error! Please contact your system administrator."
			}
		}
		if message == "" {
			session.Clear()
			session.Set("user_id", user.ID)
			session.Set("username", user.Username)

			rememberme := c.PostForm("rememberme")
			if rememberme == "on" {
				session.Options(sessions.Options{Path: "/", HttpOnly: true, MaxAge: 856400 * 365})
			} else {
				session.Options(sessions.Options{Path: "/", HttpOnly: true, MaxAge: 0})
			}

			if err := session.Save(); err != nil {
				log.Println(err)
				c.HTML(200, "login.html", gin.H{"error": "Failed to save session."})
				return
			}
			c.Redirect(302, "/")
			return
		}
	}
	c.HTML(200, "login.html", gin.H{"error": message})
}

func setting(c *gin.Context) {
	db, err := getDB()
	if err != nil {
		log.Println("Failed to connect to database:", err)
		c.HTML(200, "setting.html", gin.H{"error": "Failed to connect to database."})
		return
	}
	defer db.Close()

	session := sessions.Default(c)
	userID := session.Get("user_id")

	password := c.PostForm("password")
	password1 := c.PostForm("password1")
	password2 := c.PostForm("password2")

	var oldPassword string
	if err := db.QueryRow("SELECT password FROM user WHERE id = ?", userID).Scan(&oldPassword); err != nil {
		log.Println(err)
		c.HTML(200, "setting.html", gin.H{"error": "Failed to get current user password."})
		return
	}

	var message string
	err = bcrypt.CompareHashAndPassword([]byte(oldPassword), []byte(password))
	switch {
	case err != nil:
		if (strings.Contains(err.Error(), "too short") && password != oldPassword) || strings.Contains(err.Error(), "is not") {
			message = "Incorrect password."
		} else if password != oldPassword {
			log.Println(err)
			message = "Failed to compare password."
		}
	case password1 != password2:
		message = "Confirm password doesn't match new password."
	case password1 == password:
		message = "New password cannot be the same as your current password."
	case password1 == "":
		message = "New password cannot be blank."
	}

	if message == "" {
		newPassword, err := bcrypt.GenerateFromPassword([]byte(password1), bcrypt.MinCost)
		if err != nil {
			log.Println(err)
			c.HTML(200, "setting.html", gin.H{"error": "Failed to encrypt new password."})
			return
		}
		if _, err := db.Exec("UPDATE user SET password = ? WHERE id = ?", string(newPassword), userID); err != nil {
			log.Println(err)
			c.HTML(200, "setting.html", gin.H{"error": "Failed to update new password."})
			return
		}
		session.Clear()
		if err := session.Save(); err != nil {
			log.Println(err)
			c.HTML(200, "setting.html", gin.H{"error": "Failed to save session."})
			return
		}
		c.Redirect(302, "/")
		return
	}
	c.HTML(200, "setting.html", gin.H{"error": message})
}
