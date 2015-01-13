package news_model

import (
	"database/sql"
	"github.com/curt-labs/GoAPI/helpers/apicontext"
	"github.com/curt-labs/GoAPI/helpers/apicontextmock"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestNews(t *testing.T) {
	var n News
	var err error
	n.Title = "test news"

	MockedDTX := &apicontext.DataContext{}
	if MockedDTX, err = apicontextmock.Mock(); err != nil {
		return
	}

	Convey("Test Create", t, func() {
		err = n.Create(MockedDTX)
		So(err, ShouldBeNil)
	})
	Convey("Test Update", t, func() {
		n.Title = "Different Title"
		err = n.Update(MockedDTX)
		So(err, ShouldBeNil)
	})
	Convey("Test Get", t, func() {
		err = n.Get(MockedDTX)
		So(err, ShouldBeNil)
		So(n.Title, ShouldEqual, "Different Title")

		obj, err := Search(n.Title, "", "", "", "", "", "", "1", "1", MockedDTX)
		So(len(obj.Objects), ShouldEqual, 1)
		So(err, ShouldBeNil)

		ns, err := GetAll(MockedDTX)
		if err != sql.ErrNoRows {
			So(err, ShouldBeNil)
			So(len(ns), ShouldBeGreaterThan, 0)
		}
		ts, err := GetTitles("1", "1", MockedDTX)
		if err != sql.ErrNoRows {
			So(err, ShouldBeNil)
			So(len(ts.Objects), ShouldBeGreaterThan, 0)
		}
		ls, err := GetLeads("1", "1", MockedDTX)
		if err != sql.ErrNoRows {
			So(err, ShouldBeNil)
			So(len(ls.Objects), ShouldBeGreaterThan, 0)
		}

	})
	Convey("Test Delete", t, func() {
		err = n.Delete(MockedDTX)
		So(err, ShouldBeNil)
	})

}

func BenchmarkGetAllNews(b *testing.B) {
	var err error
	MockedDTX := &apicontext.DataContext{}
	if MockedDTX, err = apicontextmock.Mock(); err != nil {
		return
	}
	for i := 0; i < b.N; i++ {
		GetAll(MockedDTX)
	}
}

func BenchmarkGetTitles(b *testing.B) {
	var err error
	MockedDTX := &apicontext.DataContext{}
	if MockedDTX, err = apicontextmock.Mock(); err != nil {
		return
	}
	for i := 0; i < b.N; i++ {
		GetTitles("1", "1", MockedDTX)
	}
}

func BenchmarkGetLeads(b *testing.B) {
	var err error
	MockedDTX := &apicontext.DataContext{}
	if MockedDTX, err = apicontextmock.Mock(); err != nil {
		return
	}
	for i := 0; i < b.N; i++ {
		GetLeads("1", "1", MockedDTX)
	}
}

func BenchmarkSearch(b *testing.B) {
	var err error
	MockedDTX := &apicontext.DataContext{}
	if MockedDTX, err = apicontextmock.Mock(); err != nil {
		return
	}
	for i := 0; i < b.N; i++ {
		Search("Title", "", "", "", "", "", "", "1", "1", MockedDTX)
	}
}

func BenchmarkGetNews(b *testing.B) {
	var err error
	MockedDTX := &apicontext.DataContext{}
	if MockedDTX, err = apicontextmock.Mock(); err != nil {
		return
	}
	n := setupDummyNews()
	n.Create(MockedDTX)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		n.Get(MockedDTX)
	}
	b.StopTimer()
	n.Delete(MockedDTX)
}

func BenchmarkCreateNews(b *testing.B) {
	var err error
	MockedDTX := &apicontext.DataContext{}
	if MockedDTX, err = apicontextmock.Mock(); err != nil {
		return
	}
	n := setupDummyNews()
	for i := 0; i < b.N; i++ {
		n.Create(MockedDTX)
		b.StopTimer()
		n.Delete(MockedDTX)
		b.StartTimer()
	}
}

func BenchmarkUpdateNews(b *testing.B) {
	var err error
	MockedDTX := &apicontext.DataContext{}
	if MockedDTX, err = apicontextmock.Mock(); err != nil {
		return
	}
	n := setupDummyNews()
	n.Create(MockedDTX)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		n.Title = "TEST TIME"
		n.Content = "This is a awesome test."
		n.Update(MockedDTX)
	}
	b.StopTimer()
	n.Delete(MockedDTX)
}

func BenchmarkDeleteNews(b *testing.B) {
	var err error
	MockedDTX := &apicontext.DataContext{}
	if MockedDTX, err = apicontextmock.Mock(); err != nil {
		return
	}
	n := setupDummyNews()
	n.Create(MockedDTX)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		n.Delete(MockedDTX)
	}
	b.StopTimer()
	n.Delete(MockedDTX)
}

func setupDummyNews() *News {
	return &News{
		Title:   "TESTER",
		Content: "TEST TEST TEST",
	}
}
