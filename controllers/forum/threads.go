package forum_ctlr

import (
	"net/http"
	"strconv"

	"github.com/curt-labs/GoAPI/helpers/encoding"
	"github.com/curt-labs/GoAPI/models/forum"
	"github.com/go-martini/martini"
)

func GetAllThreads(rw http.ResponseWriter, req *http.Request, enc encoding.Encoder) string {
	threads, err := forum.GetAllThreads()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return err.Error()
	}
	return encoding.Must(enc.Encode(threads))
}

func GetThread(rw http.ResponseWriter, req *http.Request, params martini.Params, enc encoding.Encoder) string {
	var err error
	var thread forum.Thread

	if thread.ID, err = strconv.Atoi(params["id"]); err != nil {
		http.Error(rw, "Invalid Thread ID", http.StatusInternalServerError)
		return "Invalid Thread ID"
	}

	if err := thread.Get(); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return err.Error()
	}
	return encoding.Must(enc.Encode(thread))
}

func DeleteThread(rw http.ResponseWriter, req *http.Request, params martini.Params, enc encoding.Encoder) string {
	var err error
	var thread forum.Thread

	if thread.ID, err = strconv.Atoi(params["id"]); err != nil {
		http.Error(rw, "Invalid Thread ID", http.StatusInternalServerError)
		return "Invalid Thread ID"
	}

	if err = thread.Delete(); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return err.Error()
	}

	return encoding.Must(enc.Encode(thread))
}
