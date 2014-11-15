package customer_ctlr_new

import (
	"encoding/json"
	"github.com/curt-labs/GoAPI/helpers/encoding"
	"github.com/curt-labs/GoAPI/models/customer_new"
	"github.com/curt-labs/GoAPI/models/products"
	"github.com/go-martini/martini"
	"io/ioutil"
	"strings"

	// "log"
	"net/http"
	"strconv"
)

func GetCustomer(w http.ResponseWriter, r *http.Request, enc encoding.Encoder) string {
	qs := r.URL.Query()
	key := qs.Get("key")
	if key == "" {
		key = r.FormValue("key")
	}

	if key == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return ""
	}

	var c customer_new.Customer
	var err error

	//get id from key
	err = c.GetCustomerIdFromKey(key)
	if err != nil {
		http.Error(w, "Error getting customer.", http.StatusServiceUnavailable)
		return ""
	}

	err = c.GetCustomer(key)
	if err != nil {
		http.Error(w, "Error getting customer.", http.StatusServiceUnavailable)
		return ""
	}

	lowerKey := strings.ToLower(key)
	for i, u := range c.Users {
		for _, k := range u.Keys {
			if strings.ToLower(k.Key) == lowerKey {
				c.Users[i].Current = true
			}
		}
	}

	return encoding.Must(enc.Encode(c))
}

func GetLocations(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	qs := r.URL.Query()
	key := qs.Get("key")
	if key == "" {
		key = r.FormValue("key")
	}

	if key == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return ""
	}

	cu, err := customer_new.GetCustomerUserFromKey(key)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return ""
	}

	c, err := cu.GetCustomer(key)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return ""
	}

	return encoding.Must(enc.Encode(c.Locations))
}

//TODO - redundant
func GetUsers(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	qs := r.URL.Query()
	key := qs.Get("key")
	if key == "" {
		key = r.FormValue("key")
	}

	if key == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return ""
	}
	user, err := customer_new.GetCustomerUserFromKey(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}
	if !user.Sudo {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return ""
	}

	cust, err := user.GetCustomer(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}
	err = cust.GetUsers(key)
	if err != nil {
		return err.Error()
	}
	return encoding.Must(enc.Encode(cust.Users))
}

//Todo redundant
func GetUser(w http.ResponseWriter, r *http.Request, enc encoding.Encoder) string {
	key := r.FormValue("key")

	if key == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return ""
	}

	user, err := customer_new.GetCustomerUserFromKey(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	return encoding.Must(enc.Encode(user))
}

func GetCustomerPrice(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	var err error
	qs := r.URL.Query()
	key := qs.Get("key")
	if key == "" {
		key = r.FormValue("key")
	}

	if key == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return ""
	}

	id := params["id"]
	var p products.Part
	if id != "" {
		p.ID, err = strconv.Atoi(params["id"])
	}

	part, err := customer_new.GetCustomerPrice(key, p.ID)
	if err != nil {
		return err.Error()
	}
	return encoding.Must(enc.Encode(part))
}

func GetCustomerCartReference(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	var err error
	qs := r.URL.Query()
	key := qs.Get("key")
	if key == "" {
		key = r.FormValue("key")
	}

	if key == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return ""
	}

	id := params["id"]
	var p products.Part
	if id != "" {
		p.ID, err = strconv.Atoi(params["id"])
	}

	ref, err := customer_new.GetCustomerCartReference(key, p.ID)
	if err != nil {
		return err.Error()
	}
	return encoding.Must(enc.Encode(ref))
}

func SaveCustomer(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	qs := r.URL.Query()
	key := qs.Get("key")
	if key == "" {
		key = r.FormValue("key")
	}

	if key == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return ""
	}

	var c customer_new.Customer
	var err error
	idStr := params["id"]
	if idStr != "" {
		c.Id, err = strconv.Atoi(idStr)
		err = c.FindCustomerIdFromCustId()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return ""
		}

		err = c.Basics(key)
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
	var c customer_new.Customer
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
