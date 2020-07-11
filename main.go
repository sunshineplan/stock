package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/vharitonsky/iniflags"
)

// OS is the running program's operating system
const OS = runtime.GOOS

var self, sqlite string
var unix, host, port, logPath *string

var (
	joinPath = filepath.Join
	dir      = filepath.Dir
)

func init() {
	self, err := os.Executable()
	if err != nil {
		log.Fatalf("Failed to get self path: %v", err)
	}
	os.MkdirAll(joinPath(dir(self), "instance"), 0755)
	sqlite = joinPath(dir(self), "instance", "mystocks.db")
}

func main() {
	unix = flag.String("unix", "", "UNIX-domain Socket")
	host = flag.String("host", "127.0.0.1", "Server Host")
	port = flag.String("port", "12345", "Server Port")
	logPath = flag.String("log", joinPath(dir(self), "access.log"), "Log Path")
	iniflags.SetConfigFile(joinPath(dir(self), "config.ini"))
	iniflags.SetAllowMissingConfigFile(true)
	iniflags.Parse()

	switch flag.NArg() {
	case 0:
		run()
	case 1:
		switch flag.Arg(0) {
		case "run":
			run()
		case "backup":
			backup()
		case "init":
			restore("")
		default:
			log.Fatalf("Unknown argument: %s", flag.Arg(0))
		}
	case 2:
		switch flag.Arg(0) {
		case "add":
			addUser(flag.Arg(1))
		case "delete":
			deleteUser(flag.Arg(1))
		case "restore":
			restore(flag.Arg(1))
		default:
			log.Fatalf("Unknown arguments: %s", strings.Join(flag.Args(), " "))
		}
	default:
		log.Fatalf("Unknown arguments: %s", strings.Join(flag.Args(), " "))
	}
}
