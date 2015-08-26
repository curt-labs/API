package cartIntegration

import (
	"github.com/curt-labs/GoAPI/helpers/apicontextmock"
	. "github.com/smartystreets/goconvey/convey"

	"encoding/csv"
	"mime/multipart"
	"os"
	"testing"
)

func TestCartIntegration(t *testing.T) {
	var err error
	dtx, err := apicontextmock.Mock()
	if err != nil {
		return
	}

	Convey("Testing CustomerPrices", t, func() {
		cp := CustomerPrice{
			CustID: dtx.CustomerID,
			PartID: 11000,
			Price:  123,
			IsSale: 0,
		}

		err = cp.Create()
		So(err, ShouldBeNil)

		cp.CustomerPartID = 1
		err = cp.Update()
		So(err, ShouldBeNil)

		custprices, err := GetCustomerPrices(dtx)
		So(err, ShouldBeNil)
		So(len(custprices), ShouldBeGreaterThan, 0)

		custprices, err = GetPricingPaged(1, 1, dtx)
		So(err, ShouldBeNil)
		So(len(custprices), ShouldBeGreaterThan, 0)

		count, err := GetPricingCount(dtx)
		So(err, ShouldBeNil)
		So(count, ShouldBeGreaterThan, -1)

		prices, err := GetPartPrices(dtx)
		So(err, ShouldBeNil)
		So(len(prices), ShouldBeGreaterThan, 0)

		prices, err = GetPartPricesByPartID(cp.PartID, dtx)
		So(err, ShouldBeNil)
		So(len(prices), ShouldBeGreaterThan, 0)

		prices, err = GetMAPPartPrices(dtx)
		So(err, ShouldBeNil)
		So(len(prices), ShouldBeGreaterThan, 0)

		err = cp.Delete()
		So(err, ShouldBeNil)
	})

	Convey("Testing CartIntegrations", t, func() {
		cp := CustomerPrice{
			CustID:         dtx.CustomerID,
			PartID:         11000,
			CustomerPartID: 200,
			Price:          123.00,
			IsSale:         0,
		}

		err = cp.InsertCartIntegration()
		So(err, ShouldBeNil)

		err = cp.UpdateCartIntegration()
		So(err, ShouldBeNil)

		custprices, err := GetCustomerCartIntegrations(dtx)
		So(err, ShouldBeNil)
		So(len(custprices), ShouldBeGreaterThan, 0)

		err = cp.DeleteCartIntegration()
		So(err, ShouldBeNil)
	})

	Convey("Testing PriceTypes", t, func() {
		types, err := GetAllPriceTypes()
		So(err, ShouldBeNil)
		So(len(types), ShouldBeGreaterThanOrEqualTo, 1)
	})

	Convey("Testing FileIO", t, func() {
		file, err := newfileUploadRequest()
		So(err, ShouldBeNil)
		t.Log(file.Read(nil))

		err = UploadFile(file, dtx)
		So(err, ShouldBeNil)

		os.Remove("test.csv") //cleanup

	})
	apicontextmock.DeMock(dtx)

}

func newfileUploadRequest() (multipart.File, error) {
	file, err := os.Create("test.csv")
	if err != nil {
		return nil, err
	}

	writer := csv.NewWriter(file)
	err = writer.WriteAll([][]string{{"11000", "201", "100.00", "2020-01-01", "2021-01-01"}, {"11001", "202", "100.00", "2020-01-01", "2021-01-01"}})
	if err != nil {
		return nil, err
	}

	return file, nil
}
