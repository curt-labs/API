package customer

import (
	// "github.com/curt-labs/goacesapi/helpers/pagination"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

const (
	inputTimeFormat = "01/02/2006"
)

func TestCustomerPriceModel(t *testing.T) {
	Convey("Testing Gets", t, func() {
		Convey("Testing Get()", func() {
			var p Price
			p.ID = 50
			err := p.Get()
			So(p, ShouldNotBeNil)
			So(err, ShouldBeNil)
		})
		Convey("Testing GetAllPrices()", func() {
			ps, err := GetAllPrices()
			So(len(ps), ShouldBeGreaterThan, 200000)
			So(err, ShouldBeNil)
		})
		Convey("Testing GetPricesByCustomer()", func() {
			var c Customer
			c.Id = 10439386
			custPrices, err := c.GetPricesByCustomer()
			So(custPrices, ShouldNotBeNil)
			So(err, ShouldBeNil)
		})
		Convey("Testing GetPricesByPart()", func() {
			partID := 11000
			prices, err := GetPricesByPart(partID)
			So(len(prices), ShouldNotBeNil)
			So(err, ShouldBeNil)
		})
		Convey("Testing GetPricesBySaleRange", func() {
			var s time.Time
			var e time.Time
			c := Customer{Id: 10439386}
			var err error
			s, err = time.Parse(inputTimeFormat, "2006-01-02 15:04:05")
			e, err = time.Parse(inputTimeFormat, "2016-01-02 15:04:05")
			prices, err := c.GetPricesBySaleRange(s, e)
			So(err, ShouldBeNil)
			So(len(prices), ShouldBeGreaterThan, 0)
			So(prices, ShouldNotBeNil)

		})
	})
	Convey("Testing CUD", t, func() {
		Convey("Testing Create() Update() Delete() Price", func() {
			var p Price
			var err error
			p.CustID = 666
			p.SaleEnd, err = time.Parse(inputTimeFormat, "02/12/2006")
			p.IsSale = 1
			p.Price = 666.00
			err = p.Create()
			So(err, ShouldBeNil)
			p.SaleStart, err = time.Parse(inputTimeFormat, "01/02/2007")
			err = p.Update()
			So(err, ShouldBeNil)
			err = p.Get()
			So(err, ShouldBeNil)
			t, err := time.Parse(inputTimeFormat, "02/12/2006")
			So(p.SaleStart, ShouldHaveSameTypeAs, t)
			err = p.Delete()
			So(err, ShouldBeNil)
		})

	})
}
