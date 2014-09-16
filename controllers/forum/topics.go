package forum_ctlr

import (
	"net/http"
	"strconv"

	"github.com/curt-labs/GoAPI/helpers/encoding"
	"github.com/curt-labs/GoAPI/models/forum"
	"github.com/go-martini/martini"
)

func GetAllTopics(rw http.ResponseWriter, req *http.Request, enc encoding.Encoder) string {
	topics, err := forum.GetAllTopics()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return err.Error()
	}
	return encoding.Must(enc.Encode(topics))
}

func GetTopic(rw http.ResponseWriter, req *http.Request, params martini.Params, enc encoding.Encoder) string {
	var err error
	var topic forum.Topic

	if topic.ID, err = strconv.Atoi(params["id"]); err != nil {
		http.Error(rw, "Invalid Topic ID", http.StatusInternalServerError)
		return "Invalid Topic ID"
	}
	if err := topic.Get(); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return err.Error()
	}
	return encoding.Must(enc.Encode(topic))
}

func AddTopic(rw http.ResponseWriter, req *http.Request, params martini.Params, enc encoding.Encoder) string {
	var err error
	var topic forum.Topic

	if topic.GroupID, err = strconv.Atoi(req.FormValue("groupID")); err != nil {
		http.Error(rw, "Invalid Group ID", http.StatusInternalServerError)
		return "Invalid Group ID"
	}

	if req.FormValue("closed") != "" {
		if topic.Closed, err = strconv.ParseBool(req.FormValue("closed")); err != nil {
			http.Error(rw, "Invalid boolean for topic.Closed", http.StatusInternalServerError)
			return "Invalid boolean for topic.Closed"
		}
	}

	topic.Name = req.FormValue("name")
	topic.Description = req.FormValue("description")
	topic.Image = req.FormValue("image")
	topic.Active = true

	if err = topic.Add(); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return err.Error()
	}

	return encoding.Must(enc.Encode(topic))
}

func UpdateTopic(rw http.ResponseWriter, req *http.Request, params martini.Params, enc encoding.Encoder) string {
	var err error
	var topic forum.Topic

	if topic.ID, err = strconv.Atoi(params["id"]); err != nil {
		http.Error(rw, "Invalid Topic ID", http.StatusInternalServerError)
		return "Invalid Topic ID"
	}

	if err = topic.Get(); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return err.Error()
	}

	if req.FormValue("groupID") != "" {
		if topic.GroupID, err = strconv.Atoi(req.FormValue("groupID")); err != nil {
			http.Error(rw, "Invalid Group ID", http.StatusInternalServerError)
			return "Invalid Group ID"
		}
	}

	if req.FormValue("closed") != "" {
		if topic.Closed, err = strconv.ParseBool(req.FormValue("closed")); err != nil {
			http.Error(rw, "Invalid boolean for topic.Closed", http.StatusInternalServerError)
			return "Invalid boolean for topic.Closed"
		}
	}

	if req.FormValue("active") != "" {
		if topic.Active, err = strconv.ParseBool(req.FormValue("active")); err != nil {
			http.Error(rw, "Invalid boolean for topic.Active", http.StatusInternalServerError)
			return "Invalid boolean for topic.Active"
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
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return err.Error()
	}

	return encoding.Must(enc.Encode(topic))
}

func DeleteTopic(rw http.ResponseWriter, req *http.Request, params martini.Params, enc encoding.Encoder) string {
	var err error
	var topic forum.Topic

	if topic.ID, err = strconv.Atoi(params["id"]); err != nil {
		http.Error(rw, "Invalid Topic ID", http.StatusInternalServerError)
		return "Invalid Topic ID"
	}

	if err = topic.Delete(); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return err.Error()
	}

	return encoding.Must(enc.Encode(topic))
}
