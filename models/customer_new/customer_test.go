package customer_new

import (
	"database/sql"
	// "github.com/curt-labs/GoAPI/helpers/database"
	"github.com/curt-labs/GoAPI/models/apiKeyType"
	"github.com/curt-labs/GoAPI/models/products"
	_ "github.com/go-sql-driver/mysql"
	. "github.com/smartystreets/goconvey/convey"
	// "strings"
	"testing"
)

const (
	inputTimeFormat = "01/02/2006"
)

func BenchmarkCustomerGet(b *testing.B) {
	Convey("testing get", b, func() {
		var c Customer
		c.Id = 1
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = c.Get()
		}

	})
}
func BenchmarkCustomerBasics(b *testing.B) {
	Convey("testing basics ", b, func() {
		var c Customer
		c.Id = 1

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = c.Basics()
		}

	})
}
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

		//API KEY types
		var pub apiKeyType.ApiKeyType
		var pri apiKeyType.ApiKeyType
		var aut apiKeyType.ApiKeyType

		pub.Type = "Public"
		pri.Type = "Private"
		aut.Type = "Authentication"
		err = pub.Create()
		err = pri.Create()
		err = aut.Create()
		So(err, ShouldBeNil)
		err = cu.Create()
		So(err, ShouldBeNil)

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
		err = c.GetCustomer() //kills c.Locations
		So(err, ShouldBeNil)

		err = c.Basics()
		So(err, ShouldBeNil)

		err = c.Get() //New
		So(err, ShouldBeNil)

		err = c.GetLocations()
		So(err, ShouldBeNil)
		So(len(c.Locations), ShouldBeGreaterThan, 0)

		err = c.FindCustomerIdFromCustId()
		So(err, ShouldBeNil)

		users, err := c.GetUsers()
		So(err, ShouldBeNil)
		So(users, ShouldHaveSameTypeAs, []CustomerUser{})

		//Create Part
		var part products.Part
		var custPrice products.Price
		custPrice.Price = 123
		part.Pricing = append(part.Pricing, custPrice)
		err = part.Create()

		if len(cu.Keys) > 0 {
			price, err := GetCustomerPrice(cu.Keys[0].Key, part.ID)
			if err != sql.ErrNoRows {
				So(err, ShouldBeNil)
				So(price, ShouldEqual, 123)
			}

			ref, err := GetCustomerCartReference(cu.Keys[0].Key, part.ID)
			if err != sql.ErrNoRows {
				So(err, ShouldBeNil)
				So(ref, ShouldNotBeNil)
			}
		}

		//Delete
		err = c.Delete()
		So(err, ShouldBeNil)

		err = cl.Delete()
		So(err, ShouldBeNil)

		//clean apiKeyTypes

		err = pub.Delete()
		err = pri.Delete()
		err = aut.Delete()

	})

	Convey("testing general gets", t, func() {
		Convey("Testing GetEtailers()", func() {
			var err error
			dealers, err := GetEtailers()
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
			customers, err := GetWhereToBuyDealers()
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

// func TestCustomerModel(t *testing.T) {
// 	//From the NEW Customer Model
// 	Convey("Testing Customer", t, func() {
// 		var c Customer
// 		dealers, err := GetEtailers()
// 		if err == nil && len(dealers) > 0 {
// 			c = dealers[0]
// 		}

// 		randCust, partID, apiKey, err := getRandomCustWithLocParts()
// 		So(err, ShouldBeNil)
// 		So(partID, ShouldNotBeNil)

// 		Convey("Testing GetCustomer()", func() {
// 			err := c.GetCustomer()
// 			So(err, ShouldBeNil)
// 			So(c.Name, ShouldNotBeNil)
// 		})
// 		Convey("Testing Basics()", func() {
// 			err := c.Basics()
// 			So(err, ShouldBeNil)
// 			So(c.Name, ShouldNotBeNil)
// 		})
// 		Convey("Testing GetLocations()", func() {
// 			err := randCust.GetLocations()
// 			if err != sql.ErrNoRows {
// 				So(err, ShouldBeNil)
// 				So(randCust.Locations, ShouldHaveSameTypeAs, []CustomerLocation{})
// 			}
// 		})
// 		Convey("Testing GetUsers()", func() {
// 			users, err := c.GetUsers()
// 			So(err, ShouldBeNil)
// 			So(users, ShouldHaveSameTypeAs, []CustomerUser{})
// 		})
// 		Convey("Testing GetCustomerPrice()", func() {
// 			price, err := GetCustomerPrice(apiKey, partID)
// 			if err == nil {
// 				So(err, ShouldBeNil)
// 				So(price, ShouldNotBeNil)
// 			} else {
// 				So(err.Error(), ShouldResemble, "sql: no rows in result set")
// 			}

// 		})
// 		Convey("Testing GetCustomerCartReference())", func() {
// 			var err error
// 			ref, err := GetCustomerCartReference(apiKey, partID)
// 			if ref > 0 {
// 				So(err, ShouldBeNil)
// 			}
// 			So(err, ShouldNotBeNil)
// 		})
// 		Convey("Testing GetEtailers()", func() {
// 			var err error
// 			dealers, err := GetEtailers()
// 			So(err, ShouldBeNil)
// 			So(dealers, ShouldHaveSameTypeAs, []Customer{})
// 		})
// 		Convey("Testing GetLocalDealers()", func() {
// 			var err error
// 			latlng := "43.853282,-95.571675,45.800981,-90.468526&"
// 			center := "44.83536,-93.0201"
// 			dealers, err := GetLocalDealers(center, latlng)
// 			So(err, ShouldBeNil)
// 			So(dealers, ShouldHaveSameTypeAs, []DealerLocation{})
// 		})
// 		Convey("Testing GetLocalRegions()", func() {
// 			var err error
// 			regions, err := GetLocalRegions()
// 			So(err, ShouldBeNil)
// 			So(regions, ShouldHaveSameTypeAs, []StateRegion{})
// 		})
// 		Convey("Testing GetLocalDealerTiers()", func() {
// 			var err error
// 			tiers, err := GetLocalDealerTiers()
// 			So(err, ShouldBeNil)
// 			So(tiers, ShouldHaveSameTypeAs, []DealerTier{})
// 		})
// 		Convey("Testing GetLocalDealerTypes()", func() {
// 			var err error
// 			graphics, err := GetLocalDealerTypes()
// 			So(err, ShouldBeNil)
// 			So(graphics, ShouldHaveSameTypeAs, []MapGraphics{})
// 		})
// 		Convey("Testing GetWhereToBuyDealers()", func() {
// 			var err error
// 			customers, err := GetWhereToBuyDealers()
// 			So(err, ShouldBeNil)
// 			So(customers, ShouldHaveSameTypeAs, []Customer{})
// 		})
// 		Convey("Testing GetLocationById()", func() {
// 			var err error
// 			id := 1
// 			location, err := GetLocationById(id)
// 			So(err, ShouldBeNil)
// 			So(location, ShouldHaveSameTypeAs, CustomerLocation{})
// 		})
// 		Convey("Testing SearchLocations()", func() {
// 			var err error
// 			term := "hitch"
// 			locations, err := SearchLocations(term)
// 			So(err, ShouldBeNil)
// 			So(locations, ShouldHaveSameTypeAs, []DealerLocation{})
// 		})
// 		Convey("Testing SearchLocationsByType()", func() {
// 			var err error
// 			term := "hitch"
// 			locations, err := SearchLocationsByType(term)
// 			So(err, ShouldBeNil)
// 			So(locations, ShouldHaveSameTypeAs, DealerLocations{})
// 		})
// 		Convey("Testing SearchLocationsByLatLng()", func() {
// 			var err error
// 			latlng := GeoLocation{
// 				Latitude:  43.853282,
// 				Longitude: -95.571675,
// 			}
// 			locations, err := SearchLocationsByLatLng(latlng)
// 			So(err, ShouldBeNil)
// 			So(locations, ShouldHaveSameTypeAs, []DealerLocation{})
// 		})
// 	})

// 	Convey("Testing User", t, func() {
// 		var c Customer
// 		c.Id = 1

// 		auth_key := ""
// 		userID := ""
// 		api := "8AEE0620-412E-47FC-900A-947820EA1C1D"
// 		users, err := c.GetUsers()
// 		if err == nil && len(users) > 0 {
// 			userID = users[0].Id
// 			if err = users[0].GetKeys(); err == nil {
// 				for _, k := range users[0].Keys {
// 					if strings.ToLower(k.Type) == "public" {
// 						api = k.Key
// 					} else if strings.ToLower(k.Type) == "authentication" {
// 						auth_key = k.Key
// 					}
// 				}
// 			}
// 		}
// 		Convey("Testing UserAuthentication()", func() {
// 			var u CustomerUser
// 			var err error
// 			u.Id = userID
// 			password := "test"
// 			c, err := u.UserAuthentication(password)
// 			So(err, ShouldEqual, AuthError)
// 			So(c, ShouldBeZeroValue)
// 		})
// 		Convey("Testing UserAuthenticationByKey()", func() {
// 			c, err := UserAuthenticationByKey(auth_key)
// 			So(err, ShouldNotBeNil)
// 			So(c, ShouldBeZeroValue)
// 		})
// 		Convey("Testing GetCustomer()", func() {
// 			var u CustomerUser
// 			var err error
// 			u.Id = userID
// 			c, err := u.GetCustomer()
// 			So(err, ShouldBeNil)
// 			So(c.Name, ShouldNotBeNil)
// 			So(c.Id, ShouldNotEqual, 0)
// 		})
// 		Convey("Testing AuthenticateUser()", func() {
// 			var u CustomerUser
// 			var err error
// 			u.Id = userID
// 			password := "wrongPassword"
// 			err = u.AuthenticateUser(password)
// 			So(err, ShouldNotBeNil) //TODO - update user and auth

// 		})

// 		Convey("Testing AuthenticateUserByKey()", func() {
// 			var err error
// 			u, err := AuthenticateUserByKey(auth_key)
// 			if err != nil {
// 				So(err, ShouldEqual, sql.ErrNoRows)
// 				So(u.Id, ShouldEqual, "")
// 			} else {
// 				So(err, ShouldBeNil)
// 				So(u, ShouldNotBeNil)
// 			}
// 			Convey("Testing ResetAuthentication()", func() {
// 				var u CustomerUser
// 				var err error
// 				u.Id = userID
// 				err = u.ResetAuthentication()
// 				if err != nil {
// 					So(err.Error(), ShouldEqual, "faield to retrieve key type reference")
// 				} else {
// 					So(err, ShouldBeNil)
// 				}
// 			})
// 		})
// 		Convey("GetKeys()", func() {
// 			var u CustomerUser
// 			var err error
// 			u.Id = userID
// 			err = u.GetKeys()
// 			So(u.Keys, ShouldHaveSameTypeAs, []ApiCredentials{})
// 			So(err, ShouldBeNil)
// 		})
// 		Convey("GetLocation()", func() {
// 			var u CustomerUser
// 			var err error
// 			u.Id = userID
// 			err = u.GetLocation()
// 			So(u.Location, ShouldNotBeNil)
// 			So(err, ShouldBeNil)
// 		})

// 		Convey("GetCustomerIdFromKey()", func() {
// 			var err error
// 			id, err := GetCustomerIdFromKey(api)
// 			So(id, ShouldNotBeNil)
// 			So(err, ShouldBeNil)
// 		})
// 		Convey("GetCustomerUserFromKey()", func() {
// 			var err error
// 			user, err := GetCustomerUserFromKey(api)
// 			So(user, ShouldNotBeNil)
// 			So(err, ShouldBeNil)
// 		})
// 		Convey("GetCustomerUserFromId()", func() {
// 			if len(users) > 0 {
// 				err = users[0].Get(api)
// 				So(users[0], ShouldNotBeNil)
// 				So(err, ShouldBeNil)
// 			}
// 		})
// 	})
// 	//Comparative Tests - Old Customer Model to New One
// 	Convey("Testing Existing User object to the New One", t, func() {
// 		Convey("Testing GetCustomer()", func() {
// 			var cc customer.CustomerUser
// 			var c CustomerUser
// 			var err error
// 			c.Id = "F43F5B82-D7AF-4905-BD07-FC8BCF4C82FE"
// 			cc.Id = "F43F5B82-D7AF-4905-BD07-FC8BCF4C82FE"
// 			user, err := c.GetCustomer()
// 			user2, err := cc.GetCustomer()
// 			So(err, ShouldBeNil)
// 			So(user.Name, ShouldNotBeNil)
// 			So(user2.Name, ShouldNotBeNil)
// 			So(user.Name, ShouldEqual, user2.Name)
// 			So(user.DealerType.Id, ShouldEqual, user2.DealerType.Id)
// 		})
// 		Convey("Testing Basics()", func() {
// 			var cc customer.Customer
// 			var c Customer
// 			var err error
// 			c.Id = 1
// 			cc.Id = 1
// 			err = c.Basics()
// 			err = cc.Basics()
// 			So(err, ShouldBeNil)
// 			So(c.Name, ShouldNotBeNil)
// 			So(cc.Name, ShouldNotBeNil)
// 			So(cc.Name, ShouldEqual, c.Name)
// 			So(cc.DealerType.Id, ShouldEqual, c.DealerType.Id)
// 		})
// 		Convey("Testing GetLocations()", func() {
// 			var cc customer.Customer
// 			var c Customer
// 			var err error
// 			c.Id = 1
// 			cc.Id = 1
// 			err = c.GetLocations()
// 			err = cc.GetLocations()
// 			So(err, ShouldBeNil)
// 			So(c.Name, ShouldNotBeNil)
// 			So(cc.Name, ShouldNotBeNil)
// 			So(cc.Name, ShouldEqual, c.Name)
// 			So(cc.DealerType.Id, ShouldEqual, c.DealerType.Id)
// 		})
// 		// //The existing GetCustomerPrice errors out - doesn't work to compare
// 		// Convey("Testing GetCustomerPrice()", func() {
// 		// 	var cc customer.Customer
// 		// 	var c Customer
// 		// 	var err error
// 		// 	api, part := getAPIKeyAndPart()
// 		// 	t.Log(api, " ", part)
// 		// 	c.Id = 1
// 		// 	cc.Id = 1
// 		// 	price, err := GetCustomerPrice(api, part)
// 		// 	price2, err := customer.GetCustomerPrice(api, part)
// 		// 	t.Log(customer.GetCustomerPrice(api, part))
// 		// 	So(err, ShouldBeNil)
// 		// 	So(price, ShouldNotBeNil)
// 		// 	So(price2, ShouldNotBeNil)
// 		// 	So(price, ShouldEqual, price2)
// 		// })
// 		Convey("GetLocation()", func() {
// 			var u CustomerUser
// 			var u2 customer.CustomerUser
// 			var err error
// 			u.Id = "023E68B9-9B62-4E5D-84F7-EDB88428B4F8"
// 			u2.Id = "023E68B9-9B62-4E5D-84F7-EDB88428B4F8"
// 			err = u.GetLocation()
// 			So(u.Location, ShouldNotBeNil)
// 			So(err, ShouldBeNil)
// 			err = u2.GetLocation()
// 			So(u2.Location, ShouldNotBeNil)
// 			So(err, ShouldBeNil)
// 			So(u.Location.State.State, ShouldResemble, u2.Location.State.State)
// 		})

// 	})

// 	Convey("Test CRUD Customer", t, func() {
// 		var c Customer
// 		var err error
// 		c.Name = "test"
// 		c.Address = "Nowhere"
// 		err = c.Create()
// 		So(err, ShouldBeNil)

// 		c.Name = "New Name"
// 		c.MapixCode.ID = 1
// 		err = c.Update()
// 		So(err, ShouldBeNil)

// 		err = c.GetCustomer()
// 		So(err, ShouldBeNil)

// 		err = c.Delete()
// 		So(err, ShouldBeNil)

// 	})
// }
