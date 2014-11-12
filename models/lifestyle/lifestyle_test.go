package lifestyle

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestGetLifestyles(t *testing.T) {
	var l Lifestyle

	Convey("Testing CRUD", t, func() {

		l.Name = "testName"
		l.LongDesc = "Long description"
		err := l.Create()
		So(err, ShouldBeNil)
		err = l.Get()
		So(err, ShouldBeNil)
		So(l, ShouldNotBeNil)
		So(l.Name, ShouldEqual, "testName")
		So(l.LongDesc, ShouldEqual, "Long description")

		//Update
		l.Name = "newName"
		l.Image = "image"
		l.ShortDesc = "Desc"
		err = l.Update()
		So(err, ShouldBeNil)
		err = l.Get()

		So(err, ShouldBeNil)
		So(l, ShouldNotBeNil)
		So(l.Name, ShouldEqual, "newName")
		So(l.Image, ShouldEqual, "image")
		So(l.ShortDesc, ShouldEqual, "Desc")

		//Gets
		ls, err := GetAll()
		So(err, ShouldBeNil)
		So(ls, ShouldHaveSameTypeAs, Lifestyles{})

		err = l.Get()
		So(err, ShouldBeNil)
		So(l, ShouldNotBeNil)
		So(l.Name, ShouldHaveSameTypeAs, "str")
		So(l.ShortDesc, ShouldHaveSameTypeAs, "str")

		//Delete
		err = l.Delete()
		So(err, ShouldBeNil)

	})
}
