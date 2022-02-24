package scan

import (
	"testing"

	. "github.com/echox/servicedef/result"
)

func build_hosts() []Host {

	return []Host{
		Host{Ip: "127.0.0.1"},
		Host{Ip: "192.168.1.1", Ports: []Port{Port{Number: 80}, Port{Number: 443}}},
		Host{Ip: "1.1.1.1", Ports: []Port{Port{Number: 8080}}},
		Host{Ip: "10.10.10.10", Ports: []Port{Port{Number: 8080}}},
	}
}

func TestContainsHost(t *testing.T) {

	hosts := build_hosts()

	var h *Host = containsHost(&hosts, Host{Ip: "1.1.1.1"})
	if h == nil || h.Ip != "1.1.1.1" {
		t.Errorf("Couldn't find host by IP")
	}

	var notFound *Host = containsHost(&hosts, Host{Ip: "192.168.2.2"})
	if notFound != nil {
		t.Errorf("Should not find host but got %v", notFound)
	}
}

func TestMergeHosts(t *testing.T) {

	hosts := build_hosts()

	mergeHosts(&hosts[1], hosts[3])

	if len(hosts[1].Ports) != 3 {
		t.Errorf("Merge failed, port numbers don't match")
	} else {
		for _, p := range hosts[1].Ports {
			switch p.Number {
			case 80:
			case 443:
			case 8080:
				continue
			default:
				t.Errorf("Merge failed, port numbers don't match, got %d", p.Number)
			}

		}
	}
}

func TestMergeHostsWithPortConflict(t *testing.T) {

	hosts := build_hosts()

	mergeHosts(&hosts[2], hosts[3])

	if len(hosts[2].Ports) != 1 {
		t.Errorf("Merge failed, port numbers don't match")
	} else {
		for _, p := range hosts[2].Ports {
			switch p.Number {
			case 8080:
				continue
			default:
				t.Errorf("Merge failed, port numbers don't match, got %d", p.Number)
			}

		}
	}
}
