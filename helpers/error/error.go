package apierror

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

type ApiErr struct {
	Message     string     `json:"message" xml:"message"`
	RequestBody string     `json:"request_body" xml:"request_body"`
	QueryString url.Values `json:"query_string" xml:"query_string"`
}

func GenerateError(msg string, err error, res http.ResponseWriter, r *http.Request) {
	e := ApiErr{
		Message: "",
	}
	if msg != "" {
		e.Message = msg
	} else if err != nil {
		e.Message = err.Error()
	}

	if r != nil && r.Body != nil {
		defer r.Body.Close()

		data, readErr := ioutil.ReadAll(r.Body)
		if readErr == nil {
			e.RequestBody = string(data)
		}
		e.QueryString = r.URL.Query()
	}

	js, jsErr := json.Marshal(e)
	if jsErr != nil {
		http.Error(res, e.Message, http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusInternalServerError)
	res.Write(js)
	return
}
