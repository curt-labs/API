package applicationGuide

import (
	"github.com/curt-labs/GoAPI/helpers/apicontext"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

var (
	MockedDTX = &apicontext.DataContext{BrandID: 1, WebsiteID: 1, APIKey: "NOT_GENERATED_YET"}
)

func TestAppGuides(t *testing.T) {

	if err := MockedDTX.Mock(); err != nil {
		return
	}
	Convey("Test Create AppGuide", t, func() {
		var err error
		var ag ApplicationGuide

		//create
		ag.FileType = "pdf"
		ag.Url = "test.com"
		ag.Website.ID = 1
		err = ag.Create(MockedDTX)
		So(err, ShouldBeNil)

		//get
		err = ag.Get(MockedDTX)
		So(err, ShouldBeNil)

		//get by site
		ags, err := ag.GetBySite(MockedDTX)

		So(err, ShouldBeNil)
		So(len(ags), ShouldBeGreaterThanOrEqualTo, 1)

		//delete
		err = ag.Delete()
		So(err, ShouldBeNil)

	})
	MockedDTX.DeMock()
}

func BenchmarkGetAppGuide(b *testing.B) {
	if err := MockedDTX.Mock(); err != nil {
		return
	}
	var ag ApplicationGuide
	ag.FileType = "pdf"
	ag.Url = "http://google.com"
	ag.Website.ID = 1

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		ag.Create(MockedDTX)
		b.StartTimer()
		ag.Get(MockedDTX)
		b.StopTimer()
		ag.Delete()
	}
	MockedDTX.DeMock()
}

func BenchmarkGetBySite(b *testing.B) {
	if err := MockedDTX.Mock(); err != nil {
		return
	}
	var ag ApplicationGuide
	ag.FileType = "pdf"
	ag.Url = "http://google.com"
	ag.Website.ID = 1

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		ag.Create(MockedDTX)
		b.StartTimer()
		ag.GetBySite(MockedDTX)
		b.StopTimer()
		ag.Delete()
	}
	MockedDTX.DeMock()
}

func BenchmarkDeleteAppGuide(b *testing.B) {
	if err := MockedDTX.Mock(); err != nil {
		return
	}
	var ag ApplicationGuide
	ag.FileType = "pdf"
	ag.Url = "http://google.com"
	ag.Website.ID = 1

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		ag.Create(MockedDTX)
		b.StartTimer()
		ag.Delete()
	}
	MockedDTX.DeMock()
}
