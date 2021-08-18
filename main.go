package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"

	"github.com/echox/servicedef/export"
	. "github.com/echox/servicedef/model"
	"github.com/echox/servicedef/scan"
	"github.com/echox/servicedef/util"

	"github.com/fatih/color"
)

func main() {
	log.Println("servicedef v0")

	cfg := config()

	var services []ServiceDef
	if cfg.Services {
		log.Println("parsing services file...")
		services = parse_services(cfg.Services_File)
		log.Printf("Services #: %v", len(services))
		for _, s := range services {
			print_service(s)
		}
		log.Println("parsing services file finished")
	} else {
		log.Println("no service definitions - scanning only")
	}

	log.Println("parsing hosts file...")

	hosts := parse_hosts(cfg.Hosts_File)
	color.Set(color.FgYellow)
	for _, h := range hosts {
		log.Printf("Host: %v %v", h.Ip, h.Name)
	}
	color.Unset()

	log.Println("parsing hosts file finished")

	color.Set(color.FgGreen)
	log.Println("portscanning hosts, this might take a really long time...")
	color.Unset()
	results := scan.Scan_hosts(hosts, cfg.Progress_Seconds)
	color.Set(color.FgGreen)
	log.Println("scanning hosts finished")
	color.Unset()

	if cfg.Services {
		log.Println("checking services...")
		check_services(results, services)
		color.Set(color.FgGreen)
		log.Println("finished checking services")
		color.Unset()
	}

	if cfg.Graphviz != "" {
		log.Println("writing graphviz dot file...")
		export.Write_graphviz(results, services, cfg.Graphviz)
		log.Println("finished writing graphviz dot file")
	}

	log.Println("finished")
}

func check_services(results []Host, services []ServiceDef) {

	for _, h := range results {
		if len(h.Ports) == 0 {
			log.Printf("[%v] no exposed services", h.Ip)
			continue
		}

		for _, p := range h.Ports {
			if p.State == "open" {
				s, err := util.Find_service(p.Number, fmt.Sprintf("%s", h.Ip), services)
				if err == nil {
					log.Printf("[%v] %v - %v (%v)",
						h.Ip,
						p.Number,
						s.Id,
						s.Description)
				} else {
					color.Set(color.FgRed)
					log.Printf("! [%v] %v %v: no service definition found (%v %v)",
						h.Ip,
						p.Number,
						p.State,
						p.Name,
						p.Version)
					color.Unset()
				}
			}
		}
	}
}

func parse_hosts(jsonFile io.Reader) []HostDef {

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var hosts []HostDef
	json.Unmarshal(byteValue, &hosts)
	return hosts
}

func print_service(service ServiceDef) {

	var ports []int
	for _, p := range service.Ports {
		ports = append(ports, p.Port)
	}

	log.Printf("Service: %v %v - %v", service.Id, ports, service.Description)
}

func parse_services(jsonFile io.Reader) []ServiceDef {

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var services []ServiceDef
	json.Unmarshal(byteValue, &services)
	return services
}
