package customer_ctlr

import (
	"github.com/curt-labs/GoAPI/helpers/apicontext"
	"github.com/curt-labs/GoAPI/helpers/encoding"
	apierr "github.com/curt-labs/GoAPI/helpers/error"
	"github.com/curt-labs/GoAPI/models/customer"
	"github.com/curt-labs/GoAPI/models/products"
	"github.com/go-martini/martini"

	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

func GetCustomer(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params, dtx *apicontext.DataContext) string {
	var err error
	var c customer.Customer
	idStr := params["id"]
	if idStr != "" {
		c.Id, err = strconv.Atoi(idStr)
	} else {
		apierr.GenerateError("No Id Supplied.", errors.New("No results."), w, r)
	}

	err = c.GetCustomer(dtx.APIKey)
	if err != nil {
		http.Error(w, "Error getting customer.", http.StatusServiceUnavailable)
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

func GetLocations(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params, dtx *apicontext.DataContext) string {
	cu, err := customer.GetCustomerUserFromKey(dtx.APIKey)
	if err != nil {
		http.Error(w, "Unauthorized!", http.StatusUnauthorized)
		return ""
	}
	c, err := cu.GetCustomer(dtx.APIKey)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return ""
	}
	return encoding.Must(enc.Encode(c.Locations))
}

//TODO - redundant
func GetUsers(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params, dtx *apicontext.DataContext) string {
	user, err := customer.GetCustomerUserFromKey(dtx.APIKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}
	if !user.Sudo {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return ""
	}

	cust, err := user.GetCustomer(dtx.APIKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}
	err = cust.GetUsers(dtx.APIKey)
	if err != nil {
		return err.Error()
	}
	return encoding.Must(enc.Encode(cust.Users))
}

//Todo redundant
func GetUser(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, dtx *apicontext.DataContext) string {
	user, err := customer.GetCustomerUserFromKey(dtx.APIKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}
	return encoding.Must(enc.Encode(user))
}

func GetCustomerPrice(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params, dtx *apicontext.DataContext) string {
	var err error
	id := params["id"]
	var p products.Part
	if id != "" {
		p.ID, err = strconv.Atoi(params["id"])
	}

	price, err := customer.GetCustomerPrice(dtx, p.ID)
	if err != nil {
		return err.Error()
	}
	return encoding.Must(enc.Encode(price))
}

func GetCustomerCartReference(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params, dtx *apicontext.DataContext) string {
	var err error
	id := params["id"]
	var p products.Part
	if id != "" {
		p.ID, err = strconv.Atoi(params["id"])
	}

	ref, err := customer.GetCustomerCartReference(dtx.APIKey, p.ID)
	if err != nil {
		return err.Error()
	}
	return encoding.Must(enc.Encode(ref))
}

func SaveCustomer(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params, dtx *apicontext.DataContext) string {
	var c customer.Customer
	var err error
	idStr := params["id"]
	if idStr != "" {
		c.Id, err = strconv.Atoi(idStr)
		err = c.FindCustomerIdFromCustId()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return ""
		}

		err = c.Basics(dtx.APIKey)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return ""
		}
	}

	//json
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return encoding.Must(enc.Encode(false))
	}
	err = json.Unmarshal(requestBody, &c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return encoding.Must(enc.Encode(false))
	}

	//create or update
	if c.Id > 0 {
		err = c.Update()
	} else {
		err = c.Create()
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}
	return encoding.Must(enc.Encode(c))
}

func DeleteCustomer(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	var c customer.Customer
	var err error
	idStr := params["id"]
	if idStr == "" {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}
	c.Id, err = strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	err = c.Delete()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}
	return encoding.Must(enc.Encode(c))
}
