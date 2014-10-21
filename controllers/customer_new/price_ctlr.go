package customer_ctlr_new

import (
	"github.com/curt-labs/GoAPI/helpers/encoding"
	"github.com/curt-labs/GoAPI/helpers/sortutil"
	"github.com/curt-labs/GoAPI/models/customer_new"
	"github.com/go-martini/martini"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	inputTimeFormat = "01/02/2006"
)

func GetPrice(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	var c customer_new.Price
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
	var c customer_new.Prices
	var err error

	c, err = customer_new.GetAllPrices()
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
	var w customer_new.Price
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
	}
	if partID != "" {
		w.PartID, err = strconv.Atoi(partID)
	}
	if price != "" {
		w.Price, err = strconv.ParseFloat(price, 64)
	}
	if isSale != "" {
		w.IsSale, err = strconv.Atoi(isSale)
	}
	if saleStart != "" {
		w.SaleStart, err = time.Parse(inputTimeFormat, saleStart)
	}
	if saleEnd != "" {
		w.SaleEnd, err = time.Parse(inputTimeFormat, saleEnd)
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
	var w customer_new.Price
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
	var ps customer_new.Prices
	var partID int
	qs := url.Query
	log.Print(qs)

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
	ps, err = customer_new.GetPricesByPart(partID)

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
	var ps customer_new.Prices
	var c customer_new.Customer

	id := r.FormValue("id")
	if id != "" {
		c.Id, err = strconv.Atoi(id)
		if err != nil {
			return err.Error()
		}
	}
	if params["id"] != "" {
		c.Id, err = strconv.Atoi(params["id"])
		if err != nil {
			return err.Error()
		}
	}

	start := r.FormValue("start")
	end := r.FormValue("end")
	startDate, err := time.Parse(inputTimeFormat, start)
	endDate, err := time.Parse(inputTimeFormat, end)
	if err != nil {
		return err.Error()
	}

	ps, err = c.GetPricesBySaleRange(startDate, endDate)
	if err != nil {
		return err.Error()
	}
	return encoding.Must(enc.Encode(ps))
}

func GetPriceByCustomer(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	var err error
	var ps customer_new.CustomerPrices
	var c customer_new.Customer

	id := r.FormValue("id")
	if id != "" {
		c.Id, err = strconv.Atoi(id)
		if err != nil {
			return err.Error()
		}
	}
	if params["id"] != "" {
		c.Id, err = strconv.Atoi(params["id"])
		if err != nil {
			return err.Error()
		}
	}

	ps, err = c.GetPricesByCustomer()
	if err != nil {
		return err.Error()
	}
	return encoding.Must(enc.Encode(ps))
}
