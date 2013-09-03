package customer_ctlr

import (
	"../../helpers/plate"
	. "../../models"
	"net/http"
)

func UserAuthentication(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	pass := r.FormValue("password")

	user := CustomerUser{
		Email: email,
	}
	cust, err := user.UserAuthentication(pass)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	plate.ServeFormatted(w, r, cust)
}

func KeyedUserAuthentication(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	key := params.Get("key")

	cust, err := UserAuthenticationByKey(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	plate.ServeFormatted(w, r, cust)
}

func GetCustomer(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	key := params.Get("key")
	if key == "" {
		key = r.FormValue("key")
	}

	if key == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	id, err := GetCustomerIdFromKey(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c := Customer{
		Id: id,
	}

	err = c.GetCustomer()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	plate.ServeFormatted(w, r, c)
}

func GetLocations(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	key := params.Get("key")
	if key == "" {
		key = r.FormValue("key")
	}

	if key == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	id, err := GetCustomerIdFromKey(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c := Customer{
		Id: id,
	}

	err = c.GetLocations()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	plate.ServeFormatted(w, r, c.Locations)
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	key := params.Get("key")
	if key == "" {
		key = r.FormValue("key")
	}

	if key == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := GetCustomerUserFromKey(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if user.Sudo {
		cust, err := user.GetCustomer()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		users, err := cust.GetUsers()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		plate.ServeFormatted(w, r, users)
	} else {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	key := r.FormValue("key")

	if key == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := GetCustomerUserFromKey(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	plate.ServeFormatted(w, r, user)
}
