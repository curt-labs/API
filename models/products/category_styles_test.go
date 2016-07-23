package products

import (
	"testing"

	"gopkg.in/mgo.v2"

	"github.com/curt-labs/API/helpers/database"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGetYears(t *testing.T) {
	Convey("Test getYears(ctx *LookupContext)", t, func() {
		Convey("with no context", func() {
			years, err := getYears(nil)
			So(err, ShouldNotBeNil)
			So(years, ShouldHaveSameTypeAs, []string{})
		})

		Convey("with no mongo session", func() {
			years, err := getYears(&LookupContext{})
			So(err, ShouldNotBeNil)
			So(years, ShouldHaveSameTypeAs, []string{})
		})

		Convey("with a bad collection name", func() {
			ctx := &LookupContext{}
			var err error
			tmp := database.ProductCollectionName
			database.ProductCollectionName = ""

			ctx.Session, err = mgo.DialWithInfo(database.MongoPartConnectionString())
			So(err, ShouldBeNil)

			years, err := getYears(ctx)
			So(err, ShouldNotBeNil)
			So(years, ShouldHaveSameTypeAs, []string{})

			database.ProductCollectionName = tmp
		})

		Convey("with no brand or status", func() {
			ctx := &LookupContext{}
			var err error

			ctx.Session, err = mgo.DialWithInfo(database.MongoPartConnectionString())
			So(err, ShouldBeNil)

			years, err := getYears(ctx)
			So(err, ShouldBeNil)
			So(years, ShouldHaveSameTypeAs, []string{})
			So(len(years), ShouldEqual, 0)
		})

		Convey("with no brand", func() {
			ctx := &LookupContext{
				Statuses: []int{800, 900},
			}
			var err error

			ctx.Session, err = mgo.DialWithInfo(database.MongoPartConnectionString())
			So(err, ShouldBeNil)

			years, err := getYears(ctx)
			So(err, ShouldBeNil)
			So(years, ShouldHaveSameTypeAs, []string{})
			So(len(years), ShouldEqual, 0)
		})

		Convey("success", func() {
			ctx := &LookupContext{
				Statuses: []int{800, 900},
				Brands:   []int{1},
			}
			var err error

			ctx.Session, err = mgo.DialWithInfo(database.MongoPartConnectionString())
			So(err, ShouldBeNil)

			years, err := getYears(ctx)
			So(err, ShouldBeNil)
			So(years, ShouldHaveSameTypeAs, []string{})
		})
	})
}

func TestGetStyles(t *testing.T) {
	Convey("Test getStyles(ctx *LookupContext)", t, func() {
		Convey("success", func() {
			ctx := &LookupContext{
				Statuses: []int{800, 900},
				Brands:   []int{3},
			}
			var err error

			ctx.Session, err = mgo.DialWithInfo(database.MongoPartConnectionString())
			So(err, ShouldBeNil)

			_, _, err = getStyles(ctx, "2015", "jeep", "wrangler", "")
			So(err, ShouldBeNil)
		})
	})
}
