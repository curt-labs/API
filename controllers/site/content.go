package site

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/curt-labs/GoAPI/helpers/apicontext"
	"github.com/curt-labs/GoAPI/helpers/encoding"
	"github.com/curt-labs/GoAPI/helpers/error"
	"github.com/curt-labs/GoAPI/models/site"
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
		}
	} else {
		//Thar be a slug
		c.Slug = idStr
		err = c.GetBySlug()
		if err != nil {
			apierror.GenerateError("Trouble getting site content by slug.", err, rw, req)
		}
	}
	return encoding.Must(enc.Encode(c))
}

func GetAllContents(rw http.ResponseWriter, req *http.Request, enc encoding.Encoder, params martini.Params, dtx *apicontext.DataContext) string {
	m, err := site.GetAllContents()
	if err != nil {
		apierror.GenerateError("Trouble getting all site contents", err, rw, req)
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
