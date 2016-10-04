package products

import (
	"fmt"
	"testing"

	mgo "gopkg.in/mgo.v2"

	"github.com/curt-labs/API/helpers/database"
	"github.com/curt-labs/API/helpers/redis"
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

func TestLuverneQuery(t *testing.T) {
	Convey("LuverneQuery(ctx *LuverneLookupContext, args ...string)", t, func() {

		Convey("with no params", func() {
			ctx := &LuverneLookupContext{
				Statuses: []int{800, 900},
			}
			var err error

			ctx.Session, err = mgo.DialWithInfo(database.MongoPartConnectionString())
			So(err, ShouldBeNil)

			val, err := LuverneQuery(ctx)
			So(err, ShouldBeNil)
			So(val, ShouldHaveSameTypeAs, &LuverneCategoryVehicle{})
		})

		Convey("with one empty param", func() {
			ctx := &LuverneLookupContext{
				Statuses: []int{800, 900},
			}
			var err error

			ctx.Session, err = mgo.DialWithInfo(database.MongoPartConnectionString())
			So(err, ShouldBeNil)

			val, err := LuverneQuery(ctx, "")
			So(err, ShouldBeNil)
			So(val, ShouldHaveSameTypeAs, &LuverneCategoryVehicle{})
			err = CleanLookupRedis("")
			So(err, ShouldBeNil)
		})

		Convey("with two empty params", func() {
			ctx := &LuverneLookupContext{
				Statuses: []int{800, 900},
			}
			var err error

			ctx.Session, err = mgo.DialWithInfo(database.MongoPartConnectionString())
			So(err, ShouldBeNil)

			val, err := LuverneQuery(ctx, "", "")
			So(err, ShouldBeNil)
			So(val, ShouldHaveSameTypeAs, &LuverneCategoryVehicle{})
			err = CleanLookupRedis("", "")
			So(err, ShouldBeNil)
		})

		Convey("with three empty params", func() {
			ctx := &LuverneLookupContext{
				Statuses: []int{800, 900},
			}
			var err error

			ctx.Session, err = mgo.DialWithInfo(database.MongoPartConnectionString())
			So(err, ShouldBeNil)

			val, err := LuverneQuery(ctx, "", "", "")
			So(err, ShouldBeNil)
			So(val, ShouldHaveSameTypeAs, &LuverneCategoryVehicle{})
			err = CleanLookupRedis("", "", "")
			So(err, ShouldBeNil)
		})

		Convey("with empty vehicle", func() {
			ctx := &LuverneLookupContext{
				Statuses: []int{800, 900},
			}
			var err error

			ctx.Session, err = mgo.DialWithInfo(database.MongoPartConnectionString())
			So(err, ShouldBeNil)

			val, err := LuverneQuery(ctx, "", "", "", "")
			So(err, ShouldBeNil)
			So(val, ShouldHaveSameTypeAs, &LuverneCategoryVehicle{})
			err = CleanLookupRedis("", "", "", "")
			So(err, ShouldBeNil)
		})

		Convey("with no make, model or category", func() {
			ctx := &LuverneLookupContext{
				Statuses: []int{800, 900},
			}
			var err error
			var vehicleArgs = []string{"2016", "", "", ""}

			ctx.Session, err = mgo.DialWithInfo(database.MongoPartConnectionString())
			So(err, ShouldBeNil)

			val, err := LuverneQuery(ctx, vehicleArgs...)
			So(err, ShouldBeNil)
			So(val, ShouldHaveSameTypeAs, &LuverneCategoryVehicle{})
			err = CleanLookupRedis(vehicleArgs...)
			So(err, ShouldBeNil)
		})

		Convey("with no model or category", func() {
			ctx := &LuverneLookupContext{
				Statuses: []int{800, 900},
			}
			var err error
			var vehicleArgs = []string{"2016", "Ram", "", ""}
			ctx.Session, err = mgo.DialWithInfo(database.MongoPartConnectionString())
			So(err, ShouldBeNil)

			val, err := LuverneQuery(ctx, vehicleArgs...)
			So(err, ShouldBeNil)
			So(val, ShouldHaveSameTypeAs, &LuverneCategoryVehicle{})
			err = CleanLookupRedis(vehicleArgs...)
			So(err, ShouldBeNil)
		})

		Convey("with no category", func() {
			ctx := &LuverneLookupContext{
				Statuses: []int{800, 900},
			}
			var err error
			var vehicleArgs = []string{"2016", "Ram", "Ram 1500", ""}
			ctx.Session, err = mgo.DialWithInfo(database.MongoPartConnectionString())
			So(err, ShouldBeNil)

			val, err := LuverneQuery(ctx, vehicleArgs...)
			So(err, ShouldBeNil)
			So(val, ShouldHaveSameTypeAs, &LuverneCategoryVehicle{})
			err = CleanLookupRedis(vehicleArgs...)
			So(err, ShouldBeNil)
		})

		Convey("success", func() {
			ctx := &LuverneLookupContext{
				Statuses: []int{800, 900},
			}
			var vehicleArgs = []string{"2016", "Ram", "Ram 1500", "Aluminum Oval Bed Rails"}
			var err error

			ctx.Session, err = mgo.DialWithInfo(database.MongoPartConnectionString())
			So(err, ShouldBeNil)

			val, err := LuverneQuery(ctx, vehicleArgs...)
			So(err, ShouldBeNil)
			So(val, ShouldHaveSameTypeAs, &LuverneCategoryVehicle{})
			err = CleanLookupRedis(vehicleArgs...)
			So(err, ShouldBeNil)
		})

		Convey("success different vehicle, no category", func() {
			ctx := &LuverneLookupContext{
				Statuses: []int{800, 900},
			}
			var err error
			var vehicleArgs = []string{"2014", "Chevrolet", "Silverado 2500HD", ""}
			ctx.Session, err = mgo.DialWithInfo(database.MongoPartConnectionString())
			So(err, ShouldBeNil)

			val, err := LuverneQuery(ctx, vehicleArgs...)
			So(err, ShouldBeNil)
			So(val, ShouldHaveSameTypeAs, &LuverneCategoryVehicle{})
			err = CleanLookupRedis(vehicleArgs...)
			So(err, ShouldBeNil)
		})

		Convey("success vehicle with body types", func() {
			ctx := &LuverneLookupContext{
				Statuses: []int{800, 900},
			}
			var err error
			var vehicleArgs = []string{"2015", "Ram", "Ram 3500", ""}

			ctx.Session, err = mgo.DialWithInfo(database.MongoPartConnectionString())
			So(err, ShouldBeNil)

			val, err := LuverneQuery(ctx, vehicleArgs...)
			So(err, ShouldBeNil)
			So(val, ShouldHaveSameTypeAs, &LuverneCategoryVehicle{})
			err = CleanLookupRedis(vehicleArgs...)
			So(err, ShouldBeNil)
		})
	})
}

func TestGetLuverneStyles(t *testing.T) {
	Convey("Test getLuverneStyles(ctx *LuverneLookupContext, year, vehicleMake, model, category string)", t, func() {
		Convey("with no context", func() {
			vals1, vals2, err := getLuverneStyles(nil, "", "", "", "")
			So(err, ShouldNotBeNil)
			So(vals1, ShouldHaveSameTypeAs, []Part{})
			So(vals2, ShouldHaveSameTypeAs, []LuverneLookupCategory{})
		})

		Convey("with no mongo session", func() {
			vals1, vals2, err := getLuverneStyles(&LuverneLookupContext{}, "", "", "", "")
			So(err, ShouldNotBeNil)
			So(vals1, ShouldHaveSameTypeAs, []Part{})
			So(vals2, ShouldHaveSameTypeAs, []LuverneLookupCategory{})
		})

		Convey("with a bad collection name", func() {
			ctx := &LuverneLookupContext{}
			var err error
			tmp := database.ProductCollectionName
			database.ProductCollectionName = ""

			ctx.Session, err = mgo.DialWithInfo(database.MongoPartConnectionString())
			So(err, ShouldBeNil)

			vals1, vals2, err := getLuverneStyles(&LuverneLookupContext{}, "", "", "", "")
			So(err, ShouldNotBeNil)
			So(vals1, ShouldHaveSameTypeAs, []Part{})
			So(vals2, ShouldHaveSameTypeAs, []LuverneLookupCategory{})

			database.ProductCollectionName = tmp
		})

		Convey("with no status", func() {
			ctx := &LuverneLookupContext{}
			var err error

			ctx.Session, err = mgo.DialWithInfo(database.MongoPartConnectionString())
			So(err, ShouldBeNil)

			vals1, vals2, err := getLuverneStyles(&LuverneLookupContext{}, "", "", "", "")
			So(err, ShouldNotBeNil)
			So(vals1, ShouldHaveSameTypeAs, []Part{})
			So(vals2, ShouldHaveSameTypeAs, []LuverneLookupCategory{})
			So(len(vals1), ShouldEqual, 0)
			So(len(vals2), ShouldEqual, 0)
		})

		Convey("with no vehicle info", func() {
			ctx := &LuverneLookupContext{
				Statuses: []int{800, 900},
			}
			var err error

			ctx.Session, err = mgo.DialWithInfo(database.MongoPartConnectionString())
			So(err, ShouldBeNil)

			vals1, vals2, err := getLuverneStyles(ctx, "", "", "", "")
			So(err, ShouldBeNil)
			So(vals1, ShouldHaveSameTypeAs, []Part{})
			So(vals2, ShouldHaveSameTypeAs, []LuverneLookupCategory{})
			So(len(vals1), ShouldEqual, 0)
			So(len(vals2), ShouldEqual, 0)
		})

		Convey("with no make, model or category", func() {
			ctx := &LuverneLookupContext{
				Statuses: []int{800, 900},
			}
			var err error

			ctx.Session, err = mgo.DialWithInfo(database.MongoPartConnectionString())
			So(err, ShouldBeNil)

			vals1, vals2, err := getLuverneStyles(ctx, "2016", "", "", "")
			So(err, ShouldBeNil)
			So(vals1, ShouldHaveSameTypeAs, []Part{})
			So(vals2, ShouldHaveSameTypeAs, []LuverneLookupCategory{})
			So(len(vals1), ShouldEqual, 0)
			So(len(vals2), ShouldEqual, 0)
		})

		Convey("with no model or category", func() {
			ctx := &LuverneLookupContext{
				Statuses: []int{800, 900},
			}
			var err error

			ctx.Session, err = mgo.DialWithInfo(database.MongoPartConnectionString())
			So(err, ShouldBeNil)

			vals1, vals2, err := getLuverneStyles(ctx, "2016", "Ram", "", "")
			So(err, ShouldBeNil)
			So(vals1, ShouldHaveSameTypeAs, []Part{})
			So(vals2, ShouldHaveSameTypeAs, []LuverneLookupCategory{})
			So(len(vals1), ShouldEqual, 0)
			So(len(vals2), ShouldEqual, 0)
		})

		Convey("with no category", func() {
			ctx := &LuverneLookupContext{
				Statuses: []int{800, 900},
			}
			var err error

			ctx.Session, err = mgo.DialWithInfo(database.MongoPartConnectionString())
			So(err, ShouldBeNil)

			vals1, vals2, err := getLuverneStyles(ctx, "2016", "Ram", "Ram 1500", "")
			So(err, ShouldBeNil)
			So(vals1, ShouldHaveSameTypeAs, []Part{})
			So(vals2, ShouldHaveSameTypeAs, []LuverneLookupCategory{})
			So(len(vals1), ShouldNotEqual, 0)
			So(len(vals2), ShouldNotEqual, 0)
		})

		Convey("success", func() {
			ctx := &LuverneLookupContext{
				Statuses: []int{800, 900},
			}
			var err error

			ctx.Session, err = mgo.DialWithInfo(database.MongoPartConnectionString())
			So(err, ShouldBeNil)

			vals1, vals2, err := getLuverneStyles(ctx, "2016", "Ram", "Ram 1500", "Aluminum Oval Bed Rails")
			So(err, ShouldBeNil)
			So(vals1, ShouldHaveSameTypeAs, []Part{})
			So(vals2, ShouldHaveSameTypeAs, []LuverneLookupCategory{})
			So(len(vals1), ShouldNotEqual, 0)
			So(len(vals2), ShouldNotEqual, 0)
		})
	})
}

func CleanLookupRedis(args ...string) error {
	var redisKey string
	for i, arg := range args {
		if i == 0 {
			redisKey = fmt.Sprintf("luverne:%s", arg)
		} else {
			redisKey = fmt.Sprintf("%s:%s", redisKey, arg)
		}
	}
	// clear redis key
	return redis.Delete(redisKey)
}
