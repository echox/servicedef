package definition

import (
	. "github.com/echox/servicedef/result"
	"testing"
)

func mockResults() ResultHosts {
	return ResultHosts{
		Host{Ip: "10.20.20.20",
			Dns:   "Host1",
			Ports: []Port{},
			Tags:  []string{"host1"},
		},
		Host{Ip: "10.1.1.42",
			Dns: "Host2",
			Ports: []Port{
				Port{80, "", "open", "", nil, ""},
			},
			Tags: []string{"host2"},
		},
		Host{Ip: "10.1.1.43",
			Dns: "Host3",
			Ports: []Port{
				Port{80, "", "open", "", nil, ""},
				Port{8080, "", "open", "", nil, ""},
			},
			Tags: []string{"host3"},
		},
		Host{Ip: "10.1.1.44",
			Dns: "Host4",
			Ports: []Port{
				Port{80, "", "open", "", nil, ""},
				Port{8080, "", "open", "", nil, ""},
			},
			Tags: nil,
		},
	}
}

func mockServiceDefinitions() ServiceDefs {
	return ServiceDefs{
		ServiceDef{"dummyservice", "unused", []PortDef{PortDef{80, "http", "unused", "uri", nil, []string{"tag:unused", "10.10.10.10"}}}},
		ServiceDef{"webservice1", "description", []PortDef{PortDef{80, "http", "description", "uri", nil, []string{"tag:host3"}}}},
		ServiceDef{"webservice2", "description", []PortDef{PortDef{8080, "http", "description", "uri", nil, []string{"10.1.1.44"}}}},
	}

}

func TestMapServices(t *testing.T) {

	results := mockResults()
	services := mockServiceDefinitions()

	MapServices(results, services, make([]RulesDef, 0))

	host3 := results[2]
	if host3.Ports[0].ServiceId != "webservice1" {
		t.Errorf("could not map webservice1")
	}
	if host3.Ports[1].ServiceId != "" {
		t.Errorf("port should not be mapped but is %s", host3.Ports[1].ServiceId)
	}

	host4 := results[3]
	if host4.Ports[1].ServiceId != "webservice2" {
		t.Errorf("could not map webservice2")
	}
}
