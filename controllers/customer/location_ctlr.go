package customer_ctlr

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/curt-labs/GoAPI/helpers/apicontext"
	"github.com/curt-labs/GoAPI/helpers/encoding"
	"github.com/curt-labs/GoAPI/helpers/error"
	"github.com/curt-labs/GoAPI/models/customer"
	"github.com/go-martini/martini"
)

const (
	timeFormat = "2006-01-02 15:04:05"
)

func GetLocation(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	var cl customer.CustomerLocation
	var err error

	id := r.FormValue("id")
	if id == "" {
		id = params["id"]
	}

	if cl.Id, err = strconv.Atoi(id); err != nil {
		apierror.GenerateError("Trouble getting location ID", err, rw, r)
		return ""
	}

	if err = cl.Get(); err != nil {
		apierror.GenerateError("Trouble getting location", err, rw, r)
		return ""
	}

	return encoding.Must(enc.Encode(cl))
}

func SaveLocation(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params, dtx *apicontext.DataContext) string {
	var cl customer.CustomerLocation
	var err error

	if r.FormValue("id") != "" || params["id"] != "" {
		id := r.FormValue("id")
		if id == "" {
			id = params["id"]
		}

		if cl.Id, err = strconv.Atoi(id); err != nil {
			apierror.GenerateError("Trouble getting location ID", err, rw, r)
			return ""
		}

		if err = cl.Get(); err != nil {
			apierror.GenerateError("Trouble getting location", err, rw, r)
			return ""
		}
	}

	name := r.FormValue("name")
	address := r.FormValue("address")
	city := r.FormValue("city")
	state := r.FormValue("stateId")
	email := r.FormValue("email")
	phone := r.FormValue("phone")
	fax := r.FormValue("fax")
	latitude := r.FormValue("latitude")
	longitude := r.FormValue("longitude")
	customerID := r.FormValue("customerId")
	contactPerson := r.FormValue("contactPerson")
	isPrimary := r.FormValue("isPrimary")
	postalCode := r.FormValue("postalCode")
	shippingDefault := r.FormValue("shippingDefault")

	if name != "" {
		cl.Name = name
	}

	if address != "" {
		cl.Address = address
	}

	if city != "" {
		cl.City = city
	}

	if state != "" {
		if cl.State.Id, err = strconv.Atoi(state); err != nil {
			apierror.GenerateError("Trouble setting state ID", err, rw, r)
			return ""
		}
	}

	if email != "" {
		cl.Email = email
	}

	if phone != "" {
		cl.Phone = phone
	}

	if fax != "" {
		cl.Fax = fax
	}

	if latitude != "" {
		if cl.Coordinates.Latitude, err = strconv.ParseFloat(latitude, 64); err != nil {
			cl.Coordinates.Latitude = 0
		}
	}

	if longitude != "" {
		if cl.Coordinates.Longitude, err = strconv.ParseFloat(longitude, 64); err != nil {
			cl.Coordinates.Longitude = 0
		}
	}

	if customerID != "" {
		if cl.CustomerId, err = strconv.Atoi(customerID); err != nil {
			apierror.GenerateError("Trouble getting customer ID", err, rw, r)
			return ""
		}
	}

	if contactPerson != "" {
		cl.ContactPerson = contactPerson
	}

	if isPrimary != "" {
		if cl.IsPrimary, err = strconv.ParseBool(isPrimary); err != nil {
			cl.IsPrimary = false
		}
	}

	if postalCode != "" {
		cl.PostalCode = postalCode
	}

	if shippingDefault != "" {
		if cl.ShippingDefault, err = strconv.ParseBool(shippingDefault); err != nil {
			cl.ShippingDefault = false
		}
	}

	if cl.Id > 0 {
		err = cl.Update(dtx)
	} else {
		err = cl.Create(dtx)
	}

	if err != nil {
		msg := "Trouble creating location"
		if cl.Id > 0 {
			msg = "Trouble updating location"
		}
		apierror.GenerateError(msg, err, rw, r)
		return ""
	}

	return encoding.Must(enc.Encode(cl))
}

func SaveLocationJson(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params, dtx *apicontext.DataContext) string {
	var cl customer.CustomerLocation
	var err error

	id := r.FormValue("id")
	if id == "" {
		id = params["id"]
	}

	if cl.Id, err = strconv.Atoi(id); err != nil {
		apierror.GenerateError("Trouble getting location ID", err, rw, r)
		return ""
	}

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		apierror.GenerateError("Trouble reading request body while saving location", err, rw, r)
		return ""
	}

	if err = json.Unmarshal(body, &cl); err != nil {
		apierror.GenerateError("Trouble unmarshalling json request body while saving location", err, rw, r)
		return ""
	}

	if cl.Id > 0 {
		err = cl.Update(dtx)
	} else {
		err = cl.Create(dtx)
	}

	if err != nil {
		msg := "Trouble creating location"
		if cl.Id > 0 {
			msg = "Trouble updating location"
		}
		apierror.GenerateError(msg, err, rw, r)
		return ""
	}

	return encoding.Must(enc.Encode(cl))
}

func DeleteLocation(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params, dtx *apicontext.DataContext) string {
	var cl customer.CustomerLocation
	var err error

	id := r.FormValue("id")
	if id == "" {
		id = params["id"]
	}

	if cl.Id, err = strconv.Atoi(id); err != nil {
		apierror.GenerateError("Trouble getting location ID", err, rw, r)
		return ""
	}

	if err = cl.Delete(dtx); err != nil {
		apierror.GenerateError("Trouble deleting location", err, rw, r)
		return ""
	}

	return encoding.Must(enc.Encode(cl))
}
