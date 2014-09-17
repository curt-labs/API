package forum

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestTopics(t *testing.T) {
	var top Topic
	var lastTopicID int
	var err error

	Convey("Testing Gets", t, func() {
		Convey("Testing GetAll()", func() {
			topics, err := GetAllTopics()

			So(topics, ShouldNotBeNil)
			So(len(topics), ShouldBeGreaterThanOrEqualTo, 0)
			So(err, ShouldBeNil)
		})

		Convey("Testing Get()", func() {
			Convey("Topic with ID of 0", func() {
				err = top.Get()
				So(top.ID, ShouldEqual, 0)
				So(top.GroupID, ShouldEqual, 0)
				So(top.Name, ShouldEqual, "")
				So(top.Description, ShouldEqual, "")
				So(top.Image, ShouldEqual, "")
				So(len(top.Threads), ShouldEqual, 0)
				So(err, ShouldNotBeNil)
			})

			Convey("Topic with ID of 1", func() {
				top.ID = 1
				err = top.Get()

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
			top.GroupID = 1
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
			top.GroupID = 1
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
			top.GroupID = 1
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
			err = top.Delete()

			So(err, ShouldNotBeNil)
		})

		Convey("Last Updated Topic", func() {
			top = Topic{ID: lastTopicID}
			err = top.Delete()

			So(err, ShouldBeNil)
		})
	})
}