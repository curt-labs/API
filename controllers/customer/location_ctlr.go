package customer_ctlr

import (
	"encoding/json"
	"github.com/curt-labs/GoAPI/helpers/apicontext"
	"github.com/curt-labs/GoAPI/helpers/encoding"
	"github.com/curt-labs/GoAPI/helpers/sortutil"
	"github.com/curt-labs/GoAPI/models/customer"
	"github.com/go-martini/martini"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

const (
	timeFormat = "2006-01-02 15:04:05"
)

func GetLocation(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	var c customer.CustomerLocation
	var err error

	c.Id, err = strconv.Atoi(params["id"])
	if err != nil {
		return err.Error()
	}

	err = c.Get()
	if err != nil {
		return err.Error()
	}

	return encoding.Must(enc.Encode(c))
}

func GetAllLocations(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder, dtx *apicontext.DataContext) string {
	var c customer.CustomerLocations
	var err error
	apiKey := dtx.APIKey
	brandID := dtx.BrandID
	c, err = customer.GetAllLocations(apiKey, brandID)
	if err != nil {
		return err.Error()
	}
	sort := r.FormValue("sort")
	direction := r.FormValue("direction")
	if sort != "" {
		if strings.ContainsAny(direction, "esc") {
			sortutil.CiDescByField(c, sort)
		} else {
			sortutil.CiAscByField(c, sort)
		}
	}
	return encoding.Must(enc.Encode(c))
}

func SaveLocation(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	var w customer.CustomerLocation
	var err error
	err = r.ParseForm()

	id := params["id"]
	if id == "" {
		id = r.Form.Get("id")
		if id != "" {
			w.Id, err = strconv.Atoi(id)
			w.Get()
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
		w.Name = name
	}
	if address != "" {
		w.Address = address
	}
	if city != "" {
		w.City = city
	}
	if state != "" {
		w.State.Id, err = strconv.Atoi(state)
	}
	if email != "" {
		w.Email = email
	}
	if phone != "" {
		w.Phone = phone
	}
	if fax != "" {
		w.Fax = fax
	}
	if latitude != "" {
		w.Latitude, err = strconv.ParseFloat(latitude, 64)
	}
	if longitude != "" {
		w.Longitude, err = strconv.ParseFloat(longitude, 64)
	}
	if customerID != "" {
		w.CustomerId, err = strconv.Atoi(customerID)
	}
	if contactPerson != "" {
		w.ContactPerson = contactPerson
	}
	if isPrimary != "" {
		w.IsPrimary, err = strconv.ParseBool(isPrimary)
	}
	if postalCode != "" {
		w.PostalCode = postalCode
	}
	if shippingDefault != "" {
		w.ShippingDefault, err = strconv.ParseBool(shippingDefault)
	} else {
		w.ShippingDefault = false
	}

	if id != "" {
		err = w.Update()
	} else {
		err = w.Create()
	}

	if err != nil {
		return err.Error()
	}
	return encoding.Must(enc.Encode(w))
}

func SaveLocationJson(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	var l customer.CustomerLocation
	var err error
	id := params["id"]
	if id != "" {
		l.Id, err = strconv.Atoi(id)
		err = l.Get()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return ""
		}
	}

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}
	err = json.Unmarshal(body, &l)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}
	if l.Id != 0 {
		err = l.Update()
	} else {
		err = l.Create()
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}
	return encoding.Must(enc.Encode(l))
}

func DeleteLocation(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	var w customer.CustomerLocation
	var err error

	id := params["id"]
	if id != "" {
		w.Id, err = strconv.Atoi(id)
		if err != nil {
			return err.Error()
		}
		w.Delete()
	}
	return encoding.Must(enc.Encode(w))
}
