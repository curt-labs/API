package customer_ctlr_new

import (
	"github.com/curt-labs/GoAPI/helpers/encoding"
	"github.com/curt-labs/GoAPI/models/customer_new"
	"github.com/curt-labs/GoAPI/models/part"
	"github.com/go-martini/martini"
	"log"
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

	id, err := customer_new.GetCustomerIdFromKey(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	c := customer_new.Customer{
		Id: id,
	}

	err = c.GetCustomer_New()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
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

	id, err := customer_new.GetCustomerIdFromKey(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}
	log.Print(id)
	c := customer_new.Customer{
		Id: id,
	}
	err = c.GetLocations_New()
	if err != nil {
		return err.Error()
	}
	return encoding.Must(enc.Encode(c.Locations))
}

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
	//TODO getCustmerUserFromKey - check if Sudo
	id, err := customer_new.GetCustomerIdFromKey(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	c := customer_new.Customer{
		Id: id,
	}
	log.Print("CTRL")
	users, err := c.GetUsers_New()
	if err != nil {
		return err.Error()
	}
	return encoding.Must(enc.Encode(users))
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
	var p part.Part
	if id != "" {
		p.PartId, err = strconv.Atoi(params["id"])
	}

	part, err := customer_new.GetCustomerPrice_New(key, p.PartId)
	if err != nil {
		return err.Error()
	}
	log.Print(part)
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
	var p part.Part
	if id != "" {
		p.PartId, err = strconv.Atoi(params["id"])
	}

	ref, err := customer_new.GetCustomerCartReference_New(key, p.PartId)
	if err != nil {
		return err.Error()
	}
	return encoding.Must(enc.Encode(ref))
}
