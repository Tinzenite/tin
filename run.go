package main

/*
loadTinzenite loads an existing Tinzenite directory and runs it.
*/import (
	"fmt"
	"github.com/tinzenite/bootstrap"
	"github.com/tinzenite/core"
	"github.com/tinzenite/shared"
	"log"
	"os"
	"os/signal"
	"time"
)

func bootstrapTinzenite(path string) {
	var boot *bootstrap.Bootstrap
	var err error
	done := make(chan bool)
	if shared.IsTinzenite(path) {
		boot, err = bootstrap.Load(path, func() {
			done <- true
		})
		if err != nil {
			logMain("Bootstrap load error:", err.Error())
			return
		}
	} else {
		peerName := getString("Enter the peer name for this Tinzenite directory:")
		boot, err = bootstrap.Create(path, peerName, func() {
			done <- true
		})
		if err != nil {
			logMain("Bootstrap create error:", err.Error())
			return
		}
		// connect to:
		address := getString("Please enter the address of the peer to connect to:")
		err = boot.Start(address)
		if err != nil {
			logMain("Bootstrap start error:", err.Error())
			// return because we don't want to store a faulty bootstrap
			return
		}
		err = boot.Store()
		if err != nil {
			logMain("Bootstrap store error:", err.Error())
		}
	}
	// print information
	address, _ := boot.Address()
	fmt.Printf("Bootstrapping.\nID: %s\n", address)
	// wait for successful bootstrap
	<-done
	log.Println("Closing Bootstrap.")
	// this is required before closing boot because ToxCore may still need to
	// notify the other client that the file transfers are complete - this can
	// take a few iterations, so we delay for a second to give it time to do that.
	<-time.Tick(1 * time.Second)
	// manually close boot if we're done! It won't close itself!
	boot.Close()
	// continue to executing the directory
	fmt.Println("Bootstrapping was successful. Loading Tinzenite.")
	// load tinzenite with password
	password := getPassword()
	loadTinzenite(path, password)
}

func createTinzenite(path string) {
	if shared.IsTinzenite(path) {
		logMain("Directory is already a valid Tinzenite directory!")
		return
	}
	// get options
	peerName := getString("Enter the peer name for this Tinzenite directory:")
	userName := getString("Enter your username:")
	password := getString("Enter a password for this Tinzenite network:")
	relPath := shared.CreatePathRoot(path)
	tinzenite, err := core.CreateTinzenite(relPath.LastElement(), relPath.FullPath(), peerName, userName, password)
	if err != nil {
		logMain("Creation error:", err.Error())
		return
	}
	err = tinzenite.SyncLocal()
	if err != nil {
		logMain("Initial model sync error:", err.Error())
	}
	// run tinzenite until killed
	runTinzenite(tinzenite)
}

func loadTinzenite(path, password string) {
	if !shared.IsTinzenite(path) {
		logMain("Directory is not a valid Tinzenite directory!")
		return
	}
	tinzenite, err := core.LoadTinzenite(path, password)
	if err != nil {
		// TODO catch wrong password and allow retry
		logMain("Loading error:", err.Error())
		return
	}
	err = tinzenite.SyncLocal()
	if err != nil {
		logMain("Initial model sync error:", err.Error())
	}
	// run tinzenite until killed
	runTinzenite(tinzenite)
}

/*
runTinzenite runs the given Tinzenite instance.
*/
func runTinzenite(t *core.Tinzenite) {
	// do this here so that it is guaranteed to be set
	t.RegisterPeerValidation(allowPeer)
	// print important info
	address, _ := t.Address()
	fmt.Printf("Running peer <%s>.\nID: %s\n", t.Name(), address)
	// run update and sync in intervalls
	var counter int
	// only build this once instead of every time
	tickSpan := time.Duration(10) * time.Second
	ticker := time.Tick(tickSpan)
	// prepare quitting via ctrl-c
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	// loop until close
	for {
		select {
		case <-ticker:
			if counter >= 5 {
				// TODO remove once Merge bug is fixed
				log.Println("DEBUG: Model sync ---------------------------")
				counter = 0
				err := t.SyncRemote()
				if err != nil {
					logMain("SyncRemote error:", err.Error())
				}
				continue
			}
			// log.Println("DEBUG: Update")
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
	// TODO actually ask! NOTE: shouldn't block... how? Call silent add friend on success?
	log.Println("TODO: Actually ask, for now accepting everything!", address)
	return true
}
