package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/webkeydev/logger"
)

var (
	log      = logger.NewLogger("ota-promoter")
	shutdown = make(chan bool, 1)
)

func init() {
	logger.SetTxtLogger()
	listenDoneSignal()
}

func listenDoneSignal() {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-done
		log.Infof("close watcher")
		closeWatcher()
		log.Infof("shutdown server")
		shutDownServer()
		shutdown <- true
	}()
}

func main() {

	if err := readFiles(); err != nil {
		log.Fatalf("failed to read promoted files: %v", err)
	}

	if err := watch(); err != nil {
		log.Fatalf("%v", err)
	} else {
		log.Infof("watching: %s", promotedFolder)
	}

	listen()

	<-shutdown
	log.Infof("bye")
}
