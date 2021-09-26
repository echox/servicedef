// Package config is responsible for parsing the command line arguments and
// holding the configured state of the application.
package config

import (
	"flag"
	"log"
	"os"
)

// Config state of the application
type Config struct {
	Hosts_File        *os.File //Filehandle for the json hosts inventory, created during Init()
	Services_File     *os.File //Filehandle for the json service catalogs, created during Init()
	Rules_File        *os.File //Filehandle for the json rule file, created during Init()
	Rules             string   //Relative filesystem path to the rules json file
	Services          bool     //Relative filesystem path to the servies json file
	Graphviz          string   //Relative filesystem path for the graphviz export file
	Progress_Seconds  int      //Seconds after scanning progress gets logged to console
	Connect_Scan      bool     //Should a TCP connect scan used during scanning
	Default_Port_Scan bool     //Should only a small range of default ports be scanned
	Threads           int      //Number of concurrent scanning threads
	Quiet             bool     //Should the log be omitted
}

// Init is used to parse the command line parameters, open file handles and
// return the state of the current configuration
// It uses golangs "flag" package for parsing command line flags
func Init() Config {

	var cfg Config

	flag.StringVar(&cfg.Graphviz, "g", "", "graphviz dot file export")
	flag.StringVar(&cfg.Rules, "r", "", "use rules.json file")
	flag.IntVar(&cfg.Progress_Seconds, "p", 60, "print nmap scanning progress every x seconds")
	flag.BoolVar(&cfg.Connect_Scan, "c", false, "do a nmap connect scan (doesn't require root privileges)")
	flag.BoolVar(&cfg.Default_Port_Scan, "f", false, "scan only nmap default ports instead of all (faster)")
	flag.IntVar(&cfg.Threads, "t", 3, "number of parallel nmap scanning threads")
	flag.BoolVar(&cfg.Quiet, "q", false, "quiet - don't print to stdout")
	flag.Parse()

	args := len(flag.Args())
	if args != 1 && args != 2 {
		log.Printf("wrong arguments (%v) use: ./servicedef hosts.json [services.json]", args)
		os.Exit(1)
	}

	if h, err := os.Open(flag.Arg(0)); err != nil {
		log.Println(err)
		os.Exit(1)
	} else {
		cfg.Hosts_File = h
	}

	if args == 2 {
		if s, err := os.Open(flag.Arg(1)); err != nil {
			log.Println(err)
			os.Exit(1)
		} else {
			cfg.Services_File = s
			cfg.Services = true
		}
	}

	if cfg.Rules != "" {
		if r, err := os.Open(cfg.Rules); err != nil {
			log.Println(err)
			os.Exit(1)
		} else {
			cfg.Rules_File = r
		}
	}

	return cfg
}
