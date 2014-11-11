package lifestyle

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestGetLifestyles(t *testing.T) {
	var l Lifestyle

	Convey("Testing CRUD", t, func() {
		Convey("Testing Create()", func() {

			l.Name = "testName"
			l.LongDesc = "Long description"
			err := l.Create()
			So(err, ShouldBeNil)
			err = l.Get()
			So(err, ShouldBeNil)
			So(l, ShouldNotBeNil)
			So(l.Name, ShouldEqual, "testName")
			So(l.LongDesc, ShouldEqual, "Long description")
		})

		Convey("Testing Update()", func() {
			l.Name = "newName"
			l.Image = "image"
			l.ShortDesc = "Desc"
			err := l.Update()
			So(err, ShouldBeNil)
			err = l.Get()
			t.Log(l)
			So(err, ShouldBeNil)
			So(l, ShouldNotBeNil)
			So(l.Name, ShouldEqual, "newName")
			So(l.Image, ShouldEqual, "image")
			So(l.ShortDesc, ShouldEqual, "Desc")
		})

		Convey("Test Gets", func() {
			ls, err := GetAll()
			So(err, ShouldBeNil)
			So(len(ls), ShouldBeGreaterThan, 0)

			err = l.Get()
			So(err, ShouldBeNil)
			So(l, ShouldNotBeNil)
			So(l.Name, ShouldHaveSameTypeAs, "str")
			So(l.ShortDesc, ShouldHaveSameTypeAs, "str")

		})

		Convey("Testing Delete()", func() {
			l.Get()
			err := l.Delete()
			So(err, ShouldBeNil)

		})

	})
}
