package cart

import (
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/mgo.v2/bson"
	"os"
	"testing"
)

func TestOrderValidation(t *testing.T) {
	Convey("Testing validate()", t, func() {
		Convey("with no data", func() {
			var o Order
			err := o.validate()
			So(err, ShouldNotBeNil)
		})

		Convey("with no variant_id", func() {
			var o Order
			o.LineItems = append(o.LineItems, LineItem{
				Quantity: 1,
			})
			err := o.validate()
			So(err, ShouldNotBeNil)
		})

		Convey("with billing and no email", func() {
			var o Order
			o.LineItems = append(o.LineItems, LineItem{
				Quantity:  1,
				VariantId: 1000,
			})
			o.BillingAddress = &CustomerAddress{}
			err := o.validate()
			So(err, ShouldNotBeNil)
		})

		Convey("with no billing and no email", func() {
			var o Order
			o.LineItems = append(o.LineItems, LineItem{
				Quantity:  1,
				VariantId: 1000,
			})
			err := o.validate()
			So(err, ShouldBeNil)
		})

		Convey("with billing and email", func() {
			var o Order
			o.LineItems = append(o.LineItems, LineItem{
				Quantity:  1,
				VariantId: 1000,
			})
			o.BillingAddress = &CustomerAddress{}
			o.Email = "test@example.com"
			err := o.validate()
			So(err, ShouldBeNil)
		})
	})
}

func TestOrderCount(t *testing.T) {
	Convey("Testing GetOrderCount", t, func() {
		Convey("with bad connection", func() {
			os.Setenv("MONGO_URL", "0.0.0.1")
			count, err := getOrderCount(bson.NewObjectId())
			So(err, ShouldNotBeNil)
			So(count, ShouldEqual, 0)
			os.Setenv("MONGO_URL", "")
		})

		Convey("with good connection", func() {
			count, err := getOrderCount(bson.NewObjectId())
			So(err, ShouldBeNil)
			So(count, ShouldEqual, 0)
		})
	})
}

func TestBindCustomer(t *testing.T) {
	Convey("Testing bindCustomer()", t, func() {
		var o Order
		o.Customer = &Customer{}
		o.Email = "test@example.com"
		err := o.bindCustomer()
		So(err, ShouldNotBeNil)

		o.Email = "test@example.com"
		o.ShopId = bson.NewObjectId()
		err = o.bindCustomer()
		So(err, ShouldNotBeNil)

		if id := InsertTestData(); id != nil {
			o.ShopId = *id
			o.Customer.ShopId = o.ShopId
			o.Customer.Email = "ninnemana@gmail.com"
			o.Customer.FirstName = "Alex"
			o.Customer.LastName = "Ninneman"
			o.Customer.Password = "password"
			err = o.Customer.Insert("http://www.example.com")
			So(err, ShouldBeNil)

			So(o.bindCustomer(), ShouldBeNil)

			t.Log(o.Customer.Email)
			o.Email = o.Customer.Email
			o.Customer = nil
			So(o.bindCustomer(), ShouldBeNil)
		}

		o.Customer.Email = ""
		So(err, ShouldBeNil)
	})
}

func TestCreate(t *testing.T) {
	Convey("Testing Order.Create()", t, func() {
		var o Order
		So(o.Create(), ShouldNotBeNil)

		o.LineItems = append(o.LineItems, LineItem{
			Quantity:  1,
			VariantId: 1000,
		})
		So(o.Create(), ShouldNotBeNil)

		if id := InsertTestData(); id != nil {
			o.ShopId = *id
			So(o.Create(), ShouldBeNil)
		}
	})
}
