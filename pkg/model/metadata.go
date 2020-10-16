package model

// Metadata model
type Metadata struct {
	IP       string `json:"ip" example:"192.168.1.100"`
	Hostname string `json:"hostname" example:"dobby"`
}
