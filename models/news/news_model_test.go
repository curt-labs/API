package news_model

import (
	"github.com/curt-labs/GoAPI/helpers/pagination"
	. "github.com/smartystreets/goconvey/convey"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func TestNews(t *testing.T) {
	Convey("Testing Gets", t, func() {
		Convey("Testing GetAll()", func() {
			var fs Newses
			var err error
			fs, err = GetAll()
			So(len(fs), ShouldBeGreaterThan, 0)
			So(err, ShouldBeNil)
		})
		Convey("Gets News and a random News, page, resultsperpage", func() {
			fs, err := GetAll()
			So(err, ShouldBeNil)
			if len(fs) > 0 {
				x := rand.Intn(len(fs))
				f := fs[x]
				page := "2"
				resultsperpage := strconv.Itoa(len(fs) / 2)
				Convey("Testing Get()", func() {
					var err error
					err = f.Get()
					So(err, ShouldBeNil)
					So(f, ShouldNotBeNil)
					So(f.Title, ShouldHaveSameTypeAs, "")
					So(f.Lead, ShouldHaveSameTypeAs, "")
					So(f.Content, ShouldHaveSameTypeAs, "")
					aTime := time.Now()
					So(f.PublishStart, ShouldHaveSameTypeAs, aTime)
					So(f.PublishEnd, ShouldHaveSameTypeAs, aTime)
				})

				Convey("Testing GetTitles()", func() {
					var l pagination.Objects
					var err error
					l, err = GetTitles(page, resultsperpage)
					num, err := strconv.Atoi(resultsperpage)
					So(len(l.Objects), ShouldEqual, num)
					So(err, ShouldBeNil)
				})
				Convey("Testing GetLeads()", func() {
					var l pagination.Objects
					var err error
					l, err = GetLeads(page, resultsperpage)
					num, err := strconv.Atoi(resultsperpage)
					So(len(l.Objects), ShouldEqual, num)
					So(err, ShouldBeNil)
				})
				Convey("Testing Search()", func() {
					var err error
					var l pagination.Objects
					l, err = Search("", "", f.Content, "", "", "", "", "1", "3")
					So(len(l.Objects), ShouldBeGreaterThan, 0)
					So(err, ShouldBeNil)
				})
			}
		})
	})
	Convey("Testing C_UD", t, func() {
		Convey("Testing Create/Delete", func() {
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

			Convey("Testing update", func() {
				n.Lead = "Pickles"
				err = n.Update()
				err = n.Get()
				So(n.Lead, ShouldEqual, "Pickles")
				So(err, ShouldBeNil)

				Convey("Testing Delete", func() {
					var err error
					err = n.Delete()
					So(err, ShouldBeNil)
				})
			})
		})
	})

}
