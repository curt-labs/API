package category_ctlr

import (
	"encoding/json"
	"github.com/curt-labs/GoAPI/controllers/vehicle"
	"github.com/curt-labs/GoAPI/helpers/apifilter"
	"github.com/curt-labs/GoAPI/helpers/encoding"
	"github.com/curt-labs/GoAPI/models/products"
	"github.com/go-martini/martini"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

var (
	NoFilterCategories = map[int]int{1: 1, 3: 3, 4: 4, 5: 5, 8: 8, 9: 9, 254: 254, 2: 2, 11: 11, 12: 12, 13: 13, 14: 14, 273: 273}
	NoFilterKeys       = map[string]string{"key": "key", "page": "page", "count": "count"}
)

type FilterSpecifications struct {
	Key    string   `json:"key" xml:"key"`
	Values []string `json:"values" xml:"values"`
}

func GetCategory(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	var page int
	var count int
	var cat products.Category
	var l products.Lookup
	data, _ := ioutil.ReadAll(r.Body)

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
		cat, err = products.GetCategoryByTitle(params["id"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return ""
		}
	} else { // get by id
		cat.ID = id
		if err = cat.GetCategory(key, page, count, false, &l.Vehicle, &specs); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return ""
		}
	}

	if _, ignore := NoFilterCategories[cat.ID]; !ignore {
		if filters, err := apifilter.CategoryFilter(cat, &specs); err == nil {
			cat.Filter = filters
		}
	}

	return encoding.Must(enc.Encode(cat))
}

func Parents(w http.ResponseWriter, r *http.Request, enc encoding.Encoder) string {
	qs := r.URL.Query()
	key := qs.Get("key")

	cats, err := products.TopTierCategories(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	return encoding.Must(enc.Encode(cats))
}

func SubCategories(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	id, err := strconv.Atoi(params["id"])

	var cat products.Category
	if err != nil {
		cat, err = products.GetCategoryByTitle(params["id"])
		if err != nil {

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return ""
		}
	} else {
		cat.ID = id
	}

	subs, err := cat.GetSubCategories()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	return encoding.Must(enc.Encode(subs))
}

func GetParts(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	data, _ := ioutil.ReadAll(r.Body)
	key := params["key"]
	catID, err := strconv.Atoi(params["id"])
	qs := r.URL.Query()

	var cat products.Category
	if err != nil {
		title := params["id"]
		if title == "" {
			http.Error(w, "Invalid Category", http.StatusInternalServerError)
			return ""
		}
		cat, err = products.GetCategoryByTitle(title)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return ""
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

	log.Println(specs)

	if err := cat.GetParts(key, page, count, nil, &specs); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	return encoding.Must(enc.Encode(cat.ProductListing))
}
