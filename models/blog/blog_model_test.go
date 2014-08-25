package blog_model

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func TestGetBlogs(t *testing.T) {
	Convey("Testing Gets", t, func() {
		Convey("Testing Get()", func() {
			var b Blog
			b.ID = 1
			b.Get()
			So(b, ShouldNotBeNil)
			So(b.Title, ShouldEqual, "Cyclocross for Hunger")
			So(b.Slug, ShouldEqual, "cyclocross_for_hunger")
		})
		Convey("Testing Get()", func() {
			var b Blog
			var err error
			b.ID = 2
			err = b.Get()
			So(b, ShouldNotBeNil)
			var t time.Time
			So(b.PublishedDate, ShouldHaveSameTypeAs, t)
			So(err, ShouldBeNil)
		})
		Convey("Testing GetAll()", func() {
			var b Blogs
			var err error
			b, err = GetAll()
			So(b, ShouldNotBeNil)
			So(err, ShouldBeNil)
			So(len(b), ShouldNotBeNil)
		})
		Convey("Testing GetAllCategories()", func() {
			qs, err := GetAllCategories()
			So(qs, ShouldNotBeNil)
			So(err, ShouldBeNil)
		})
		Convey("Testing Search()", func() {
			as, err := Search("test", "", "", "", "", "", "", "", "", "", "", "1", "0")
			So(as, ShouldNotBeNil)
			So(err, ShouldBeNil)
			So(as.Pagination.Page, ShouldEqual, 1)
			So(as.Pagination.ReturnedCount, ShouldNotBeNil)
			So(as.Pagination.PerPage, ShouldNotBeNil)
			So(as.Pagination.PerPage, ShouldEqual, len(as.Objects))
		})

	})
	Convey("Testing CUD", t, func() {
		Convey("Testing Create()", func() {
			var f Blog
			var cs BlogCategories
			var c BlogCategory
			var err error
			f.Title = "testTitle"
			f.Slug = "testSlug"
			f.Text = "test"
			f.PublishedDate, err = time.Parse(timeFormat, "2004-03-03 9:15:00")
			f.UserID = 1
			f.MetaTitle = "test"
			f.MetaDescription = "test"
			f.Keywords = "test"
			f.Active = true
			c.Category.Active = true
			c.Category.Name = "testCat"
			c.Category.Slug = "catSlug"
			cs = append(cs, c)
			f.BlogCategories = cs

			f.Create()
			f.Get()
			So(f, ShouldNotBeNil)
			So(err, ShouldBeNil)
			So(f.Title, ShouldEqual, "testTitle")
			So(f.Slug, ShouldEqual, "testSlug")
			var t time.Time
			So(f.PublishedDate, ShouldHaveSameTypeAs, t)
			f.Delete()
		})
		Convey("Testing Update()", func() {
			var f Blog
			var err error
			f.ID = 17
			f.Title = "testTitle222"
			f.Slug = "testSlug222"
			f.PublishedDate, err = time.Parse(timeFormat, "2004-03-03 09:15:00")

			c := make(chan int)
			go func() {
				f.Update()
				c <- 1
			}()
			<-c
			f.Get()
			So(err, ShouldBeNil)
			So(f, ShouldNotBeNil)
			So(f.Title, ShouldEqual, "testTitle222")
			So(f.Slug, ShouldEqual, "testSlug222")

		})
		Convey("Testing Delete()", func() {
			var f Blog
			f.ID = 15
			c := make(chan int)
			go func() {
				f.Delete()
				c <- 1
			}()
			<-c
			go f.Get()
			So(f.Title, ShouldBeBlank)
			So(f.Text, ShouldBeBlank)
			So(f.CreatedDate, ShouldBeZeroValue)
			So(f.LastModified, ShouldBeZeroValue)

		})
		Convey("Testing CreateCategory()/DeleteCategory()", func() {
			var c Category
			var err error
			c.Name = "testTitle"
			c.Slug = "testSlug"
			c.Active = true
			c.Create()
			c.Get()
			So(c, ShouldNotBeNil)
			So(err, ShouldBeNil)
			So(c.Name, ShouldEqual, "testTitle")
			So(c.Slug, ShouldEqual, "testSlug")
			So(c.Active, ShouldBeTrue)
			err = c.Delete()
			So(err, ShouldBeNil)

		})

	})

}
