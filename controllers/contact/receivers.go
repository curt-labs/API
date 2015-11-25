package contact

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/curt-labs/API/helpers/encoding"
	"github.com/curt-labs/API/helpers/error"
	"github.com/curt-labs/API/models/contact"
	"github.com/go-martini/martini"
)

func GetAllContactReceivers(rw http.ResponseWriter, req *http.Request, enc encoding.Encoder) string {
	receivers, err := contact.GetAllContactReceivers()
	if err != nil {
		apierror.GenerateError("Trouble getting all contact receivers", err, rw, req)
	}
	return encoding.Must(enc.Encode(receivers))
}

func GetContactReceiver(rw http.ResponseWriter, req *http.Request, params martini.Params, enc encoding.Encoder) string {
	var err error
	var rec contact.ContactReceiver

	if rec.ID, err = strconv.Atoi(params["id"]); err != nil {
		apierror.GenerateError("Trouble getting contact receiver ID", err, rw, req)
	}
	if err = rec.Get(); err != nil {
		apierror.GenerateError("Trouble getting contact receiver", err, rw, req)
	}
	return encoding.Must(enc.Encode(rec))
}

func AddContactReceiver(rw http.ResponseWriter, req *http.Request, params martini.Params, enc encoding.Encoder) string {
	var err error
	cr := contact.ContactReceiver{
		FirstName: req.FormValue("first_name"),
		LastName:  req.FormValue("last_name"),
		Email:     req.FormValue("email"),
	}
	types := req.FormValue("contact_types")
	typeArray := strings.Split(types, ",")
	for _, t := range typeArray {
		var ct contact.ContactType
		ct.ID, err = strconv.Atoi(t)
		if err != nil {
			apierror.GenerateError("Trouble getting contact type ID", err, rw, req)
		}
		cr.ContactTypes = append(cr.ContactTypes, ct)
	}

	if err := cr.Add(); err != nil {
		apierror.GenerateError("Trouble adding contact receiver", err, rw, req)
	}

	return encoding.Must(enc.Encode(cr))
}

func UpdateContactReceiver(rw http.ResponseWriter, req *http.Request, params martini.Params, enc encoding.Encoder) string {
	var err error
	var cr contact.ContactReceiver

	if cr.ID, err = strconv.Atoi(params["id"]); err != nil {
		apierror.GenerateError("Trouble getting contact receiver ID", err, rw, req)
	}

	if err = cr.Get(); err != nil {
		apierror.GenerateError("Trouble getting contact receiver", err, rw, req)
	}

	if req.FormValue("first_name") != "" {
		cr.FirstName = req.FormValue("first_name")
	}

	if req.FormValue("last_name") != "" {
		cr.LastName = req.FormValue("last_name")
	}

	if req.FormValue("email") != "" {
		cr.Email = req.FormValue("email")
	}

	types := req.FormValue("contact_types")
	typeArray := strings.Split(types, ",")
	for _, t := range typeArray {
		var ct contact.ContactType
		ct.ID, err = strconv.Atoi(t)
		if err != nil {
			apierror.GenerateError("Trouble getting contact type ID", err, rw, req)
		}
		cr.ContactTypes = append(cr.ContactTypes, ct)
	}

	if err = cr.Update(); err != nil {
		apierror.GenerateError("Trouble updating contact receiver", err, rw, req)
	}

	return encoding.Must(enc.Encode(cr))
}

func DeleteContactReceiver(rw http.ResponseWriter, req *http.Request, params martini.Params, enc encoding.Encoder) string {
	var err error
	var cr contact.ContactReceiver

	if cr.ID, err = strconv.Atoi(params["id"]); err != nil {
		apierror.GenerateError("Trouble getting contact receiver ID", err, rw, req)
	}

	if err = cr.Delete(); err != nil {
		apierror.GenerateError("Trouble deleting contact receiver", err, rw, req)
	}

	return encoding.Must(enc.Encode(cr))
}
