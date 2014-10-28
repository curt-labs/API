package applicationGuide

import (
	"encoding/json"
	"github.com/curt-labs/GoAPI/helpers/encoding"
	"github.com/curt-labs/GoAPI/models/applicationGuide"
	"github.com/go-martini/martini"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

func GetApplicationGuide(rw http.ResponseWriter, req *http.Request, enc encoding.Encoder, params martini.Params) string {
	var err error
	var ag applicationGuide.ApplicationGuide
	id := params["id"]
	ag.ID, err = strconv.Atoi(id)

	err = ag.Get()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return err.Error()
	}
	return encoding.Must(enc.Encode(ag))
}

func GetApplicationGuidesByWebsite(rw http.ResponseWriter, req *http.Request, enc encoding.Encoder, params martini.Params) string {
	var err error
	var ag applicationGuide.ApplicationGuide
	id := params["id"]
	ag.Website.ID, err = strconv.Atoi(id)

	ags, err := ag.GetBySite()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return err.Error()
	}
	return encoding.Must(enc.Encode(ags))
}

func CreateApplicationGuide(rw http.ResponseWriter, req *http.Request, enc encoding.Encoder, params martini.Params) string {
	contType := req.Header.Get("Content-Type")

	var ag applicationGuide.ApplicationGuide
	var err error

	// if contType == "application/json" {
	if strings.Contains(contType, "application/json") {
		//json
		requestBody, err := ioutil.ReadAll(req.Body)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return encoding.Must(enc.Encode(false))
		}

		err = json.Unmarshal(requestBody, &ag)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return encoding.Must(enc.Encode(false))
		}
	} else {
		//else, form
		ag.Url = req.FormValue("url")
		ag.Website.ID, err = strconv.Atoi(req.FormValue("website_id"))
		ag.FileType = req.FormValue("file_type")
		ag.Category.ID, err = strconv.Atoi(req.FormValue("category_id"))

		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return err.Error()
		}
	}
	err = ag.Create()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return err.Error()
	}

	//Return JSON
	return encoding.Must(enc.Encode(ag))
}

func DeleteApplicationGuide(rw http.ResponseWriter, req *http.Request, enc encoding.Encoder, params martini.Params) string {
	var err error
	var ag applicationGuide.ApplicationGuide
	id, err := strconv.Atoi(params["id"])
	ag.ID = id
	err = ag.Delete()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return err.Error()
	}

	//Return JSON
	return encoding.Must(enc.Encode(ag))
}
