package products

import (
	"github.com/curt-labs/API/helpers/apicontextmock"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestPart(t *testing.T) {
	MockedDTX, err := apicontextmock.Mock()
	if err != nil {
		return
	}
	Convey("Testing Basics", t, func() {
		p := Part{
			ID: 11000,
		}
		err := p.FromDatabase([]int{1, 3})
		So(err, ShouldBeNil)
	})

	Convey("Testing All", t, func() {
		parts, err := All(0, 1, MockedDTX)
		So(err, ShouldBeNil)
		So(len(parts), ShouldEqual, 1)
		So(parts, ShouldHaveSameTypeAs, []Part{})
	})

	Convey("Testing GetLatest", t, func() {
		parts, err := Latest(10, MockedDTX, 1)
		So(err, ShouldBeNil)
		So(len(parts), ShouldEqual, 10)
		So(parts, ShouldHaveSameTypeAs, []Part{})
		So(parts[0].DateAdded.String(), ShouldBeGreaterThan, parts[8].DateAdded.String())
	})

	Convey("Testing Featured", t, func() {
		parts, err := Featured(3, MockedDTX, 1)
		So(err, ShouldBeNil)
		So(len(parts), ShouldEqual, 3)
		So(parts, ShouldHaveSameTypeAs, []Part{})
	})

	Convey("Testing Related", t, func() {
		p := Part{
			ID: 11000,
		}
		p.Get(MockedDTX)
		parts, err := p.GetRelated(MockedDTX)
		So(err, ShouldBeNil)
		So(parts, ShouldHaveSameTypeAs, []Part{})
	})

	Convey("Testing BindCustomerToSeveralParts", t, func() {
		p := Part{
			ID: 11000,
		}
		p1 := Part{
			ID: 110003,
		}
		p2 := Part{
			ID: 110013,
		}
		_, err := BindCustomerToSeveralParts([]Part{p, p1, p2}, MockedDTX)
		So(err, ShouldBeNil)

	})
	_ = apicontextmock.DeMock(MockedDTX)
}
