package vehicle

import (
	"errors"

	"github.com/curt-labs/API/helpers/apicontext"
	"github.com/curt-labs/API/helpers/encoding"
	"github.com/curt-labs/API/helpers/error"
	"github.com/curt-labs/API/models/products"

	"net/http"
	"net/url"
)

func CurtLookup(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, dtx *apicontext.DataContext) string {
	var v products.CurtVehicle
	var err error

	// Get vehicle year
	v.Year = r.FormValue("year")
	delete(r.Form, "year")

	// Get vehicle make
	v.Make = r.FormValue("make")
	delete(r.Form, "make")

	// Get vehicle model
	v.Model = r.FormValue("model")
	delete(r.Form, "model")

	// Get vehicle submodel
	v.Style = r.FormValue("style")
	delete(r.Form, "style")

	cl := products.CurtLookup{
		CurtVehicle: v,
	}

	// heavyduty
	var heavyduty bool
	hdstr := r.URL.Query().Get("heavyduty")
	if hdstr == "true" {
		heavyduty = true
	}
	// determine if you are going to get customer prices for each part
	var getCustomerPrices bool
	custPricingStr := r.URL.Query().Get("customerPrices")
	if custPricingStr == "true" || custPricingStr == "True" {
		getCustomerPrices = true
	}

	cl, err = CurtLookupWorker(cl, heavyduty, dtx, getCustomerPrices)

	if err != nil {
		apierror.GenerateError("Trouble finding vehicles.", err, w, r)
		return ""
	}

	return encoding.Must(enc.Encode(cl))
}

//CurtLookupGet is the exact same as CurtLookup, except in a GET request as it
//should be, as defined in RFC 7231. This will also help mitigate Google Cloud
//related 502s
func CurtLookupGet(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, dtx *apicontext.DataContext) string {
	var v products.CurtVehicle
	var err error

	v.Year = r.URL.Query().Get("year")
	v.Make = r.URL.Query().Get("make")
	v.Model = r.URL.Query().Get("model")
	v.Style = r.URL.Query().Get("style")

	v.Year, err = url.QueryUnescape(v.Year)
	v.Make, err = url.QueryUnescape(v.Make)
	v.Model, err = url.QueryUnescape(v.Model)
	v.Style, err = url.QueryUnescape(v.Style)

	if err != nil {
		return err.Error()
	}

	cl := products.CurtLookup{
		CurtVehicle: v,
	}

	// heavyduty
	var heavyduty bool
	hdstr := r.URL.Query().Get("heavyduty")
	if hdstr == "true" {
		heavyduty = true
	}

	// determine if you are going to get customer prices for each part
	var getCustomerPrices bool
	custPricingStr := r.URL.Query().Get("customerprices")
	if custPricingStr == "true" || custPricingStr == "True" {
		getCustomerPrices = true
	}

	cl, err = CurtLookupWorker(cl, heavyduty, dtx, getCustomerPrices)

	if err != nil {
		apierror.GenerateError("Trouble finding vehicles.", err, w, r)
		return ""
	}

	return encoding.Must(enc.Encode(cl))
}

//CurtLookupWorker is a function that is used by both the POST and GET versions
//of CurtLookup. Just here so that there is as little code duplication as possible
func CurtLookupWorker(cl products.CurtLookup, heavyduty bool, dtx *apicontext.DataContext, getCustomerPricing bool) (products.CurtLookup, error) {
	var err error
	if cl.CurtVehicle.Year == "" {
		err = cl.GetYears(heavyduty)
	} else if cl.CurtVehicle.Make == "" {
		err = cl.GetMakes(heavyduty)
	} else if cl.CurtVehicle.Model == "" {
		err = cl.GetModels(heavyduty)
	} else {
		err = cl.GetStyles(heavyduty)
		if err != nil {
			return cl, errors.New("Trouble finding styles.")
		}
		err = cl.GetParts(dtx, heavyduty)
		if getCustomerPricing {
			cl.Parts, err = products.BindCustomerToSeveralParts(cl.Parts, dtx)
		}
	}

	return cl, err
}
