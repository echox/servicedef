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

func Write_graphviz(hosts []Host, services []ServiceDef, file string) {

	g := dot.NewGraph(dot.Directed)
	internet := g.Node("internet")

	for _, h := range hosts {

		host_lbl := fmt.Sprintf("%v\n%v", h.Ip, h.Dns)
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
