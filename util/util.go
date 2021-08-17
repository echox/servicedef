package util

import (
	"github.com/echox/servicedef/model"

	"errors"
)

func Find_service(port int, host string, services []model.ServiceDef) (model.ServiceDef, error) {

	var r model.ServiceDef
	for _, s := range services {
		for _, p := range s.Ports {
			if p.Port == port {
				for _, h := range p.Hosts {
					if host == h {
						return s, nil
					}
				}
			}
		}
	}

	return r, errors.New("NOT_FOUND")
}
