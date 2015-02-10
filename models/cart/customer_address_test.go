package cart

import (
	. "github.com/smartystreets/goconvey/convey"
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

		t.Log(addr.Id)

		addr.City = ""
		err = customer.SaveAddress(addr)
		t.Log(err)
		So(err, ShouldNotBeNil)

		addr.City = "Altoona"
		err = customer.SaveAddress(addr)
		So(err, ShouldBeNil)
	})
}
