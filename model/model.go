// Package model holds the different types, for example the definitions of the
// services, hosts or the scanning results
package model

// ServiceDef of a service
type ServiceDef struct {
	Id          string    `json:"id"`          // Identifier
	Description string    `json:"description"` //Short text description of the services purpose
	Ports       []PortDef `json:"ports"`       // Port definitions of the service
}

// PortDef for a service
type PortDef struct {
	Port        int      `json:"port"` // Port used by the service
	Protocol    string   `json:"protocol"`
	Description string   `json:"description"` // Short text description what is handled at the specific port
	Uri         string   `json:"uri"`         // Optional URI for reaching the service at the given port
	Rules       []string `json:"rules"`       // Rules which should be fulfilled by the service
	Hosts       []string `json:"hosts"`       // Hosts where the service is deployed. tags:tag can also be used
}

// HostDef of the inventory
type HostDef struct {
	Address     string   `json:"address"`     // Adress, hostname or subnet in CIDR notation
	Description string   `json:"description"` // Short text description of the hosts purpose
	Tags        []string `json:"tags"`        // Tags for referencing a service
}

// RulesDef collection which can be referenced as Rules within a PortDef
type RulesDef struct {
	Name   string    `json:"name"`   // Identifier of the rule
	Type_  string    `json:"type"`   // Type, currently 'http' is the only supported type
	Status int       `json:"status"` // Optional, expected http status
	Rules  []RuleDef `json:"rules"`  // Optional, Rules which must match
}

// RuleDef which must be matched by a service
type RuleDef struct {
	Name     string `json:"name"`     // HTTP header with given name must exist
	Contains string `json:"contains"` // Optional, HTTP header must contain given value
}

// Host scanning result
type Host struct {
	Ip    string   // IP address of the host
	Dns   string   // Reverse DNS name of the host
	Ports []Port   // Available ports
	Tags  []string // Tags of the host
}

// Port represents an discovered port on a scanned host
type Port struct {
	Number       int             // Numeric port representation
	Version      string          // Optional, version string if discovered
	State        string          // State of port either open, closed or filtered
	Name         string          // Name according to /etc/services
	Rule_Results map[string]bool // Rule evaluation results
}
