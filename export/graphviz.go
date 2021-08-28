package export

import (
	"fmt"
	"io/ioutil"
	"log"

	. "github.com/echox/servicedef/model"
	"github.com/echox/servicedef/util"

	"github.com/emicklei/dot"
	"github.com/google/uuid"
)

func add_host_description(host_label string, host Host, definitions []HostDef) string {

	if host_def, err := util.Find_host(host.Ip, definitions); err == nil {

		if host_def.Description != "" {
			return fmt.Sprintf("%v\n%v", host_label, host_def.Description)
		}
	}

	return host_label
}

func Write_graphviz(hosts []Host, services []ServiceDef, hosts_def []HostDef, file string) {

	g := dot.NewGraph(dot.Directed)
	internet := g.Node("internet")

	for _, h := range hosts {

		if len(h.Ports) == 0 {
			continue
		}

		host_lbl := fmt.Sprintf("%v\nPTR:%v", h.Ip, h.Dns)
		host_lbl = add_host_description(host_lbl, h, hosts_def)
		host_node := g.Node(host_lbl)
		g.Edge(internet, host_node)

		for _, p := range h.Ports {
			if p.State == "open" {

				s, err := util.Find_service(p.Number, h.Ip, services)
				if err == nil {
					lbl := fmt.Sprintf("Port %v\n%v", p.Number, s.Id)
					service_node := g.Node(uuid.NewString()).Attr("label", lbl)
					g.Edge(host_node, service_node)

				} else {
					lbl := fmt.Sprintf("Port %v\n(%v) %v", p.Number, p.Name, p.Version)
					service_node := g.Node(uuid.NewString()).
						Attr("style", "filled").
						Attr("color", "red").
						Attr("label", lbl)
					g.Edge(host_node, service_node)
				}

			}
		}

	}

	b := []byte(g.String())
	fwerr := ioutil.WriteFile(file, b, 0644)
	if fwerr != nil {
		log.Println(fwerr)
	}
}
