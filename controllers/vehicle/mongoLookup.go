package vehicle

import (
	"github.com/curt-labs/API/helpers/apicontext"
	"github.com/curt-labs/API/helpers/database"
	"github.com/curt-labs/API/helpers/encoding"
	"github.com/curt-labs/API/helpers/error"
	"github.com/curt-labs/API/models/products"

	"net/http"
	"sort"
	"strconv"
	"strings"
)

func Collections(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, dtx *apicontext.DataContext) string {
	if err := database.Init(); err != nil {
		apierror.GenerateError(err.Error(), err, w, r)
		return ""
	}

	cols, err := products.GetAriesVehicleCollections(database.ProductMongoSession)
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

func ByCategoryLuverne(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, dtx *apicontext.DataContext) string {
	catID, err := strconv.Atoi(r.FormValue("catID"))
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

	res, err := products.FindAppsLuverne(catID, offset, limit)
	if err != nil {
		apierror.GenerateError("Trouble finding vehicles.", err, w, r)
		return ""
	}

	return encoding.Must(enc.Encode(res))
}

//Hack version that slowly traverses all the collection and aggregates results
func AllCollectionsLookup(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, dtx *apicontext.DataContext) string {
	if err := database.Init(); err != nil {
		apierror.GenerateError(err.Error(), err, w, r)
		return ""
	}

	sess := database.ProductMongoSession.Copy()
	defer sess.Close()

	//Get all collections
	cols, err := products.GetAriesVehicleCollections(sess)
	if err != nil {
		apierror.GenerateError(err.Error(), err, w, r)
		return ""
	}

	var v products.NoSqlVehicle

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

	if r.FormValue("collection") != "" {
		tmpCols := cols
		cols = []string{}
		for _, tc := range tmpCols {
			if strings.ToLower(tc) == strings.ToLower(r.FormValue("collection")) {
				cols = []string{tc}
				break
			}
		}
	}

	var collectionVehicleArray []products.NoSqlLookup

	for _, col := range cols {
		noSqlLookup, err := products.FindVehiclesWithParts(v, col, dtx, sess)
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
	partmap := make(map[int]products.Part)

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
		for _, part := range lookup.Parts {
			partmap[part.ID] = part
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
	for _, part := range partmap {
		l.Parts = append(l.Parts, part)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(l.Years)))
	sort.Strings(l.Makes)
	sort.Strings(l.Models)
	sort.Strings(l.Styles)

	return l
}

//return parts for a vehicle(incl style) within a specific category
func AllCollectionsLookupCategory(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, dtx *apicontext.DataContext) string {
	var v products.NoSqlVehicle
	noSqlLookup := make(map[string]products.NoSqlLookup)
	var err error

	if err := database.Init(); err != nil {
		apierror.GenerateError(err.Error(), err, w, r)
		return ""
	}

	sess := database.ProductMongoSession.Copy()
	defer sess.Close()

	// Get vehicle year
	v.Year = r.FormValue("year")
	delete(r.Form, "year")

	// Get vehicle make
	v.Make = r.FormValue("make")
	delete(r.Form, "make")

	// Get vehicle model
	v.Model = r.FormValue("model")
	delete(r.Form, "model")

	// // Get vehicle submodel
	v.Style = r.FormValue("style")
	delete(r.Form, "style")

	collection := r.FormValue("collection")
	if collection == "" {
		noSqlLookup, err = products.FindVehiclesFromAllCategories(v, dtx, sess)
	} else {
		noSqlLookup, err = products.FindPartsFromOneCategory(v, collection, dtx, sess)
	}
	if err != nil {
		apierror.GenerateError("Trouble finding vehicles.", err, w, r)
		return ""
	}

	return encoding.Must(enc.Encode(noSqlLookup))
}
