package news_model

import (
	"github.com/curt-labs/goacesapi/helpers/pagination"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func TestNews(t *testing.T) {
	Convey("Testing Gets", t, func() {
		Convey("Testing Get()", func() {
			var f News
			f.ID = 1
			f.Get()
			So(f, ShouldNotBeNil)
			So(f.Title, ShouldEqual, "Growth & Expansion Continues at CURT Manufacturing ")
			So(f.Lead, ShouldEqual, "Company Posts Record Sales in First-Half of 2011 - Set to Open New DC")
			So(f.Content, ShouldNotBeNil)

			strTime := f.PublishStart.String()
			So(strTime, ShouldContainSubstring, "2011-09-29 09:22:00")
			strTime2 := f.PublishEnd.String()
			So(strTime2, ShouldContainSubstring, "0001-01-01 00:00:00")
		})
		Convey("Testing GetAll()", func() {
			var fs Newses
			var err error
			fs, err = GetAll()
			So(len(fs), ShouldBeGreaterThan, 0)
			So(err, ShouldBeNil)
		})
		Convey("Testing GetTitles()", func() {
			var l pagination.Objects
			var err error
			l, err = GetTitles("1", "4")
			So(len(l.Objects), ShouldEqual, 4)
			So(err, ShouldBeNil)
		})
		Convey("Testing GetLeads()", func() {
			var l pagination.Objects
			var err error
			l, err = GetLeads("1", "2")
			So(len(l.Objects), ShouldEqual, 2)
			So(err, ShouldBeNil)
		})
		Convey("Testing Search()", func() {
			var err error
			var n News
			var l pagination.Objects
			n.Lead = "Curt"
			l, err = Search("", "", n.Lead, "", "", "", "", "1", "3")
			So(len(l.Objects), ShouldEqual, 3)
			So(err, ShouldBeNil)
		})
	})
	Convey("Testing CUD", t, func() {
		Convey("Testing Create", func() {
			var n News
			var err error
			n.Title = "test"
			n.Lead = "testlead"
			n.Content = "content!"
			n.PublishStart, err = time.Parse(timeFormat, "2011-09-29 09:22:00")
			n.PublishEnd, err = time.Parse(timeFormat, "2011-09-29 09:22:00")
			So(err, ShouldBeNil)
			err = n.Create()
			So(err, ShouldBeNil)
		})
		Convey("Testing update", func() {
			var n News
			var err error
			n.ID = 13
			n.Lead = "Pickles"
			err = n.Update()
			err = n.Get()
			So(n.Lead, ShouldEqual, "Pickles")
			So(err, ShouldBeNil)
		})
		Convey("Testing Delete", func() {
			n := News{ID: 14}
			var err error
			err = n.Delete()
			So(err, ShouldNotBeNil)
		})
	})
}
