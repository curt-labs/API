package part_ctlr

import (
	"encoding/json"
	"github.com/curt-labs/GoAPI/helpers/apicontext"
	"github.com/curt-labs/GoAPI/helpers/encoding"
	"github.com/curt-labs/GoAPI/models/products"
	"github.com/go-martini/martini"
	"io/ioutil"
	// "github.com/ninnemana/analytics-go"
	// "log"
	"net/http"
	"strconv"
)

func GetAllReviews(w http.ResponseWriter, r *http.Request, params martini.Params, enc encoding.Encoder, dtx *apicontext.DataContext) string {
	revs, err := products.GetAllReviews(dtx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	return encoding.Must(enc.Encode(revs))
}
func GetReview(w http.ResponseWriter, r *http.Request, params martini.Params, enc encoding.Encoder, dtx *apicontext.DataContext) string {
	id, err := strconv.Atoi(params["id"])
	var rev products.Review
	rev.Id = id

	err = rev.Get(dtx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	return encoding.Must(enc.Encode(rev))
}

func SaveReview(rw http.ResponseWriter, req *http.Request, enc encoding.Encoder, params martini.Params, dtx *apicontext.DataContext) string {
	var r products.Review

	var err error
	idStr := params["id"]
	if idStr != "" {
		r.Id, err = strconv.Atoi(idStr)
		err = r.Get(dtx)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return ""
		}
	}

	//json
	requestBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return ""
	}
	err = json.Unmarshal(requestBody, &r)
	if err != nil {

		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return ""
	}

	//create or update
	if r.Id > 0 {
		err = r.Update()
	} else {
		err = r.Create()
	}

	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return ""
	}
	return encoding.Must(enc.Encode(r))
}

func DeleteReview(w http.ResponseWriter, r *http.Request, params martini.Params, enc encoding.Encoder) string {
	id, err := strconv.Atoi(params["id"])
	var rev products.Review
	rev.Id = id

	err = rev.Delete()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	return encoding.Must(enc.Encode(rev))
}
