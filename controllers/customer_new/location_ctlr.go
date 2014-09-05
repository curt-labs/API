package customer_ctlr_new

import (
	"github.com/curt-labs/GoAPI/helpers/encoding"
	"github.com/curt-labs/GoAPI/helpers/sortutil"
	"github.com/curt-labs/GoAPI/models/customer_new"
	"github.com/go-martini/martini"
	"net/http"
	"strconv"
	"strings"
)

const (
	timeFormat = "2006-01-02 15:04:05"
)

func GetLocation(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	var c customer_new.CustomerLocation
	var err error
	c.Id, err = strconv.Atoi(r.FormValue("id"))
	if err != nil {
		return err.Error()
	}

	err = c.Get()
	if err != nil {
		return err.Error()
	}
	return encoding.Must(enc.Encode(c))
}
func GetAllLocations(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder) string {
	var c customer_new.CustomerLocations
	var err error

	c, err = customer_new.GetAllLocations()
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

func CreateUpdateLocation(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder) string {
	var w customer_new.CustomerLocation
	var err error

	id := r.FormValue("id")
	if id != "" {
		w.Id, err = strconv.Atoi(id)
		w.Get()
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

func DeleteLocation(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder) string {
	var w customer_new.CustomerLocation
	var err error

	id := r.FormValue("id")
	if id != "" {
		w.Id, err = strconv.Atoi(id)
		if err != nil {
			return err.Error()
		}
		w.Delete()
	}
	return encoding.Must(enc.Encode(w))
}
