package products

import (
	"database/sql"
	"github.com/curt-labs/GoAPI/helpers/apicontext"
	"github.com/curt-labs/GoAPI/helpers/apicontextmock"
	"github.com/curt-labs/GoAPI/helpers/database"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func generateAPIkey() (apiKey string) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return ""
	}
	defer db.Close()
	stmt, err := db.Prepare("SELECT api_key FROM ApiKey ORDER BY RAND() LIMIT 1")
	if err != nil {
		return ""
	}
	defer stmt.Close()
	err = stmt.QueryRow().Scan(&apiKey)
	if err != nil {
		return ""
	}
	return apiKey
}

func TestTopTierCategories(t *testing.T) {
	var err error
	MockedDTX := &apicontext.DataContext{}
	if MockedDTX, err = apicontextmock.Mock(); err != nil {
		return
	}

	var cat Category
	cat.ParentID = 0
	cat.IsLifestyle = false
	cat.Title = "test cat"

	Convey("Create", t, func() {
		err := cat.Create()
		So(err, ShouldBeNil)
	})

	Convey("Test TopTierCategories()", t, func() {
		cats, err := TopTierCategories(generateAPIkey(), MockedDTX.BrandID)
		So(cats, ShouldNotBeNil)
		So(err, ShouldBeNil)
		So(cats, ShouldHaveSameTypeAs, []Category{})
	})

	Convey("Test GetCategoryByTitle", t, func() {
		Convey("with ``", func() {
			cat, err := GetCategoryByTitle("", MockedDTX)
			So(cat, ShouldNotBeNil)
			So(cat.ID, ShouldBeGreaterThanOrEqualTo, 0)
			So(err, ShouldBeNil)
		})
		Convey("with `test`", func() {
			cat, err := GetCategoryByTitle("test", MockedDTX)
			So(cat, ShouldNotBeNil)
			So(cat.ID, ShouldEqual, 0)
			So(err, ShouldBeNil)
		})
		Convey("with `Trailer Hitches`", func() {
			cat, err := GetCategoryByTitle("Trailer Hitches", MockedDTX)
			So(cat, ShouldNotBeNil)
			So(cat, ShouldHaveSameTypeAs, Category{})
			So(err, ShouldBeNil)
		})
	})

	Convey("Test GetCategoryById", t, func() {
		Convey("with 0", func() {
			cat, err := GetCategoryById(MockedDTX.BrandID, 0)
			So(cat, ShouldNotBeNil)
			So(cat.ID, ShouldEqual, 0)
			So(err, ShouldBeNil)
		})

		Convey("with 1", func() {
			cat, err := GetCategoryById(MockedDTX.BrandID, 1)
			So(cat, ShouldNotBeNil)
			So(cat, ShouldHaveSameTypeAs, Category{})
			So(err, ShouldBeNil)
		})
	})

	Convey("Test Subcategories()", t, func() {
		Convey("with invalid category", func() {
			var cat Category
			subs, err := cat.GetSubCategories()
			So(subs, ShouldNotBeNil)
			So(err, ShouldBeNil)
		})
		Convey("with valid category `1`", func() {
			cat, err := GetCategoryById(MockedDTX.BrandID, 1)
			So(cat, ShouldNotBeNil)
			So(cat, ShouldHaveSameTypeAs, Category{})
			So(err, ShouldBeNil)

			subs, err := cat.GetSubCategories()
			So(subs, ShouldNotBeNil)
			So(err, ShouldBeNil)
		})

	})

	Convey("Test GetCategory()", t, func() {
		Convey("with empty category", func() {
			cat := Category{
				ID: 0,
			}

			err := cat.GetCategory("", 0, 10, false, nil, nil)
			So(cat, ShouldNotBeNil)
			So(cat.ID, ShouldEqual, 0)
			So(err, ShouldNotBeNil)
		})

		Convey("with `1` category without parts", func() {
			cat := Category{
				ID: 1,
			}

			err := cat.GetCategory("9300f7bc-2ca6-11e4-8758-42010af0fd79", 0, 10, true, nil, nil)
			So(cat, ShouldNotBeNil)
			So(cat, ShouldHaveSameTypeAs, Category{})
			So(err, ShouldBeNil)
			So(cat.ProductListing, ShouldBeNil)
		})

		// Convey("with `1` category with parts", func() {
		// 	cat := Category{
		// 		ID: 1,
		// 	}

		// 	err := cat.GetCategory("9300f7bc-2ca6-11e4-8758-42010af0fd79", 0, 10, false, nil, nil)
		// 	So(cat, ShouldHaveSameTypeAs, Category{})
		// 	So(err, ShouldBeNil)
		// 	So(cat.ProductListing, ShouldNotBeNil)
		// 	So(len(cat.ProductListing.Parts), ShouldBeLessThanOrEqualTo, 10)
		// })

	})
	Convey("Delete", t, func() {
		err := cat.Delete()
		So(err, ShouldBeNil)

	})
	_ = apicontextmock.DeMock(MockedDTX)
}
