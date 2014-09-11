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

func GetAllTopics(rw http.ResponseWriter, req *http.Request, enc encoding.Encoder) string {
	topics, err := forum.GetAllTopics()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return err.Error()
	}
	return encoding.Must(enc.Encode(topics))
}

func GetAllThreads(rw http.ResponseWriter, req *http.Request, enc encoding.Encoder) string {
	threads, err := forum.GetAllThreads()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return err.Error()
	}
	return encoding.Must(enc.Encode(threads))
}

func GetAllPosts(rw http.ResponseWriter, req *http.Request, enc encoding.Encoder) string {
	posts, err := forum.GetAllPosts()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return err.Error()
	}
	return encoding.Must(enc.Encode(posts))
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

func GetPost(rw http.ResponseWriter, req *http.Request, params martini.Params, enc encoding.Encoder) string {
	var err error
	var post forum.Post

	if post.ID, err = strconv.Atoi(params["id"]); err != nil {
		http.Error(rw, "Invalid Post ID", http.StatusInternalServerError)
		return "Invalid Post ID"
	}

	if err := post.Get(); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return err.Error()
	}
	return encoding.Must(enc.Encode(post))
}
