package main

import (
	"fmt"
	"os"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	promotedFolder string
	listenAddress  string
	listenPort     int
)

func init() {

	pflag.StringVar(&promotedFolder, "path", "", "promoted folder")
	pflag.StringVar(&listenAddress, "address", "192.168.0.10", "listen address")
	pflag.IntVar(&listenPort, "port", 9090, "listen port")

	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		log.Fatal(err)
	}

	parseCmdLine()
}

func parseCmdLine() {
	pflag.Parse()

	if promotedFolder == "" {
		fmt.Println("invalid promoted folder")
		pflag.PrintDefaults()
		os.Exit(1)
	}
}
