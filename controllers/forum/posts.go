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
	var err error
	var post forum.Post

	//we add a post either by topic (new post) or by identifying the parentID (replying to an existing post)
	if req.FormValue("topicID") == "" && req.FormValue("parentID") == "" {
		http.Error(rw, "Missing topic ID or parent post ID", http.StatusInternalServerError)
		return "Missing topic ID or parent post ID"
	}

	//we're adding a new post, so we use the topicID in order to create a new thread
	if req.FormValue("topicID") != "" {
		var topic forum.Topic
		if topic.ID, err = strconv.Atoi(req.FormValue("topicID")); err != nil {
			http.Error(rw, "Invalid topic ID", http.StatusInternalServerError)
			return "Invalid topic ID"
		}
		//verify that this topic exists
		if err = topic.Get(); err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return err.Error()
		}

		//create a new thread and save it
		thread := forum.Thread{
			TopicID: topic.ID,
			Active:  true,
		}
		if err = thread.Add(); err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return err.Error()
		}
		post.ThreadID = thread.ID
	} else {
		//we're replying to an existing post
		var parentPost forum.Post
		if parentPost.ID, err = strconv.Atoi(req.FormValue("parentID")); err != nil {
			http.Error(rw, "Invalid parent post ID", http.StatusInternalServerError)
			return "Invalid parent post ID"
		}
		if err = parentPost.Get(); err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return err.Error()
		}
		post.ParentID = parentPost.ID
		post.ThreadID = parentPost.ThreadID
	}

	if req.FormValue("notify") != "" {
		if post.Notify, err = strconv.ParseBool(req.FormValue("notify")); err != nil {
			http.Error(rw, "Invalid boolean for post.Notify", http.StatusInternalServerError)
			return "Invalid boolean for post.Notify"
		}
	}

	if req.FormValue("sticky") != "" {
		if post.Sticky, err = strconv.ParseBool(req.FormValue("sticky")); err != nil {
			http.Error(rw, "Invalid boolean for post.Sticky", http.StatusInternalServerError)
			return "Invalid boolean for post.Sticky"
		}
	}

	post.Approved = true
	post.Active = true
	post.Title = req.FormValue("title")
	post.Post = req.FormValue("post")
	post.Name = req.FormValue("name")
	post.Email = req.FormValue("email")
	post.Company = req.FormValue("company")
	//post.IPAddress = ""

	if err = post.Add(); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return err.Error()
	}

	return encoding.Must(enc.Encode(post))
}

func UpdatePost(rw http.ResponseWriter, req *http.Request, params martini.Params, enc encoding.Encoder) string {
	var err error
	var post forum.Post

	if post.ID, err = strconv.Atoi(params["id"]); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return err.Error()
	}

	if err = post.Get(); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return err.Error()
	}

	if req.FormValue("parentID") != "" {
		var parentPost forum.Post
		if parentPost.ID, err = strconv.Atoi(req.FormValue("parentID")); err != nil {
			http.Error(rw, "Invalid parent post ID", http.StatusInternalServerError)
			return "Invalid parent post ID"
		}
		if err = parentPost.Get(); err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return err.Error()
		}
		post.ParentID = parentPost.ID
	}

	if req.FormValue("threadID") != "" {
		var thread forum.Thread
		if thread.ID, err = strconv.Atoi(req.FormValue("threadID")); err != nil {
			http.Error(rw, "Invalid thread ID", http.StatusInternalServerError)
			return "Invalid thread ID"
		}
		if err = thread.Get(); err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return err.Error()
		}
		post.ThreadID = thread.ID
	}

	if req.FormValue("approved") != "" {
		if post.Approved, err = strconv.ParseBool(req.FormValue("approved")); err != nil {
			http.Error(rw, "Invalid boolean for post.Approved", http.StatusInternalServerError)
			return "Invalid boolean for post.Approved"
		}
	}

	if req.FormValue("active") != "" {
		if post.Approved, err = strconv.ParseBool(req.FormValue("active")); err != nil {
			http.Error(rw, "Invalid boolean for post.Active", http.StatusInternalServerError)
			return "Invalid boolean for post.Active"
		}
	}

	if req.FormValue("notify") != "" {
		if post.Notify, err = strconv.ParseBool(req.FormValue("notify")); err != nil {
			http.Error(rw, "Invalid boolean for post.Notify", http.StatusInternalServerError)
			return "Invalid boolean for post.Notify"
		}
	}

	if req.FormValue("sticky") != "" {
		if post.Sticky, err = strconv.ParseBool(req.FormValue("sticky")); err != nil {
			http.Error(rw, "Invalid boolean for post.Sticky", http.StatusInternalServerError)
			return "Invalid boolean for post.Sticky"
		}
	}

	if req.FormValue("flag") != "" {
		if post.Flag, err = strconv.ParseBool(req.FormValue("flag")); err != nil {
			http.Error(rw, "Invalid boolean for post.Flag", http.StatusInternalServerError)
			return "Invalid boolean for post.Flag"
		}
	}

	post.Title = req.FormValue("title")
	post.Post = req.FormValue("post")
	post.Name = req.FormValue("name")
	post.Email = req.FormValue("email")
	post.Company = req.FormValue("company")
	//post.IPAddress = ""

	if err = post.Update(); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return err.Error()
	}

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
