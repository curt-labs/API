package forum_ctlr

import (
	"github.com/curt-labs/GoAPI/helpers/apicontextmock"
	"github.com/curt-labs/GoAPI/helpers/httprunner"
	"github.com/curt-labs/GoAPI/helpers/testThatHttp"
	"github.com/curt-labs/GoAPI/models/forum"
	. "github.com/smartystreets/goconvey/convey"

	"encoding/json"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestForums(t *testing.T) {
	var err error
	var g forum.Group
	var to forum.Topic
	var th forum.Thread
	var p forum.Post
	var gs forum.Groups
	var tos forum.Topics
	var ths forum.Threads
	var ps forum.Posts

	dtx, err := apicontextmock.Mock()
	if err != nil {
		t.Log(err)
	}

	Convey("Testing Forums", t, func() {
		//test add group
		form := url.Values{"name": {"Posts About Ponies"}, "description": {"The wonderful world of ponies."}}
		v := form.Encode()
		body := strings.NewReader(v)
		thyme := time.Now()
		testThatHttp.RequestWithDtx("post", "/forum/groups", "", "?key="+dtx.APIKey, AddGroup, body, "application/x-www-form-urlencoded", dtx)
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &g)
		So(err, ShouldBeNil)
		So(g, ShouldHaveSameTypeAs, forum.Group{})
		So(g.ID, ShouldBeGreaterThan, 0)

		//test add topic
		form = url.Values{"name": {"The Prettiest Ponies"}, "description": {"We rank them by mane."}, "closed": {"false"}, "groupID": {strconv.Itoa(g.ID)}}
		v = form.Encode()
		body = strings.NewReader(v)
		thyme = time.Now()
		testThatHttp.RequestWithDtx("post", "/forum/topics", "", "?key="+dtx.APIKey, AddTopic, body, "application/x-www-form-urlencoded", dtx)
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &to)
		So(err, ShouldBeNil)
		So(to, ShouldHaveSameTypeAs, forum.Topic{})
		So(to.ID, ShouldBeGreaterThan, 0)

		//test add post
		form = url.Values{"title": {"Ponies"}, "post": {"I like pink and yellow ones the best."}, "name": {"Michael Jordan"}, "topicID": {strconv.Itoa(to.ID)}}
		v = form.Encode()
		body = strings.NewReader(v)
		thyme = time.Now()
		testThatHttp.RequestWithDtx("post", "/forum/posts", "", "?key="+dtx.APIKey, AddPost, body, "application/x-www-form-urlencoded", dtx)
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &p)
		So(err, ShouldBeNil)
		So(p, ShouldHaveSameTypeAs, forum.Post{})
		So(p.ID, ShouldBeGreaterThan, 0)

		//test update group
		form = url.Values{"description": {"Ponies are exciting"}}
		v = form.Encode()
		body = strings.NewReader(v)
		thyme = time.Now()
		testThatHttp.RequestWithDtx("put", "/forum/groups/", ":id", strconv.Itoa(g.ID)+"?key="+dtx.APIKey, UpdateGroup, body, "application/x-www-form-urlencoded", dtx)
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &g)
		So(err, ShouldBeNil)
		So(g, ShouldHaveSameTypeAs, forum.Group{})
		So(g.Description, ShouldNotEqual, "The wonderful world of ponies.")

		//test update topic
		form = url.Values{"description": {"We rank them by mane color."}, "closed": {"false"}, "groupID": {strconv.Itoa(g.ID)}, "active": {"true"}}
		v = form.Encode()
		body = strings.NewReader(v)
		thyme = time.Now()
		testThatHttp.RequestWithDtx("put", "/forum/topics/", ":id", strconv.Itoa(to.ID)+"?key="+dtx.APIKey, UpdateTopic, body, "application/x-www-form-urlencoded", dtx)
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &to)
		So(err, ShouldBeNil)
		So(to, ShouldHaveSameTypeAs, forum.Topic{})
		So(to.Description, ShouldNotEqual, "We rank them by mane.")

		//test update post
		form = url.Values{"title": {"Ponies"}, "post": {"I like pink and yellow ones the best."}, "name": {"Michael Jordan"}, "topicID": {strconv.Itoa(to.ID)}}
		v = form.Encode()
		body = strings.NewReader(v)
		thyme = time.Now()
		testThatHttp.RequestWithDtx("put", "/forum/posts/", ":id", strconv.Itoa(p.ID)+"?key="+dtx.APIKey, UpdatePost, body, "application/x-www-form-urlencoded", dtx)
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &p)
		So(err, ShouldBeNil)
		So(p, ShouldHaveSameTypeAs, forum.Post{})
		So(p.ID, ShouldBeGreaterThan, 0)

		//test get group
		thyme = time.Now()
		testThatHttp.RequestWithDtx("get", "/forum/groups/", ":id", strconv.Itoa(g.ID)+"?key="+dtx.APIKey, GetGroup, nil, "", dtx)
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &g)
		So(err, ShouldBeNil)
		So(g, ShouldHaveSameTypeAs, forum.Group{})

		//test get all groups
		thyme = time.Now()
		testThatHttp.RequestWithDtx("get", "/forum/groups", "", "?key="+dtx.APIKey, GetAllGroups, nil, "", dtx)
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &gs)
		So(err, ShouldBeNil)
		So(gs, ShouldHaveSameTypeAs, forum.Groups{})
		So(len(gs), ShouldBeGreaterThan, 0)

		//test get topic
		thyme = time.Now()
		testThatHttp.RequestWithDtx("get", "/forum/topics/", ":id", strconv.Itoa(to.ID)+"?key="+dtx.APIKey, GetTopic, nil, "", dtx)
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &to)
		So(err, ShouldBeNil)
		So(to, ShouldHaveSameTypeAs, forum.Topic{})

		//test get all topics
		thyme = time.Now()
		testThatHttp.RequestWithDtx("get", "/forum/topics", "", "?key="+dtx.APIKey, GetAllTopics, nil, "", dtx)
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &tos)
		So(err, ShouldBeNil)
		So(tos, ShouldHaveSameTypeAs, forum.Topics{})
		So(len(tos), ShouldBeGreaterThan, 0)

		//test get thread
		thyme = time.Now()
		testThatHttp.RequestWithDtx("get", "/forum/threads/", ":id", strconv.Itoa(p.ThreadID)+"?key="+dtx.APIKey, GetThread, nil, "", dtx)
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &th)
		So(err, ShouldBeNil)
		So(th, ShouldHaveSameTypeAs, forum.Thread{})

		//test get all threads
		thyme = time.Now()
		testThatHttp.RequestWithDtx("get", "/forum/threads", "", "?key="+dtx.APIKey, GetAllThreads, nil, "", dtx)
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &ths)
		So(err, ShouldBeNil)
		So(ths, ShouldHaveSameTypeAs, forum.Threads{})
		So(len(ths), ShouldBeGreaterThan, 0)

		//test get post
		thyme = time.Now()
		testThatHttp.RequestWithDtx("get", "/forum/posts/", ":id", strconv.Itoa(p.ID)+"?key="+dtx.APIKey, GetPost, nil, "", dtx)
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &p)
		So(err, ShouldBeNil)
		So(p, ShouldHaveSameTypeAs, forum.Post{})

		//test get all posts
		thyme = time.Now()
		testThatHttp.RequestWithDtx("get", "/forum/posts", "", "?key="+dtx.APIKey, GetAllPosts, nil, "", dtx)
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &ps)
		So(err, ShouldBeNil)
		So(ps, ShouldHaveSameTypeAs, forum.Posts{})
		So(len(ps), ShouldBeGreaterThan, 0)

		//test delete post
		thyme = time.Now()
		testThatHttp.RequestWithDtx("delete", "/forum/posts/", ":id", strconv.Itoa(p.ID)+"?key="+dtx.APIKey, DeletePost, nil, "", dtx)
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &p)
		So(err, ShouldBeNil)
		So(p, ShouldHaveSameTypeAs, forum.Post{})

		//test delete thread
		thyme = time.Now()
		testThatHttp.RequestWithDtx("delete", "/forum/threads/", ":id", strconv.Itoa(th.ID)+"?key="+dtx.APIKey, DeleteThread, nil, "", dtx)
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &th)
		So(err, ShouldBeNil)
		So(th, ShouldHaveSameTypeAs, forum.Thread{})

		//test delete topic
		thyme = time.Now()
		testThatHttp.RequestWithDtx("delete", "/forum/topics/", ":id", strconv.Itoa(to.ID)+"?key="+dtx.APIKey, DeleteTopic, body, "application/x-www-form-urlencoded", dtx)
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &to)
		So(err, ShouldBeNil)
		So(to, ShouldHaveSameTypeAs, forum.Topic{})

		//test delete group
		thyme = time.Now()
		testThatHttp.RequestWithDtx("delete", "/forum/groups/", ":id", strconv.Itoa(g.ID)+"?key="+dtx.APIKey, DeleteGroup, nil, "", dtx)
		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &g)
		So(err, ShouldBeNil)
		So(g, ShouldHaveSameTypeAs, forum.Group{})

	})
	_ = apicontextmock.DeMock(dtx)
}

func BenchmarkCRUDForum(b *testing.B) {
	qs := make(url.Values, 0)

	formGroup := url.Values{"name": {"Posts About Ponies"}, "description": {"The wonderful world of ponies."}}
	formTopic := url.Values{"name": {"The Prettiest Ponies"}, "description": {"We rank them by mane."}, "closed": {"false"}, "groupID": {"1"}}
	formPost := url.Values{"title": {"Ponies"}, "post": {"I like pink and yellow ones the best."}, "name": {"Michael Jordan"}, "topicID": {"1"}}
	Convey("Faqs", b, func() {
		//create group
		(&httprunner.BenchmarkOptions{
			Method:             "POST",
			Route:              "/forum/groups",
			ParameterizedRoute: "/forum/groups",
			Handler:            AddGroup,
			QueryString:        &qs,
			JsonBody:           nil,
			FormBody:           formGroup,
			Runs:               b.N,
		}).RequestBenchmark()
		//create group
		(&httprunner.BenchmarkOptions{
			Method:             "POST",
			Route:              "/forum/topics",
			ParameterizedRoute: "/forum/topics",
			Handler:            AddTopic,
			QueryString:        &qs,
			JsonBody:           nil,
			FormBody:           formTopic,
			Runs:               b.N,
		}).RequestBenchmark()
		//create group
		(&httprunner.BenchmarkOptions{
			Method:             "POST",
			Route:              "/forum/posts",
			ParameterizedRoute: "/forum/posts",
			Handler:            AddPost,
			QueryString:        &qs,
			JsonBody:           nil,
			FormBody:           formPost,
			Runs:               b.N,
		}).RequestBenchmark()

		//delete group
		(&httprunner.BenchmarkOptions{
			Method:             "DELETE",
			Route:              "/forum/groups",
			ParameterizedRoute: "/forum/groups/1",
			Handler:            DeleteGroup,
			QueryString:        &qs,
			JsonBody:           nil,
			FormBody:           nil,
			Runs:               b.N,
		}).RequestBenchmark()

		//delete group
		(&httprunner.BenchmarkOptions{
			Method:             "DELETE",
			Route:              "/forum/topics",
			ParameterizedRoute: "/forum/topics/1",
			Handler:            DeleteTopic,
			QueryString:        &qs,
			JsonBody:           nil,
			FormBody:           nil,
			Runs:               b.N,
		}).RequestBenchmark()

		//delete group
		(&httprunner.BenchmarkOptions{
			Method:             "DELETE",
			Route:              "/forum/posts",
			ParameterizedRoute: "/forum/posts/1",
			Handler:            DeletePost,
			QueryString:        &qs,
			JsonBody:           nil,
			FormBody:           nil,
			Runs:               b.N,
		}).RequestBenchmark()
	})
}
