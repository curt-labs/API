package cart

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestGet(t *testing.T) {
	Convey("Testing GetShop", t, func() {
		Convey("with no Id", func() {
			shop, err := GetShop("")
			So(shop, ShouldBeNil)
			So(err, ShouldNotBeNil)
		})
		Convey("create a test shop", func() {
			if id := insertTestData(); id != "" {
				Convey("with Id", func() {
					shop, err := GetShop("")
					So(shop, ShouldBeNil)
					So(err, ShouldNotBeNil)
				})
			}

		})
	})

}
