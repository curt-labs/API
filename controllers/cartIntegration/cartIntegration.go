package cartIntegration

import (
	"encoding/json"
	"github.com/curt-labs/GoAPI/helpers/encoding"
	"github.com/curt-labs/GoAPI/models/cartIntegration"
	"github.com/go-martini/martini"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
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

func GetCI(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	var ci cartIntegration.CartIntegration
	var err error
	id := params["id"]
	if id == "" {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}
	ci.ID, err = strconv.Atoi(id)

	err = ci.Get()
	if err != nil {
		http.Error(w, err.Error(), http.StatusNoContent)
		return ""
	}
	return encoding.Must(enc.Encode(ci))
}

func GetCIbyPart(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	var ci cartIntegration.CartIntegration
	var err error
	id := params["id"]
	if id == "" {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}
	ci.PartID, err = strconv.Atoi(id)

	cis, err := cartIntegration.GetCartIntegrationsByPart(ci)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNoContent)
		return ""
	}
	return encoding.Must(enc.Encode(cis))
}

func GetCIbyCustomer(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	var ci cartIntegration.CartIntegration
	var err error
	id := params["id"]
	if id == "" {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}
	ci.CustID, err = strconv.Atoi(id)

	cis, err := cartIntegration.GetCartIntegrationsByCustomer(ci)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNoContent)
		return ""
	}
	return encoding.Must(enc.Encode(cis))
}

func SaveCI(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	var ci cartIntegration.CartIntegration
	var err error
	id := params["id"]
	if id != "" {
		ci.ID, err = strconv.Atoi(id)
		err = ci.Get()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return ""
		}
	}

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}
	err = json.Unmarshal(body, &ci)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}
	if ci.ID != 0 {
		err = ci.Update()
	} else {
		err = ci.Create()
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}
	return encoding.Must(enc.Encode(ci))
}

func DeleteCI(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	var ci cartIntegration.CartIntegration
	var err error
	id := params["id"]

	ci.ID, err = strconv.Atoi(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	err = ci.Delete()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}
	return encoding.Must(enc.Encode(ci))

}

func GetCustomerPricing(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	var err error

	custID, err := strconv.Atoi(params["custID"])
	if custID == 0 {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	prices, err := cartIntegration.GetPricesByCustomerID(custID)
	if custID == 0 {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}
	return encoding.Must(enc.Encode(prices))
}
func GetCustomerPricingPaged(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	var err error

	custID, err := strconv.Atoi(params["custID"])
	if custID == 0 {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	page, err := strconv.Atoi(params["page"])
	if page == 0 {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	count, err := strconv.Atoi(params["count"])
	if count == 0 {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	prices, err := cartIntegration.GetPricesByCustomerIDPaged(custID, page, count)
	if custID == 0 {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}
	return encoding.Must(enc.Encode(prices))

}

//kind of dumb, like the tedious version of len(prices); it existed in the cartIntegration project, so maybe it's needed for something
func GetCustomerPricingCount(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	var err error
	custID, err := strconv.Atoi(params["custID"])
	if custID == 0 {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}
	count, err := cartIntegration.GetPricingCount(custID)
	if custID == 0 {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}
	return encoding.Must(enc.Encode(count))
}
