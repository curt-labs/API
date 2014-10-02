package site

import (
	"github.com/curt-labs/GoAPI/helpers/encoding"
	"github.com/curt-labs/GoAPI/models/site"
	"github.com/go-martini/martini"
	"net/http"
	"strconv"
)

func GetPrimaryMenu(rw http.ResponseWriter, req *http.Request, enc encoding.Encoder) string {
	var m site.MenuWithContent
	err := m.GetPrimaryMenu()

	if err != nil {
		http.Error(rw, err.Error(), http.StatusNoContent)
		return ""
	}

	return encoding.Must(enc.Encode(m))
}

func GetMenuWithContent(rw http.ResponseWriter, req *http.Request, enc encoding.Encoder, params martini.Params) string {
	var err error
	var m site.MenuWithContent
	err = req.ParseForm()

	auth, err := strconv.ParseBool(req.FormValue("auth"))
	if err != nil || auth == false {
		http.Error(rw, err.Error(), http.StatusForbidden)
		return ""
	}

	id, err := strconv.Atoi(params["idOrName"])
	if err != nil {
		name := req.FormValue("idOrName")
		err = m.GetMenuWithContentByName(name)
	} else {
		err = m.GetMenuByContentId(id, auth)
	}

	if err != nil {
		http.Error(rw, err.Error(), http.StatusNoContent)
		return ""
	}

	return encoding.Must(enc.Encode(m))
}

func GetMenuByContentId(rw http.ResponseWriter, req *http.Request, enc encoding.Encoder, params martini.Params) string {
	var err error
	var m site.MenuWithContent
	err = req.ParseForm()

	auth, err := strconv.ParseBool(req.FormValue("auth"))
	if err != nil || auth == false {
		http.Error(rw, err.Error(), http.StatusForbidden)
		return ""
	}

	contentId, err := strconv.Atoi(req.FormValue("contentId"))
	menuId, err := strconv.Atoi(req.FormValue("menuId"))
	m.Menu.Id = menuId

	err = m.GetMenuByContentId(contentId, auth)

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

func GetMenuSitemap(rw http.ResponseWriter, req *http.Request, enc encoding.Encoder) string {
	menus, err := site.GetMenuSitemap()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusNoContent)
		return ""
	}

	return encoding.Must(enc.Encode(menus))
}
