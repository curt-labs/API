package custcontent

import (
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
			x := rand.Intn(len(content))
			c := content[x]

			Convey("Testing GetCategoryContent", func() {
				var cus []CustomerContent
				custContent, err := GetCategoryContent(c.CategoryId, key)
				So(err, ShouldBeNil)
				So(custContent, ShouldNotBeNil)
				So(custContent, ShouldHaveSameTypeAs, cus)
			})
		})
	})
}
