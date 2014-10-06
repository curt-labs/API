package site_new

import (
	"github.com/curt-labs/GoAPI/helpers/encoding"
	"github.com/curt-labs/GoAPI/models/site_new"
	"github.com/go-martini/martini"
	// "log"
	// "encoding/json"
	// "io/ioutil"
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

	err = c.Get()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusNoContent)
		return ""
	}

	return encoding.Must(enc.Encode(c))
}
