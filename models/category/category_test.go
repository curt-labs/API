package category

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestTopTierCategories(t *testing.T) {
	Convey("Test TopTierCategories()", t, func() {
		cats, err := TopTierCategories()
		So(cats, ShouldNotBeNil)
		So(err, ShouldBeNil)
		So(cats, ShouldNotBeEmpty)
	})
}

func TestGetByTitle(t *testing.T) {

	Convey("Test GetByTitle", t, func() {
		Convey("with ``", func() {
			cat, err := GetByTitle("")
			So(cat, ShouldNotBeNil)
			So(cat.CategoryId, ShouldEqual, 0)
			So(err, ShouldBeNil)
		})
		Convey("with `test`", func() {
			cat, err := GetByTitle("test")
			So(cat, ShouldNotBeNil)
			So(cat.CategoryId, ShouldEqual, 0)
			So(err, ShouldBeNil)
		})
		Convey("with `Trailer Hitches`", func() {
			cat, err := GetByTitle("Trailer Hitches")
			So(cat, ShouldNotBeNil)
			So(cat.CategoryId, ShouldNotEqual, 0)
			So(err, ShouldBeNil)
		})
	})

}

func TestGetById(t *testing.T) {
	Convey("Test GetById", t, func() {
		Convey("with 0", func() {
			cat, err := GetById(0)
			So(cat, ShouldNotBeNil)
			So(cat.CategoryId, ShouldEqual, 0)
			So(err, ShouldBeNil)
		})

		Convey("with 1", func() {
			cat, err := GetById(1)
			So(cat, ShouldNotBeNil)
			So(cat.CategoryId, ShouldNotEqual, 0)
			So(err, ShouldBeNil)
		})
	})
}

func TestSubCategories(t *testing.T) {
	Convey("Test Subcategories()", t, func() {
		Convey("with invalid category", func() {
			var cat Category
			subs, err := cat.SubCategories()
			So(subs, ShouldBeNil)
			So(err, ShouldBeNil)
		})
		Convey("with valid category `1`", func() {
			cat, err := GetById(1)
			So(cat, ShouldNotBeNil)
			So(cat.CategoryId, ShouldNotEqual, 0)
			So(err, ShouldBeNil)

			subs, err := cat.SubCategories()
			So(subs, ShouldNotBeNil)
			So(err, ShouldBeNil)
		})

	})
}

func TestGetCategory(t *testing.T) {
	Convey("Test GetCategory()", t, func() {
		Convey("with empty category", func() {
			cat := Category{
				CategoryId: 0,
			}

			ext, err := cat.GetCategory("")
			So(ext, ShouldNotBeNil)
			So(ext.CategoryId, ShouldEqual, 0)
			So(err, ShouldNotBeNil)
		})

		Convey("with `1` category", func() {
			cat := Category{
				CategoryId: 1,
			}

			ext, err := cat.GetCategory("8AEE0620-412E-47FC-900A-947820EA1C1D")
			So(ext, ShouldNotBeNil)
			So(ext.CategoryId, ShouldNotEqual, 0)
			So(err, ShouldNotBeNil)
		})

	})
}