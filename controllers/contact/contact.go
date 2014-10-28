package contact

import (
	"encoding/json"
	"fmt"
	"github.com/curt-labs/GoAPI/helpers/encoding"
	"github.com/curt-labs/GoAPI/models/contact"
	"github.com/curt-labs/GoAPI/models/geography"
	"github.com/go-martini/martini"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
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
	var c contact.Contact
	var ct contact.ContactType
	var sendEmail bool

	contType := req.Header.Get("Content-Type")

	if strings.Contains(contType, "application/json") {
		//this is our json payload
		var formData map[string]interface{}

		requestBody, err := ioutil.ReadAll(req.Body)
		if err != nil {

			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return err.Error()
		}

		if err = json.Unmarshal(requestBody, &formData); err != nil {

			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return err.Error()
		}

		//require contact type
		if str_id, found := formData["contactType"]; !found {
			http.Error(rw, "Invalid Contact Type ID", http.StatusInternalServerError)
			return "Invalid Contact Type ID"
		} else {
			// if ct.ID, err = strconv.Atoi(str_id.(string)); err != nil {
			if ct.ID = int(str_id.(float64)); err != nil {
				return "Invalid Contact Type ID"
			}
			if err = ct.Get(); err != nil {
				return err.Error()
			}
			c.Type = ct.Name
		}

		//require email
		if email, found := formData["email"]; !found {
			http.Error(rw, "Email is required", http.StatusInternalServerError)
			return "Email is required"
		} else {
			c.Email = email.(string)
		}

		//require first name
		if first, found := formData["firstName"]; !found {
			http.Error(rw, "First name is required", http.StatusInternalServerError)
			return "First name is required"
		} else {
			c.FirstName = first.(string)
		}

		//require last name
		if last, found := formData["lastName"]; !found {
			http.Error(rw, "Last name is required", http.StatusInternalServerError)
			return "Last name is required"
		} else {
			c.LastName = last.(string)
		}

		if phone, found := formData["phoneNumber"]; found {
			c.Phone = phone.(string)
		}
		if address1, found := formData["address1"]; found {
			c.Address1 = address1.(string)
		}
		if address1, found := formData["address2"]; found {
			c.Address1 = address1.(string)
		}
		if city, found := formData["city"]; found {
			c.City = city.(string)
		}
		if st, found := formData["state"]; found {
			id, err := strconv.Atoi(st.(string))
			if id > 0 && err == nil {
				countries, err := geography.GetAllCountriesAndStates()
				if err == nil {
					for _, ctry := range countries {
						for _, state := range *ctry.States {
							if state.Id == id {
								c.State = state.State
								c.Country = ctry.Country
								break
							}
						}
					}
				}
			}
		}

		if postal, found := formData["postalCode"]; found {
			c.PostalCode = postal.(string)
		}

		switch ct.ID {
		case 15: //become a dealer

			//require business name
			businessName, found := formData["businessName"]
			if !found {
				http.Error(rw, "Business name is required", http.StatusInternalServerError)
				return "Business name is required"
			}
			//require business type
			businessType, found := formData["businessType"]
			if !found {
				http.Error(rw, "Business type is required", http.StatusInternalServerError)
				return "Business type is required"
			}

			formMessage, found := formData["message"]
			if found {
				c.Message = formMessage.(string)
			}

			c.Subject = "Become a Dealer"
			c.Message = fmt.Sprintf(
				`This contact is interested in becoming a Dealer.\n
				Name: %s\n
				Business: %s\n
				Business Type: %s\n
				Email: %s\n
				Phone: %s\n
				Address 1: %s\n
				Address 2: %s\n
				City, State, Zip: %s, %s %s\n
				Country: %s\n
				Message: %s\n`,
				c.FirstName+" "+c.LastName,
				businessName,
				businessType,
				c.Email,
				c.Phone,
				c.Address1,
				c.Address2,
				c.City, c.State, c.PostalCode,
				c.Country,
				c.Message,
			)

		default: //everything else
			if subject, found := formData["subject"]; found {
				c.Subject = subject.(string)
			}
			if message, found := formData["message"]; found {
				c.Message = message.(string)
			}
		}

		if send_email, found := formData["sendEmail"]; found {
			sendEmail = send_email.(bool)
		}
	} else { //form post parameters
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

		//TODO: this needs work

		if req.FormValue("send_email") != "" {
			sendEmail, _ = strconv.ParseBool(req.FormValue("send_email"))
		}
	}
	if err := c.Add(); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return err.Error()
	}
	if sendEmail {
		var emailBody string

		switch ct.ID {
		case 15: //Become a dealer
			emailBody = c.Message
		default: //everything else
			emailBody = fmt.Sprintf(
				`From: %s %s\n
				 Phone: %s\n
				 Subject: %s\n
				 Time: %s\n
				 Type: %s\n
				 Address1: %s\n
				 Address2: %s\n
				 City: %s\n
				 State: %s\n
				 PostalCode: %s\n
				 Country: %s\n
				 Message: %s\n`,
				c.Email,
				c.Phone,
				c.Subject,
				c.Created.String(),
				c.Type,
				c.Address1,
				c.Address2,
				c.City,
				c.State,
				c.PostalCode,
				c.Country,
				c.Message,
			)
		}
		if emailBody != "" {
			subject := "Email from Contact Form"
			if err := contact.SendEmail(ct, subject, emailBody); err != nil {
				http.Error(rw, err.Error(), http.StatusInternalServerError)
				return err.Error()
			}

		}
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
