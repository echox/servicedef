package model

type ServiceDef struct {
	Id          string    `json:"id"`
	Type_       string    `json:"type"`
	Description string    `json:"description"`
	Ports       []PortDef `json:"ports"`
}

type PortDef struct {
	Port        int      `json:"port"`
	Type_       string   `json:"type"`
	Protocol    string   `json:"protocol"`
	Description string   `json:"description"`
	Uri         string   `json:"uri"`
	Rules       []string `json:"rules"`
	Hosts       []string `json:"hosts"`
}

type HostDef struct {
	Address     string `json:"address"`
	Description string `json:"description"`
}

type RulesDef struct {
	Name   string    `json:"name"`
	Type_  string    `json:"type"`
	Status int       `json:"status"`
	Rules  []RuleDef `json:"rules"`
}

type RuleDef struct {
	Name     string `json:"name"`
	Contains string `json:"value"`
}

type Host struct {
	Ip    string
	Dns   string
	Ports []Port
}

type IPort interface {
	IsOpen() bool
}

type Port struct {
	Number       int
	Version      string
	State        string
	Name         string
	Rule_Results map[string]bool
}
