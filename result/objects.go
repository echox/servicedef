// Packge result contains the types used to represent the finshed evaluation
package result

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
