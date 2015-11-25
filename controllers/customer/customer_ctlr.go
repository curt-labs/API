package customer_ctlr

import (
	"github.com/curt-labs/API/helpers/apicontext"
	"github.com/curt-labs/API/helpers/encoding"
	"github.com/curt-labs/API/helpers/error"
	"github.com/curt-labs/API/models/customer"
	"github.com/curt-labs/API/models/products"
	"github.com/go-martini/martini"

	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func GetCustomer(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params, dtx *apicontext.DataContext) string {
	var err error
	var c customer.Customer

	if err = c.GetCustomerIdFromKey(dtx.APIKey); err != nil {
		apierror.GenerateError("Trouble getting customer ID", err, rw, r)
		return ""
	}

	if err = c.GetCustomer(dtx.APIKey); err != nil {
		apierror.GenerateError("Trouble getting customer", err, rw, r, http.StatusServiceUnavailable)
		return ""
	}

	lowerKey := strings.ToLower(dtx.APIKey)
	for i, u := range c.Users {
		for _, k := range u.Keys {
			if strings.ToLower(k.Key) == lowerKey {
				c.Users[i].Current = true
			}
		}
	}

	return encoding.Must(enc.Encode(c))
}

func GetLocations(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params, dtx *apicontext.DataContext) string {
	cu, err := customer.GetCustomerUserFromKey(dtx.APIKey)
	if err != nil {
		err = errors.New("Unauthorized!")
		apierror.GenerateError("Unauthorized!", err, rw, r, http.StatusUnauthorized)
		return ""
	}
	c, err := cu.GetCustomer(dtx.APIKey)
	if err != nil {
		err = errors.New("Unauthorized!")
		apierror.GenerateError("Unauthorized!", err, rw, r, http.StatusUnauthorized)
		return ""
	}
	return encoding.Must(enc.Encode(c.Locations))
}

//TODO - redundant
func GetUsers(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params, dtx *apicontext.DataContext) string {
	user, err := customer.GetCustomerUserFromKey(dtx.APIKey)
	if err != nil {
		apierror.GenerateError("Trouble getting customer user", err, rw, r)
		return ""
	}

	if !user.Sudo {
		err = errors.New("Unauthorized!")
		apierror.GenerateError("Unauthorized!", err, rw, r, http.StatusUnauthorized)
		return ""
	}

	cust, err := user.GetCustomer(dtx.APIKey)
	if err != nil {
		apierror.GenerateError("Trouble getting customer", err, rw, r)
		return ""
	}

	if err = cust.GetUsers(dtx.APIKey); err != nil {
		apierror.GenerateError("Trouble getting users", err, rw, r)
		return ""
	}

	return encoding.Must(enc.Encode(cust.Users))
}

//Todo redundant
//Hacky like this to work with old forms of authentication
func GetUser(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder) string {
	qs, err := url.Parse(r.URL.String())
	if err != nil {
		apierror.GenerateError("err parsing url", err, rw, r)
		return ""
	}

	key := qs.Query().Get("key")
	if key == "" {
		key = r.FormValue("key")
	}

	user, err := customer.GetCustomerUserFromKey(key)
	if err != nil {
		apierror.GenerateError("Trouble getting customer user", err, rw, r)
		return ""
	}

	return encoding.Must(enc.Encode(user))
}

func GetCustomerPrice(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params, dtx *apicontext.DataContext) string {
	var err error
	var p products.Part

	id := r.FormValue("id")
	if id == "" {
		id = params["id"]
	}

	if p.ID, err = strconv.Atoi(id); err != nil {
		apierror.GenerateError("Trouble getting part ID", err, rw, r)
		return ""
	}

	var price float64
	if price, err = customer.GetCustomerPrice(dtx, p.ID); err != nil {
		apierror.GenerateError("Trouble getting price", err, rw, r)
		return ""
	}

	return encoding.Must(enc.Encode(price))
}

func GetCustomerCartReference(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params, dtx *apicontext.DataContext) string {
	var err error
	var p products.Part

	id := r.FormValue("id")
	if id == "" {
		id = params["id"]
	}

	if p.ID, err = strconv.Atoi(id); err != nil {
		apierror.GenerateError("Trouble getting part ID", err, rw, r)
		return ""
	}

	var ref int
	if ref, err = customer.GetCustomerCartReference(dtx.APIKey, p.ID); err != nil {
		apierror.GenerateError("Trouble getting customer cart reference", err, rw, r)
		return ""
	}

	return encoding.Must(enc.Encode(ref))
}

func SaveCustomer(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params, dtx *apicontext.DataContext) string {
	var c customer.Customer
	var err error

	if r.FormValue("id") != "" || params["id"] != "" {
		id := r.FormValue("id")
		if id == "" {
			id = params["id"]
		}

		if c.Id, err = strconv.Atoi(id); err != nil {
			apierror.GenerateError("Trouble getting customer ID", err, rw, r)
			return ""
		}

		if err = c.Basics(dtx.APIKey); err != nil {
			apierror.GenerateError("Trouble getting customer", err, rw, r)
			return ""
		}
	}

	//json
	var requestBody []byte
	if requestBody, err = ioutil.ReadAll(r.Body); err != nil {
		apierror.GenerateError("Trouble reading request body while saving customer", err, rw, r)
		return ""
	}

	if err = json.Unmarshal(requestBody, &c); err != nil {
		apierror.GenerateError("Trouble unmarshalling json request body while saving customer", err, rw, r)
		return ""
	}

	//create or update
	if c.Id > 0 {
		err = c.Update()
	} else {
		err = c.Create()
	}

	if err != nil {
		msg := "Trouble creating customer"
		if c.Id > 0 {
			msg = "Trouble updating customer"
		}
		apierror.GenerateError(msg, err, rw, r)
		return ""
	}

	return encoding.Must(enc.Encode(c))
}

func DeleteCustomer(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	var c customer.Customer
	var err error

	id := r.FormValue("id")
	if id == "" {
		id = params["id"]
	}

	if c.Id, err = strconv.Atoi(id); err != nil {
		apierror.GenerateError("Trouble getting customer ID", err, rw, r)
		return ""
	}

	if err = c.Delete(); err != nil {
		apierror.GenerateError("Trouble deleting customer", err, rw, r)
		return ""
	}

	return encoding.Must(enc.Encode(c))
}
