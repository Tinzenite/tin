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

func main() {
	tinzenite()
}

func tinzenite() bool {
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
			err := tinzenite.Sync()
			if err != nil {
				log.Println(err.Error())
			} else {
				log.Println("OK")
			}
		case "update":
			err := tinzenite.SyncLocal()
			if err != nil {
				log.Println(err.Error())
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
	return false
}
