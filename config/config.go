// Package config is responsible for parsing the command line arguments and
// holding the configured state of the application.
package config

import (
	"flag"
	"fmt"
)

// Config state of the application
type Config struct {
	HostsPath         string //Relative filesystem path for the json hosts inventory
	ServicesPath      string //Relative filesystem path for the json service catalogs
	RulesPath         string //Relative filesystem path for the json rules file
	Graphviz          string //Relative filesystem path for the graphviz export file
	Progress_Seconds  int    //Seconds after scanning progress gets logged to console
	Connect_Scan      bool   //Should a TCP connect scan used during scanning
	Default_Port_Scan bool   //Should only a small range of default ports be scanned
	Threads           int    //Number of concurrent scanning threads
	Loglevel          string //Loglevel for logrus
	JSONLog           bool   //enable JSON logging
}

// Init is used to parse the command line parameters
// It uses golangs "flag" package for parsing
func Init() (Config, error) {

	var cfg Config

	flag.StringVar(&cfg.Graphviz, "g", "", "graphviz dot file export")
	flag.StringVar(&cfg.RulesPath, "r", "", "use rules.json file")
	flag.IntVar(&cfg.Progress_Seconds, "p", 60, "print nmap scanning progress every x seconds")
	flag.BoolVar(&cfg.Connect_Scan, "c", false, "do a nmap connect scan (doesn't require root privileges)")
	flag.BoolVar(&cfg.Default_Port_Scan, "f", false, "scan only nmap default ports instead of all (faster)")
	flag.IntVar(&cfg.Threads, "t", 3, "number of parallel nmap scanning threads")
	flag.StringVar(&cfg.Loglevel, "ll", "info", "loglevel: debug, info, warn, error, fatal - use panic to disable output")
	flag.BoolVar(&cfg.JSONLog, "j", false, "enable structured JSON logging")
	flag.Parse()

	args := len(flag.Args())
	if args != 1 && args != 2 {
		return cfg, fmt.Errorf("wrong arguments (%v) use: ./servicedef hosts.json [services.json]", args)
	}

	cfg.HostsPath = flag.Arg(0)

	if args == 2 {
		cfg.ServicesPath = flag.Arg(1)
	}

	return cfg, nil
}
