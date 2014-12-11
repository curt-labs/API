package faq_controller

import (
	"github.com/curt-labs/GoAPI/helpers/apicontext"
	"github.com/curt-labs/GoAPI/helpers/encoding"
	"github.com/curt-labs/GoAPI/helpers/sortutil"
	"github.com/curt-labs/GoAPI/models/faq"
	"github.com/go-martini/martini"
	"net/http"
	"strconv"
	"strings"
)

func GetAll(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder, dtx *apicontext.DataContext) string {
	var fs faq_model.Faqs
	var err error

	fs, err = faq_model.GetAll(dtx)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return err.Error()
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
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return err.Error()
		}
	} else {
		f.ID, err = strconv.Atoi(params["id"])
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return err.Error()
		}
	}

	err = f.Get(dtx)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return err.Error()
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
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return err.Error()
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
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return err.Error()
		}
	} else {
		f.ID, err = strconv.Atoi(params["id"])
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return err.Error()
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
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return err.Error()
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
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return err.Error()
		}
	} else {
		f.ID, err = strconv.Atoi(params["id"])
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return err.Error()
		}
	}
	err = f.Delete()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return err.Error()
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
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return err.Error()
	}
	return encoding.Must(enc.Encode(l))
}
