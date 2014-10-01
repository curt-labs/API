package site

import (
	"github.com/curt-labs/GoAPI/helpers/encoding"
	"github.com/curt-labs/GoAPI/models/site"
	"log"
	"net/http"
	"strconv"
)

func GetPrimaryMenu(rw http.ResponseWriter, req *http.Request, enc encoding.Encoder) string {
	var m site.Menu
	err := m.GetPrimaryMenu()

	if err != nil {
		http.Error(rw, err.Error(), http.StatusNoContent)
		return ""
	}

	return encoding.Must(enc.Encode(m))
}

func GetMenu(rw http.ResponseWriter, req *http.Request, enc encoding.Encoder) string {
	var err error
	var m site.Menu
	err = req.ParseForm()

	id, err := strconv.Atoi(req.FormValue("id"))
	if err != nil {
		name := req.FormValue("name")
		m.Name = name
		err = m.GetByName()
	} else {
		m.Id = id
		err = m.Get()
	}

	if err != nil {
		http.Error(rw, err.Error(), http.StatusNoContent)
		return ""
	}

	return encoding.Must(enc.Encode(m))
}

func GetMenuContent(rw http.ResponseWriter, req *http.Request, enc encoding.Encoder) string {
	var err error
	var m site.Menu
	err = req.ParseForm()

	id, err := strconv.Atoi(req.FormValue("id"))
	m.Id = id

	err = m.GetMenuContents()
	// err = m.GetContentPages()

	if err != nil {
		log.Print(err)
		http.Error(rw, err.Error(), http.StatusNoContent)
		return ""
	}

	return encoding.Must(enc.Encode(m))
}

func GetMenuWithContent(rw http.ResponseWriter, req *http.Request, enc encoding.Encoder) string {
	var err error
	var m site.Menu
	err = req.ParseForm()

	auth, err := strconv.ParseBool(req.FormValue("auth"))
	if err != nil {
		http.Error(rw, err.Error(), http.StatusForbidden)
		return ""
	}
	m.RequireAuthentication = auth
	id, err := strconv.Atoi(req.FormValue("id"))
	if err != nil {
		slug := req.FormValue("name")
		err = m.GetMenuWithContentByName(slug)
		log.Print(slug)
	} else {
		err = m.GetMenuByContentId(id)
		log.Print(id)
	}

	if err != nil {
		http.Error(rw, err.Error(), http.StatusNoContent)
		return ""
	}

	return encoding.Must(enc.Encode(m))
}

func GetFooterSitemap(rw http.ResponseWriter, req *http.Request, enc encoding.Encoder) string {
	menus, err := site.GetFooterSitemap()

	if err != nil {
		http.Error(rw, err.Error(), http.StatusNoContent)
		return ""
	}

	return encoding.Must(enc.Encode(menus))
}
