package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"os/exec"
)

func execScript(file string) {
	var cmd string
	switch OS {
	case "windows":
		cmd = "python"
	case "linux":
		cmd = "python3"
	default:
		log.Fatal("Unsupported operating system.")
	}

	var args []string
	args = append(args, sqlitePy)
	args = append(args, "restore")
	args = append(args, sqlite)
	args = append(args, file)

	c := exec.Command(cmd, args...)
	var stderr bytes.Buffer
	c.Stderr = &stderr
	if err := c.Run(); err != nil {
		log.Fatalf("Failed to execute sqlite3 script: %s\n%v", stderr.String(), err)
	}
}

func dump() string {
	tmpfile, err := ioutil.TempFile("", "tmp")
	if err != nil {
		log.Fatalf("Failed to create temporary file: %v", err)
	}
	tmpfile.Close()

	var cmd string
	switch OS {
	case "windows":
		cmd = "python"
	case "linux":
		cmd = "python3"
	default:
		log.Fatal("Unsupported operating system.")
	}

	var args []string
	args = append(args, sqlitePy)
	args = append(args, "backup")
	args = append(args, sqlite)
	args = append(args, tmpfile.Name())

	c := exec.Command(cmd, args...)
	var stderr bytes.Buffer
	c.Stderr = &stderr
	if err := c.Run(); err != nil {
		log.Fatalf("Failed to run backup command: %s\n%v", stderr.String(), err)
	}
	return tmpfile.Name()
}
