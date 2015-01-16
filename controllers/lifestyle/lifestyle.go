package lifestyle

import (
	"encoding/json"
	"github.com/curt-labs/GoAPI/helpers/apicontext"
	"github.com/curt-labs/GoAPI/helpers/encoding"
	"github.com/curt-labs/GoAPI/models/lifestyle"
	"github.com/go-martini/martini"
	"io/ioutil"
	"net/http"
	"strconv"
)

func GetAll(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params, dtx *apicontext.DataContext) string {
	lifestyles, err := lifestyle.GetAll(dtx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	return encoding.Must(enc.Encode(lifestyles))
}
func Get(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params, dtx *apicontext.DataContext) string {
	var l lifestyle.Lifestyle
	var err error
	l.ID, err = strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	err = l.Get()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
		err = l.Get()
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return ""
		}
	}

	//json
	requestBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return ""
	}
	err = json.Unmarshal(requestBody, &l)
	if err != nil {

		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return ""
	}

	//create or update
	if l.ID > 0 {
		err = l.Update()
	} else {
		err = l.Create()
	}

	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
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
		err = l.Get()
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return ""
		}
	}
	err = l.Delete()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return ""
	}
	return encoding.Must(enc.Encode(l))

}
