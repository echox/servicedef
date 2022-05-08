package definition

import (
	. "github.com/echox/servicedef/result"

	"github.com/fatih/color"
	log "github.com/sirupsen/logrus"
)

// map serviceIds to the scan results
func MapServices(results ResultHosts, services ServiceDefs, rules []RulesDef) {

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
