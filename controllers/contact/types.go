package contact

import (
	"net/http"
	"strconv"

	"github.com/curt-labs/GoAPI/helpers/apicontext"
	"github.com/curt-labs/GoAPI/helpers/encoding"
	"github.com/curt-labs/GoAPI/helpers/error"
	"github.com/curt-labs/GoAPI/models/contact"
	"github.com/go-martini/martini"
)

func GetAllContactTypes(rw http.ResponseWriter, req *http.Request, enc encoding.Encoder, dtx *apicontext.DataContext) string {
	types, err := contact.GetAllContactTypes(dtx)
	if err != nil {
		apierror.GenerateError("Trouble getting all contact types", err, rw, req)
	}
	return encoding.Must(enc.Encode(types))
}

func GetContactType(rw http.ResponseWriter, req *http.Request, params martini.Params, enc encoding.Encoder) string {
	var err error
	var ctype contact.ContactType

	if ctype.ID, err = strconv.Atoi(params["id"]); err != nil {
		apierror.GenerateError("Trouble getting contact type ID", err, rw, req)
	}
	if err = ctype.Get(); err != nil {
		apierror.GenerateError("Trouble getting contact type", err, rw, req)
	}
	return encoding.Must(enc.Encode(ctype))
}

//Get receivers of a certain contact type
func GetReceiversByContactType(rw http.ResponseWriter, req *http.Request, params martini.Params, enc encoding.Encoder) string {
	var err error
	var ctype contact.ContactType
	var crs contact.ContactReceivers

	if ctype.ID, err = strconv.Atoi(params["id"]); err != nil {
		apierror.GenerateError("Trouble getting contact type ID", err, rw, req)
	}
	crs, err = ctype.GetReceivers()
	if err != nil {
		apierror.GenerateError("Trouble getting receivers for contact type", err, rw, req)
	}
	return encoding.Must(enc.Encode(crs))
}

func AddContactType(rw http.ResponseWriter, req *http.Request, params martini.Params, enc encoding.Encoder) string {
	ct := contact.ContactType{
		Name: req.FormValue("name"),
	}
	if req.FormValue("show") != "" {
		if show, err := strconv.ParseBool(req.FormValue("show")); err == nil {
			ct.ShowOnWebsite = show
		}
	}

	if req.FormValue("brandId") != "" {
		if brandId, err := strconv.Atoi(req.FormValue("brandId")); err == nil {
			ct.BrandID = brandId
		}
	}

	if err := ct.Add(); err != nil {
		apierror.GenerateError("Trouble adding contact type", err, rw, req)
	}

	return encoding.Must(enc.Encode(ct))
}

func UpdateContactType(rw http.ResponseWriter, req *http.Request, params martini.Params, enc encoding.Encoder) string {
	var err error
	var ct contact.ContactType

	if ct.ID, err = strconv.Atoi(params["id"]); err != nil {
		apierror.GenerateError("Trouble getting contact type ID", err, rw, req)
	}

	//json?

	if err = ct.Get(); err != nil {
		apierror.GenerateError("Trouble getting contact type", err, rw, req)
	}

	if req.FormValue("name") != "" {
		ct.Name = req.FormValue("name")
	}

	if req.FormValue("show") != "" {
		if show, err := strconv.ParseBool(req.FormValue("show")); err == nil {
			ct.ShowOnWebsite = show
		}
	}
	if req.FormValue("brandId") != "" {
		if brandId, err := strconv.Atoi(req.FormValue("brandId")); err == nil {
			ct.BrandID = brandId
		}
	}

	if err = ct.Update(); err != nil {
		apierror.GenerateError("Trouble updating contact type", err, rw, req)
	}

	return encoding.Must(enc.Encode(ct))
}

func DeleteContactType(rw http.ResponseWriter, req *http.Request, params martini.Params, enc encoding.Encoder) string {
	var err error
	var ct contact.ContactType

	if ct.ID, err = strconv.Atoi(params["id"]); err != nil {
		apierror.GenerateError("Trouble getting contact type ID", err, rw, req)
	}

	if err = ct.Delete(); err != nil {
		apierror.GenerateError("Trouble deleting contact type ID", err, rw, req)
	}

	return encoding.Must(enc.Encode(ct))
}
