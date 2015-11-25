package customer

import (
	"github.com/curt-labs/API/helpers/apicontext"
	"github.com/curt-labs/API/helpers/apicontextmock"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestBusiness(t *testing.T) {
	var err error
	MockedDTX := &apicontext.DataContext{}
	if MockedDTX, err = apicontextmock.Mock(); err != nil {
		return
	}
	Convey("Testing Gets", t, func() {
		Convey("Testing GetAllBusinessClasses()", func() {
			classes, err := GetAllBusinessClasses(MockedDTX)
			So(len(classes), ShouldBeGreaterThanOrEqualTo, 0)
			So(err, ShouldBeNil)
		})
	})
	_ = apicontextmock.DeMock(MockedDTX)
}
