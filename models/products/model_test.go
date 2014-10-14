package products

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestGetModels(t *testing.T) {
	var l Lookup
	Convey("Testing GetModels() without year/make", t, func() {
		err := l.GetModels()
		So(err, ShouldEqual, nil)
		So(l.Models, ShouldNotEqual, nil)
		So(len(l.Models), ShouldEqual, 0)
	})

	Convey("Testing GetModels() with bogus data", t, func() {
		l.Vehicle.Base.Year = 1
		l.Vehicle.Base.Make = "KD"
		err := l.GetModels()
		So(err, ShouldEqual, nil)
		So(l.Models, ShouldNotEqual, nil)
		So(len(l.Models), ShouldEqual, 0)
	})

	Convey("Testing GetModels() with year", t, func() {
		l.Vehicle.Base.Year = 2010
		err := l.GetModels()
		So(err, ShouldEqual, nil)
		So(l.Models, ShouldNotEqual, nil)
		So(len(l.Models), ShouldEqual, 0)
	})

	Convey("Testing GetModels()", t, func() {
		l.Vehicle.Base.Year = 2010
		l.Vehicle.Base.Make = "Ford"
		err := l.GetModels()
		So(err, ShouldEqual, nil)
		So(l.Models, ShouldNotEqual, nil)
		So(l.Models, ShouldHaveSameTypeAs, []string{})
		So(l.Vehicle.Base.Year, ShouldEqual, 2010)
		So(l.Vehicle.Base.Make, ShouldEqual, "Ford")
	})
}
