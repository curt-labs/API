package testimonials

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/curt-labs/GoAPI/helpers/apicontext"
	"github.com/curt-labs/GoAPI/helpers/encoding"
	"github.com/curt-labs/GoAPI/models/testimonials"
	"github.com/go-martini/martini"
)

func GetAllTestimonials(rw http.ResponseWriter, req *http.Request, enc encoding.Encoder, dtx *apicontext.DataContext) string {
	var err error
	var page int
	var count int
	var randomize bool

	qs := req.URL.Query()

	log.Println(dtx.BrandID) // example of how to use the data context.

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

	tests, err := testimonials.GetAllTestimonials(page, count, randomize)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return err.Error()
	}
	return encoding.Must(enc.Encode(tests))
}

func GetTestimonial(rw http.ResponseWriter, req *http.Request, params martini.Params, enc encoding.Encoder) string {
	var err error
	var test testimonials.Testimonial

	if test.ID, err = strconv.Atoi(params["id"]); err != nil {
		http.Error(rw, "Invalid Testimonial ID", http.StatusInternalServerError)
		return "Invalid Testimonial ID"
	}
	if err := test.Get(); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return err.Error()
	}
	return encoding.Must(enc.Encode(test))
}

func Save(rw http.ResponseWriter, req *http.Request, enc encoding.Encoder, params martini.Params) string {
	var a testimonials.Testimonial
	var err error
	idStr := params["id"]
	if idStr != "" {
		a.ID, err = strconv.Atoi(idStr)
		err = a.Get()
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return ""
		}
	}

	//json
	requestBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return encoding.Must(enc.Encode(false))
	}
	err = json.Unmarshal(requestBody, &a)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return encoding.Must(enc.Encode(false))
	}
	//create or update
	if a.ID > 0 {
		err = a.Update()
	} else {
		err = a.Create()
	}

	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return ""
	}
	return encoding.Must(enc.Encode(a))
}

func Delete(rw http.ResponseWriter, req *http.Request, enc encoding.Encoder, params martini.Params) string {
	var err error
	var a testimonials.Testimonial

	idStr := params["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return ""
	}
	a.ID = id
	err = a.Delete()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return ""
	}

	return encoding.Must(enc.Encode(a))
}
