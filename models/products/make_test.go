package products

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestGetMakes(t *testing.T) {
	var l Lookup
	l.Brands = append(l.Brands, 1)
	Convey("Testing GetMakes() without year", t, func() {
		err := l.GetMakes()
		So(err, ShouldEqual, nil)
		So(l.Makes, ShouldNotEqual, nil)
		So(len(l.Makes), ShouldEqual, 0)
	})

	Convey("Testing GetMakes() with bogus year", t, func() {
		l.Vehicle.Base.Year = 1
		err := l.GetMakes()
		So(err, ShouldEqual, nil)
		So(l.Makes, ShouldNotEqual, nil)
		So(len(l.Makes), ShouldEqual, 0)
		So(l.Vehicle.Base.Year, ShouldEqual, 1)
	})

	Convey("Testing GetMakes() with year", t, func() {
		l.Vehicle.Base.Year = 2010
		err := l.GetMakes()
		So(err, ShouldEqual, nil)
		So(l.Makes, ShouldNotEqual, nil)
		So(l.Makes, ShouldHaveSameTypeAs, []string{})
		So(l.Vehicle.Base.Year, ShouldEqual, 2010)
	})
}
