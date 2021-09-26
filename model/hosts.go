package model

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net"
)

// Collection of host definitions from json file
type HostDefs []HostDef

// Init loads the HostDefs form a json file
func (defs *HostDefs) Init(jsonFile io.Reader) error {

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return err
	}
	if json_error := json.Unmarshal(byteValue, defs); json_error != nil {
		return json_error
	}
	return nil
}

// Find returns a host definition by name or an error
func (defs *HostDefs) Find(host string) (HostDef, error) {

	var h HostDef

	//TODO refactor to allow hostnames
	if ip := net.ParseIP(host); ip != nil {

		for _, h := range *defs {

			if h.Address == host {
				return h, nil
			} else {
				if _, cidr, cidr_err := net.ParseCIDR(h.Address); cidr_err == nil {
					if cidr != nil && cidr.Contains(ip) {
						return h, nil
					}
				}
			}
		}

	} else {
		return h, errors.New("NO_IP")
	}

	return h, errors.New("NOT_FOUND")
}
