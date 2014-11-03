package main

type PassivScanRequest struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	UseHTTPS bool   `json:"useHttps"`
	Request  string `json:"request"`
	Response string `json:"response"`
}
