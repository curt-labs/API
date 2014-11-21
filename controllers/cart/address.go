package cart_ctlr

import (
	"github.com/curt-labs/GoAPI/helpers/encoding"
	"github.com/curt-labs/GoAPI/models/cart"
	"github.com/go-martini/martini"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"strconv"
)

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
