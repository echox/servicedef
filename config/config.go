package config

import (
	"flag"
	"log"
	"os"
)

type Config struct {
	Hosts_File    *os.File
	Services_File *os.File
	Services      bool
	Graphviz      string
	Progress_Seconds int
	Connect_Scan  bool
}

func Init() Config {

	var cfg Config

	flag.StringVar(&cfg.Graphviz, "g", "", "graphviz dot file export")
	flag.IntVar(&cfg.Progress_Seconds, "p", 60, "print nmap scanning progress every x seconds")
	flag.BoolVar(&cfg.Connect_Scan, "c", false, "do a nmap connect scan (doesn't require root privileges)")
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

	return cfg
}
