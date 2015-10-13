package vehicle

import (
	"github.com/curt-labs/GoAPI/helpers/apicontext"
	"github.com/curt-labs/GoAPI/helpers/encoding"
	"github.com/curt-labs/GoAPI/helpers/error"
	"github.com/curt-labs/GoAPI/models/products"
	"net/http"
	"sort"
	"strconv"
)

func Collections(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, dtx *apicontext.DataContext) string {
	cols, err := products.GetAriesVehicleCollections()
	if err != nil {
		apierror.GenerateError(err.Error(), err, w, r)
		return ""
	}

	return encoding.Must(enc.Encode(cols))
}

func Lookup(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, dtx *apicontext.DataContext) string {
	var v products.NoSqlVehicle
	var collection string //e.g. interior/exterior

	//Get collection
	collection = r.FormValue("collection")
	delete(r.Form, "collection")

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

	l, err := products.FindVehicles(v, collection, dtx)
	if err != nil {
		apierror.GenerateError("Trouble finding vehicles.", err, w, r)
		return ""
	}

	return encoding.Must(enc.Encode(l))
}

func ByCategory(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, dtx *apicontext.DataContext) string {
	collection := r.FormValue("collection")
	page, _ := strconv.Atoi(r.FormValue("page"))
	limit, _ := strconv.Atoi(r.FormValue("limit"))
	var offset int
	if page == 0 {
		offset = 0
	} else if page == 1 {
		offset = 101
	} else {
		offset = page*limit + 1
	}

	res, err := products.FindApplications(collection, offset, limit)
	if err != nil {
		apierror.GenerateError("Trouble finding vehicles.", err, w, r)
		return ""
	}

	return encoding.Must(enc.Encode(res))
}

//Hack version that slowly traverses all the collection and aggregates results
func AllCollectionsLookup(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, dtx *apicontext.DataContext) string {
	var v products.NoSqlVehicle

	//Get all collections
	cols, err := products.GetAriesVehicleCollections()
	if err != nil {
		apierror.GenerateError(err.Error(), err, w, r)
		return ""
	}

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

	var collectionVehicleArray []products.NoSqlLookup

	for _, col := range cols {
		noSqlLookup, err := products.FindVehicles(v, col, dtx)
		if err != nil {
			apierror.GenerateError("Trouble finding vehicles.", err, w, r)
			return ""
		}
		collectionVehicleArray = append(collectionVehicleArray, noSqlLookup)
	}
	l := makeLookupFrommanyLookups(collectionVehicleArray)

	return encoding.Must(enc.Encode(l))
}

func makeLookupFrommanyLookups(lookupArrays []products.NoSqlLookup) (l products.NoSqlLookup) {
	yearmap := make(map[string]string)
	makemap := make(map[string]string)
	modelmap := make(map[string]string)
	stylemap := make(map[string]string)
	for _, lookup := range lookupArrays {
		for _, year := range lookup.Years {
			yearmap[year] = year
		}
		for _, mk := range lookup.Makes {
			makemap[mk] = mk
		}
		for _, model := range lookup.Models {
			modelmap[model] = model
		}
		for _, style := range lookup.Styles {
			stylemap[style] = style
		}
	}
	for year, _ := range yearmap {
		l.Years = append(l.Years, year)
	}
	for mk, _ := range makemap {
		l.Makes = append(l.Makes, mk)
	}
	for model, _ := range modelmap {
		l.Models = append(l.Models, model)
	}
	for style, _ := range stylemap {
		l.Styles = append(l.Styles, style)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(l.Years)))
	sort.Strings(l.Makes)
	sort.Strings(l.Models)
	sort.Strings(l.Styles)

	return l
}
