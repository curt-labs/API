package vehicle_ctlr

import (
	"github.com/curt-labs/GoAPI/helpers/encoding"
	. "github.com/curt-labs/GoAPI/models"
	"github.com/go-martini/martini"
	"net/http"
	"strconv"
	"strings"
)

func Year(w http.ResponseWriter, r *http.Request, enc encoding.Encoder) string {
	var l Lookup

	config := ConfigResponse{
		ConfigOption: l.GetYears(),
		Matched:      new(ProductMatch),
	}

	return encoding.Must(enc.Encode(config))
}

func Make(w http.ResponseWriter, r *http.Request, params martini.Params, enc encoding.Encoder) string {
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

	return encoding.Must(enc.Encode(config))
}

func Model(w http.ResponseWriter, r *http.Request, params martini.Params, enc encoding.Encoder) string {
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

	return encoding.Must(enc.Encode(config))
}

func Submodel(w http.ResponseWriter, r *http.Request, params martini.Params, enc encoding.Encoder) string {
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

	return encoding.Must(enc.Encode(config))
}

func Config(w http.ResponseWriter, r *http.Request, params martini.Params, enc encoding.Encoder) string {
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

	return encoding.Must(enc.Encode(config))
}

func Connector(w http.ResponseWriter, r *http.Request, params martini.Params, enc encoding.Encoder) string {
	year, err := strconv.ParseFloat(params["year"], 64)
	if err != nil {
		http.Error(w, "Failed to process vehicle", http.StatusInternalServerError)
		return ""
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
		return ""
	}

	return encoding.Must(enc.Encode(lookup.Parts))
}
