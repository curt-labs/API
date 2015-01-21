package customer

import (
	"github.com/curt-labs/GoAPI/helpers/apicontextmock"
	_ "github.com/go-sql-driver/mysql"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestCustomerLocations(t *testing.T) {
	var l CustomerLocation
	var err error

	MockedDTX, err := apicontextmock.Mock()
	if err != nil {
		return
	}
	Convey("Testing Locations", t, func() {

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
			l.CustomerId = MockedDTX.CustomerID
			l.ContactPerson = "Tes"
			l.IsPrimary = true
			l.PostalCode = "Tes"
			l.ShippingDefault = false
			err := l.Create()
			So(err, ShouldBeNil)
		})

		Convey("Update", func() {
			l.Name = "Chuck Norris"
			err := l.Update()
			So(err, ShouldBeNil)
			So(l.Name, ShouldNotEqual, "test")

		})

		Convey("Testing GetAll()", func() {
			locations, err := GetAllLocations(MockedDTX)
			So(len(locations), ShouldBeGreaterThanOrEqualTo, 0)
			So(err, ShouldBeNil)
		})
		Convey("Get", func() {
			err = l.Get()
			So(err, ShouldBeNil)
			So(l, ShouldNotBeNil)
		})

		Convey("Delete", func() {
			err := l.Delete()
			So(err, ShouldBeNil)

		})

	})
	_ = apicontextmock.DeMock(MockedDTX)
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
