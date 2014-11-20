package cart_ctlr

import (
	"encoding/json"
	"github.com/curt-labs/GoAPI/helpers/encoding"
	"github.com/curt-labs/GoAPI/models/cart"
	"github.com/go-martini/martini"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

const (
	TimeLayout = "2006-01-02 15:04"
)

func GetCustomers(w http.ResponseWriter, req *http.Request, params martini.Params, enc encoding.Encoder) string {
	var since_id bson.ObjectId
	var created_at_min *time.Time
	var created_at_max *time.Time
	var updated_at_min *time.Time
	var updated_at_max *time.Time
	limit := 50
	page := 0
	qs := req.URL.Query()

	if since := qs.Get("since_id"); since != "" {
		if bson.IsObjectIdHex(since) {
			since_id = bson.ObjectIdHex(since)
		}
	}

	if created_min, err := time.Parse(TimeLayout, qs.Get("created_at_min")); err == nil {
		if created_max, err := time.Parse(TimeLayout, qs.Get("created_at_max")); err == nil {
			created_at_min = &created_min
			created_at_max = &created_max
		}
	}
	if updated_min, err := time.Parse(TimeLayout, qs.Get("updated_at_min")); err == nil {
		if updated_max, err := time.Parse(TimeLayout, qs.Get("updated_at_max")); err == nil {
			updated_at_min = &updated_min
			updated_at_max = &updated_max
		}
	}

	if l := qs.Get("limit"); l != "" {
		lmt, err := strconv.Atoi(l)
		if err == nil && lmt != 0 {
			limit = lmt
		}
	}
	if p := qs.Get("page"); p != "" {
		pg, err := strconv.Atoi(p)
		if err == nil && pg != 0 {
			page = pg
		}
	}

	custs, err := cart.CustomersSinceId(since_id, page, limit, created_at_min, created_at_max, updated_at_min, updated_at_max)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	return encoding.Must(enc.Encode(custs))
}

func GetCustomer(w http.ResponseWriter, req *http.Request, params martini.Params, enc encoding.Encoder) string {

	customerId := params["id"]

	if !bson.IsObjectIdHex(customerId) {
		http.Error(w, "invalid customer reference", http.StatusInternalServerError)
		return ""
	}

	c := cart.Customer{
		Id: bson.ObjectIdHex(customerId),
	}
	if err := c.Get(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	return encoding.Must(enc.Encode(c))
}

func AddCustomer(w http.ResponseWriter, req *http.Request, params martini.Params, enc encoding.Encoder) string {

	var c cart.Customer
	defer req.Body.Close()

	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	if err = json.Unmarshal(data, &c); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	if err = c.Insert(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	return encoding.Must(enc.Encode(c))
}

func EditCustomer(w http.ResponseWriter, req *http.Request, params martini.Params, enc encoding.Encoder) string {

	var c cart.Customer
	defer req.Body.Close()

	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	if err = json.Unmarshal(data, &c); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	customerId := params["id"]

	if !bson.IsObjectIdHex(customerId) {
		http.Error(w, "invalid customer reference", http.StatusInternalServerError)
		return ""
	}
	c.Id = bson.ObjectIdHex(customerId)

	if err = c.Update(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	return encoding.Must(enc.Encode(c))
}

func DeleteCustomer(w http.ResponseWriter, req *http.Request, params martini.Params, enc encoding.Encoder) string {

	var c cart.Customer
	customerId := params["id"]

	if !bson.IsObjectIdHex(customerId) {
		http.Error(w, "invalid customer reference", http.StatusInternalServerError)
		return ""
	}
	c.Id = bson.ObjectIdHex(customerId)

	if err := c.Delete(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	return encoding.Must(enc.Encode(c))
}

func GetAddresses(w http.ResponseWriter, req *http.Request, params martini.Params, enc encoding.Encoder) string {
	customerId := params["id"]
	limit := 50
	page := 1
	qs := req.URL.Query()

	if l := qs.Get("limit"); l != "" {
		lmt, err := strconv.Atoi(l)
		if err == nil && lmt != 0 {
			limit = lmt
		}
	}
	if p := qs.Get("page"); p != "" {
		pg, err := strconv.Atoi(p)
		if err == nil && pg != 0 {
			page = pg
		}
	}

	if !bson.IsObjectIdHex(customerId) {
		http.Error(w, "invalid customer reference", http.StatusInternalServerError)
		return ""
	}

	c := cart.Customer{
		Id: bson.ObjectIdHex(customerId),
	}
	if err := c.Get(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	addr := c.Addresses[:limit]
	if page > 1 && len(c.Addresses) >= ((page-1)*limit) {
		addr = c.Addresses[((page - 1) / limit):limit]
	}
	return encoding.Must(enc.Encode(addr))
}
