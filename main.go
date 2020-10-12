package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sunshineplan/metadata"
	"github.com/sunshineplan/utils"
	"github.com/vharitonsky/iniflags"
	"golang.org/x/sys/windows/svc"
)

// OS is the running program's operating system
const OS = runtime.GOOS

var metadataConfig metadata.Config

var self string
var unix, host, port, logPath *string

var (
	joinPath = filepath.Join
	dir      = filepath.Dir
)

func init() {
	var err error
	self, err = os.Executable()
	if err != nil {
		log.Fatalln("Failed to get self path:", err)
	}
	os.MkdirAll(joinPath(dir(self), "instance"), 0755)
	sqlite = joinPath(dir(self), "instance/mystocks.db")
	sqlitePy = joinPath(dir(self), "scripts/sqlite.py")
}

func usage(errmsg string) {
	fmt.Fprintf(os.Stderr,
		`%s

usage: %s <command>
       where <command> is one of install, remove, start, stop.
`, errmsg, os.Args[0])
	os.Exit(2)
}

func main() {
	flag.StringVar(&metadataConfig.Server, "server", "", "Metadata Server Address")
	flag.StringVar(&metadataConfig.VerifyHeader, "header", "", "Verify Header Header Name")
	flag.StringVar(&metadataConfig.VerifyValue, "value", "", "Verify Header Value")
	unix = flag.String("unix", "", "UNIX-domain Socket")
	host = flag.String("host", "0.0.0.0", "Server Host")
	port = flag.String("port", "8888", "Server Port")
	//logPath = flag.String("log", joinPath(dir(self), "access.log"), "Log Path")
	logPath = flag.String("log", "", "Log Path")
	iniflags.SetConfigFile(joinPath(dir(self), "config.ini"))
	iniflags.SetAllowMissingConfigFile(true)
	iniflags.Parse()

	isIntSess, err := svc.IsAnInteractiveSession()
	if err != nil {
		log.Fatalln("failed to determine if we are running in an interactive session:", err)
	}
	if !isIntSess {
		runService(svcName, false)
		return
	}

	switch flag.NArg() {
	case 0:
		run()
	case 1:
		switch flag.Arg(0) {
		case "run", "debug":
			run()
		case "install":
			err = installService(svcName, svcDescription)
		case "remove":
			err = removeService(svcName)
		case "start":
			err = startService(svcName)
		case "stop":
			err = controlService(svcName, svc.Stop, svc.Stopped)
		case "backup":
			backup()
		case "init":
			if utils.Confirm("Do you want to initialize database?", 3) {
				restore("")
			}
		default:
			usage(fmt.Sprintf("Unknown argument: %s", flag.Arg(0)))
		}
	case 2:
		switch flag.Arg(0) {
		case "add":
			addUser(flag.Arg(1))
		case "delete":
			if utils.Confirm(fmt.Sprintf("Do you want to delete user %s?", flag.Arg(1)), 3) {
				deleteUser(flag.Arg(1))
			}
		case "restore":
			if utils.Confirm("Do you want to restore database?", 3) {
				restore(flag.Arg(1))
			}
		default:
			log.Fatalf("Unknown arguments: %s", strings.Join(flag.Args(), " "))
		}
	default:
		usage(fmt.Sprintf("Unknown arguments: %s", strings.Join(flag.Args(), " ")))
	}
	if err != nil {
		log.Fatalf("failed to %s %s: %v", flag.Arg(0), svcName, err)
	}
}
