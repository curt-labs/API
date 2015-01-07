package customer

import (
	//"database/sql"
	//"github.com/curt-labs/GoAPI/models/products"
	_ "github.com/go-sql-driver/mysql"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

const (
	inputTimeFormat = "01/02/2006"
)

// func BenchmarkCustomerGet(b *testing.B) {
// 	Convey("testing get", b, func() {
// 		var c Customer
// 		c.Id = 1
// 		b.ResetTimer()
// 		for i := 0; i < b.N; i++ {
// 			_ = c.Get()
// 		}

// 	})
// }
// func BenchmarkCustomerBasics(b *testing.B) {
// 	Convey("testing basics ", b, func() {
// 		var c Customer
// 		c.Id = 1

// 		b.ResetTimer()
// 		for i := 0; i < b.N; i++ {
// 			_ = c.Basics()
// 		}

// 	})
// }
func TestCustomerModel(t *testing.T) {
	Convey("Testing Customer Model", t, func() {
		var c Customer
		var err error

		//Create
		c.Name = "testCustomer"
		c.Address = "Nowhere"
		c.CustomerId = 666
		c.State.Id = 1 //TODO
		err = c.Create()
		So(err, ShouldBeNil)

		//create location
		var cl CustomerLocation
		cl.Name = "testLocation"
		cl.CustomerId = c.Id
		err = cl.Create()
		So(err, ShouldBeNil)

		//get Location
		err = cl.Get()
		So(err, ShouldBeNil)

		c.Locations = append(c.Locations, cl)

		//create User
		var cu CustomerUser
		cu.Name = "testUser"
		cu.Password = "test"
		cu.OldCustomerID = c.Id
		cu.Active = true
		cu.Location.Id = cl.Id
		cu.Sudo = false
		cu.CustomerID = c.CustomerId
		cu.Current = false

		// err = cu.Create()
		// So(err, ShouldBeNil)

		// cu = *someuser
		c.Users = append(c.Users, cu)

		//Upate
		c.Name = "New Name"
		c.MapixCode.ID = 1
		err = c.Update()
		So(err, ShouldBeNil)

		err = c.GetLocations()
		So(err, ShouldBeNil)
		So(len(c.Locations), ShouldBeGreaterThan, 0)

		//Gets
		err = c.GetCustomer("") //kills c.Locations
		So(err, ShouldBeNil)

		err = c.Basics("")
		So(err, ShouldBeNil)

		err = c.Get() //New
		So(err, ShouldBeNil)

		err = c.GetLocations()
		So(err, ShouldBeNil)
		So(len(c.Locations), ShouldBeGreaterThan, 0)

		err = c.FindCustomerIdFromCustId()
		So(err, ShouldBeNil)

		err = c.GetUsers("")
		So(err, ShouldBeNil)
		So(c.Users, ShouldHaveSameTypeAs, []CustomerUser{})

		//Create Part
		//var part products.Part
		//var custPrice products.Price
		//custPrice.Price = 123
		//part.Pricing = append(part.Pricing, custPrice)
		//err = part.Create()

		if len(cu.Keys) > 0 {
			//price, err := GetCustomerPrice(cu.Keys[0].Key, part.ID)
			//if err != sql.ErrNoRows {
			//	So(err, ShouldBeNil)
			//	So(price, ShouldEqual, 123)
			//}

			//ref, err := GetCustomerCartReference(cu.Keys[0].Key, part.ID)
			//if err != sql.ErrNoRows {
			//	So(err, ShouldBeNil)
			//	So(ref, ShouldNotBeNil)
			//}
		}

		//Delete
		err = c.Delete()
		So(err, ShouldBeNil)

		err = cl.Delete()
		So(err, ShouldBeNil)

		err = cu.Delete()
		So(err, ShouldBeNil)

	})

	Convey("testing general gets", t, func() {
		Convey("Testing GetEtailers()", func() {
			var err error
			dealers, err := GetEtailers("")
			So(err, ShouldBeNil)
			So(dealers, ShouldHaveSameTypeAs, []Customer{})
		})
		Convey("Testing GetLocalDealers()", func() {
			var err error
			latlng := "43.853282,-95.571675,45.800981,-90.468526&"
			center := "44.83536,-93.0201"
			dealers, err := GetLocalDealers(center, latlng)
			So(err, ShouldBeNil)
			So(dealers, ShouldHaveSameTypeAs, []DealerLocation{})
		})
		Convey("Testing GetLocalRegions()", func() {
			var err error
			regions, err := GetLocalRegions()
			So(err, ShouldBeNil)
			So(regions, ShouldHaveSameTypeAs, []StateRegion{})
		})
		Convey("Testing GetLocalDealerTiers()", func() {
			var err error
			tiers, err := GetLocalDealerTiers()
			So(err, ShouldBeNil)
			So(tiers, ShouldHaveSameTypeAs, []DealerTier{})
		})
		Convey("Testing GetLocalDealerTypes()", func() {
			var err error
			graphics, err := GetLocalDealerTypes()
			So(err, ShouldBeNil)
			So(graphics, ShouldHaveSameTypeAs, []MapGraphics{})
		})
		Convey("Testing GetWhereToBuyDealers()", func() {
			var err error
			customers, err := GetWhereToBuyDealers("")
			So(err, ShouldBeNil)
			So(customers, ShouldHaveSameTypeAs, []Customer{})
		})
		Convey("Testing SearchLocations()", func() {
			var err error
			term := "hitch"
			locations, err := SearchLocations(term)
			So(err, ShouldBeNil)
			So(locations, ShouldHaveSameTypeAs, []DealerLocation{})
		})
		Convey("Testing SearchLocationsByType()", func() {
			var err error
			term := "hitch"
			locations, err := SearchLocationsByType(term)
			So(err, ShouldBeNil)
			So(locations, ShouldHaveSameTypeAs, DealerLocations{})
		})
		Convey("Testing SearchLocationsByLatLng()", func() {
			var err error
			latlng := GeoLocation{
				Latitude:  43.853282,
				Longitude: -95.571675,
			}
			locations, err := SearchLocationsByLatLng(latlng)
			So(err, ShouldBeNil)
			So(locations, ShouldHaveSameTypeAs, []DealerLocation{})
		})
	})

}
