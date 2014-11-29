package cartIntegration

import (
	"database/sql"
	"github.com/curt-labs/GoAPI/models/customer"
	"github.com/curt-labs/GoAPI/models/products"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestCI(t *testing.T) {
	//setup
	var part products.Part
	var price customer.Price
	var cust customer.Customer
	var err error
	cust.CustomerId = 666
	cust.Create()

	part.ShortDesc = "test"
	part.ID = 123456789
	part.Status = 800
	err = part.Create()
	if err != nil {
		err = nil
		err = part.Update()
	}
	price.CustID = cust.CustomerId
	price.PartID = part.ID
	price.Create()

	Convey("Testing CartIntegration", t, func() {
		var ci CartIntegration
		var p PricePoint
		ci.PartID = part.ID
		ci.CustID = cust.Id
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
		So(pricesList, ShouldHaveSameTypeAs, []PricePoint{})

		pagedPricesList, err := GetPricesByCustomerIDPaged(ci.CustID, 1, 1)
		So(err, ShouldBeNil)
		So(pagedPricesList, ShouldHaveSameTypeAs, []PricePoint{})

		count, err := GetPricingCount(ci.CustID)
		So(err, ShouldBeNil)
		So(count, ShouldBeGreaterThanOrEqualTo, 0)

		p.CartIntegration = ci
		err = p.GetCustPriceID()
		So(err, ShouldBeNil)
		So(p.CartIntegration.CustPartID, ShouldEqual, ci.CustPartID)

		err = ci.Delete()
		So(err, ShouldBeNil)

		cis, err = GetAllCartIntegrations()
		if err != sql.ErrNoRows {
			So(err, ShouldBeNil)
			So(cis, ShouldHaveSameTypeAs, []CartIntegration{})
		}
	})
	//cleanup
	cust.Delete()
	part.Delete()
	price.Delete()
}

func BenchmarkGetAllCartIntegrations(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetAllCartIntegrations()
	}
}

func BenchmarkGetCartIntegration(b *testing.B) {
	ci := setupDummyCartIntegration()
	ci.Create()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ci.Get()
	}
	b.StopTimer()
	ci.Delete()
}

func BenchmarkGetCartIntegrationByPart(b *testing.B) {
	ci := setupDummyCartIntegration()
	ci.Create()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetCartIntegrationsByPart(*ci)
	}
	b.StopTimer()
	ci.Delete()
}

func BenchmarkGetCartIngegrationByCustomer(b *testing.B) {
	ci := setupDummyCartIntegration()
	ci.Create()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetCartIntegrationsByCustomer(*ci)
	}
	b.StopTimer()
	ci.Delete()
}

func BenchmarkGetPricesByCustomerID(b *testing.B) {
	ci := setupDummyCartIntegration()
	ci.Create()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetPricesByCustomerID(ci.CustID)
	}
	b.StopTimer()
	ci.Delete()
}

func BenchmarkGetPricesByCustomerIDPaged(b *testing.B) {
	ci := setupDummyCartIntegration()
	ci.Create()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetPricesByCustomerIDPaged(ci.CustID, 1, 1)
	}
	b.StopTimer()
	ci.Delete()
}

func BenchmarkGetPricingCount(b *testing.B) {
	ci := setupDummyCartIntegration()
	ci.Create()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetPricingCount(ci.CustID)
	}
	b.StopTimer()
	ci.Delete()
}

func BenchmarkCreateCartIntegration(b *testing.B) {
	ci := setupDummyCartIntegration()
	for i := 0; i < b.N; i++ {
		ci.Create()
		b.StopTimer()
		ci.Delete()
		b.StartTimer()
	}
}

func BenchmarkUpdateCartIntegration(b *testing.B) {
	ci := setupDummyCartIntegration()
	ci.Create()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ci.PartID = 987654
		ci.Update()
	}
	b.StopTimer()
	ci.Delete()
}

func BenchmarkDeleteCartIntegration(b *testing.B) {
	ci := setupDummyCartIntegration()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		ci.Create()
		b.StartTimer()
		ci.Delete()
	}
}

func setupDummyCartIntegration() *CartIntegration {
	return &CartIntegration{
		PartID:     999999,
		CustID:     999999,
		CustPartID: 123456789,
	}
}
