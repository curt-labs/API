package vinLookup

import (
	"github.com/curt-labs/GoAPI/helpers/encoding"
	"github.com/curt-labs/GoAPI/models/vinLookup"
	"github.com/go-martini/martini"
	"net/http"
	"strconv"
)

func GetParts(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	vin := params["vin"]

	parts, err := vinLookup.VinPartLookup(vin)
	if err != nil {
		return err.Error()
	}
	//return array of vehicles containing array of parts
	return encoding.Must(enc.Encode(parts))

}
func GetConfigs(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	vin := params["vin"]

	configs, err := vinLookup.GetVehicleConfigs(vin)
	if err != nil {
		return err.Error()

	}

	return encoding.Must(enc.Encode(configs))
}

func GetPartsFromVehicleID(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	vehicleID := params["vehicleID"]
	id, err := strconv.Atoi(vehicleID)
	if err != nil {
		return ""
	}
	var v vinLookup.CurtVehicle
	v.ID = id
	parts, err := v.GetPartsFromVehicleConfig()
	if err != nil {
		return err.Error()

	}
	return encoding.Must(enc.Encode(parts))

}
