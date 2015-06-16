package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/tinzenite/core"
)

const name = "music"
const path = "/home/tamino/Music"
const user = "Xamino"

func main() {

	if !test() {
		return
	}

	channel, err := core.CreateChannel("TestMe", nil)
	if err != nil {
		panic(err)
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	address, _ := channel.Address()
	log.Println("ID:\n" + address)
	select {
	case <-c:
		channel.Close()
	}
}

func test() bool {
	_, err := core.CreateTinzenite(name, path, "shana", user, false)
	if err != nil {
		log.Println(err.Error())
	}
	return false
}
