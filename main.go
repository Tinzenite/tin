package main

import "github.com/tinzenite/core"

// TODO check what Golang offers us to help here
const configPath = ""

func main() {
	context, _ := loadContext()
	// just run core for now
	core.Run(*context)
}

// loadContext or create a new one if required. If new will save it.
func loadContext() (*core.Context, error) {
	// create new context to use
	context, err := core.NewContext("NewTest")
	// TODO check if one already exists, if yes use it. Otherwise create new one and save it.
	return context, err
}
