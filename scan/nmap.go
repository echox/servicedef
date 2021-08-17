package scan

import (
	"log"
	"sync"
	"time"

	. "github.com/echox/servicedef/model"

	"github.com/Ullaakut/nmap/v2"
	"github.com/fatih/color"
)

func Scan_hosts(hosts []HostDef) []Host {

	p := make(chan string, len(hosts))
	result_queue := make(chan Host, 10)

	var result_hosts []Host

	for _, h := range hosts {
		p <- h.Ip
	}
	close(p)

	var wg_collector sync.WaitGroup
	wg_collector.Add(1)
	go func() {
		for result_host := range result_queue {
			result_hosts = append(result_hosts, result_host)
		}
		defer wg_collector.Done()
	}()

	var wg sync.WaitGroup
	m := &sync.Mutex{}

	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go scan_host_worker(i, p, &wg, m, result_queue)
	}
	wg.Wait()
	close(result_queue)
	wg_collector.Wait()

	log.Printf("gathered %v hosts", len(result_hosts))

	return result_hosts

}

func scan_host(id int, host string) *nmap.Run {

	log.Printf("[worker_%v] scanning %v...", id, host)

	s, err := nmap.NewScanner(
		nmap.WithTargets(host),
		nmap.WithTimingTemplate(nmap.TimingAggressive),
		nmap.WithServiceInfo(),
		nmap.WithSYNScan(),
		//nmap.WithPorts("-"),
		nmap.WithVerbosity(3),
		nmap.WithFastMode(),
	)

	if err != nil {
		log.Fatalf("[worker_%v] unable to create nmap scanner: %v", id, err)
	}

	progress := make(chan float32, 1)

	ts := time.Now()

	go func() {
		for p := range progress {
			if time.Now().After(ts.Add(60 * time.Second)) {
				ts = time.Now()
				log.Printf("[worker_%v] portscan progress: %v %%", id, p)
			}
		}
	}()

	result, w, e := s.RunWithProgress(progress)
	if e != nil {
		color.Set(color.FgRed)
		log.Fatalf("[worker_%v]unable to run nmap scan: %v", id, e)
		color.Unset()
	}

	if w != nil {
		color.Set(color.FgYellow)
		log.Printf("[worker_%v] Warnings: \n %v", id, w)
		color.Unset()
	}

	color.Set(color.FgGreen)
	log.Printf("[worker_%v] [%v] nmap done: %d hosts up scanned in %3f seconds\n", id, host, len(result.Hosts), result.Stats.Finished.Elapsed)
	color.Unset()

	return result
}

func scan_host_worker(id int, pool chan string, wg *sync.WaitGroup, m *sync.Mutex, result_queue chan Host) {

	for ip := range pool {
		sr := scan_host(id, ip)
		parse_nmap(sr, result_queue)
		m.Lock()
		m.Unlock()
	}

	log.Printf("[worker_%v] finished queue", id)
	defer wg.Done()
}

func parse_nmap_host(h nmap.Host) Host {

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

			parsed_host.Ports = append(parsed_host.Ports, parsed_port)
		}

	}

	return parsed_host
}

func parse_nmap(scan *nmap.Run, result_queue chan Host) []Host {

	var hosts []Host

	for _, h := range scan.Hosts {

		parsed_host := parse_nmap_host(h)
		if parsed_host.Ip != "" {
			result_queue <- parsed_host
		}

		if len(h.Ports) == 0 || len(h.Addresses) == 0 {
			continue
		}
	}

	return hosts
}
