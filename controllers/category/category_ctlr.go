package category_ctlr

import (
	"../../helpers/plate"
	. "../../models"
	"net/http"
	"strconv"
)

func GetCategory(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	key := params.Get(":key")
	id, err := strconv.Atoi(params.Get(":id"))

	var cat Category
	if err != nil {
		cat, err = GetByTitle(params.Get(":id"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		cat.CategoryId = id
	}

	ext, err := cat.GetCategory(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	plate.ServeFormatted(w, r, ext)
}

func Parents(w http.ResponseWriter, r *http.Request) {

	cats, err := TopTierCategories()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	plate.ServeFormatted(w, r, cats)
}

func SubCategories(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	id, err := strconv.Atoi(params.Get(":id"))

	var cat Category
	if err != nil {
		cat, err = GetByTitle(params.Get(":id"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		cat.CategoryId = id
	}

	subs, err := cat.SubCategories()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	plate.ServeFormatted(w, r, subs)
}

func GetParts(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	key := params.Get("key")
	catID, err := strconv.Atoi(params.Get(":id"))

	var cat Category
	if err != nil {
		title := params.Get(":id")
		if title == "" {
			http.Error(w, "Invalid Category", http.StatusInternalServerError)
			return
		}
		cat, err = GetByTitle(title)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		cat.CategoryId = catID
	}

	count, err := strconv.Atoi(params.Get(":count"))
	if err != nil {
		count = 10
	}

	page, err := strconv.Atoi(params.Get(":page"))
	if err != nil {
		page = 0
	} else {
		page = page - 1
	}

	parts, err := cat.GetCategoryParts(key, page, count)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if parts == nil {
		parts = make([]Part, 0)
	}

	plate.ServeFormatted(w, r, parts)

}
