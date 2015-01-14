package products

import (
	"github.com/curt-labs/GoAPI/helpers/apicontext"
	"github.com/curt-labs/GoAPI/helpers/apicontextmock"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestGetYears(t *testing.T) {
	var l Lookup
	l.Brands = append(l.Brands, 1)
	var err error
	l.Brands = append(l.Brands, 1)
	MockedDTX := &apicontext.DataContext{}
	if MockedDTX, err = apicontextmock.Mock(); err != nil {
		return
	}
	Convey("Testing GetYears()", t, func() {
		err := l.GetYears(MockedDTX)
		So(err, ShouldEqual, nil)
		So(l.Years, ShouldNotEqual, nil)
		So(l.Years, ShouldHaveSameTypeAs, []int{})
	})
	_ = apicontextmock.DeMock(MockedDTX)
}
