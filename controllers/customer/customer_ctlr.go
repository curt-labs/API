package customer_ctlr

import (
	"github.com/curt-labs/GoAPI/helpers/encoding"
	"github.com/curt-labs/GoAPI/models/customer"
	"net/http"
)

func UserAuthentication(w http.ResponseWriter, r *http.Request, enc encoding.Encoder) string {
	email := r.FormValue("email")
	pass := r.FormValue("password")

	user := customer.CustomerUser{
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

	cust, err := customer.UserAuthenticationByKey(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	return encoding.Must(enc.Encode(cust))
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

	id, err := customer.GetCustomerIdFromKey(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	c := customer.Customer{
		Id: id,
	}

	err = c.GetCustomer()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	return encoding.Must(enc.Encode(c))
}

func GetLocations(w http.ResponseWriter, r *http.Request, enc encoding.Encoder) string {
	qs := r.URL.Query()
	key := qs.Get("key")
	if key == "" {
		key = r.FormValue("key")
	}

	if key == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return ""
	}

	id, err := customer.GetCustomerIdFromKey(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	c := customer.Customer{
		Id: id,
	}

	err = c.GetLocations()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	return encoding.Must(enc.Encode(c.Locations))
}

func GetUsers(w http.ResponseWriter, r *http.Request, enc encoding.Encoder) string {
	params := r.URL.Query()
	key := params.Get("key")
	if key == "" {
		key = r.FormValue("key")
	}

	if key == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return ""
	}

	user, err := customer.GetCustomerUserFromKey(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}
	if !user.Sudo {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return ""
	}

	cust, err := user.GetCustomer()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	if err = cust.GetUsers(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	return encoding.Must(enc.Encode(cust.Users))
}

func GetUser(w http.ResponseWriter, r *http.Request, enc encoding.Encoder) string {
	key := r.FormValue("key")
	if key == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return ""
	}

	user, err := customer.GetCustomerUserFromKey(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	return encoding.Must(enc.Encode(user))
}
