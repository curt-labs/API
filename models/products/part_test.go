package products

import (
	"github.com/curt-labs/GoAPI/helpers/apicontextmock"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestPart(t *testing.T) {
	MockedDTX, err := apicontextmock.Mock()
	if err != nil {
		return
	}

	Convey("Testing All", t, func() {
		parts, err := All(0, 1, MockedDTX)
		So(err, ShouldBeNil)
		So(parts, ShouldHaveSameTypeAs, []Part{})
	})
	Convey("Testing GetLatest", t, func() {
		parts, err := Latest(10, MockedDTX)
		So(err, ShouldBeNil)
		So(parts, ShouldHaveSameTypeAs, []Part{})
	})
	Convey("Testing Featured", t, func() {
		parts, err := Featured(3, MockedDTX)
		So(err, ShouldBeNil)
		So(parts, ShouldHaveSameTypeAs, []Part{})
	})
	_ = apicontextmock.DeMock(MockedDTX)
}
