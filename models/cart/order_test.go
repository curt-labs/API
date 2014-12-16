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
