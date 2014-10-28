package applicationGuide

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestTechSupport(t *testing.T) {
	Convey("Test Create AppGuide", t, func() {
		var err error
		var ag ApplicationGuide

		//create
		ag.FileType = "pdf"
		ag.Url = "test.com"
		ag.Website.ID = 1
		err = ag.Create()
		So(err, ShouldBeNil)

		//get
		err = ag.Get()
		So(err, ShouldBeNil)

		//get by site
		ags, err := ag.GetBySite()

		So(err, ShouldBeNil)
		So(len(ags), ShouldBeGreaterThanOrEqualTo, 1)

		//delete
		err = ag.Delete()
		So(err, ShouldBeNil)

	})
}
