package salesrep

import (
	"net/http"
	"strconv"

	"github.com/curt-labs/GoAPI/helpers/encoding"
	"github.com/curt-labs/GoAPI/models/salesrep"
	"github.com/go-martini/martini"
)

func GetAllSalesReps(rw http.ResponseWriter, req *http.Request, enc encoding.Encoder) string {
	reps, err := salesrep.GetAllSalesReps()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return err.Error()
	}
	return encoding.Must(enc.Encode(reps))
}

func GetSalesRep(rw http.ResponseWriter, req *http.Request, params martini.Params, enc encoding.Encoder) string {
	var err error
	var rep salesrep.SalesRep

	if rep.ID, err = strconv.Atoi(params["id"]); err != nil {
		http.Error(rw, "Invalid SalesRep ID", http.StatusInternalServerError)
		return "Invalid SalesRep ID"
	}
	if err := rep.Get(); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return err.Error()
	}
	return encoding.Must(enc.Encode(rep))
}

func AddSalesRep(rw http.ResponseWriter, req *http.Request, params martini.Params, enc encoding.Encoder) string {
	rep := salesrep.SalesRep{
		Name: req.FormValue("name"),
		Code: req.FormValue("code"),
	}

	if err := rep.Add(); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return err.Error()
	}

	return encoding.Must(enc.Encode(rep))
}

func UpdateSalesRep(rw http.ResponseWriter, req *http.Request, params martini.Params, enc encoding.Encoder) string {
	var err error
	var rep salesrep.SalesRep

	if rep.ID, err = strconv.Atoi(params["id"]); err != nil {
		http.Error(rw, "Invalid SalesRep ID", http.StatusInternalServerError)
		return "Invalid SalesRep ID"
	}

	if err = rep.Get(); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return err.Error()
	}

	if req.FormValue("name") != "" {
		rep.Name = req.FormValue("name")
	}

	if req.FormValue("code") != "" {
		rep.Code = req.FormValue("code")
	}

	if err := rep.Update(); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return err.Error()
	}

	return encoding.Must(enc.Encode(rep))
}

func DeleteSalesRep(rw http.ResponseWriter, req *http.Request, params martini.Params, enc encoding.Encoder) string {
	var err error
	var rep salesrep.SalesRep

	if rep.ID, err = strconv.Atoi(params["id"]); err != nil {
		http.Error(rw, "Invalid SalesRep ID", http.StatusInternalServerError)
		return "Invalid SalesRep ID"
	}

	if err = rep.Delete(); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return err.Error()
	}

	return encoding.Must(enc.Encode(rep))
}
