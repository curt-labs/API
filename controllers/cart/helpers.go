package cart_ctlr

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

type CartErr struct {
	Message     string     `json:"message" xml:"message"`
	Error       error      `json:"error" xml:"error"`
	RequestBody string     `json:"request_body" xml:"request_body"`
	QueryString url.Values `json:"query_string" xml:"query_string"`
}

func generateError(msg string, err error, res http.ResponseWriter, r *http.Request) {
	var e CartErr
	if msg != "" {
		e.Message = msg
	} else if err != nil {
		e.Message = err.Error()
	}
	defer r.Body.Close()

	e.Error = err
	data, readErr := ioutil.ReadAll(r.Body)
	if readErr == nil {
		e.RequestBody = string(data)
	}
	e.QueryString = r.URL.Query()

	js, jsErr := json.Marshal(e)
	if jsErr != nil {
		http.Error(res, e.Message, http.StatusInternalServerError)
		return
	}

	res.Write(js)
	res.Header().Set("StatusCode", "500")
	return
}
