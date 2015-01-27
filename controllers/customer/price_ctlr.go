package customer_ctlr

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/curt-labs/GoAPI/helpers/encoding"
	"github.com/curt-labs/GoAPI/helpers/error"
	"github.com/curt-labs/GoAPI/helpers/sortutil"
	"github.com/curt-labs/GoAPI/models/customer"
	"github.com/go-martini/martini"
)

const (
	inputTimeFormat = "01/02/2006"
)

func GetPrice(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	var c customer.Price
	var err error

	id := r.FormValue("id")
	if id == "" {
		id = params["id"]
	}

	if c.ID, err = strconv.Atoi(id); err != nil {
		apierror.GenerateError("Trouble getting price ID", err, rw, r)
		return ""
	}

	if err = c.Get(); err != nil {
		apierror.GenerateError("Trouble getting price", err, rw, r)
		return ""
	}

	return encoding.Must(enc.Encode(c))
}

func GetAllPrices(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder) string {
	var c customer.Prices
	var err error

	if c, err = customer.GetAllPrices(); err != nil {
		apierror.GenerateError("Trouble getting all prices", err, rw, r)
		return ""
	}

	sort := r.FormValue("sort")
	direction := r.FormValue("direction")

	if sort != "" {
		if strings.ContainsAny(direction, "esc") {
			sortutil.CiDescByField(c, sort)
		} else {
			sortutil.CiAscByField(c, sort)
		}
	}

	return encoding.Must(enc.Encode(c))
}

func CreateUpdatePrice(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	var w customer.Price
	var err error

	if r.FormValue("id") != "" || params["id"] != "" {
		id := r.FormValue("id")
		if id == "" {
			id = params["id"]
		}

		if w.ID, err = strconv.Atoi(id); err != nil {
			apierror.GenerateError("Trouble getting price ID", err, rw, r)
			return ""
		}

		if err = w.Get(); err != nil {
			apierror.GenerateError("Trouble getting price", err, rw, r)
			return ""
		}
	}

	custID := r.FormValue("custID")
	partID := r.FormValue("partID")
	price := r.FormValue("price")
	isSale := r.FormValue("isSale")
	saleStart := r.FormValue("saleStart")
	saleEnd := r.FormValue("saleEnd")

	if custID != "" {
		if w.CustID, err = strconv.Atoi(custID); err != nil {
			apierror.GenerateError("Trouble getting customer ID", err, rw, r)
			return ""
		}
	}

	if partID != "" {
		if w.PartID, err = strconv.Atoi(partID); err != nil {
			apierror.GenerateError("Trouble getting part ID", err, rw, r)
			return ""
		}
	}

	if price != "" {
		if w.Price, err = strconv.ParseFloat(price, 64); err != nil {
			apierror.GenerateError("Trouble getting price", err, rw, r)
			return ""
		}
	}

	if isSale != "" {
		saleBool := false
		if saleBool, err = strconv.ParseBool(isSale); err != nil {
			apierror.GenerateError("Trouble setting sale", err, rw, r)
			return ""
		}
		if saleBool {
			w.IsSale = 1
		}
	}

	if saleStart != "" {
		if w.SaleStart, err = time.Parse(inputTimeFormat, saleStart); err != nil {
			apierror.GenerateError("Trouble getting sale start", err, rw, r)
			return ""
		}
	}

	if saleEnd != "" {
		if w.SaleEnd, err = time.Parse(inputTimeFormat, saleEnd); err != nil {
			apierror.GenerateError("Trouble getting sale end", err, rw, r)
			return ""
		}
	}

	if w.ID > 0 {
		err = w.Update()
	} else {
		err = w.Create()
	}

	if err != nil {
		msg := "Trouble creating customer price"
		if w.ID > 0 {
			msg = "Trouble updating customer price"
		}
		apierror.GenerateError(msg, err, rw, r)
		return ""
	}

	return encoding.Must(enc.Encode(w))
}

func DeletePrice(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	var w customer.Price
	var err error

	id := r.FormValue("id")
	if id == "" {
		id = params["id"]
	}

	if w.ID, err = strconv.Atoi(id); err != nil {
		apierror.GenerateError("Trouble gett price ID", err, rw, r)
		return ""
	}

	if err = w.Delete(); err != nil {
		apierror.GenerateError("Trouble deleting price", err, rw, r)
		return ""
	}

	return encoding.Must(enc.Encode(w))
}

func GetPricesByPart(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	var err error
	var ps customer.Prices
	var partID int

	id := r.FormValue("id")
	if id == "" {
		id = params["id"]
	}

	if partID, err = strconv.Atoi(id); err != nil {
		apierror.GenerateError("Trouble getting part ID", err, rw, r)
		return ""
	}

	if ps, err = customer.GetPricesByPart(partID); err != nil {
		apierror.GenerateError("Trouble getting prices by part", err, rw, r)
		return ""
	}

	sort := r.FormValue("sort")
	direction := r.FormValue("direction")
	if sort != "" {
		if strings.ContainsAny(direction, "esc") {
			sortutil.DescByField(ps, sort)
		} else {
			sortutil.AscByField(ps, sort)
		}
	}

	return encoding.Must(enc.Encode(ps))
}

func GetSales(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	var err error
	var ps customer.Prices
	var c customer.Customer

	id := r.FormValue("id")
	if id == "" {
		id = params["id"]
	}

	if c.Id, err = strconv.Atoi(id); err != nil {
		apierror.GenerateError("Trouble getting customer ID", err, rw, r)
		return ""
	}

	start := r.FormValue("start")
	end := r.FormValue("end")

	startDate, err := time.Parse(inputTimeFormat, start)
	if err != nil {
		apierror.GenerateError("Trouble getting sales start date", err, rw, r)
		return ""
	}

	endDate, err := time.Parse(inputTimeFormat, end)
	if err != nil {
		apierror.GenerateError("Trouble getting sales end date", err, rw, r)
		return ""
	}

	ps, err = c.GetPricesBySaleRange(startDate, endDate)
	if err != nil {
		apierror.GenerateError("Trouble getting prices by sales range", err, rw, r)
		return ""
	}

	return encoding.Must(enc.Encode(ps))
}

func GetPriceByCustomer(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	var err error
	var ps customer.CustomerPrices
	var c customer.Customer

	id := r.FormValue("id")
	if id == "" {
		id = params["id"]
	}

	if c.Id, err = strconv.Atoi(id); err != nil {
		apierror.GenerateError("Trouble getting customer ID", err, rw, r)
		return ""
	}

	if ps, err = c.GetPricesByCustomer(); err != nil {
		apierror.GenerateError("Trouble getting prices by customer ID", err, rw, r)
		return ""
	}

	return encoding.Must(enc.Encode(ps))
}
