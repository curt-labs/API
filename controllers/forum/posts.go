package forum_ctlr

import (
	"net/http"
	"strconv"

	"github.com/curt-labs/GoAPI/helpers/encoding"
	"github.com/curt-labs/GoAPI/models/forum"
	"github.com/go-martini/martini"
)

func GetAllPosts(rw http.ResponseWriter, req *http.Request, enc encoding.Encoder) string {
	posts, err := forum.GetAllPosts()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return err.Error()
	}
	return encoding.Must(enc.Encode(posts))
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

func AddPost(rw http.ResponseWriter, req *http.Request, params martini.Params, enc encoding.Encoder) string {
	var post forum.Post

	return encoding.Must(enc.Encode(post))
}

func UpdatePost(rw http.ResponseWriter, req *http.Request, params martini.Params, enc encoding.Encoder) string {
	var post forum.Post

	return encoding.Must(enc.Encode(post))
}

func DeletePost(rw http.ResponseWriter, req *http.Request, params martini.Params, enc encoding.Encoder) string {
	var err error
	var post forum.Post

	if post.ID, err = strconv.Atoi(params["id"]); err != nil {
		http.Error(rw, "Invalid Post ID", http.StatusInternalServerError)
		return "Invalid Post ID"
	}

	if err = post.Delete(); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return err.Error()
	}

	return encoding.Must(enc.Encode(post))
}
