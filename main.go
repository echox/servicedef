package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/echox/servicedef/config"
	. "github.com/echox/servicedef/definition"
	"github.com/echox/servicedef/export"
	. "github.com/echox/servicedef/result"
	"github.com/echox/servicedef/scan"

	"github.com/fatih/color"
	log "github.com/sirupsen/logrus"
)

func initLog(cfg config.Config) {

	if cfg.JSONLog {
		log.SetFormatter(&log.JSONFormatter{})
	} else {
		log.SetFormatter(&log.TextFormatter{
			FullTimestamp: true,
			DisableColors: true, // use custom formatter, see https://github.com/sirupsen/logrus/issues/1194
		})
	}

	if lvl, err := log.ParseLevel(cfg.Loglevel); err != nil {
		log.Warnf("Couldn't parse loglevel \"%s\" using info instead.", cfg.Loglevel)
		log.SetLevel(log.InfoLevel)
	} else {
		log.SetLevel(lvl)
	}

}

func main() {

	cfg, err := config.Init()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	initLog(cfg)

	log.Println("running servicedef v0")

	var (
		hosts    HostDefs
		services ServiceDefs
		rules    Rules
	)

	log.Tracef("parsing hosts file...")
	if json_error := hosts.Init(cfg.HostsPath); json_error != nil {
		log.Fatalf("parsing hosts error: %v", json_error)
	}
	color.Set(color.FgYellow)
	for _, h := range hosts {
		log.Debugf("Host: %v %v", h.Address, h.Description)
	}
	color.Unset()
	log.Tracef("parsing hosts file finished")
	log.Printf("host definitions: %v", len(hosts))

	if cfg.ServicesPath != "" {
		log.Tracef("parsing services file...")
		if json_error := services.Init(cfg.ServicesPath); json_error != nil {
			log.Fatalf("parsing services error: %v", json_error)
		}
		log.Printf("service definitions: %v", len(services))
		for _, s := range services {
			s.Print()
		}
		log.Tracef("parsing services file finished")
	} else {
		log.Println("no service definitions - scanning only")
	}

	if cfg.RulesPath == "" {
		log.Println("no rules supplied, use (-r) if needed")
	} else {
		log.Printf("loading rule file %v", cfg.RulesPath)

		if err := rules.Init(cfg.RulesPath); err != nil {
			log.Fatalf("error loading rules: %v", err)
		}
		log.Printf("Rules definitions: %v", len(rules))
	}

	color.Set(color.FgGreen)
	log.Println("portscanning hosts, this might take a really long time...")
	color.Unset()
	results := scan.ScanHosts(hosts, cfg)
	color.Set(color.FgGreen)
	log.Println("scanning hosts finished")
	color.Unset()

	// sort result list before exports to make things more comparable
	sort.SliceStable(results, func(i, j int) bool { return results[i].Ip < results[j].Ip })

	if len(services) != 0 {

		log.Println("checking services...")
		mapServices(results, services, rules)
		color.Set(color.FgGreen)
		log.Tracef("finished checking services")
		color.Unset()
	} else {
		log.Println("no services to check - only printing open ports...")
		results.PrintOpenPorts()
	}

	if cfg.ResultPath != "" {
		export.WriteJSON(results, cfg.ResultPath)
	}

	if cfg.Graphviz != "" {
		log.Println("writing graphviz dot file...")
		export.WriteGraphviz(results, services, hosts, cfg.Graphviz)
		log.Tracef("finished writing graphviz dot file")
	}

	log.Println("finished")
}

func mapServices(results ResultHosts, services ServiceDefs, rules []RulesDef) {

	servicesNotUsed := make(ServiceDefs, len(services))
	copy(servicesNotUsed, services)

	for i, h := range results {
		if len(h.Ports) == 0 {
			log.Printf("[%v] no exposed services", h.Ip)
			continue
		}

		for portIdx, p := range h.Ports {
			if p.State == "open" {
				h.Ports[portIdx], servicesNotUsed = mapPort(services, servicesNotUsed, rules, &p, &h)
			}
		}
		results[i] = h

	}

	color.Set(color.FgYellow)
	log.Printf("services in the catalog but not found during scan: %d", len(servicesNotUsed))
	color.Unset()
	for _, s := range servicesNotUsed {
		log.Printf("Service defined but not found: %s - %s", s.Id, s.Description)
	}
}

func mapPort(services ServiceDefs, servicesNotUsed ServiceDefs, rules Rules, port *Port, host *Host) (Port, ServiceDefs) {

	service, err := services.Find(port.Number, *host)
	if err == nil {
		port.ServiceId = service.Id
		servicesNotUsed = servicesNotUsed.Remove(service)
		log.Printf("[%v] %v - %v (%v)",
			host.Ip,
			port.Number,
			service.Id,
			service.Description)

		for _, pDef := range service.Ports {
			if len(pDef.Rules) != 0 && pDef.Port == port.Number && host.Inside(pDef.Hosts) {
				rules.Check(pDef, port, service, host.Ip)
			}
		}
	} else {
		color.Set(color.FgRed)
		log.Printf("! [%v] %v %v: no service definition found (%v %v)",
			host.Ip,
			port.Number,
			port.State,
			port.Name,
			port.Version)
		color.Unset()
	}

	return *port, servicesNotUsed

}
