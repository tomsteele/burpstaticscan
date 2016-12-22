package main

// PassiveScanRequest is used to submit a JSON to burpbuddy.
type PassiveScanRequest struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	UseHTTPS bool   `json:"useHttps"`
	Request  string `json:"request"`
	Response string `json:"response"`
}
