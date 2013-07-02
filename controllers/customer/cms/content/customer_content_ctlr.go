package customer_content_ctlr

import (
	"../../../../helpers/plate"
	"../../../../models/cms/customer"
	"encoding/json"
	"io/ioutil"
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

	var content customer_cms.CustomerContent
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

	var content customer_cms.CustomerContent
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

	var content customer_cms.CustomerContent
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

	var content customer_cms.CustomerContent
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
