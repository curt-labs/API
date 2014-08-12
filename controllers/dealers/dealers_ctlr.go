package dealers_ctlr

import (
	"github.com/curt-labs/GoAPI/helpers/encoding"
	"github.com/curt-labs/GoAPI/models/customer"
	"github.com/go-martini/martini"
	"net/http"
	"strconv"
)

func Etailers(w http.ResponseWriter, r *http.Request, enc encoding.Encoder) string {

	dealers, err := customer.GetEtailers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	return encoding.Must(enc.Encode(dealers))
}

// Sample Data
//
// latlng: 43.853282,-95.571675,45.800981,-90.468526
// center: 44.83536,-93.0201
//
// Old Path: http://curtmfg.com/WhereToBuy/getLocalDealersJSON?latlng=43.853282,-95.571675,45.800981,-90.468526&center=44.83536,-93.0201
func LocalDealers(w http.ResponseWriter, r *http.Request, enc encoding.Encoder) string {

	qs := r.URL.Query()
	latlng := qs.Get("latlng")
	center := qs.Get("center")

	dealers, err := customer.GetLocalDealers(center, latlng)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	return encoding.Must(enc.Encode(dealers))
}

func LocalRegions(w http.ResponseWriter, r *http.Request, enc encoding.Encoder) string {
	regions, err := customer.GetLocalRegions()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	return encoding.Must(enc.Encode(regions))
}

func LocalDealerTiers(w http.ResponseWriter, r *http.Request, enc encoding.Encoder) string {
	tiers := customer.GetLocalDealerTiers()

	return encoding.Must(enc.Encode(tiers))
}

func LocalDealerTypes(w http.ResponseWriter, r *http.Request, enc encoding.Encoder) string {
	types := customer.GetLocalDealerTypes()

	return encoding.Must(enc.Encode(types))
}

func PlatinumEtailers(w http.ResponseWriter, r *http.Request, enc encoding.Encoder) string {
	custs := customer.GetWhereToBuyDealers()
	return encoding.Must(enc.Encode(custs))
}

func GetLocation(w http.ResponseWriter, r *http.Request, params martini.Params, enc encoding.Encoder) string {
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

	loc, err := customer.GetLocationById(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	return encoding.Must(enc.Encode(loc))
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
