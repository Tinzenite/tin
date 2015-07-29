package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"strings"

	"github.com/tinzenite/bootstrap"
	"github.com/tinzenite/core"
	"github.com/tinzenite/shared"
)

const user = "Xamino"
const password = "hunter2"

var path string
var name string
var flagBoot bool

func main() {
	parseFlags()
	if flagBoot {
		bootstrapDirectory()
		return
	}
	tinzeniteDirectory()
}

func bootstrapDirectory() {
	var boot *bootstrap.Bootstrap
	var err error
	if shared.IsTinzenite(path) {
		log.Println("Loading bootstrap")
		boot, err = bootstrap.Load(path, onSuccessfulBootstrap)
	} else {
		log.Println("Creating bootstrap")
		boot, err = bootstrap.Create(path, name, onSuccessfulBootstrap)
	}
	if err != nil {
		log.Println("Bootstrap:", err)
		return
	}
	boot.Store()
	// read input
	reader := bufio.NewReader(os.Stdin)
	run := true
	for run {
		input, _ := reader.ReadString('\n')
		input = strings.Trim(input, "\n")
		if strings.HasPrefix(input, "connect") {
			address := strings.Split(input, " ")[1]
			err := boot.Start(address)
			if err != nil {
				log.Println("Start:", err)
			}
			log.Println("Connecting.")
			continue
		}
		switch input {
		case "store":
			boot.Store()
			log.Println("Stored.")
		case "exit":
			boot.Store()
			boot.Close()
			run = false
		case "status":
			log.Println(boot.PrintStatus())
		default:
			log.Println("CMD UNKNOWN")
		}
	}
	log.Println("DONE")
}

/*
For now just start tinzenite
*/
func onSuccessfulBootstrap() {
	tinzeniteDirectory()
}

func tinzeniteDirectory() {
	var tinzenite *core.Tinzenite
	var err error
	if shared.IsTinzenite(path) {
		log.Println("Loading existing Tinzenite.")
		tinzenite, err = core.LoadTinzenite(path, password)
	} else {
		log.Println("Creating new Tinzenite.")
		tinzenite, err = core.CreateTinzenite("test", path, name, user, password)
	}
	if err != nil {
		log.Println("Failed to start:", err)
		return
	}
	log.Println("Ready.")
	// prepare global console reader (before register because it may directly need it)
	reader := bufio.NewReader(os.Stdin)
	// if all ok, register callback
	tinzenite.RegisterPeerValidation(acceptPeer)
	// now allow manual operations
	run := true
	for run {
		input, _ := reader.ReadString('\n')
		input = strings.Trim(input, "\n")
		switch input {
		case "id":
			address, _ := tinzenite.Address()
			log.Println("ID:\n" + address)
		case "info":
			log.Println("Path:", tinzenite.Path)
		case "store":
			err := tinzenite.Store()
			if err != nil {
				log.Println(err.Error())
			} else {
				log.Println("OK")
			}
		case "sync":
			err := tinzenite.Sync()
			if err != nil {
				log.Println("Sync:", err)
			} else {
				log.Println("OK")
			}
		case "update":
			err := tinzenite.SyncLocal()
			if err != nil {
				log.Println("SyncLocal:", err)
			} else {
				log.Println("OK")
			}
		case "clear":
			os.RemoveAll(tinzenite.Path + "/.tinzenite/temp")
			os.Mkdir(tinzenite.Path+"/.tinzenite/temp", 0777)
			log.Println("OK")
		case "status":
			log.Println(tinzenite.PrintStatus())
		case "exit":
			log.Println("Exiting!")
			run = false
		default:
			log.Println("Unknown command.")
		}
	}
	tinzenite.Close()
}

func acceptPeer(address string, wantsTrust bool) bool {
	log.Printf("Accepting <%s>, wants trust: %+v.\n", address, wantsTrust)
	return true
}

func parseFlags() {
	// define
	flag.BoolVar(&flagBoot, "bootstrap", false, "Flag whether to bootstrap to a network.")
	flag.StringVar(&path, "path", "/home/tamino/Music", "Path of where to run Tinzenite.")
	backup, _ := shared.NewIdentifier()
	flag.StringVar(&name, "name", backup, "Name of the Tinzenite peer.")
	// important: apply
	flag.Parse()
}
