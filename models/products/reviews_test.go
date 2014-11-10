package products

import (
	// "database/sql"
	// "github.com/curt-labs/GoAPI/helpers/database"
	. "github.com/smartystreets/goconvey/convey"
	"math/rand"
	// "strconv"
	"testing"
)

func TestGetReviews(t *testing.T) {
	Convey("Testing REviews", t, func() {
		Convey("Gets reviews and a random review", func() {
			ls, err := GetAll()
			So(err, ShouldBeNil)
			if len(ls) > 0 {
				x := rand.Intn(len(ls))
				l := ls[x]

				Convey("Testing Get()", func() {
					l.Get()
					So(l, ShouldNotBeNil)
					So(l.Name, ShouldHaveSameTypeAs, "str")
					So(l.Subject, ShouldHaveSameTypeAs, "str")

				})
			}
		})

		Convey("Testing C_UD", func() {
			p := Part{
				ID:        999999,
				Status:    900,
				ShortDesc: "TEST",
				PriceCode: 129,
				Class: Class{
					ID: 1,
				},
				Featured:       false,
				AcesPartTypeID: 1212,
			}
			p.Create()

			Convey("Testing Create()", func() {
				var l Review
				l.PartID = 999999
				l.Name = "testName"
				l.ReviewText = "Long description"
				err := l.Create()
				So(err, ShouldBeNil)
				err = l.Get()
				So(err, ShouldBeNil)
				So(l, ShouldNotBeNil)
				So(l.Name, ShouldEqual, "testName")
				So(l.ReviewText, ShouldEqual, "Long description")

				Convey("Testing Update()", func() {
					l.Name = "newName"
					l.Email = "email"
					l.Subject = "Desc"
					err := l.Update()
					So(err, ShouldBeNil)
					err = l.Get()
					t.Log(l)
					So(err, ShouldBeNil)
					So(l, ShouldNotBeNil)
					So(l.Name, ShouldEqual, "newName")
					So(l.Email, ShouldEqual, "email")
					So(l.Subject, ShouldEqual, "Desc")

					Convey("Testing Delete()", func() {
						l.Get()
						err := l.Delete()
						So(err, ShouldBeNil)

					})
				})
			})

			p.Delete()
		})
		Convey("Testing Bad Get()", func() {
			var l Review
			getReview = "Bad Query Stmt"
			err := l.Get()
			So(err, ShouldNotBeNil)
		})
		Convey("Testing ActiveApprovedReviews", func() {
			var l Part //will be no rows
			err := l.GetActiveApprovedReviews()
			So(err, ShouldBeNil)
		})
	})
}
