package products

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestLatest(t *testing.T) {
	Convey("Testing GetLatest", t, func() {
		parts, err := Latest(generateAPIkey(), 10)
		So(err, ShouldBeNil)
		So(parts, ShouldHaveSameTypeAs, []Part{})
	})
}
