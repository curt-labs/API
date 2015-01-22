package customer_ctlr

import (
	"github.com/curt-labs/GoAPI/helpers/apicontext"
	"github.com/curt-labs/GoAPI/helpers/encoding"
	"github.com/curt-labs/GoAPI/helpers/error"
	"github.com/curt-labs/GoAPI/models/customer/content"
	"github.com/go-martini/martini"

	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
)

// Get it all
func GetAllContent(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder, dtx *apicontext.DataContext) string {
	content, err := custcontent.AllCustomerContent(dtx.APIKey)
	if err != nil {
		apierror.GenerateError("Trouble getting customer content", err, rw, r)
		return ""
	}

	return encoding.Must(enc.Encode(content))
}

// Get Content by Content Id
// Returns: CustomerContent
func GetContentById(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params, dtx *apicontext.DataContext) string {
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		apierror.GenerateError("Trouble getting customer content ID", err, rw, r)
		return ""
	}

	content, err := custcontent.GetCustomerContent(id, dtx.APIKey)
	if err != nil {
		apierror.GenerateError("Trouble getting customer content", err, rw, r)
		return ""
	}

	return encoding.Must(enc.Encode(content))
}

// Get Content by Content Id
// Returns: CustomerContent
func GetContentRevisionsById(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params, dtx *apicontext.DataContext) string {
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		apierror.GenerateError("Trouble getting customer content ID", err, rw, r)
		return ""
	}

	revs, err := custcontent.GetCustomerContentRevisions(id, dtx.APIKey)
	if err != nil {
		apierror.GenerateError("Trouble getting customer content revisions", err, rw, r)
		return ""
	}

	return encoding.Must(enc.Encode(revs))
}

// Part Content Endpoints
func AllPartContent(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder, dtx *apicontext.DataContext) string {
	content, err := custcontent.GetAllPartContent(dtx.APIKey)
	if err != nil {
		apierror.GenerateError("Trouble getting all part content", err, rw, r)
		return ""
	}

	return encoding.Must(enc.Encode(content))
}

func UniquePartContent(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params, dtx *apicontext.DataContext) string {
	partID, err := strconv.Atoi(params["id"])
	if err != nil {
		apierror.GenerateError("Trouble getting part ID", err, rw, r)
		return ""
	}

	content, err := custcontent.GetPartContent(partID, dtx.APIKey)
	if err != nil {
		apierror.GenerateError("Trouble getting part content", err, rw, r)
		return ""
	}

	return encoding.Must(enc.Encode(content))
}

func CreatePartContent(rw http.ResponseWriter, r *http.Request, params martini.Params, enc encoding.Encoder, dtx *apicontext.DataContext) string {
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		apierror.GenerateError("Trouble getting part ID", err, rw, r)
		return ""
	}

	// Defer the body closing until we're finished
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		apierror.GenerateError("Trouble reading request body while creating part content", err, rw, r)
		return ""
	}

	var content custcontent.CustomerContent
	if err = json.Unmarshal(body, &content); err != nil {
		apierror.GenerateError("Trouble unmarshalling json request body while creating part content", err, rw, r)
		return ""
	}

	if err = content.Save(id, 0, dtx.APIKey); err != nil {
		apierror.GenerateError("Trouble creating part content", err, rw, r)
		return ""
	}

	return encoding.Must(enc.Encode(content))
}
func UpdatePartContent(rw http.ResponseWriter, r *http.Request, params martini.Params, enc encoding.Encoder, dtx *apicontext.DataContext) string {
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		apierror.GenerateError("Trouble getting part ID", err, rw, r)
		return ""
	}

	// Defer the body closing until we're finished
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		apierror.GenerateError("Trouble reading request body while updating part content", err, rw, r)
		return ""
	}

	var content custcontent.CustomerContent
	if err = json.Unmarshal(body, &content); err != nil {
		apierror.GenerateError("Trouble unmarshalling json request body while updating part content", err, rw, r)
		return ""
	}

	if err = content.Save(id, 0, dtx.APIKey); err != nil {
		apierror.GenerateError("Trouble updating part content", err, rw, r)
		return ""
	}

	return encoding.Must(enc.Encode(content))
}

func DeletePartContent(rw http.ResponseWriter, r *http.Request, params martini.Params, enc encoding.Encoder, dtx *apicontext.DataContext) string {
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		apierror.GenerateError("Trouble getting part ID", err, rw, r)
		return ""
	}

	// Defer the body closing until we're finished
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		apierror.GenerateError("Trouble reading request body while deleting part content", err, rw, r)
		return ""
	}

	var content custcontent.CustomerContent
	if err = json.Unmarshal(body, &content); err != nil {
		apierror.GenerateError("Trouble unmarshalling json request body while deleting part content", err, rw, r)
		return ""
	}

	if err = content.Delete(id, 0, dtx.APIKey); err != nil {
		apierror.GenerateError("Trouble deleting part content", err, rw, r)
		return ""
	}

	return encoding.Must(enc.Encode(content))
}

// Category Content Endpoints
func AllCategoryContent(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder, dtx *apicontext.DataContext) string {
	content, err := custcontent.GetAllCategoryContent(dtx.APIKey)
	if err != nil {
		apierror.GenerateError("Trouble getting all category content", err, rw, r)
		return ""
	}

	return encoding.Must(enc.Encode(content))
}

func UniqueCategoryContent(rw http.ResponseWriter, r *http.Request, params martini.Params, enc encoding.Encoder, dtx *apicontext.DataContext) string {
	catID, err := strconv.Atoi(params["id"])
	if err != nil {
		apierror.GenerateError("Trouble getting category ID", err, rw, r)
		return ""
	}

	content, err := custcontent.GetCategoryContent(catID, dtx.APIKey)
	if err != nil {
		apierror.GenerateError("Trouble getting category content", err, rw, r)
		return ""
	}

	return encoding.Must(enc.Encode(content))
}

func CreateCategoryContent(rw http.ResponseWriter, r *http.Request, params martini.Params, enc encoding.Encoder, dtx *apicontext.DataContext) string {
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		apierror.GenerateError("Trouble getting category ID", err, rw, r)
		return ""
	}

	// Defer the body closing until we're finished
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		apierror.GenerateError("Trouble reading request body while creating category content", err, rw, r)
		return ""
	}

	var content custcontent.CustomerContent
	if err = json.Unmarshal(body, &content); err != nil {
		apierror.GenerateError("Trouble unmarshalling json request body while creating category content", err, rw, r)
		return ""
	}

	if err = content.Save(0, id, dtx.APIKey); err != nil {
		apierror.GenerateError("Trouble creating category content", err, rw, r)
		return ""
	}

	return encoding.Must(enc.Encode(content))
}

func UpdateCategoryContent(rw http.ResponseWriter, r *http.Request, params martini.Params, enc encoding.Encoder, dtx *apicontext.DataContext) string {
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		apierror.GenerateError("Trouble getting category ID", err, rw, r)
		return ""
	}

	// Defer the body closing until we're finished
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		apierror.GenerateError("Trouble reading request body while updating category content", err, rw, r)
		return ""
	}

	var content custcontent.CustomerContent
	if err = json.Unmarshal(body, &content); err != nil {
		apierror.GenerateError("Trouble unmarshalling json request body while updating category content", err, rw, r)
		return ""
	}

	if err = content.Save(0, id, dtx.APIKey); err != nil {
		apierror.GenerateError("Trouble updating category content", err, rw, r)
		return ""
	}

	return encoding.Must(enc.Encode(content))
}

func DeleteCategoryContent(rw http.ResponseWriter, r *http.Request, params martini.Params, enc encoding.Encoder, dtx *apicontext.DataContext) string {
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		apierror.GenerateError("Trouble getting category ID", err, rw, r)
		return ""
	}

	// Defer the body closing until we're finished
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		apierror.GenerateError("Trouble reading request body while deleting category content", err, rw, r)
		return ""
	}

	var content custcontent.CustomerContent
	if err = json.Unmarshal(body, &content); err != nil {
		apierror.GenerateError("Trouble unmarshalling json request body while deleting category content", err, rw, r)
		return ""
	}

	if err = content.Delete(0, id, dtx.APIKey); err != nil {
		apierror.GenerateError("Trouble deleting category content", err, rw, r)
		return ""
	}

	return encoding.Must(enc.Encode(content))
}

// Content Types
func GetAllContentTypes(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder) string {
	types, err := custcontent.AllCustomerContentTypes()
	if err != nil {
		apierror.GenerateError("Trouble getting all content types", err, rw, r)
		return ""
	}

	return encoding.Must(enc.Encode(types))
}
