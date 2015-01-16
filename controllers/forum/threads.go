package forum_ctlr

import (
	"net/http"
	"strconv"

	"github.com/curt-labs/GoAPI/helpers/apicontext"
	"github.com/curt-labs/GoAPI/helpers/encoding"
	"github.com/curt-labs/GoAPI/helpers/error"
	"github.com/curt-labs/GoAPI/models/forum"
	"github.com/go-martini/martini"
)

func GetAllThreads(rw http.ResponseWriter, req *http.Request, enc encoding.Encoder, dtx *apicontext.DataContext) string {
	threads, err := forum.GetAllThreads(dtx)
	if err != nil {
		apierror.GenerateError("Trouble getting all forum threads", err, rw, req)
	}
	return encoding.Must(enc.Encode(threads))
}

func GetThread(rw http.ResponseWriter, req *http.Request, params martini.Params, enc encoding.Encoder, dtx *apicontext.DataContext) string {
	var err error
	var thread forum.Thread

	if thread.ID, err = strconv.Atoi(params["id"]); err != nil {
		apierror.GenerateError("Trouble getting forum thread ID", err, rw, req)
	}
	if err := thread.Get(dtx); err != nil {
		apierror.GenerateError("Trouble getting forum thread", err, rw, req)
	}
	return encoding.Must(enc.Encode(thread))
}

func DeleteThread(rw http.ResponseWriter, req *http.Request, params martini.Params, enc encoding.Encoder) string {
	var err error
	var thread forum.Thread

	if thread.ID, err = strconv.Atoi(params["id"]); err != nil {
		apierror.GenerateError("Trouble getting forum thread ID", err, rw, req)
	}

	if err = thread.Delete(); err != nil {
		apierror.GenerateError("Trouble deleting forum thread", err, rw, req)
	}

	return encoding.Must(enc.Encode(thread))
}
