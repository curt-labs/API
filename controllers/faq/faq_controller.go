package faq_controller

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/curt-labs/API/helpers/apicontext"
	"github.com/curt-labs/API/helpers/encoding"
	"github.com/curt-labs/API/helpers/error"
	"github.com/curt-labs/API/helpers/sortutil"
	"github.com/curt-labs/API/models/faq"
	"github.com/go-martini/martini"
)

func GetAll(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder, dtx *apicontext.DataContext) string {
	var fs faq_model.Faqs
	var err error

	fs, err = faq_model.GetAll(dtx)
	if err != nil {
		apierror.GenerateError("Trouble getting all faqs", err, rw, r)
	}

	sort := r.FormValue("sort")
	direction := r.FormValue("direction")
	if sort != "" {
		if strings.ContainsAny(direction, "esc") {
			sortutil.DescByField(fs, sort)
		} else {
			sortutil.AscByField(fs, sort)
		}

	}

	return encoding.Must(enc.Encode(fs))
}

func Get(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params, dtx *apicontext.DataContext) string {
	var f faq_model.Faq
	var err error

	idStr := r.FormValue("id")
	if idStr != "" {
		f.ID, err = strconv.Atoi(idStr)
		if err != nil {
			apierror.GenerateError("Trouble getting faq ID", err, rw, r)
		}
	} else {
		f.ID, err = strconv.Atoi(params["id"])
		if err != nil {
			apierror.GenerateError("Trouble getting faq ID", err, rw, r)
		}
	}

	err = f.Get(dtx)
	if err != nil {
		apierror.GenerateError("Trouble getting faq", err, rw, r)
	}

	return encoding.Must(enc.Encode(f))
}

func Create(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder, dtx *apicontext.DataContext) string {
	var f faq_model.Faq
	var err error

	f.Question = r.FormValue("question")
	f.Answer = r.FormValue("answer")

	err = f.Create(dtx)
	if err != nil {
		apierror.GenerateError("Trouble creating faq", err, rw, r)
	}
	return encoding.Must(enc.Encode(f))
}

func Update(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params, dtx *apicontext.DataContext) string {
	var f faq_model.Faq
	var err error

	idStr := r.FormValue("id")
	if idStr != "" {
		f.ID, err = strconv.Atoi(idStr)
		if err != nil {
			apierror.GenerateError("Trouble getting faq ID", err, rw, r)
		}
	} else {
		f.ID, err = strconv.Atoi(params["id"])
		if err != nil {
			apierror.GenerateError("Trouble getting faq ID", err, rw, r)
		}
	}
	f.Get(dtx)
	question := r.FormValue("question")
	answer := r.FormValue("answer")
	if question != "" {
		f.Question = question
	}
	if answer != "" {
		f.Answer = answer
	}

	err = f.Update(dtx)
	if err != nil {
		apierror.GenerateError("Trouble updating faq", err, rw, r)
	}
	return encoding.Must(enc.Encode(f))
}

func Delete(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	var f faq_model.Faq
	var err error

	idStr := r.FormValue("id")
	if idStr != "" {
		f.ID, err = strconv.Atoi(idStr)
		if err != nil {
			apierror.GenerateError("Trouble getting faq ID", err, rw, r)
		}
	} else {
		f.ID, err = strconv.Atoi(params["id"])
		if err != nil {
			apierror.GenerateError("Trouble getting faq ID", err, rw, r)
		}
	}
	err = f.Delete()
	if err != nil {
		apierror.GenerateError("Trouble deleting faq", err, rw, r)
	}
	return encoding.Must(enc.Encode(f))
}

func Search(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder, dtx *apicontext.DataContext) string {
	var err error

	question := r.FormValue("question")
	answer := r.FormValue("answer")
	page := r.FormValue("page")
	results := r.FormValue("results")

	l, err := faq_model.Search(dtx, question, answer, page, results)
	if err != nil {
		apierror.GenerateError("Trouble searching faq", err, rw, r)
	}
	return encoding.Must(enc.Encode(l))
}
