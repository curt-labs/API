package cart_ctlr

import (
	"encoding/json"
	"github.com/curt-labs/GoAPI/helpers/encoding"
	"github.com/curt-labs/GoAPI/helpers/error"
	"github.com/curt-labs/GoAPI/models/cart"
	"github.com/go-martini/martini"
	"io/ioutil"
	"net/http"
)

// Login a customer for a given shop.
func AccountLogin(w http.ResponseWriter, req *http.Request, params martini.Params, enc encoding.Encoder, shop *cart.Shop) string {

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

	if err = c.Login(req.Referer()); err != nil {
		apierror.GenerateError(err.Error(), err, w, req)
		return ""
	}

	return encoding.Must(enc.Encode(c))
}

// Create a customer for a
// given shop.
func AddAccount(w http.ResponseWriter, req *http.Request, params martini.Params, enc encoding.Encoder, shop *cart.Shop) string {

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

// Get an account for a given shop
func GetAccount(w http.ResponseWriter, req *http.Request, params martini.Params, enc encoding.Encoder, shop *cart.Shop, token string) string {
	cust := cart.Customer{
		ShopId: shop.Id,
	}
	var err error

	cust.Id, err = cart.IdentifierFromToken(token)
	if err != nil {
		apierror.GenerateError(err.Error(), err, w, req)
		return ""
	}

	if err = cust.Get(); err != nil {
		apierror.GenerateError(err.Error(), err, w, req)
		return ""
	}

	return encoding.Must(enc.Encode(cust))
}

// Edit an account for a given shop.
func EditAccount(w http.ResponseWriter, req *http.Request, params martini.Params, enc encoding.Encoder, shop *cart.Shop, token string) string {

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

	c.Id, err = cart.IdentifierFromToken(token)
	if err != nil {
		apierror.GenerateError(err.Error(), err, w, req)
		return ""
	}

	if err = c.Update(); err != nil {
		apierror.GenerateError(err.Error(), err, w, req)
		return ""
	}

	return encoding.Must(enc.Encode(c))
}
