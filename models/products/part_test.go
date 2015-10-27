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
		parts, err := Latest(10, MockedDTX)
		So(err, ShouldBeNil)
		So(len(parts), ShouldEqual, 10)
		So(parts, ShouldHaveSameTypeAs, []Part{})
		So(parts[0].DateAdded.String(), ShouldBeGreaterThan, parts[8].DateAdded.String())
	})

	Convey("Testing Featured", t, func() {
		parts, err := Featured(3, MockedDTX)
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

	Convey("Get BY Old Part Number", t, func() {
		p := Part{
			OldPartNumber: "BM01821501",
		}
		err = p.GetPartByOldPartNumber()
		So(err, ShouldBeNil)
	})
	_ = apicontextmock.DeMock(MockedDTX)
}
