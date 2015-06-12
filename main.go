package main

import (
	"log"

	"github.com/tinzenite/core"
)

const path = "/home/tamino/Music"

func main() {
	var context *core.Context
	var err error
	if core.IsTinzenite(path) {
		context, err = core.Load(path)
	} else {
		context, err = core.Create("Test", path)
	}
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	err = context.Run()
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	context.Kill()
}
