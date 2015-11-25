package forum_ctlr

import (
	"net/http"
	"strconv"

	"github.com/curt-labs/API/helpers/apicontext"
	"github.com/curt-labs/API/helpers/encoding"
	"github.com/curt-labs/API/helpers/error"
	"github.com/curt-labs/API/models/forum"
	"github.com/go-martini/martini"
)

func GetAllGroups(rw http.ResponseWriter, req *http.Request, enc encoding.Encoder, dtx *apicontext.DataContext) string {
	groups, err := forum.GetAllGroups(dtx)
	if err != nil {
		apierror.GenerateError("Trouble getting all forum groups", err, rw, req)
	}
	return encoding.Must(enc.Encode(groups))
}

func GetGroup(rw http.ResponseWriter, req *http.Request, params martini.Params, enc encoding.Encoder, dtx *apicontext.DataContext) string {
	var err error
	var group forum.Group

	if group.ID, err = strconv.Atoi(params["id"]); err != nil {
		apierror.GenerateError("Trouble getting forum group ID", err, rw, req)
	}
	if err := group.Get(dtx); err != nil {
		apierror.GenerateError("Trouble getting forum group", err, rw, req)
	}
	return encoding.Must(enc.Encode(group))
}

func AddGroup(rw http.ResponseWriter, req *http.Request, params martini.Params, enc encoding.Encoder, dtx *apicontext.DataContext) string {
	group := forum.Group{
		Name:        req.FormValue("name"),
		Description: req.FormValue("description"),
	}

	if err := group.Add(dtx); err != nil {
		apierror.GenerateError("Trouble adding forum group", err, rw, req)
	}

	return encoding.Must(enc.Encode(group))
}

func UpdateGroup(rw http.ResponseWriter, req *http.Request, params martini.Params, enc encoding.Encoder, dtx *apicontext.DataContext) string {
	var err error
	var group forum.Group

	if group.ID, err = strconv.Atoi(params["id"]); err != nil {
		apierror.GenerateError("Trouble getting forum group ID", err, rw, req)
	}

	if err = group.Get(dtx); err != nil {
		apierror.GenerateError("Trouble getting forum group", err, rw, req)
	}

	if req.FormValue("name") != "" {
		group.Name = req.FormValue("name")
	}

	if req.FormValue("description") != "" {
		group.Description = req.FormValue("description")
	}

	if err := group.Update(dtx); err != nil {
		apierror.GenerateError("Trouble updating forum group", err, rw, req)
	}

	return encoding.Must(enc.Encode(group))
}

func DeleteGroup(rw http.ResponseWriter, req *http.Request, params martini.Params, enc encoding.Encoder, dtx *apicontext.DataContext) string {
	var err error
	var group forum.Group

	if group.ID, err = strconv.Atoi(params["id"]); err != nil {
		apierror.GenerateError("Trouble getting forum group ID", err, rw, req)
	}

	if err = group.Delete(dtx); err != nil {
		apierror.GenerateError("Trouble deleting forum group", err, rw, req)
	}

	return encoding.Must(enc.Encode(group))
}
