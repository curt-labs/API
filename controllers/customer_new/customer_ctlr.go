package customer_ctlr_new

import (
	"encoding/json"
	"github.com/curt-labs/GoAPI/helpers/encoding"
	"github.com/curt-labs/GoAPI/models/customer_new"
	"github.com/curt-labs/GoAPI/models/products"
	"github.com/go-martini/martini"
	"io/ioutil"
	"log"
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
		log.Print(err)
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

	userChan := make(chan error)
	go func() {
		user, err := customer_new.GetCustomerUserFromKey(key)
		if err == nil {
			user.Current = true
		}
		c.Users = append(c.Users, user)
		userChan <- err
	}()

	err = c.GetCustomer()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	<-userChan

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
	user, err := customer_new.GetCustomerUserFromKey(key)
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
	users, err := cust.GetUsers()
	if err != nil {
		return err.Error()
	}
	return encoding.Must(enc.Encode(users))
}

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
	var c customer_new.Customer
	var err error
	idStr := params["id"]
	if idStr != "" {
		c.Id, err = strconv.Atoi(idStr)
		log.Print(c.Id)
		err = c.FindCustomerIdFromCustId()
		log.Print(c.Id, " ", c.CustomerId)
		err = c.Basics()
		log.Print(c.Id, " ", c.CustomerId)
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
