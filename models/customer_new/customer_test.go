package customer_new

import (
	"database/sql"
	"github.com/curt-labs/GoAPI/helpers/database"
	"github.com/curt-labs/GoAPI/models/customer"
	_ "github.com/go-sql-driver/mysql"
	. "github.com/smartystreets/goconvey/convey"
	"math/rand"
	"strings"
	"testing"
	"time"
)

const (
	inputTimeFormat = "01/02/2006"
)

var (
	getCustWithLocationAndParts = `SELECT customerID, cp.partID, apiKey FROM Customer AS c
								LEFT JOIN CustomerLocations AS cl on cl.cust_id = c.CustomerID
								LEFT JOIN CustomerPricing AS cp ON cp.cust_id = c.CustomerID
								WHERE cl.locationID IS NOT NULL
								AND cp.partID IS NOT NULL`
)

func getRandomCustWithLocParts() (cust Customer, partID int, apiKey string, err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return cust, partID, apiKey, err
	}
	defer db.Close()

	stmt, err := db.Prepare(getCustWithLocationAndParts)
	if err != nil {
		return cust, partID, apiKey, err
	}
	defer stmt.Close()
	var cs []Customer
	var api string
	res, err := stmt.Query()
	for res.Next() {
		var c Customer
		err = res.Scan(&c.Id, &partID, &api)
		cs = append(cs, c)
	}
	if len(cs) > 0 {
		x := rand.Intn(len(cs))
		cust = cs[x]
	}

	users, err := cust.GetUsers()
	if err == nil && len(users) > 0 {
		if err = users[0].GetKeys(); err == nil {
			for _, k := range users[0].Keys {
				if strings.ToLower(k.Type) == "public" {
					apiKey = k.Key
					break
				}
			}
		}
	}

	return cust, partID, apiKey, err
}

func TestCustomerPriceModel(t *testing.T) {
	Convey("Testing Price - Gets", t, func() {
		Convey("Testing GetAllPrices()", func() {
			ps, err := GetAllPrices()
			So(len(ps), ShouldBeGreaterThan, 200000)
			So(err, ShouldBeNil)
		})
		Convey("Gets random CustomerPrice", func() {
			ps, err := GetAllPrices()
			So(err, ShouldBeNil)
			if len(ps) > 0 {
				x := rand.Intn(len(ps))
				p := ps[x]

				Convey("Testing Get()", func() {
					err := p.Get()
					So(p.Price, ShouldHaveSameTypeAs, 0.00)
					So(p, ShouldNotBeNil)
					So(err, ShouldBeNil)
				})

				Convey("Testing GetPricesByCustomer()", func() {
					var c Customer
					c.Id = p.CustID
					custPrices, err := c.GetPricesByCustomer()
					So(custPrices, ShouldNotBeNil)
					So(err, ShouldBeNil)
				})
				Convey("Testing GetPricesByPart()", func() {
					partID := p.PartID
					prices, err := GetPricesByPart(partID)
					So(len(prices), ShouldNotBeNil)
					So(err, ShouldBeNil)
				})
				Convey("Testing GetPricesBySaleRange", func() {
					var s time.Time
					var e time.Time
					c := Customer{Id: p.CustID}
					var err error
					s, err = time.Parse(inputTimeFormat, "2006-01-02 15:04:05")
					e, err = time.Parse(inputTimeFormat, "2016-01-02 15:04:05")
					prices, err := c.GetPricesBySaleRange(s, e)
					So(err, ShouldBeNil)
					So(len(prices), ShouldBeGreaterThan, 0)
					So(prices, ShouldNotBeNil)
				})

				Convey("Testing Price -  CUD", func() {
					Convey("Testing Create() Update() Delete() Price", func() {
						var pr Price
						var err error
						pr.CustID = p.CustID
						pr.SaleEnd, err = time.Parse(inputTimeFormat, "02/12/2006")
						pr.IsSale = 1
						pr.Price = 666.00
						err = pr.Create()
						So(err, ShouldBeNil)
						pr.SaleStart, err = time.Parse(inputTimeFormat, "01/02/2007")
						err = pr.Update()
						So(err, ShouldBeNil)
						err = pr.Get()
						So(err, ShouldBeNil)
						t, err := time.Parse(inputTimeFormat, "02/12/2006")
						So(pr.SaleStart, ShouldHaveSameTypeAs, t)
						err = pr.Delete()
						So(err, ShouldBeNil)
					})

				})
			}
		})
	})
}

func TestCustomerModel(t *testing.T) {
	//From the NEW Customer Model
	Convey("Testing Customer", t, func() {
		var c Customer
		dealers, err := GetEtailers()
		if err == nil && len(dealers) > 0 {
			c = dealers[0]
		}

		randCust, partID, apiKey, err := getRandomCustWithLocParts()
		So(err, ShouldBeNil)
		So(partID, ShouldNotBeNil)

		Convey("Testing GetCustomer()", func() {
			err := c.GetCustomer()
			So(err, ShouldBeNil)
			So(c.Name, ShouldNotBeNil)
		})
		Convey("Testing Basics()", func() {
			err := c.Basics()
			So(err, ShouldBeNil)
			So(c.Name, ShouldNotBeNil)
		})
		Convey("Testing GetLocations()", func() {
			err := randCust.GetLocations()
			So(err, ShouldBeNil)
			So(len(randCust.Locations), ShouldBeGreaterThan, 0)
		})
		Convey("Testing GetUsers()", func() {
			users, err := c.GetUsers()
			So(err, ShouldBeNil)
			So(users, ShouldNotBeNil)
		})
		Convey("Testing GetCustomerPrice()", func() {

			price, err := GetCustomerPrice(apiKey, partID)
			So(err, ShouldBeNil)
			So(price, ShouldNotBeNil)
		})
		Convey("Testing GetCustomerCartReference())", func() {
			var err error
			ref, err := GetCustomerCartReference(apiKey, partID)
			if ref > 0 {
				So(err, ShouldBeNil)
			}
			So(err, ShouldNotBeNil)
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
		var c Customer
		c.Id = 1

		auth_key := ""
		userID := ""
		api := "8AEE0620-412E-47FC-900A-947820EA1C1D"
		users, err := c.GetUsers()
		if err == nil && len(users) > 0 {
			userID = users[0].Id
			if err = users[0].GetKeys(); err == nil {
				for _, k := range users[0].Keys {
					if strings.ToLower(k.Type) == "public" {
						api = k.Key
					} else if strings.ToLower(k.Type) == "authentication" {
						auth_key = k.Key
					}
				}
			}
		}
		Convey("Testing UserAuthentication()", func() {
			var u CustomerUser
			var err error
			u.Id = userID
			password := "test"
			c, err := u.UserAuthentication(password)
			So(err, ShouldEqual, AuthError)
			So(c, ShouldBeZeroValue)
		})
		Convey("Testing UserAuthenticationByKey()", func() {
			c, err := UserAuthenticationByKey(auth_key)
			So(err, ShouldNotBeNil)
			So(c, ShouldBeZeroValue)
		})
		Convey("Testing GetCustomer()", func() {
			var u CustomerUser
			var err error
			u.Id = userID
			c, err := u.GetCustomer()
			So(err, ShouldBeNil)
			So(c.Name, ShouldNotBeNil)
			So(c.Name, ShouldNotEqual, "")
			So(c.Id, ShouldNotEqual, 0)
		})
		Convey("Testing AuthenticateUser()", func() {
			var u CustomerUser
			var err error
			u.Id = userID
			password := "wrongPassword"
			err = u.AuthenticateUser(password)
			So(err, ShouldNotBeNil) //TODO - update user and auth

		})

		Convey("Testing AuthenticateUserByKey()", func() {
			var err error
			u, err := AuthenticateUserByKey(auth_key)
			if err != nil {
				So(err, ShouldEqual, AuthError)
				So(u.Id, ShouldEqual, "")
			} else {
				So(err, ShouldBeNil)
				So(u, ShouldNotBeNil)
			}
			Convey("Testing ResetAuthentication()", func() {
				var u CustomerUser
				var err error
				u.Id = userID
				err = u.ResetAuthentication()
				So(err, ShouldBeNil)
			})
		})
		Convey("GetKeys()", func() {
			var u CustomerUser
			var err error
			u.Id = userID
			err = u.GetKeys()
			So(u.Keys, ShouldNotBeNil)
			So(err, ShouldBeNil)
		})
		Convey("GetLocation()", func() {
			var u CustomerUser
			var err error
			u.Id = userID
			err = u.GetLocation()
			So(u.Location, ShouldNotBeNil)
			So(err, ShouldBeNil)
		})

		Convey("GetCustomerIdFromKey()", func() {
			var err error
			id, err := GetCustomerIdFromKey(api)
			So(id, ShouldNotBeNil)
			So(err, ShouldBeNil)
		})
		Convey("GetCustomerUserFromKey()", func() {
			var err error
			user, err := GetCustomerUserFromKey(api)
			So(user, ShouldNotBeNil)
			So(err, ShouldBeNil)
		})
		Convey("GetCustomerUserFromId()", func() {
			user, err := GetCustomerUserById(users[0].Id)
			So(user, ShouldNotBeNil)
			So(err, ShouldBeNil)
		})
	})
	//Comparative Tests - Old Customer Model to New One
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
