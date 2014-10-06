package site_new

import (
	"github.com/curt-labs/GoAPI/helpers/encoding"
	"github.com/curt-labs/GoAPI/models/site_new"
	"github.com/go-martini/martini"
	// "log"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
)

func GetContent(rw http.ResponseWriter, req *http.Request, enc encoding.Encoder, params martini.Params) string {
	var c site_new.Content
	var err error
	idStr := params["id"]
	id, err := strconv.Atoi(idStr)

	if err == nil {
		//Thar be an Id int
		c.Id = id
		err = c.Get()
		if err != nil {
			http.Error(rw, err.Error(), http.StatusNoContent)
			return ""
		}
	} else {
		//Thar be a slug
		c.Slug = idStr
		err = c.GetBySlug()
		if err != nil {
			http.Error(rw, err.Error(), http.StatusNoContent)
			return ""
		}
	}
	return encoding.Must(enc.Encode(c))
}

func GetAllContents(rw http.ResponseWriter, req *http.Request, enc encoding.Encoder, params martini.Params) string {
	m, err := site_new.GetAllContents()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusNoContent)
		return ""
	}
	return encoding.Must(enc.Encode(m))
}

func GetContentRevisions(rw http.ResponseWriter, req *http.Request, enc encoding.Encoder, params martini.Params) string {
	var c site_new.Content
	var err error
	idStr := params["id"]
	c.Id, err = strconv.Atoi(idStr)

	err = c.GetContentRevisions()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusNoContent)
		return ""
	}

	return encoding.Must(enc.Encode(c))
}

func SaveContent(rw http.ResponseWriter, req *http.Request, enc encoding.Encoder, params martini.Params) string {
	var c site_new.Content
	var err error
	idStr := params["id"]
	if idStr != "" {
		c.Id, err = strconv.Atoi(idStr)
		err = c.Get()
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
	err = json.Unmarshal(requestBody, &c)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return encoding.Must(enc.Encode(false))
	}
	//create or update
	if c.Id > 0 {
		err = c.Update()
	} else {
		err = c.Create()
	}

	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return ""
	}
	return encoding.Must(enc.Encode(c))
}

func DeleteContent(rw http.ResponseWriter, req *http.Request, enc encoding.Encoder, params martini.Params) string {
	var err error
	var c site_new.Content

	idStr := params["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return ""
	}
	c.Id = id
	err = c.Delete()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return ""
	}

	return encoding.Must(enc.Encode(c))
}
