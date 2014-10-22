package contact

import (
	"encoding/json"
	"github.com/curt-labs/GoAPI/helpers/encoding"
	"github.com/curt-labs/GoAPI/models/contact"
	"github.com/go-martini/martini"
	"io/ioutil"
	"net/http"
	"strconv"
)

func GetAllContacts(rw http.ResponseWriter, req *http.Request, enc encoding.Encoder) string {
	var err error
	var page int = 1
	var count int = 50

	if req.FormValue("page") != "" {
		if page, err = strconv.Atoi(req.FormValue("page")); err != nil {
			page = 1
		}
	}
	if req.FormValue("count") != "" {
		if count, err = strconv.Atoi(req.FormValue("count")); err != nil {
			count = 50
		}
	}

	contacts, err := contact.GetAllContacts(page, count)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return err.Error()
	}
	return encoding.Must(enc.Encode(contacts))
}

func GetContact(rw http.ResponseWriter, req *http.Request, params martini.Params, enc encoding.Encoder) string {
	var err error
	var c contact.Contact

	if c.ID, err = strconv.Atoi(params["id"]); err != nil {
		http.Error(rw, "Invalid Contact ID", http.StatusInternalServerError)
		return "Invalid Contact ID"
	}
	if err = c.Get(); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return err.Error()
	}
	return encoding.Must(enc.Encode(c))
}

func AddContact(rw http.ResponseWriter, req *http.Request, enc encoding.Encoder) string {
	contType := req.Header.Get("Content-Type")

	var c contact.Contact
	if contType == "application/json" {
		//json
		requestBody, err := ioutil.ReadAll(req.Body)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return encoding.Must(enc.Encode(false))
		}

		err = json.Unmarshal(requestBody, &c)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return encoding.Must(enc.Encode(false))
		}
	} else {
		//else, form
		c = contact.Contact{
			FirstName:  req.FormValue("first_name"),
			LastName:   req.FormValue("last_name"),
			Email:      req.FormValue("email"),
			Phone:      req.FormValue("phone"),
			Subject:    req.FormValue("subject"),
			Message:    req.FormValue("message"),
			Type:       req.FormValue("type"),
			Address1:   req.FormValue("address1"),
			Address2:   req.FormValue("address2"),
			City:       req.FormValue("city"),
			State:      req.FormValue("state"),
			PostalCode: req.FormValue("postal_code"),
			Country:    req.FormValue("country"),
		}
	}
	if err := c.Add(); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return err.Error()
	}
	return encoding.Must(enc.Encode(c))
}

func UpdateContact(rw http.ResponseWriter, req *http.Request, params martini.Params, enc encoding.Encoder) string {
	var err error
	var c contact.Contact
	if c.ID, err = strconv.Atoi(params["id"]); err != nil {
		http.Error(rw, "Invalid Contact ID", http.StatusInternalServerError)
		return "Invalid Contact ID"
	}

	if err = c.Get(); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return err.Error()
	}

	contType := req.Header.Get("Content-Type")
	if contType == "application/json" {
		//json
		requestBody, err := ioutil.ReadAll(req.Body)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return encoding.Must(enc.Encode(false))
		}

		err = json.Unmarshal(requestBody, &c)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return encoding.Must(enc.Encode(false))
		}
	} else {
		if req.FormValue("first_name") != "" {
			c.FirstName = req.FormValue("first_name")
		}

		if req.FormValue("last_name") != "" {
			c.LastName = req.FormValue("last_name")
		}

		if req.FormValue("email") != "" {
			c.Email = req.FormValue("email")
		}

		if req.FormValue("phone") != "" {
			c.Phone = req.FormValue("phone")
		}

		if req.FormValue("subject") != "" {
			c.Subject = req.FormValue("subject")
		}

		if req.FormValue("message") != "" {
			c.Message = req.FormValue("message")
		}

		if req.FormValue("type") != "" {
			c.Type = req.FormValue("type")
		}

		if req.FormValue("address1") != "" {
			c.Address1 = req.FormValue("address1")
		}

		if req.FormValue("address2") != "" {
			c.Address2 = req.FormValue("address2")
		}

		if req.FormValue("city") != "" {
			c.City = req.FormValue("city")
		}

		if req.FormValue("state") != "" {
			c.State = req.FormValue("state")
		}

		if req.FormValue("postal_code") != "" {
			c.PostalCode = req.FormValue("postal_code")
		}

		if req.FormValue("country") != "" {
			c.Country = req.FormValue("country")
		}
	}
	if err = c.Update(); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return err.Error()
	}
	return encoding.Must(enc.Encode(c))
}

func DeleteContact(rw http.ResponseWriter, req *http.Request, params martini.Params, enc encoding.Encoder) string {
	var err error
	var c contact.Contact

	if c.ID, err = strconv.Atoi(params["id"]); err != nil {
		http.Error(rw, "Invalid Contact ID", http.StatusInternalServerError)
		return "Invalid Contact ID"
	}

	if err = c.Delete(); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return err.Error()
	}

	return encoding.Must(enc.Encode(c))
}
