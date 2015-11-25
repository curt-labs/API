package showcase

import (
	"github.com/curt-labs/API/helpers/apicontextmock"
	. "github.com/smartystreets/goconvey/convey"

	"net/url"
	"testing"
)

func TestShowcases(t *testing.T) {
	MockedDTX, err := apicontextmock.Mock()
	if err != nil {
		return
	}
	var s Showcase
	s.BrandID = MockedDTX.BrandID
	Convey("Testing Create Showcase", t, func() {
		u, _ := url.Parse("www.test.com")
		i := Image{
			Path: u,
		}
		s.Text = "Test Content"
		s.Images = append(s.Images, i)
		err = s.Create()
		So(err, ShouldBeNil)
	})
	Convey("Update", t, func() {
		s.Text = "New Content"
		s.Active = true
		s.Approved = true
		err = s.Update()
		So(err, ShouldBeNil)
	})

	Convey("Get showcase", t, func() {
		err = s.Get(MockedDTX)
		So(err, ShouldBeNil)
	})
	Convey("GetAll - No paging", t, func() {
		shows, err := GetAllShowcases(0, 1, false, MockedDTX)
		So(err, ShouldBeNil)
		So(len(shows), ShouldBeGreaterThan, 0)
	})

	Convey("GetAll - Paged", t, func() {
		shows, err := GetAllShowcases(0, 1, false, MockedDTX)
		So(err, ShouldBeNil)
		So(len(shows), ShouldBeGreaterThan, 0)
	})

	Convey("GetAll - randomized", t, func() {
		shows, err := GetAllShowcases(0, 1, true, MockedDTX)
		So(err, ShouldBeNil)
		So(len(shows), ShouldBeGreaterThan, 0)

	})
	Convey("Delete", t, func() {
		err = s.Delete()
		So(err, ShouldBeNil)
	})

	_ = apicontextmock.DeMock(MockedDTX)
}

func BenchmarkGetAllShowcases(b *testing.B) {
	MockedDTX, err := apicontextmock.Mock()
	if err != nil {
		return
	}
	for i := 0; i < b.N; i++ {
		GetAllShowcases(0, 1, false, MockedDTX)
	}
	_ = apicontextmock.DeMock(MockedDTX)
}

func BenchmarkGetShowcase(b *testing.B) {
	MockedDTX, err := apicontextmock.Mock()
	if err != nil {
		return
	}
	show := setupDummyShowcases()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		show.Create()
		b.StartTimer()
		show.Get(MockedDTX)
		b.StopTimer()
		show.Delete()
	}
	_ = apicontextmock.DeMock(MockedDTX)
}

func BenchmarkCreateShowcases(b *testing.B) {
	MockedDTX, err := apicontextmock.Mock()
	if err != nil {
		return
	}
	show := setupDummyShowcases()
	for i := 0; i < b.N; i++ {
		b.StartTimer()
		show.Create()
		b.StopTimer()
		show.Delete()
	}
	_ = apicontextmock.DeMock(MockedDTX)
}

func BenchmarkUpdateShowcases(b *testing.B) {
	MockedDTX, err := apicontextmock.Mock()
	if err != nil {
		return
	}
	show := setupDummyShowcases()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		show.Create()
		b.StartTimer()
		show.Text = "This is a good test."
		show.Update()
		b.StopTimer()
		show.Delete()
	}
	_ = apicontextmock.DeMock(MockedDTX)
}

func BenchmarkDeleteShowcases(b *testing.B) {
	MockedDTX, err := apicontextmock.Mock()
	if err != nil {
		return
	}
	show := setupDummyShowcases()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		show.Create()
		b.StartTimer()
		show.Delete()
	}
	_ = apicontextmock.DeMock(MockedDTX)
}

func setupDummyShowcases() *Showcase {
	return &Showcase{
		Rating:    5,
		Title:     "Test Test",
		Text:      "This is a test.",
		Approved:  true,
		Active:    true,
		FirstName: "TESTER",
		LastName:  "TESTER",
		Location:  "Testville, Oklahoma",
		BrandID:   1,
	}
}
