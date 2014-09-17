package forum

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestThreads(t *testing.T) {
	var th Thread
	var lastThreadID int
	var err error

	Convey("Testing Gets", t, func() {
		Convey("Testing GetAll()", func() {
			threads, err := GetAllThreads()

			So(len(threads), ShouldBeGreaterThanOrEqualTo, 0)
			So(err, ShouldBeNil)
		})

		Convey("Testing Get()", func() {
			Convey("Thread with ID of 0", func() {
				err = th.Get()
				So(th.ID, ShouldEqual, 0)
				So(err, ShouldNotBeNil)
			})

			Convey("Thread with non-zero ID", func() {
				//create a bogus thread
				th = Thread{}
				th.TopicID = 1
				th.Add()
				id := th.ID

				//lets try getting
				th = Thread{ID: id}
				err = th.Get()

				//let's remove it...we're done with it
				th.Delete()

				So(th.ID, ShouldNotEqual, 0)
				So(err, ShouldBeNil)
			})
		})
	})

	Convey("Testing Add()", t, func() {
		Convey("Add Empty Thread", func() {
			th = Thread{}
			err = th.Add()

			So(th.ID, ShouldEqual, 0)
			So(err, ShouldNotBeNil)
		})
		Convey("Add Valid Thread", func() {
			th = Thread{}
			th.TopicID = 1

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
			th.TopicID = 1

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
}
