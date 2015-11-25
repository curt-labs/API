package salesrep

import (
	"net/http"
	"strconv"

	"github.com/curt-labs/API/helpers/encoding"
	"github.com/curt-labs/API/helpers/error"
	"github.com/curt-labs/API/models/salesrep"
	"github.com/go-martini/martini"
)

func GetAllSalesReps(rw http.ResponseWriter, req *http.Request, enc encoding.Encoder) string {
	reps, err := salesrep.GetAllSalesReps()
	if err != nil {
		apierror.GenerateError("Trouble getting all sales reps", err, rw, req)
	}
	return encoding.Must(enc.Encode(reps))
}

func GetSalesRep(rw http.ResponseWriter, req *http.Request, params martini.Params, enc encoding.Encoder) string {
	var err error
	var rep salesrep.SalesRep

	if rep.ID, err = strconv.Atoi(params["id"]); err != nil {
		apierror.GenerateError("Trouble getting sales rep ID", err, rw, req)
	}
	if err := rep.Get(); err != nil {
		apierror.GenerateError("Trouble getting sales rep", err, rw, req)
	}
	return encoding.Must(enc.Encode(rep))
}

func AddSalesRep(rw http.ResponseWriter, req *http.Request, params martini.Params, enc encoding.Encoder) string {
	rep := salesrep.SalesRep{
		Name: req.FormValue("name"),
		Code: req.FormValue("code"),
	}

	if err := rep.Add(); err != nil {
		apierror.GenerateError("Trouble adding sales rep", err, rw, req)
	}

	return encoding.Must(enc.Encode(rep))
}

func UpdateSalesRep(rw http.ResponseWriter, req *http.Request, params martini.Params, enc encoding.Encoder) string {
	var err error
	var rep salesrep.SalesRep

	if rep.ID, err = strconv.Atoi(params["id"]); err != nil {
		apierror.GenerateError("Trouble getting sales rep ID", err, rw, req)
	}

	if err = rep.Get(); err != nil {
		apierror.GenerateError("Trouble getting sales rep", err, rw, req)
	}

	if req.FormValue("name") != "" {
		rep.Name = req.FormValue("name")
	}

	if req.FormValue("code") != "" {
		rep.Code = req.FormValue("code")
	}

	if err := rep.Update(); err != nil {
		apierror.GenerateError("Trouble updating sales rep", err, rw, req)
	}

	return encoding.Must(enc.Encode(rep))
}

func DeleteSalesRep(rw http.ResponseWriter, req *http.Request, params martini.Params, enc encoding.Encoder) string {
	var err error
	var rep salesrep.SalesRep

	if rep.ID, err = strconv.Atoi(params["id"]); err != nil {
		apierror.GenerateError("Trouble getting sales rep ID", err, rw, req)
	}

	if err = rep.Delete(); err != nil {
		apierror.GenerateError("Trouble deleting sales rep", err, rw, req)
	}

	return encoding.Must(enc.Encode(rep))
}
