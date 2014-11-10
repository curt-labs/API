package customer_new

import (
	_ "github.com/go-sql-driver/mysql"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestUser(t *testing.T) {
	var cu CustomerUser
	var err error
	var cust Customer
	cust.Name = "IMA TESTER"

	Convey("Testing Create", t, func() {
		//make customer to use
		err = cust.Create()
		So(err, ShouldBeNil)
		err = cust.Get()
		So(err, ShouldBeNil)

		cu.Name = "testname"
		cu.Email = "test@test.com"
		cu.Password = "test"
		cu.Active = true
		cu.Current = true
		cu.Sudo = true
		cu.CustomerID = cust.Id
		cu.Location.Id = 1
		err = cu.Create()
		So(err, ShouldBeNil)
		So(len(cu.Keys), ShouldEqual, 3)

		err = cu.AuthenticateUser()
		So(err, ShouldBeNil)
	})
	Convey("Testing Update", t, func() {
		cu.Name = "new name"
		err = cu.UpdateCustomerUser()
		So(err, ShouldBeNil)

	})
	Convey("Testing Get", t, func() {
		c, err := cu.GetCustomer()
		So(err, ShouldBeNil)
		So(c, ShouldNotBeNil)

	})
	Convey("Testing Auth", t, func() {
		var authKey string
		for _, key := range cu.Keys {
			if key.Type == "AUTHENTICATION" {
				authKey = key.Key
			}
		}

		err := cu.Get(authKey)
		So(err, ShouldBeNil)

		user, err := AuthenticateUserByKey(authKey)
		So(err, ShouldBeNil)
		So(user, ShouldNotBeNil)

		customer, err := AuthenticateAndGetCustomer(authKey)
		So(err, ShouldBeNil)
		So(customer, ShouldNotBeNil)

		user2, err := GetCustomerUserFromKey(authKey)
		So(err, ShouldBeNil)
		So(user2, ShouldNotBeNil)

		cust, err := user2.GetCustomer()
		So(err, ShouldBeNil)
		So(cust.Id, ShouldBeGreaterThan, 0)

		err = user.GetKeys()
		So(err, ShouldBeNil)

		err = user.GetLocation()
		So(err, ShouldBeNil)

		err = user.ResetAuthentication()
		So(err, ShouldBeNil)

		err = cu.ChangePass(cu.Password, cu.Password)
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
}
