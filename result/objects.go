// Packge result contains the types used to represent the finshed evaluation
package result

import (
	"strings"

	"github.com/fatih/color"
	log "github.com/sirupsen/logrus"
)

// Array of scan results
type ResultHosts []Host

// Host scan result
type Host struct {
	Ip    string   // IP address of the host
	Dns   string   // Reverse DNS name of the host
	Ports []Port   // Available ports
	Tags  []string // Tags of the host
}

// Port represents a discovered port on a scanned host
type Port struct {
	Number       int             // Numeric port representation
	Version      string          // Optional, version string if discovered
	State        string          // State of port either open, closed or filtered
	Name         string          // Name according to /etc/services
	Rule_Results map[string]bool // Rule evaluation results
}

func (results ResultHosts) PrintOpenPorts() {

	for _, h := range results {
		if len(h.Ports) == 0 {
			log.Printf("[%v] no exposed ports", h.Ip)
			continue
		}

		for _, p := range h.Ports {
			if p.State == "open" {
				color.Set(color.FgRed)
				log.Printf("[%v] %v %v: %v %v",
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

func (host Host) Inside(list []string) bool {

	for _, item := range list {
		if item == host.Ip || item == host.Dns {
			return true
		} else if strings.HasPrefix(item, "tag:") {
			for _, tag := range host.Tags {
				if "tag:"+tag == item {
					return true
				}
			}
		}
	}

	return false
}
