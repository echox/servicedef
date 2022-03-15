package definition

import (
	"testing"

	. "github.com/echox/servicedef/result"
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
		ServiceDef{Id: "Proxy", Description: "proxy test description",
			Ports: []PortDef{
				PortDef{Port: 8080,
					Hosts: []string{"tag:proxy"},
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

	//tagged host
	if s, err := services.Find(8080, Host{Tags: []string{"proxy"}}); err == nil {
		if s.Id != "Proxy" {
			t.Errorf("Should find 'proxy' service but got %s", s.Id)
		}
	} else {
		t.Errorf("Should find 'proxy' service but got %v", err)
	}

	//tag does not match
	if s, err := services.Find(8080, Host{Tags: []string{"proxyy"}}); err == nil {
		t.Errorf("Should not find service but got %s", s.Id)
	}
}

func TestServiceRemoveLast(t *testing.T) {

	services := build_servicedefs()

	proxy := ServiceDef{Id: "Proxy"}
	services = services.Remove(proxy)

	for _, s := range services {
		if s.Id == proxy.Id {
			t.Errorf("Remove(Id=%s) did fail", proxy.Id)
		}
	}

	if _, err := services.FindById("Proxy"); err == nil {
			t.Errorf("Proxy should be there")
	}

	if _, err := services.FindById("HTTP Service"); err != nil {
			t.Errorf("HTTP Service should be there")
	}

	if _, err := services.FindById("SSH"); err != nil {
			t.Errorf("SSH should be there")
	}

	if len(services) != 2 {
			t.Errorf("services total should be 2 but is %v", len(services))
	}
}

func TestServiceRemoveFirst(t *testing.T) {

	services := build_servicedefs()

	http := ServiceDef{Id: "HTTP Service"}
	services = services.Remove(http)

	for _, s := range services {
		if s.Id == http.Id {
			t.Errorf("Remove(Id=%s) did fail", http.Id)
		}
	}

	if _, err := services.FindById("HTTP Service"); err == nil {
			t.Errorf("HTTP Service shouldn't be there")
	}

	if _, err := services.FindById("SSH"); err != nil {
			t.Errorf("SSH should be there")
	}

	if _, err := services.FindById("Proxy"); err != nil {
			t.Errorf("Proxy should be there")
	}
	if len(services) != 2 {
			t.Errorf("services total should be 2 but is %v", len(services))
	}
}
