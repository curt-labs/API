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
