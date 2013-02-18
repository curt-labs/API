package controllers

import (
	"../models/vehicle"
	"../plate"

	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request) {
	var v vehicle.Vehicle
	plate.ServeFormatted(w, r, v.GetYears())
	return
}
