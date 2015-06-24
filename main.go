package main

import (
	"bufio"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"time"

	"github.com/tinzenite/core"
)

const name = "music"
const path = "/home/tamino/Music"
const user = "Xamino"

func main() {
	model()
	// walkTest()
	// channel()
	// tinzenite()
}

func model() {
	start := time.Now()
	m, err := core.LoadModel(path)
	if err != nil {
		log.Println(err.Error())
	}
	elapsed := time.Since(start)
	elapsed = elapsed / time.Millisecond
	log.Printf("Output:\n\n%s\n", m.String())
	log.Printf("Took %d msec\n", elapsed)
	// register channel for updates
	updates := make(chan core.UpdateMessage, 1)
	m.Register(updates)
	go func() {
		for {
			update := <-updates
			log.Printf("Update received! Type: %s\n%+v\n", update.Operation, update.Object)
		}
	}()
	// now allow manual operations
	reader := bufio.NewReader(os.Stdin)
	for {
		input, _ := reader.ReadString('\n')
		if strings.Contains(input, "exit") {
			break
		}
		err := m.Update()
		if err != nil {
			log.Println(err.Error())
		}
		// log.Printf("Output:\n\n%s\n", m.String())
		log.Println("Updated")
	}
}

type t struct {
}

func (*t) CallbackMessage(address, message string) {
	log.Println("Incomming: " + message)
}

func (*t) CallbackNewConnection(address, message string) {
	log.Println("Accepting: " + message)
}

func channel() {
	channel, err := core.CreateChannel("test", nil, &t{})
	if err != nil {
		log.Println(err.Error())
		return
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	address, _ := channel.Address()
	log.Println("ID:\n" + address)
	for {
		select {
		case <-c:
			channel.Close()
			goto exit
		case <-time.Tick(10 * time.Second):
			err := channel.Send("ed284a9fa07142cb8f6fa8c821d7f722cf63d2c7f74390566c6949bdb898b33e", "Tick!")
			if err != nil {
				log.Println(err.Error())
			}
		}
	}
exit:
}

func tinzenite() bool {
	tinzenite, err := core.CreateTinzenite(name, path, "shana", user, false)
	if err == core.ErrIsTinzenite {
		err = core.RemoveTinzenite(path)
		if err != nil {
			log.Println(err.Error())
		}
		return false
	}
	if err != nil {
		log.Println(err.Error())
		return false
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	address := tinzenite.Address()
	log.Println("ID:\n" + address)
	select {
	case <-c:
		tinzenite.Close()
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
