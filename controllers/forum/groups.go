package forum_ctlr

import (
	"net/http"
	"strconv"

	"github.com/curt-labs/GoAPI/helpers/encoding"
	"github.com/curt-labs/GoAPI/models/forum"
	"github.com/go-martini/martini"
)

func GetAllGroups(rw http.ResponseWriter, req *http.Request, enc encoding.Encoder) string {
	groups, err := forum.GetAllGroups()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return err.Error()
	}
	return encoding.Must(enc.Encode(groups))
}

func GetGroup(rw http.ResponseWriter, req *http.Request, params martini.Params, enc encoding.Encoder) string {
	var err error
	var group forum.Group

	if group.ID, err = strconv.Atoi(params["id"]); err != nil {
		http.Error(rw, "Invalid Group ID", http.StatusInternalServerError)
		return "Invalid Group ID"
	}
	if err := group.Get(); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return err.Error()
	}
	return encoding.Must(enc.Encode(group))
}

func AddGroup(rw http.ResponseWriter, req *http.Request, params martini.Params, enc encoding.Encoder) string {
	group := forum.Group{
		Name:        req.FormValue("name"),
		Description: req.FormValue("description"),
	}

	if err := group.Add(); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return err.Error()
	}

	return encoding.Must(enc.Encode(group))
}

func UpdateGroup(rw http.ResponseWriter, req *http.Request, params martini.Params, enc encoding.Encoder) string {
	var err error
	var group forum.Group

	if group.ID, err = strconv.Atoi(params["id"]); err != nil {
		http.Error(rw, "Invalid Group ID", http.StatusInternalServerError)
		return "Invalid Group ID"
	}

	if err = group.Get(); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return err.Error()
	}

	if req.FormValue("name") != "" {
		group.Name = req.FormValue("name")
	}

	if req.FormValue("description") != "" {
		group.Description = req.FormValue("description")
	}

	if err := group.Update(); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return err.Error()
	}

	return encoding.Must(enc.Encode(group))
}

func DeleteGroup(rw http.ResponseWriter, req *http.Request, params martini.Params, enc encoding.Encoder) string {
	var err error
	var group forum.Group

	if group.ID, err = strconv.Atoi(params["id"]); err != nil {
		http.Error(rw, "Invalid Group ID", http.StatusInternalServerError)
		return "Invalid Group ID"
	}

	if err = group.Delete(); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return err.Error()
	}

	return encoding.Must(enc.Encode(group))
}
