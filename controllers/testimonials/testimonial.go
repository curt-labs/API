package testimonials

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/curt-labs/GoAPI/helpers/apicontext"
	"github.com/curt-labs/GoAPI/helpers/encoding"
	"github.com/curt-labs/GoAPI/helpers/error"
	"github.com/curt-labs/GoAPI/models/testimonials"
	"github.com/go-martini/martini"
)

func GetAllTestimonials(rw http.ResponseWriter, req *http.Request, enc encoding.Encoder, dtx *apicontext.DataContext) string {
	var err error
	var page int
	var count int
	var randomize bool

	qs := req.URL.Query()

	if qs.Get("page") != "" {
		if pg, err := strconv.Atoi(qs.Get("page")); err == nil {
			page = pg
		}
	}
	if qs.Get("count") != "" {
		if c, err := strconv.Atoi(qs.Get("count")); err == nil {
			count = c
		}
	}

	if qs.Get("randomize") != "" {
		randomize, err = strconv.ParseBool(qs.Get("randomize"))
	}

	tests, err := testimonials.GetAllTestimonials(page, count, randomize, dtx)
	if err != nil {
		apierror.GenerateError("Trouble getting all testimonials", err, rw, req)
	}
	return encoding.Must(enc.Encode(tests))
}

func GetTestimonial(rw http.ResponseWriter, req *http.Request, params martini.Params, enc encoding.Encoder, dtx *apicontext.DataContext) string {
	var err error
	var test testimonials.Testimonial

	if test.ID, err = strconv.Atoi(params["id"]); err != nil {
		apierror.GenerateError("Trouble getting testimonial ID", err, rw, req)
	}
	if err := test.Get(dtx); err != nil {
		apierror.GenerateError("Trouble getting testimonial", err, rw, req)
	}
	return encoding.Must(enc.Encode(test))
}

func Save(rw http.ResponseWriter, req *http.Request, enc encoding.Encoder, params martini.Params, dtx *apicontext.DataContext) string {
	var a testimonials.Testimonial
	var err error
	idStr := params["id"]
	if idStr != "" {
		a.ID, err = strconv.Atoi(idStr)
		err = a.Get(dtx)
		if err != nil {
			apierror.GenerateError("Trouble getting testimonial", err, rw, req)
		}
	}

	//json
	requestBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		apierror.GenerateError("Trouble reading request body for saving testimonial", err, rw, req)
	}
	err = json.Unmarshal(requestBody, &a)
	if err != nil {
		apierror.GenerateError("Trouble unmarshalling request body for saving testimonial", err, rw, req)
	}
	//create or update
	if a.ID > 0 {
		err = a.Update(dtx)
	} else {
		err = a.Create(dtx)
	}

	if err != nil {
		apierror.GenerateError("Trouble saving testimonial", err, rw, req)
	}
	return encoding.Must(enc.Encode(a))
}

func Delete(rw http.ResponseWriter, req *http.Request, enc encoding.Encoder, params martini.Params) string {
	var err error
	var a testimonials.Testimonial

	idStr := params["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		apierror.GenerateError("Trouble getting testimonial ID", err, rw, req)
	}
	a.ID = id
	err = a.Delete()
	if err != nil {
		apierror.GenerateError("Trouble deleting testimonial", err, rw, req)
	}

	return encoding.Must(enc.Encode(a))
}
