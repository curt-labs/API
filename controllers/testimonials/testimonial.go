package testimonials

import (
	"net/http"
	"strconv"

	"github.com/curt-labs/GoAPI/helpers/encoding"
	"github.com/curt-labs/GoAPI/models/testimonials"
	"github.com/go-martini/martini"
)

func GetAllTestimonials(rw http.ResponseWriter, req *http.Request, enc encoding.Encoder) string {
	tests, err := testimonials.GetAllTestimonials()
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
