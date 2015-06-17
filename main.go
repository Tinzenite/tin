package main

import (
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"

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
	list, err := core.ReadTinIgnore(ignorepath)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	var dir int
	var count, file int64
	filepath.Walk(testpath, func(path string, info os.FileInfo, err error) error {
		for _, line := range list {
			isDirLine := strings.HasPrefix(line, "/")
			// check dir only against dir candidates
			if isDirLine && info.IsDir() {
				if strings.HasSuffix(path, line) {
					dir++
					// no need to walk the dir if we ignore it
					return filepath.SkipDir
				}
			} else if !isDirLine && !info.IsDir() {
				// check file only against file ^^
				if strings.HasSuffix(path, line) {
					file += info.Size()
					return nil
				}
			}
		}
		count += info.Size()
		return nil
	})
	log.Printf("Kept %dkb worth of objects, ignored %d directories and %dkb worth of files\n", count/1024, dir, file/1024)
	return false
}
