package vehicle_ctlr

import (
	"../../models/vehicle"
	"../../plate"
	"net/http"
	"strconv"
	"strings"
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

	v := vehicle.Vehicle{
		Year: year,
		Make: params.Get(":make"),
	}

	config := vehicle.ConfigResponse{
		ConfigOption: v.GetModels(),
		Matched:      new(vehicle.ProductMatch),
	}

	plate.ServeFormatted(w, r, config)
	return
}

func Submodel(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	year, _ := strconv.ParseFloat(params.Get(":year"), 64)

	v := vehicle.Vehicle{
		Year:  year,
		Make:  params.Get(":make"),
		Model: params.Get(":model"),
	}

	var subs vehicle.ConfigOption
	var matched *vehicle.ProductMatch

	subsChan := make(chan int)
	matchedChan := make(chan int)
	go func() {
		subs = v.GetSubmodels()
		subsChan <- 1
	}()
	go func() {
		matched = v.GetProductMatch()
		matchedChan <- 1
	}()
	<-subsChan
	<-matchedChan

	config := vehicle.ConfigResponse{
		ConfigOption: subs,
		Matched:      matched,
	}

	plate.ServeFormatted(w, r, config)
	return
}

func Config(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	year, _ := strconv.ParseFloat(params.Get(":year"), 64)

	config_vals := strings.Split(params.Get(":config"), "/")

	v := vehicle.Vehicle{
		Year:          year,
		Make:          params.Get(":make"),
		Model:         params.Get(":model"),
		Submodel:      params.Get(":submodel"),
		Configuration: config_vals,
	}

	var config_opts vehicle.ConfigOption
	var matched *vehicle.ProductMatch

	configChan := make(chan int)
	matchedChan := make(chan int)
	go func() {
		config_opts = v.GetConfiguration()
		configChan <- 1
	}()
	go func() {
		matched = v.GetProductMatch()
		matchedChan <- 1
	}()
	<-configChan
	<-matchedChan

	config := vehicle.ConfigResponse{
		ConfigOption: config_opts,
		Matched:      matched,
	}

	plate.ServeFormatted(w, r, config)
	return
}
