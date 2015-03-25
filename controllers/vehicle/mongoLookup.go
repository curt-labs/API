package vehicle

import (
	"github.com/curt-labs/GoAPI/helpers/apicontext"
	"github.com/curt-labs/GoAPI/helpers/encoding"
	"github.com/curt-labs/GoAPI/helpers/error"
	"github.com/curt-labs/GoAPI/models/products"

	"net/http"
	"strconv"
)

func Lookup(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, dtx *apicontext.DataContext) string {
	var v products.NoSqlVehicle
	var collection string //e.g. interior/exterior

	//Get collection
	collection = r.FormValue("collection")
	delete(r.Form, "collection")

	// Get vehicle year
	y_str := r.FormValue("year")
	v.Year, _ = strconv.Atoi(y_str)
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

	l, err := products.FindVehicles(v, collection, dtx)
	if err != nil {
		apierror.GenerateError("Trouble finding vehicles.", err, w, r)
		return ""
	}

	return encoding.Must(enc.Encode(l))
}
