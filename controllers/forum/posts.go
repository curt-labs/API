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

func GetAllPosts(rw http.ResponseWriter, req *http.Request, enc encoding.Encoder, dtx *apicontext.DataContext) string {
	posts, err := forum.GetAllPosts(dtx)
	if err != nil {
		apierror.GenerateError("Trouble getting all forum posts", err, rw, req)
	}
	return encoding.Must(enc.Encode(posts))
}

func GetPost(rw http.ResponseWriter, req *http.Request, params martini.Params, enc encoding.Encoder, dtx *apicontext.DataContext) string {
	var err error
	var post forum.Post

	if post.ID, err = strconv.Atoi(params["id"]); err != nil {
		apierror.GenerateError("Trouble getting forum post ID", err, rw, req)
	}

	if err := post.Get(dtx); err != nil {
		apierror.GenerateError("Trouble getting forum post", err, rw, req)
	}
	return encoding.Must(enc.Encode(post))
}

func AddPost(rw http.ResponseWriter, req *http.Request, params martini.Params, enc encoding.Encoder, dtx *apicontext.DataContext) string {
	var err error
	var post forum.Post

	//we add a post either by topic (new post) or by identifying the parentID (replying to an existing post)
	if req.FormValue("topicID") == "" && req.FormValue("parentID") == "" {
		apierror.GenerateError("Trouble adding forum post -- Missing topic ID or parent post ID", err, rw, req)
	}

	//we're adding a new post, so we use the topicID in order to create a new thread
	if req.FormValue("topicID") != "" {
		var topic forum.Topic
		if topic.ID, err = strconv.Atoi(req.FormValue("topicID")); err != nil {
			apierror.GenerateError("Trouble adding forum post -- invalid topic ID", err, rw, req)
		}
		//verify that this topic exists
		if err = topic.Get(dtx); err != nil {
			apierror.GenerateError("Trouble getting forum post -- topic doesn't exist", err, rw, req)
		}

		//create a new thread and save it
		thread := forum.Thread{
			TopicID: topic.ID,
			Active:  true,
		}
		if err = thread.Add(); err != nil {
			apierror.GenerateError("Trouble adding forum post -- failed to add thread", err, rw, req)
		}
		post.ThreadID = thread.ID
	} else {
		//we're replying to an existing post
		var parentPost forum.Post
		if parentPost.ID, err = strconv.Atoi(req.FormValue("parentID")); err != nil {
			apierror.GenerateError("Trouble adding forum post reply -- parentID is invalid", err, rw, req)
		}
		if err = parentPost.Get(dtx); err != nil {
			apierror.GenerateError("Trouble getting parent forum post", err, rw, req)
		}
		post.ParentID = parentPost.ID
		post.ThreadID = parentPost.ThreadID
	}

	if req.FormValue("notify") != "" {
		if post.Notify, err = strconv.ParseBool(req.FormValue("notify")); err != nil {
			apierror.GenerateError("Trouble adding post -- boolean notify parameter is invalid", err, rw, req)
		}
	}

	if req.FormValue("sticky") != "" {
		if post.Sticky, err = strconv.ParseBool(req.FormValue("sticky")); err != nil {
			apierror.GenerateError("Trouble adding post -- boolean sticky parameter is invalid", err, rw, req)
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
		apierror.GenerateError("Trouble adding post", err, rw, req)
	}

	return encoding.Must(enc.Encode(post))
}

func UpdatePost(rw http.ResponseWriter, req *http.Request, params martini.Params, enc encoding.Encoder, dtx *apicontext.DataContext) string {
	var err error
	var post forum.Post

	if post.ID, err = strconv.Atoi(params["id"]); err != nil {
		apierror.GenerateError("Trouble getting post ID", err, rw, req)
	}

	if err = post.Get(dtx); err != nil {
		apierror.GenerateError("Trouble getting post", err, rw, req)
	}

	if req.FormValue("parentID") != "" {
		var parentPost forum.Post
		if parentPost.ID, err = strconv.Atoi(req.FormValue("parentID")); err != nil {
			apierror.GenerateError("Trouble updating post -- parent post ID is invalid", err, rw, req)
		}
		if err = parentPost.Get(dtx); err != nil {
			apierror.GenerateError("Trouble updating post -- failed to get parent post", err, rw, req)
		}
		post.ParentID = parentPost.ID
	}

	if req.FormValue("threadID") != "" {
		var thread forum.Thread
		if thread.ID, err = strconv.Atoi(req.FormValue("threadID")); err != nil {
			apierror.GenerateError("Trouble updating post -- thread ID is invalid", err, rw, req)
		}
		if err = thread.Get(dtx); err != nil {
			apierror.GenerateError("Trouble updating post -- thread was not found", err, rw, req)
		}
		post.ThreadID = thread.ID
	}

	if req.FormValue("approved") != "" {
		if post.Approved, err = strconv.ParseBool(req.FormValue("approved")); err != nil {
			apierror.GenerateError("Trouble updating post -- boolean approved parameter is invalid", err, rw, req)
		}
	}

	if req.FormValue("active") != "" {
		if post.Approved, err = strconv.ParseBool(req.FormValue("active")); err != nil {
			apierror.GenerateError("Trouble updating post -- boolean active parameter is invalid", err, rw, req)
		}
	}

	if req.FormValue("notify") != "" {
		if post.Notify, err = strconv.ParseBool(req.FormValue("notify")); err != nil {
			apierror.GenerateError("Trouble updating post -- boolean notify parameter is invalid", err, rw, req)
		}
	}

	if req.FormValue("sticky") != "" {
		if post.Sticky, err = strconv.ParseBool(req.FormValue("sticky")); err != nil {
			apierror.GenerateError("Trouble updating post -- boolean sticky parameter is invalid", err, rw, req)
		}
	}

	if req.FormValue("flag") != "" {
		if post.Flag, err = strconv.ParseBool(req.FormValue("flag")); err != nil {
			apierror.GenerateError("Trouble updating post -- boolean flag parameter is invalid", err, rw, req)
		}
	}

	if req.FormValue("title") != "" {
		post.Title = req.FormValue("title")
	}

	if req.FormValue("post") != "" {
		post.Post = req.FormValue("post")
	}

	if req.FormValue("name") != "" {
		post.Name = req.FormValue("name")
	}

	if req.FormValue("email") != "" {
		post.Email = req.FormValue("email")
	}

	if req.FormValue("company") != "" {
		post.Company = req.FormValue("company")
	}
	//post.IPAddress = ""

	if err = post.Update(); err != nil {
		apierror.GenerateError("Trouble updating post", err, rw, req)
	}

	return encoding.Must(enc.Encode(post))
}

func DeletePost(rw http.ResponseWriter, req *http.Request, params martini.Params, enc encoding.Encoder) string {
	var err error
	var post forum.Post

	if post.ID, err = strconv.Atoi(params["id"]); err != nil {
		apierror.GenerateError("Trouble getting forum post ID post", err, rw, req)
	}

	if err = post.Delete(); err != nil {
		apierror.GenerateError("Trouble deleting forum post", err, rw, req)
	}

	return encoding.Must(enc.Encode(post))
}
