package vehicle_ctlr

import (
	"../../helpers/plate"
	. "../../models"
	"log"
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
	log.Println("fetching makes")

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
		matched = lookup.Vehicle.GetProductMatch(k)
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

	config_vals := strings.Split(params.Get(":config"), "/")

	v := Vehicle{
		Year:          year,
		Make:          params.Get(":make"),
		Model:         params.Get(":model"),
		Submodel:      params.Get(":submodel"),
		Configuration: config_vals,
	}

	var config_opts ConfigOption
	var matched *ProductMatch

	configChan := make(chan int)
	matchedChan := make(chan int)
	go func() {
		config_opts = v.GetConfiguration()
		configChan <- 1
	}()
	go func(k string) {
		matched = v.GetProductMatch(k)
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
