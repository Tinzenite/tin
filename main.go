package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/tinzenite/core"
)

const path = "/home/tamino/Music"

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
	_, err := core.CreateTinzenite(path, false)
	log.Println(err.Error())
	return false
}
