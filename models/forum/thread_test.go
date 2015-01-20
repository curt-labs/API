package forum

import (
	"github.com/curt-labs/GoAPI/helpers/apicontextmock"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestThreads(t *testing.T) {
	var th Thread
	var lastThreadID int
	MockedDTX, err := apicontextmock.Mock()
	if err != nil {
		return
	}
	//setup a testing group, a topic, and an existing thread
	g := Group{}
	g.Name = "Test Group"
	g.Description = "This is a test group."
	g.Add(MockedDTX)

	top := Topic{}
	top.GroupID = g.ID
	top.Name = "test-topic"
	top.Description = "This is a updated test topic"
	top.Add()

	thready := Thread{}
	thready.TopicID = top.ID
	thready.Add()

	//run our tests
	Convey("Testing Gets", t, func() {
		Convey("Testing GetAll()", func() {
			threads, err := GetAllThreads(MockedDTX)

			So(len(threads), ShouldBeGreaterThanOrEqualTo, 0)
			So(err, ShouldBeNil)
		})

		Convey("Testing Get()", func() {
			Convey("Thread with ID of 0", func() {
				err = th.Get(MockedDTX)
				So(th.ID, ShouldEqual, 0)
				So(err, ShouldNotBeNil)
			})

			Convey("Thread with non-zero ID", func() {
				th = Thread{ID: thready.ID}
				err = th.Get(MockedDTX)

				So(th.ID, ShouldNotEqual, 0)
				So(err, ShouldBeNil)
			})
		})
	})

	Convey("Testing Add()", t, func() {
		Convey("Add Empty Thread", func() {
			th = Thread{}
			th.TopicID = top.ID
			err = th.Add()
			So(err, ShouldBeNil)
		})
		Convey("Add Valid Thread", func() {
			th = Thread{}
			th.TopicID = top.ID

			err = th.Add()

			So(th.ID, ShouldBeGreaterThan, 0)
			So(err, ShouldBeNil)

			lastThreadID = th.ID
		})
	})

	Convey("Testing Update()", t, func() {
		Convey("Empty Thread", func() {
			th = Thread{}
			err = th.Update()

			So(th.ID, ShouldEqual, 0)
			So(err, ShouldNotBeNil)
		})

		Convey("Last Added Thread", func() {
			th = Thread{ID: lastThreadID}
			th.TopicID = top.ID

			err = th.Update()

			So(th.ID, ShouldNotEqual, 0)
			So(th.TopicID, ShouldNotEqual, 0)
			So(err, ShouldBeNil)
		})
	})

	Convey("Testing Delete()", t, func() {
		Convey("Empty Thread", func() {
			th = Thread{}
			err = th.Delete()

			So(err, ShouldNotBeNil)
		})

		Convey("Last Updated Thread", func() {
			th = Thread{ID: lastThreadID}
			err = th.Delete()

			So(err, ShouldBeNil)
		})
	})

	//destroy the group and everything tied to it
	g.Delete(MockedDTX)

	_ = apicontextmock.DeMock(MockedDTX)
}
