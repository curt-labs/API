package blog_model

import (
	"github.com/curt-labs/GoAPI/helpers/apicontext"
	"github.com/curt-labs/GoAPI/helpers/apicontextmock"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

// var (
// 	MockedDTX = &apicontext.DataContext{BrandID: 1, WebsiteID: 1, APIKey: "NOT_GENERATED_YET"}
// )

func TestGetBlogs(t *testing.T) {
	var f Blog
	var cs BlogCategories
	var bc BlogCategory
	var c Category
	var err error
	MockedDTX := &apicontext.DataContext{}
	if MockedDTX, err = apicontextmock.Mock(); err != nil {
		return
	}

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

		err = c.Create(MockedDTX)
		So(err, ShouldBeNil)

	})
	Convey("Testing Gets", t, func() {
		err = c.Get(MockedDTX)
		So(c, ShouldNotBeNil)
		So(err, ShouldBeNil)
		So(c.Name, ShouldEqual, "testTitle")
		So(c.Slug, ShouldEqual, "testSlug")
		So(c.Active, ShouldBeTrue)

		err = f.Get(MockedDTX)
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
		err = f.Update(MockedDTX)
		So(err, ShouldBeNil)
		So(err, ShouldBeNil)
		So(f, ShouldNotBeNil)
		So(f.Title, ShouldEqual, "testTitle222")
		So(f.Slug, ShouldEqual, "testSlug222")

	})
	Convey("Testing Gets", t, func() {

		var bs Blogs
		var err error
		bs, err = GetAll(MockedDTX)
		So(bs, ShouldHaveSameTypeAs, Blogs{})
		So(err, ShouldBeNil)
		So(len(bs), ShouldNotBeNil)

	})

	Convey("Testing GetAllCategories()", t, func() {
		qs, err := GetAllCategories(MockedDTX)
		So(qs, ShouldHaveSameTypeAs, Categories{})
		So(err, ShouldBeNil)
	})
	Convey("Testing Search()", t, func() {
		as, err := Search("test", "", "", "", "", "", "", "", "", "", "", "1", "0", MockedDTX)
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

		err = c.Delete(MockedDTX)
		So(err, ShouldBeNil)

	})
	_ = apicontextmock.DeMock(MockedDTX)
}

func BenchmarkGetAllBlogs(b *testing.B) {
	MockedDTX := &apicontext.DataContext{}
	var err error
	if MockedDTX, err = apicontextmock.Mock(); err != nil {
		return
	}
	for i := 0; i < b.N; i++ {
		GetAll(MockedDTX)
	}
	_ = apicontextmock.DeMock(MockedDTX)
}

func BenchmarkGetAllBlogCategories(b *testing.B) {
	MockedDTX := &apicontext.DataContext{}
	var err error
	if MockedDTX, err = apicontextmock.Mock(); err != nil {
		return
	}
	for i := 0; i < b.N; i++ {
		GetAllCategories(MockedDTX)
	}
	_ = apicontextmock.DeMock(MockedDTX)
}

/**
TODO: I think Redis is making it so where these don't run right
func BenchmarkGetBlog(b *testing.B) {
	blog := Blog{
		UserID:          1,
		MetaTitle:       "Test",
		MetaDescription: "Test",
		Keywords:        "Test",
		Active:          true,
	}
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		blog.Create()
		b.StartTimer()
		blog.Get()
		b.StopTimer()
		blog.Delete()
	}
}
func BenchmarkGetCategory(b *testing.B) {
	cat := Category{
		Name:   "TESTER",
		Slug:   "TESTER",
		Active: false,
	}
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		cat.Create()
		b.StartTimer()
		cat.Get()
		b.StopTimer()
		cat.Delete()
	}
}
**/

func BenchmarkUpdateBlog(b *testing.B) {
	MockedDTX := &apicontext.DataContext{}
	var err error
	if MockedDTX, err = apicontextmock.Mock(); err != nil {
		return
	}
	blog := setupDummyBlog()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		blog.Create()
		b.StartTimer()
		blog.Text = "Blog post magic. Whoop! Whoop!"
		blog.Update(MockedDTX)
		b.StopTimer()
		blog.Delete()
	}
}

func BenchmarkDeleteBlog(b *testing.B) {
	MockedDTX := &apicontext.DataContext{}
	var err error
	if MockedDTX, err = apicontextmock.Mock(); err != nil {
		return
	}
	blog := setupDummyBlog()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		blog.Create()
		b.StartTimer()
		blog.Delete()
	}
	_ = apicontextmock.DeMock(MockedDTX)
}

func BenchmarkDeleteCategory(b *testing.B) {
	MockedDTX := &apicontext.DataContext{}
	var err error
	if MockedDTX, err = apicontextmock.Mock(); err != nil {
		return
	}
	cat := Category{
		Name: "TESTER",
		Slug: "TESTER",
	}
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		cat.Create(MockedDTX)
		b.StartTimer()
		cat.Delete(MockedDTX)
	}
	_ = apicontextmock.DeMock(MockedDTX)
}

func setupDummyBlog() *Blog {
	return &Blog{
		UserID:          1,
		MetaTitle:       "Test",
		MetaDescription: "Test",
		Keywords:        "Test",
		Active:          true,
	}
}
