package forum

import (
	"github.com/curt-labs/GoAPI/helpers/apicontextmock"
	. "github.com/smartystreets/goconvey/convey"

	"testing"
)

func TestPosts(t *testing.T) {
	var p Post
	var lastPostID int
	MockedDTX, err := apicontextmock.Mock()
	if err != nil {
		return
	}

	//setup a group, topic, and thread
	g := Group{}
	g.Name = "Test Group"
	g.Description = "This is a test group."
	g.Add(MockedDTX)

	top := Topic{}
	top.GroupID = g.ID
	top.Name = "test-topic"
	top.Description = "This is a updated test topic"
	top.Add()

	th := Thread{}
	th.TopicID = top.ID
	th.Add()

	//run our tests
	Convey("Testing Gets", t, func() {
		Convey("Testing GetAll()", func() {
			posts, err := GetAllPosts(MockedDTX)

			So(len(posts), ShouldBeGreaterThanOrEqualTo, 0)
			So(err, ShouldBeNil)
		})

		Convey("Testing Get()", func() {
			Convey("Post with ID of 0", func() {
				err = p.Get(MockedDTX)
				So(p.ID, ShouldEqual, 0)
				So(err, ShouldNotBeNil)
			})
		})
	})

	Convey("Testing Add()", t, func() {
		Convey("Empty Post", func() {
			p = Post{}
			err = p.Add()

			So(p.ID, ShouldEqual, 0)
			So(err, ShouldNotBeNil)
		})
		Convey("Valid Post", func() {
			p = Post{}
			p.ThreadID = th.ID
			p.Title = "Test Post"
			p.Post = "This is a test post."

			err = p.Add()

			So(p.ID, ShouldNotEqual, 0)
			So(err, ShouldBeNil)

			lastPostID = p.ID
		})
	})

	Convey("Testing Update()", t, func() {
		Convey("Empty Post", func() {
			p = Post{}
			err = p.Update()

			So(p.ID, ShouldEqual, 0)
			So(err, ShouldNotBeNil)
		})
		Convey("Missing Thread ID", func() {
			p = Post{ID: lastPostID}
			p.ThreadID = 0
			p.Title = "Test Post"
			p.Post = "This is a updated test post."

			err = p.Update()

			So(err, ShouldNotBeNil)
		})
		Convey("Missing Title", func() {
			p = Post{ID: lastPostID}
			p.ThreadID = th.ID
			p.Title = ""
			p.Post = "This is a updated test post."

			err = p.Update()

			So(err, ShouldNotBeNil)
		})
		Convey("Missing Post Message", func() {
			p = Post{ID: lastPostID}
			p.ThreadID = th.ID
			p.Title = "Test Post"
			p.Post = ""

			err = p.Update()

			So(err, ShouldNotBeNil)
		})
		Convey("Notify - But Missing Email", func() {
			p = Post{ID: lastPostID}
			p.ThreadID = th.ID
			p.Title = "Test Post"
			p.Post = "This is a updated test post."
			p.Email = ""
			p.Notify = true

			err = p.Update()

			So(err, ShouldNotBeNil)
		})
		Convey("Last Added Post", func() {
			p = Post{ID: lastPostID}
			p.ThreadID = th.ID
			p.Title = "Test Post"
			p.Post = "This is a updated test post."

			err = p.Update()
			So(p.ID, ShouldNotEqual, 0)
			So(err, ShouldBeNil)
		})
	})

	Convey("Testing Delete()", t, func() {
		Convey("Empty Post", func() {
			p = Post{}
			err = p.Delete()

			So(err, ShouldNotBeNil)
		})
		Convey("Last Updated Post", func() {
			p = Post{ID: lastPostID}
			err = p.Delete()

			So(err, ShouldBeNil)
		})
	})

	//destroy the group and everything tied to it
	g.Delete(MockedDTX)

	_ = apicontextmock.DeMock(MockedDTX)
}
