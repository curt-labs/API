package controllers

import (
	"../../helpers/plate"
	. "../models"

	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request) {
	var v Vehicle

	config := ConfigResponse{
		ConfigOption: v.GetYears(),
		Matched:      new(ProductMatch),
	}

	plate.ServeFormatted(w, r, config)
	return
}
