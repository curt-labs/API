package cartIntegration

import (
	"github.com/curt-labs/GoAPI/helpers/apicontextmock"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestCartIntegration(t *testing.T) {
	var err error
	dtx, err := apicontextmock.Mock()
	if err != nil {
		return
	}

	Convey("Testing", t, func() {
		//Create
		cp := CustomerPrice{
			CustID: 1,
			PartID: 11000,
			Price:  123,
			IsSale: 0,
		}

		err = cp.Create()
		So(err, ShouldBeNil)

		//Get
		prices, err := GetCustomerPrices(dtx)
		So(err, ShouldBeNil)

		err = cp.Delete()
		So(err, ShouldBeNil)

	})

	apicontextmock.DeMock(dtx)

}
