package vinLookup

import (
	"github.com/curt-labs/GoAPI/helpers/apicontextmock"
	. "github.com/smartystreets/goconvey/convey"

	"database/sql"
	"testing"
)

var (
	buickVin  = "1g4ha5em2au000001"
	taurusVin = "1fahp2fw5ag100583"
	caddyVin  = "1g6da5egxa0100211"
	bogusVin  = "123456789"
)

func TestVinLookup(t *testing.T) {
	dtx, err := apicontextmock.Mock()
	if err != nil {
		t.Log(err)
	}

	Convey("Testing VinPartLookup", t, func() {
		vs, err := VinPartLookup(buickVin, dtx)
		if err != sql.ErrNoRows {
			So(err, ShouldBeNil)
			So(len(vs.Parts), ShouldBeGreaterThanOrEqualTo, 1)

		}
		if err != sql.ErrNoRows {
			vs, err = VinPartLookup(taurusVin, dtx)
			So(err, ShouldBeNil)
			So(len(vs.Parts), ShouldBeGreaterThanOrEqualTo, 1)

			//Make sure it's a Taurus - VINs should be constant
			So(vs.Vehicle.Base.Model, ShouldEqual, "Taurus")

			//We have 2010 Taurus Hitches
			So(len(vs.Parts), ShouldBeGreaterThanOrEqualTo, 1)
		}

	})
	_ = apicontextmock.DeMock(dtx)
}
