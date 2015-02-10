package cart

import (
	. "github.com/smartystreets/goconvey/convey"
	"os"
	"testing"
)

func Test_IdentifierFromToken(t *testing.T) {
	clearMongo()

	id := InsertTestData()
	if id == nil {
		return
	}

	var customer Customer
	customer.ShopId = *id
	customer.FirstName = "Alex"
	customer.LastName = "Ninneman"
	customer.Password = "password"
	customer.Email = "ninnemana@gmail.com"
	err := customer.Insert("http://example.com")
	if err != nil {
		t.Log(err)
		return
	}

	Convey("Test Account Lookup via Token", t, func() {
		Convey("with empty token", func() {
			id, err := IdentifierFromToken("")
			So(id, ShouldBeEmpty)
			So(err, ShouldNotBeNil)
		})

		Convey("with bad connection", func() {
			os.Setenv("MONGO_URL", "0.0.0.1")
			id, err := IdentifierFromToken(customer.Token)
			So(id, ShouldBeEmpty)
			So(err, ShouldNotBeNil)
			os.Setenv("MONGO_URL", "")
		})

		Convey("with good token and good connection", func() {
			id, err := IdentifierFromToken(customer.Token)
			So(id, ShouldNotBeEmpty)
			So(err, ShouldBeNil)
		})

	})
}

func Test_AuthenticateAccount(t *testing.T) {
	clearMongo()

	id := InsertTestData()
	if id == nil {
		return
	}

	var customer Customer
	customer.ShopId = *id
	customer.FirstName = "Alex"
	customer.LastName = "Ninneman"
	customer.Password = "password"
	customer.Email = "ninnemana@gmail.com"
	err := customer.Insert("http://example.com")
	if err != nil {
		t.Log(err)
		return
	}

	Convey("Test Account Lookup via Token", t, func() {
		Convey("with empty token", func() {
			cust, err := AuthenticateAccount("")
			So(cust.Id, ShouldBeEmpty)
			So(err, ShouldNotBeNil)
		})

		Convey("with bad connection", func() {
			os.Setenv("MONGO_URL", "0.0.0.1")
			cust, err := AuthenticateAccount(customer.Token)
			So(cust.Id, ShouldBeEmpty)
			So(err, ShouldNotBeNil)
			os.Setenv("MONGO_URL", "")
		})

		Convey("with good token and good connection", func() {
			cust, err := AuthenticateAccount(customer.Token)
			So(cust.Id, ShouldEqual, customer.Id)
			So(err, ShouldBeNil)
		})

	})
}
