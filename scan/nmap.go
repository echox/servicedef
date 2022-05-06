// Package scan contains the scanning functionality of servicedef. It uses the
// an array of host definitions HostDef created from the hosts.json file as
// input.
// Currently only scanning via nmap is supported and it is unlikely we will
// ever support other methods. Servicedef is currently strongly coupled with
// the nmap scanning mehtods and results.
// At the moment nmap is utilized via the github.com/Ullaakut/nmap package.
package scan

import (
	log "github.com/sirupsen/logrus"

	"sync"
	"time"

	"github.com/echox/servicedef/config"
	. "github.com/echox/servicedef/definition"
	. "github.com/echox/servicedef/result"

	"github.com/Ullaakut/nmap/v2"
	"github.com/fatih/color"
)

type ScannedCounter struct {
	mu      sync.Mutex
	counter int
	max     int
}

func containsHost(results *ResultHosts, host Host) *Host {
	for _, resultHost := range *results {
		if resultHost.Ip == host.Ip {
			return &resultHost
		}
	}
	return nil
}

func mergeHosts(target *Host, source Host) {
	for _, port := range source.Ports {
		if !portExists(target, port) {
			target.Ports = append(target.Ports, port)
		}
	}
	for _, tag := range source.Tags {
		if !tagExists(target, tag) {
			target.Tags = append(target.Tags, tag)
		}
	}
}

func tagExists(host *Host, tag string) bool {
	for _, current := range host.Tags {
		if current == tag {
			return true
		}
	}
	return false
}

func portExists(target *Host, port Port) bool {
	for _, target_port := range target.Ports {
		if port.Number == target_port.Number {
			return true
		}
	}
	return false
}

// Scan_hosts evaluates a given list of host definitions
// The configuration is needed for knowing if scanning progress should be
// logged
func ScanHosts(hosts []HostDef, cfg config.Config) ResultHosts {

	scanned := ScannedCounter{counter: 0, max: len(hosts)}
	p := make(chan HostDef, scanned.max)
	result_queue := make(chan Host, 10)

	var result_hosts ResultHosts

	for _, h := range hosts {
		p <- h
	}
	close(p)

	var wg_collector sync.WaitGroup
	wg_collector.Add(1)
	go func() {
		for result_host := range result_queue {
			existing := containsHost(&result_hosts, result_host)
			if existing != nil {
				mergeHosts(existing, result_host)
			} else {
				result_hosts = append(result_hosts, result_host)
			}
		}
		defer wg_collector.Done()
	}()

	var wg sync.WaitGroup

	for i := 1; i <= cfg.Threads; i++ {
		wg.Add(1)
		go scanHostWorker(i, p, &scanned, &wg, result_queue, cfg)
	}
	wg.Wait()
	close(result_queue)
	wg_collector.Wait()

	log.Printf("gathered %v hosts", len(result_hosts))

	return result_hosts

}

func scanHost(id int, host string, cfg config.Config) (*nmap.Run, error) {

	log.Printf("[worker_%v] scanning %v...", id, host)

	options := []nmap.Option{
		nmap.WithTargets(host),
		nmap.WithTimingTemplate(nmap.TimingAggressive),
		nmap.WithServiceInfo(),
		nmap.WithVerbosity(3),
	}

	if cfg.Connect_Scan {
		options = append(options, nmap.WithConnectScan())
	} else {
		options = append(options, nmap.WithSYNScan())
	}

	if !cfg.Default_Port_Scan {
		options = append(options, nmap.WithPorts("-"))
	}

	s, err := nmap.NewScanner(options...)
	if err != nil {
		log.Printf("[worker_%v] unable to create nmap scanner: %v", id, err)
		return nil, err
	}

	var result *nmap.Run
	var w []string
	var e error

	if cfg.Progress_Seconds > 0 {
		progress := make(chan float32, 1)
		ts := time.Now()
		go func() {
			for p := range progress {
				if time.Now().After(ts.Add(time.Duration(cfg.Progress_Seconds) * time.Second)) {
					ts = time.Now()
					log.Printf("[worker_%v] [%v] portscan progress: %v %%", id, host, p)
				}
			}
		}()

		result, w, e = s.RunWithProgress(progress)
	} else {
		result, w, e = s.Run()
	}

	if e != nil {
		color.Set(color.FgRed)
		log.Errorf("[worker_%v]unable to run nmap scan: %v", id, e)
		color.Unset()
		return nil, e
	}

	if w != nil {
		color.Set(color.FgYellow)
		log.Warnf("[worker_%v] Warnings: %v", id, w)
		color.Unset()
	}

	color.Set(color.FgGreen)
	log.Printf("[worker_%v] [%v] nmap done: %d hosts up scanned in %3f seconds", id, host, len(result.Hosts), result.Stats.Finished.Elapsed)
	color.Unset()

	return result, nil
}

func scanHostWorker(id int, pool chan HostDef, scanned *ScannedCounter, wg *sync.WaitGroup, result_queue chan Host, cfg config.Config) {

	for hostDef := range pool {
		if sr, err := scanHost(id, hostDef.Address, cfg); err == nil {
			hosts := parseNmap(sr)
			for _, h := range hosts {
				h.Tags = hostDef.Tags
				result_queue <- h
			}
		}
		scanned.mu.Lock()
		scanned.counter++
		scanned.mu.Unlock()
		log.Printf("%d/%d definitions scanned", scanned.counter, scanned.max)
	}

	log.Printf("[worker_%v] finished queue", id)
	defer wg.Done()
}

func parseNmapHost(h nmap.Host) Host {

	var parsed_host Host

	// just use the first adress and hostname if available
	var hostname = ""
	if len(h.Hostnames) != 0 {
		hostname = h.Hostnames[0].Name
	}
	parsed_host.Dns = hostname

	if len(h.Addresses) != 0 {
		parsed_host.Ip = h.Addresses[0].Addr
	}

	if len(h.Ports) != 0 {

		for _, p := range h.Ports {
			var parsed_port Port
			parsed_port.Number = int(p.ID)
			parsed_port.State = p.State.State
			parsed_port.Version = p.Service.Version
			parsed_port.Name = p.Service.Name
			//TODO refactor to constructor
			parsed_port.RuleResults = make(map[string]bool)

			parsed_host.Ports = append(parsed_host.Ports, parsed_port)
		}

	}

	return parsed_host
}

func parseNmap(scan *nmap.Run) []Host {

	var hosts []Host

	for _, h := range scan.Hosts {

		parsed_host := parseNmapHost(h)
		if parsed_host.Ip != "" {
			hosts = append(hosts, parsed_host)
		}

		if len(h.Ports) == 0 || len(h.Addresses) == 0 {
			continue
		}
	}

	return hosts
}
