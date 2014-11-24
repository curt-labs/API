package brand_ctlr

import (
	"net/http"
	"strconv"

	"github.com/curt-labs/GoAPI/helpers/encoding"
	"github.com/curt-labs/GoAPI/models/brand"
	"github.com/go-martini/martini"
)

func GetAllBrands(rw http.ResponseWriter, req *http.Request, enc encoding.Encoder) string {
	brands, err := brand.GetAllBrands()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return err.Error()
	}
	return encoding.Must(enc.Encode(brands))
}

func GetBrand(rw http.ResponseWriter, req *http.Request, params martini.Params, enc encoding.Encoder) string {
	var err error
	var br brand.Brand

	if br.ID, err = strconv.Atoi(params["id"]); err != nil {
		http.Error(rw, "Invalid Brand ID", http.StatusInternalServerError)
		return "Invalid Brand ID"
	}
	if err := br.Get(); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return err.Error()
	}
	return encoding.Must(enc.Encode(br))
}

func CreateBrand(rw http.ResponseWriter, req *http.Request, params martini.Params, enc encoding.Encoder) string {
	br := brand.Brand{
		Name: req.FormValue("name"),
		Code: req.FormValue("code"),
	}

	if err := br.Create(); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return err.Error()
	}

	return encoding.Must(enc.Encode(br))
}

func UpdateBrand(rw http.ResponseWriter, req *http.Request, params martini.Params, enc encoding.Encoder) string {
	var err error
	var br brand.Brand

	if br.ID, err = strconv.Atoi(params["id"]); err != nil {
		http.Error(rw, "Invalid Brand ID", http.StatusInternalServerError)
		return "Invalid Brand ID"
	}

	if err = br.Get(); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return err.Error()
	}

	if req.FormValue("name") != "" {
		br.Name = req.FormValue("name")
	}

	if req.FormValue("code") != "" {
		br.Code = req.FormValue("code")
	}

	if err := br.Update(); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return err.Error()
	}

	return encoding.Must(enc.Encode(br))
}

func DeleteBrand(rw http.ResponseWriter, req *http.Request, params martini.Params, enc encoding.Encoder) string {
	var err error
	var br brand.Brand

	if br.ID, err = strconv.Atoi(params["id"]); err != nil {
		http.Error(rw, "Invalid Brand ID", http.StatusInternalServerError)
		return "Invalid Brand ID"
	}

	if err = br.Delete(); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return err.Error()
	}

	return encoding.Must(enc.Encode(br))
}
