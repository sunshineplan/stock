package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

func myStocks(c *gin.Context) {
	db, err := sql.Open("sqlite3", sqlite)
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		c.String(503, "")
		return
	}
	defer db.Close()
	session := sessions.Default(c)
	userID := session.Get("user_id")

	if userID == nil {
		userID = 0
	}

	rows, err := db.Query(`SELECT idx, code FROM stock WHERE user_id = ? ORDER BY seq`, userID)
	if err != nil {
		log.Printf("Failed to get all stocks: %v", err)
		c.String(500, "")
		return
	}
	defer rows.Close()
	var stocks []stock
	for rows.Next() {
		var index, code string
		if err := rows.Scan(&index, &code); err != nil {
			log.Printf("Failed to scan all stocks: %v", err)
			c.String(500, "")
			return
		}
		stock := initStock(index, code)
		stocks = append(stocks, stock)
	}
	doGetRealtimes(stocks)

	c.JSON(200, "index.html", gin.H{"category": category, "stocks": stocks})
}

func addstock(c *gin.Context) {
	db, err := sql.Open("sqlite3", sqlite)
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		c.String(503, "")
		return
	}
	defer db.Close()
	session := sessions.Default(c)
	userID := session.Get("user_id")

	var stock stock
	var categories []string
	categoryID, err := strconv.Atoi(c.Query("category"))
	if err == nil {
		db.QueryRow("SELECT category FROM category WHERE id = ? AND user_id = ?", categoryID, userID).Scan(&stock.Category)
	} else {
		stock.Category = []byte("")
	}
	rows, err := db.Query("SELECT category FROM category WHERE user_id = ? ORDER BY category", userID)
	if err != nil {
		log.Printf("Failed to get categories name: %v", err)
		c.String(500, "")
		return
	}
	defer rows.Close()
	for rows.Next() {
		var categoryName string
		if err := rows.Scan(&categoryName); err != nil {
			log.Printf("Failed to scan category name: %v", err)
			c.String(500, "")
			return
		}
		categories = append(categories, categoryName)
	}
	c.HTML(200, "stock.html", gin.H{"id": 0, "stock": stock, "categories": categories})
}

func doAddstock(c *gin.Context) {
	db, err := sql.Open("sqlite3", sqlite)
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		c.String(503, "")
		return
	}
	defer db.Close()
	session := sessions.Default(c)
	userID := session.Get("user_id")
	category := strings.TrimSpace(c.PostForm("category"))
	stock := strings.TrimSpace(c.PostForm("stock"))
	url := strings.TrimSpace(c.PostForm("url"))
	categoryID, err := getCategoryID(category, userID.(int), db)
	if err != nil {
		log.Printf("Failed to get category id: %v", err)
		c.String(500, "")
		return
	}

	var exist, message string
	var errorCode int
	if stock == "" {
		message = "stock name is empty."
		errorCode = 1
	} else if err = db.QueryRow("SELECT id FROM stock WHERE stock = ? AND user_id = ?", stock, userID).Scan(&exist); err == nil {
		message = fmt.Sprintf("stock name %s is already existed.", stock)
		errorCode = 1
	} else if err = db.QueryRow("SELECT id FROM stock WHERE url = ? AND user_id = ?", url, userID).Scan(&exist); err == nil {
		message = fmt.Sprintf("stock url %s is already existed.", url)
		errorCode = 2
	} else if categoryID == -1 {
		message = "Category name exceeded length limit."
		errorCode = 3
	} else {
		_, err = db.Exec("INSERT INTO stock (stock, url, user_id, category_id) VALUES (?, ?, ?, ?)", stock, url, userID, categoryID)
		if err != nil {
			log.Printf("Failed to add stock: %v", err)
			c.String(500, "")
			return
		}
		c.JSON(200, gin.H{"status": 1})
		return
	}
	c.JSON(200, gin.H{"status": 0, "message": message, "error": errorCode})
}

func editstock(c *gin.Context) {
	db, err := sql.Open("sqlite3", sqlite)
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		c.String(503, "")
		return
	}
	defer db.Close()
	session := sessions.Default(c)
	userID := session.Get("user_id")
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Printf("Failed to get id param: %v", err)
		c.String(400, "")
		return
	}

	var stock stock
	err = db.QueryRow(`
SELECT stock, url, category FROM stock
LEFT JOIN category ON category_id = category.id
WHERE stock.id = ? AND stock.user_id = ?
`, id, userID).Scan(&stock.Name, &stock.URL, &stock.Category)
	if err != nil {
		log.Printf("Failed to scan stock: %v", err)
		c.String(500, "")
		return
	}

	var categories []string
	rows, err := db.Query("SELECT category FROM category WHERE user_id = ? ORDER BY category", userID)
	if err != nil {
		log.Printf("Failed to get categories: %v", err)
		c.String(500, "")
		return
	}
	defer rows.Close()
	for rows.Next() {
		var category string
		if err := rows.Scan(&category); err != nil {
			log.Printf("Failed to scan category: %v", err)
			c.String(500, "")
			return
		}
		categories = append(categories, category)
	}
	c.HTML(200, "stock.html", gin.H{"id": id, "stock": stock, "categories": categories})
}

func doEditstock(c *gin.Context) {
	db, err := sql.Open("sqlite3", sqlite)
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		c.String(503, "")
		return
	}
	defer db.Close()
	session := sessions.Default(c)
	userID := session.Get("user_id")
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Printf("Failed to get id param: %v", err)
		c.String(400, "")
		return
	}

	var old stock
	err = db.QueryRow(`
SELECT stock, url, category FROM stock
LEFT JOIN category ON category_id = category.id
WHERE stock.id = ? AND stock.user_id = ?
`, id, userID).Scan(&old.Name, &old.URL, &old.Category)
	if err != nil {
		log.Printf("Failed to scan stock: %v", err)
		c.String(500, "")
		return
	}
	stock := strings.TrimSpace(c.PostForm("stock"))
	url := strings.TrimSpace(c.PostForm("url"))
	category := strings.TrimSpace(c.PostForm("category"))
	categoryID, err := getCategoryID(category, userID.(int), db)
	if err != nil {
		log.Printf("Failed to get category id: %v", err)
		c.String(500, "")
		return
	}

	var exist, message string
	var errorCode int
	if stock == "" {
		message = "stock name is empty."
		errorCode = 1
	} else if old.Name == stock && old.URL == url && string(old.Category) == category {
		message = "New stock is same as old stock."
	} else if err = db.QueryRow("SELECT id FROM stock WHERE stock = ? AND id != ? AND user_id = ?", stock, id, userID).Scan(&exist); err == nil {
		message = fmt.Sprintf("stock name %s is already existed.", stock)
		errorCode = 1
	} else if err = db.QueryRow("SELECT id FROM stock WHERE url = ? AND id != ? AND user_id = ?", url, id, userID).Scan(&exist); err == nil {
		message = fmt.Sprintf("stock url %s is already existed.", url)
		errorCode = 2
	} else if categoryID == -1 {
		message = "Category name exceeded length limit."
		errorCode = 3
	} else {
		_, err = db.Exec("UPDATE stock SET stock = ?, url = ?, category_id = ? WHERE id = ? AND user_id = ?", stock, url, categoryID, id, userID)
		if err != nil {
			log.Printf("Failed to edit stock: %v", err)
			c.String(500, "")
			return
		}
		c.JSON(200, gin.H{"status": 1})
		return
	}
	c.JSON(200, gin.H{"status": 0, "message": message, "error": errorCode})
}

func doDeletestock(c *gin.Context) {
	db, err := sql.Open("sqlite3", sqlite)
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		c.String(503, "")
		return
	}
	defer db.Close()
	session := sessions.Default(c)
	userID := session.Get("user_id")
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Printf("Failed to get id param: %v", err)
		c.String(400, "")
		return
	}

	_, err = db.Exec("DELETE FROM stock WHERE id = ? and user_id = ?", id, userID)
	if err != nil {
		log.Printf("Failed to delete stock: %v", err)
		c.String(500, "")
		return
	}
	c.JSON(200, gin.H{"status": 1})
}

func reorder(c *gin.Context) {
	db, err := sql.Open("sqlite3", sqlite)
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		c.String(503, "")
		return
	}
	defer db.Close()
	session := sessions.Default(c)
	userID := session.Get("user_id")

	orig := c.PostForm("orig")
	dest := c.PostForm("dest")
	next := c.PostForm("next")

	var origSeq, destSeq int

	err = db.QueryRow("SELECT seq FROM seq WHERE stock_id = ? AND user_id = ?", orig, userID).Scan(&origSeq)
	if err != nil {
		log.Printf("Failed to scan orig seq: %v", err)
		c.String(500, "")
		return
	}
	if dest != "#TOP_POSITION#" {
		err = db.QueryRow("SELECT seq FROM seq WHERE stock_id = ? AND user_id = ?", dest, userID).Scan(&destSeq)
	} else {
		err = db.QueryRow("SELECT seq FROM seq WHERE stock_id = ? AND user_id = ?", next, userID).Scan(&destSeq)
		destSeq--
	}
	if err != nil {
		log.Printf("Failed to scan dest seq: %v", err)
		c.String(500, "")
		return
	}

	if origSeq > destSeq {
		destSeq++
		_, err = db.Exec("UPDATE seq SET seq = seq+1 WHERE seq >= ? AND user_id = ? AND seq < ?", destSeq, userID, origSeq)
	} else {
		_, err = db.Exec("UPDATE seq SET seq = seq-1 WHERE seq <= ? AND user_id = ? AND seq > ?", destSeq, userID, origSeq)
	}
	if err != nil {
		log.Printf("Failed to update other seq: %v", err)
		c.String(500, "")
		return
	}
	_, err = db.Exec("UPDATE seq SET seq = ? WHERE stock_id = ? AND user_id = ?", destSeq, orig, userID)
	if err != nil {
		log.Printf("Failed to update orig seq: %v", err)
		c.String(500, "")
		return
	}
	c.String(200, "1")
}
