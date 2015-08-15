package main

/*
loadTinzenite loads an existing Tinzenite directory and runs it.
*/import (
	"fmt"
	"github.com/tinzenite/core"
	"github.com/tinzenite/shared"
	"log"
)

func bootstrapTinzenite(path string) {

}

func createTinzenite(path string) {
	if shared.IsTinzenite(path) {
		logMain("Directory is already a valid Tinzenite directory!")
		return
	}
	// get options
	peerName := getString("Enter the peer name for this Tinzenite directory:")
	userName := getString("Your username:")
	password := getString("Enter a directory password:")
	relPath := shared.CreatePathRoot(path)
	tinzenite, err := core.CreateTinzenite(relPath.LastElement(), relPath.FullPath(), peerName, userName, password)
	if err != nil {
		logMain("Creation error:", err.Error())
		return
	}
	runTinzenite(tinzenite)
}

func loadTinzenite(path string) {
	if !shared.IsTinzenite(path) {
		logMain("Directory is not a valid Tinzenite directory!")
		return
	}
	password := getString("Please enter the directory password:")
	tinzenite, err := core.LoadTinzenite(path, password)
	if err != nil {
		// TODO catch wrong password and allow retry
		logMain("Loading error:", err.Error())
		return
	}
	tinzenite.RegisterPeerValidation(allowPeer)
	// run tinzenite until killed
	runTinzenite(tinzenite)
}

func runTinzenite(t *core.Tinzenite) {
	address, _ := t.Address()
	fmt.Printf("Running peer <%s>.\nID: %s\n", t.Name(), address)
	// TODO print all relevant info and start bg thread that keeps it running until killed
	log.Println("TODO: implement RUN")
}

/*
allowPeer asks the user whether to accept the given peer.
*/
func allowPeer(address string, wantsTrust bool) bool {
	// TODO actually ask!
	log.Println("TODO: Actually ask, for now accepting everything!")
	return true
}
