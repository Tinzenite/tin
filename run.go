package main

/*
loadTinzenite loads an existing Tinzenite directory and runs it.
*/import (
	"github.com/tinzenite/core"
	"github.com/tinzenite/shared"
	"log"
)

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
