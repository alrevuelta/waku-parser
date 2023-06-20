package config

import (
	"flag"
	"os"

	log "github.com/sirupsen/logrus"
)

type CliConfig struct {
	DockerHost       string
	TimeoutInMilisec uint64
	LogLevel         string
}

// By default the release is a custom build. CI takes care of upgrading it with
// go build -v -ldflags="-X 'github.com/dappnode/mev-sp-oracle/config.ReleaseVersion=x.y.z'"
var ReleaseVersion = "custom-build"

func NewCliConfig() (*CliConfig, error) {
	// Optional flags
	var version = flag.Bool("version", false, "Prints the release version and exits")
	var dockerHost = flag.String("docker-host", "unix:///var/run/docker.sock", "Endpoint of the docker host, see default")
	var timeOutMilisec = flag.Uint64("timeout-ms", 5*1000, "Timeout in milliseconds to consider a message lost")
	var logLevel = flag.String("log-level", "info", "Logging verbosity (trace, debug, info=default, warn, error, fatal, panic)")

	// Mandatory flags
	//var xxx = flag.String("xxx", "", "yyy")

	flag.Parse()

	if *version {
		log.Info("Version: ", ReleaseVersion)
		os.Exit(0)
	}

	cliConf := &CliConfig{
		DockerHost:       *dockerHost,
		TimeoutInMilisec: *timeOutMilisec,
		LogLevel:         *logLevel,
	}
	logConfig(cliConf)
	return cliConf, nil
}

func logConfig(cfg *CliConfig) {
	log.WithFields(log.Fields{
		"ReleaseVersion":   ReleaseVersion,
		"DockerHost":       cfg.DockerHost,
		"TimeoutInMilisec": cfg.TimeoutInMilisec,
		"LogLevel":         cfg.LogLevel,
	}).Info("Cli Config:")
}
