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
	// done and onSuccess are used to determine when a bootstrap has completed
	done := make(chan bool, 1)
	onSuccess := func() { done <- true }
	// if tinzenite OR encrypted we can just load the previous bootstrap
	if shared.IsTinzenite(path) || shared.IsEncrypted(path) {
		boot, err = bootstrap.Load(path, onSuccess)
		if err != nil {
			logMain("Bootstrap load error:", err.Error())
			return
		}
	} else {
		// ask whether this is supposed to be a trusted peer
		question := shared.CreateYesNo("Is this a TRUSTED peer?")
		trusted := question.Ask() > 0
		// get peer name
		peerName := shared.GetString("Enter the peer name for this Tinzenite directory:")
		// get address to connect to BEFORE starting boot to avoid terminal clutter
		address := shared.GetString("Please enter the address of the peer to connect to:")
		// build object
		boot, err = bootstrap.Create(path, peerName, trusted, onSuccess)
		if err != nil {
			logMain("Bootstrap create error:", err.Error())
			return
		}
		// connect to:
		err = boot.Start(address)
		if err != nil {
			logMain("Bootstrap start error:", err.Error())
			// return because we don't want to store a faulty bootstrap
			return
		}
		// if everything ok, try storing this bootstrap
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
	// IF w bootstrapped an encrypted peer, write that to the log and quit.
	if !boot.IsTrusted() {
		fmt.Println("Bootstrapping was successful. Run server to start encrypted peer.")
		return
	}
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
	peerName := shared.GetString("Enter the peer name for this Tinzenite directory:")
	userName := shared.GetString("Enter your username:")
	password := shared.GetString("Enter a password for this Tinzenite network:")
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
	t.RegisterPeerValidation(func(address string, wantsTrust bool) {
		var allow bool
		if wantsTrust {
			question := shared.CreateYesNo("Add peer " + address[:8] + " as TRUSTED peer?")
			allow = question.Ask() > 0
		} else {
			question := shared.CreateYesNo("Add peer " + address[:8] + " as ENCRYPTED peer?")
			allow = question.Ask() > 0
		}
		if !allow {
			log.Println("Tin: will not add peer, as requested.")
			return
		}
		// allow peer
		err := t.AllowPeer(address)
		if err != nil {
			log.Println("Tinzenite: failed to allow peer:", err)
		}
		log.Println("Tin: will allow peer, as requested.")
	})
	// print important info
	address, _ := t.Address()
	fmt.Printf("Running peer <%s>.\nID: %s\n", t.Name(), address)
	// build ticks only once instead of every time
	// FIXME: for now using prime numbers to keep them from all ticking at the same time
	tickUpdate := time.Tick(time.Duration(7) * time.Second)
	tickRemote := time.Tick(time.Duration(29) * time.Second)
	tickEncrypted := time.Tick(time.Duration(53) * time.Second)
	// prepare quitting via ctrl-c
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	// loop until close
	for {
		select {
		case <-tickUpdate:
			err := t.SyncLocal()
			if err != nil {
				logMain("SyncLocal error:", err.Error())
			}
		case <-tickRemote:
			err := t.SyncRemote()
			if err != nil {
				logMain("SyncRemote error:", err.Error())
			}
		case <-tickEncrypted:
			err := t.SyncEncrypted()
			if err != nil {
				logMain("SyncEncrypted error:", err.Error())
			}
		case <-c:
			// on interrupt close tinzenite
			t.Close()
			return
		} // select
	} // for
}
