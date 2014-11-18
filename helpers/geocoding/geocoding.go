/**
	This is a geocoding library which uses the Google Geocoding API.

	Lookups are done either by Address or by Location (latitude and longitude pair)
}
**/
package geocoding

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

//constants
const (
	GOOGLEMAPS_API = "http://maps.googleapis.com/maps/api/geocode/json"
)

//this is the Google repsonse structure we get back from Google
type GoogleResponse struct {
	Results []GoogleResults `json:"results"`
	Status  string          `json:"status"`
}

//these are the various result components we get back from the response from Google
type GoogleResults struct {
	FormattedAddress  string             `json:"formatted_address"`
	AddressComponents []AddressComponent `json:"address_components"`
	Geometry          Geometry           `json:"geometry"`
	Types             []string           `json:"types"`
}

//this is a Point or Location
type Point struct {
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lng"`
}

//the bounds of the map geometry
type Bounds struct {
	NorthEast Point `json:"northeast"`
	NorthWest Point `json:"southwest"`
}

//various items that Google uses for address classifications
type AddressComponent struct {
	LongName  string   `json:"long_name"`
	ShortName string   `json:"short_name"`
	Types     []string `json:"types"`
}

//the map geometry
type Geometry struct {
	Bounds       Bounds `json:"bounds"`
	ViewPort     Bounds `json:"viewport"`
	Location     Point  `json:"location"`
	LocationType string `json:"location_type"`
}

//this is our lookup object
type Lookup struct {
	Address  string
	Location *Point //by using a pointer this allows us to null check the location
}

//search returns an object represenation of the Google response
func (l *Lookup) Search() (res GoogleResponse, err error) {
	var buf []byte
	var resp *http.Response
	var vals = make(url.Values, 0)

	//TODO: we need to add the "key" or api key parameter, otherwise we'll be fighting with
	//the default api request quotas set by Google

	switch {
	case l.Address != "":
		vals.Add("address", l.Address)
		if resp, err = http.Get(GOOGLEMAPS_API + "?" + vals.Encode()); err != nil {
			return
		}
	case l.Location != nil:
		vals.Add("latlng", strconv.FormatFloat(l.Location.Latitude, 'f', 7, 64)+","+
			strconv.FormatFloat(l.Location.Longitude, 'f', 7, 64))
		if resp, err = http.Get(GOOGLEMAPS_API + "?" + vals.Encode()); err != nil {
			return
		}
	default:
		err = errors.New("Must search by either address or location!")
		return
	}

	defer resp.Body.Close()

	if buf, err = ioutil.ReadAll(resp.Body); err != nil {
		return
	}

	if err = json.Unmarshal(buf, &res); err != nil {
		return
	}

	return
}