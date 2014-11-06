package blog_model

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func TestGetBlogs(t *testing.T) {
	var f Blog
	var cs BlogCategories
	var bc BlogCategory
	var c Category
	var err error
	Convey("Testing Create", t, func() {
		f.Title = "testTitle"
		f.Slug = "testSlug"
		f.Text = "test"
		f.PublishedDate, err = time.Parse(timeFormat, "2004-03-03 9:15:00")
		f.UserID = 1
		f.MetaTitle = "test"
		f.MetaDescription = "test"
		f.Keywords = "test"
		f.Active = true
		c.Name = "testTitle"
		c.Slug = "testSlug"
		c.Active = true
		bc.Category = c
		cs = append(cs, bc)
		f.BlogCategories = cs

		err = f.Create()
		So(err, ShouldBeNil)

		err = c.Create()
		So(err, ShouldBeNil)

	})
	Convey("Testing Gets", t, func() {
		err = c.Get()
		So(c, ShouldNotBeNil)
		So(err, ShouldBeNil)
		So(c.Name, ShouldEqual, "testTitle")
		So(c.Slug, ShouldEqual, "testSlug")
		So(c.Active, ShouldBeTrue)

		err = f.Get()
		So(f, ShouldNotBeNil)
		So(err, ShouldBeNil)
		So(f.Title, ShouldEqual, "testTitle")
		So(f.Slug, ShouldEqual, "testSlug")
		var t time.Time
		So(f.PublishedDate, ShouldHaveSameTypeAs, t)

	})
	Convey("Testing Update", t, func() {
		f.Title = "testTitle222"
		f.Slug = "testSlug222"
		f.PublishedDate, err = time.Parse(timeFormat, "2004-03-03 09:15:00")
		err = f.Update()
		So(err, ShouldBeNil)
		So(err, ShouldBeNil)
		So(f, ShouldNotBeNil)
		So(f.Title, ShouldEqual, "testTitle222")
		So(f.Slug, ShouldEqual, "testSlug222")

	})
	Convey("Testing Gets", t, func() {

		var bs Blogs
		var err error
		bs, err = GetAll()
		So(bs, ShouldHaveSameTypeAs, Blogs{})
		So(err, ShouldBeNil)
		So(len(bs), ShouldNotBeNil)

	})

	Convey("Testing GetAllCategories()", t, func() {
		qs, err := GetAllCategories()
		So(qs, ShouldHaveSameTypeAs, Categories{})
		So(err, ShouldBeNil)
	})
	Convey("Testing Search()", t, func() {
		as, err := Search("test", "", "", "", "", "", "", "", "", "", "", "1", "0")
		So(as, ShouldNotBeNil)
		So(err, ShouldBeNil)
		So(as.Pagination.Page, ShouldEqual, 1)
		So(as.Pagination.ReturnedCount, ShouldNotBeNil)
		So(as.Pagination.PerPage, ShouldNotBeNil)
		So(as.Pagination.PerPage, ShouldEqual, len(as.Objects))
	})

	Convey("Testing Delete", t, func() {
		err = f.Delete()
		So(err, ShouldBeNil)

		err = c.Delete()
		So(err, ShouldBeNil)

	})

}
