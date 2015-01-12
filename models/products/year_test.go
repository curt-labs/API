package products

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestGetYears(t *testing.T) {
	var l Lookup
	l.Brands = append(l.Brands, 1)
	Convey("Testing GetYears()", t, func() {
		err := l.GetYears()
		So(err, ShouldEqual, nil)
		So(l.Years, ShouldNotEqual, nil)
		So(l.Years, ShouldHaveSameTypeAs, []int{})
	})
}
