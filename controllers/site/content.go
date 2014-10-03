package site

import (
	"github.com/curt-labs/GoAPI/helpers/encoding"
	"github.com/curt-labs/GoAPI/models/site"
	"github.com/go-martini/martini"
	"net/http"
	"strconv"
)

func GetContentPage(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	var err error
	var cp site.ContentPage
	err = r.ParseForm()
	authenticated, err := strconv.ParseBool(r.FormValue("auth"))
	if err != nil || authenticated == false {
		http.Error(w, err.Error(), http.StatusForbidden)
		return ""
	}
	name := r.FormValue("name")
	menuId, err := strconv.Atoi(r.FormValue("menuId"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return ""
	}

	cp.SiteContent.Slug = name
	err = cp.GetContentPageByName(menuId, authenticated)

	return encoding.Must(enc.Encode(cp))

}

func GetPrimaryContentPage(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	var err error
	var cp site.ContentPage
	err = cp.GetPrimaryContentPage()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return ""
	}
	return encoding.Must(enc.Encode(cp))
}

//super clumsy function; borrowed from v2; TODO--trash it or fix it
func GetSitemapCP(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	var err error
	var cps site.ContentPages
	cps, err = site.GetSitemapCP()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return ""
	}
	return encoding.Must(enc.Encode(cps))
}

func GetLandingPage(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	var err error
	var lp site.LandingPage
	err = r.ParseForm()
	lp.Id, err = strconv.Atoi(r.FormValue("id"))

	err = lp.Get()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return ""
	}

	return encoding.Must(enc.Encode(lp))
}
