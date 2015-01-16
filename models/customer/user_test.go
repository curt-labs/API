package customer

import (
	"github.com/curt-labs/GoAPI/helpers/apicontextmock"
	_ "github.com/go-sql-driver/mysql"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestUser(t *testing.T) {
	var cu CustomerUser
	var err error
	var cust Customer
	cust.Name = "IMA TESTER"
	dtx, err := apicontextmock.Mock()
	if err != nil {
		t.Log(err)
	}

	Convey("Testing Create", t, func() {
		//make customer to use
		err = cust.Create()
		So(err, ShouldBeNil)

		cu.Name = "testname"
		cu.Email = "test@test.com"
		cu.Password = "test"
		cu.Active = true
		cu.Current = true
		cu.Sudo = true
		cu.CustomerID = cust.Id
		cu.Location.Id = 1
		err = cu.Create(dtx.BrandArray)
		if err != nil {
			errorString := "failed to retrieve auth type"
			So(err.Error(), ShouldEqual, errorString)
		} else {
			So(err, ShouldBeNil)
			So(len(cu.Keys), ShouldEqual, 3)
		}
		err = cu.AuthenticateUser(dtx.BrandArray)
		So(err, ShouldBeNil)
	})
	Convey("Testing Update", t, func() {
		cu.Name = "new name"
		err = cu.UpdateCustomerUser()
		So(err, ShouldBeNil)

	})
	Convey("Testing Get", t, func() {
		c, err := cu.GetCustomer("")
		if cu.Id == "" {
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "error: user not bound to customer")
		} else {
			So(err, ShouldBeNil)
			So(c, ShouldNotBeNil)
		}
	})
	Convey("Testing Auth", t, func() {
		var authKey string
		var pubKey string
		for _, key := range cu.Keys {
			if key.Type == "AUTHENTICATION" {
				authKey = key.Key
			}
			if key.Type == "PUBLIC" {
				pubKey = key.Key
			}
		}

		err := cu.Get(authKey)
		if authKey == "" {
			So(err.Error(), ShouldEqual, "error: user does not exist")
		} else {
			So(err, ShouldBeNil)
		}

		user, err := AuthenticateUserByKey(authKey, dtx)
		if authKey == "" {
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "error: user does not exist")
		} else {
			//So(err, ShouldBeNil) // bad test, user needs to be created for this test to pass
			So(user, ShouldNotBeNil)
		}

		customer, err := AuthenticateAndGetCustomer(authKey, dtx)
		if authKey == "" {
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "failed to authenticate")
		} else {
			//So(err, ShouldBeNil) // bad test user needs to be created for this test to pass
			So(customer, ShouldNotBeNil)
		}

		user2, err := GetCustomerUserFromKey(pubKey)
		if pubKey == "" {
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "error: user does not exist")
		} else {
			So(err, ShouldBeNil)
			So(user2, ShouldNotBeNil)
		}

		cust, err := user2.GetCustomer(pubKey)
		if user2.Id == "" {
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "error: user not bound to customer")
			So(cust.Id, ShouldEqual, 0)
		} else {
			So(err, ShouldBeNil)
			So(cust.Id, ShouldBeGreaterThan, 0)
		}

		err = user.GetKeys()
		So(err, ShouldBeNil)

		err = user.GetLocation()
		So(err, ShouldBeNil)

		err = user.ResetAuthentication(dtx.BrandArray)
		if user.Id == "" {
			So(err, ShouldNotBeNil)
			//So(err.Error(), ShouldEqual, "error: failed to retrieve key type reference")
		} else {
			So(err, ShouldBeNil)
		}

		err = cu.ChangePass(cu.Password, cu.Password, dtx)
		So(err, ShouldBeNil)

		randPassword, err := cu.ResetPass()
		So(err, ShouldBeNil)
		So(randPassword, ShouldNotEqual, cu.Password)

		err = DeleteCustomerUsersByCustomerID(cust.Id)
		So(err, ShouldBeNil)

		//cleanup
		user.Delete()
		user2.Delete()
		customer.Delete()

	})
	Convey("Testing Delete", t, func() {
		err = cu.Delete()
		So(err, ShouldBeNil)
	})

	//cleanup
	cust.Delete()
	_ = apicontextmock.DeMock(dtx)
}
