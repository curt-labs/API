package customer_ctlr_new

import (
	"encoding/json"
	"github.com/curt-labs/GoAPI/helpers/encoding"
	"github.com/curt-labs/GoAPI/models/customer_new/content"
	"github.com/go-martini/martini"
	"io/ioutil"
	"net/http"
	"strconv"
)

// Get it all
func GetAllContent(w http.ResponseWriter, r *http.Request, enc encoding.Encoder) string {
	params := r.URL.Query()
	key := params.Get("key")

	content, err := custcontent.AllCustomerContent(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	return encoding.Must(enc.Encode(content))
}

// Get Content by Content Id
// Returns: CustomerContent
func GetContentById(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	qs := r.URL.Query()
	key := qs.Get("key")
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	content, err := custcontent.GetCustomerContent(id, key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	return encoding.Must(enc.Encode(content))
}

// Get Content by Content Id
// Returns: CustomerContent
func GetContentRevisionsById(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	qs := r.URL.Query()
	key := qs.Get("key")
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	revs, err := custcontent.GetCustomerContentRevisions(id, key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	return encoding.Must(enc.Encode(revs))
}

// Part Content Endpoints
func AllPartContent(w http.ResponseWriter, r *http.Request, enc encoding.Encoder) string {
	params := r.URL.Query()
	key := params.Get("key")

	content, err := custcontent.GetAllPartContent(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	return encoding.Must(enc.Encode(content))
}

func UniquePartContent(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	p := r.URL.Query()
	key := p.Get("key")
	partID, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	content, err := custcontent.GetPartContent(partID, key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	return encoding.Must(enc.Encode(content))
}

func CreatePartContent(w http.ResponseWriter, r *http.Request, params martini.Params, enc encoding.Encoder) string {
	// Get the key from the query string
	qs := r.URL.Query()
	key := qs.Get("key")
	id, err := strconv.Atoi(params["id"])
	// Defer the body closing until we're finished
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	var content custcontent.CustomerContent
	err = json.Unmarshal(body, &content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	if err = content.Save(id, 0, key); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	return encoding.Must(enc.Encode(content))
}
func UpdatePartContent(w http.ResponseWriter, r *http.Request, params martini.Params, enc encoding.Encoder) string {
	// Get the key from the query string
	qs := r.URL.Query()
	key := qs.Get("key")
	id, err := strconv.Atoi(params["id"])

	// Defer the body closing until we're finished
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	var content custcontent.CustomerContent
	err = json.Unmarshal(body, &content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	if err = content.Save(id, 0, key); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	return encoding.Must(enc.Encode(content))
}

func DeletePartContent(w http.ResponseWriter, r *http.Request, params martini.Params, enc encoding.Encoder) string {
	// Get the key from the query string
	qs := r.URL.Query()
	key := qs.Get("key")
	id, err := strconv.Atoi(params["id"])

	// Defer the body closing until we're finished
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	var content custcontent.CustomerContent
	err = json.Unmarshal(body, &content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	if err = content.Delete(id, 0, key); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	return encoding.Must(enc.Encode(content))
}

// Category Content Endpoints
func AllCategoryContent(w http.ResponseWriter, r *http.Request, enc encoding.Encoder) string {
	params := r.URL.Query()
	key := params.Get("key")

	content, err := custcontent.GetAllCategoryContent(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	return encoding.Must(enc.Encode(content))
}

func UniqueCategoryContent(w http.ResponseWriter, r *http.Request, params martini.Params, enc encoding.Encoder) string {
	qs := r.URL.Query()
	key := qs.Get("key")
	catID, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	content, err := custcontent.GetCategoryContent(catID, key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	return encoding.Must(enc.Encode(content))
}

func UpdateCategoryContent(w http.ResponseWriter, r *http.Request, params martini.Params, enc encoding.Encoder) string {
	// Get the key from the query string
	qs := r.URL.Query()
	key := qs.Get("key")
	id, err := strconv.Atoi(params["id"])

	// Defer the body closing until we're finished
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	var content custcontent.CustomerContent
	err = json.Unmarshal(body, &content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	if err = content.Save(0, id, key); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	return encoding.Must(enc.Encode(content))
}

func DeleteCategoryContent(w http.ResponseWriter, r *http.Request, params martini.Params, enc encoding.Encoder) string {
	// Get the key from the query string
	qs := r.URL.Query()
	key := qs.Get("key")
	id, err := strconv.Atoi(params["id"])

	// Defer the body closing until we're finished
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	var content custcontent.CustomerContent
	err = json.Unmarshal(body, &content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	if err = content.Delete(0, id, key); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	return encoding.Must(enc.Encode(content))
}

// Content Types
func GetAllContentTypes(w http.ResponseWriter, r *http.Request, enc encoding.Encoder) string {
	types, err := custcontent.AllCustomerContentTypes()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	return encoding.Must(enc.Encode(types))
}
