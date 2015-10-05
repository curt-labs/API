package category_ctlr

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/curt-labs/GoAPI/controllers/vehicle"
	"github.com/curt-labs/GoAPI/helpers/apicontext"
	"github.com/curt-labs/GoAPI/helpers/apifilter"
	"github.com/curt-labs/GoAPI/helpers/encoding"
	"github.com/curt-labs/GoAPI/helpers/error"
	"github.com/curt-labs/GoAPI/models/products"
	"github.com/go-martini/martini"
)

var (
	// NoFilterCategories = map[int]int{1: 1, 3: 3, 4: 4, 5: 5, 8: 8, 9: 9, 254: 254, 2: 2, 11: 11, 12: 12, 13: 13, 14: 14, 273: 273}
	NoFilterCategories = map[int]int{}
	NoFilterKeys       = map[string]string{"key": "key", "page": "page", "count": "count"}
)

type FilterSpecifications struct {
	Key    string   `json:"key" xml:"key"`
	Values []string `json:"values" xml:"values"`
}

func GetCategory(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params, dtx *apicontext.DataContext) string {
	var page int
	var count int
	var cat products.Category
	var l products.Lookup

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		apierror.GenerateError("Trouble reading request body while getting category", err, rw, r)
	}

	qs := r.URL.Query()
	key := qs.Get("key")
	id, err := strconv.Atoi(params["id"])
	page, _ = strconv.Atoi(qs.Get("page"))
	count, _ = strconv.Atoi(qs.Get("count"))

	// Load Vehicle from Request
	l.Vehicle = vehicle.LoadVehicle(r)
	defer r.Body.Close()
	specs := make(map[string][]string, 0)
	if strings.Contains(r.Header.Get("Content-Type"), "json") && len(data) > 0 {
		var fs []FilterSpecifications
		if err = json.Unmarshal(data, &fs); err == nil {
			for _, f := range fs {
				specs[f.Key] = f.Values
			}
		}
	} else {
		r.ParseForm()
		if _, ignore := NoFilterCategories[cat.ID]; !ignore {
			for k, v := range r.Form {
				if _, excluded := NoFilterKeys[strings.ToLower(k)]; !excluded {
					if _, ok := specs[k]; !ok {
						specs[k] = make([]string, 0)
					}
					specs[k] = append(specs[k], v...)
				}
			}
		}
	}
	// Get Category
	if err != nil { // get by title
		cat, err = products.GetCategoryByTitle(params["id"], dtx)
		if err != nil {
			apierror.GenerateError("Trouble getting category by title", err, rw, r)
		}
	} else { // get by id
		cat.ID = id
		if err = cat.GetCategory(key, page, count, false, &l.Vehicle, &specs, dtx); err != nil {
			apierror.GenerateError("Trouble getting category", err, rw, r)
		}
	}

	if _, ignore := NoFilterCategories[cat.ID]; !ignore {
		if filters, err := apifilter.CategoryFilter(cat, &specs); err == nil {
			cat.Filter = filters
		}
	}

	return encoding.Must(enc.Encode(cat))
}

func Parents(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder, dtx *apicontext.DataContext) string {
	// var err error
	var c []products.Category

	c, err := products.TopTierCategories(dtx)
	if err != nil {
		apierror.GenerateError("Trouble getting parent categories", err, rw, r)
	}
	return encoding.Must(enc.Encode(c))
}

func Tree(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder, dtx *apicontext.DataContext) string {
	var c []products.Category

	c, err := products.CategoryTree(dtx)
	if err != nil {
		apierror.GenerateError("Trouble getting category tree", err, rw, r)
	}
	return encoding.Must(enc.Encode(c))
}

func SubCategories(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params, dtx *apicontext.DataContext) string {
	id, err := strconv.Atoi(params["id"])

	var cat products.Category
	if err != nil {
		cat, err = products.GetCategoryByTitle(params["id"], dtx)
		if err != nil {
			apierror.GenerateError("Trouble getting category for sub categories", err, rw, r)
		}
	} else {
		cat.ID = id
	}

	subs, err := cat.GetSubCategories(dtx)
	if err != nil {
		apierror.GenerateError("Trouble getting sub categories", err, rw, r)
	}

	return encoding.Must(enc.Encode(subs))
}

func GetParts(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params, dtx *apicontext.DataContext) string {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		apierror.GenerateError("Trouble reading request body for get category parts", err, rw, r)
	}
	key := params["key"]
	catID, err := strconv.Atoi(params["id"])
	qs := r.URL.Query()

	var cat products.Category
	if err != nil {
		title := params["id"]
		if title == "" {
			apierror.GenerateError("Trouble getting category title for get category parts", err, rw, r)
		}
		cat, err = products.GetCategoryByTitle(title, dtx)
		if err != nil {
			apierror.GenerateError("Trouble getting category by title for get category parts", err, rw, r)
		}
	} else {
		cat.ID = catID
	}
	count, _ := strconv.Atoi(qs.Get("count"))
	page, _ := strconv.Atoi(qs.Get("page"))

	defer r.Body.Close()
	specs := make(map[string][]string, 0)
	if strings.Contains(r.Header.Get("Content-Type"), "json") && len(data) > 0 {
		var fs []FilterSpecifications
		if err = json.Unmarshal(data, &fs); err == nil {

			for _, f := range fs {
				specs[f.Key] = f.Values
			}
		}
	} else {
		r.ParseForm()
		if _, ignore := NoFilterCategories[cat.ID]; !ignore {
			for k, v := range r.Form {
				if _, excluded := NoFilterKeys[strings.ToLower(k)]; !excluded {
					if _, ok := specs[k]; !ok {
						specs[k] = make([]string, 0)
					}
					specs[k] = append(specs[k], v...)
				}
			}
		}
	}

	if err = cat.GetParts(key, page, count, nil, &specs, dtx); err != nil {
		apierror.GenerateError("Trouble getting category parts", err, rw, r)
	}
	return encoding.Must(enc.Encode(cat.ProductListing))
}
