package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"strings"

	"net/http"

	"github.com/echox/servicedef/config"
	"github.com/echox/servicedef/export"
	. "github.com/echox/servicedef/model"
	"github.com/echox/servicedef/scan"
	"github.com/echox/servicedef/util"

	"github.com/fatih/color"
)

func main() {
	log.Println("running servicedef v0")

	cfg := config.Init()

	if cfg.Quiet {
		log.SetOutput(ioutil.Discard)
	}

	var rules []RulesDef
	if cfg.Rules == "" {
		log.Println("no rules supplied, use (-r) if needed")
	} else {
		rules = loadRules(cfg)
	}

	var services []ServiceDef
	if cfg.Services {
		log.Println("parsing services file...")
		var json_error error
		services, json_error = parse_services(cfg.Services_File)
		if json_error != nil {
			log.Printf("parsing services error: %v", json_error)
			return
		}
		log.Printf("Services #: %v", len(services))
		for _, s := range services {
			print_service(s)
		}
		log.Println("parsing services file finished")
	} else {
		log.Println("no service definitions - scanning only")
	}

	log.Println("parsing hosts file...")

	var json_error error
	hosts, json_error := parse_hosts(cfg.Hosts_File)
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

func loadRules(cfg config.Config) []RulesDef {

	log.Printf("loading rule file %v", cfg.Rules)

	var rules []RulesDef
	byteValue, err := ioutil.ReadAll(cfg.Rules_File)
	if err != nil {
		log.Printf("error loading rules: %v", err)
		return rules
	}

	json.Unmarshal(byteValue, &rules)
	log.Printf("Rules #: %v", len(rules))

	return rules
}

func check_services(results []Host, services []ServiceDef, rules []RulesDef) {

	for _, h := range results {
		if len(h.Ports) == 0 {
			log.Printf("[%v] no exposed services", h.Ip)
			continue
		}

		for _, p := range h.Ports {
			if p.State == "open" {
				s, err := util.Find_service(p.Number, h, services)
				if err == nil {
					log.Printf("[%v] %v - %v (%v)",
						h.Ip,
						p.Number,
						s.Id,
						s.Description)

					for _, pDef := range s.Ports {
						if len(pDef.Rules) != 0 && pDef.Port == p.Number {
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

func parse_hosts(jsonFile io.Reader) ([]HostDef, error) {

	var hosts []HostDef
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return hosts, err
	}
	if json_error := json.Unmarshal(byteValue, &hosts); json_error != nil {
		return hosts, json_error
	}
	return hosts, nil
}

func print_service(service ServiceDef) {

	var ports []int
	for _, p := range service.Ports {
		ports = append(ports, p.Port)
	}

	log.Printf("Service: %v %v - %v", service.Id, ports, service.Description)
}

func parse_services(jsonFile io.Reader) ([]ServiceDef, error) {

	var services []ServiceDef
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return services, err
	}
	if json_error := json.Unmarshal(byteValue, &services); json_error != nil {
		return services, json_error
	}
	return services, nil
}
