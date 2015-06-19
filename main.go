package main

import (
	"log"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/tinzenite/core"
)

const name = "music"
const path = "/home/tamino/Music"
const user = "Xamino"

func main() {

	if !walkTest() {
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
	if err == core.ErrIsTinzenite {
		err = core.RemoveTinzenite(path)
		if err != nil {
			log.Println(err.Error())
		}
	}
	if err != nil {
		log.Println(err.Error())
	}
	return false
}

func walkTest() bool {
	testpath := "/home/tamino/Programming"
	ignorepath := "/home/tamino"
	matcher, err := core.CreateMatcher(ignorepath)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	var objects int
	var count int64
	filepath.Walk(testpath, func(path string, info os.FileInfo, err error) error {
		if matcher.Ignore(path) {
			objects++
			return filepath.SkipDir
		}
		count += info.Size()
		return nil
	})
	log.Printf("Kept %dkb worth of objects, ignored %d objects\n", count/1024, objects)
	return false
}
