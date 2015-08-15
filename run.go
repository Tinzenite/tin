package main

/*
loadTinzenite loads an existing Tinzenite directory and runs it.
*/import (
	"fmt"
	"github.com/tinzenite/core"
	"github.com/tinzenite/shared"
	"log"
	"os"
	"os/signal"
	"time"
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

/*
runTinzenite runs the given Tinzenite instance.

TODO implement interrupts and close behaviour!
*/
func runTinzenite(t *core.Tinzenite) {
	// print important info
	address, _ := t.Address()
	fmt.Printf("Running peer <%s>.\nID: %s\n", t.Name(), address)
	// run update and sync in intervalls
	var counter int
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	for {
		select {
		case <-time.Tick(time.Duration(5) * time.Second):
			if counter >= 5 {
				counter = 0
				err := t.SyncRemote()
				if err != nil {
					logMain("SyncRemote error:", err.Error())
				}
				continue
			}
			counter++
			err := t.SyncLocal()
			if err != nil {
				logMain("SyncLocal error:", err.Error())
			}
		case <-c:
			// on interrupt close tinzenite
			t.Close()
			return
		} // select
	} // for
}

/*
allowPeer asks the user whether to accept the given peer.
*/
func allowPeer(address string, wantsTrust bool) bool {
	// TODO actually ask!
	log.Println("TODO: Actually ask, for now accepting everything!")
	return true
}
