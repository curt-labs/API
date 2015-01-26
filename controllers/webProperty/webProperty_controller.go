package webProperty_controller

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/curt-labs/GoAPI/helpers/apicontext"
	"github.com/curt-labs/GoAPI/helpers/encoding"
	"github.com/curt-labs/GoAPI/helpers/error"
	"github.com/curt-labs/GoAPI/helpers/sortutil"
	"github.com/curt-labs/GoAPI/models/customer"
	"github.com/curt-labs/GoAPI/models/webProperty"
	"github.com/go-martini/martini"
)

const (
	timeFormat = "2006-01-02 15:04:05"
)

func GetAll(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder, dtx *apicontext.DataContext) string {
	props, err := webProperty_model.GetAll(dtx)
	if err != nil {
		apierror.GenerateError("Trouble getting all web properties", err, rw, r)
		return ""
	}
	sort := r.FormValue("sort")
	direction := r.FormValue("direction")
	if sort != "" {
		if strings.ContainsAny(direction, "esc") {
			sortutil.DescByField(props, sort)
		} else {
			sortutil.AscByField(props, sort)
		}
	}
	return encoding.Must(enc.Encode(props))
}

func Get(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params, dtx *apicontext.DataContext) string {
	var w webProperty_model.WebProperty
	var err error

	idStr := r.FormValue("id")
	if idStr == "" {
		idStr = params["id"]
	}

	if w.ID, err = strconv.Atoi(idStr); err != nil {
		apierror.GenerateError("Trouble getting web property ID", err, rw, r)
		return ""
	}

	if err = w.Get(dtx); err != nil {
		apierror.GenerateError("Trouble getting web property", err, rw, r)
		return ""
	}

	sort := r.FormValue("sort")
	direction := r.FormValue("direction")
	if sort != "" {
		if strings.ContainsAny(direction, "esc") {
			sortutil.DescByField(w, sort)
		} else {
			sortutil.AscByField(w, sort)
		}

	}
	return encoding.Must(enc.Encode(w))
}

func GetByPrivateKey(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params, dtx *apicontext.DataContext) string {
	var err error
	var webProperties webProperty_model.WebProperties

	privateKey := r.FormValue("key")
	cust := customer.Customer{}

	if err = cust.GetCustomerIdFromKey(privateKey); err != nil {
		apierror.GenerateError("Trouble getting customer", err, rw, r)
		return ""
	}

	if webProperties, err = webProperty_model.GetByCustomer(cust.Id, dtx); err != nil {
		apierror.GenerateError("Trouble getting web property", err, rw, r)
		return ""
	}

	return encoding.Must(enc.Encode(webProperties))
}

func GetAllTypes(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder, dtx *apicontext.DataContext) string {
	props, err := webProperty_model.GetAllWebPropertyTypes(dtx)
	if err != nil {
		apierror.GenerateError("Trouble getting web property types", err, rw, r)
		return ""
	}

	sort := r.FormValue("sort")
	direction := r.FormValue("direction")
	if sort != "" {
		if strings.ContainsAny(direction, "esc") {
			sortutil.DescByField(props, sort)
		} else {
			sortutil.AscByField(props, sort)
		}

	}
	return encoding.Must(enc.Encode(props))
}
func GetAllNotes(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder, dtx *apicontext.DataContext) string {
	props, err := webProperty_model.GetAllWebPropertyNotes(dtx)
	if err != nil {
		apierror.GenerateError("Trouble getting all web property notes", err, rw, r)
		return ""
	}

	sort := r.FormValue("sort")
	direction := r.FormValue("direction")
	if sort != "" {
		if strings.ContainsAny(direction, "esc") {
			sortutil.DescByField(props, sort)
		} else {
			sortutil.AscByField(props, sort)
		}

	}
	return encoding.Must(enc.Encode(props))
}
func GetAllRequirements(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder, dtx *apicontext.DataContext) string {
	props, err := webProperty_model.GetAllWebPropertyRequirements(dtx)
	if err != nil {
		apierror.GenerateError("Trouble getting all web property requirements", err, rw, r)
		return ""
	}

	sort := r.FormValue("sort")
	direction := r.FormValue("direction")
	if sort != "" {
		if strings.ContainsAny(direction, "esc") {
			sortutil.DescByField(props, sort)
		} else {
			sortutil.AscByField(props, sort)
		}

	}
	return encoding.Must(enc.Encode(props))
}

func CreateUpdateWebProperty(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params, dtx *apicontext.DataContext) string {
	var w webProperty_model.WebProperty
	var err error

	cust := customer.Customer{}

	if err = cust.GetCustomerIdFromKey(dtx.APIKey); err != nil {
		apierror.GenerateError("Trouble getting customer ID from API Key", err, rw, r)
		return ""
	}

	w.CustID = cust.Id

	//determine content type
	contType := r.Header.Get("Content-Type")
	if contType == "application/json" {
		//json
		requestBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			apierror.GenerateError("Trouble reading request body while saving web property", err, rw, r)
			return ""
		}

		if err = json.Unmarshal(requestBody, &w); err != nil {
			apierror.GenerateError("Trouble unmarshalling request body while saving web property", err, rw, r)
			return ""
		}
	} else {
		if r.FormValue("id") != "" || params["id"] != "" {
			idStr := r.FormValue("id")
			if idStr == "" {
				idStr = params["id"]
			}

			if w.ID, err = strconv.Atoi(idStr); err != nil {
				apierror.GenerateError("Trouble getting web property ID", err, rw, r)
				return ""
			}

			if err = w.Get(dtx); err != nil {
				apierror.GenerateError("Trouble getting web property", err, rw, r)
				return ""
			}
		}

		name := r.FormValue("name")
		url := r.FormValue("url")
		isEnabled := r.FormValue("isEnabled")
		sellerID := r.FormValue("sellerID")
		webPropertyTypeID := r.FormValue("webPropertyTypeID")
		isFinalApproved := r.FormValue("isApproved")
		enabledDate := r.FormValue("enabledDate")
		isDenied := r.FormValue("isDenied")
		requestedDate := r.FormValue("requestedDate")
		typeID := r.FormValue("typeID")

		if name != "" {
			w.Name = name
		}
		if url != "" {
			w.Url = url
		}
		if isEnabled != "" {
			if w.IsEnabled, err = strconv.ParseBool(isEnabled); err != nil {
				apierror.GenerateError("Trouble parsing boolean value for webproperty.isEnabled", err, rw, r)
				return ""
			}
		}
		if sellerID != "" {
			w.SellerID = sellerID
		}
		if webPropertyTypeID != "" {
			if w.WebPropertyType.ID, err = strconv.Atoi(webPropertyTypeID); err != nil {
				apierror.GenerateError("Trouble getting web property type ID", err, rw, r)
				return ""
			}
		}
		if isFinalApproved != "" {
			if w.IsFinalApproved, err = strconv.ParseBool(isFinalApproved); err != nil {
				apierror.GenerateError("Trouble parsing boolean value for webproperty.isApproved", err, rw, r)
				return ""
			}
		}
		if enabledDate != "" {
			en, err := time.Parse(timeFormat, enabledDate)
			if err != nil {
				apierror.GenerateError("Trouble parsing date for webproperty.enabledDate", err, rw, r)
				return ""
			}
			w.IsEnabledDate = &en
		}
		if isDenied != "" {
			if w.IsDenied, err = strconv.ParseBool(isDenied); err != nil {
				apierror.GenerateError("Trouble parsing boolean value for webproperty.isDenied", err, rw, r)
				return ""
			}
		}
		if requestedDate != "" {
			req, err := time.Parse(timeFormat, requestedDate)
			if err != nil {
				apierror.GenerateError("Trouble parsing boolean value for webproperty.requestedDate", err, rw, r)
				return ""
			}
			w.RequestedDate = &req
		}
		if typeID != "" {
			if w.WebPropertyType.ID, err = strconv.Atoi(typeID); err != nil {
				apierror.GenerateError("Trouble getting type ID", err, rw, r)
				return ""
			}
		}
	}

	if w.ID > 0 {
		err = w.Update()
	} else {
		//notes (text) and property requirements (reqID) can be created when web property is created
		notes := r.Form["note"]
		for _, v := range notes {
			var n webProperty_model.WebPropertyNote
			n.Text = v
			w.WebPropertyNotes = append(w.WebPropertyNotes, n)
		}
		reqIDs := r.Form["requirement"]
		for _, v := range reqIDs {
			var n webProperty_model.WebPropertyRequirement
			n.RequirementID, err = strconv.Atoi(v)
			w.WebPropertyRequirements = append(w.WebPropertyRequirements, n)
		}
		err = w.Create()
	}

	if err != nil {
		msg := "Trouble creating web property"
		if w.ID > 0 {
			msg = "Trouble updating web property"
		}
		apierror.GenerateError(msg, err, rw, r)
		return ""
	}
	return encoding.Must(enc.Encode(w))
}

func DeleteWebProperty(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	var w webProperty_model.WebProperty
	var err error

	idStr := r.FormValue("id")
	if idStr == "" {
		idStr = params["id"]
	}

	if w.ID, err = strconv.Atoi(idStr); err != nil {
		apierror.GenerateError("Trouble getting web property ID", err, rw, r)
		return ""
	}

	if err = w.Delete(); err != nil {
		apierror.GenerateError("Trouble deleting web property", err, rw, r)
		return ""
	}

	return encoding.Must(enc.Encode(w))
}

func GetWebPropertyNote(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	var wn webProperty_model.WebPropertyNote
	var err error

	idStr := r.FormValue("id")
	if idStr == "" {
		idStr = params["id"]
	}

	if wn.ID, err = strconv.Atoi(idStr); err != nil {
		apierror.GenerateError("Trouble getting web property note ID", err, rw, r)
		return ""
	}

	if err = wn.Get(); err != nil {
		apierror.GenerateError("Trouble getting web property note", err, rw, r)
		return ""
	}

	return encoding.Must(enc.Encode(wn))
}

func CreateUpdateWebPropertyNote(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	var w webProperty_model.WebProperty
	var n webProperty_model.WebPropertyNote
	var err error

	contType := r.Header.Get("Content-Type")
	if contType == "application/json" {
		//json
		requestBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			apierror.GenerateError("Trouble reading request body while saving web property note", err, rw, r)
			return ""
		}

		if err = json.Unmarshal(requestBody, &n); err != nil {
			apierror.GenerateError("Trouble unmarshalling request body while saving web property note", err, rw, r)
			return ""
		}
	} else {

		if r.FormValue("id") != "" || params["id"] != "" {
			idStr := r.FormValue("id")
			if idStr == "" {
				idStr = params["id"]
			}

			if n.ID, err = strconv.Atoi(idStr); err != nil {
				apierror.GenerateError("Trouble getting web property ID", err, rw, r)
				return ""
			}

			if err = n.Get(); err != nil {
				apierror.GenerateError("Trouble getting web property", err, rw, r)
				return ""
			}
		}

		text := r.FormValue("text")
		if text != "" {
			n.Text = text
		}
		propID := r.FormValue("propertyID")
		if propID != "" {
			w.ID, err = strconv.Atoi(propID)
			n.WebPropID = w.ID
		}
	}
	if n.ID > 0 {
		err = n.Update()
	} else {
		err = n.Create()
	}

	if err != nil {
		msg := "Trouble creating web property note"
		if n.ID > 0 {
			msg = "Trouble updating web property note"
		}
		apierror.GenerateError(msg, err, rw, r)
		return ""
	}

	return encoding.Must(enc.Encode(n))
}

func DeleteWebPropertyNote(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	var n webProperty_model.WebPropertyNote
	var err error

	idStr := r.FormValue("id")
	if idStr == "" {
		idStr = params["id"]
	}

	if n.ID, err = strconv.Atoi(idStr); err != nil {
		apierror.GenerateError("Trouble getting web property note ID", err, rw, r)
		return ""
	}

	if err = n.Delete(); err != nil {
		apierror.GenerateError("Trouble deleting web property note", err, rw, r)
		return ""
	}

	return encoding.Must(enc.Encode(n))
}

func GetWebPropertyRequirement(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	var wr webProperty_model.WebPropertyRequirement
	var err error

	idStr := r.FormValue("id")
	if idStr == "" {
		idStr = params["id"]
	}

	if wr.RequirementID, err = strconv.Atoi(idStr); err != nil {
		apierror.GenerateError("Trouble getting web property requirement ID", err, rw, r)
		return ""
	}

	if err = wr.Get(); err != nil {
		apierror.GenerateError("Trouble getting web property requirement", err, rw, r)
		return ""
	}

	return encoding.Must(enc.Encode(wr))
}

func CreateUpdateWebPropertyRequirement(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	var wr webProperty_model.WebPropertyRequirement
	var err error

	contType := r.Header.Get("Content-Type")
	if contType == "application/json" {
		//json
		requestBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			apierror.GenerateError("Trouble reading request body while saving web property requirement", err, rw, r)
			return ""
		}

		if err = json.Unmarshal(requestBody, &wr); err != nil {
			apierror.GenerateError("Trouble unmarshalling request body while saving web property requirement", err, rw, r)
			return ""
		}
	} else {

		if r.FormValue("id") != "" || params["id"] != "" {
			idStr := r.FormValue("id")
			if idStr == "" {
				idStr = params["id"]
			}

			if wr.RequirementID, err = strconv.Atoi(idStr); err != nil {
				apierror.GenerateError("Trouble getting web property requirement ID", err, rw, r)
				return ""
			}

			if err = wr.Get(); err != nil {
				apierror.GenerateError("Trouble getting web property requirement", err, rw, r)
				return ""
			}
		}

		reqType := r.FormValue("reqType")
		requirement := r.FormValue("requirement")

		if reqType != "" {
			wr.ReqType = reqType
		}

		if requirement != "" {
			wr.Requirement = requirement
		}
	}

	if wr.RequirementID > 0 {
		err = wr.Update()
	} else {
		err = wr.Create()
	}

	if err != nil {
		msg := "Trouble creating web property requirement"
		if wr.RequirementID > 0 {
			msg = "Trouble updating web property requirement"
		}
		apierror.GenerateError(msg, err, rw, r)
		return ""
	}

	return encoding.Must(enc.Encode(wr))
}

func DeleteWebPropertyRequirement(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	var wr webProperty_model.WebPropertyRequirement
	var err error

	idStr := r.FormValue("id")
	if idStr == "" {
		idStr = params["id"]
	}

	if wr.RequirementID, err = strconv.Atoi(idStr); err != nil {
		apierror.GenerateError("Trouble getting web property requirement ID", err, rw, r)
		return ""
	}

	if err = wr.Delete(); err != nil {
		log.Print(err)
		apierror.GenerateError("Trouble deleting web property requirement", err, rw, r)
		return ""
	}

	return encoding.Must(enc.Encode(wr))
}

func GetWebPropertyType(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	var t webProperty_model.WebPropertyType
	var err error

	idStr := r.FormValue("id")
	if idStr == "" {
		idStr = params["id"]
	}

	if t.ID, err = strconv.Atoi(idStr); err != nil {
		apierror.GenerateError("Trouble getting web property type ID", err, rw, r)
		return ""
	}

	if err = t.Get(); err != nil {
		apierror.GenerateError("Trouble getting web property type", err, rw, r)
		return ""
	}

	return encoding.Must(enc.Encode(t))
}

func CreateUpdateWebPropertyType(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	var t webProperty_model.WebPropertyType
	var err error

	contType := r.Header.Get("Content-Type")
	if contType == "application/json" {
		//json
		requestBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			apierror.GenerateError("Trouble reading request body while saving web property type", err, rw, r)
			return ""
		}

		if err = json.Unmarshal(requestBody, &t); err != nil {
			apierror.GenerateError("Trouble unmarshalling request body while saving web property type", err, rw, r)
			return ""
		}
	} else {

		if r.FormValue("id") != "" || params["id"] != "" {
			idStr := r.FormValue("id")
			if idStr == "" {
				idStr = params["id"]
			}

			if t.ID, err = strconv.Atoi(idStr); err != nil {
				apierror.GenerateError("Trouble getting web property type ID", err, rw, r)
				return ""
			}

			if err = t.Get(); err != nil {
				apierror.GenerateError("Trouble getting web property type", err, rw, r)
				return ""
			}
		}

		typeID := r.FormValue("typeID")
		theType := r.FormValue("type")

		if typeID != "" {
			t.TypeID, err = strconv.Atoi(typeID)
		}

		if theType != "" {
			t.Type = theType
		}
	}

	if t.ID > 0 {
		err = t.Update()
	} else {
		err = t.Create()
	}

	if err != nil {
		msg := "Trouble creating web property type"
		if t.ID > 0 {
			msg = "Trboule updating web property type"
		}
		apierror.GenerateError(msg, err, rw, r)
		return ""
	}
	return encoding.Must(enc.Encode(t))
}

func DeleteWebPropertyType(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	var t webProperty_model.WebPropertyType
	var err error

	idStr := r.FormValue("id")
	if idStr == "" {
		idStr = params["id"]
	}

	if t.ID, err = strconv.Atoi(idStr); err != nil {
		apierror.GenerateError("Trouble getting web property type ID", err, rw, r)
		return ""
	}

	if err = t.Delete(); err != nil {
		apierror.GenerateError("Trouble deleting web property type", err, rw, r)
		return ""
	}

	return encoding.Must(enc.Encode(t))
}

func Search(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder) string {
	var err error

	name := r.FormValue("name")
	custID := r.FormValue("custID")
	badgeID := r.FormValue("badgeID")
	url := r.FormValue("url")
	isEnabled := r.FormValue("isEnabled")
	sellerID := r.FormValue("sellerID")
	webPropertyTypeID := r.FormValue("webPropertyTypeID")
	isFinalApproved := r.FormValue("isApproved")
	enabledDate := r.FormValue("enabledDate")
	isDenied := r.FormValue("isDenied")
	requestedDate := r.FormValue("requestedDate")
	typeID := r.FormValue("typeID")

	page := r.FormValue("page")
	results := r.FormValue("results")

	l, err := webProperty_model.Search(name, custID, badgeID, url, isEnabled, sellerID, webPropertyTypeID, isFinalApproved, enabledDate, isDenied, requestedDate, typeID, page, results)
	if err != nil {
		apierror.GenerateError("Trouble searching for web properties", err, rw, r)
		return ""
	}

	return encoding.Must(enc.Encode(l))
}
