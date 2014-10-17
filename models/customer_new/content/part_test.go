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
			var c PartContent
			var c2 PartContent
			if len(content) > 0 {
				c = content[rand.Intn(len(content))]
				c2 = content[rand.Intn(len(content))]
			}
			id1 := strconv.Itoa(c.PartId)
			id2 := strconv.Itoa(c2.PartId)
			ids := []string{id1, id2}

			Convey("Testing GetPartContent", func() {
				var cus []CustomerContent
				partContent, err := GetPartContent(c.PartId, key)
				So(err, ShouldBeNil)
				So(partContent, ShouldNotBeNil)
				So(partContent, ShouldHaveSameTypeAs, cus)
			})
			Convey("Testing GetGroupedPartContent", func() {

				partContent, err := GetGroupedPartContent(ids, key)
				So(err, ShouldBeNil)
				So(partContent, ShouldNotBeNil)
			})

			//Tests compare part content to old Part content model
			// Convey("Comparative Tests", func() {
			// 	Convey("All Part Content", func() {
			// 		content, err := GetAllPartContent(key)
			// 		So(err, ShouldBeNil)
			// 		old, err := custcontent.GetAllPartContent(key)
			// 		So(err, ShouldBeNil)
			// 		So(len(content), ShouldResemble, len(old))
			// 	})
			// 	Convey("Part Content", func() {
			// 		_, key := getApiKey(allCustContent)
			// 		content, err := GetPartContent(c.PartId, key)
			// 		So(err, ShouldBeNil)
			// 		old, err := custcontent.GetPartContent(c.PartId, key)
			// 		So(err, ShouldBeNil)
			// 		So(len(content), ShouldResemble, len(old))
			// 	})
			// 	Convey("Grouped Part Content", func() {
			// 		_, key := getApiKey(allCustContent)
			// 		content, err := GetGroupedPartContent(ids, key)
			// 		So(err, ShouldBeNil)
			// 		old, err := custcontent.GetGroupedPartContent(ids, key)
			// 		So(err, ShouldBeNil)
			// 		So(len(content), ShouldResemble, len(old))
			// 	})
			// })
		})
	})

}
