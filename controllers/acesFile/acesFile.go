package acesFile

import (
	"encoding/xml"
	"strconv"

	"github.com/curt-labs/API/helpers/encoding"
	"github.com/curt-labs/API/helpers/error"
	"github.com/curt-labs/API/models/acesFile"
	"github.com/curt-labs/API/models/brand"
	"github.com/go-martini/martini"

	"net/http"
)

func GetAcesFile(rw http.ResponseWriter, req *http.Request, enc encoding.Encoder, params martini.Params) string {
	version := params["version"]

	var brandObj brand.Brand
	brandID, err := strconv.Atoi(req.URL.Query().Get("brandID"))
	if err != nil {
		apierror.GenerateError("Invalid brand ID", err, rw, req, http.StatusBadRequest)
		return ""
	}

	brandObj.ID = brandID
	err = brandObj.Get()
	if err != nil {
		apierror.GenerateError("Invalid brand ID", err, rw, req, http.StatusBadRequest)
		return ""
	}

	file, err := acesFile.GetAcesFile(brandObj, version)
	if err != nil {
		apierror.GenerateError("Could not fetch Aces File", err, rw, req, http.StatusBadRequest)
		return ""
	}

	var aces acesFile.Aces

	err = xml.Unmarshal([]byte(file), &aces)
	if err != nil {
		apierror.GenerateError("Could not fetch Aces File", err, rw, req, http.StatusBadRequest)
		return ""
	}

	return encoding.Must(enc.Encode(aces))
}
