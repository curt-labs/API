package cartIntegration

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/curt-labs/GoAPI/helpers/apicontext"
	"github.com/curt-labs/GoAPI/helpers/encoding"
	"github.com/curt-labs/GoAPI/helpers/error"
	"github.com/curt-labs/GoAPI/models/cartIntegration"
	"github.com/curt-labs/GoAPI/models/customer"
	"github.com/go-martini/martini"
)

const dateFormat = "Jan 2, 2006"

//Came From CartIntegration Model - not sure of usefulness
func ParsePricePointFields(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	r.ParseForm()
	var p cartIntegration.PricePoint

	if custID, err := strconv.ParseInt(r.FormValue("custID"), 10, 64); err == nil {
		p.CartIntegration.CustID = int(custID)
	}

	if partID, err := strconv.ParseInt(r.FormValue("partID"), 10, 64); err == nil {
		p.CartIntegration.PartID = int(partID)
	}

	if price, err := strconv.ParseFloat(r.FormValue("price"), 64); err == nil {
		p.Price = price
	}

	if custPartId, err := strconv.ParseInt(r.FormValue("custPartID"), 10, 64); err == nil {
		p.CartIntegration.CustPartID = int(custPartId)
	}

	if isSale, err := strconv.ParseInt(r.FormValue("isSale"), 10, 64); err == nil {
		if isSale >= 1 {
			p.IsSale = 1
		} else if isSale <= 0 {
			p.IsSale = 0
		}
	}

	if startd, err := time.Parse(dateFormat, r.FormValue("sale_start")); err == nil {
		p.Sale_start = startd
	}

	if endd, err := time.Parse(dateFormat, r.FormValue("sale_end")); err == nil {
		p.Sale_end = endd
	}
	return encoding.Must(enc.Encode(p))
}

func GetCI(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params, dtx *apicontext.DataContext) string {
	var ci cartIntegration.CartIntegration
	var err error
	id := params["id"]
	if id == "" {
		apierror.GenerateError("Trouble getting cart integration ID", err, rw, r)
	}
	ci.ID, err = strconv.Atoi(id)

	err = ci.Get()
	if err != nil {
		apierror.GenerateError("Trouble getting cart integration", err, rw, r)
	}
	return encoding.Must(enc.Encode(ci))
}

func GetCIbyPart(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params, dtx *apicontext.DataContext) string {
	var ci cartIntegration.CartIntegration
	var err error
	id := params["id"]
	if id == "" {
		apierror.GenerateError("Trouble getting part number for cart integration", err, rw, r)
	}
	ci.PartID, err = strconv.Atoi(id)

	cis, err := cartIntegration.GetCartIntegrationsByPart(ci)
	if err != nil {
		apierror.GenerateError("Trouble getting cart integrations by part number", err, rw, r)
	}
	return encoding.Must(enc.Encode(cis))
}

func GetCIbyCustomer(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params, dtx *apicontext.DataContext) string {
	var ci cartIntegration.CartIntegration
	var err error
	id := params["id"]
	if id == "" {
		apierror.GenerateError("Trouble getting customer ID for cart integration", err, rw, r)
	}
	//get cust_id from (old) customerID
	var c customer.Customer
	c.CustomerId, err = strconv.Atoi(id)
	err = c.FindCustIdFromCustomerId()
	if err != nil {
		apierror.GenerateError("Trouble finding cust ID from CustomerID", err, rw, r)
	}
	//set cartIntegration customer ID to cust_id
	ci.CustID = c.Id

	cis, err := cartIntegration.GetCartIntegrationsByCustomer(ci)
	if err != nil {
		apierror.GenerateError("Trouble getting cart integrations by customer", err, rw, r)
	}
	return encoding.Must(enc.Encode(cis))
}

func SaveCI(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params, dtx *apicontext.DataContext) string {
	var ci cartIntegration.CartIntegration
	var err error
	id := params["id"]
	if id != "" {
		ci.ID, err = strconv.Atoi(id)
		err = ci.Get()
		if err != nil {
			apierror.GenerateError("Trouble getting cart integration", err, rw, r)
		}
	}

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		apierror.GenerateError("Trouble reading request body for saving cart integration", err, rw, r)
	}
	err = json.Unmarshal(body, &ci)
	if err != nil {
		apierror.GenerateError("Trouble unmarshaling request body for cart integration", err, rw, r)
	}
	if ci.ID != 0 {
		err = ci.Update()
	} else {
		err = ci.Create()
	}
	if err != nil {
		apierror.GenerateError("Trouble creating/updating cart integration", err, rw, r)
	}
	return encoding.Must(enc.Encode(ci))
}

func DeleteCI(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params, dtx *apicontext.DataContext) string {
	var ci cartIntegration.CartIntegration
	var err error
	id := params["id"]

	ci.ID, err = strconv.Atoi(id)
	if err != nil {
		apierror.GenerateError("Trouble getting cart integration ID", err, rw, r)
	}

	err = ci.Delete()
	if err != nil {
		apierror.GenerateError("Trouble deleting cart integration", err, rw, r)
	}

	return encoding.Must(enc.Encode(ci))
}

func GetCustomerPricing(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params, dtx *apicontext.DataContext) string {
	var err error
	custID, err := strconv.Atoi(params["custID"])
	if custID == 0 {
		apierror.GenerateError("Trouble getting custID for customer pricing", err, rw, r)
	}
	prices, err := cartIntegration.GetPricesByCustomerID(custID)
	if err != nil {
		apierror.GenerateError("Trouble getting prices by customer ID", err, rw, r)
	}
	return encoding.Must(enc.Encode(prices))
}
func GetCustomerPricingPaged(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params, dtx *apicontext.DataContext) string {
	var err error

	custID, err := strconv.Atoi(params["custID"])
	if custID == 0 {
		apierror.GenerateError("Trouble getting custID for paged customer pricing ", err, rw, r)
	}

	page, err := strconv.Atoi(params["page"])
	if page < 1 {
		apierror.GenerateError("Trouble getting page number for paged customer pricing", err, rw, r)
	}

	count, err := strconv.Atoi(params["count"])
	if count < 1 {
		apierror.GenerateError("Trouble getting count for paged customer pricing", err, rw, r)
	}

	prices, err := cartIntegration.GetPricesByCustomerIDPaged(custID, page, count)
	if custID == 0 {
		apierror.GenerateError("Trouble getting prices for paged customer pricing", err, rw, r)
	}

	return encoding.Must(enc.Encode(prices))
}

//kind of dumb, like the tedious version of len(prices); it existed in the cartIntegration project, so maybe it's needed for something
func GetCustomerPricingCount(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params, dtx *apicontext.DataContext) string {
	var err error
	custID, err := strconv.Atoi(params["custID"])
	if custID == 0 {
		apierror.GenerateError("Trouble getting custID for pricing count", err, rw, r)
	}
	count, err := cartIntegration.GetPricingCount(custID)
	if custID == 0 {
		apierror.GenerateError("Trouble getting pricing count", err, rw, r)
	}
	return encoding.Must(enc.Encode(count))
}
