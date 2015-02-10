package cart

import (
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/mgo.v2/bson"
	"os"
	"testing"
)

func TestValidate(t *testing.T) {
	clearMongo()

	Convey("Testing Validate", t, func() {
		var addr CustomerAddress
		err := addr.Validate()
		So(err, ShouldNotBeNil)

		addr.Address1 = "1 Main St"
		err = addr.Validate()
		So(err, ShouldNotBeNil)

		addr.City = "Hometown"
		err = addr.Validate()
		So(err, ShouldNotBeNil)

		addr.Province = "AB"
		err = addr.Validate()
		So(err, ShouldNotBeNil)

		addr.CountryCode = "CD"
		err = addr.Validate()
		So(err, ShouldNotBeNil)

		addr.Zip = "12345"
		err = addr.Validate()
		So(err, ShouldBeNil)
	})
}

func Test_AddAddress(t *testing.T) {
	clearMongo()

	id := InsertTestData()
	if id == nil {
		return
	}

	var customer Customer
	customer.ShopId = *id
	customer.FirstName = "Alex"
	customer.LastName = "Ninneman"
	customer.Password = "password"
	customer.Email = "ninnemana@gmail.com"
	err := customer.Insert("http://example.com")
	if err != nil {
		t.Log(err)
		return
	}

	Convey("Add Address", t, func() {

		var id bson.ObjectId = customer.Id

		customer.Id = ""
		Convey("with bad customer Id", func() {
			addr := CustomerAddress{
				Company:   "CURT Manufacturing, LLC",
				Name:      "Alex's House",
				FirstName: "Alex",
				LastName:  "Ninneman",
				Phone:     "7153082604",
				Address1:  "1119 Sunset Lane",
			}

			err = customer.AddAddress(addr)
			So(err, ShouldNotBeNil)
		})
		customer.Id = id

		Convey("with bad address1", func() {
			addr := CustomerAddress{
				Company:   "CURT Manufacturing, LLC",
				Name:      "Alex's House",
				FirstName: "Alex",
				LastName:  "Ninneman",
				Phone:     "7153082604",
			}

			err = customer.AddAddress(addr)
			So(err, ShouldNotBeNil)
		})
		Convey("with bad city", func() {
			addr := CustomerAddress{
				Company:   "CURT Manufacturing, LLC",
				Name:      "Alex's House",
				FirstName: "Alex",
				LastName:  "Ninneman",
				Phone:     "7153082604",
				Address1:  "1119 Sunset Lane",
			}

			err = customer.AddAddress(addr)
			So(err, ShouldNotBeNil)
		})
		Convey("with bad province", func() {
			addr := CustomerAddress{
				Company:   "CURT Manufacturing, LLC",
				Name:      "Alex's House",
				FirstName: "Alex",
				LastName:  "Ninneman",
				Phone:     "7153082604",
				Address1:  "1119 Sunset Lane",
				City:      "Altoona",
			}

			err = customer.AddAddress(addr)
			So(err, ShouldNotBeNil)
		})
		Convey("with bad province_code", func() {
			addr := CustomerAddress{
				Company:   "CURT Manufacturing, LLC",
				Name:      "Alex's House",
				FirstName: "Alex",
				LastName:  "Ninneman",
				Phone:     "7153082604",
				Address1:  "1119 Sunset Lane",
				City:      "Altoona",
				Province:  "Wisconsin",
			}

			err = customer.AddAddress(addr)
			So(err, ShouldNotBeNil)
		})
		Convey("with bad country", func() {
			addr := CustomerAddress{
				Company:      "CURT Manufacturing, LLC",
				Name:         "Alex's House",
				FirstName:    "Alex",
				LastName:     "Ninneman",
				Phone:        "7153082604",
				Address1:     "1119 Sunset Lane",
				City:         "Altoona",
				Province:     "Wisconsin",
				ProvinceCode: "WI",
			}

			err = customer.AddAddress(addr)
			So(err, ShouldNotBeNil)
		})
		Convey("with bad country_name", func() {
			addr := CustomerAddress{
				Company:      "CURT Manufacturing, LLC",
				Name:         "Alex's House",
				FirstName:    "Alex",
				LastName:     "Ninneman",
				Phone:        "7153082604",
				Address1:     "1119 Sunset Lane",
				City:         "Altoona",
				Province:     "Wisconsin",
				ProvinceCode: "WI",
				Country:      "United States of America",
			}

			err = customer.AddAddress(addr)
			So(err, ShouldNotBeNil)
		})
		Convey("with bad country_code", func() {
			addr := CustomerAddress{
				Company:      "CURT Manufacturing, LLC",
				Name:         "Alex's House",
				FirstName:    "Alex",
				LastName:     "Ninneman",
				Phone:        "7153082604",
				Address1:     "1119 Sunset Lane",
				City:         "Altoona",
				Province:     "Wisconsin",
				ProvinceCode: "WI",
				Country:      "United States of America",
				CountryName:  "United States",
			}

			err = customer.AddAddress(addr)
			So(err, ShouldNotBeNil)
		})
		Convey("with bad zip", func() {
			addr := CustomerAddress{
				Company:      "CURT Manufacturing, LLC",
				Name:         "Alex's House",
				FirstName:    "Alex",
				LastName:     "Ninneman",
				Phone:        "7153082604",
				Address1:     "1119 Sunset Lane",
				City:         "Altoona",
				Province:     "Wisconsin",
				ProvinceCode: "WI",
				Country:      "United States of America",
				CountryName:  "United States",
				CountryCode:  "USA",
			}

			err = customer.AddAddress(addr)
			So(err, ShouldNotBeNil)
		})
		Convey("with good data and bad connection", func() {
			addr := CustomerAddress{
				Company:      "CURT Manufacturing, LLC",
				Name:         "Alex's House",
				FirstName:    "Alex",
				LastName:     "Ninneman",
				Phone:        "7153082604",
				Address1:     "1119 Sunset Lane",
				City:         "Altoona",
				Province:     "Wisconsin",
				ProvinceCode: "WI",
				Country:      "United States of America",
				CountryName:  "United States",
				CountryCode:  "USA",
				Zip:          "54720",
			}

			os.Setenv("MONGO_URL", "0.0.0.1")
			err = customer.AddAddress(addr)
			So(err, ShouldNotBeNil)
			os.Setenv("MONGO_URL", "")
		})

		Convey("with good data", func() {
			addr := CustomerAddress{
				Company:      "CURT Manufacturing, LLC",
				Name:         "Alex's House",
				FirstName:    "Alex",
				LastName:     "Ninneman",
				Phone:        "7153082604",
				Address1:     "1119 Sunset Lane",
				City:         "Altoona",
				Province:     "Wisconsin",
				ProvinceCode: "WI",
				Country:      "United States of America",
				CountryName:  "United States",
				CountryCode:  "USA",
				Zip:          "54720",
			}

			err = customer.AddAddress(addr)
			So(err, ShouldBeNil)
		})
	})
}

func Test_EditAddress(t *testing.T) {
	clearMongo()

	id := InsertTestData()
	if id == nil {
		return
	}

	var customer Customer
	customer.ShopId = *id
	customer.FirstName = "Alex"
	customer.LastName = "Ninneman"
	customer.Password = "password"
	customer.Email = "ninnemana@gmail.com"
	err := customer.Insert("http://example.com")
	if err != nil {
		t.Log(err)
		return
	}

	addr := CustomerAddress{
		Company:      "CURT Manufacturing, LLC",
		Name:         "Alex's House",
		FirstName:    "Alex",
		LastName:     "Ninneman",
		Phone:        "7153082604",
		Address1:     "1119 Sunset Lane",
		City:         "Altoona",
		Province:     "Wisconsin",
		ProvinceCode: "WI",
		Country:      "United States of America",
		CountryName:  "United States",
		CountryCode:  "USA",
		Zip:          "54720",
	}

	Convey("Edit Address", t, func() {

		err = customer.AddAddress(addr)
		So(err, ShouldBeNil)

		err = customer.Get()
		So(err, ShouldBeNil)
		So(len(customer.Addresses), ShouldBeGreaterThan, 0)

		addr = customer.Addresses[0]

		addr.City = ""
		err = customer.SaveAddress(addr)
		So(err, ShouldNotBeNil)

		os.Setenv("MONGO_URL", "0.0.0.1")
		addr.City = "Altoona"
		err = customer.SaveAddress(addr)
		So(err, ShouldNotBeNil)
		os.Setenv("MONGO_URL", "")

		id := customer.Id
		customer.Id = ""
		addr.City = "Altoona"
		err = customer.SaveAddress(addr)
		So(err, ShouldNotBeNil)

		customer.Id = id
		err = customer.SaveAddress(addr)
		So(err, ShouldBeNil)
	})
}

func Test_deepEqual(t *testing.T) {
	Convey("Augmenting properties to run deep equal", t, func() {
		var a1 *CustomerAddress
		var a2 *CustomerAddress
		res := a1.deepEqual(a2)
		So(res, ShouldBeTrue)

		a1 = &CustomerAddress{}
		res = a1.deepEqual(a2)
		So(res, ShouldBeFalse)

		a1.Address1 = "100 Main"
		a2 = &CustomerAddress{}
		res = a1.deepEqual(a2)
		So(res, ShouldBeFalse)

		a2.Address1 = a1.Address1
		a1.Address2 = "Apt 1"
		res = a1.deepEqual(a2)
		So(res, ShouldBeFalse)

		a2.Address2 = a1.Address2
		a1.City = "Brown"
		res = a1.deepEqual(a2)
		So(res, ShouldBeFalse)

		a2.City = a1.City
		a2.Company = "Some Company"
		res = a1.deepEqual(a2)
		So(res, ShouldBeFalse)

		a2.Company = a1.Company
		a1.Name = "Some Name"
		res = a1.deepEqual(a2)
		So(res, ShouldBeFalse)

		a2.Name = a1.Name
		a1.FirstName = "Joe"
		res = a1.deepEqual(a2)
		So(res, ShouldBeFalse)

		a2.FirstName = a1.FirstName
		a1.LastName = "Smith"
		res = a1.deepEqual(a2)
		So(res, ShouldBeFalse)

		a2.LastName = a1.LastName
		a1.Phone = "555-555-5555"
		res = a1.deepEqual(a2)
		So(res, ShouldBeFalse)

		a2.Phone = a1.Phone
		a1.Province = "Wiconsin"
		res = a1.deepEqual(a2)
		So(res, ShouldBeFalse)

		a2.Province = a1.Province
		a1.ProvinceCode = "WI"
		res = a1.deepEqual(a2)
		So(res, ShouldBeFalse)

		a2.ProvinceCode = a1.ProvinceCode
		a1.Country = "United States of America"
		res = a1.deepEqual(a2)
		So(res, ShouldBeFalse)

		a2.Country = a1.Country
		a1.CountryCode = "USA"
		res = a1.deepEqual(a2)
		So(res, ShouldBeFalse)

		a2.CountryCode = a1.CountryCode
		a1.CountryName = "United States"
		res = a1.deepEqual(a2)
		So(res, ShouldBeFalse)

		a2.CountryName = a1.CountryName
		a1.Zip = "12345"
		res = a1.deepEqual(a2)
		So(res, ShouldBeFalse)

		a2.Zip = a1.Zip
		res = a1.deepEqual(a2)
		So(res, ShouldBeTrue)
	})
}
