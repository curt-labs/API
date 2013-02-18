package vehicle_ctlr

import (
	"../../models/vehicle"
	"../../plate"
	"net/http"
	"strconv"
)

func Year(w http.ResponseWriter, r *http.Request) {
	var v vehicle.Vehicle

	config := vehicle.ConfigResponse{
		ConfigOption: v.GetYears(),
		Matched:      new(vehicle.ProductMatch),
	}

	plate.ServeFormatted(w, r, config)
	return
}

func Make(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	year, _ := strconv.ParseFloat(params.Get(":year"), 64)
	v := vehicle.Vehicle{
		Year: year,
	}

	config := vehicle.ConfigResponse{
		ConfigOption: v.GetMakes(),
		Matched:      new(vehicle.ProductMatch),
	}

	plate.ServeFormatted(w, r, config)
	return
}

func Model(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	year, _ := strconv.ParseFloat(params.Get(":year"), 64)
	make := params.Get(":make")
	v := vehicle.Vehicle{
		Year: year,
		Make: make,
	}

	config := vehicle.ConfigResponse{
		ConfigOption: v.GetModels(),
		Matched:      new(vehicle.ProductMatch),
	}

	plate.ServeFormatted(w, r, config)
	return
}
