package vehicle_ctlr

import (
	"github.com/curt-labs/GoAPI/helpers/plate"
	. "github.com/curt-labs/GoAPI/models"
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

func Make(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	year, _ := strconv.ParseFloat(params.Get(":year"), 64)
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

func Model(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	year, _ := strconv.ParseFloat(params.Get(":year"), 64)

	lookup := Lookup{
		Vehicle: Vehicle{
			Year: year,
			Make: params.Get(":make"),
		},
	}

	config := ConfigResponse{
		ConfigOption: lookup.GetModels(),
		Matched:      new(ProductMatch),
	}

	plate.ServeFormatted(w, r, config)
	return
}

func Submodel(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	year, _ := strconv.ParseFloat(params.Get(":year"), 64)
	key := params.Get("key")

	lookup := Lookup{
		Vehicle: Vehicle{
			Year:  year,
			Make:  params.Get(":make"),
			Model: params.Get(":model"),
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

func Config(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	year, _ := strconv.ParseFloat(params.Get(":year"), 64)
	key := params.Get("key")

	config_vals := strings.Split(strings.TrimSpace(params.Get(":config")), "/")

	if len(config_vals) == 1 && config_vals[0] == "" {
		config_vals = nil
	}

	lookup := Lookup{
		Vehicle: Vehicle{
			Year:          year,
			Make:          params.Get(":make"),
			Model:         params.Get(":model"),
			Submodel:      params.Get(":submodel"),
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

func Connector(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	year, err := strconv.ParseFloat(params.Get(":year"), 64)
	if err != nil {
		http.Error(w, "Failed to process vehicle", http.StatusInternalServerError)
		return
	}
	key := params.Get("key")

	config_vals := strings.Split(strings.TrimSpace(params.Get(":config")), "/")

	if len(config_vals) == 1 && config_vals[0] == "" {
		config_vals = nil
	}

	lookup := Lookup{
		Vehicle: Vehicle{
			Year:          year,
			Make:          params.Get(":make"),
			Model:         params.Get(":model"),
			Submodel:      params.Get(":submodel"),
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
