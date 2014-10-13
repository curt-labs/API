package cartIntegration

import (
	"database/sql"
	// "github.com/curt-labs/GoAPI/helpers/database"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestCI(t *testing.T) {
	Convey("Testing CartIntegration", t, func() {
		var ci CartIntegration
		var p PricePoint
		var err error
		ci.PartID = 11000
		ci.CustID = 1
		ci.CustPartID = 123456789

		err = ci.Create()
		So(err, ShouldBeNil)

		ci.CustPartID = 1234567890
		err = ci.Update()
		So(err, ShouldBeNil)

		err = ci.Get()
		So(err, ShouldBeNil)

		cis, err := GetCartIntegrationsByPart(ci)
		So(err, ShouldBeNil)
		So(len(cis), ShouldBeGreaterThan, 0)

		cis, err = GetCartIntegrationsByCustomer(ci)
		So(err, ShouldBeNil)
		So(len(cis), ShouldBeGreaterThan, 0)

		pricesList, err := GetPricesByCustomerID(ci.CustID)
		So(err, ShouldBeNil)
		So(len(pricesList), ShouldBeGreaterThan, 0)

		pagedPricesList, err := GetPricesByCustomerIDPaged(ci.CustID, 1, 1)
		So(err, ShouldBeNil)
		So(len(pagedPricesList), ShouldEqual, 1)

		count, err := GetPricingCount(ci.CustID)
		So(err, ShouldBeNil)
		So(count, ShouldBeGreaterThan, 0)

		p.CartIntegration = ci
		err = p.GetCustPriceID()
		So(err, ShouldBeNil)
		So(p.CartIntegration.CustPartID, ShouldEqual, ci.CustPartID)

		err = ci.Delete()
		So(err, ShouldBeNil)

		cis, err = GetAllCartIntegrations()
		if err != sql.ErrNoRows {
			So(err, ShouldBeNil)
			So(len(cis), ShouldBeGreaterThan, 0)
		}

	})
}
