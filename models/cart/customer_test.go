package cart

import (
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/mgo.v2/bson"
	"os"
	"testing"
	"time"
)

func TestSinceId(t *testing.T) {
	Convey("Testing Customer Gets with no shop", t, func() {
		Convey("Testing SinceId with bad connection", func() {
			os.Setenv("MONGO_URL", "0.0.0.1")
			custs, err := CustomersSinceId(bson.NewObjectId(), 0, 0, nil, nil, nil, nil)
			So(err, ShouldNotBeNil)
			So(custs, ShouldHaveSameTypeAs, []Customer{})
			os.Setenv("MONGO_URL", "")
		})
		Convey("Testing SinceId with no Id", func() {
			custs, err := CustomersSinceId(bson.NewObjectId(), 0, 0, nil, nil, nil, nil)
			So(err, ShouldBeNil)
			So(custs, ShouldHaveSameTypeAs, []Customer{})
		})
		Convey("Testing SinceId with no Id and created dates", func() {
			created_min, _ := time.Parse(time.RFC3339Nano, time.Now().AddDate(-1, 0, 0).Format(time.RFC3339Nano))
			created_max, _ := time.Parse(time.RFC3339Nano, time.Now().Format(time.RFC3339Nano))

			custs, err := CustomersSinceId(bson.NewObjectId(), 0, 0, &created_min, &created_max, nil, nil)
			So(err, ShouldBeNil)
			So(custs, ShouldHaveSameTypeAs, []Customer{})
		})
		Convey("Testing SinceId with no Id and created, updated dates", func() {
			created_min, _ := time.Parse(time.RFC3339Nano, time.Now().AddDate(-1, 0, 0).Format(time.RFC3339Nano))
			created_max, _ := time.Parse(time.RFC3339Nano, time.Now().Format(time.RFC3339Nano))

			custs, err := CustomersSinceId(bson.NewObjectId(), 1, 0, &created_min, &created_max, &created_min, &created_max)
			So(err, ShouldBeNil)
			So(custs, ShouldHaveSameTypeAs, []Customer{})
		})
	})
}

func TestGetCustomer(t *testing.T) {
	Convey("Testing Customer Gets with no shop", t, func() {
		Convey("Testing GetCustomers with bad connection", func() {
			os.Setenv("MONGO_URL", "0.0.0.1")
			custs, err := GetCustomers(bson.NewObjectId(), 0, 0, nil, nil, nil, nil)
			So(err, ShouldNotBeNil)
			So(custs, ShouldHaveSameTypeAs, []Customer{})
			os.Setenv("MONGO_URL", "")
		})
		Convey("Testing GetCustomers with no Id", func() {
			custs, err := GetCustomers(bson.NewObjectId(), 0, 0, nil, nil, nil, nil)
			So(err, ShouldBeNil)
			So(custs, ShouldHaveSameTypeAs, []Customer{})
		})
		Convey("Testing GetCustomers with no Id and created dates", func() {
			created_min, _ := time.Parse(time.RFC3339Nano, time.Now().AddDate(-1, 0, 0).Format(time.RFC3339Nano))
			created_max, _ := time.Parse(time.RFC3339Nano, time.Now().Format(time.RFC3339Nano))

			custs, err := GetCustomers(bson.NewObjectId(), 0, 0, &created_min, &created_max, nil, nil)
			So(err, ShouldBeNil)
			So(custs, ShouldHaveSameTypeAs, []Customer{})
		})
		Convey("Testing GetCustomers with no Id and created, updated dates", func() {
			created_min, _ := time.Parse(time.RFC3339Nano, time.Now().AddDate(-1, 0, 0).Format(time.RFC3339Nano))
			created_max, _ := time.Parse(time.RFC3339Nano, time.Now().Format(time.RFC3339Nano))

			custs, err := GetCustomers(bson.NewObjectId(), 1, 0, &created_min, &created_max, &created_min, &created_max)
			So(err, ShouldBeNil)
			So(custs, ShouldHaveSameTypeAs, []Customer{})
		})
	})
}

func TestCustomer(t *testing.T) {
	Convey("Testing Get with no Id", t, func() {
		cust := Customer{}
		err := cust.Get()
		So(err, ShouldNotBeNil)
	})

	Convey("Generating a Test Shop and testing Customer functions", t, func() {
		if id := insertTestData(); id != nil {
			shop := Shop{
				Id: *id,
			}
			var customer Customer
			customer.ShopId = shop.Id
			err := customer.Insert()
			So(err, ShouldNotBeNil)

			customer.Email = "ninnemana@gmail.com"
			err = customer.Insert()
			So(err, ShouldNotBeNil)

			customer.FirstName = "Alex"
			err = customer.Insert()
			So(err, ShouldNotBeNil)

			customer.LastName = "Ninneman"
			os.Setenv("MONGO_URL", "0.0.0.1")
			err = customer.Insert()
			So(err, ShouldNotBeNil)
			os.Setenv("MONGO_URL", "")

			err = customer.Insert()
			So(err, ShouldBeNil)

			tmpID := customer.Id
			var blankId bson.ObjectId
			customer.Id = blankId
			customer.FirstName = ""
			customer.LastName = ""
			customer.Email = ""
			err = customer.Update()
			So(err, ShouldNotBeNil)
			customer.Id = tmpID

			customer.FirstName = ""
			customer.LastName = ""
			customer.Email = ""
			err = customer.Update()
			So(err, ShouldNotBeNil)

			customer.FirstName = ""
			customer.LastName = ""
			customer.Email = "ninnemana@gmail.com"
			err = customer.Update()
			So(err, ShouldNotBeNil)

			customer.FirstName = "Alex"
			customer.LastName = ""
			customer.Email = "ninnemana@gmail.com"
			err = customer.Update()
			So(err, ShouldNotBeNil)

			customer.AcceptsMarketing = true
			customer.FirstName = "Alex"
			customer.LastName = "Ninneman"
			customer.Email = "ninnemana@gmail.com"
			addr := CustomerAddress{
				Id:           bson.NewObjectId(),
				Address1:     "1119 Sunset Lane",
				City:         "Altoona",
				Company:      "CURT Manufacturing, LLC",
				Name:         "Alex's House",
				FirstName:    "Alex",
				LastName:     "Ninneman",
				Phone:        "7153082604",
				Province:     "Wisconsin",
				ProvinceCode: "WI",
				Country:      "USA",
				CountryCode:  "US",
				CountryName:  "United States of America",
				Zip:          "54720",
			}
			customer.Addresses = append(customer.Addresses, addr)
			customer.DefaultAddress = &addr
			customer.Note = "Holy shit this is easy"

			os.Setenv("MONGO_URL", "0.0.0.1")
			err = customer.Update()
			So(err, ShouldNotBeNil)
			os.Setenv("MONGO_URL", "")

			err = customer.Update()
			So(err, ShouldBeNil)

			os.Setenv("MONGO_URL", "0.0.0.1")
			err = customer.Get()
			So(err, ShouldNotBeNil)
			os.Setenv("MONGO_URL", "")

			os.Setenv("MONGO_URL", "0.0.0.1")
			err = customer.Delete()
			So(err, ShouldNotBeNil)
			os.Setenv("MONGO_URL", "")

			customer.Id = bson.NewObjectId()
			err = customer.Delete()
			So(err, ShouldNotBeNil)

			customer.Id = blankId
			err = customer.Delete()
			So(err, ShouldNotBeNil)
			customer.Id = tmpID

			err = customer.Delete()
			So(err, ShouldBeNil)
		}
	})
}

func TestCustomerCount(t *testing.T) {
	Convey("Test CustomerCount", t, func() {
		Convey("Bad Connection", func() {
			os.Setenv("MONGO_URL", "0.0.0.1")
			count, err := CustomerCount(bson.NewObjectId())
			So(err, ShouldNotBeNil)
			So(count, ShouldEqual, 0)
			os.Setenv("MONGO_URL", "")
		})
		Convey("Empty Shop ID", func() {
			var id bson.ObjectId
			count, err := CustomerCount(id)
			So(err, ShouldNotBeNil)
			So(count, ShouldEqual, 0)
		})
		Convey("With Shop ID", func() {
			count, err := CustomerCount(bson.NewObjectId())
			So(err, ShouldBeNil)
			So(count, ShouldEqual, 0)
		})
	})
}
