package acesFile

import (
	"bytes"
	"time"

	"github.com/curt-labs/API/helpers/apicontext"
	"github.com/curt-labs/API/helpers/encoding"
	"github.com/curt-labs/API/helpers/error"
	"github.com/curt-labs/API/models/acesFile"
	"github.com/curt-labs/API/models/brand"
	"github.com/go-martini/martini"

	"net/http"
)

func GetAcesFile(rw http.ResponseWriter, req *http.Request, enc encoding.Encoder, params martini.Params, dtx *apicontext.DataContext) string {
	var err error
	version := params["version"]

	var brandObj brand.Brand
	brandObj.ID = dtx.BrandID
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

	rw.Header().Set("Content-Type", "application/xml")

	http.ServeContent(rw, req, "aces.xml", time.Now(), bytes.NewReader([]byte(file)))

	return ""
}
