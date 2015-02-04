package cart

import (
	"github.com/curt-labs/GoAPI/helpers/database"
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"os"
	"testing"
	"time"
)

func TestSinceId(t *testing.T) {
	clearMongo()

	Convey("Testing CustomerSinceId with no shop", t, func() {
		Convey("with bad connection", func() {
			os.Setenv("MONGO_URL", "0.0.0.1")
			custs, err := CustomersSinceId(bson.NewObjectId(), bson.NewObjectId(), 0, 0, nil, nil, nil, nil)
			So(err, ShouldNotBeNil)
			So(custs, ShouldHaveSameTypeAs, []Customer{})
			os.Setenv("MONGO_URL", "")
		})
		Convey("with no Id", func() {
			custs, err := CustomersSinceId(bson.NewObjectId(), bson.NewObjectId(), 0, 0, nil, nil, nil, nil)
			So(err, ShouldBeNil)
			So(custs, ShouldHaveSameTypeAs, []Customer{})
		})
		Convey("with no Id and created dates", func() {
			created_min, _ := time.Parse(time.RFC3339Nano, time.Now().AddDate(-1, 0, 0).Format(time.RFC3339Nano))
			created_max, _ := time.Parse(time.RFC3339Nano, time.Now().Format(time.RFC3339Nano))

			custs, err := CustomersSinceId(bson.NewObjectId(), bson.NewObjectId(), 0, 0, &created_min, &created_max, nil, nil)
			So(err, ShouldBeNil)
			So(custs, ShouldHaveSameTypeAs, []Customer{})
		})
		Convey("with no Id and created, updated dates", func() {
			created_min, _ := time.Parse(time.RFC3339Nano, time.Now().AddDate(-1, 0, 0).Format(time.RFC3339Nano))
			created_max, _ := time.Parse(time.RFC3339Nano, time.Now().Format(time.RFC3339Nano))

			custs, err := CustomersSinceId(bson.NewObjectId(), bson.NewObjectId(), 1, 0, &created_min, &created_max, &created_min, &created_max)
			So(err, ShouldBeNil)
			So(custs, ShouldHaveSameTypeAs, []Customer{})
			for _, cust := range custs {
				So(cust.Password, ShouldEqual, "")
			}
		})
	})
}

func TestGetCustomer(t *testing.T) {
	clearMongo()

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
			for _, cust := range custs {
				So(cust.Password, ShouldEqual, "")
			}
		})
	})
}

func TestCustomer(t *testing.T) {
	clearMongo()

	Convey("Testing Get with no Id", t, func() {
		cust := Customer{}
		err := cust.Get()
		So(err, ShouldNotBeNil)

		err = cust.GetByEmail()
		So(err, ShouldNotBeNil)
	})

	Convey("Generating a Test Shop and testing Customer functions", t, func() {
		if id := InsertTestData(); id != nil {
			shop := Shop{
				Id: *id,
			}
			var customer Customer
			customer.ShopId = shop.Id
			err := customer.Insert("http://www.example.com")
			So(err, ShouldNotBeNil)

			customer.Email = "ninnemana@gmail.com"
			customer.Password = "password"
			err = customer.Insert("http://www.example.com")
			So(err, ShouldNotBeNil)

			customer.FirstName = "Alex"
			err = customer.Insert("http://www.example.com")
			So(err, ShouldNotBeNil)

			customer.LastName = "Ninneman"
			os.Setenv("MONGO_URL", "0.0.0.1")
			err = customer.Insert("http://www.example.com")
			So(err, ShouldNotBeNil)
			os.Setenv("MONGO_URL", "")

			customer.Password = ""
			err = customer.Insert("http://www.example.com")
			So(err, ShouldNotBeNil)

			customer.Password = "password"
			err = customer.Insert("http://www.example.com")
			So(err, ShouldBeNil)
			So(customer.Password, ShouldEqual, "")

			customer.Password = "bad_password"
			err = customer.Login("http://example.com")
			So(err, ShouldNotBeNil)

			customer.Password = "password"
			customer.Email = "alex@ninneman.org"
			err = customer.Login("http://example.com")
			So(err, ShouldNotBeNil)

			customer.Email = "ninnemana@gmail.com"
			os.Setenv("MONGO_URL", "0.0.0.1")
			err = customer.Login("http://example.com")
			So(err, ShouldNotBeNil)
			os.Setenv("MONGO_URL", "")

			customer.Password = "password"
			err = customer.Login("http://example.com")
			So(err, ShouldBeNil)
			So(customer.Password, ShouldEqual, "")

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
			addrId := bson.NewObjectId()
			addr := CustomerAddress{
				Id:           &addrId,
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

			tmpCust := customer
			os.Setenv("MONGO_URL", "0.0.0.1")
			err = customer.Update()
			So(err, ShouldNotBeNil)
			os.Setenv("MONGO_URL", "")
			customer = tmpCust

			err = customer.Update()
			So(err, ShouldBeNil)
			So(customer.Password, ShouldEqual, "")

			os.Setenv("MONGO_URL", "0.0.0.1")
			err = customer.Get()
			So(err, ShouldNotBeNil)
			err = customer.GetByEmail()
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
	clearMongo()
	Convey("Test CustomerCount", t, func() {
		Convey("Bad Connection", func() {
			os.Setenv("MONGO_URL", "0.0.0.1")
			count, err := CustomerCount(bson.NewObjectId())
			So(err, ShouldNotBeNil)
			So(count, ShouldEqual, 0)
			os.Setenv("MONGO_URL", "")
		})
		Convey("Empty Shop ID", func() {
			count, err := CustomerCount("")
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

func TestSearchCustomers(t *testing.T) {
	clearMongo()
	Convey("Bad Connection", t, func() {
		os.Setenv("MONGO_URL", "0.0.0.1")
		custs, err := SearchCustomers("", bson.NewObjectId())
		So(err, ShouldNotBeNil)
		So(custs, ShouldHaveSameTypeAs, []Customer{})
		os.Setenv("MONGO_URL", "")
	})
	Convey("Empty query", t, func() {
		custs, err := SearchCustomers("", bson.NewObjectId())
		So(err, ShouldNotBeNil)
		So(custs, ShouldHaveSameTypeAs, []Customer{})
	})
	Convey("With Query", t, func() {
		custs, err := SearchCustomers("test", bson.NewObjectId())
		t.Log("// TODO: This should not return an error eventually")
		So(err, ShouldNotBeNil)
		So(custs, ShouldHaveSameTypeAs, []Customer{})
	})
}

func clearMongo() {
	sess, err := mgo.DialWithInfo(database.MongoConnectionString())
	if err != nil {
		return
	}
	defer sess.Close()
	sess.DB("CurtCart").C("customer").DropCollection()
	sess.DB("CurtCart").C("order").DropCollection()
	sess.DB("CurtCart").C("shop").DropCollection()

}
