package category_ctlr

import (
	"github.com/curt-labs/GoAPI/helpers/apifilter"
	"github.com/curt-labs/GoAPI/helpers/encoding"
	"github.com/curt-labs/GoAPI/models/products"
	"github.com/go-martini/martini"
	"net/http"
	"strconv"
)

func GetCategory(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	qs := r.URL.Query()
	key := qs.Get("key")
	id, err := strconv.Atoi(params["id"])

	var cat products.Category
	if err != nil {
		cat, err = products.GetCategoryByTitle(params["id"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return ""
		}
	} else {
		cat.CategoryId = id
	}

	ext, err := cat.GetCategory(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	filters, err := apifilter.CategoryFilter(ext, nil)
	if err == nil {
		ext.Filter = filters
	}

	return encoding.Must(enc.Encode(ext))
}

func Parents(w http.ResponseWriter, r *http.Request, enc encoding.Encoder) string {

	cats, err := products.TopTierCategories()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	return encoding.Must(enc.Encode(cats))
}

func SubCategories(w http.ResponseWriter, r *http.Request, enc encoding.Encoder) string {
	params := r.URL.Query()
	id, err := strconv.Atoi(params.Get(":id"))

	var cat products.Category
	if err != nil {
		cat, err = products.GetCategoryByTitle(params.Get(":id"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return ""
		}
	} else {
		cat.CategoryId = id
	}

	subs, err := cat.GetSubCategories()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	return encoding.Must(enc.Encode(subs))
}

func GetParts(w http.ResponseWriter, r *http.Request, enc encoding.Encoder) string {
	params := r.URL.Query()
	key := params.Get("key")
	catID, err := strconv.Atoi(params.Get(":id"))

	var cat products.Category
	if err != nil {
		title := params.Get(":id")
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
		cat.CategoryId = catID
	}

	count, err := strconv.Atoi(params.Get(":count"))
	if err != nil {
		count = 5
	}

	page, err := strconv.Atoi(params.Get(":page"))
	if err != nil {
		page = 0
	} else {
		page = page - 1
	}

	parts, err := products.GetCategoryParts(cat, key, page, count)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	if parts == nil {
		parts = make([]products.Part, 0)
	}

	return encoding.Must(enc.Encode(parts))
}
