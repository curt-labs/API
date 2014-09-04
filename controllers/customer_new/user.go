package customer_ctlr_new

import (
	"github.com/curt-labs/GoAPI/helpers/encoding"
	"github.com/curt-labs/GoAPI/models/customer_new"
	// "github.com/curt-labs/GoAPI/models/part"
	"github.com/go-martini/martini"
	"log"
	"net/http"
	"strconv"
)

func ResetPassword(w http.ResponseWriter, r *http.Request, enc encoding.Encoder) string {
	email := r.FormValue("email")
	custID, err := strconv.Atoi(r.FormValue("customerID"))
	if email == "" {
		log.Print("Email err")
	}
	if custID == 0 {
		log.Print("custid is zero")
	}

	var user customer_new.CustomerUser
	user.Email = email

	resp, err := user.ResetPass(custID)
	if err != nil || resp == "" {
		http.Error(w, err.Error(), http.StatusUnauthorized)
	}
	return encoding.Must(enc.Encode(resp))
}

func ChangePassword(w http.ResponseWriter, r *http.Request, enc encoding.Encoder) string {
	email := r.FormValue("email")
	custID, err := strconv.Atoi(r.FormValue("customerID"))
	oldPass := r.FormValue("oldPass")
	newPass := r.FormValue("newPass")

	var user customer_new.CustomerUser
	user.Email = email

	resp, err := user.ChangePass(oldPass, newPass, custID)
	if err != nil || resp == "" {
		http.Error(w, err.Error(), http.StatusUnauthorized)
	}
	return encoding.Must(enc.Encode(resp))
}

func RegisterUser(w http.ResponseWriter, r *http.Request, enc encoding.Encoder) string {
	name := r.FormValue("name")
	email := r.FormValue("email")
	pass := r.FormValue("pass")
	customerID, err := strconv.Atoi(r.FormValue("customerID"))
	isActive, err := strconv.ParseBool(r.FormValue("isActive"))
	locationID, err := strconv.Atoi(r.FormValue("locationID"))
	isSudo, err := strconv.ParseBool(r.FormValue("isSudo"))
	cust_ID, err := strconv.Atoi(r.FormValue("cust_ID"))
	notCustomer, err := strconv.ParseBool(r.FormValue("notCustomer"))

	var user customer_new.CustomerUser
	user.Email = email
	user.Name = name
	cu, err := user.Register(pass, customerID, isActive, locationID, isSudo, cust_ID, notCustomer)
	if err != nil || cu == nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return ""
	}

	return encoding.Must(enc.Encode(user))
}
func DeleteCustomerUser(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	id := params["id"]

	var cu customer_new.CustomerUser
	cu.Id = id
	err := cu.Delete()
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return ""
	}

	return encoding.Must(enc.Encode(cu))
}

func AuthenticateUser(w http.ResponseWriter, r *http.Request, enc encoding.Encoder) string {
	email := r.FormValue("email")
	pass := r.FormValue("pass")

	var user customer_new.CustomerUser
	user.Email = email

	cust, err := user.UserAuthentication(pass)
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return ""
	}

	return encoding.Must(enc.Encode(cust))
}
