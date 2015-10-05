package showcase

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/curt-labs/GoAPI/helpers/apicontext"
	"github.com/curt-labs/GoAPI/helpers/encoding"
	"github.com/curt-labs/GoAPI/helpers/error"
	"github.com/curt-labs/GoAPI/models/showcase"
	"github.com/go-martini/martini"
)

func GetAllShowcases(rw http.ResponseWriter, req *http.Request, enc encoding.Encoder, dtx *apicontext.DataContext) string {
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

	shows, err := showcase.GetAllShowcases(page, count, randomize, dtx)
	if err != nil {
		apierror.GenerateError("Trouble getting all", err, rw, req)
	}
	return encoding.Must(enc.Encode(shows))
}

func GetShowcase(rw http.ResponseWriter, req *http.Request, params martini.Params, enc encoding.Encoder, dtx *apicontext.DataContext) string {
	var err error
	var show showcase.Showcase

	if show.ID, err = strconv.Atoi(params["id"]); err != nil {
		apierror.GenerateError("Trouble getting testimonial ID", err, rw, req)
	}
	if err := show.Get(dtx); err != nil {
		apierror.GenerateError("Trouble getting testimonial", err, rw, req)
	}
	return encoding.Must(enc.Encode(show))
}

func Save(rw http.ResponseWriter, req *http.Request, enc encoding.Encoder, params martini.Params, dtx *apicontext.DataContext) string {
	var show showcase.Showcase
	var err error
	idStr := params["id"]
	if idStr != "" {
		show.ID, err = strconv.Atoi(idStr)
		err = show.Get(dtx)
		if err != nil {
			apierror.GenerateError("Trouble getting testimonial", err, rw, req)
		}
	}
	//json
	requestBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		apierror.GenerateError("Trouble reading request body for saving testimonial", err, rw, req)
	}
	err = json.Unmarshal(requestBody, &show)
	if err != nil {
		apierror.GenerateError("Trouble unmarshalling request body for saving testimonial", err, rw, req)
	}
	//create or update
	if show.ID > 0 {
		err = show.Update()
	} else {
		err = show.Create()
	}

	if err != nil {
		apierror.GenerateError("Trouble saving testimonial", err, rw, req)
	}
	return encoding.Must(enc.Encode(show))
}

func Delete(rw http.ResponseWriter, req *http.Request, enc encoding.Encoder, params martini.Params) string {
	var err error
	var a showcase.Showcase

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
