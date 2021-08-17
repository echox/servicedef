package model

type ServiceDef struct {
	Id          string    `json:"id"`
	Type_       string    `json:"type"`
	Description string    `json:"description"`
	Ports       []PortDef `json:"ports"`
}

type PortDef struct {
	Port          int      `json:"port"`
	Type_         string   `json:"type"`
	Protocol      string   `json:"protocol"`
	Description   string   `json:"description"`
	Auth          bool     `json:"auth"`
	Auth_provider string   `json:"auth-provider"`
	Uri           string   `json:"uri"`
	Hosts         []string `json:"hosts"`
}

type Auth_service struct {
	Id           string `json:"id"`
	Type_        string `json:"type"`
	Auth_type    string `json:"auth-type"`
	Description  string `json:"description"`
	Provider_uri string `json:"provider-uri"`
}

type HostDef struct {
	Ip   string `json:"ip"`
	Name string `json:"name"`
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
	Number  int
	Version string
	State   string
	Name    string
}
