package customer_new

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
	//From the NEW Customer Model
	Convey("Testing Customer_New", t, func() {
		Convey("Testing GetCustomer_New()", func() {
			var c Customer
			var err error
			c.Id = 108664501
			err = c.GetCustomer_New()
			So(err, ShouldBeNil)
			So(c.Name, ShouldNotBeNil)
		})
		Convey("Testing Basics_New()", func() {
			var c Customer
			var err error
			c.Id = 108664501
			err = c.Basics_New()
			So(err, ShouldBeNil)
			So(c.Name, ShouldNotBeNil)
		})
		Convey("Testing GetLocations_New()", func() {
			var c Customer
			var err error
			c.Id = 1 //choose customer with locations
			err = c.GetLocations_New()
			So(err, ShouldBeNil)
			So(len(c.Locations), ShouldBeGreaterThan, 0)
		})
		Convey("Testing GetUsers_New()", func() {
			var c Customer
			var err error
			c.Id = 10579901
			users, err := c.GetUsers_New()
			So(err, ShouldBeNil)
			So(users, ShouldNotBeNil)
		})
		Convey("Testing GetCustomerPrice_New()", func() {
			var c Customer
			var err error
			c.Id = 1
			partId := 11000
			api := "8AEE0620-412E-47FC-900A-947820EA1C1D"
			price, err := GetCustomerPrice_New(api, partId)
			So(err, ShouldBeNil)
			So(price, ShouldNotBeNil)
		})
		Convey("Testing GetCustomerCartReference_New()", func() {
			var c Customer
			var err error
			c.Id = 1
			partId := 11000
			api := "8AEE0620-412E-47FC-900A-947820EA1C1D"
			ref, err := GetCustomerCartReference_New(api, partId)
			So(err, ShouldBeNil)
			So(ref, ShouldNotBeNil)
		})
		Convey("Testing GetEtailers_New()", func() {
			var err error
			dealers, err := GetEtailers_New()
			So(err, ShouldBeNil)
			So(dealers, ShouldNotBeNil)
			So(len(dealers), ShouldBeGreaterThan, 0)
		})
		Convey("Testing GetLocalDealers_New()", func() {
			var err error
			latlng := "43.853282,-95.571675,45.800981,-90.468526&"
			center := "44.83536,-93.0201"
			dealers, err := GetLocalDealers_New(center, latlng)
			So(err, ShouldBeNil)
			So(dealers, ShouldNotBeNil)
			So(len(dealers), ShouldBeGreaterThan, 0)
		})
		Convey("Testing GetLocalRegions_New()", func() {
			var err error
			regions, err := GetLocalRegions_New()
			So(err, ShouldBeNil)
			So(regions, ShouldNotBeNil)
			So(len(regions), ShouldBeGreaterThan, 0)
		})
		Convey("Testing GetLocalDealerTiers_New()", func() {
			var err error
			tiers, err := GetLocalDealerTiers_New()
			So(err, ShouldBeNil)
			So(tiers, ShouldNotBeNil)
			So(len(tiers), ShouldBeGreaterThan, 0)
		})
		Convey("Testing GetLocalDealerTypes_New()", func() {
			var err error
			graphics, err := GetLocalDealerTypes_New()
			So(err, ShouldBeNil)
			So(graphics, ShouldNotBeNil)
			So(len(graphics), ShouldBeGreaterThan, 0)
		})
		Convey("Testing GetWhereToBuyDealers_New()", func() {
			var err error
			customers, err := GetWhereToBuyDealers_New()
			So(err, ShouldBeNil)
			So(customers, ShouldNotBeNil)
			So(len(customers), ShouldBeGreaterThan, 0)
		})
		Convey("Testing GetLocationById_New()", func() {
			var err error
			id := 1
			location, err := GetLocationById_New(id)
			So(err, ShouldBeNil)
			So(location, ShouldNotBeNil)
		})
		Convey("Testing SearchLocationsByType_New()", func() {
			var err error
			term := "hitch"
			locations, err := SearchLocationsByType_New(term)
			So(err, ShouldBeNil)
			So(locations, ShouldNotBeNil)
			So(len(locations), ShouldBeGreaterThan, 0)
		})
		Convey("Testing SearchLocationsByLatLng_New()", func() {
			var err error
			latlng := GeoLocation{
				Latitude:  43.853282,
				Longitude: -95.571675,
			}
			locations, err := SearchLocationsByLatLng_New(latlng)
			So(err, ShouldBeNil)
			So(locations, ShouldNotBeNil)
			So(len(locations), ShouldBeGreaterThan, 0)
		})
		//key=8AEE0620-412E-47FC-900A-947820EA1C1D
		//&latitude=43.853282&longitude=-95.571675
	})
	Convey("Testing User", t, func() {
		Convey("Testing GetCustomer()", func() {
			var u CustomerUser
			var err error
			u.Id = "F43F5B82-D7AF-4905-BD07-FC8BCF4C82FE"
			c, err := u.GetCustomer()
			So(err, ShouldBeNil)
			So(c.Name, ShouldNotBeNil)
			// So(c.Name, ShouldEqual, "Alex Ninneman")
			// So(c.Id, ShouldEqual, 1)
		})
	})

}
