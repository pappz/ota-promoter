package main

import (
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
	formatter "github.com/webkeydev/logger"

	"github.com/pappz/ota-promoter/promoter"
	"github.com/pappz/ota-promoter/web"
)

var (
	shutdown = make(chan bool, 1)

	webServer     web.Server
	promo         *promoter.Promoter
	changeWatcher promoter.ChangeWatcher
)

func init() {
	prepareLogFormatter()
	listenDoneSignal()
}

func prepareLogFormatter() {
	formatter.SetTxtFormatterForLogger(log.StandardLogger())
	log.StandardLogger().SetLevel(log.DebugLevel)
}

func listenDoneSignal() {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-done

		log.Infof("close watcher")
		changeWatcher.CloseWatcher()

		log.Infof("shutdown web server")
		err := webServer.ShutDownServer()
		if err != nil {
			log.Errorf("%s", err.Error())
		}
		shutdown <- true
	}()
}

func watcherCb() {
	err := promo.ReadFiles()
	if err != nil {
		log.Errorf("failed to read promoted files: %v", err)
	}
}

func main() {
	cfg, err := newConfig()
	if err != nil {
		os.Exit(1)
	}
	promo = promoter.NewPromoter(cfg.promotedFolder)
	err = promo.ReadFiles()
	if err != nil {
		log.Fatalf("failed to read files: %s", err)
	}

	changeWatcher, err = promoter.NewChangeWatcher()
	if err != nil {
		log.Fatalf("failed to setup watcher: %s", err)
	}

	err = changeWatcher.Watch(cfg.promotedFolder, watcherCb, watcherError)
	if err != nil {
		log.Fatalf("failed to start watcher: %s, %s", err)
	}
	log.Infof("watching: %s", cfg.promotedFolder)

	webServer = web.NewServer(cfg.listenAddress, promo)
	webServer.Listen()
	log.Infof("server is lisening on: %s", cfg.listenAddress)

	<-shutdown
	log.Infof("bye")
}

func watcherError(err error) {
	log.Errorf("watch error: %s", err.Error())
}
