package products

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestFindAppsLuverne(t *testing.T) {
	Convey("Test FindAppsLuverne(catID, skip, limit int)", t, func() {
		Convey("with no cat", func() {
			res, err := FindAppsLuverne(0, 0, 0)
			So(err, ShouldBeNil)
			So(res, ShouldHaveSameTypeAs, LuverneResult{})
		})

		Convey("with valid cat and limit", func() {
			res, err := FindAppsLuverne(364, 0, 50)
			So(err, ShouldBeNil)
			So(res, ShouldHaveSameTypeAs, LuverneResult{})
		})

		Convey("with valid cat, limit and second page", func() {
			res, err := FindAppsLuverne(364, 2, 50)
			So(err, ShouldBeNil)
			So(res, ShouldHaveSameTypeAs, LuverneResult{})
		})
	})
}
