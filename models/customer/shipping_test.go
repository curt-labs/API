package customer

import (
	"github.com/curt-labs/GoAPI/helpers/apicontext"
	"github.com/curt-labs/GoAPI/helpers/apicontextmock"
	. "github.com/smartystreets/goconvey/convey"
	"os"
	"testing"
)

func TestShipping(t *testing.T) {
	var err error
	MockedDTX := &apicontext.DataContext{}
	if MockedDTX, err = apicontextmock.Mock(); err != nil {
		return
	}
	Convey("Testing Customer 1 Mapics Shipping", t, func() {
		c := Customer{
			Id: 1,
		}
		err = c.GetShippingInfo()
		So(err, ShouldBeNil)
		So(c.ShippingInfo.CustomerInfo.CustName, ShouldEqual, "INVALID Customer Number")
	})

	Convey("Testing getWarehouses", t, func() {
		warehouses, err := getWarehouseCodes()
		So(err, ShouldBeNil)
		So(warehouses, ShouldHaveSameTypeAs, map[string]int{})
	})

	Convey("Testing Customer 1 Adjust to Mapics - DON'T RUN LIVE", t, func() {
		addr := os.Getenv("DATABASE_HOST")
		if addr == "" {
			c := Customer{
				Id: 1,
			}
			err := c.CompareCustomerShippingInfo()
			So(err, ShouldBeNil)
		}
	})

	_ = apicontextmock.DeMock(MockedDTX)
}
