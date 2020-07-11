package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"os/exec"
	"strings"
)

var sqlitePy = joinPath(dir(self), "scripts/sqlite.py")

func execScript(file string) {
	var cmd, arg string
	switch OS {
	case "windows":
		cmd = "cmd"
		arg = "/c"
	case "linux":
		cmd = "bash"
		arg = "-c"
	default:
		log.Fatal("Unsupported operating system.")
	}

	var args []string
	args = append(args, sqlitePy)
	args = append(args, "restore")
	args = append(args, sqlite)
	args = append(args, file)

	c := exec.Command(cmd, arg, strings.Join(args, " "))
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

	var cmd, arg string
	switch OS {
	case "windows":
		cmd = "cmd"
		arg = "/c"
	case "linux":
		cmd = "bash"
		arg = "-c"
	default:
		log.Fatal("Unsupported operating system.")
	}

	var args []string
	args = append(args, sqlitePy)
	args = append(args, "backup")
	args = append(args, sqlite)
	args = append(args, tmpfile.Name())

	dump := exec.Command(cmd, arg, strings.Join(args, " "))
	var stderr bytes.Buffer
	dump.Stderr = &stderr
	if err := dump.Run(); err != nil {
		log.Fatalf("Failed to run backup command: %s\n%v", stderr.String(), err)
	}
	return tmpfile.Name()
}
