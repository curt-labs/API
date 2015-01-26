package warranty

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/curt-labs/GoAPI/helpers/apicontext"
	"github.com/curt-labs/GoAPI/helpers/encoding"
	"github.com/curt-labs/GoAPI/helpers/error"
	"github.com/curt-labs/GoAPI/models/contact"
	"github.com/curt-labs/GoAPI/models/warranty"
	"github.com/go-martini/martini"
)

const (
	timeFormat = "2006-01-02"
)

func GetAllWarranties(rw http.ResponseWriter, req *http.Request, enc encoding.Encoder, dtx *apicontext.DataContext) string {
	var err error

	ws, err := warranty.GetAllWarranties(dtx)
	if err != nil {
		apierror.GenerateError("Trouble getting all warranties", err, rw, req)
		return ""
	}
	return encoding.Must(enc.Encode(ws))
}

func GetWarranty(rw http.ResponseWriter, req *http.Request, enc encoding.Encoder, params martini.Params) string {
	var err error
	var w warranty.Warranty
	id := params["id"]
	w.ID, err = strconv.Atoi(id)

	err = w.Get()
	if err != nil {
		apierror.GenerateError("Trouble getting warranty", err, rw, req)
		return ""
	}
	return encoding.Must(enc.Encode(w))
}

func GetWarrantyByContact(rw http.ResponseWriter, req *http.Request, enc encoding.Encoder, params martini.Params) string {
	var err error
	var w warranty.Warranty
	id := params["id"]
	w.Contact.ID, err = strconv.Atoi(id)
	if err != nil {
		apierror.GenerateError("Trouble parsing contact ID.", err, rw, req)
		return ""
	}

	ws, err := w.GetByContact()
	if err != nil {
		apierror.GenerateError("Trouble getting warranty by contact.", err, rw, req)
		return ""
	}
	return encoding.Must(enc.Encode(ws))
}

func CreateWarranty(rw http.ResponseWriter, req *http.Request, enc encoding.Encoder, params martini.Params) string {
	contType := req.Header.Get("Content-Type")
	var w warranty.Warranty
	var err error

	contactTypeID, err := strconv.Atoi(params["contactReceiverTypeID"]) //to whom the emails go
	if err != nil {
		apierror.GenerateError("Trouble parsing contact type ID.", err, rw, req)
		return ""
	}
	sendEmail, err := strconv.ParseBool(params["sendEmail"])
	if err != nil {
		apierror.GenerateError("Trouble parsing send email.", err, rw, req)
		return ""
	}
	if strings.Contains(contType, "application/json") {
		//json
		requestBody, err := ioutil.ReadAll(req.Body)
		if err != nil {
			apierror.GenerateError("Trouble reading request body for creating warranty", err, rw, req)
			return ""
		}

		err = json.Unmarshal(requestBody, &w)
		if err != nil {
			apierror.GenerateError("Trouble unmarshalling request body for creating warranty", err, rw, req)
			return ""
		}

	} else {
		//else, form
		w.PartNumber, err = strconv.Atoi(req.FormValue("part_number"))
		w.OldPartNumber = req.FormValue("old_part_number")
		date, err := time.Parse(timeFormat, req.FormValue("date"))
		if err != nil {
			apierror.GenerateError("Trouble creating warranty", err, rw, req)
			return ""
		}
		w.Date = &date
		w.SerialNumber = req.FormValue("serial_number")

		w.Contact.FirstName = req.FormValue("first_name")
		w.Contact.LastName = req.FormValue("last_name")
		w.Contact.Email = req.FormValue("email")
		w.Contact.Phone = req.FormValue("phone")
		w.Contact.Type = req.FormValue("type")
		w.Contact.Address1 = req.FormValue("address1")
		w.Contact.Address2 = req.FormValue("address2")
		w.Contact.City = req.FormValue("city")
		w.Contact.State = req.FormValue("state")
		w.Contact.PostalCode = req.FormValue("postal_code")
		w.Contact.Country = req.FormValue("country")
	}
	err = w.Create()
	if err != nil {
		apierror.GenerateError("Trouble creating warranty", err, rw, req)
		return ""
	}
	if sendEmail == true {
		//Send Email
		body :=
			"Name: " + w.Contact.FirstName + " " + w.Contact.LastName + "\n" +
				"Email: " + w.Contact.Email + "\n" +
				"Phone: " + w.Contact.Phone + "\n" +
				"Serial Number: " + w.SerialNumber + "\n" +
				"Date: " + w.Date.String() + "\n" +
				"Part Number: " + strconv.Itoa(w.PartNumber) + "\n"

		var ct contact.ContactType
		ct.ID = contactTypeID
		subject := "Email from Warranty Applications Form"
		err = contact.SendEmail(ct, subject, body) //contact type id, subject, techSupport
		if err != nil {
			apierror.GenerateError("Trouble sending email to receivers while creating warranty", err, rw, req)
			return ""
		}
	}
	//Return JSON
	return encoding.Must(enc.Encode(w))
}

func DeleteWarranty(rw http.ResponseWriter, req *http.Request, enc encoding.Encoder, params martini.Params) string {
	var err error
	var w warranty.Warranty
	id := params["id"]

	if w.ID, err = strconv.Atoi(id); err != nil {
		apierror.GenerateError("Trouble getting warranty ID", err, rw, req)
		return ""
	}

	if err = w.Delete(); err != nil {
		apierror.GenerateError("Trouble deleting warranty", err, rw, req)
		return ""
	}

	//Return JSON
	return encoding.Must(enc.Encode(w))
}
