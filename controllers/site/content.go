package site

import (
	"github.com/curt-labs/GoAPI/helpers/encoding"
	"github.com/curt-labs/GoAPI/models/site"
	"github.com/go-martini/martini"
	// "log"
	"net/http"
	"strconv"
)

func GetContentPage(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	var err error
	var cp site.ContentPage
	err = r.ParseForm()
	authenticated, err := strconv.ParseBool(r.FormValue("authenticated"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return ""
	}

	if authenticated == false {
		http.Error(w, "Unauthorized.", http.StatusUnauthorized)
		return ""
	}
	cp.Menu.RequireAuthentication = authenticated
	cp.SiteContent.Slug = r.FormValue("name")
	cp.SiteContent.Id, err = strconv.Atoi(r.FormValue("id"))

	if cp.SiteContent.Slug != "" {
		err = cp.SiteContent.GetBySlug()
		if err != nil {
			http.Error(w, err.Error(), http.StatusNoContent)
			return ""
		}
	}
	if cp.SiteContent.Id > 0 {
		err = cp.SiteContent.Get()
		if err != nil {
			http.Error(w, err.Error(), http.StatusNoContent)
			return ""
		}
	}

	err = cp.GetRevision()
	if err != nil {
		http.Error(w, err.Error(), http.StatusNoContent)
		return ""
	}

	err = cp.Menu.GetMenuByContentId(cp.SiteContent.Id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNoContent)
		return ""
	}

	return encoding.Must(enc.Encode(cp))
}
