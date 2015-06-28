package main

import (
	"bufio"
	"log"
	"os"
	"strings"

	"github.com/tinzenite/core"
)

const name = "music"
const path = "/home/tamino/Music"
const user = "Xamino"

func main() {
	tinzenite()
}

func tinzenite() bool {
	var tinzenite *core.Tinzenite
	var err error
	if core.IsTinzenite(path) {
		log.Println("Loading existing.")
		tinzenite, err = core.LoadTinzenite(path)
	} else {
		log.Println("Creating new.")
		tinzenite, err = core.CreateTinzenite(name, path, "shana", user)
	}
	if err != nil {
		log.Println("Failed to start: " + err.Error())
		return false
	}
	address := tinzenite.Address()
	log.Println("ID:\n" + address)
	// now allow manual operations
	reader := bufio.NewReader(os.Stdin)
	run := true
	for run {
		input, _ := reader.ReadString('\n')
		input = strings.Trim(input, "\n")
		switch input {
		case "store":
			err := tinzenite.Store()
			if err != nil {
				log.Println(err.Error())
			} else {
				log.Println("OK")
			}
		case "sync":
			err := tinzenite.SyncModel()
			if err != nil {
				log.Println(err.Error())
			} else {
				log.Println("OK")
			}
		case "exit":
			log.Println("Exiting!")
			run = false
		default:
			log.Println("Unknown command.")
		}
	}
	tinzenite.Close()
	return false
}
