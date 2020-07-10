package main

import (
	"database/sql"
	"log"
	"strconv"

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
	c.JSON(200, doGetRealtimes(stocks))
}

func indices(c *gin.Context) {
	indices := []stock{&sse{code: "000001"}, &szse{code: "399001"}, &szse{code: "399006"}, &szse{code: "399005"}}
	c.JSON(200, doGetRealtimes(indices))
}

func getStock(c *gin.Context) {
	index := c.Param("index")
	code := c.Param("code ")
	q := c.Param("q")

	stock := initStock(index, code)

	if q == "realtime" {
		c.JSON(200, stock.realtime())
		return
	} else if q == "chart" {
		c.JSON(200, stock.chart())
		return
	}
	c.String(400, "")
}

func getSuggest(c *gin.Context) {
	keyword := c.Param("keyword")
	sse := sseSuggest(keyword)
	szse := szseSuggest(keyword)
	c.JSON(200, append(sse, szse...))
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
