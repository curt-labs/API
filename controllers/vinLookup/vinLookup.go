package vinLookup

import (
	"net/http"
	"strconv"

	"github.com/curt-labs/GoAPI/helpers/apicontext"
	"github.com/curt-labs/GoAPI/helpers/encoding"
	"github.com/curt-labs/GoAPI/helpers/error"
	"github.com/curt-labs/GoAPI/models/vinLookup"
	"github.com/go-martini/martini"
)

func GetParts(rw http.ResponseWriter, req *http.Request, enc encoding.Encoder, params martini.Params, dtx *apicontext.DataContext) string {
	vin := params["vin"]

	parts, err := vinLookup.VinPartLookup(vin, dtx)
	if err != nil {
		apierror.GenerateError("Trouble getting parts", err, rw, req)
	}
	//return array of vehicles containing array of parts
	return encoding.Must(enc.Encode(parts))

}
func GetConfigs(rw http.ResponseWriter, req *http.Request, enc encoding.Encoder, params martini.Params) string {
	vin := params["vin"]

	configs, err := vinLookup.GetVehicleConfigs(vin)
	if err != nil {
		apierror.GenerateError("Trouble getting vehicle configurations", err, rw, req)
	}

	return encoding.Must(enc.Encode(configs))
}

func GetPartsFromVehicleID(rw http.ResponseWriter, req *http.Request, enc encoding.Encoder, params martini.Params, dtx *apicontext.DataContext) string {
	vehicleID := params["vehicleID"]
	id, err := strconv.Atoi(vehicleID)
	if err != nil {
		apierror.GenerateError("Trouble getting vehicle ID", err, rw, req)
		return ""
	}
	var v vinLookup.CurtVehicle
	v.ID = id
	parts, err := v.GetPartsFromVehicleConfig(dtx)
	if err != nil {
		apierror.GenerateError("Trouble getting parts from vehicle configuration", err, rw, req)

	}
	return encoding.Must(enc.Encode(parts))
}
