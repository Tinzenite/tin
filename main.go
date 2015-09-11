package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime/pprof"
	"strings"

	"github.com/tinzenite/shared"
)

func main() {
	log.Println("Starting client.")
	// declare flags
	var commandString string
	var path string
	var password string
	var cpuProfileFile string
	// write flag stuff
	flag.StringVar(&path, "path", "", "File directory path in which to run the client.")
	flag.StringVar(&password, "pwd", "", "Password for loading or creating peers.")
	flag.StringVar(&commandString, "cmd", "load", "Command for the path: create, load, or bootstrap (short: boot). Default is load.")
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
	// replace short form with long form
	if commandString == "boot" {
		commandString = "bootstrap"
	}
	// need to do some additional work because flag doesn't allow custom enumeration variables
	command := cmdParse(commandString)
	// do some path work
	if path == "" {
		// if empty get it
		path = getPath()
	}
	// make sure the path is clean and absolute
	path, _ = filepath.Abs(filepath.Clean(path))
	// path may not be empty OR only contain '.' (which means error in filepath)
	if path == "" || path == "." {
		logMain("No path given!")
		return
	}
	// check if path exists
	if exists, _ := shared.DirectoryExists(path); !exists {
		// offer creating it
		if createYesNo("Path <"+path+"> doesn't exist. Create it?").ask() < 0 {
			// explain why we're quitting
			logMain("Can not run Tinzenite without valid path.")
			return
		}
		// create it
		logMain("Creating path <" + path + ">.")
		shared.MakeDirectory(path)
		// and continue below
	}
	logMain("Will", command.String(), "Tinzenite at", path, ".")
	switch command {
	case cmdLoad:
		// get password if not given yet (the other options will ask themselves)
		if password == "" {
			password = getPassword()
		}
		loadTinzenite(path, password)
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

func getPassword() string {
	return getString("Please enter the directory password:")
}

/*
Log function that respects the AllowLogging flag.
*/
func logMain(msg ...string) {
	toPrint := []string{"Main:"}
	toPrint = append(toPrint, msg...)
	log.Println(strings.Join(toPrint, " "))
}
