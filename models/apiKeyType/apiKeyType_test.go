package apiKeyType

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestTechSupport(t *testing.T) {
	Convey("Test Create AppGuide", t, func() {
		var err error
		var akt ApiKeyType

		//create
		akt.Type = "testType"

		err = akt.Create()
		So(err, ShouldBeNil)

		//get
		err = akt.Get()
		So(err, ShouldBeNil)

		as, err := GetAllApiKeyTypes()
		So(err, ShouldBeNil)
		So(len(as), ShouldBeGreaterThan, 0)

		//delete
		err = akt.Delete()
		So(err, ShouldBeNil)

	})

}
