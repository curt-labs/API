package vinLookup

import (
	"github.com/curt-labs/GoAPI/helpers/encoding"
	"github.com/curt-labs/GoAPI/models/vinLookup"
	"github.com/go-martini/martini"
	"net/http"
)

func GetVehicleConfigs(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	vin := params["vin"]
	v, err := vinLookup.Lookup(vin)
	if err != nil {
		http.Error(w, err.Error(), http.StatusExpectationFailed)
		return ""
	}
	return encoding.Must(enc.Encode(v))
}

func GetParts(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	// vin := params["vin"]
	// v, err := vinLookup.Lookup(vin)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusExpectationFailed)
	// 	return ""
	// }
	v := "More to come"
	return encoding.Must(enc.Encode(v))

}
