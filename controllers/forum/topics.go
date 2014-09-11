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
	var topic forum.Topic

	return encoding.Must(enc.Encode(topic))
}

func UpdateTopic(rw http.ResponseWriter, req *http.Request, params martini.Params, enc encoding.Encoder) string {
	var topic forum.Topic

	return encoding.Must(enc.Encode(topic))
}

func DeleteTopic(rw http.ResponseWriter, req *http.Request, params martini.Params, enc encoding.Encoder) string {
	var topic forum.Thread

	return encoding.Must(enc.Encode(topic))
}
