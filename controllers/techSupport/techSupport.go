package techSupport

import (
	"encoding/json"
	"github.com/curt-labs/GoAPI/helpers/encoding"
	"github.com/curt-labs/GoAPI/models/contact"
	"github.com/curt-labs/GoAPI/models/techSupport"
	"github.com/go-martini/martini"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

const (
	timeFormat = "2006-01-02"
)

func GetAllTechSupport(rw http.ResponseWriter, req *http.Request, enc encoding.Encoder) string {
	var err error

	ts, err := techSupport.GetAllTechSupport()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return err.Error()
	}
	return encoding.Must(enc.Encode(ts))
}

func GetTechSupport(rw http.ResponseWriter, req *http.Request, enc encoding.Encoder, params martini.Params) string {
	var err error
	var t techSupport.TechSupport
	id := params["id"]
	t.ID, err = strconv.Atoi(id)

	err = t.Get()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return err.Error()
	}
	return encoding.Must(enc.Encode(t))
}
func GetTechSupportByContact(rw http.ResponseWriter, req *http.Request, enc encoding.Encoder, params martini.Params) string {
	var err error
	var t techSupport.TechSupport
	id := params["id"]
	t.Contact.ID, err = strconv.Atoi(id)

	ts, err := t.GetByContact()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return err.Error()
	}
	return encoding.Must(enc.Encode(ts))
}

func CreateTechSupport(rw http.ResponseWriter, req *http.Request, enc encoding.Encoder, params martini.Params) string {
	contType := req.Header.Get("Content-Type")

	var t techSupport.TechSupport
	var err error

	contactTypeID, err := strconv.Atoi(params["contactReceiverTypeID"]) //to whom the emails go

	if contType == "application/json" {
		//json
		requestBody, err := ioutil.ReadAll(req.Body)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return encoding.Must(enc.Encode(false))
		}

		err = json.Unmarshal(requestBody, &t)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return encoding.Must(enc.Encode(false))
		}
	} else {
		//else, form
		t.VehicleMake = req.FormValue("vehicle_make")
		t.VehicleModel = req.FormValue("vehicle_model")
		t.VehicleYear, err = strconv.Atoi(req.FormValue("vehicle_year"))
		t.PurchaseDate, err = time.Parse(timeFormat, req.FormValue("purchase_date"))
		t.PurchasedFrom = req.FormValue("purchased_from")
		t.DealerName = req.FormValue("dealer_name")
		t.ProductCode = req.FormValue("product_code")
		t.DateCode = req.FormValue("date_code")
		t.Issue = req.FormValue("issue")

		t.Contact.FirstName = req.FormValue("first_name")
		t.Contact.LastName = req.FormValue("last_name")
		t.Contact.Email = req.FormValue("email")
		t.Contact.Phone = req.FormValue("phone")
		t.Contact.Subject = req.FormValue("subject")
		t.Contact.Message = req.FormValue("message")
		t.Contact.Type = req.FormValue("type")
		t.Contact.Address1 = req.FormValue("address1")
		t.Contact.Address2 = req.FormValue("address2")
		t.Contact.City = req.FormValue("city")
		t.Contact.State = req.FormValue("state")
		t.Contact.PostalCode = req.FormValue("postal_code")
		t.Contact.Country = req.FormValue("country")
	}
	err = t.Create()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return err.Error()
	}

	//Send Email
	body :=
		"Name: " + t.Contact.FirstName + " " + t.Contact.LastName + "\n" +
			"Email: " + t.Contact.Email + "\n" +
			"Phone: " + t.Contact.Phone + "\n" +
			"Make: " + t.VehicleMake + "\n" +
			"Model: " + t.VehicleModel + "\n" +
			"Year: " + strconv.Itoa(t.VehicleYear) + "\n" +
			"Purchase Date: " + t.PurchaseDate.String() + "\n" +
			"Purchased From: " + t.PurchasedFrom + "\n" +
			"Dealer Name: " + t.DealerName + "\n" +
			"Product Code: " + t.ProductCode + "\n" +
			"Date Code: " + t.DateCode + "\n\n" +
			"Issue: " + t.Issue + "\n"

	var ct contact.ContactType
	ct.ID = contactTypeID
	subject := "Email from Aries Tech Support"
	err = contact.SendEmail(ct, subject, body) //contact type id, subject, techSupport
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return err.Error()
	}

	//Return JSON
	return encoding.Must(enc.Encode(t))
}

func DeleteTechSupport(rw http.ResponseWriter, req *http.Request, enc encoding.Encoder, params martini.Params) string {
	var err error
	var t techSupport.TechSupport
	id, err := strconv.Atoi(params["id"])
	t.ID = id
	err = t.Delete()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return err.Error()
	}

	//Return JSON
	return encoding.Must(enc.Encode(t))
}
