package contact

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/curt-labs/API/helpers/apicontext"
	"github.com/curt-labs/API/helpers/encoding"
	"github.com/curt-labs/API/helpers/error"
	"github.com/curt-labs/API/models/brand"
	"github.com/curt-labs/API/models/contact"
	"github.com/curt-labs/API/models/geography"
	"github.com/go-martini/martini"
)

var (
	noEmail = flag.Bool("noEmail", false, "Do not send email")
)

func GetAllContacts(rw http.ResponseWriter, req *http.Request, enc encoding.Encoder, dtx *apicontext.DataContext) string {
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
	contacts, err := contact.GetAllContacts(page, count, dtx)
	if err != nil {
		apierror.GenerateError("Trouble getting all contacts", err, rw, req)
	}
	return encoding.Must(enc.Encode(contacts))
}

func GetContact(rw http.ResponseWriter, req *http.Request, params martini.Params, enc encoding.Encoder) string {
	var err error
	var c contact.Contact

	if c.ID, err = strconv.Atoi(params["id"]); err != nil {
		apierror.GenerateError("Trouble getting contact ID", err, rw, req)
	}
	if err = c.Get(); err != nil {
		apierror.GenerateError("Trouble getting contact", err, rw, req)
	}
	return encoding.Must(enc.Encode(c))
}

func AddDealerContact(rw http.ResponseWriter, req *http.Request, enc encoding.Encoder, params martini.Params, dtx *apicontext.DataContext) string {
	flag.Parse()

	var d contact.DealerContact
	var ct contact.ContactType
	var subject string
	var err error
	var brandName string

	brand := brand.Brand{ID: dtx.BrandID}
	err = brand.Get()
	if err != nil {
		brandName = "Unknown Brand"
	} else {
		brandName = brand.Name
	}
	d.Brand = brand

	ct.ID, err = strconv.Atoi(params["contactTypeID"]) //determines to whom emails go
	if err != nil {
		apierror.GenerateError("Trouble getting contact type ID", err, rw, req)
		return ""
	}

	contType := req.Header.Get("Content-Type")

	if strings.Contains(contType, "application/json") {
		//this is our json payload
		requestBody, err := ioutil.ReadAll(req.Body)
		if err != nil {
			apierror.GenerateError("Trouble reading JSON request body for adding contact", err, rw, req)
			return ""
		}

		if err = json.Unmarshal(requestBody, &d); err != nil {
			apierror.GenerateError("Trouble unmarshalling JSON request body into contact", err, rw, req)
			return ""
		}
	}

	d.Type, err = contact.GetContactTypeNameFromId(ct.ID)
	if err != nil {
		apierror.GenerateError("Trouble getting contact type from Id.", err, rw, req)
		return ""
	}

	if d.Type == "" {
		d.Type = "New Customer"
		d.Subject = "Becoming a Dealer"
		subject = d.Subject
	} else {
		d.Subject = "Contact Submission regarding " + d.Type + ", from " + d.FirstName + " " + d.LastName + "."
		subject = d.Subject
	}

	//state/country
	st := d.State
	id, err := strconv.Atoi(st)
	if id > 0 && err == nil {
		countries, err := geography.GetAllCountriesAndStates()
		if err == nil {
			for _, ctry := range countries {
				for _, state := range *ctry.States {
					if state.Id == id {
						d.State = state.State
						d.Country = ctry.Country
						break
					}
				}
			}
		}
	}
	if err := d.Add(dtx); err != nil {
		apierror.GenerateError("Trouble adding contact", err, rw, req)
		return ""
	}

	emailBody := fmt.Sprintf(
		"This %s contact is inquiring about: %s.\n"+
			"Name: %s\n"+
			"Email: %s\n"+
			"Phone: %s\n"+
			"Address 1: %s\n"+
			"Address 2: %s\n"+
			"City, State, Zip: %s, %s %s\n"+
			"Country: %s\n"+
			"Subject: %s\n"+
			"Message: %s\n",
		brandName,
		d.Type,
		d.FirstName+" "+d.LastName,
		d.Email,
		d.Phone,
		d.Address1,
		d.Address2,
		d.City, d.State, d.PostalCode,
		d.Country,
		d.Subject,
		d.Message,
	)

	emailAppend1 := fmt.Sprintf(
		`Business: %s
		Business Type: %s`,
		d.BusinessName,
		d.BusinessType.Type,
	)
	if d.BusinessName != "" || d.BusinessType.Type != "" {
		emailBody += emailAppend1
	}

	if emailBody != "" && *noEmail == false {
		if err := contact.SendEmail(ct, subject, emailBody); err != nil {
			apierror.GenerateError("Trouble sending email to receivers", err, rw, req)
			return ""
		}
	}
	return encoding.Must(enc.Encode(d))
}

func UpdateContact(rw http.ResponseWriter, req *http.Request, params martini.Params, enc encoding.Encoder) string {
	var err error
	var c contact.Contact
	if c.ID, err = strconv.Atoi(params["id"]); err != nil {
		apierror.GenerateError("Trouble getting contact ID", err, rw, req)
	}

	if err = c.Get(); err != nil {
		apierror.GenerateError("Trouble getting contact", err, rw, req)
	}

	contType := req.Header.Get("Content-Type")
	if contType == "application/json" {
		//json
		requestBody, err := ioutil.ReadAll(req.Body)
		if err != nil {
			apierror.GenerateError("Trouble getting JSON request body for updating contact", err, rw, req)
		}

		err = json.Unmarshal(requestBody, &c)
		if err != nil {
			apierror.GenerateError("Trouble unmarshalling JSON request body for updating contact", err, rw, req)
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
		if req.FormValue("brandID") != "" {
			c.Brand.ID, err = strconv.Atoi(req.FormValue("brandID"))
		}
	}
	if err = c.Update(); err != nil {
		apierror.GenerateError("Trouble updating contact", err, rw, req)
	}
	return encoding.Must(enc.Encode(c))
}

func DeleteContact(rw http.ResponseWriter, req *http.Request, params martini.Params, enc encoding.Encoder) string {
	var err error
	var c contact.Contact

	if c.ID, err = strconv.Atoi(params["id"]); err != nil {
		apierror.GenerateError("Trouble getting contact ID", err, rw, req)
	}

	if err = c.Delete(); err != nil {
		apierror.GenerateError("Trouble deleting contact", err, rw, req)
	}

	return encoding.Must(enc.Encode(c))
}
