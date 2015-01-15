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

func GetAllTopics(rw http.ResponseWriter, req *http.Request, enc encoding.Encoder, dtx *apicontext.DataContext) string {
	topics, err := forum.GetAllTopics(dtx)
	if err != nil {
		apierror.GenerateError("Trouble getting all forum topics", err, rw, req)
	}
	return encoding.Must(enc.Encode(topics))
}

func GetTopic(rw http.ResponseWriter, req *http.Request, params martini.Params, enc encoding.Encoder, dtx *apicontext.DataContext) string {
	var err error
	var topic forum.Topic

	if topic.ID, err = strconv.Atoi(params["id"]); err != nil {
		apierror.GenerateError("Trouble getting forum topic ID", err, rw, req)
	}
	if err := topic.Get(dtx); err != nil {
		apierror.GenerateError("Trouble getting forum topic", err, rw, req)
	}
	return encoding.Must(enc.Encode(topic))
}

func AddTopic(rw http.ResponseWriter, req *http.Request, params martini.Params, enc encoding.Encoder, dtx *apicontext.DataContext) string {
	var err error
	var topic forum.Topic

	if topic.GroupID, err = strconv.Atoi(req.FormValue("groupID")); err != nil {
		apierror.GenerateError("Trouble getting forum group ID for new topic", err, rw, req)
	}

	if req.FormValue("closed") != "" {
		if topic.Closed, err = strconv.ParseBool(req.FormValue("closed")); err != nil {
			apierror.GenerateError("Trouble adding forum topic -- boolean closed parameter is invalid", err, rw, req)
		}
	}

	topic.Name = req.FormValue("name")
	topic.Description = req.FormValue("description")
	topic.Image = req.FormValue("image")
	topic.Active = true

	if err = topic.Add(); err != nil {
		apierror.GenerateError("Trouble adding forum topic", err, rw, req)
	}

	return encoding.Must(enc.Encode(topic))
}

func UpdateTopic(rw http.ResponseWriter, req *http.Request, params martini.Params, enc encoding.Encoder, dtx *apicontext.DataContext) string {
	var err error
	var topic forum.Topic

	if topic.ID, err = strconv.Atoi(params["id"]); err != nil {
		apierror.GenerateError("Trouble getting forum topic ID", err, rw, req)
	}

	if err = topic.Get(dtx); err != nil {
		apierror.GenerateError("Trouble getting forum topic", err, rw, req)
	}

	if req.FormValue("groupID") != "" {
		if topic.GroupID, err = strconv.Atoi(req.FormValue("groupID")); err != nil {
			apierror.GenerateError("Trouble updating forum topic -- invalid forum group ID", err, rw, req)
		}
	}

	if req.FormValue("closed") != "" {
		if topic.Closed, err = strconv.ParseBool(req.FormValue("closed")); err != nil {
			apierror.GenerateError("Trouble updating forum topic -- boolean closed parameter is invalid", err, rw, req)
		}
	}

	if req.FormValue("active") != "" {
		if topic.Active, err = strconv.ParseBool(req.FormValue("active")); err != nil {
			apierror.GenerateError("Trouble updating forum topic -- boolean active parameter is invalid", err, rw, req)
		}
	}

	if req.FormValue("name") != "" {
		topic.Name = req.FormValue("name")
	}

	if req.FormValue("description") != "" {
		topic.Description = req.FormValue("description")
	}

	if req.FormValue("image") != "" {
		topic.Image = req.FormValue("image")
	}

	if err = topic.Update(); err != nil {
		apierror.GenerateError("Trouble updating forum topic", err, rw, req)
	}

	return encoding.Must(enc.Encode(topic))
}

func DeleteTopic(rw http.ResponseWriter, req *http.Request, params martini.Params, enc encoding.Encoder, dtx *apicontext.DataContext) string {
	var err error
	var topic forum.Topic

	if topic.ID, err = strconv.Atoi(params["id"]); err != nil {
		apierror.GenerateError("Trouble getting forum topic ID", err, rw, req)
	}

	if err = topic.Delete(dtx); err != nil {
		apierror.GenerateError("Trouble deleting forum topic", err, rw, req)
	}

	return encoding.Must(enc.Encode(topic))
}
