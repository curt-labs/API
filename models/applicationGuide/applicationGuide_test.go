package applicationGuide

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestAppGuides(t *testing.T) {
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

func BenchmarkGetAppGuide(b *testing.B) {
	var ag ApplicationGuide
	ag.FileType = "pdf"
	ag.Url = "http://google.com"
	ag.Website.ID = 1

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		ag.Create()
		b.StartTimer()
		ag.Get()
		b.StopTimer()
		ag.Delete()
	}
}

func BenchmarkGetBySite(b *testing.B) {
	var ag ApplicationGuide
	ag.FileType = "pdf"
	ag.Url = "http://google.com"
	ag.Website.ID = 1

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		ag.Create()
		b.StartTimer()
		ag.GetBySite()
		b.StopTimer()
		ag.Delete()
	}
}

func BenchmarkDeleteAppGuide(b *testing.B) {
	var ag ApplicationGuide
	ag.FileType = "pdf"
	ag.Url = "http://google.com"
	ag.Website.ID = 1

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		ag.Create()
		b.StartTimer()
		ag.Delete()
	}
}
