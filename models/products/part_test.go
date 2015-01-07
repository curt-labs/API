package products

import (
	"github.com/curt-labs/GoAPI/helpers/apicontext"
	"github.com/curt-labs/GoAPI/helpers/apicontextmock"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestPart(t *testing.T) {
	MockedDTX := &apicontext.DataContext{}
	var err error
	if MockedDTX, err = apicontextmock.Mock(); err != nil {
		return
	}
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
