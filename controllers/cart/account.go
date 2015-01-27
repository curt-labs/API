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

// Login a specific customer for a
// given shop.
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
