package main

import (
	"errors"
	"fmt"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	errMandatoryPath = errors.New("path is mandatory parameter")
)

type config struct {
	promotedFolder string
	listenAddress  string
}

func newConfig() (config, error) {
	cfg := config{}
	if err := cfg.bindVars(); err != nil {
		return cfg, err
	}

	if err := cfg.parseCmdLine(); err != nil {
		return cfg, err
	}

	return cfg, nil
}

func (c *config) bindVars() error {

	pflag.StringVar(&c.promotedFolder, "path", "", "promoted folder")
	pflag.StringVar(&c.listenAddress, "address", ":9090", "listen address")

	return viper.BindPFlags(pflag.CommandLine)
}

func (c *config) parseCmdLine() error {
	pflag.Parse()

	if c.promotedFolder == "" {
		fmt.Printf("%s\nusage:\n", errMandatoryPath)
		pflag.PrintDefaults()
		return errMandatoryPath
	}
	return nil
}
