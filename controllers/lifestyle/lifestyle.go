package lifestyle

import (
	"github.com/curt-labs/API/helpers/apicontext"
	"github.com/curt-labs/API/helpers/encoding"
	"github.com/curt-labs/API/helpers/error"
	"github.com/curt-labs/API/models/lifestyle"
	"github.com/go-martini/martini"

	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
)

func GetAll(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params, dtx *apicontext.DataContext) string {
	lifestyles, err := lifestyle.GetAll(dtx)
	if err != nil {
		apierror.GenerateError("Trouble getting all lifestyles", err, w, r)
		return ""
	}

	return encoding.Must(enc.Encode(lifestyles))
}

func Get(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params, dtx *apicontext.DataContext) string {
	var l lifestyle.Lifestyle
	var err error
	l.ID, err = strconv.Atoi(params["id"])
	if err != nil {
		apierror.GenerateError("Trouble getting lifestyle ID", err, w, r)
		return ""
	}

	err = l.Get(dtx)
	if err != nil {
		apierror.GenerateError("Trouble getting lifestyle", err, w, r)
		return ""
	}

	return encoding.Must(enc.Encode(l))
}

func Save(rw http.ResponseWriter, req *http.Request, enc encoding.Encoder, params martini.Params, dtx *apicontext.DataContext) string {
	var l lifestyle.Lifestyle
	var err error
	idStr := params["id"]
	if idStr != "" {
		l.ID, err = strconv.Atoi(idStr)
		err = l.Get(dtx)
		if err != nil {
			apierror.GenerateError("Trouble getting lifestyle", err, rw, req)
			return ""
		}
	}

	//json
	requestBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		apierror.GenerateError("Trouble reading request body for saving lifestyle", err, rw, req)
		return ""
	}
	err = json.Unmarshal(requestBody, &l)
	if err != nil {
		apierror.GenerateError("Trouble unmarshalling request body response for saving lifestyle", err, rw, req)
		return ""
	}

	//create or update
	if l.ID > 0 {
		err = l.Update(dtx)
	} else {
		err = l.Create(dtx)
	}

	if err != nil {
		apierror.GenerateError("Trouble saving lifestyle", err, rw, req)
		return ""
	}
	return encoding.Must(enc.Encode(l))
}

func Delete(rw http.ResponseWriter, req *http.Request, enc encoding.Encoder, params martini.Params, dtx *apicontext.DataContext) string {
	var l lifestyle.Lifestyle
	var err error
	idStr := params["id"]
	if idStr != "" {
		l.ID, err = strconv.Atoi(idStr)
		err = l.Get(dtx)
		if err != nil {
			apierror.GenerateError("Trouble getting lifestyle", err, rw, req)
			return ""
		}
	}
	err = l.Delete()
	if err != nil {
		apierror.GenerateError("Trouble deleting lifestyle", err, rw, req)
		return ""
	}
	return encoding.Must(enc.Encode(l))
}
