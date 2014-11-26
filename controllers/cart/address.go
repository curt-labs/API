package cart_ctlr

import (
	"encoding/json"
	"github.com/curt-labs/GoAPI/helpers/encoding"
	"github.com/curt-labs/GoAPI/helpers/error"
	"github.com/curt-labs/GoAPI/models/cart"
	"github.com/go-martini/martini"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"net/http"
	"strconv"
)

func GetAddresses(w http.ResponseWriter, req *http.Request, params martini.Params, enc encoding.Encoder, shop *cart.Shop) string {
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

	addr := c.Addresses
	if len(c.Addresses) > 0 {
		addr = c.Addresses[:limit]
		if page > 1 && len(c.Addresses) >= ((page-1)*limit) {
			addr = c.Addresses[((page - 1) / limit):limit]
		}
	}

	return encoding.Must(enc.Encode(addr))
}

func GetAddress(w http.ResponseWriter, req *http.Request, params martini.Params, enc encoding.Encoder, shop *cart.Shop) string {
	customerId := params["id"]
	addressId := params["address"]

	if !bson.IsObjectIdHex(customerId) {
		apierror.GenerateError("invalid customer reference", nil, w, req)
		return ""
	}
	if !bson.IsObjectIdHex(addressId) {
		apierror.GenerateError("invalid address reference", nil, w, req)
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

	var address *cart.CustomerAddress
	for _, addr := range c.Addresses {
		if addr.Id.Hex() == addressId {
			address = &addr
			break
		}
	}

	if address == nil {
		apierror.GenerateError("invalid address reference", nil, w, req)
		return ""
	}

	return encoding.Must(enc.Encode(address))
}

func AddAddress(w http.ResponseWriter, req *http.Request, params martini.Params, enc encoding.Encoder, shop *cart.Shop) string {

	var address cart.CustomerAddress
	var c cart.Customer

	customerId := params["id"]
	if !bson.IsObjectIdHex(customerId) {
		apierror.GenerateError("invalid customer reference", nil, w, req)
		return ""
	}

	c.Id = bson.ObjectIdHex(customerId)
	c.ShopId = shop.Id
	if err := c.Get(); err != nil {
		apierror.GenerateError("", err, w, req)
		return ""
	}

	defer req.Body.Close()

	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		apierror.GenerateError(err.Error(), err, w, req)
		return ""
	}

	if err = json.Unmarshal(data, &address); err != nil {
		apierror.GenerateError(err.Error(), err, w, req)
		return ""
	}

	// audit the address
	if err := address.Validate(); err != nil {
		apierror.GenerateError(err.Error(), err, w, req)
		return ""
	}

	addressId := bson.NewObjectId()
	address.Id = &addressId
	c.Addresses = append(c.Addresses, address)

	if err = c.Update(); err != nil {
		apierror.GenerateError(err.Error(), err, w, req)
		return ""
	}

	return encoding.Must(enc.Encode(address))
}

func EditAddress(w http.ResponseWriter, req *http.Request, params martini.Params, enc encoding.Encoder, shop *cart.Shop) string {

	var address cart.CustomerAddress
	var c cart.Customer

	customerId := params["id"]
	if !bson.IsObjectIdHex(customerId) {
		apierror.GenerateError("invalid customer reference", nil, w, req)
		return ""
	}

	c.Id = bson.ObjectIdHex(customerId)
	c.ShopId = shop.Id
	if err := c.Get(); err != nil {
		apierror.GenerateError("", err, w, req)
		return ""
	}

	addressId := params["address"]
	if !bson.IsObjectIdHex(addressId) {
		apierror.GenerateError("invalid address reference", nil, w, req)
		return ""
	}
	defer req.Body.Close()

	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		apierror.GenerateError(err.Error(), err, w, req)
		return ""
	}

	if err = json.Unmarshal(data, &address); err != nil {
		apierror.GenerateError(err.Error(), err, w, req)
		return ""
	}

	// audit the address
	if err := address.Validate(); err != nil {
		apierror.GenerateError(err.Error(), err, w, req)
		return ""
	}

	found := false
	for i, addr := range c.Addresses {
		if addr.Id.Hex() == addressId {
			c.Addresses[i] = address
			found = true
			break
		}
	}

	if !found {
		apierror.GenerateError("invalid address reference", nil, w, req)
		return ""
	}

	if err = c.Update(); err != nil {
		apierror.GenerateError(err.Error(), err, w, req)
		return ""
	}

	return encoding.Must(enc.Encode(address))
}

func SetDefaultAddress(w http.ResponseWriter, req *http.Request, params martini.Params, enc encoding.Encoder, shop *cart.Shop) string {

	var c cart.Customer

	customerId := params["id"]
	if !bson.IsObjectIdHex(customerId) {
		apierror.GenerateError("invalid customer reference", nil, w, req)
		return ""
	}

	c.Id = bson.ObjectIdHex(customerId)
	c.ShopId = shop.Id
	if err := c.Get(); err != nil {
		apierror.GenerateError("", err, w, req)
		return ""
	}

	addressId := params["address"]
	if !bson.IsObjectIdHex(addressId) {
		apierror.GenerateError("invalid address reference", nil, w, req)
		return ""
	}

	found := false
	for _, addr := range c.Addresses {
		if addr.Id.Hex() == addressId {
			c.DefaultAddress = &addr
			found = true
			break
		}
	}

	if !found {
		apierror.GenerateError("invalid address reference", nil, w, req)
		return ""
	}

	if err := c.Update(); err != nil {
		apierror.GenerateError(err.Error(), err, w, req)
		return ""
	}

	return encoding.Must(enc.Encode(c))
}

func DeleteAddress(w http.ResponseWriter, req *http.Request, params martini.Params, enc encoding.Encoder, shop *cart.Shop) string {

	var c cart.Customer

	customerId := params["id"]
	if !bson.IsObjectIdHex(customerId) {
		apierror.GenerateError("invalid customer reference", nil, w, req)
		return ""
	}

	c.Id = bson.ObjectIdHex(customerId)
	c.ShopId = shop.Id
	if err := c.Get(); err != nil {
		apierror.GenerateError("", err, w, req)
		return ""
	}

	addressId := params["address"]
	if !bson.IsObjectIdHex(addressId) {
		apierror.GenerateError("invalid address reference", nil, w, req)
		return ""
	}

	// Make sure this isn't the default address
	if c.DefaultAddress != nil && c.DefaultAddress.Id.Hex() == addressId {
		apierror.GenerateError("removing a default address is not allowed", nil, w, req)
		return ""
	}

	found := false
	for i, addr := range c.Addresses {
		if addr.Id.Hex() == addressId {
			c.Addresses = append(c.Addresses[:i], c.Addresses[i+1:]...)
			found = true
			break
		}
	}

	if !found {
		apierror.GenerateError("invalid address reference", nil, w, req)
		return ""
	}

	if err := c.Update(); err != nil {
		apierror.GenerateError(err.Error(), err, w, req)
		return ""
	}

	return ""
}
