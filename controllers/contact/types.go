package contact

import (
	"net/http"
	"strconv"

	"github.com/curt-labs/GoAPI/helpers/encoding"
	"github.com/curt-labs/GoAPI/models/contact"
	"github.com/go-martini/martini"
)

func GetAllContactTypes(rw http.ResponseWriter, req *http.Request, enc encoding.Encoder) string {
	types, err := contact.GetAllContactTypes()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return err.Error()
	}
	return encoding.Must(enc.Encode(types))
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

func AddContactType(rw http.ResponseWriter, req *http.Request, params martini.Params, enc encoding.Encoder) string {
	ct := contact.ContactType{
		Name: req.FormValue("name"),
	}

	if err := ct.Add(); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return err.Error()
	}

	return encoding.Must(enc.Encode(ct))
}

func UpdateContactType(rw http.ResponseWriter, req *http.Request, params martini.Params, enc encoding.Encoder) string {
	var err error
	var ct contact.ContactType

	if ct.ID, err = strconv.Atoi(params["id"]); err != nil {
		http.Error(rw, "Invalid ContactType ID", http.StatusInternalServerError)
		return "Invalid ContactType ID"
	}

	if err = ct.Get(); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return err.Error()
	}

	if req.FormValue("name") != "" {
		ct.Name = req.FormValue("name")
	}

	if req.FormValue("show") != "" {
		if show, err := strconv.ParseBool(req.FormValue("show")); err == nil {
			ct.ShowOnWebsite = show
		}
	}

	if err = ct.Update(); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return err.Error()
	}

	return encoding.Must(enc.Encode(ct))
}

func DeleteContactType(rw http.ResponseWriter, req *http.Request, params martini.Params, enc encoding.Encoder) string {
	var err error
	var ct contact.ContactType

	if ct.ID, err = strconv.Atoi(params["id"]); err != nil {
		http.Error(rw, "Invalid ContactType ID", http.StatusInternalServerError)
		return "Invalid ContactType ID"
	}

	if err = ct.Delete(); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return err.Error()
	}

	return encoding.Must(enc.Encode(ct))
}
