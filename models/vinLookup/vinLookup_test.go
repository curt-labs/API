package vinLookup

import (
	"database/sql"
	"github.com/curt-labs/GoAPI/helpers/apicontext"
	"github.com/curt-labs/GoAPI/models/products"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

var (
	buickVin  = "1g4ha5em2au000001"
	taurusVin = "1fahp2fw5ag100583"
	caddyVin  = "1g6da5egxa0100211"
	bogusVin  = "123456789"
)

func TestVinLookup(t *testing.T) {
	Convey("Testing VinPartLookup", t, func() {
		vs, err := VinPartLookup(buickVin, &apicontext.DataContext{})
		if err != sql.ErrNoRows {
			So(err, ShouldBeNil)
			So(len(vs.Parts), ShouldBeGreaterThanOrEqualTo, 1)

		}
		if err != sql.ErrNoRows {
			vs, err = VinPartLookup(taurusVin, &apicontext.DataContext{})
			So(err, ShouldBeNil)
			So(len(vs.Parts), ShouldBeGreaterThanOrEqualTo, 1)

			//Make sure it's a Taurus - VINs should be constant
			So(vs.Vehicle.Base.Model, ShouldEqual, "Taurus")

			//We have 2010 Taurus Hitches
			So(len(vs.Parts), ShouldBeGreaterThanOrEqualTo, 1)
		}

	})
	Convey("Testing GetVehicleConfigs->GetParts", t, func() {
		v, err := GetVehicleConfigs(caddyVin)
		if err != sql.ErrNoRows {
			So(err, ShouldBeNil)
			So(len(v.Configurations), ShouldBeGreaterThanOrEqualTo, 1)

			//get random vehicleConfig
			// i := rand.Intn(len(v.Configurations))

			//get parts
			// parts, err := v.
			// if err != sql.ErrNoRows {
			// 	So(err, ShouldBeNil)
			// 	So(len(parts), ShouldBeGreaterThanOrEqualTo, 1)
			// }
		}

	})
	Convey("Testing Bad Vin", t, func() {
		vs, err := VinPartLookup(bogusVin, &apicontext.DataContext{})
		So(err, ShouldNotBeNil)
		So(vs, ShouldHaveSameTypeAs, products.Lookup{})

		vcs, err := GetVehicleConfigs(bogusVin)
		So(err, ShouldNotBeNil)
		So(vcs, ShouldHaveSameTypeAs, products.Lookup{})
	})
}

func BenchmarkVinPartLookup(b *testing.B) {
	for i := 0; i < b.N; i++ {
		VinPartLookup(buickVin, &apicontext.DataContext{})
	}
}

func BenchmarkGetVehicleConfigs(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetVehicleConfigs(buickVin)
	}
}
