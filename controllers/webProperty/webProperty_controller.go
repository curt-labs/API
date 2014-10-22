package webProperty_controller

import (
	"encoding/json"
	// "errors"
	"github.com/curt-labs/GoAPI/helpers/encoding"
	"github.com/curt-labs/GoAPI/helpers/sortutil"
	"github.com/curt-labs/GoAPI/models/customer"
	"github.com/curt-labs/GoAPI/models/webProperty"
	"github.com/go-martini/martini"
	"io/ioutil"
	// "log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	timeFormat = "2006-01-02 15:04:05"
)

func GetAll(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder) string {
	props, err := webProperty_model.GetAll()
	if err != nil {
		return err.Error()
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

func Get(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	var w webProperty_model.WebProperty
	var err error
	idStr := r.FormValue("id")
	if idStr != "" {
		w.ID, err = strconv.Atoi(idStr)
		if err != nil {
			return err.Error()
		}
	} else {
		w.ID, err = strconv.Atoi(params["id"])
		if err != nil {
			return err.Error()
		}
	}

	err = w.Get()
	if err != nil {
		return err.Error()
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

func GetByPrivateKey(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	var w webProperty_model.WebProperty
	var err error
	//get private key & custID
	privateKey := r.FormValue("key")
	custID, err := customer.GetCustomerIdFromKey(privateKey)
	w.CustID = custID
	err = w.GetByCust()
	if err != nil {
		return err
	}
	return encoding.Must(enc.Encode(w))
}

func GetAllTypes(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder) string {
	props, err := webProperty_model.GetAllWebPropertyTypes()
	if err != nil {
		return err.Error()
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
func GetAllNotes(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder) string {
	props, err := webProperty_model.GetAllWebPropertyNotes()
	if err != nil {
		return err.Error()
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
func GetAllRequirements(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder) string {
	props, err := webProperty_model.GetAllWebPropertyRequirements()
	if err != nil {
		return err.Error()
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

func CreateUpdateWebProperty(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	var w webProperty_model.WebProperty
	var err error
	//get private key & custID
	privateKey := r.FormValue("key")
	custID, err := customer.GetCustomerIdFromKey(privateKey)
	w.CustID = custID

	//determine content type
	contType := r.Header.Get("Content-Type")
	if contType == "application/json" {
		//json
		requestBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return encoding.Must(enc.Encode(false))
		}

		err = json.Unmarshal(requestBody, &w)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return encoding.Must(enc.Encode(false))
		}
	} else {
		idStr := r.FormValue("id")
		if idStr != "" {
			w.ID, err = strconv.Atoi(idStr)
			if err != nil {
				return err.Error()
			}
		}
		if params["id"] != "" {
			w.ID, err = strconv.Atoi(params["id"])
			if err != nil {
				return err.Error()
			}
		}
		if w.ID > 0 {
			w.Get()
		}

		name := r.FormValue("name")
		url := r.FormValue("url")
		isEnabled := r.FormValue("isEnabled")
		sellerID := r.FormValue("sellerID")
		webPropertyTypeID := r.FormValue("webPropertyTypeID")
		isFinalApproved := r.FormValue("metaDescription")
		isEnabledDate := r.FormValue("metaDescription")
		isDenied := r.FormValue("metaDescription")
		requestedDate := r.FormValue("metaDescription")
		typeID := r.FormValue("typeID")

		if name != "" {
			w.Name = name
		}
		if url != "" {
			w.Url = url
		}
		if isEnabled != "" {
			w.IsEnabled, err = strconv.ParseBool(isEnabled)
		}
		if sellerID != "" {
			w.SellerID = sellerID
		}
		if webPropertyTypeID != "" {
			w.WebPropertyType.ID, err = strconv.Atoi(webPropertyTypeID)
		}
		if isFinalApproved != "" {
			w.IsFinalApproved, err = strconv.ParseBool(isFinalApproved)
		}
		if isEnabledDate != "" {
			w.IsEnabledDate, err = time.Parse(timeFormat, isEnabledDate)
		}
		if isDenied != "" {
			w.IsDenied, err = strconv.ParseBool(isDenied)
		}
		if requestedDate != "" {
			w.RequestedDate, err = time.Parse(timeFormat, requestedDate)
		}
		if typeID != "" {
			w.WebPropertyType.ID, err = strconv.Atoi(typeID)
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
		return err.Error()
	}
	return encoding.Must(enc.Encode(w))
}

func DeleteWebProperty(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	var p webProperty_model.WebProperty
	var err error

	idStr := r.FormValue("id")
	if idStr != "" {
		p.ID, err = strconv.Atoi(idStr)
		if err != nil {
			return err.Error()
		}
	} else {
		p.ID, err = strconv.Atoi(params["id"])
		if err != nil {
			return err.Error()
		}
	}
	err = p.Delete()
	if err != nil {
		return err.Error()
	}
	return encoding.Must(enc.Encode(p))
}

func GetWebPropertyNote(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	var n webProperty_model.WebPropertyNote
	var err error

	idStr := r.FormValue("id")
	if idStr != "" {
		n.ID, err = strconv.Atoi(idStr)
		if err != nil {
			return err.Error()
		}
	} else {
		n.ID, err = strconv.Atoi(params["id"])
		if err != nil {
			return err.Error()
		}
	}
	err = n.Get()
	if err != nil {
		return err.Error()
	}
	return encoding.Must(enc.Encode(n))
}

func CreateUpdateWebPropertyNote(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	var w webProperty_model.WebProperty
	var n webProperty_model.WebPropertyNote
	var err error

	idStr := r.FormValue("id")
	if idStr != "" {
		n.ID, err = strconv.Atoi(idStr)
		if err != nil {
			return err.Error()
		}
	}
	if params["id"] != "" {
		n.ID, err = strconv.Atoi(params["id"])
		if err != nil {
			return err.Error()
		}
	}
	if n.ID > 0 {
		n.Get()
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
	if n.ID > 0 {
		err = n.Update()
	} else {
		err = n.Create()
	}

	if err != nil {
		return err.Error()
	}
	return encoding.Must(enc.Encode(n))
}
func DeleteWebPropertyNote(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	var n webProperty_model.WebPropertyNote
	var err error

	idStr := r.FormValue("id")
	if idStr != "" {
		n.ID, err = strconv.Atoi(idStr)
		if err != nil {
			return err.Error()
		}
	} else {
		n.ID, err = strconv.Atoi(params["id"])
		if err != nil {
			return err.Error()
		}
	}
	err = n.Delete()
	if err != nil {
		return err.Error()
	}
	return encoding.Must(enc.Encode(n))
}

func GetWebPropertyRequirement(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	var req webProperty_model.WebPropertyRequirement
	var err error
	idStr := r.FormValue("id")
	if idStr != "" {
		req.RequirementID, err = strconv.Atoi(idStr)
		if err != nil {
			return err.Error()
		}
	} else {
		req.RequirementID, err = strconv.Atoi(params["id"])
		if err != nil {
			return err.Error()
		}
	}
	err = req.Get()
	if err != nil {
		return err.Error()
	}
	return encoding.Must(enc.Encode(req))

}

func CreateUpdateWebPropertyRequirement(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	var req webProperty_model.WebPropertyRequirement
	var err error
	idStr := r.FormValue("id")
	if idStr != "" {
		req.RequirementID, err = strconv.Atoi(idStr)
		if err != nil {
			return err.Error()
		}
	}
	if params["id"] != "" {
		req.RequirementID, err = strconv.Atoi(params["id"])
		if err != nil {
			return err.Error()
		}
	}
	if req.RequirementID > 0 {
		err = req.Get()
	}
	reqType := r.FormValue("reqType")
	requirement := r.FormValue("requirement")

	if reqType != "" {
		req.ReqType = reqType
	}
	if requirement != "" {
		req.Requirement = requirement
	}
	if req.RequirementID > 0 {
		err = req.Update()
	} else {
		err = req.Create()
	}
	if err != nil {
		return err.Error()
	}
	return encoding.Must(enc.Encode(req))

}

func DeleteWebPropertyRequirement(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	var req webProperty_model.WebPropertyRequirement
	var err error
	idStr := r.FormValue("id")
	if idStr != "" {
		req.RequirementID, err = strconv.Atoi(idStr)
		if err != nil {
			return err.Error()
		}
	} else {
		req.RequirementID, err = strconv.Atoi(params["id"])
		if err != nil {
			return err.Error()
		}
	}
	err = req.Delete()
	if err != nil {
		return err.Error()
	}
	return encoding.Must(enc.Encode(req))
}

func GetWebPropertyType(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	var t webProperty_model.WebPropertyType
	var err error
	idStr := r.FormValue("id")
	if idStr != "" {
		t.ID, err = strconv.Atoi(idStr)
		if err != nil {
			return err.Error()
		}
	} else {
		t.ID, err = strconv.Atoi(params["id"])
		if err != nil {
			return err.Error()
		}
	}
	err = t.Get()
	if err != nil {
		return err.Error()
	}
	return encoding.Must(enc.Encode(t))

}

func CreateUpdateWebPropertyType(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	var t webProperty_model.WebPropertyType
	var err error
	idStr := r.FormValue("id")
	if idStr != "" {
		t.ID, err = strconv.Atoi(idStr)
		if err != nil {
			return err.Error()
		}
	}
	if params["id"] != "" {
		t.ID, err = strconv.Atoi(params["id"])
		if err != nil {
			return err.Error()
		}
	}
	if t.ID > 0 {
		err = t.Get()
	}
	typeID := r.FormValue("typeID")
	theType := r.FormValue("type")

	if typeID != "" {
		t.TypeID, err = strconv.Atoi(typeID)
	}
	if theType != "" {
		t.Type = theType
	}
	if t.ID > 0 {
		err = t.Update()
	} else {
		err = t.Create()
	}
	if err != nil {
		return err.Error()
	}
	return encoding.Must(enc.Encode(t))

}

func DeleteWebPropertyType(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	var t webProperty_model.WebPropertyType
	var err error
	idStr := r.FormValue("id")
	if idStr != "" {
		t.ID, err = strconv.Atoi(idStr)
		if err != nil {
			return err.Error()
		}
	} else {
		t.ID, err = strconv.Atoi(params["id"])
		if err != nil {
			return err.Error()
		}
	}
	err = t.Delete()
	if err != nil {
		return err.Error()
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
	isFinalApproved := r.FormValue("metaDescription")
	isEnabledDate := r.FormValue("metaDescription")
	isDenied := r.FormValue("metaDescription")
	requestedDate := r.FormValue("metaDescription")
	typeID := r.FormValue("typeID")

	page := r.FormValue("page")
	results := r.FormValue("results")

	l, err := webProperty_model.Search(name, custID, badgeID, url, isEnabled, sellerID, webPropertyTypeID, isFinalApproved, isEnabledDate, isDenied, requestedDate, typeID, page, results)
	if err != nil {
		return err.Error()
	}

	return encoding.Must(enc.Encode(l))
}
