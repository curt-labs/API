package customer_new

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestBusiness(t *testing.T) {
	Convey("Testing Gets", t, func() {
		Convey("Testing GetAllBusinessClasses()", func() {
			classes, err := GetAllBusinessClasses()
			So(len(classes), ShouldBeGreaterThanOrEqualTo, 0)
			So(err, ShouldBeNil)
		})
	})
}
