package applicationGuide

import (
	"encoding/json"
	"github.com/curt-labs/GoAPI/helpers/apicontext"
	"github.com/curt-labs/GoAPI/helpers/encoding"
	"github.com/curt-labs/GoAPI/helpers/error"
	"github.com/curt-labs/GoAPI/models/applicationGuide"
	"github.com/go-martini/martini"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func GetApplicationGuide(rw http.ResponseWriter, req *http.Request, enc encoding.Encoder, params martini.Params, dtx *apicontext.DataContext) string {
	var err error
	var ag applicationGuide.ApplicationGuide
	id := params["id"]
	ag.ID, err = strconv.Atoi(id)
	if err != nil {
		apierror.GenerateError("Trouble converting ID parameter", err, rw, req)
	}

	err = ag.Get(dtx)
	if err != nil {
		apierror.GenerateError("Error getting Application Guide", err, rw, req)
	}
	return encoding.Must(enc.Encode(ag))
}

func GetApplicationGuidesByWebsite(rw http.ResponseWriter, req *http.Request, enc encoding.Encoder, params martini.Params, dtx *apicontext.DataContext) string {
	var err error
	var ag applicationGuide.ApplicationGuide
	id := params["id"]
	ag.Website.ID, err = strconv.Atoi(id)

	ags, err := ag.GetBySite(dtx)
	if err != nil {
		apierror.GenerateError("Error getting Application Guides", err, rw, req)
	}
	return encoding.Must(enc.Encode(ags))
}

func CreateApplicationGuide(rw http.ResponseWriter, req *http.Request, enc encoding.Encoder, params martini.Params, dtx *apicontext.DataContext) string {
	contType := req.Header.Get("Content-Type")

	var ag applicationGuide.ApplicationGuide
	var err error

	// if contType == "application/json" {
	if strings.Contains(contType, "application/json") {
		//json
		requestBody, err := ioutil.ReadAll(req.Body)
		if err != nil {
			apierror.GenerateError("Error reading request body", err, rw, req)
		}

		err = json.Unmarshal(requestBody, &ag)
		if err != nil {
			apierror.GenerateError("Error decoding request body", err, rw, req)
		}
	} else {
		//else, form
		ag.Url = req.FormValue("url")
		web := req.FormValue("website_id")
		ag.FileType = req.FormValue("file_type")
		cat := req.FormValue("category_id")

		if err != nil {
			apierror.GenerateError("Error parsing form", err, rw, req)
		}
		if web != "" {
			ag.Website.ID, err = strconv.Atoi(web)
		}
		if cat != "" {
			ag.Category.ID, err = strconv.Atoi(cat)
		}
		if err != nil {
			apierror.GenerateError("Error parsing category ID or website ID", err, rw, req)
		}
	}
	err = ag.Create(dtx)
	if err != nil {
		log.Print("HERe", err)
		apierror.GenerateError("Error creating Application Guide", err, rw, req)
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
		apierror.GenerateError("Error deleting Application Guide", err, rw, req)
	}

	//Return JSON
	return encoding.Must(enc.Encode(ag))
}
