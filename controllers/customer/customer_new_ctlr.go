package customer_ctlr

import (
	"github.com/curt-labs/GoAPI/helpers/encoding"
	"github.com/curt-labs/GoAPI/models/customer"
	"github.com/curt-labs/GoAPI/models/part"
	"github.com/go-martini/martini"
	"log"
	"net/http"
	"strconv"
)

func GetCustomer_New(w http.ResponseWriter, r *http.Request, enc encoding.Encoder) string {
	qs := r.URL.Query()
	key := qs.Get("key")
	if key == "" {
		key = r.FormValue("key")
	}

	if key == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return ""
	}

	id, err := customer.GetCustomerIdFromKey(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	c := customer.Customer_New{
		Id: id,
	}

	err = c.GetCustomer_New()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	return encoding.Must(enc.Encode(c))
}

func GetLocations_New(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	qs := r.URL.Query()
	key := qs.Get("key")
	if key == "" {
		key = r.FormValue("key")
	}

	if key == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return ""
	}

	id, err := customer.GetCustomerIdFromKey(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}
	log.Print(id)
	c := customer.Customer_New{
		Id: id,
	}
	err = c.GetLocations_New()
	if err != nil {
		return err.Error()
	}
	return encoding.Must(enc.Encode(c))
}

func GetUsers_New(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	qs := r.URL.Query()
	key := qs.Get("key")
	if key == "" {
		key = r.FormValue("key")
	}

	if key == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return ""
	}

	id, err := customer.GetCustomerIdFromKey(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	c := customer.Customer_New{
		Id: id,
	}
	log.Print("CTRL")
	users, err := c.GetUsers_New()
	if err != nil {
		return err.Error()
	}
	return encoding.Must(enc.Encode(users))
}

func GetCustomerPrice_New(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	var err error
	qs := r.URL.Query()
	key := qs.Get("key")
	if key == "" {
		key = r.FormValue("key")
	}

	if key == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return ""
	}

	id := params["id"]
	var p part.Part
	if id != "" {
		p.PartId, err = strconv.Atoi(params["id"])
	}

	part, err := customer.GetCustomerPrice_New(key, p.PartId)
	if err != nil {
		return err.Error()
	}
	log.Print(part)
	return encoding.Must(enc.Encode(part))
}

func GetCustomerCartReference_New(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	var err error
	qs := r.URL.Query()
	key := qs.Get("key")
	if key == "" {
		key = r.FormValue("key")
	}

	if key == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return ""
	}

	id := params["id"]
	var p part.Part
	if id != "" {
		p.PartId, err = strconv.Atoi(params["id"])
	}

	ref, err := customer.GetCustomerCartReference_New(key, p.PartId)
	if err != nil {
		return err.Error()
	}
	return encoding.Must(enc.Encode(ref))
}

func GetEtailers_New(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	var err error
	qs := r.URL.Query()
	key := qs.Get("key")
	if key == "" {
		key = r.FormValue("key")
	}

	if key == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return ""
	}

	dealers, err := customer.GetEtailers_New()
	if err != nil {
		return err.Error()
	}
	return encoding.Must(enc.Encode(dealers))
}

//Goes in Dealers
// Sample Data
//
// latlng: 43.853282,-95.571675,45.800981,-90.468526
// center: 44.83536,-93.0201
//
// Old Path: http://curtmfg.com/WhereToBuy/getLocalDealersJSON?latlng=43.853282,-95.571675,45.800981,-90.468526&center=44.83536,-93.0201
// TODO - this method found in Dealers ctlr
func GetLocalDealers_New(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	var err error
	qs := r.URL.Query()
	key := qs.Get("key")
	if key == "" {
		key = r.FormValue("key")
	}

	if key == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return ""
	}
	// Get the latlng
	latlng := qs.Get("latlng")
	if latlng == "" {
		latlng = r.FormValue("latlng")
	}
	// Get the center
	center := qs.Get("center")
	if center == "" {
		center = r.FormValue("center")
	}

	dealers, err := customer.GetLocalDealers_New(center, latlng)
	if err != nil {
		return err.Error()
	}
	return encoding.Must(enc.Encode(dealers))

}

//Goes in Dealers
func GetLocalRegions_New(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	regions, err := customer.GetLocalRegions_New()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}
	return encoding.Must(enc.Encode(regions))
}
