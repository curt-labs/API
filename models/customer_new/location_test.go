package customer_new

import (
	_ "github.com/go-sql-driver/mysql"
	. "github.com/smartystreets/goconvey/convey"
	"math/rand"

	"testing"
)

func TestCustomerLocations(t *testing.T) {
	Convey("Testing Locations", t, func() {
		Convey("Testing GetAll()", func() {
			Locations, err := GetAllLocations()
			So(Locations, ShouldNotBeNil)
			So(err, ShouldBeNil)
			x := rand.Intn(len(Locations))
			Location := Locations[x]
			var testLoc CustomerLocation
			testLoc.Id = Location.Id

			Convey("Get", func() {
				err = testLoc.Get()
				So(err, ShouldBeNil)
				So(testLoc, ShouldNotBeNil)
			})
			Convey("Get Bad stmt", func() {
				getLocation = "bad"
				err = testLoc.Get()
				So(err, ShouldNotBeNil)
			})
			Convey("GetAll Bad stmt", func() {
				getLocations = "bad"
				locs, err := GetAllLocations()
				So(err, ShouldNotBeNil)
				So(locs, ShouldBeNil)
			})

		})
		var l CustomerLocation
		Convey("Testing Create", func() {
			l.Name = "test"
			l.Address = "testA"
			l.City = "Tes"
			l.State.Id = 12
			l.IsPrimary = true
			l.Email = "Tes"
			l.Fax = "Tes"
			l.Phone = "Tes"
			l.Latitude = 44.913687
			l.Longitude = -91.89981
			l.CustomerId = 1
			l.ContactPerson = "Tes"
			l.IsPrimary = true
			l.PostalCode = "Tes"
			l.ShippingDefault = false
			err := l.Create()
			So(err, ShouldBeNil)

			Convey("Update", func() {
				l.Name = "Chuck Norris"
				err := l.Update()
				So(err, ShouldBeNil)
				So(l.Name, ShouldNotEqual, "test")

			})
			Convey("Delete", func() {
				err := l.Delete()
				So(err, ShouldBeNil)

			})
		})
	})

}

func TestBadCrudStmts(t *testing.T) {
	Convey("Testing bad statements", t, func() {
		var l CustomerLocation
		l.Name = "test"
		l.Address = "testA"
		l.City = "Tes"
		l.State.Id = 12
		l.IsPrimary = true
		createLocation = "Bad Query Stmt"
		updateLocation = "Bad Query Stmt"
		deleteLocation = "Bad Query Stmt"
		Convey("Testing Bad Create()", func() {
			err := l.Create()
			So(err, ShouldNotBeNil)
		})

		Convey("Testing Bad Update()", func() {
			l.Name = "test2"
			err := l.Update()
			So(err, ShouldNotBeNil)
		})
		Convey("Testing Bad Delete()", func() {
			err := l.Delete()
			So(err, ShouldNotBeNil)
		})

	})
}
