package lifestyle

import (
	// "database/sql"
	// "github.com/curt-labs/GoAPI/helpers/database"
	. "github.com/smartystreets/goconvey/convey"
	"math/rand"
	// "strconv"
	"testing"
)

func TestGetLifestyles(t *testing.T) {
	Convey("Testing Lifestyle", t, func() {
		Convey("Gets lifestyles and a random lifestyle", func() {
			ls, err := GetAll()
			So(err, ShouldBeNil)
			if len(ls) > 0 {
				x := rand.Intn(len(ls))
				l := ls[x]

				Convey("Testing Get()", func() {
					l.Get()
					So(l, ShouldNotBeNil)
					So(l.Name, ShouldHaveSameTypeAs, "str")
					So(l.ShortDesc, ShouldHaveSameTypeAs, "str")

				})
			}
		})

		Convey("Testing C_UD", func() {
			Convey("Testing Create()", func() {
				var l Lifestyle
				l.Name = "testName"
				l.LongDesc = "Long description"
				err := l.Create()
				So(err, ShouldBeNil)
				err = l.Get()
				So(err, ShouldBeNil)
				So(l, ShouldNotBeNil)
				So(l.Name, ShouldEqual, "testName")
				So(l.LongDesc, ShouldEqual, "Long description")

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

					Convey("Testing Delete()", func() {
						l.Get()
						err := l.Delete()
						So(err, ShouldBeNil)

					})
				})
			})
		})
		Convey("Testing Bad Get()", func() {
			var l Lifestyle
			getLifestyle = "Bad Query Stmt"
			err := l.Get()
			So(err, ShouldNotBeNil)
		})
	})
}
