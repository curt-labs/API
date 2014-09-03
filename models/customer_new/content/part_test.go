package custcontent

import (
	. "github.com/smartystreets/goconvey/convey"
	"math/rand"
	"strconv"
	"testing"
)

func TestPart(t *testing.T) {
	Convey("Testing Part", t, func() {
		Convey("Testing  GetAllPartContent()", func() {
			_, key := getApiKey(allCustContent)
			var con []PartContent
			content, err := GetAllPartContent(key)
			So(err, ShouldBeNil)
			So(content, ShouldHaveSameTypeAs, con)
			x := rand.Intn(len(content))
			c := content[x]
			x2 := rand.Intn(len(content))
			c2 := content[x2]

			Convey("Testing GetPartContent", func() {
				var cus []CustomerContent
				partContent, err := GetPartContent(c.PartId, key)
				So(err, ShouldBeNil)
				So(partContent, ShouldNotBeNil)
				So(partContent, ShouldHaveSameTypeAs, cus)
			})
			Convey("Testing GetGroupedPartContent", func() {
				id1 := strconv.Itoa(c.PartId)
				id2 := strconv.Itoa(c2.PartId)
				ids := []string{id1, id2}
				partContent, err := GetGroupedPartContent(ids, key)
				So(err, ShouldBeNil)
				So(partContent, ShouldNotBeNil)
			})
		})
	})
}
