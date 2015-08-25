package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"strings"

	"github.com/tinzenite/shared"
)

func main() {
	log.Println("Starting client.")
	// declare flags
	var commandString string
	var path string
	var cpuProfileFile string
	// write flag stuff
	flag.StringVar(&path, "path", "", "File directory path in which to run the client.")
	flag.StringVar(&commandString, "cmd", "load", "Command for the path: create, load, or bootstrap. Default is load.")
	flag.StringVar(&cpuProfileFile, "profile", "", "By using this flag with a path, a cpu profile will be written to the given path.")
	// parse them
	flag.Parse()
	// cpu profiling
	if cpuProfileFile != "" {
		f, err := os.Create(cpuProfileFile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
	// need to do some additional work because flag doesn't allow custom enumeration variables
	command := cmdParse(commandString)
	if path == "" {
		path = getPath()
	}
	if !shared.FileExists(path) {
		/*TODO offer creating it?*/
		logMain("Path", path, "doesn't exist!")
		return
	}
	logMain("Will", command.String(), "Tinzenite at", path, ".")
	switch command {
	case cmdLoad:
		loadTinzenite(path)
	case cmdCreate:
		createTinzenite(path)
	case cmdBootstrap:
		bootstrapTinzenite(path)
	default:
		logMain("No command was chosen, so we'll do nothing.")
	}
	// if we reach this it means the client is closing.
	logMain("Quitting.")
}

/*
getPath gets a path from the user, either manually entered or chosen from the
known paths.

TODO: rewrite, I call the same manual request in 3 places!
*/
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
		return getString("Enter path for Tinzenite directory:")
	}
	newQueston := createYesNo("Choose from existing paths?")
	// if no --> manual entry
	if newQueston.ask() < 0 {
		return getString("Enter path for Tinzenite directory:")
	}
	// if only one choice then that is all they have
	if len(options) == 1 {
		useQuestion := createYesNo("Only one candidate: " + options[0] + ". Use it?")
		if useQuestion.ask() < 0 {
			return getString("Enter path for Tinzenite directory:")
		}
		return options[0]
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
	toPrint := []string{"Main:"}
	toPrint = append(toPrint, msg...)
	log.Println(strings.Join(toPrint, " "))
}
