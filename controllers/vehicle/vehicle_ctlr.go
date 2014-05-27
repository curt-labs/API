package vehicle_ctlr

import (
	"github.com/curt-labs/GoAPI/helpers/plate"
	. "github.com/curt-labs/GoAPI/models"
	"github.com/go-martini/martini"
	"net/http"
	"strconv"
	"strings"
)

func Year(w http.ResponseWriter, r *http.Request) {
	var l Lookup

	config := ConfigResponse{
		ConfigOption: l.GetYears(),
		Matched:      new(ProductMatch),
	}

	plate.ServeFormatted(w, r, config)
	return
}

func Make(w http.ResponseWriter, r *http.Request, params martini.Params) {
	year, _ := strconv.ParseFloat(params["year"], 64)
	lookup := Lookup{
		Vehicle: Vehicle{
			Year: year,
		},
	}

	config := ConfigResponse{
		ConfigOption: lookup.GetMakes(),
		Matched:      new(ProductMatch),
	}

	plate.ServeFormatted(w, r, config)
	return
}

func Model(w http.ResponseWriter, r *http.Request, params martini.Params) {
	year, _ := strconv.ParseFloat(params["year"], 64)

	lookup := Lookup{
		Vehicle: Vehicle{
			Year: year,
			Make: params["make"],
		},
	}

	config := ConfigResponse{
		ConfigOption: lookup.GetModels(),
		Matched:      new(ProductMatch),
	}

	plate.ServeFormatted(w, r, config)
	return
}

func Submodel(w http.ResponseWriter, r *http.Request, params martini.Params) {
	qs := r.URL.Query()
	year, _ := strconv.ParseFloat(params["year"], 64)
	key := qs.Get("key")

	lookup := Lookup{
		Vehicle: Vehicle{
			Year:  year,
			Make:  params["make"],
			Model: params["model"],
		},
	}

	var subs ConfigOption
	var matched *ProductMatch

	subsChan := make(chan int)
	matchedChan := make(chan int)
	go func() {
		subs = lookup.GetSubmodels()
		subsChan <- 1
	}()
	go func(k string) {
		matched = lookup.GetProductMatch(k)
		matchedChan <- 1
	}(key)
	<-subsChan
	<-matchedChan

	config := ConfigResponse{
		ConfigOption: subs,
		Matched:      matched,
	}

	plate.ServeFormatted(w, r, config)
	return
}

func Config(w http.ResponseWriter, r *http.Request, params martini.Params) {
	qs := r.URL.Query()
	year, _ := strconv.ParseFloat(params["year"], 64)
	key := qs.Get("key")

	config_vals := strings.Split(strings.TrimSpace(params["config"]), "/")

	if len(config_vals) == 1 && config_vals[0] == "" {
		config_vals = nil
	}

	lookup := Lookup{
		Vehicle: Vehicle{
			Year:          year,
			Make:          params["make"],
			Model:         params["model"],
			Submodel:      params["submodel"],
			Configuration: config_vals,
		},
	}

	var config_opts ConfigOption
	var matched *ProductMatch

	configChan := make(chan int)
	matchedChan := make(chan int)
	go func() {
		config_opts = lookup.GetConfiguration()
		configChan <- 1
	}()
	go func(k string) {
		matched = lookup.GetProductMatch(k)
		matchedChan <- 1
	}(key)
	<-configChan
	<-matchedChan

	config := ConfigResponse{
		ConfigOption: config_opts,
		Matched:      matched,
	}

	plate.ServeFormatted(w, r, config)
	return
}

func Connector(w http.ResponseWriter, r *http.Request, params martini.Params) {
	year, err := strconv.ParseFloat(params["year"], 64)
	if err != nil {
		http.Error(w, "Failed to process vehicle", http.StatusInternalServerError)
		return
	}
	qs := r.URL.Query()
	key := qs.Get("key")

	config_vals := strings.Split(strings.TrimSpace(params["config"]), "/")

	if len(config_vals) == 1 && config_vals[0] == "" {
		config_vals = nil
	}

	lookup := Lookup{
		Vehicle: Vehicle{
			Year:          year,
			Make:          params["make"],
			Model:         params["model"],
			Submodel:      params["submodel"],
			Configuration: config_vals,
		},
	}

	err = lookup.GetConnector(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	plate.ServeFormatted(w, r, lookup.Parts)
}
