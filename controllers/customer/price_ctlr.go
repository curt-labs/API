package customer_ctlr

import (
	"github.com/curt-labs/GoAPI/helpers/encoding"
	"github.com/curt-labs/GoAPI/helpers/sortutil"
	"github.com/curt-labs/GoAPI/models/customer"
	"github.com/go-martini/martini"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	inputTimeFormat = "01/02/2006"
)

func GetPrice(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	var c customer.Price
	var err error
	id := r.FormValue("id")
	if id != "" {
		c.ID, err = strconv.Atoi(id)
		if err != nil {
			return err.Error()
		}
	}
	if params["id"] != "" {
		c.ID, err = strconv.Atoi(params["id"])
		if err != nil {
			return err.Error()
		}
	}

	err = c.Get()
	if err != nil {
		return err.Error()
	}
	return encoding.Must(enc.Encode(c))
}
func GetAllPrices(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder) string {
	var c customer.Prices
	var err error

	c, err = customer.GetAllPrices()
	if err != nil {
		return err.Error()
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

	id := r.FormValue("id")
	if id != "" {
		w.ID, err = strconv.Atoi(id)
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

	custID := r.FormValue("custID")
	partID := r.FormValue("partID")
	price := r.FormValue("price")
	isSale := r.FormValue("isSale")
	saleStart := r.FormValue("saleStart")
	saleEnd := r.FormValue("saleEnd")

	if custID != "" {
		w.CustID, err = strconv.Atoi(custID)
		if err != nil {
			return err.Error()
		}
	}
	if partID != "" {
		w.PartID, err = strconv.Atoi(partID)
		if err != nil {
			return err.Error()
		}
	}
	if price != "" {
		w.Price, err = strconv.ParseFloat(price, 64)
		if err != nil {
			return err.Error()
		}
	}
	if isSale != "" {
		saleBool, err := strconv.ParseBool(isSale)
		if err != nil {
			return err.Error()
		}
		if saleBool == true {
			w.IsSale = 1
		}
	}
	if saleStart != "" {
		w.SaleStart, err = time.Parse(inputTimeFormat, saleStart)
		if err != nil {
			return err.Error()
		}
	}
	if saleEnd != "" {
		w.SaleEnd, err = time.Parse(inputTimeFormat, saleEnd)
		if err != nil {
			return err.Error()
		}
	}

	if w.ID > 0 {
		err = w.Update()
	} else {

		err = w.Create()
	}

	if err != nil {
		return err.Error()
	}
	return encoding.Must(enc.Encode(w))
}

func DeletePrice(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	var w customer.Price
	var err error

	id := r.FormValue("id")
	if id != "" {
		w.ID, err = strconv.Atoi(id)
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
	w.Delete()

	return encoding.Must(enc.Encode(w))
}
func GetPricesByPart(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	var err error
	var ps customer.Prices
	var partID int

	id := r.FormValue("id")
	if id != "" {
		partID, err = strconv.Atoi(id)
		if err != nil {
			return err.Error()
		}
	}
	if params["id"] != "" {
		partID, err = strconv.Atoi(params["id"])
		if err != nil {
			return err.Error()
		}
	}
	ps, err = customer.GetPricesByPart(partID)

	sort := r.FormValue("sort")
	direction := r.FormValue("direction")
	if sort != "" {
		if strings.ContainsAny(direction, "esc") {
			sortutil.DescByField(ps, sort)
		} else {
			sortutil.AscByField(ps, sort)
		}
	}
	if err != nil {
		return err.Error()
	}
	return encoding.Must(enc.Encode(ps))
}

func GetSales(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	var err error
	var ps customer.Prices
	var c customer.Customer
	id := r.FormValue("id")
	if id != "" {
		c.Id, err = strconv.Atoi(id)
		if err != nil {
			return err.Error()
		}
	} else {
		id = params["id"]
		if id == "" {
			return "No Customer Id Supplied."
		}
		c.Id, err = strconv.Atoi(id)
	}

	start := r.FormValue("start")
	end := r.FormValue("end")

	startDate, err := time.Parse(inputTimeFormat, start)
	endDate, err := time.Parse(inputTimeFormat, end)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return ""
	}

	ps, err = c.GetPricesBySaleRange(startDate, endDate)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return ""
	}
	return encoding.Must(enc.Encode(ps))
}

func GetPriceByCustomer(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	var err error
	var ps customer.CustomerPrices
	var c customer.Customer

	id := r.FormValue("id")

	if id != "" {
		c.Id, err = strconv.Atoi(id)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return ""
		}
	}
	if params["id"] != "" {
		c.Id, err = strconv.Atoi(params["id"])
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return ""
		}
	}
	ps, err = c.GetPricesByCustomer()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return ""
	}
	return encoding.Must(enc.Encode(ps))
}
