package customer_ctlr_new

import (
	"github.com/curt-labs/GoAPI/helpers/encoding"
	"github.com/curt-labs/GoAPI/models/customer_new"
	"github.com/curt-labs/GoAPI/models/part"
	"github.com/go-martini/martini"
	"net/http"
	"strconv"
)

func UserAuthentication(w http.ResponseWriter, r *http.Request, enc encoding.Encoder) string {
	email := r.FormValue("email")
	pass := r.FormValue("password")

	user := customer_new.CustomerUser{
		Email: email,
	}
	cust, err := user.UserAuthentication(pass)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	return encoding.Must(enc.Encode(cust))
}

func KeyedUserAuthentication(w http.ResponseWriter, r *http.Request, enc encoding.Encoder) string {
	qs := r.URL.Query()
	key := qs.Get("key")
	cust, err := customer_new.UserAuthenticationByKey(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	return encoding.Must(enc.Encode(cust))
}
func ResetAuthentication(w http.ResponseWriter, r *http.Request, enc encoding.Encoder) string { //Testing only
	qs := r.URL.Query()
	id := qs.Get("id")
	var u customer_new.CustomerUser
	u.Id = id
	err := u.ResetAuthentication()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	return "Success"
}

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

	err = c.GetCustomer()
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
	c := customer_new.Customer{
		Id: id,
	}
	err = c.GetLocations()
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
	users, err := c.GetUsers()
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

	part, err := customer_new.GetCustomerPrice(key, p.PartId)
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
	var p part.Part
	if id != "" {
		p.PartId, err = strconv.Atoi(params["id"])
	}

	ref, err := customer_new.GetCustomerCartReference(key, p.PartId)
	if err != nil {
		return err.Error()
	}
	return encoding.Must(enc.Encode(ref))
}
