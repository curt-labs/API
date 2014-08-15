package aces

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestGetYears(t *testing.T) {
	var l Lookup
	Convey("Testing GetYears()", t, func() {
		err := l.GetYears()
		So(err, ShouldEqual, nil)
		So(l.Years, ShouldNotEqual, nil)
		So(len(l.Years), ShouldNotEqual, 0)
	})
}
