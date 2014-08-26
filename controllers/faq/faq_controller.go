package faq_controller

import (
	"github.com/curt-labs/GoAPI/helpers/encoding"
	"github.com/curt-labs/GoAPI/helpers/pagination"
	"github.com/curt-labs/GoAPI/helpers/sortutil"
	"github.com/curt-labs/GoAPI/models/faq"
	"github.com/go-martini/martini"
	"net/http"
	"strconv"
	"strings"
)

func GetAll(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder) string {
	var fs faq_model.Faqs
	var err error

	fs, err = faq_model.GetAll()
	if err != nil {
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

func Get(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	var f faq_model.Faq
	var err error

	idStr := r.FormValue("id")
	if idStr != "" {
		f.ID, err = strconv.Atoi(idStr)
		if err != nil {
			return err.Error()
		}
	} else {
		f.ID, err = strconv.Atoi(params["id"])
		if err != nil {
			return err.Error()
		}
	}

	err = f.Get()
	if err != nil {
		return err.Error()
	}
	return encoding.Must(enc.Encode(f))
}

func Create(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder) string {
	var f faq_model.Faq
	var err error

	f.Question = r.FormValue("question")
	f.Answer = r.FormValue("answer")

	err = f.Create()
	if err != nil {
		return err.Error()
	}
	return encoding.Must(enc.Encode(f))
}

func Update(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	var f faq_model.Faq
	var err error

	idStr := r.FormValue("id")
	if idStr != "" {
		f.ID, err = strconv.Atoi(idStr)
		if err != nil {
			return err.Error()
		}
	} else {
		f.ID, err = strconv.Atoi(params["id"])
		if err != nil {
			return err.Error()
		}
	}
	f.Get()
	question := r.FormValue("question")
	answer := r.FormValue("answer")
	if question != "" {
		f.Question = question
	}
	if answer != "" {
		f.Answer = answer
	}

	err = f.Update()
	if err != nil {
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
			return err.Error()
		}
	} else {
		f.ID, err = strconv.Atoi(params["id"])
		if err != nil {
			return err.Error()
		}
	}
	err = f.Delete()
	if err != nil {
		return err.Error()
	}
	return encoding.Must(enc.Encode(f))
}

func GetQuestions(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder) string {
	var l pagination.Objects
	var err error
	page := r.FormValue("page")
	results := r.FormValue("results")

	l, err = faq_model.GetQuestions(page, results)
	if err != nil {
		return err.Error()
	}
	return encoding.Must(enc.Encode(l))
}

func GetAnswers(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder) string {
	var l pagination.Objects
	var err error
	page := r.FormValue("page")
	results := r.FormValue("results")

	l, err = faq_model.GetAnswers(page, results)
	if err != nil {
		return err.Error()
	}
	return encoding.Must(enc.Encode(l))
}

func Search(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder) string {
	var err error

	question := r.FormValue("question")
	answer := r.FormValue("answer")
	page := r.FormValue("page")
	results := r.FormValue("results")

	l, err := faq_model.Search(question, answer, page, results)
	if err != nil {
		return err.Error()
	}
	return encoding.Must(enc.Encode(l))
}
