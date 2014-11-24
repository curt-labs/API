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

func GetCustomers(w http.ResponseWriter, req *http.Request, params martini.Params, enc encoding.Encoder, shop *cart.Shop) string {

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

	var custs []cart.Customer
	var err error
	if since_id.Hex() != "" {
		custs, err = cart.CustomersSinceId(shop.Id, since_id, page, limit, created_at_min, created_at_max, updated_at_min, updated_at_max)
	} else {
		custs, err = cart.GetCustomers(shop.Id, page, limit, created_at_min, created_at_max, updated_at_min, updated_at_max)
	}

	if err != nil {
		generateError("", err, w, req)
		return ""
	}

	return encoding.Must(enc.Encode(custs))
}

func GetCustomer(w http.ResponseWriter, req *http.Request, params martini.Params, enc encoding.Encoder, shop *cart.Shop) string {

	customerId := params["id"]

	if !bson.IsObjectIdHex(customerId) {
		generateError("invalid customer reference", nil, w, req)
		return ""
	}

	c := cart.Customer{
		Id:     bson.ObjectIdHex(customerId),
		ShopId: shop.Id,
	}
	if err := c.Get(); err != nil {
		generateError(err.Error(), err, w, req)
		return ""
	}

	return encoding.Must(enc.Encode(c))
}

func AddCustomer(w http.ResponseWriter, req *http.Request, params martini.Params, enc encoding.Encoder, shop *cart.Shop) string {

	var c cart.Customer
	defer req.Body.Close()

	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		generateError(err.Error(), err, w, req)
		return ""
	}

	if err = json.Unmarshal(data, &c); err != nil {
		generateError(err.Error(), err, w, req)
		return ""
	}

	c.ShopId = shop.Id

	if err = c.Insert(); err != nil {
		generateError(err.Error(), err, w, req)
		return ""
	}

	return encoding.Must(enc.Encode(c))
}

func EditCustomer(w http.ResponseWriter, req *http.Request, params martini.Params, enc encoding.Encoder, shop *cart.Shop) string {

	var c cart.Customer
	defer req.Body.Close()

	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		generateError(err.Error(), err, w, req)
		return ""
	}

	if err = json.Unmarshal(data, &c); err != nil {
		generateError(err.Error(), err, w, req)
		return ""
	}

	customerId := params["id"]

	if !bson.IsObjectIdHex(customerId) {
		generateError("invalid customer reference", nil, w, req)
		return ""
	}
	c.Id = bson.ObjectIdHex(customerId)
	c.ShopId = shop.Id

	if err = c.Update(); err != nil {
		generateError(err.Error(), err, w, req)
		return ""
	}

	return encoding.Must(enc.Encode(c))
}

func DeleteCustomer(w http.ResponseWriter, req *http.Request, params martini.Params, enc encoding.Encoder, shop *cart.Shop) string {

	var c cart.Customer
	customerId := params["id"]

	if !bson.IsObjectIdHex(customerId) {
		generateError("invalid customer reference", nil, w, req)
		return ""
	}
	c.Id = bson.ObjectIdHex(customerId)
	c.ShopId = shop.Id

	if err := c.Delete(); err != nil {
		generateError(err.Error(), err, w, req)
		return ""
	}

	return encoding.Must(enc.Encode(c))
}

func SearchCustomer(w http.ResponseWriter, req *http.Request, params martini.Params, enc encoding.Encoder, shop *cart.Shop) string {
	qs := req.URL.Query()
	custs, err := cart.SearchCustomers(qs.Get("query"), shop.Id)
	if err != nil {
		generateError(err.Error(), err, w, req)
		return ""
	}

	return encoding.Must(enc.Encode(custs))
}

func GetCustomerOrders(w http.ResponseWriter, req *http.Request, params martini.Params, enc encoding.Encoder, shop *cart.Shop) string {

	customerId := params["id"]

	if !bson.IsObjectIdHex(customerId) {
		generateError("invalid customer reference", nil, w, req)
		return ""
	}

	c := cart.Customer{
		Id:     bson.ObjectIdHex(customerId),
		ShopId: shop.Id,
	}
	if err := c.Get(); err != nil {
		generateError(err.Error(), err, w, req)
		return ""
	}

	return encoding.Must(enc.Encode(c.Orders))
}
