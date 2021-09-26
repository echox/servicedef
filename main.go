package main

import (
	"io/ioutil"
	"log"
	"strings"

	"net/http"

	"github.com/echox/servicedef/config"
	"github.com/echox/servicedef/export"
	. "github.com/echox/servicedef/model"
	"github.com/echox/servicedef/scan"

	"github.com/fatih/color"
)

func main() {
	log.Println("running servicedef v0")

	cfg := config.Init()

	if cfg.Quiet {
		log.SetOutput(ioutil.Discard)
	}

	var rules Rules
	if cfg.Rules == "" {
		log.Println("no rules supplied, use (-r) if needed")
	} else {
		if err := rules.Init(cfg); err != nil {
			log.Printf("error loading rules: %v", err)
		}
	}

	var services ServiceDefs
	if cfg.Services {
		log.Println("parsing services file...")
		var json_error error
		json_error = services.Init(cfg.Services_File)
		if json_error != nil {
			log.Printf("parsing services error: %v", json_error)
			return
		}
		log.Printf("Services #: %v", len(services))
		for _, s := range services {
			s.Print()
		}
		log.Println("parsing services file finished")
	} else {
		log.Println("no service definitions - scanning only")
	}

	log.Println("parsing hosts file...")

	var hosts HostDefs
	json_error := hosts.Init(cfg.Hosts_File)
	if json_error != nil {
		log.Printf("parsing hosts error: %v", json_error)
		return
	}
	color.Set(color.FgYellow)
	for _, h := range hosts {
		log.Printf("Host: %v %v", h.Address, h.Description)
	}
	color.Unset()

	log.Println("parsing hosts file finished")

	color.Set(color.FgGreen)
	log.Println("portscanning hosts, this might take a really long time...")
	color.Unset()
	results := scan.Scan_hosts(hosts, cfg)
	color.Set(color.FgGreen)
	log.Println("scanning hosts finished")
	color.Unset()

	if cfg.Services {

		log.Println("checking services...")
		check_services(results, services, rules)
		color.Set(color.FgGreen)
		log.Println("finished checking services")
		color.Unset()
	}

	if cfg.Graphviz != "" {
		log.Println("writing graphviz dot file...")
		export.Write_graphviz(results, services, hosts, cfg.Graphviz)
		log.Println("finished writing graphviz dot file")
	}

	log.Println("finished")
}

func contains(hosts []string, host Host) bool {
	for _, h := range hosts {
		if h == host.Ip || h == host.Dns {
			return true
		}
	}
	return false
}

func check_services(results []Host, services ServiceDefs, rules []RulesDef) {

	for _, h := range results {
		if len(h.Ports) == 0 {
			log.Printf("[%v] no exposed services", h.Ip)
			continue
		}

		for _, p := range h.Ports {
			if p.State == "open" {
				s, err := services.Find(p.Number, h)
				if err == nil {
					log.Printf("[%v] %v - %v (%v)",
						h.Ip,
						p.Number,
						s.Id,
						s.Description)

					for _, pDef := range s.Ports {
						if len(pDef.Rules) != 0 && pDef.Port == p.Number && contains(pDef.Hosts, h) {
							check_rules(rules, pDef, &p, s, h.Ip)
						}
					}
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

func check_rules(rules []RulesDef, port PortDef, port_result *Port, service ServiceDef, ip string) {

	for _, pRule := range port.Rules {
		for _, r := range rules {
			if r.Name == pRule {
				if r.Type_ == "http" {
					s := eval_http(r, port.Uri)
					port_result.Rule_Results[r.Name] = s
					if s == false {
						color.Set(color.FgRed)
						log.Printf("! [%v] rule %v doesn't match %v", ip, r.Name, port.Uri)
						color.Unset()
					}
				}
			}
		}
	}
}

func eval_http(rules RulesDef, uri string) bool {

	if uri == "" {
		log.Printf("URI needed for checking %v", rules.Name)
		return false
	}

	hc := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	r, err := hc.Get(uri)
	if err != nil {
		log.Println(err)
		return false
	}
	defer r.Body.Close()
	//body, err := io.ReadAll(r.Body)

	if rules.Status != 0 {

		if r.StatusCode != rules.Status {
			log.Printf("! [%v] Status doesn't match. Expected %v got %v", uri, rules.Status, r.Status)
			return false
		}
	}

	for _, rule := range rules.Rules {
		v := r.Header.Get(rule.Name)
		if v == "" {
			return false
		} else {
			if rule.Contains == "" {
				log.Printf("rule %v on %v matches header-rule '%v'", rules.Name, uri, rule.Name)
				continue
			} else {
				if strings.Contains(v, rule.Contains) {
					log.Printf("[%v] matches %v", uri, rule.Contains)
					continue
				} else {
					log.Printf("! [%v] Header mismatch for %v. Expected %v got %v", uri, rules.Name, rule.Contains, v)
					return false
				}

			}
		}

	}

	return true
}
