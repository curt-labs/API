package vehicle

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/curt-labs/GoAPI/helpers/apicontext"
	"github.com/curt-labs/GoAPI/helpers/apifilter"
	"github.com/curt-labs/GoAPI/helpers/encoding"
	"github.com/curt-labs/GoAPI/helpers/error"
	"github.com/curt-labs/GoAPI/models/products"
)

var (
	ignoredFormParams = []string{"key"}
)

// Finds further configuration options and parts that match
// the given configuration. Doesn't start looking for parts
// until the model is provided.
func Query(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, dtx *apicontext.DataContext) string {
	var l products.Lookup
	l.Vehicle = LoadVehicle(r)
	l.Brands = dtx.BrandArray

	qs := r.URL.Query()
	if qs.Get("key") != "" {
		l.CustomerKey = qs.Get("key")
	} else if r.FormValue("key") != "" {
		l.CustomerKey = r.FormValue("key")
		delete(r.Form, "key")
	} else {
		l.CustomerKey = r.Header.Get("key")
	}

	if l.Vehicle.Base.Year == 0 { // Get Years
		if err := l.GetYears(dtx); err != nil {
			apierror.GenerateError("Trouble getting years for vehicle lookup", err, w, r)
		}
	} else if l.Vehicle.Base.Make == "" { // Get Makes
		if err := l.GetMakes(dtx); err != nil {
			apierror.GenerateError("Trouble getting makes for vehicle lookup", err, w, r)
		}
	} else if l.Vehicle.Base.Model == "" { // Get Models
		if err := l.GetModels(); err != nil {
			apierror.GenerateError("Trouble getting models for vehicle lookup", err, w, r)
		}
	} else {

		// Kick off part getter
		partChan := make(chan []products.Part)
		go l.LoadParts(partChan, dtx)

		if l.Vehicle.Submodel == "" { // Get Submodels
			if err := l.GetSubmodels(); err != nil {
				apierror.GenerateError("Trouble getting submodels for vehicle lookup", err, w, r)
			}
		} else { // Get configurations
			if err := l.GetConfigurations(); err != nil {
				apierror.GenerateError("Trouble getting configurations for vehicle lookup", err, w, r)
			}
		}

		select {
		case parts := <-partChan:
			if len(parts) > 0 {
				l.Parts = parts
				l.Filter, _ = apifilter.PartFilter(l.Parts, nil)
			}
		case <-time.After(5 * time.Second):

		}
	}

	return encoding.Must(enc.Encode(l))
}

// Parses the vehicle data out of the request
// body. It will first check for Content-Type as
// JSON and parse accordingly.
func LoadVehicle(r *http.Request) (v products.Vehicle) {
	defer r.Body.Close()

	if strings.Contains(strings.ToLower(r.Header.Get("Content-Type")), "json") {
		if data, err := ioutil.ReadAll(r.Body); err == nil || len(data) > 0 {
			err = json.Unmarshal(data, &v)
			if err == nil && v.Base.Year > 0 {
				return
			}
		}
	}

	// Get vehicle year
	y_str := r.FormValue("year")
	if y_str == "" {
		return
	}
	v.Base.Year, _ = strconv.Atoi(y_str)
	if v.Base.Year == 0 {
		return
	}
	delete(r.Form, "year")

	// Get vehicle make
	v.Base.Make = r.FormValue("make")
	if v.Base.Make == "" {
		return
	}
	delete(r.Form, "make")

	// Get vehicle model
	v.Base.Model = r.FormValue("model")
	if v.Base.Model == "" {
		return
	}
	delete(r.Form, "model")

	// Get vehicle submodel
	v.Submodel = r.FormValue("submodel")
	if v.Submodel == "" {
		return
	}
	delete(r.Form, "submodel")

	// Get vehicle configuration options
	for key, opt := range r.Form {
		ignore := false
		for _, param := range ignoredFormParams {
			if param == strings.ToLower(key) {
				ignore = true
				break
			}
		}
		if !ignore && len(opt) > 0 {
			conf := products.Configuration{
				Key:   key,
				Value: opt[0],
			}
			v.Configurations = append(v.Configurations, conf)
		}
	}

	return
}
