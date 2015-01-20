package applicationGuide

import (
	"github.com/curt-labs/GoAPI/helpers/apicontext"
	"github.com/curt-labs/GoAPI/helpers/apicontextmock"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestAppGuides(t *testing.T) {
	MockedDTX := &apicontext.DataContext{}
	var err error
	if MockedDTX, err = apicontextmock.Mock(); err != nil {
		return
	}
	Convey("Test Create AppGuide", t, func() {
		var err error
		var ag ApplicationGuide

		//create
		ag.FileType = "pdf"
		ag.Url = "test.com"
		ag.Website.ID = MockedDTX.WebsiteID
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
	_ = apicontextmock.DeMock(MockedDTX)
}

func BenchmarkGetAppGuide(b *testing.B) {
	MockedDTX := &apicontext.DataContext{}
	var err error
	if MockedDTX, err = apicontextmock.Mock(); err != nil {
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
	_ = apicontextmock.DeMock(MockedDTX)
}

func BenchmarkGetBySite(b *testing.B) {
	MockedDTX := &apicontext.DataContext{}
	var err error
	if MockedDTX, err = apicontextmock.Mock(); err != nil {
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
	_ = apicontextmock.DeMock(MockedDTX)
}

func BenchmarkDeleteAppGuide(b *testing.B) {
	MockedDTX := &apicontext.DataContext{}
	var err error
	if MockedDTX, err = apicontextmock.Mock(); err != nil {
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
	_ = apicontextmock.DeMock(MockedDTX)
}
