package dealers_ctlr

import (
	"github.com/curt-labs/GoAPI/helpers/apicontext"
	"github.com/curt-labs/GoAPI/helpers/encoding"
	apierr "github.com/curt-labs/GoAPI/helpers/error"
	"github.com/curt-labs/GoAPI/models/customer"
	"github.com/go-martini/martini"
	"net/http"
	"strconv"
)

func GetEtailers(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
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

	dealers, err := customer.GetEtailers(key)
	if err != nil {
		return err.Error()
	}
	return encoding.Must(enc.Encode(dealers))
}

// Sample Data
//
// latlng: 43.853282,-95.571675,45.800981,-90.468526
// center: 44.83536,-93.0201
//
// Old Path: http://curtmfg.com/WhereToBuy/getLocalDealersJSON?latlng=43.853282,-95.571675,45.800981,-90.468526&center=44.83536,-93.0201
// TODO - this method found in Dealers ctlr

func GetLocalDealers(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
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
	latlng := params["latlng"]
	if latlng == "" {
		latlng = qs.Get("latlng")
	}
	// Get the center
	center := params["center"]
	if center == "" {
		center = qs.Get("center")
	}

	dealerLocations, err := customer.GetLocalDealers(center, latlng)
	if err != nil {
		return err.Error()
	}
	return encoding.Must(enc.Encode(dealerLocations))

}

func GetLocalRegions(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	regions, err := customer.GetLocalRegions()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}
	return encoding.Must(enc.Encode(regions))
}

func GetLocalDealerTiers(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	tiers, err := customer.GetLocalDealerTiers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}
	return encoding.Must(enc.Encode(tiers))
}

func GetLocalDealerTypes(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	types, err := customer.GetLocalDealerTypes()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}
	return encoding.Must(enc.Encode(types))
}

func PlatinumEtailers(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	qs := r.URL.Query()
	key := qs.Get("key")
	if key == "" {
		key = r.FormValue("key")
	}

	if key == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return ""
	}

	cust, err := customer.GetWhereToBuyDealers(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}
	return encoding.Must(enc.Encode(cust))
}

func GetLocationById(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	str_id := params["id"]
	if str_id == "" {
		http.Error(w, "You must supply a location identification number.", http.StatusInternalServerError)
		return ""
	}
	id, err := strconv.Atoi(str_id)
	if err != nil {
		http.Error(w, "You must supply a location identification number.", http.StatusInternalServerError)
		return ""
	}
	var l customer.CustomerLocation
	l.Id = id
	// loc, err := customer.GetLocationById(id)
	err = l.Get()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	return encoding.Must(enc.Encode(l))
}

func GetAllBusinessClasses(w http.ResponseWriter, r *http.Request, params martini.Params, enc encoding.Encoder, dtx *apicontext.DataContext) string {
	classes, err := customer.GetAllBusinessClasses(dtx)
	if len(classes) == 0 {
		apierr.GenerateError("No results", err, w, r)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err.Error()
	}
	return encoding.Must(enc.Encode(classes))
}

func SearchLocations(w http.ResponseWriter, r *http.Request, params martini.Params, enc encoding.Encoder) string {
	search_term := params["search"]
	qs := r.URL.Query()
	if search_term == "" {
		search_term = qs.Get("search")
	}
	locs, err := customer.SearchLocations(search_term)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	return encoding.Must(enc.Encode(locs))
}

func SearchLocationsByType(w http.ResponseWriter, r *http.Request, params martini.Params, enc encoding.Encoder) string {
	qs := r.URL.Query()
	search_term := params["search"]
	if search_term == "" {
		search_term = qs.Get("search")
	}
	locs, err := customer.SearchLocationsByType(search_term)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	return encoding.Must(enc.Encode(locs))
}
func SearchLocationsByLatLng(w http.ResponseWriter, r *http.Request, params martini.Params, enc encoding.Encoder) string {
	qs := r.URL.Query()

	// Get the latitude
	latitude := params["latitude"]
	if latitude == "" {
		latitude = qs.Get("latitude")
	}
	// Get the longitude
	longitude := params["longitude"]
	if longitude == "" {
		longitude = qs.Get("longitude")
	}

	latFloat, _ := strconv.ParseFloat(latitude, 64)
	lngFloat, _ := strconv.ParseFloat(longitude, 64)

	latlng := customer.GeoLocation{
		Latitude:  latFloat,
		Longitude: lngFloat,
	}

	locs, err := customer.SearchLocationsByLatLng(latlng)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	return encoding.Must(enc.Encode(locs))
}
