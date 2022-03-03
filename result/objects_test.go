package result

import (
	"testing"
)

func TestInside(t *testing.T) {

	h := Host{Ip: "127.0.0.1", Dns: "localhost", Tags: []string{"someTag", "anotherTag"}}
	listIp := []string{"127.0.0.2", "someHostname", "127.0.0.1", "tag:someHost"}
	listDNS := []string{"someHost", "192.168.1.1", "tag:someHost", "localhost"}
	listNoMatch := []string{"192.168.1.1", "tag:tag", "someHostname", "127.0.0.2"}
	listTag := []string{"192.168.1.1", "someHostname", "127.0.0.2", "tag:anotherTag"}

	if !h.Inside(listIp) {
		t.Errorf("%v should be found in %v", h, listIp)
	}
	if !h.Inside(listDNS) {
		t.Errorf("%v should be found in %v", h, listDNS)
	}
	if h.Inside(listNoMatch) {
		t.Errorf("%v should not be found in %v", h, listNoMatch)
	}
	if !h.Inside(listTag) {
		t.Errorf("%v should be found in %v", h, listTag)
	}
}
