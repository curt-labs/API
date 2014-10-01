package category_ctlr

import (
	"github.com/curt-labs/GoAPI/helpers/apifilter"
	"github.com/curt-labs/GoAPI/helpers/encoding"
	"github.com/curt-labs/GoAPI/models/products"
	"github.com/go-martini/martini"
	"net/http"
	"strconv"
)

var (
	NoFilterCategories = map[int]int{1: 1, 3: 3, 4: 4, 5: 5, 8: 8, 9: 9, 254: 254, 2: 2, 11: 11, 12: 12, 13: 13, 14: 14, 273: 273}
)

func GetCategory(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	qs := r.URL.Query()
	key := qs.Get("key")
	id, err := strconv.Atoi(params["id"])
	var page int
	var count int
	page, _ = strconv.Atoi(qs.Get("page"))
	count, _ = strconv.Atoi(qs.Get("count"))

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

	if err = cat.GetCategory(key, page, count, false); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	if _, ignore := NoFilterCategories[cat.ID]; !ignore {
		if filters, err := apifilter.CategoryFilter(cat, nil); err == nil {
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
	key := params["key"]
	catID, err := strconv.Atoi(params["id"])

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

	count, err := strconv.Atoi(params["count"])
	if err != nil {
		count = 5
	}

	page, err := strconv.Atoi(params["page"])
	if err != nil {
		page = 0
	} else {
		page = page - 1
	}

	if err := cat.GetParts(key, page, count, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	return encoding.Must(enc.Encode(cat.ProductListing.Parts))
}
