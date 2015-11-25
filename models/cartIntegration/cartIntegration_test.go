package cartIntegration

import (
	"github.com/curt-labs/API/helpers/database"
	_ "github.com/go-sql-driver/mysql"

	"database/sql"
	. "github.com/smartystreets/goconvey/convey"

	"encoding/csv"
	"mime/multipart"
	"os"
	"testing"
)

func TestCartIntegration(t *testing.T) {
	var err error
	Brand_ID = 1
	Customer_ID = 1

	key, _ := getCustomerKey()

	Convey("Testing CustomerPrices", t, func() {
		cp := CustomerPrice{
			CustID: 1,
			PartID: 11000,
			Price:  123,
			IsSale: 0,
		}

		err = cp.Create()
		So(err, ShouldBeNil)

		cp.CustomerPartID = 1
		err = cp.Update()
		So(err, ShouldBeNil)

		custprices, err := GetCustomerPrices()
		So(err, ShouldBeNil)
		So(len(custprices), ShouldBeGreaterThan, 0)

		custprices, err = GetPricingPaged(1, 1)
		So(err, ShouldBeNil)
		So(len(custprices), ShouldBeGreaterThan, 0)

		count, err := GetPricingCount()
		So(err, ShouldBeNil)
		So(count, ShouldBeGreaterThan, 0)

		prices, err := GetPartPrices()
		So(err, ShouldBeNil)
		So(len(prices), ShouldBeGreaterThan, 0)

		prices, err = GetPartPricesByPartID(cp.PartID)
		So(err, ShouldBeNil)
		So(len(prices), ShouldBeGreaterThanOrEqualTo, 1)

		prices, err = GetMAPPartPrices()
		So(err, ShouldBeNil)
		So(len(prices), ShouldBeGreaterThanOrEqualTo, 1)

		err = cp.Delete()
		So(err, ShouldBeNil)
	})

	Convey("Testing CartIntegrations", t, func() {
		cp := CustomerPrice{
			CustID:         1,
			PartID:         11000,
			CustomerPartID: 200,
			Price:          123.00,
			IsSale:         0,
		}

		err = cp.InsertCartIntegration()
		So(err, ShouldBeNil)

		err = cp.UpdateCartIntegration()
		So(err, ShouldBeNil)

		custprices, err := GetCustomerCartIntegrations(key)
		So(err, ShouldBeNil)
		So(len(custprices), ShouldBeGreaterThanOrEqualTo, 0)

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

		err = UploadFile(file, key)
		So(err, ShouldBeNil)

		os.Remove("test.csv") //cleanup

	})

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

func getCustomerKey() (string, error) {
	var key string
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return key, err
	}
	defer db.Close()

	stmt, err := db.Prepare("select api_key from ApiKey a join CustomerUser c on c.id = a.user_id join ApiKeyType at on at.id = a.type_id where c.cust_id = 1 and at.type = 'Public' limit 1")
	if err != nil {
		return key, err
	}
	defer stmt.Close()
	err = stmt.QueryRow().Scan(&key)
	return key, err
}
