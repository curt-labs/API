package controllers

import (
	"../models/vehicle"
	"../plate"

	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request) {
	var v vehicle.Vehicle

	config := vehicle.ConfigResponse{
		ConfigOption: v.GetYears(),
		Matched:      new(vehicle.ProductMatch),
	}

	plate.ServeFormatted(w, r, config)
	return
}
