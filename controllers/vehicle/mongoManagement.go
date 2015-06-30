package vehicle

import (
	"github.com/curt-labs/GoAPI/helpers/apicontext"
	"github.com/curt-labs/GoAPI/helpers/encoding"
	"github.com/curt-labs/GoAPI/helpers/error"
	"github.com/curt-labs/GoAPI/models/products"
	"github.com/go-martini/martini"

	"encoding/json"
	"io/ioutil"
	"net/http"
)

func GetAllCollectionApplications(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, dtx *apicontext.DataContext, params martini.Params) string {
	collection := params["collection"]
	if collection == "" {
		apierror.GenerateError("No Collection in URL", nil, w, r)
		return ""
	}
	apps, err := products.GetAllCollectionApplications(collection)
	if err != nil {
		apierror.GenerateError(err.Error(), err, w, r)
		return ""
	}
	return encoding.Must(enc.Encode(apps))
}

func UpdateApplication(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, dtx *apicontext.DataContext, params martini.Params) string {
	var app products.NoSqlVehicle
	collection := params["collection"]
	if collection == "" {
		apierror.GenerateError("No Collection in URL", nil, w, r)
		return ""
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		apierror.GenerateError("Error reading request body", nil, w, r)
		return ""
	}

	if err = json.Unmarshal(body, &app); err != nil {
		apierror.GenerateError("Error decoding vehicle application", nil, w, r)
		return ""
	}

	if err = app.Update(collection); err != nil {
		apierror.GenerateError("Error updating vehicle", nil, w, r)
		return ""
	}
	return encoding.Must(enc.Encode(app))
}

func DeleteApplication(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, dtx *apicontext.DataContext, params martini.Params) string {
	var app products.NoSqlVehicle
	collection := params["collection"]
	if collection == "" {
		apierror.GenerateError("No Collection in URL", nil, w, r)
		return ""
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		apierror.GenerateError("Error reading request body", nil, w, r)
		return ""
	}

	if err = json.Unmarshal(body, &app); err != nil {
		apierror.GenerateError("Error decoding vehicle application", nil, w, r)
		return ""
	}

	if err = app.Delete(collection); err != nil {
		apierror.GenerateError("Error updating vehicle", nil, w, r)
		return ""
	}
	return encoding.Must(enc.Encode(app))
}

func CreateApplication(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, dtx *apicontext.DataContext, params martini.Params) string {
	var app products.NoSqlVehicle
	collection := params["collection"]
	if collection == "" {
		apierror.GenerateError("No Collection in URL", nil, w, r)
		return ""
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		apierror.GenerateError("Error reading request body", nil, w, r)
		return ""
	}

	if err = json.Unmarshal(body, &app); err != nil {
		apierror.GenerateError("Error decoding vehicle application", nil, w, r)
		return ""
	}

	if err = app.Create(collection); err != nil {
		apierror.GenerateError("Error updating vehicle", nil, w, r)
		return ""
	}
	return encoding.Must(enc.Encode(app))
}
