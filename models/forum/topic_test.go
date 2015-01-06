package forum

import (
	"github.com/curt-labs/GoAPI/helpers/apicontext"
	"github.com/curt-labs/GoAPI/helpers/apicontextmock"

	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestTopics(t *testing.T) {
	var top Topic
	var lastTopicID int
	var err error

	MockedDTX := &apicontext.DataContext{}
	if MockedDTX, err = apicontextmock.Mock(); err != nil {
		return
	}

	//setup a test group
	g := Group{}
	g.Name = "Test Group"
	g.Description = "This is a test group."
	g.Add(MockedDTX)

	toppy := Topic{}
	toppy.GroupID = g.ID
	toppy.Name = "test-topic"
	toppy.Description = "This is a test topic"
	toppy.Add()

	Convey("Testing Gets", t, func() {
		Convey("Testing GetAll()", func() {
			topics, err := GetAllTopics(MockedDTX)

			So(len(topics), ShouldBeGreaterThanOrEqualTo, 0)
			So(err, ShouldBeNil)
		})

		Convey("Testing Get()", func() {
			Convey("Empty Topic", func() {
				top = Topic{}
				err = top.Get(MockedDTX)
				So(top.ID, ShouldEqual, 0)
				So(top.GroupID, ShouldEqual, 0)
				So(top.Name, ShouldEqual, "")
				So(top.Description, ShouldEqual, "")
				So(top.Image, ShouldEqual, "")
				So(len(top.Threads), ShouldEqual, 0)
				So(err, ShouldNotBeNil)
			})

			Convey("Topic with non-zero ID", func() {
				top = Topic{ID: toppy.ID}
				err = top.Get(MockedDTX)

				So(top.ID, ShouldNotEqual, 0)
				So(len(top.Threads), ShouldBeGreaterThanOrEqualTo, 0)
				So(err, ShouldBeNil)
			})
		})
	})

	Convey("Testing Add()", t, func() {
		Convey("Add Empty Topic", func() {
			top = Topic{}
			err = top.Add()

			So(top.ID, ShouldEqual, 0)
			So(err, ShouldNotBeNil)
		})
		Convey("Add Valid Topic", func() {
			top = Topic{}
			top.GroupID = g.ID
			top.Name = "test-topic"
			top.Description = "This is a test topic."

			err = top.Add()

			So(top.ID, ShouldBeGreaterThan, 0)
			So(top.GroupID, ShouldBeGreaterThan, 0)
			So(top.Name, ShouldNotEqual, "")
			So(top.Description, ShouldNotEqual, "")
			So(err, ShouldBeNil)

			lastTopicID = top.ID
		})
	})

	Convey("Testing Update()", t, func() {
		Convey("Empty Topic", func() {
			top = Topic{}
			err = top.Update()

			So(top.ID, ShouldEqual, 0)
			So(top.GroupID, ShouldEqual, 0)
			So(top.Name, ShouldEqual, "")
			So(top.Description, ShouldEqual, "")
			So(err, ShouldNotBeNil)
		})

		Convey("Missing Name", func() {
			top = Topic{ID: lastTopicID}
			top.GroupID = g.ID
			top.Description = "This is a updated test topic."
			top.Name = ""

			err = top.Update()

			So(top.ID, ShouldNotEqual, 0)
			So(top.GroupID, ShouldNotEqual, 0)
			So(top.Description, ShouldNotEqual, "")
			So(top.Name, ShouldEqual, "")
			So(err, ShouldNotBeNil)
		})

		Convey("Missing Group ID", func() {
			top = Topic{ID: lastTopicID}
			top.GroupID = 0
			top.Name = "test-topic"
			top.Description = "This is a updated test topic."

			err = top.Update()

			So(top.ID, ShouldNotEqual, 0)
			So(top.Name, ShouldNotEqual, "")
			So(top.Description, ShouldNotEqual, "")
			So(top.GroupID, ShouldEqual, 0)
			So(err, ShouldNotBeNil)
		})

		Convey("Last Added Topic", func() {
			top = Topic{ID: lastTopicID}
			top.GroupID = g.ID
			top.Name = "test-topic"
			top.Description = "This is a updated test topic"

			err = top.Update()

			So(top.ID, ShouldNotEqual, 0)
			So(top.GroupID, ShouldNotEqual, 0)
			So(top.Name, ShouldNotEqual, "")
			So(top.Description, ShouldNotEqual, "")
			So(err, ShouldBeNil)
		})
	})

	Convey("Testing Delete()", t, func() {
		Convey("Empty Topic", func() {
			top = Topic{}
			err = top.Delete(MockedDTX)

			So(err, ShouldNotBeNil)
		})

		Convey("Last Updated Topic", func() {
			top = Topic{ID: lastTopicID}
			err = top.Delete(MockedDTX)

			So(err, ShouldBeNil)
		})
	})

	g.Delete(MockedDTX)

	_ = apicontextmock.DeMock(MockedDTX)
}
