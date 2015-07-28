package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"strings"

	"github.com/tinzenite/core"
	"github.com/tinzenite/shared"
)

const user = "Xamino"
const password = "hunter2"

var path string
var name string

var reader *bufio.Reader

func main() {
	parseFlags()
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
	reader = bufio.NewReader(os.Stdin)
	// if all ok, register callback
	tinzenite.RegisterPeerValidation(acceptPeer)
	// now allow manual operations
	run := true
	for run {
		input, _ := reader.ReadString('\n')
		input = strings.Trim(input, "\n")
		// special case to connect because we need to read the address
		if strings.HasPrefix(input, "connect") {
			address := strings.Split(input, " ")[1]
			err := tinzenite.Connect(address)
			if err != nil {
				log.Println(err)
			}
			log.Println("Requested to", address)
			continue
		}
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
	flag.StringVar(&path, "path", "/home/tamino/Music", "Path of where to run Tinzenite.")
	backup, _ := shared.NewIdentifier()
	flag.StringVar(&name, "name", backup, "Name of the Tinzenite peer.")
	// important: apply
	flag.Parse()
	log.Println("Starting at", path, "as", name)
}
