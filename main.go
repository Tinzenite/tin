package main

import (
	"flag"
	"log"
	"strings"

	"github.com/tinzenite/shared"
)

func main() {
	log.Println("Starting client.")
	// declare flags
	var path = *flag.String("path", "", "File directory path in which to run the client.")
	// read them
	flag.Parse()
	// ask which operation to do
	opQuestion := createQuestion("(L)oad a Tinzenite directory, (C)reate one, or (B)ootstrap to an existing one?")
	opQuestion.createAnswer(0, "l", "load")
	opQuestion.createAnswer(1, "c", "create")
	opQuestion.createAnswer(2, "b", "bootstrap")
	switch opQuestion.ask() {
	case 0:
		/*TODO load*/
		log.Println("Load")
	case 1:
		/*TODO create*/
		log.Println("Create")
	case 2:
		/*TODO bootstrap*/
		log.Println("Bootstrap")
	default:
		log.Println("Question returned unknown operation!")
		return
	}
	/*TODO continue here*/
	log.Println("Done for now")
	return
	// if no path was given we need to read the directory list and let the user choose which dir to run
	if path == "" {
		options, err := shared.ReadDirectoryList()
		if err != nil {
			logMain(err.Error())
			return
		}
		if len(options) == 0 {
			log.Println("NONE AVAILABLE")
		}
		log.Println("Choose which ")
	}
}

func loadTinzenite(path string) {
	log.Println("TODO")
}

/*
Log function that respects the AllowLogging flag.
*/
func logMain(msg ...string) {
	toPrint := []string{"MAIN:"}
	toPrint = append(toPrint, msg...)
	log.Println(strings.Join(toPrint, " "))
}
