package apierror

import (
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/curt-labs/API/helpers/database"
	"gopkg.in/mgo.v2"
)

type ApiErr struct {
	Message        string     `json:"message" xml:"message"`
	MessageDetails string     `json:"messageDetails" xml:"message_details"`
	RequestBody    string     `json:"request_body" xml:"request_body"`
	QueryString    url.Values `json:"query_string" xml:"query_string"`
}

func GenerateError(msg string, err error, res http.ResponseWriter, r *http.Request, errorCode ...int) {
	e := ApiErr{
		Message: "",
	}

	e.Message = msg
	if err != nil {
		if e.Message == "" {
			e.Message = err.Error()
		}
		e.MessageDetails = err.Error()
	}

	if r != nil && r.Body != nil {
		defer r.Body.Close()

		data, readErr := ioutil.ReadAll(r.Body)
		if readErr == nil {
			e.RequestBody = string(data)
		}
		e.QueryString = r.URL.Query()
	}

	//mongo
	logError(e)

	var errorResp []byte
	var marshalErr error

	res.Header().Set("Access-Control-Allow-Origin", "*")
	switch r.Header.Get("Content-Type") {
	case "application/xml":
		res.Header().Set("Content-Type", "application/xml")
		errorResp, marshalErr = xml.Marshal(e)
	case "application/json":
		//JSON is our defaulted content type encoding
		res.Header().Set("Content-Type", "application/json")
		errorResp, marshalErr = json.Marshal(e)
	default:
		//JSON is our defaulted content type encoding
		res.Header().Set("Content-Type", "application/json")
		errorResp, marshalErr = json.Marshal(e)
	}

	if marshalErr != nil {
		http.Error(res, e.Message, http.StatusInternalServerError)
		return
	}

	respCode := http.StatusInternalServerError
	if len(errorCode) > 0 {
		respCode = errorCode[0]
	}

	res.WriteHeader(respCode)
	res.Write(errorResp)
	return
}

func logError(e ApiErr) error {
	session, err := mgo.DialWithInfo(database.MongoConnectionString())
	if err != nil {
		return err
	}
	defer session.Close()
	c := session.DB("errorDB").C("errors")
	err = c.Insert(e)
	return err
}
