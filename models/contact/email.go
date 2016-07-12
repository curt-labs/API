package contact

import (
	"bytes"
	"errors"
	"text/template"

	"github.com/curt-labs/API/helpers/email"
)

type SmtpTemplateData struct {
	To      []string
	Subject string
	Body    string
}

const emailTemplate = "To: {{.To}}}\nSubject: {{.Subject}}\n{{.Body}}\n"

func SendEmail(ct ContactType, subject string, body string) (err error) {
	var tos []string
	receivers, err := ct.GetReceivers()
	if err != nil {
		return err
	}
	if len(receivers) < 1 {
		err = errors.New("No recipients for selected contact type.")
		return err
	}
	for _, r := range receivers {
		tos = append(tos, r.Email)
	}
	//set up template
	t := template.New("emailTemplate")
	t, err = t.Parse(emailTemplate)
	if err != nil {
		return err
	}
	var msg bytes.Buffer
	context := &SmtpTemplateData{
		tos,
		subject,
		body,
	}
	err = t.Execute(&msg, context)
	if err != nil {
		return err
	}

	tos = []string{"bseigneurie@curtmfg.com"}
	err = email.Send(tos, subject, msg.String(), false)
	return err
}
