package products

import (
	"database/sql"
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

// func TestTopTierCategories(t *testing.T) {
// 	Convey("Test TopTierCategories()", t, func() {
// 		// cats, err := TopTierCategories("9300f7bc-2ca6-11e4-8758-42010af0fd79")
// 		cats, err := TopTierCategories(generateAPIkey())
// 		So(cats, ShouldNotBeNil)
// 		So(err, ShouldBeNil)
// 		So(cats, ShouldNotBeEmpty)
// 	})
// }

func TestGetCategoryByTitle(t *testing.T) {

	Convey("Test GetCategoryByTitle", t, func() {
		Convey("with ``", func() {
			cat, err := GetCategoryByTitle("")
			So(cat, ShouldNotBeNil)
			So(cat.ID, ShouldEqual, 0)
			So(err, ShouldBeNil)
		})
		Convey("with `test`", func() {
			cat, err := GetCategoryByTitle("test")
			So(cat, ShouldNotBeNil)
			So(cat.ID, ShouldEqual, 0)
			So(err, ShouldBeNil)
		})
		Convey("with `Trailer Hitches`", func() {
			cat, err := GetCategoryByTitle("Trailer Hitches")
			So(cat, ShouldNotBeNil)
			t.Log(cat)
			So(cat.ID, ShouldNotEqual, 0)
			So(err, ShouldBeNil)
		})
	})

}

func TestGetCategoryById(t *testing.T) {
	Convey("Test GetCategoryById", t, func() {
		Convey("with 0", func() {
			cat, err := GetCategoryById(0)
			So(cat, ShouldNotBeNil)
			So(cat.ID, ShouldEqual, 0)
			So(err, ShouldBeNil)
		})

		Convey("with 1", func() {
			cat, err := GetCategoryById(1)
			So(cat, ShouldNotBeNil)
			So(cat.ID, ShouldNotEqual, 0)
			So(err, ShouldBeNil)
		})
	})
}

func TestSubCategories(t *testing.T) {
	Convey("Test Subcategories()", t, func() {
		Convey("with invalid category", func() {
			var cat Category
			subs, err := cat.GetSubCategories()
			So(subs, ShouldBeNil)
			So(err, ShouldBeNil)
		})
		Convey("with valid category `1`", func() {
			cat, err := GetCategoryById(1)
			So(cat, ShouldNotBeNil)
			So(cat.ID, ShouldNotEqual, 0)
			So(err, ShouldBeNil)

			subs, err := cat.GetSubCategories()
			So(subs, ShouldNotBeNil)
			So(err, ShouldBeNil)
		})

	})
}

func TestGetCategory(t *testing.T) {
	Convey("Test GetCategory()", t, func() {
		Convey("with empty category", func() {
			cat := Category{
				ID: 0,
			}

			err := cat.GetCategory("", 0, 10, false, nil)
			So(cat, ShouldNotBeNil)
			So(cat.ID, ShouldEqual, 0)
			So(err, ShouldNotBeNil)
		})

		Convey("with `1` category", func() {
			cat := Category{
				ID: 1,
			}

			err := cat.GetCategory("9300f7bc-2ca6-11e4-8758-42010af0fd79", 0, 10, false, nil)
			So(cat, ShouldNotBeNil)
			So(cat.ID, ShouldNotEqual, 0)
			So(err, ShouldBeNil)
			So(cat.ProductListing, ShouldNotBeNil)
			So(len(cat.ProductListing.Parts), ShouldBeLessThanOrEqualTo, 10)
		})

	})
}
