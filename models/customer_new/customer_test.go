package customer_new

import (
	"github.com/curt-labs/GoAPI/helpers/database"
	"github.com/curt-labs/GoAPI/models/customer"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

const (
	inputTimeFormat = "01/02/2006"
)

func TestCustomerModel(t *testing.T) {
	Convey("Testing Price - Gets", t, func() {
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
	Convey("Testing Price -  CUD", t, func() {
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
	Convey("Testing Customer", t, func() {
		Convey("Testing GetCustomer()", func() {
			var c Customer
			var err error
			c.Id = 108664501
			err = c.GetCustomer()
			So(err, ShouldBeNil)
			So(c.Name, ShouldNotBeNil)
		})
		Convey("Testing Basics()", func() {
			var c Customer
			var err error
			c.Id = 108664501
			err = c.Basics()
			So(err, ShouldBeNil)
			So(c.Name, ShouldNotBeNil)
		})
		Convey("Testing GetLocations()", func() {
			var c Customer
			var err error
			c.Id = 1 //choose customer with locations
			err = c.GetLocations()
			So(err, ShouldBeNil)
			So(len(c.Locations), ShouldBeGreaterThan, 0)
		})
		Convey("Testing GetUsers()", func() {
			var c Customer
			var err error
			c.Id = 10579901
			users, err := c.GetUsers()
			So(err, ShouldBeNil)
			So(users, ShouldNotBeNil)
		})
		Convey("Testing GetCustomerPrice()", func() {
			var c Customer
			var err error
			c.Id = 1
			partId := 11000
			api := "8AEE0620-412E-47FC-900A-947820EA1C1D"
			price, err := GetCustomerPrice(api, partId)
			So(err, ShouldBeNil)
			So(price, ShouldNotBeNil)
		})
		Convey("Testing GetCustomerCartReference())", func() {
			var c Customer
			var err error
			c.Id = 1
			partId := 11000
			api := "8AEE0620-412E-47FC-900A-947820EA1C1D"
			ref, err := GetCustomerCartReference(api, partId)
			So(err, ShouldBeNil)
			So(ref, ShouldNotBeNil)
		})
		Convey("Testing GetEtailers()", func() {
			var err error
			dealers, err := GetEtailers()
			So(err, ShouldBeNil)
			So(dealers, ShouldNotBeNil)
			So(len(dealers), ShouldBeGreaterThan, 0)
		})
		Convey("Testing GetLocalDealers()", func() {
			var err error
			latlng := "43.853282,-95.571675,45.800981,-90.468526&"
			center := "44.83536,-93.0201"
			dealers, err := GetLocalDealers(center, latlng)
			So(err, ShouldBeNil)
			So(dealers, ShouldNotBeNil)
			So(len(dealers), ShouldBeGreaterThan, 0)
		})
		Convey("Testing GetLocalRegions()", func() {
			var err error
			regions, err := GetLocalRegions()
			So(err, ShouldBeNil)
			So(regions, ShouldNotBeNil)
			So(len(regions), ShouldBeGreaterThan, 0)
		})
		Convey("Testing GetLocalDealerTiers()", func() {
			var err error
			tiers, err := GetLocalDealerTiers()
			So(err, ShouldBeNil)
			So(tiers, ShouldNotBeNil)
			So(len(tiers), ShouldBeGreaterThan, 0)
		})
		Convey("Testing GetLocalDealerTypes()", func() {
			var err error
			graphics, err := GetLocalDealerTypes()
			So(err, ShouldBeNil)
			So(graphics, ShouldNotBeNil)
			So(len(graphics), ShouldBeGreaterThan, 0)
		})
		Convey("Testing GetWhereToBuyDealers()", func() {
			var err error
			customers, err := GetWhereToBuyDealers()
			So(err, ShouldBeNil)
			So(customers, ShouldNotBeNil)
			So(len(customers), ShouldBeGreaterThan, 0)
		})
		Convey("Testing GetLocationById()", func() {
			var err error
			id := 1
			location, err := GetLocationById(id)
			So(err, ShouldBeNil)
			So(location, ShouldNotBeNil)
		})
		Convey("Testing SearchLocations()", func() {
			var err error
			term := "hitch"
			locations, err := SearchLocations(term)
			So(err, ShouldBeNil)
			So(locations, ShouldNotBeNil)
			So(len(locations), ShouldBeGreaterThan, 0)
		})
		Convey("Testing SearchLocationsByType()", func() {
			var err error
			term := "hitch"
			locations, err := SearchLocationsByType(term)
			So(err, ShouldBeNil)
			So(locations, ShouldNotBeNil)
			So(len(locations), ShouldBeGreaterThan, 0)
		})
		Convey("Testing SearchLocationsByLatLng()", func() {
			var err error
			latlng := GeoLocation{
				Latitude:  43.853282,
				Longitude: -95.571675,
			}
			locations, err := SearchLocationsByLatLng(latlng)
			So(err, ShouldBeNil)
			So(locations, ShouldNotBeNil)
			So(len(locations), ShouldBeGreaterThan, 0)
		})
	})
	Convey("Testing User", t, func() {
		Convey("Testing UserAuthentication()", func() {
			var u CustomerUser
			var err error
			u.Id = "F43F5B82-D7AF-4905-BD07-FC8BCF4C82FE"
			password := "test"
			c, err := u.UserAuthentication(password)
			So(err, ShouldNotBeNil)
			So(c, ShouldBeZeroValue)
		})
		Convey("Testing UserAuthenticationByKey()", func() {
			var err error
			key := "F43F5B82-D7AF-4905-BD07-FC8BCF4C82FE"
			c, err := UserAuthenticationByKey(key)
			So(err, ShouldNotBeNil)
			So(c, ShouldBeZeroValue)
		})
		Convey("Testing GetCustomer()", func() {
			var u CustomerUser
			var err error
			u.Id = "F43F5B82-D7AF-4905-BD07-FC8BCF4C82FE"
			c, err := u.GetCustomer()
			So(err, ShouldBeNil)
			So(c.Name, ShouldNotBeNil)
			So(c.Name, ShouldEqual, "Alex's Hitches") //TODO
			So(c.Id, ShouldEqual, 1)
		})
		Convey("Testing AuthenticateUser()", func() {
			var u CustomerUser
			var err error
			u.Id = "F43F5B82-D7AF-4905-BD07-FC8BCF4C82FE"
			password := "wrongPassword"
			err = u.AuthenticateUser(password)
			So(err, ShouldNotBeNil) //TODO - update user and auth

		})
		Convey("Testing ResetAuthentication()", func() {
			var u CustomerUser
			var err error
			u.Id = "F43F5B82-D7AF-4905-BD07-FC8BCF4C82FE"
			err = u.ResetAuthentication()
			So(err, ShouldBeNil)
		})
		Convey("Testing AuthenticateUserByKey()", func() {
			var err error
			key := "DE0A3046-380F-4816-AECC-1D239A0FF1D0" //auth type
			u, err := AuthenticateUserByKey(key)
			So(err, ShouldBeNil)
			So(u, ShouldNotBeNil)
		})
		Convey("GetKeys()", func() {
			var u CustomerUser
			var err error
			u.Id = "F43F5B82-D7AF-4905-BD07-FC8BCF4C82FE"
			err = u.GetKeys()
			So(u.Keys, ShouldNotBeNil)
			So(err, ShouldBeNil)
		})
		Convey("GetLocation()", func() {
			var u CustomerUser
			var err error
			u.Id = "023E68B9-9B62-4E5D-84F7-EDB88428B4F8"
			err = u.GetLocation()
			So(u.Location, ShouldNotBeNil)
			So(err, ShouldBeNil)
		})

		Convey("GetCustomerIdFromKey()", func() {
			var err error
			key := "BB337D2C-1613-4B4D-A2D5-D151CC96888C"
			id, err := GetCustomerIdFromKey(key)
			So(id, ShouldNotBeNil)
			So(err, ShouldBeNil)
		})
		Convey("GetCustomerUserFromKey()", func() {
			var err error
			key := "BB337D2C-1613-4B4D-A2D5-D151CC96888C"
			user, err := GetCustomerUserFromKey(key)
			So(user, ShouldNotBeNil)
			So(err, ShouldBeNil)
		})
		Convey("GetCustomerUserFromId()", func() {
			var err error
			id := "F43F5B82-D7AF-4905-BD07-FC8BCF4C82FE"
			user, err := GetCustomerUserById(id)
			So(user, ShouldNotBeNil)
			So(err, ShouldBeNil)
		})
	})

	Convey("Testing Existing User object to the New One", t, func() {
		err := database.PrepareAll()
		So(err, ShouldBeNil)
		Convey("Testing GetCustomer()", func() {
			var cc customer.CustomerUser
			var c CustomerUser
			var err error
			c.Id = "F43F5B82-D7AF-4905-BD07-FC8BCF4C82FE"
			cc.Id = "F43F5B82-D7AF-4905-BD07-FC8BCF4C82FE"
			user, err := c.GetCustomer()
			user2, err := cc.GetCustomer()
			So(err, ShouldBeNil)
			So(user.Name, ShouldNotBeNil)
			So(user2.Name, ShouldNotBeNil)
			So(user.Name, ShouldEqual, user2.Name)
			So(user.DealerType.Id, ShouldEqual, user2.DealerType.Id)
		})
		Convey("Testing Basics()", func() {
			var cc customer.Customer
			var c Customer
			var err error
			c.Id = 1
			cc.Id = 1
			err = c.Basics()
			err = cc.Basics()
			So(err, ShouldBeNil)
			So(c.Name, ShouldNotBeNil)
			So(cc.Name, ShouldNotBeNil)
			So(cc.Name, ShouldEqual, c.Name)
			So(cc.DealerType.Id, ShouldEqual, c.DealerType.Id)
		})
		Convey("Testing GetLocations()", func() {
			var cc customer.Customer
			var c Customer
			var err error
			c.Id = 1
			cc.Id = 1
			err = c.GetLocations()
			err = cc.GetLocations()
			So(err, ShouldBeNil)
			So(c.Name, ShouldNotBeNil)
			So(cc.Name, ShouldNotBeNil)
			So(cc.Name, ShouldEqual, c.Name)
			So(cc.DealerType.Id, ShouldEqual, c.DealerType.Id)
		})
		Convey("Testing GetCustomerPrice()", func() {
			var cc customer.Customer
			var c Customer
			var err error
			api := "8AEE0620-412E-47FC-900A-947820EA1C1D"
			partId := 11001
			c.Id = 1
			cc.Id = 1
			price, err := GetCustomerPrice(api, partId)
			price2, err := customer.GetCustomerPrice(api, partId)
			So(err, ShouldBeNil)
			So(price, ShouldNotBeNil)
			So(price2, ShouldNotBeNil)
			So(price, ShouldEqual, price2)
		})
		Convey("GetLocation()", func() {
			var u CustomerUser
			var u2 customer.CustomerUser
			var err error
			u.Id = "023E68B9-9B62-4E5D-84F7-EDB88428B4F8"
			u2.Id = "023E68B9-9B62-4E5D-84F7-EDB88428B4F8"
			err = u.GetLocation()
			So(u.Location, ShouldNotBeNil)
			So(err, ShouldBeNil)
			err = u2.GetLocation()
			So(u2.Location, ShouldNotBeNil)
			So(err, ShouldBeNil)
			So(u.Location.State.State, ShouldResemble, u2.Location.State.State)
		})

	})
}
