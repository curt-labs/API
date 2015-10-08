package vehicle

import (
	"github.com/curt-labs/GoAPI/helpers/apicontext"
	"github.com/curt-labs/GoAPI/helpers/encoding"
	"github.com/curt-labs/GoAPI/helpers/error"
	"github.com/curt-labs/GoAPI/models/products"

	"net/http"
)

func CurtLookup(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, dtx *apicontext.DataContext) string {
	var v products.CurtVehicle

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

	var err error
	if v.Year == "" {
		err = cl.GetYears()
	} else if v.Make == "" {
		err = cl.GetMakes()
	} else if v.Model == "" {
		err = cl.GetModels()
	} else if v.Style == "" {
		err = cl.GetStyles()
	} else {
		err = cl.GetParts(dtx)
	}

	if err != nil {
		apierror.GenerateError("Trouble finding vehicles.", err, w, r)
		return ""
	}

	return encoding.Must(enc.Encode(cl))
}
