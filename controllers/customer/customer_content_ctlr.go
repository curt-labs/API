package customer_ctlr

import (
	"../../helpers/plate"
	. "../../models"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
)

// Get it all
func GetAllContent(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	key := params.Get("key")

	content, err := AllCustomerContent(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	plate.ServeFormatted(w, r, content)
}

// Get Content by Content Id
// Returns: CustomerContent
func GetContentById(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	key := params.Get("key")
	id, err := strconv.Atoi(params.Get(":id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	content, err := GetCustomerContent(id, key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	plate.ServeFormatted(w, r, content)

}

// Get Content by Content Id
// Returns: CustomerContent
func GetContentRevisionsById(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	key := params.Get("key")
	id, err := strconv.Atoi(params.Get(":id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	revs, err := GetCustomerContentRevisions(id, key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	plate.ServeFormatted(w, r, revs)

}

// Part Content Endpoints
func AllPartContent(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	key := params.Get("key")

	content, err := GetAllPartContent(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	plate.ServeFormatted(w, r, content)
}

func UniquePartContent(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	key := params.Get("key")
	partID, err := strconv.Atoi(params.Get(":id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	content, err := GetPartContent(partID, key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	plate.ServeFormatted(w, r, content)
}

func UpdatePartContent(w http.ResponseWriter, r *http.Request) {
	// Get the key from the query string
	params := r.URL.Query()
	key := params.Get("key")
	id, err := strconv.Atoi(params.Get(":id"))

	// Defer the body closing until we're finished
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var content CustomerContent
	err = json.Unmarshal(body, &content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = content.Save(id, 0, key); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	plate.ServeFormatted(w, r, content)
}

func DeletePartContent(w http.ResponseWriter, r *http.Request) {
	// Get the key from the query string
	params := r.URL.Query()
	key := params.Get("key")
	id, err := strconv.Atoi(params.Get(":id"))

	// Defer the body closing until we're finished
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var content CustomerContent
	err = json.Unmarshal(body, &content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = content.Delete(id, 0, key); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	plate.ServeFormatted(w, r, content)
}

// Category Content Endpoints
func AllCategoryContent(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	key := params.Get("key")

	content, err := GetAllCategoryContent(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	plate.ServeFormatted(w, r, content)
}

func UniqueCategoryContent(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	key := params.Get("key")
	catID, err := strconv.Atoi(params.Get(":id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	content, err := GetCategoryContent(catID, key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	plate.ServeFormatted(w, r, content)
}

func UpdateCategoryContent(w http.ResponseWriter, r *http.Request) {
	// Get the key from the query string
	params := r.URL.Query()
	key := params.Get("key")
	id, err := strconv.Atoi(params.Get(":id"))

	// Defer the body closing until we're finished
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var content CustomerContent
	err = json.Unmarshal(body, &content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = content.Save(0, id, key); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	plate.ServeFormatted(w, r, content)
}

func DeleteCategoryContent(w http.ResponseWriter, r *http.Request) {
	// Get the key from the query string
	params := r.URL.Query()
	key := params.Get("key")
	id, err := strconv.Atoi(params.Get(":id"))

	// Defer the body closing until we're finished
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var content CustomerContent
	err = json.Unmarshal(body, &content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = content.Delete(0, id, key); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	plate.ServeFormatted(w, r, content)
}

// Content Types
func GetAllContentTypes(w http.ResponseWriter, r *http.Request) {
	types, err := AllCustomerContentTypes()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	plate.ServeFormatted(w, r, types)
}
