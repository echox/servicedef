package model

import (
	"testing"
)

func build_servicedefs() ServiceDefs {

	return ServiceDefs{
		ServiceDef{Id: "HTTP Service", Description: "http service test description",
			Ports: []PortDef{
				PortDef{Port: 80,
					Hosts: []string{"172.16.1.1", "172.16.1.2", "some.domain"},
				},
			}},
		ServiceDef{Id: "SSH", Description: "ssh test description",
			Ports: []PortDef{
				PortDef{Port: 22,
					Hosts: []string{"172.16.1.1"},
				},
			}},
	}
}

func TestServiceFind(t *testing.T) {

	services := build_servicedefs()

	//matching port and ip
	if s, err := services.Find(80, Host{Ip: "172.16.1.1"}); err == nil {
		if s.Id != "HTTP Service" {
			t.Errorf("Couldn't find service for %v, %v got %v", 80, "172.16.1.1", s.Id)
		}

	} else {
		t.Errorf("Couldn't find service for %v, %v %v", 80, "172.16.1.1", err)
	}

	//matching host only
	if s, err := services.Find(81, Host{Ip: "172.16.1.1"}); err == nil {
		t.Errorf("Shouldn't find service but got %v", s)
	}

	//matching port only
	if s, err := services.Find(80, Host{Ip: "172.16.1.0"}); err == nil {
		t.Errorf("Shouldn't find service but got %v", s)
	}

	//matching dns
	if s, err := services.Find(80, Host{Dns: "some.domain"}); err != nil {
		t.Errorf("Should find service but got %v", s)
	}

	//not matching dns
	if s, err := services.Find(80, Host{Dns: "some.unknown.domain"}); err == nil {
		t.Errorf("Shouldn't find service but got %v", s)
	}

	//zero port, empty ip
	if s, err := services.Find(0, Host{Ip: ""}); err == nil {
		t.Errorf("Shouldn't find service but got %v", s)
	}

}
