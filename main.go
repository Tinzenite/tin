package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/tinzenite/shared"
)

func main() {
	log.Println("Starting client.")
	// declare flags
	var commandString string
	var path string
	// write flag stuff
	flag.StringVar(&path, "path", "", "File directory path in which to run the client.")
	flag.StringVar(&commandString, "cmd", "none", "Command for the path: create, load, or bootstrap.")
	// parse them
	flag.Parse()
	// need to do some additional work because flag doesn't allow custom enumeration variables
	command := cmdParse(commandString)
	// TODO implement load as sane default where?
	// make sure that path and command have been given, otherwise ask explicitely
	if command == cmdNone {
		// default to load
		command = cmdLoad
	}
	if path == "" {
		path = getPath()
	}
	logMain("Will", command.String(), "Tinzenite at", path, ".")
}

func getCmd() cmd {
	opQuestion := createQuestion("(L)oad a Tinzenite directory, (C)reate one, or (B)ootstrap to an existing one?")
	opQuestion.createAnswer(0, "l", "load")
	opQuestion.createAnswer(1, "c", "create")
	opQuestion.createAnswer(2, "b", "bootstrap")
	switch opQuestion.ask() {
	case 0:
		return cmdLoad
	case 1:
		return cmdCreate
	case 2:
		return cmdBootstrap
	default:
		log.Println("Question returned unknown operation!")
		return cmdNone
	}
}

func getPath() string {
	// load available dirs
	options, err := shared.ReadDirectoryList()
	if err != nil {
		logMain(err.Error())
		return ""
	}
	// if none saved --> ask for manual entry
	if len(options) == 0 {
		fmt.Println("No previous Tinzenite directories known.")
		return getString("Enter path for new Tinzenite directory:")
	}
	newQueston := createYesNo("Choose from existing paths?")
	// if no --> manual entry
	if newQueston.ask() < 0 {
		return getString("Enter path for new Tinzenite directory:")
	}
	fmt.Println("Available paths:")
	for index, path := range options {
		// plus one for human readable numbers
		fmt.Println(index+1, ":", path)
	}
	var pathIndex int
	for {
		pathIndex = getInt("Enter the corresponding number to choose a path:")
		pathIndex-- // need to subtract one to undo human readable numbers
		if pathIndex >= 0 && pathIndex < len(options) {
			break
		}
		fmt.Println("Invalid choice. Choose between 1 and the maximum!")
	}
	return options[pathIndex]
}

/*
Log function that respects the AllowLogging flag.
*/
func logMain(msg ...string) {
	toPrint := []string{"MAIN:"}
	toPrint = append(toPrint, msg...)
	log.Println(strings.Join(toPrint, " "))
}
