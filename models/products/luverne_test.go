package products

import (
	"testing"

	mgo "gopkg.in/mgo.v2"

	"github.com/curt-labs/API/helpers/database"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGetLuverneYears(t *testing.T) {
	Convey("Test getLuverneYears(ctx *LuverneLookupContext)", t, func() {
		Convey("with no context", func() {
			years, err := getLuverneYears(nil)
			So(err, ShouldNotBeNil)
			So(years, ShouldHaveSameTypeAs, []string{})
		})

		Convey("with no mongo session", func() {
			years, err := getLuverneYears(&LuverneLookupContext{})
			So(err, ShouldNotBeNil)
			So(years, ShouldHaveSameTypeAs, []string{})
		})

		Convey("with a bad collection name", func() {
			ctx := &LuverneLookupContext{}
			var err error
			tmp := database.ProductCollectionName
			database.ProductCollectionName = ""

			ctx.Session, err = mgo.DialWithInfo(database.MongoPartConnectionString())
			So(err, ShouldBeNil)

			years, err := getLuverneYears(ctx)
			So(err, ShouldNotBeNil)
			So(years, ShouldHaveSameTypeAs, []string{})

			database.ProductCollectionName = tmp
		})

		Convey("with no status", func() {
			ctx := &LuverneLookupContext{}
			var err error

			ctx.Session, err = mgo.DialWithInfo(database.MongoPartConnectionString())
			So(err, ShouldBeNil)

			years, err := getLuverneYears(ctx)
			So(err, ShouldBeNil)
			So(years, ShouldHaveSameTypeAs, []string{})
			So(len(years), ShouldEqual, 0)
		})

		Convey("success", func() {
			ctx := &LuverneLookupContext{
				Statuses: []int{800, 900},
			}
			var err error

			ctx.Session, err = mgo.DialWithInfo(database.MongoPartConnectionString())
			So(err, ShouldBeNil)

			years, err := getLuverneYears(ctx)
			t.Log(years)
			So(err, ShouldBeNil)
			So(years, ShouldHaveSameTypeAs, []string{})
		})
	})
}

func TestGetLuverneMakes(t *testing.T) {
	Convey("Test getLuverneMakes(ctx *LuverneLookupContext)", t, func() {
		Convey("with no context", func() {
			vals, err := getLuverneMakes(nil, "")
			So(err, ShouldNotBeNil)
			So(vals, ShouldHaveSameTypeAs, []string{})
		})

		Convey("with no mongo session", func() {
			vals, err := getLuverneMakes(&LuverneLookupContext{}, "")
			So(err, ShouldNotBeNil)
			So(vals, ShouldHaveSameTypeAs, []string{})
		})

		Convey("with a bad collection name", func() {
			ctx := &LuverneLookupContext{}
			var err error
			tmp := database.ProductCollectionName
			database.ProductCollectionName = ""

			ctx.Session, err = mgo.DialWithInfo(database.MongoPartConnectionString())
			So(err, ShouldBeNil)

			vals, err := getLuverneMakes(ctx, "")
			So(err, ShouldNotBeNil)
			So(vals, ShouldHaveSameTypeAs, []string{})

			database.ProductCollectionName = tmp
		})

		Convey("with no status", func() {
			ctx := &LuverneLookupContext{}
			var err error

			ctx.Session, err = mgo.DialWithInfo(database.MongoPartConnectionString())
			So(err, ShouldBeNil)

			vals, err := getLuverneMakes(ctx, "")
			So(err, ShouldBeNil)
			So(vals, ShouldHaveSameTypeAs, []string{})
			So(len(vals), ShouldEqual, 0)
		})

		Convey("with no year", func() {
			ctx := &LuverneLookupContext{
				Statuses: []int{800, 900},
			}
			var err error

			ctx.Session, err = mgo.DialWithInfo(database.MongoPartConnectionString())
			So(err, ShouldBeNil)

			vals, err := getLuverneMakes(ctx, "")
			So(err, ShouldBeNil)
			So(vals, ShouldHaveSameTypeAs, []string{})
			So(len(vals), ShouldEqual, 0)
		})

		Convey("success", func() {
			ctx := &LuverneLookupContext{
				Statuses: []int{800, 900},
			}
			var err error

			ctx.Session, err = mgo.DialWithInfo(database.MongoPartConnectionString())
			So(err, ShouldBeNil)

			vals, err := getLuverneMakes(ctx, "2014")
			t.Log(vals)
			So(err, ShouldBeNil)
			So(vals, ShouldHaveSameTypeAs, []string{})
		})
	})
}

func TestGetLuverneModels(t *testing.T) {
	Convey("Test getLuverneModels(ctx *LuverneLookupContext)", t, func() {
		Convey("with no context", func() {
			vals, err := getLuverneModels(nil, "", "")
			So(err, ShouldNotBeNil)
			So(vals, ShouldHaveSameTypeAs, []string{})
		})

		Convey("with no mongo session", func() {
			vals, err := getLuverneModels(&LuverneLookupContext{}, "", "")
			So(err, ShouldNotBeNil)
			So(vals, ShouldHaveSameTypeAs, []string{})
		})

		Convey("with a bad collection name", func() {
			ctx := &LuverneLookupContext{}
			var err error
			tmp := database.ProductCollectionName
			database.ProductCollectionName = ""

			ctx.Session, err = mgo.DialWithInfo(database.MongoPartConnectionString())
			So(err, ShouldBeNil)

			vals, err := getLuverneModels(ctx, "", "")
			So(err, ShouldNotBeNil)
			So(vals, ShouldHaveSameTypeAs, []string{})

			database.ProductCollectionName = tmp
		})

		Convey("with no status", func() {
			ctx := &LuverneLookupContext{}
			var err error

			ctx.Session, err = mgo.DialWithInfo(database.MongoPartConnectionString())
			So(err, ShouldBeNil)

			vals, err := getLuverneModels(ctx, "", "")
			So(err, ShouldBeNil)
			So(vals, ShouldHaveSameTypeAs, []string{})
			So(len(vals), ShouldEqual, 0)
		})

		Convey("with no year", func() {
			ctx := &LuverneLookupContext{
				Statuses: []int{800, 900},
			}
			var err error

			ctx.Session, err = mgo.DialWithInfo(database.MongoPartConnectionString())
			So(err, ShouldBeNil)

			vals, err := getLuverneModels(ctx, "", "")
			So(err, ShouldBeNil)
			So(vals, ShouldHaveSameTypeAs, []string{})
			So(len(vals), ShouldEqual, 0)
		})

		Convey("with no model", func() {
			ctx := &LuverneLookupContext{
				Statuses: []int{800, 900},
			}
			var err error

			ctx.Session, err = mgo.DialWithInfo(database.MongoPartConnectionString())
			So(err, ShouldBeNil)

			vals, err := getLuverneModels(ctx, "", "")
			So(err, ShouldBeNil)
			So(vals, ShouldHaveSameTypeAs, []string{})
			So(len(vals), ShouldEqual, 0)
		})

		Convey("success", func() {
			ctx := &LuverneLookupContext{
				Statuses: []int{800, 900},
			}
			var err error

			ctx.Session, err = mgo.DialWithInfo(database.MongoPartConnectionString())
			So(err, ShouldBeNil)

			vals, err := getLuverneModels(ctx, "2014", "Ford")
			t.Log(vals)
			So(err, ShouldBeNil)
			So(vals, ShouldHaveSameTypeAs, []string{})
		})
	})
}
