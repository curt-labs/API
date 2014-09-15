package custcontent

import (
	"github.com/curt-labs/GoAPI/models/customer/content"
	. "github.com/smartystreets/goconvey/convey"
	"math/rand"
	"testing"
)

func TestCategory(t *testing.T) {
	Convey("Testing Category", t, func() {
		Convey("Testing  GetAllCategoryContent()", func() {
			_, key := getApiKey(allCustContent)
			var con []CategoryContent
			content, err := GetAllCategoryContent(key)

			So(err, ShouldBeNil)
			So(content, ShouldHaveSameTypeAs, con)
			var c CategoryContent
			if len(content) > 0 {
				c = content[rand.Intn(len(content))]
			}

			Convey("Testing GetCategoryContent", func() {
				var cus []CustomerContent
				custContent, err := GetCategoryContent(c.CategoryId, key)
				So(err, ShouldBeNil)
				So(custContent, ShouldNotBeNil)
				So(custContent, ShouldHaveSameTypeAs, cus)

				Convey("Comparative Tests", func() {
					Convey("AllCatContent", func() {
						oldContent, err := custcontent.GetAllCategoryContent(key)
						if err == nil { //does the mymysql version work?
							So(content, ShouldResemble, oldContent)
						}
					})
					Convey("CatContent", func() {
						oldCustContent, err := custcontent.GetCategoryContent(c.CategoryId, key)
						if err == nil { //does the mymysql version work?
							So(custContent, ShouldResemble, oldCustContent)
						}
					})
				})
			})
		})
	})
}
