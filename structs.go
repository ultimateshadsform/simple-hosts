package main

type Host struct {
	Hostname string `json:"host"`
	IP       string `json:"ip"`
	Comment  string `json:"comment"`
}
