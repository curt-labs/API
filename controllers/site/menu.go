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

func GetMenu(rw http.ResponseWriter, req *http.Request, enc encoding.Encoder, params martini.Params, dtx *apicontext.DataContext) string {
	var m site.Menu
	var err error
	idStr := params["id"]
	id, err := strconv.Atoi(idStr)

	if err == nil {
		//Thar be an Id int
		m.Id = id
		err = m.Get(dtx)
		if err != nil {
			apierror.GenerateError("Trouble getting site menu", err, rw, req)
		}
	} else {
		//Thar be a name
		m.Name = idStr
		err = m.GetByName(dtx)
		if err != nil {
			apierror.GenerateError("Trouble getting site menu", err, rw, req)
		}
	}
	return encoding.Must(enc.Encode(m))
}

func GetAllMenus(rw http.ResponseWriter, req *http.Request, enc encoding.Encoder, params martini.Params, dtx *apicontext.DataContext) string {
	m, err := site.GetAllMenus(dtx)
	if err != nil {
		apierror.GenerateError("Trouble getting all site menus", err, rw, req)
	}
	return encoding.Must(enc.Encode(m))
}

func GetMenuWithContents(rw http.ResponseWriter, req *http.Request, enc encoding.Encoder, params martini.Params, dtx *apicontext.DataContext) string {
	var m site.Menu
	var err error
	idStr := params["id"]
	id, err := strconv.Atoi(idStr)

	if err == nil {
		//Thar be an Id int
		m.Id = id
		err = m.Get(dtx)

	} else {
		//Thar be a name
		m.Name = idStr
		err = m.GetByName(dtx)
	}

	if err != nil {
		apierror.GenerateError("Trouble getting site menu", err, rw, req)
		return ""
	}

	err = m.GetContents()
	if err != nil {
		apierror.GenerateError("Trouble getting site menu with contents", err, rw, req)
		return ""
	}

	return encoding.Must(enc.Encode(m))
}

func SaveMenu(rw http.ResponseWriter, req *http.Request, enc encoding.Encoder, params martini.Params, dtx *apicontext.DataContext) string {
	var m site.Menu
	var err error
	idStr := params["id"]
	if idStr != "" {
		m.Id, err = strconv.Atoi(idStr)
		err = m.Get(dtx)
		if err != nil {
			apierror.GenerateError("Trouble getting site menu ID", err, rw, req)
		}
	}

	//json
	requestBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		apierror.GenerateError("Trouble reading request body for saving site menu", err, rw, req)
	}
	err = json.Unmarshal(requestBody, &m)
	if err != nil {
		apierror.GenerateError("Trouble unmarshalling request body for saving site menu", err, rw, req)
	}
	//create or update
	if m.Id > 0 {
		err = m.Update()
	} else {
		err = m.Create()
	}

	if err != nil {
		apierror.GenerateError("Trouble saving site menu", err, rw, req)
	}
	return encoding.Must(enc.Encode(m))
}

func DeleteMenu(rw http.ResponseWriter, req *http.Request, enc encoding.Encoder, params martini.Params, dtx *apicontext.DataContext) string {
	var err error
	var m site.Menu

	idStr := params["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		apierror.GenerateError("Trouble getting site menu ID", err, rw, req)
	}
	m.Id = id
	err = m.Delete()
	if err != nil {
		apierror.GenerateError("Trouble deleting site menu", err, rw, req)
	}

	return encoding.Must(enc.Encode(m))
}
