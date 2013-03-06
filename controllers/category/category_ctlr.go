package category_ctlr

import (
	. "../../models"
	"../../plate"
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
