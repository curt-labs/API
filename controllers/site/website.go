package site

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/curt-labs/API/helpers/apicontext"
	"github.com/curt-labs/API/helpers/encoding"
	"github.com/curt-labs/API/helpers/error"
	"github.com/curt-labs/API/models/site"
	"github.com/go-martini/martini"
)

func GetSiteDetails(rw http.ResponseWriter, req *http.Request, enc encoding.Encoder, params martini.Params, dtx *apicontext.DataContext) string {
	var w site.Website
	var err error
	idStr := params["id"]
	id, err := strconv.Atoi(idStr)

	if err != nil {
		apierror.GenerateError("Trouble getting site ID", err, rw, req)
	}
	w.ID = id

	err = w.GetDetails(dtx)
	if err != nil {
		apierror.GenerateError("Trouble getting site details", err, rw, req)
	}
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
			apierror.GenerateError("Trouble getting website", err, rw, req)
		}
	}

	//json
	requestBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		apierror.GenerateError("Trouble reading request body for saving website", err, rw, req)
	}
	err = json.Unmarshal(requestBody, &m)
	if err != nil {
		apierror.GenerateError("Trouble unmarshalling request body for saving website", err, rw, req)
	}
	//create or update
	if m.ID > 0 {
		err = m.Update()
	} else {
		err = m.Create()
	}

	if err != nil {
		apierror.GenerateError("Trouble saving website", err, rw, req)
	}
	return encoding.Must(enc.Encode(m))
}

func DeleteSite(rw http.ResponseWriter, req *http.Request, enc encoding.Encoder, params martini.Params) string {
	var err error
	var m site.Website

	idStr := params["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		apierror.GenerateError("Trouble getting website ID", err, rw, req)
	}
	m.ID = id
	err = m.Delete()
	if err != nil {
		apierror.GenerateError("Trouble getting website", err, rw, req)
	}

	return encoding.Must(enc.Encode(m))
}
