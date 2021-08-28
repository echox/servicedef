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

	// test host directly
	if h, err := Find_host("172.16.3.2", hosts); err == nil {
		if h.Address != "172.16.3.2" {
			t.Errorf("Expected %v but got %v", "172.16.3.2", h)
		}
	} else {
		t.Errorf("Couldn't find %v: %v", "172.16.3.2", err)

	}

}
