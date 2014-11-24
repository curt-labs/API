package cart

import (
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/mgo.v2/bson"
	"os"
	"testing"
)

func TestGet(t *testing.T) {
	Convey("Testing GetShop", t, func() {
		Convey("Bad Connection", func() {
			os.Setenv("MONGO_URL", "0.0.0.1")
			var shop Shop
			err := shop.Get()
			So(err, ShouldNotBeNil)
			os.Setenv("MONGO_URL", "")
		})
		Convey("with no Id", func() {
			var shop Shop
			err := shop.Get()
			So(bson.IsObjectIdHex(shop.Id.String()), ShouldBeFalse)
			So(err, ShouldNotBeNil)
		})
		Convey("create a test shop", func() {
			if id := InsertTestData(); id != nil {
				Convey("with Id", func() {
					shop := Shop{
						Id: *id,
					}
					err := shop.Get()
					So(shop, ShouldNotBeNil)
					So(err, ShouldBeNil)
				})
			}
		})
	})

}
