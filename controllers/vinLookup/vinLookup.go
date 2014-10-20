package vinLookup

import (
	"github.com/curt-labs/GoAPI/helpers/encoding"
	"github.com/curt-labs/GoAPI/models/vinLookup"
	"github.com/go-martini/martini"
	"log"
	"net/http"
	"strconv"
)

func GetParts(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	vin := params["vin"]

	parts, err := vinLookup.VinPartLookup(vin)
	if err != nil {
		log.Print(err)
		return ""
	}

	return encoding.Must(enc.Encode(parts))

}

func GetConfigs(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	vin := params["vin"]

	configs, err := vinLookup.GetVehicleConfigs(vin)
	if err != nil {
		log.Print(err)
		return ""
	}

	return encoding.Must(enc.Encode(configs))
}

func GetPartsFromVehicleID(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	vehicleID := params["vehicleID"]
	id, err := strconv.Atoi(vehicleID)
	var v vinLookup.CurtVehicle
	v.ID = id
	parts, err := v.GetPartsFromVehicleConfig()
	if err != nil {
		log.Print(err)
		return ""
	}
	return encoding.Must(enc.Encode(parts))

}
