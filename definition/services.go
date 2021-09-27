package definition

import (
	"errors"
	"log"
	"strings"

	. "github.com/echox/servicedef/result"
)

// ServiceDef of a service
type ServiceDef struct {
	Id          string    `json:"id"`          // Identifier
	Description string    `json:"description"` //Short text description of the services purpose
	Ports       []PortDef `json:"ports"`       // Port definitions of the service
}

// PortDef for a service
type PortDef struct {
	Port        int      `json:"port"` // Port used by the service
	Protocol    string   `json:"protocol"`
	Description string   `json:"description"` // Short text description what is handled at the specific port
	Uri         string   `json:"uri"`         // Optional URI for reaching the service at the given port
	Rules       []string `json:"rules"`       // Rules which should be fulfilled by the service
	Hosts       []string `json:"hosts"`       // Hosts where the service is deployed. tags:tag can also be used
}

// Collection of service definitions from json file
type ServiceDefs []ServiceDef

// Init loads the ServiceDefs form a json file
func (defs *ServiceDefs) Init(servicesPath string) error {
	return parseJSONFile(servicesPath, defs)
}

// Find returns a service definition for a given host and port or an error
func (defs *ServiceDefs) Find(port int, host Host) (ServiceDef, error) {

	var r ServiceDef

	for _, s := range *defs {
		for _, p := range s.Ports {
			if p.Port == port {
				for _, h := range p.Hosts {
					if strings.HasPrefix(h, "tag:") {
						v := strings.Split(h, "tag:")
						for _, tag := range host.Tags {
							if v[1] == tag {
								return s, nil
							}
						}

					} else {
						if h == host.Dns || h == host.Ip {
							return s, nil
						}
					}
				}
			}
		}
	}

	return r, errors.New("NOT_FOUND")
}

// Print in a loggin friendly way
func (service *ServiceDef) Print() {

	var ports []int
	for _, p := range service.Ports {
		ports = append(ports, p.Port)
	}

	log.Printf("Service: %v %v - %v", service.Id, ports, service.Description)
}
