package aces

type BaseVehicle struct {
	Year  int    `json:"year" xml:"year,attr"`
	Make  string `json:"make" xml:"make"`
	Model string `json:"model" xml:"model"`
}
