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

func GetContent(rw http.ResponseWriter, req *http.Request, enc encoding.Encoder, params martini.Params, dtx *apicontext.DataContext) string {
	var c site.Content
	var err error
	idStr := params["id"]
	id, err := strconv.Atoi(idStr)

	if err == nil {
		//Thar be an Id int
		c.Id = id
		err = c.Get(dtx)
		if err != nil {
			apierror.GenerateError("Trouble getting site content by Id.", err, rw, req)
			return ""
		}
	}

	//Thar be a slug
	c.Slug = idStr
	if c.Slug != "_ah/health" {
		err = c.GetBySlug(dtx)
		if err != nil {
			apierror.GenerateError("Trouble getting site content by slug.", err, rw, req, 204)
			return ""
		}
	}
	return encoding.Must(enc.Encode(c))
}

func GetAllContents(rw http.ResponseWriter, req *http.Request, enc encoding.Encoder, params martini.Params, dtx *apicontext.DataContext) string {
	var siteID int
	var err error
	siteID_str := req.URL.Query().Get("siteID")
	if siteID_str != "" {
		siteID, err = strconv.Atoi(siteID_str)
		if err != nil {
			apierror.GenerateError("Trouble getting all site contents", err, rw, req)
			return ""
		}

	}
	m, err := site.GetAllContents(dtx, siteID)
	if err != nil {
		apierror.GenerateError("Trouble getting all site contents", err, rw, req)
		return ""
	}
	return encoding.Must(enc.Encode(m))
}

func GetContentRevisions(rw http.ResponseWriter, req *http.Request, enc encoding.Encoder, params martini.Params, dtx *apicontext.DataContext) string {
	var c site.Content
	var err error
	idStr := params["id"]
	c.Id, err = strconv.Atoi(idStr)

	err = c.GetContentRevisions()
	if err != nil {
		apierror.GenerateError("Trouble getting site content revisions", err, rw, req)
	}

	return encoding.Must(enc.Encode(c))
}

func SaveContent(rw http.ResponseWriter, req *http.Request, enc encoding.Encoder, params martini.Params, dtx *apicontext.DataContext) string {
	var c site.Content
	var err error
	idStr := params["id"]
	if idStr != "" {
		c.Id, err = strconv.Atoi(idStr)
		err = c.Get(dtx)
		if err != nil {
			apierror.GenerateError("Trouble getting site content", err, rw, req)
		}
	}

	//json
	requestBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		apierror.GenerateError("Trouble reading request body for saving site content", err, rw, req)
	}
	err = json.Unmarshal(requestBody, &c)
	if err != nil {
		apierror.GenerateError("Trouble unmarshalling request body for saving site content", err, rw, req)
	}
	//create or update
	if c.Id > 0 {
		err = c.Update()
	} else {
		err = c.Create()
	}

	if err != nil {
		apierror.GenerateError("Trouble saving site content", err, rw, req)
	}
	return encoding.Must(enc.Encode(c))
}

func DeleteContent(rw http.ResponseWriter, req *http.Request, enc encoding.Encoder, params martini.Params, dtx *apicontext.DataContext) string {
	var err error
	var c site.Content

	idStr := params["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		apierror.GenerateError("Trouble getting site content ID", err, rw, req)
	}
	c.Id = id
	err = c.Delete()
	if err != nil {
		apierror.GenerateError("Trouble deleting site content", err, rw, req)
	}

	return encoding.Must(enc.Encode(c))
}
