package model

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"strings"
)

// Collection of service definitions from json file
type ServiceDefs []ServiceDef

// Init loads the ServiceDefs form a json file
func (defs *ServiceDefs) Init(jsonFile io.Reader) error {

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return err
	}
	if json_error := json.Unmarshal(byteValue, defs); json_error != nil {
		return json_error
	}
	return nil
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

func (service *ServiceDef) Print() {

	var ports []int
	for _, p := range service.Ports {
		ports = append(ports, p.Port)
	}

	log.Printf("Service: %v %v - %v", service.Id, ports, service.Description)
}
