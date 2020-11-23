package main

import (
	"log"
	"strings"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/sunshineplan/stock"
)

func myStocks(c *gin.Context) {
	db, err := getDB()
	if err != nil {
		log.Println("Failed to connect to database:", err)
		c.String(503, "")
		return
	}
	defer db.Close()

	session := sessions.Default(c)
	userID := session.Get("user_id")
	if userID == nil {
		session.Clear()
		session.Set("user_id", 0)
		session.Save()
		userID = 0
	}

	rows, err := db.Query(`SELECT idx, code FROM stock WHERE user_id = ? ORDER BY seq`, userID)
	if err != nil {
		if strings.Contains(err.Error(), "no such table") {
			restore("")
			c.String(501, "")
			return
		}
		log.Println("Failed to get all stocks:", err)
		c.String(500, "")
		return
	}
	defer rows.Close()
	var stocks []stock.Stock
	for rows.Next() {
		var index, code string
		if err := rows.Scan(&index, &code); err != nil {
			log.Println("Failed to scan all stocks:", err)
			c.String(500, "")
			return
		}
		stocks = append(stocks, stock.Init(index, code))
	}
	c.JSON(200, stock.Realtimes(stocks))
}

func indices(c *gin.Context) {
	indices := stock.Realtimes(
		[]stock.Stock{
			stock.Init("SSE", "000001"),
			stock.Init("SZSE", "399001"),
			stock.Init("SZSE", "399006"),
			stock.Init("SZSE", "399005"),
		})
	c.JSON(200, gin.H{"沪": indices[0], "深": indices[1], "创": indices[2], "中": indices[3]})
}

func getStock(c *gin.Context) {
	var r struct{ Index, Code, Q string }
	if err := c.BindJSON(&r); err != nil {
		c.String(400, "")
		return
	}

	s := stock.Init(r.Index, r.Code)

	if r.Q == "realtime" {
		realtime := s.GetRealtime()
		c.JSON(200, realtime)
		return
	} else if r.Q == "chart" {
		chart := s.GetChart()
		c.JSON(200, chart)
		return
	}
	c.String(400, "")
}

func getSuggest(c *gin.Context) {
	var r struct{ Keyword string }
	if err := c.BindJSON(&r); err != nil {
		c.String(400, "")
		return
	}
	c.JSON(200, stock.Suggests(r.Keyword))
}

func star(c *gin.Context) {
	db, err := getDB()
	if err != nil {
		log.Println("Failed to connect to database:", err)
		c.String(503, "")
		return
	}
	defer db.Close()

	session := sessions.Default(c)
	userID := session.Get("user_id")
	refer := strings.Split(c.Request.Referer(), "/")
	index := refer[len(refer)-2]
	code := refer[len(refer)-1]

	if userID != nil {
		var exist string
		if err := db.QueryRow("SELECT idx FROM stock WHERE idx = ? AND code = ? AND user_id = ?",
			index, code, userID).Scan(&exist); err == nil {
			c.String(200, "1")
			return
		}
	}
	c.String(200, "0")
}

func doStar(c *gin.Context) {
	db, err := getDB()
	if err != nil {
		log.Println("Failed to connect to database:", err)
		c.String(503, "")
		return
	}
	defer db.Close()

	session := sessions.Default(c)
	userID := session.Get("user_id")
	refer := strings.Split(c.Request.Referer(), "/")
	index := refer[len(refer)-2]
	code := refer[len(refer)-1]

	var r struct{ Action string }
	if err := c.BindJSON(&r); err != nil {
		c.String(400, "")
		return
	}

	if userID != nil {
		if r.Action == "unstar" {
			if _, err := db.Exec("DELETE FROM stock WHERE idx = ? AND code = ? AND user_id = ?",
				index, code, userID); err != nil {
				log.Println("Failed to unstar stock:", err)
				c.String(500, "")
				return
			}
		} else {
			if _, err := db.Exec("INSERT INTO stock (idx, code, user_id) VALUES (?, ?, ?)",
				index, code, userID); err != nil {
				log.Println("Failed to star stock:", err)
				c.String(500, "")
				return
			}
		}
		c.String(200, "1")
		return
	}
	c.String(200, "0")
}

func reorder(c *gin.Context) {
	db, err := getDB()
	if err != nil {
		log.Println("Failed to connect to database:", err)
		c.String(503, "")
		return
	}
	defer db.Close()

	session := sessions.Default(c)
	userID := session.Get("user_id")

	var r struct{ New, Old string }
	if err := c.BindJSON(&r); err != nil {
		c.String(400, "")
		return
	}

	orig := strings.Split(r.New, " ")
	dest := strings.Split(r.Old, " ")

	ec := make(chan error, 1)
	var origSeq, destSeq int
	go func() {
		ec <- db.QueryRow("SELECT seq FROM stock WHERE idx = ? AND code = ? AND user_id = ?",
			orig[0], orig[1], userID).Scan(&origSeq)
	}()
	if err := db.QueryRow("SELECT seq FROM stock WHERE idx = ? AND code = ? AND user_id = ?",
		dest[0], dest[1], userID).Scan(&destSeq); err != nil {
		log.Println("Failed to scan dest seq:", err)
		c.String(500, "")
		return
	}
	if err := <-ec; err != nil {
		log.Println("Failed to scan orig seq:", err)
		c.String(500, "")
		return
	}

	go func() {
		_, err := db.Exec("UPDATE stock SET seq = ? WHERE idx = ? AND code = ? AND user_id = ?",
			destSeq, orig[0], orig[1], userID)
		ec <- err
	}()
	if origSeq > destSeq {
		_, err = db.Exec("UPDATE stock SET seq = seq + 1 WHERE seq >= ? AND seq < ? AND user_id = ?",
			destSeq, origSeq, userID)
	} else {
		_, err = db.Exec("UPDATE stock SET seq = seq - 1 WHERE seq > ? AND seq <= ? AND user_id = ?",
			origSeq, destSeq, userID)
	}
	if err != nil {
		log.Println("Failed to update other seq:", err)
		c.String(500, "")
		return
	}
	if err := <-ec; err != nil {
		log.Println("Failed to update orig seq:", err)
		c.String(500, "")
		return
	}
	c.String(200, "1")
}
