package part_ctlr

// import (
// 	"encoding/json"
// 	"errors"
// 	"io/ioutil"
// 	"net/http"
// 	"strconv"

// 	"github.com/curt-labs/API/helpers/apicontext"
// 	"github.com/curt-labs/API/helpers/encoding"
// 	"github.com/curt-labs/API/helpers/error"
// 	"github.com/curt-labs/API/models/products"
// 	"github.com/go-martini/martini"
// )

// func GetAllReviews(rw http.ResponseWriter, req *http.Request, params martini.Params, enc encoding.Encoder, dtx *apicontext.DataContext) string {
// 	revs, err := products.GetAllReviews(dtx)
// 	if err != nil {
// 		apierror.GenerateError("Trouble getting all part reviews", err, rw, req)
// 		return ""
// 	}

// 	return encoding.Must(enc.Encode(revs))
// }
// func GetReview(rw http.ResponseWriter, req *http.Request, params martini.Params, enc encoding.Encoder, dtx *apicontext.DataContext) string {
// 	var rev products.Review
// 	var err error

// 	if params["id"] == "" {
// 		err = errors.New("Missing review ID parameter in query string")
// 		apierror.GenerateError("Trouble getting review ID", err, rw, req)
// 		return ""
// 	}

// 	if rev.Id, err = strconv.Atoi(params["id"]); err != nil {
// 		apierror.GenerateError("Trouble getting review ID", err, rw, req)
// 		return ""
// 	}

// 	if err = rev.Get(dtx); err != nil {
// 		apierror.GenerateError("Trouble getting review", err, rw, req)
// 		return ""
// 	}

// 	return encoding.Must(enc.Encode(rev))
// }

// func SaveReview(rw http.ResponseWriter, req *http.Request, enc encoding.Encoder, params martini.Params, dtx *apicontext.DataContext) string {
// 	var rev products.Review
// 	var err error

// 	idStr := params["id"]
// 	if idStr != "" {
// 		rev.Id, err = strconv.Atoi(idStr)
// 		if err != nil {
// 			apierror.GenerateError("Trouble getting review ID", err, rw, req)
// 			return ""
// 		}
// 		err = rev.Get(dtx)
// 		if err != nil {
// 			apierror.GenerateError("Trouble getting review", err, rw, req)
// 			return ""
// 		}
// 	}

// 	//json
// 	requestBody, err := ioutil.ReadAll(req.Body)
// 	if err != nil {
// 		apierror.GenerateError("Trouble reading request body while saving review", err, rw, req)
// 		return ""
// 	}
// 	err = json.Unmarshal(requestBody, &rev)
// 	if err != nil {
// 		apierror.GenerateError("Trouble unmarshalling json request body while saving review", err, rw, req)
// 		return ""
// 	}

// 	//create or update
// 	if rev.Id > 0 {
// 		err = rev.Update(dtx)
// 	} else {
// 		err = rev.Create(dtx)
// 	}

// 	if err != nil {
// 		msg := "Trouble creating review"
// 		if rev.Id > 0 {
// 			msg = "Trouble updating review"
// 		}
// 		apierror.GenerateError(msg, err, rw, req)
// 		return ""
// 	}
// 	return encoding.Must(enc.Encode(rev))
// }

// func DeleteReview(rw http.ResponseWriter, r *http.Request, params martini.Params, enc encoding.Encoder, dtx *apicontext.DataContext) string {
// 	var rev products.Review
// 	var err error

// 	if params["id"] == "" {
// 		err = errors.New("Missing review ID parameter in query string")
// 		apierror.GenerateError("Trouble getting review ID", err, rw, r)
// 		return ""
// 	}

// 	if rev.Id, err = strconv.Atoi(params["id"]); err != nil {
// 		apierror.GenerateError("Trouble getting review ID", err, rw, r)
// 		return ""
// 	}

// 	if err = rev.Delete(dtx); err != nil {
// 		apierror.GenerateError("Trouble deleting review", err, rw, r)
// 		return ""
// 	}

// 	return encoding.Must(enc.Encode(rev))
// }
