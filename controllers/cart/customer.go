// This section is built as a management portal for a shop
// to view customer information.
//
// The customer will use endpoints defined in the account speicification. This
// allows us to separate the logic/authentication for the customer.

package cart_ctlr

import (
	"encoding/json"
	"github.com/curt-labs/API/helpers/encoding"
	"github.com/curt-labs/API/helpers/error"
	"github.com/curt-labs/API/models/cart"
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

// Get a list of all customers for a
// given shop.
// Params (optional):
// @since_id bson.ObjectId
// @created_at_min time.Time 2006-01-02 15:04
// @created_at_max time.Time 2006-01-02 15:04
// @updated_at_min time.Time 2006-01-02 15:04
// @updated_at_max time.Time 2006-01-02 15:04
// @page int
// @limit int
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
		apierror.GenerateError("", err, w, req)
		return ""
	}

	return encoding.Must(enc.Encode(custs))
}

// Get a specific customer for a
// given shop.
func GetCustomer(w http.ResponseWriter, req *http.Request, params martini.Params, enc encoding.Encoder, shop *cart.Shop) string {

	customerId := params["id"]

	if !bson.IsObjectIdHex(customerId) {
		apierror.GenerateError("invalid customer reference", nil, w, req)
		return ""
	}

	c := cart.Customer{
		Id:     bson.ObjectIdHex(customerId),
		ShopId: shop.Id,
	}
	if err := c.Get(); err != nil {
		apierror.GenerateError(err.Error(), err, w, req)
		return ""
	}

	return encoding.Must(enc.Encode(c))
}

// Create a customer for a
// given shop.
func AddCustomer(w http.ResponseWriter, req *http.Request, params martini.Params, enc encoding.Encoder, shop *cart.Shop) string {

	var c cart.Customer
	defer req.Body.Close()

	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		apierror.GenerateError(err.Error(), err, w, req)
		return ""
	}

	if err = json.Unmarshal(data, &c); err != nil {
		apierror.GenerateError(err.Error(), err, w, req)
		return ""
	}

	c.ShopId = shop.Id

	if err = c.Insert(req.Referer()); err != nil {
		apierror.GenerateError(err.Error(), err, w, req)
		return ""
	}

	return encoding.Must(enc.Encode(c))
}

// Edit an existing customer for a
// given shop.
func EditCustomer(w http.ResponseWriter, req *http.Request, params martini.Params, enc encoding.Encoder, shop *cart.Shop) string {

	var c cart.Customer
	defer req.Body.Close()

	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		apierror.GenerateError(err.Error(), err, w, req)
		return ""
	}

	if err = json.Unmarshal(data, &c); err != nil {
		apierror.GenerateError(err.Error(), err, w, req)
		return ""
	}

	customerId := params["id"]

	if !bson.IsObjectIdHex(customerId) {
		apierror.GenerateError("invalid customer reference", nil, w, req)
		return ""
	}
	c.Id = bson.ObjectIdHex(customerId)
	c.ShopId = shop.Id

	if err = c.Update(); err != nil {
		apierror.GenerateError(err.Error(), err, w, req)
		return ""
	}

	return encoding.Must(enc.Encode(c))
}

// Delete a customer for a given shop.
// Note: Can't delete if the customer has existing orders.
func DeleteCustomer(w http.ResponseWriter, req *http.Request, params martini.Params, enc encoding.Encoder, shop *cart.Shop) string {

	var c cart.Customer
	customerId := params["id"]

	if !bson.IsObjectIdHex(customerId) {
		apierror.GenerateError("invalid customer reference", nil, w, req)
		return ""
	}
	c.Id = bson.ObjectIdHex(customerId)
	c.ShopId = shop.Id

	if err := c.Delete(); err != nil {
		apierror.GenerateError(err.Error(), err, w, req)
		return ""
	}

	return ""
}

// Search customer records for a given shop.
// Params:
// @query string
func SearchCustomer(w http.ResponseWriter, req *http.Request, params martini.Params, enc encoding.Encoder, shop *cart.Shop) string {
	qs := req.URL.Query()
	custs, err := cart.SearchCustomers(qs.Get("query"), shop.Id)
	if err != nil {
		apierror.GenerateError(err.Error(), err, w, req)
		return ""
	}

	return encoding.Must(enc.Encode(custs))
}

// Get order history for a specific customer
// of a given shop.
func GetCustomerOrders(w http.ResponseWriter, req *http.Request, params martini.Params, enc encoding.Encoder, shop *cart.Shop) string {

	customerId := params["id"]

	if !bson.IsObjectIdHex(customerId) {
		apierror.GenerateError("invalid customer reference", nil, w, req)
		return ""
	}

	c := cart.Customer{
		Id:     bson.ObjectIdHex(customerId),
		ShopId: shop.Id,
	}
	if err := c.Get(); err != nil {
		apierror.GenerateError(err.Error(), err, w, req)
		return ""
	}

	return encoding.Must(enc.Encode(c.Orders))
}
