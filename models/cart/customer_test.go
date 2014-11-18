package cart

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestSinceId(t *testing.T) {
	Convey("Testing Customer Gets", t, func() {
		Convey("Testing SinceId with no Id", func() {
			custs, err := SinceId("")
			So(err, ShouldBeNil)
			So(custs, ShouldHaveSameTypeAs, []Customer{})
		})

		Convey("Testing Get with no Id", func() {
			cust := Customer{}
			err := cust.Get()
			So(err, ShouldNotBeNil)
		})
	})
}
