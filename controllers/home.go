package controllers

import (
	"github.com/curt-labs/GoAPI/helpers/plate"
	. "github.com/curt-labs/GoAPI/models"

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
