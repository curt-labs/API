package dealers_ctlr

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/curt-labs/GoAPI/helpers/apicontext"
	"github.com/curt-labs/GoAPI/helpers/encoding"
	"github.com/curt-labs/GoAPI/helpers/error"
	"github.com/curt-labs/GoAPI/models/customer"
	"github.com/go-martini/martini"
)

func GetEtailers(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params, dtx *apicontext.DataContext) string {
	dealers, err := customer.GetEtailers(dtx)
	if err != nil {
		apierror.GenerateError("Error retrieving etailers.", err, w, r)
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
		apierror.GenerateError("Unauthorized.", err, w, r)
	}
	// Get the latlng
	latlng := params["latlng"]
	if latlng == "" {
		latlng = qs.Get("latlng")
	}

	var distance int
	if qs.Get("distince") != "" {
		distance, _ = strconv.Atoi(qs.Get("distance"))
	}
	var count int
	if qs.Get("count") != "" {
		count, _ = strconv.Atoi(qs.Get("count"))
	}
	if count == 0 {
		count = 50
	}
	var skip int
	if qs.Get("skip") != "" {
		skip, _ = strconv.Atoi(qs.Get("skip"))
	}

	dealerLocations, err := customer.GetLocalDealers(latlng, distance, skip, count)
	if err != nil {
		apierror.GenerateError("Error retrieving locations.", err, w, r)
	}

	return encoding.Must(enc.Encode(dealerLocations))

}

func GetLocalRegions(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	regions, err := customer.GetLocalRegions()
	if err != nil {
		apierror.GenerateError("Error retrieving local regions.", err, w, r)
	}
	return encoding.Must(enc.Encode(regions))
}

func GetLocalDealerTiers(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params, dtx *apicontext.DataContext) string {
	tiers, err := customer.GetLocalDealerTiers(dtx)
	if err != nil {
		apierror.GenerateError("Error retrieving dealer tiers.", err, w, r)
	}
	return encoding.Must(enc.Encode(tiers))
}

func GetLocalDealerTypes(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params, dtx *apicontext.DataContext) string {
	types, err := customer.GetLocalDealerTypes(dtx)
	if err != nil {
		apierror.GenerateError("Error retrieving dealer types.", err, w, r)
	}
	return encoding.Must(enc.Encode(types))
}

func PlatinumEtailers(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params, dtx *apicontext.DataContext) string {
	cust, err := customer.GetWhereToBuyDealers(dtx)
	if err != nil {
		apierror.GenerateError("Error retrieving platinum etailers.", err, w, r)
	}
	return encoding.Must(enc.Encode(cust))
}

func GetLocationById(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	str_id := params["id"]
	if str_id == "" {
		apierror.GenerateError("You must supply a location identification number.", errors.New("No id."), w, r)
	}
	id, err := strconv.Atoi(str_id)
	if err != nil {
		apierror.GenerateError("You must supply a location identification number.", err, w, r)
	}
	var l customer.CustomerLocation
	l.Id = id
	// loc, err := customer.GetLocationById(id)
	err = l.Get()
	if err != nil {
		apierror.GenerateError("Error retrieving locations.", err, w, r)
	}

	return encoding.Must(enc.Encode(l))
}

func GetAllBusinessClasses(w http.ResponseWriter, r *http.Request, params martini.Params, enc encoding.Encoder, dtx *apicontext.DataContext) string {
	classes, err := customer.GetAllBusinessClasses(dtx)
	if err != nil {
		apierror.GenerateError("No results.", err, w, r)
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
		apierror.GenerateError("Error searching locations.", err, w, r)
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
		apierror.GenerateError("Error searching locations.", err, w, r)
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
		apierror.GenerateError("Error searching locations.", err, w, r)
	}

	return encoding.Must(enc.Encode(locs))
}
