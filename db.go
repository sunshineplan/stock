package main

import (
	"io/ioutil"
	"log"
	"os/exec"
	"strings"
)

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
	args = append(args, "sqlite3")
	args = append(args, "<")
	args = append(args, file)

	c := exec.Command(cmd, arg, strings.Join(args, " "))
	if err := c.Run(); err != nil {
		log.Fatalf("Failed to execute sqlite3 script: %v", err)
	}
}

func dump() string {
	tmpfile, err := ioutil.TempFile("", "tmp")
	if err != nil {
		log.Fatalf("Failed to create temporary file: %v", err)
	}

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
	args = append(args, "sqlite3dump")

	dump := exec.Command(cmd, arg, strings.Join(args, " "))
	if err := dump.Run(); err != nil {
		log.Fatalf("Failed to run backup command: %v", err)
	}
	return tmpfile.Name()
}
