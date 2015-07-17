package main

import (
	"bufio"
	"log"
	"os"
	"strings"

	"github.com/tinzenite/core"
	"github.com/tinzenite/shared"
)

const name = "music"
const path = "/home/tamino/Music"
const user = "Xamino"
const password = "hunter2"

var reader *bufio.Reader

func main() {
	var tinzenite *core.Tinzenite
	var err error
	if shared.IsTinzenite(path) {
		log.Println("Loading existing.")
		tinzenite, err = core.LoadTinzenite(path, password)
	} else {
		log.Println("Creating new.")
		tinzenite, err = core.CreateTinzenite(name, path, "shana", user, password)
	}
	if err != nil {
		log.Println("Failed to start:", err)
		return
	}
	// prepare global console reader (before register because it may directly need it)
	reader = bufio.NewReader(os.Stdin)
	// if all ok, register callback
	tinzenite.RegisterPeerValidation(acceptPeer)
	// now allow manual operations
	run := true
	for run {
		input, _ := reader.ReadString('\n')
		input = strings.Trim(input, "\n")
		switch input {
		case "id":
			address := tinzenite.Address()
			log.Println("ID:\n" + address)
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
	var input string
	for {
		switch input {
		case "accept":
			return true
		case "deny":
			return false
		default:
			log.Println("Accept with <accept>, deny with <deny>. All else will be ignored.")
			input, _ = reader.ReadString('\n')
			input = strings.Trim(input, "\n")
		}
	}
	// return false <-- can never be reached
}
