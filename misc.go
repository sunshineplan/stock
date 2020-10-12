package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/sunshineplan/utils/mail"
)

func addUser(username string) {
	log.Print("Start!")
	db, err := getDB()
	if err != nil {
		log.Fatalln("Failed to connect to database:", err)
	}
	defer db.Close()

	if _, err := db.Exec("INSERT INTO user(username) VALUES (?)", strings.ToLower(username)); err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			log.Fatalf("Username %s already exists.", strings.ToLower(username))
		} else {
			log.Fatalln("Failed to add user:", err)
		}
	}
	log.Print("Done!")
}

func deleteUser(username string) {
	log.Print("Start!")
	db, err := getDB()
	if err != nil {
		log.Fatalln("Failed to connect to database:", err)
	}
	defer db.Close()

	res, err := db.Exec("DELETE FROM user WHERE username=?", strings.ToLower(username))
	if err != nil {
		log.Fatalln("Failed to delete user:", err)
	}
	if n, err := res.RowsAffected(); err != nil {
		log.Fatalln("Failed to get affected rows:", err)
	} else if n == 0 {
		log.Fatalf("User %s does not exist.", strings.ToLower(username))
	}
	log.Print("Done!")
}

func backup() {
	log.Print("Start!")
	m, err := metadataConfig.Get("mystocks_backup")
	if err != nil {
		log.Fatalln("Failed to get mystocks_backup metadata:", err)
	}
	var mailSetting mail.Setting
	err = json.Unmarshal(m, &mailSetting)
	if err != nil {
		log.Fatalln("Failed to unmarshal json:", err)
	}

	file := dump()
	defer os.Remove(file)
	if err := mailSetting.Send(
		fmt.Sprintf("My stocks Backup-%s", time.Now().Format("20060102")),
		"",
		&mail.Attachment{FilePath: file, Filename: "database"},
	); err != nil {
		log.Fatalln("Failed to send mail:", err)
	}
	log.Print("Done!")
}

func restore(file string) {
	log.Print("Start!")
	if file == "" {
		file = joinPath(dir(self), "scripts/schema.sql")
	} else {
		if _, err := os.Stat(file); err != nil {
			log.Fatalln("File not found:", err)
		}
	}
	dropAll := joinPath(dir(self), "scripts/drop_all.sql")
	execScript(dropAll)
	execScript(file)
	log.Print("Done!")
}
