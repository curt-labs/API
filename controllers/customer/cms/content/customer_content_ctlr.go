package customer_content_ctlr

import (
	"../../../../helpers/plate"
	"../../../../models/cms/customer"
	"net/http"
	"strconv"
)

// Get it all
func AllContent(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	key := params.Get("key")

	content, err := customer_cms.AllContent(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	plate.ServeFormatted(w, r, content)

}

// Part Content Endpoints
func AllPartContent(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	key := params.Get("key")

	content, err := customer_cms.AllPartContent(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	plate.ServeFormatted(w, r, content)
}

func PartContent(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	key := params.Get("key")
	partID, err := strconv.Atoi(params.Get(":id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	content, err := customer_cms.GetPartContent(partID, key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	plate.ServeFormatted(w, r, content)
}

func UpdatePartContent(w http.ResponseWriter, r *http.Request) {
}

func DeletePartContent(w http.ResponseWriter, r *http.Request) {
}

// Category Content Endpoints
func AllCategoryContent(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	key := params.Get("key")

	content, err := customer_cms.AllCategoryContent(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	plate.ServeFormatted(w, r, content)
}

func CategoryContent(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	key := params.Get("key")
	catID, err := strconv.Atoi(params.Get(":id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	content, err := customer_cms.GetCategoryContent(catID, key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	plate.ServeFormatted(w, r, content)
}

func UpdateCategoryContent(w http.ResponseWriter, r *http.Request) {
}

func DeleteCategoryContent(w http.ResponseWriter, r *http.Request) {
}
