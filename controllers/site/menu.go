package site

import (
	"github.com/curt-labs/GoAPI/helpers/encoding"
	"github.com/curt-labs/GoAPI/models/site"
	"github.com/go-martini/martini"
	// "log"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
)

func GetMenu(rw http.ResponseWriter, req *http.Request, enc encoding.Encoder, params martini.Params) string {
	var m site.Menu
	var err error
	idStr := params["id"]
	id, err := strconv.Atoi(idStr)

	if err == nil {
		//Thar be an Id int
		m.Id = id
		err = m.Get()
		if err != nil {
			http.Error(rw, err.Error(), http.StatusNoContent)
			return ""
		}
	} else {
		//Thar be a name
		m.Name = idStr
		err = m.GetByName()
		if err != nil {
			http.Error(rw, err.Error(), http.StatusNoContent)
			return ""
		}
	}
	return encoding.Must(enc.Encode(m))
}

func GetAllMenus(rw http.ResponseWriter, req *http.Request, enc encoding.Encoder, params martini.Params) string {
	m, err := site.GetAllMenus()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusNoContent)
		return ""
	}
	return encoding.Must(enc.Encode(m))
}

func GetMenuWithContents(rw http.ResponseWriter, req *http.Request, enc encoding.Encoder, params martini.Params) string {
	var m site.Menu
	var err error
	idStr := params["id"]
	id, err := strconv.Atoi(idStr)

	if err == nil {
		//Thar be an Id int
		m.Id = id
		err = m.Get()
		if err != nil {
			http.Error(rw, err.Error(), http.StatusNoContent)
			return ""
		}
	} else {
		//Thar be a name
		m.Name = idStr
		err = m.GetByName()
		if err != nil {
			http.Error(rw, err.Error(), http.StatusNoContent)
			return ""
		}
	}
	err = m.GetContents()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusNoContent)
		return ""
	}

	return encoding.Must(enc.Encode(m))
}

func SaveMenu(rw http.ResponseWriter, req *http.Request, enc encoding.Encoder, params martini.Params) string {
	var m site.Menu
	var err error
	idStr := params["id"]
	if idStr != "" {
		m.Id, err = strconv.Atoi(idStr)
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
	if m.Id > 0 {
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

func DeleteMenu(rw http.ResponseWriter, req *http.Request, enc encoding.Encoder, params martini.Params) string {
	var err error
	var m site.Menu

	idStr := params["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return ""
	}
	m.Id = id
	err = m.Delete()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return ""
	}

	return encoding.Must(enc.Encode(m))
}
