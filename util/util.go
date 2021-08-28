package util

import (
	"errors"
	"github.com/echox/servicedef/model"
	"net"
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

func Find_host(host string, hosts []model.HostDef) (model.HostDef, error) {

	var h model.HostDef

	//TODO refactor to allow hostnames
	if ip := net.ParseIP(host); ip != nil {

		for _, h := range hosts {

			if h.Address == host {
				return h, nil
			} else {
				if _, cidr, cidr_err := net.ParseCIDR(h.Address); cidr_err == nil {
					if cidr != nil && cidr.Contains(ip) {
						return h, nil
					}
				}
			}
		}

	} else {
		return h, errors.New("NO_IP")
	}

	return h, errors.New("NOT_FOUND")
}
