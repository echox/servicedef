package util

import (
	. "github.com/echox/servicedef/model"
	"testing"
)

func build_hostdefs() []HostDef {

	return []HostDef{
		HostDef{"172.16.1.0/24", "test a"},
		HostDef{"172.16.2.0/24", "test b"},
		HostDef{"172.16.3.1", "test c"},
		HostDef{"172.16.3.2", "test d"},
	}
}

func TestFind_host(t *testing.T) {

	hosts := build_hostdefs()

	// test host in subnet
	if h, err := Find_host("172.16.1.1", hosts); err == nil {
		if h.Address != "172.16.1.0/24" {
			t.Errorf("Expected %v but got %v", "172.16.1.0/24", h)
		}
	} else {
		t.Errorf("Couldn't find %v: %v", "172.16.1.1", err)

	}

	// test host not inside list
	if h, err := Find_host("127.0.0.1", hosts); err == nil {
		t.Errorf("Expected error but got %v", h)
	}

	// broken ip
	if h, err := Find_host("127.www0.0.1", hosts); err == nil {
		t.Errorf("Expected error but got %v", h)
	}

	// test host directly
	if h, err := Find_host("172.16.3.2", hosts); err == nil {
		if h.Address != "172.16.3.2" {
			t.Errorf("Expected %v but got %v", "172.16.3.2", h)
		}
	} else {
		t.Errorf("Couldn't find %v: %v", "172.16.3.2", err)

	}

}

func build_servicedefs() []ServiceDef {

	return []ServiceDef{
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

func TestFind_service(t *testing.T) {

	services := build_servicedefs()

	//matching port and ip
	if s, err := Find_service(80, Host{Ip: "172.16.1.1"}, services); err == nil {
		if s.Id != "HTTP Service" {
			t.Errorf("Couldn't find service for %v, %v got %v", 80, "172.16.1.1", s.Id)
		}

	} else {
		t.Errorf("Couldn't find service for %v, %v %v", 80, "172.16.1.1", err)
	}

	//matching host only
	if s, err := Find_service(81, Host{Ip: "172.16.1.1"}, services); err == nil {
		t.Errorf("Shouldn't find service but got %v", s)
	}

	//matching port only
	if s, err := Find_service(80, Host{Ip: "172.16.1.0"}, services); err == nil {
		t.Errorf("Shouldn't find service but got %v", s)
	}

	//matching dns
	if s, err := Find_service(80, Host{Dns: "some.domain"}, services); err != nil {
		t.Errorf("Should find service but got %v", s)
	}

	//not matching dns
	if s, err := Find_service(80, Host{Dns: "some.unknown.domain"}, services); err == nil {
		t.Errorf("Shouldn't find service but got %v", s)
	}

	//zero port, empty ip
	if s, err := Find_service(0, Host{Ip: ""}, services); err == nil {
		t.Errorf("Shouldn't find service but got %v", s)
	}

}
