package dealers_ctlr

import (
	"../../helpers/plate"
	"../../models"
	"net/http"
	"strconv"
)

func Etailers(w http.ResponseWriter, r *http.Request) {

	dealers, err := models.GetEtailers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	plate.ServeFormatted(w, r, dealers)
}

// Sample Data
//
// latlng: 43.853282,-95.571675,45.800981,-90.468526
// center: 44.83536,-93.0201
//
// Old Path: http://curtmfg.com/WhereToBuy/getLocalDealersJSON?latlng=43.853282,-95.571675,45.800981,-90.468526&center=44.83536,-93.0201
func LocalDealers(w http.ResponseWriter, r *http.Request) {

	params := r.URL.Query()
	latlng := params.Get("latlng")
	center := params.Get("center")

	dealers, err := models.GetLocalDealers(center, latlng)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	plate.ServeFormatted(w, r, dealers)
}

func LocalRegions(w http.ResponseWriter, r *http.Request) {
	regions, err := models.GetLocalRegions()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	plate.ServeFormatted(w, r, regions)
}

func LocalDealerTiers(w http.ResponseWriter, r *http.Request) {
	tiers := models.GetLocalDealerTiers()
	plate.ServeFormatted(w, r, tiers)
}

func LocalDealerTypes(w http.ResponseWriter, r *http.Request) {
	types := models.GetLocalDealerTypes()
	plate.ServeFormatted(w, r, types)
}

func PlatinumEtailers(w http.ResponseWriter, r *http.Request) {
	custs := models.GetWhereToBuyDealers()
	plate.ServeFormatted(w, r, custs)
}

func GetLocation(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	str_id := params.Get(":id")
	if str_id == "" {
		http.Error(w, "You must supply a location identification number.", http.StatusInternalServerError)
		return
	}
	id, err := strconv.Atoi(str_id)
	if err != nil {
		http.Error(w, "You must supply a location identification number.", http.StatusInternalServerError)
		return
	}

	loc, err := models.GetLocationById(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	plate.ServeFormatted(w, r, loc)
}

func SearchLocations(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	search_term := params.Get(":search")
	if search_term == "" {
		search_term = params.Get("search")
	}
	locs, err := models.SearchLocations(search_term)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	plate.ServeFormatted(w, r, locs)
}

func SearchLocationsByType(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	search_term := params.Get(":search")
	if search_term == "" {
		search_term = params.Get("search")
	}
	locs, err := models.SearchLocationsByType(search_term)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	plate.ServeFormatted(w, r, locs)
}

func SearchLocationsByLatLng(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()

	// Get the latitude
	latitude := params.Get(":latitude")
	if latitude == "" {
		latitude = params.Get("latitude")
	}
	// Get the longitude
	longitude := params.Get(":longitude")
	if longitude == "" {
		longitude = params.Get("longitude")
	}

	latFloat, _ := strconv.ParseFloat(latitude, 64)
	lngFloat, _ := strconv.ParseFloat(longitude, 64)

	latlng := models.GeoLocation{
		Latitude:  latFloat,
		Longitude: lngFloat,
	}

	locs, err := models.SearchLocationsByLatLng(latlng)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	plate.ServeFormatted(w, r, locs)
}
