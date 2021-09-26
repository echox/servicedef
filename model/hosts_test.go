package model

import (
	"testing"
)

func build_hostdefs() HostDefs {

	return HostDefs{
		HostDef{"172.16.1.0/24", "test a", []string{}},
		HostDef{"172.16.2.0/24", "test b", []string{}},
		HostDef{"172.16.3.1", "test c", []string{}},
		HostDef{"172.16.3.2", "test d", []string{}},
	}
}

func TestHostFind(t *testing.T) {

	hosts := build_hostdefs()

	// test host in subnet
	if h, err := hosts.Find("172.16.1.1"); err == nil {
		if h.Address != "172.16.1.0/24" {
			t.Errorf("Expected %v but got %v", "172.16.1.0/24", h)
		}
	} else {
		t.Errorf("Couldn't find %v: %v", "172.16.1.1", err)

	}

	// test host not inside list
	if h, err := hosts.Find("127.0.0.1"); err == nil {
		t.Errorf("Expected error but got %v", h)
	}

	// broken ip
	if h, err := hosts.Find("127.www0.0.1"); err == nil {
		t.Errorf("Expected error but got %v", h)
	}

	// test host directly
	if h, err := hosts.Find("172.16.3.2"); err == nil {
		if h.Address != "172.16.3.2" {
			t.Errorf("Expected %v but got %v", "172.16.3.2", h)
		}
	} else {
		t.Errorf("Couldn't find %v: %v", "172.16.3.2", err)

	}

}
