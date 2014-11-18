package news_model

import (
	"database/sql"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestNews(t *testing.T) {
	var n News
	var err error
	n.Title = "test news"

	Convey("Test Create", t, func() {
		err = n.Create()
		So(err, ShouldBeNil)
	})
	Convey("Test Update", t, func() {
		n.Title = "Different Title"
		err = n.Update()
		So(err, ShouldBeNil)
	})
	Convey("Test Get", t, func() {
		err = n.Get()
		So(err, ShouldBeNil)
		So(n.Title, ShouldEqual, "Different Title")

		obj, err := Search(n.Title, "", "", "", "", "", "", "1", "1")
		So(len(obj.Objects), ShouldEqual, 1)
		So(err, ShouldBeNil)

		ns, err := GetAll()
		if err != sql.ErrNoRows {
			So(err, ShouldBeNil)
			So(len(ns), ShouldBeGreaterThan, 0)
		}
		ts, err := GetTitles("1", "1")
		if err != sql.ErrNoRows {
			So(err, ShouldBeNil)
			So(len(ts.Objects), ShouldBeGreaterThan, 0)
		}
		ls, err := GetLeads("1", "1")
		if err != sql.ErrNoRows {
			So(err, ShouldBeNil)
			So(len(ls.Objects), ShouldBeGreaterThan, 0)
		}

	})
	Convey("Test Delete", t, func() {
		err = n.Delete()
		So(err, ShouldBeNil)
	})

}

func BenchmarkGetAllNews(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetAll()
	}
}

func BenchmarkGetTitles(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetTitles("1", "1")
	}
}

func BenchmarkGetLeads(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetLeads("1", "1")
	}
}

func BenchmarkSearch(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Search("Title", "", "", "", "", "", "", "1", "1")
	}
}

func BenchmarkGetNews(b *testing.B) {
	n := setupDummyNews()
	n.Create()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		n.Get()
	}
	b.StopTimer()
	n.Delete()
}

func BenchmarkCreateNews(b *testing.B) {
	n := setupDummyNews()
	for i := 0; i < b.N; i++ {
		n.Create()
		b.StopTimer()
		n.Delete()
		b.StartTimer()
	}
}

func BenchmarkUpdateNews(b *testing.B) {
	n := setupDummyNews()
	n.Create()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		n.Title = "TEST TIME"
		n.Content = "This is a awesome test."
		n.Update()
	}
	b.StopTimer()
	n.Delete()
}

func BenchmarkDeleteNews(b *testing.B) {
	n := setupDummyNews()
	n.Create()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		n.Delete()
	}
	b.StopTimer()
	n.Delete()
}

func setupDummyNews() *News {
	return &News{
		Title:   "TESTER",
		Content: "TEST TEST TEST",
	}
}
