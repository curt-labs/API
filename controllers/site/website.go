package site

import (
	"encoding/json"
	"github.com/curt-labs/GoAPI/helpers/encoding"
	"github.com/curt-labs/GoAPI/models/site"
	"github.com/go-martini/martini"
	"io/ioutil"
	"net/http"
	"strconv"
)

func GetSiteDetails(rw http.ResponseWriter, req *http.Request, enc encoding.Encoder, params martini.Params) string {
	var w site.Website
	var err error
	idStr := params["id"]
	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusNoContent)
		return ""
	}
	w.ID = id

	err = w.GetDetails()
	return encoding.Must(enc.Encode(w))
}

func SaveSite(rw http.ResponseWriter, req *http.Request, enc encoding.Encoder, params martini.Params) string {
	var m site.Website
	var err error
	idStr := params["id"]
	if idStr != "" {
		m.ID, err = strconv.Atoi(idStr)
		err = m.Get()
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return ""
		}
	}

	//json
	requestBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return encoding.Must(enc.Encode(false))
	}
	err = json.Unmarshal(requestBody, &m)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return encoding.Must(enc.Encode(false))
	}
	//create or update
	if m.ID > 0 {
		err = m.Update()
	} else {
		err = m.Create()
	}

	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return ""
	}
	return encoding.Must(enc.Encode(m))
}

func DeleteSite(rw http.ResponseWriter, req *http.Request, enc encoding.Encoder, params martini.Params) string {
	var err error
	var m site.Website

	idStr := params["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return ""
	}
	m.ID = id
	err = m.Delete()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return ""
	}

	return encoding.Must(enc.Encode(m))
}
