package contact

import (
	"bytes"
	"encoding/json"
	"github.com/curt-labs/GoAPI/helpers/email"
	"github.com/curt-labs/GoAPI/helpers/encoding"
	"github.com/curt-labs/GoAPI/models/contact"
	"github.com/go-martini/martini"
	"io/ioutil"
	"net/http"
	"strconv"
	"text/template"
)

type SmtpTemplateData struct {
	From    string
	To      []string
	Subject string
	Body    string
}

const emailTemplate = "From: {{.From}}\nTo: {{.To}}}\nSubject: {{.Subject}}\n{{.Body}}\n"

func SendEmail(rw http.ResponseWriter, req *http.Request, enc encoding.Encoder, params martini.Params) string {
	// set up recipients
	var ct contact.ContactType
	var tos []string
	var err error
	contactTypeID := params["id"] //will generate to's from contactTypeID
	ct.ID, err = strconv.Atoi(contactTypeID)
	receivers, err := ct.GetReceivers()
	for _, r := range receivers {
		tos = append(tos, r.Email)
	}

	subject := "Email from the Aries Contact Form!"

	//form contact data from requestBody
	requestBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return encoding.Must(enc.Encode(false))
	}
	var c contact.Contact
	err = json.Unmarshal(requestBody, &c)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return encoding.Must(enc.Encode(false))
	}

	//create body from contact object
	body :=
		"Email: " + c.Email + "\n" +
			"Phone: " + c.Phone + "\n" +
			"Subject: " + c.Subject + "\n" +
			"Time: " + c.Created.String() + "\n" +
			"Type: " + c.Type + "\n" +
			"Address1: " + c.Address1 + "\n" +
			"Address2: " + c.Address2 + "\n" +
			"City: " + c.City + "\n" +
			"State: " + c.State + "\n" +
			"PostalCode: " + c.PostalCode + "\n" +
			"Country: " + c.Country + "\n\n" +
			"Message: " + c.Message + "\n"

	//set up template
	t := template.New("emailTemplate")
	t, err = t.Parse(emailTemplate)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return encoding.Must(enc.Encode(false))
	}
	var msg bytes.Buffer
	context := &SmtpTemplateData{
		(c.FirstName + " " + c.LastName),
		tos,
		subject,
		body,
	}
	err = t.Execute(&msg, context)

	err = email.Send(tos, subject, msg.String(), false)

	return "Email sent"
}
