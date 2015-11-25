package brand_ctlr

import (
	"net/http"
	"strconv"

	"github.com/curt-labs/API/helpers/encoding"
	"github.com/curt-labs/API/helpers/error"
	"github.com/curt-labs/API/models/brand"
	"github.com/go-martini/martini"
)

func GetAllBrands(rw http.ResponseWriter, req *http.Request, enc encoding.Encoder) string {
	brands, err := brand.GetAllBrands()
	if err != nil {
		apierror.GenerateError("Trouble getting all brands", err, rw, req)
	}
	return encoding.Must(enc.Encode(brands))
}

func GetBrand(rw http.ResponseWriter, req *http.Request, params martini.Params, enc encoding.Encoder) string {
	var err error
	var br brand.Brand

	if br.ID, err = strconv.Atoi(params["id"]); err != nil {
		apierror.GenerateError("Trouble getting brand ID", err, rw, req)
	}
	if err := br.Get(); err != nil {
		apierror.GenerateError("Trouble getting brand", err, rw, req)
	}
	return encoding.Must(enc.Encode(br))
}

func CreateBrand(rw http.ResponseWriter, req *http.Request, params martini.Params, enc encoding.Encoder) string {
	br := brand.Brand{
		Name: req.FormValue("name"),
		Code: req.FormValue("code"),
	}

	if err := br.Create(); err != nil {
		apierror.GenerateError("Trouble creating brand", err, rw, req)
	}

	return encoding.Must(enc.Encode(br))
}

func UpdateBrand(rw http.ResponseWriter, req *http.Request, params martini.Params, enc encoding.Encoder) string {
	var err error
	var br brand.Brand

	if br.ID, err = strconv.Atoi(params["id"]); err != nil {
		apierror.GenerateError("Trouble getting brand ID", err, rw, req)
	}

	if err = br.Get(); err != nil {
		apierror.GenerateError("Trouble getting brand", err, rw, req)
	}

	if req.FormValue("name") != "" {
		br.Name = req.FormValue("name")
	}

	if req.FormValue("code") != "" {
		br.Code = req.FormValue("code")
	}

	if err := br.Update(); err != nil {
		apierror.GenerateError("Trouble updating brand", err, rw, req)
	}

	return encoding.Must(enc.Encode(br))
}

func DeleteBrand(rw http.ResponseWriter, req *http.Request, params martini.Params, enc encoding.Encoder) string {
	var err error
	var br brand.Brand

	if br.ID, err = strconv.Atoi(params["id"]); err != nil {
		apierror.GenerateError("Trouble getting brand ID", err, rw, req)
	}

	if err = br.Delete(); err != nil {
		apierror.GenerateError("Trouble deleting brand", err, rw, req)
	}

	return encoding.Must(enc.Encode(br))
}
