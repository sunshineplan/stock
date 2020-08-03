package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/sunshineplan/utils/mail"
)

func addUser(username string) {
	log.Println("Start!")
	db, err := sql.Open("sqlite3", sqlite)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if _, err := db.Exec("INSERT INTO user(username) VALUES (?)", strings.ToLower(username)); err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			log.Fatalf("Username %s already exists.", strings.ToLower(username))
		} else {
			log.Fatalf("Failed to add user: %v", err)
		}
	}
	log.Println("Done!")
}

func deleteUser(username string) {
	log.Println("Start!")
	db, err := sql.Open("sqlite3", sqlite)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	res, err := db.Exec("DELETE FROM user WHERE username=?", strings.ToLower(username))
	if err != nil {
		log.Fatalf("Failed to delete user: %v", err)
	}
	if n, err := res.RowsAffected(); err != nil {
		log.Fatalf("Failed to get affected rows: %v", err)
	} else if n == 0 {
		log.Fatalf("User %s does not exist.", strings.ToLower(username))
	}
	log.Println("Done!")
}

func backup() {
	log.Println("Start!")
	m, err := metadataConfig.Get("mystocks_backup")
	if err != nil {
		log.Fatalf("Failed to get mystocks_backup metadata: %v", err)
	}
	var mailSetting mail.Setting
	err = json.Unmarshal(m, &mailSetting)
	if err != nil {
		log.Fatalf("Failed to unmarshal json: %v", err)
	}

	file := dump()
	defer os.Remove(file)
	if err := mail.SendMail(
		&mailSetting,
		fmt.Sprintf("My stocks Backup-%s", time.Now().Format("20060102")),
		"",
		&mail.Attachment{FilePath: file, Filename: "database"},
	); err != nil {
		log.Fatalf("Failed to send mail: %v", err)
	}
	log.Println("Done!")
}

func restore(file string) {
	log.Println("Start!")
	if file == "" {
		file = joinPath(dir(self), "scripts/schema.sql")
	} else {
		if _, err := os.Stat(file); err != nil {
			log.Fatalf("File not found: %v", err)
		}
	}
	dropAll := joinPath(dir(self), "scripts/drop_all.sql")
	execScript(dropAll)
	execScript(file)
	log.Println("Done!")
}
