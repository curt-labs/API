package webProperty_controller

import (
	"encoding/json"
	"github.com/curt-labs/GoAPI/helpers/apicontext"
	"github.com/curt-labs/GoAPI/helpers/encoding"
	"github.com/curt-labs/GoAPI/models/webProperty"
	"github.com/go-martini/martini"
	"io/ioutil"
	"net/http"
	"strconv"
)

//Parses JSON input from the body, a la Angular http
func Save_Json(rw http.ResponseWriter, req *http.Request, enc encoding.Encoder, params martini.Params, dtx *apicontext.DataContext) string {
	var w webProperty_model.WebProperty
	var err error
	idStr := params["id"]
	if idStr != "" {
		w.ID, err = strconv.Atoi(idStr)
		err = w.Get(dtx)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return ""
		}
	}

	//json
	requestBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return ""
	}
	err = json.Unmarshal(requestBody, &w)
	if err != nil {

		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return ""
	}

	//create or update
	if w.ID > 0 {
		err = w.Update()
	} else {
		err = w.Create()
	}

	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return ""
	}
	return encoding.Must(enc.Encode(w))
}

func SaveType_Json(rw http.ResponseWriter, req *http.Request, enc encoding.Encoder, params martini.Params) string {
	var w webProperty_model.WebPropertyType
	var err error
	idStr := params["id"]
	if idStr != "" {
		w.ID, err = strconv.Atoi(idStr)
		err = w.Get()
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return ""
		}
	}

	//json
	requestBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return ""
	}
	err = json.Unmarshal(requestBody, &w)
	if err != nil {

		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return ""
	}

	//create or update
	if w.ID > 0 {
		err = w.Update()
	} else {
		err = w.Create()
	}

	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return ""
	}
	return encoding.Must(enc.Encode(w))
}

func SaveRequirement_Json(rw http.ResponseWriter, req *http.Request, enc encoding.Encoder, params martini.Params) string {
	var w webProperty_model.WebPropertyRequirement
	var err error
	idStr := params["id"]
	if idStr != "" {
		w.ID, err = strconv.Atoi(idStr)
		err = w.Get()
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return ""
		}
	}

	//json
	requestBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return ""
	}
	err = json.Unmarshal(requestBody, &w)
	if err != nil {

		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return ""
	}

	//create or update
	if w.ID > 0 {
		err = w.Update()
	} else {
		err = w.Create()
	}

	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return ""
	}
	return encoding.Must(enc.Encode(w))
}

func SaveNote_Json(rw http.ResponseWriter, req *http.Request, enc encoding.Encoder, params martini.Params) string {
	var w webProperty_model.WebPropertyNote
	var err error
	idStr := params["id"]
	if idStr != "" {
		w.ID, err = strconv.Atoi(idStr)
		err = w.Get()
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return ""
		}
	}

	//json
	requestBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return ""
	}
	err = json.Unmarshal(requestBody, &w)
	if err != nil {

		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return ""
	}

	//create or update
	if w.ID > 0 {
		err = w.Update()
	} else {
		err = w.Create()
	}

	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return ""
	}
	return encoding.Must(enc.Encode(w))
}
