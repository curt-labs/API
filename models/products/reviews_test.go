package products

import (
	"github.com/curt-labs/GoAPI/helpers/apicontext"
	"github.com/curt-labs/GoAPI/helpers/apicontextmock"
	. "github.com/smartystreets/goconvey/convey"

	"testing"
)

func TestGetReviews(t *testing.T) {
	var p Part
	var l Review
	var err error
	MockedDTX := &apicontext.DataContext{}
	if MockedDTX, err = apicontextmock.Mock(); err != nil {
		return
	}

	Convey("Testing Reviews", t, func() {
		Convey("Testing C_UD", func() {
			//create part to review
			p.ID = 999999
			p.Status = 900
			p.ShortDesc = "TEST"
			p.PriceCode = 129
			p.Class.ID = 1
			p.Featured = false
			p.AcesPartTypeID = 1212

			p.Create()
		})
		Convey("Testing Create()", func() {
			l.PartID = 999999
			l.Name = "testName"
			l.ReviewText = "Long description"
			err := l.Create()
			So(err, ShouldBeNil)
			err = l.Get(MockedDTX)
			So(err, ShouldBeNil)
			So(l, ShouldNotBeNil)
			So(l.Name, ShouldEqual, "testName")
			So(l.ReviewText, ShouldEqual, "Long description")
		})

		Convey("Testing Update()", func() {
			l.Name = "newName"
			l.Email = "email"
			l.Subject = "Desc"
			err := l.Update()
			So(err, ShouldBeNil)
			err = l.Get(MockedDTX)

			So(err, ShouldBeNil)
			So(l, ShouldNotBeNil)
			So(l.Name, ShouldEqual, "newName")
			So(l.Email, ShouldEqual, "email")
			So(l.Subject, ShouldEqual, "Desc")
		})

		Convey("Gets reviews and a random review", func() {
			ls, err := GetAllReviews(MockedDTX)
			So(err, ShouldBeNil)
			So(len(ls), ShouldBeGreaterThanOrEqualTo, 0)

			err = l.Get(MockedDTX)
			So(err, ShouldBeNil)
			So(l, ShouldNotBeNil)
			So(l.Name, ShouldHaveSameTypeAs, "str")
			So(l.Subject, ShouldHaveSameTypeAs, "str")

		})

		Convey("Testing Delete()", func() {
			l.Get(MockedDTX)
			err := l.Delete()
			So(err, ShouldBeNil)
			//delete part
			p.Delete()
		})

	})
	Convey("Testing Bad Get()", t, func() {
		var l Review
		getReview = "Bad Query Stmt"
		err := l.Get(MockedDTX)
		So(err, ShouldNotBeNil)
	})
	Convey("Testing ActiveApprovedReviews", t, func() {
		var l Part //will be no rows
		err := l.GetActiveApprovedReviews()
		So(err, ShouldBeNil)
	})
	_ = apicontextmock.DeMock(MockedDTX)
}
