package news_model

import (
	"database/sql"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestNews(t *testing.T) {

	Convey("Test CRUD", t, func() {
		var n News
		var err error
		n.Title = "test news"
		err = n.Create()
		So(err, ShouldBeNil)

		n.Title = "Different Title"
		err = n.Update()
		So(err, ShouldBeNil)

		err = n.Get()
		So(err, ShouldBeNil)
		So(n.Title, ShouldEqual, "Different Title")

		obj, err := Search(n.Title, "", "", "", "", "", "", "1", "1")
		So(len(obj.Objects), ShouldEqual, 1)
		So(err, ShouldBeNil)

		err = n.Delete()
		So(err, ShouldBeNil)
	})
	Convey("Test Get-Alls", t, func() {
		var err error
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
}
