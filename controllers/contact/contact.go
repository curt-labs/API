package contact

import (
	"net/http"
	"strconv"

	"github.com/curt-labs/GoAPI/helpers/encoding"
	"github.com/curt-labs/GoAPI/models/contact"
	"github.com/go-martini/martini"
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

func GetAllContactTypes(rw http.ResponseWriter, req *http.Request, enc encoding.Encoder) string {
	types, err := contact.GetAllContactTypes()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return err.Error()
	}
	return encoding.Must(enc.Encode(types))
}

func GetAllContactReceivers(rw http.ResponseWriter, req *http.Request, enc encoding.Encoder) string {
	receivers, err := contact.GetAllContactReceivers()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return err.Error()
	}
	return encoding.Must(enc.Encode(receivers))
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

func GetContactType(rw http.ResponseWriter, req *http.Request, params martini.Params, enc encoding.Encoder) string {
	var err error
	var ctype contact.ContactType

	if ctype.ID, err = strconv.Atoi(params["id"]); err != nil {
		http.Error(rw, "Invalid ContactType ID", http.StatusInternalServerError)
		return "Invalid ContactType ID"
	}
	if err = ctype.Get(); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return err.Error()
	}
	return encoding.Must(enc.Encode(ctype))
}

func GetContactReceiver(rw http.ResponseWriter, req *http.Request, params martini.Params, enc encoding.Encoder) string {
	var err error
	var rec contact.ContactReceiver

	if rec.ID, err = strconv.Atoi(params["id"]); err != nil {
		http.Error(rw, "Invalid ContactReceiver ID", http.StatusInternalServerError)
		return "Invalid ContactReceiver ID"
	}
	if err = rec.Get(); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return err.Error()
	}
	return encoding.Must(enc.Encode(rec))
}
