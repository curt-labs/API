package custcontent

import (
	// "github.com/curt-labs/GoAPI/helpers/database"
	"github.com/curt-labs/GoAPI/models/customer/content"
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

			Convey("Comparative Tests", func() {
				Convey("All Part Content", func() {
					content, err := GetAllPartContent(key)
					So(err, ShouldBeNil)
					old, err := custcontent.GetAllPartContent(key)
					So(err, ShouldBeNil)
					So(len(content), ShouldResemble, len(old))
				})
				Convey("Part Content", func() {
					_, key := getApiKey(allCustContent)
					content, err := GetPartContent(c.PartId, key)
					So(err, ShouldBeNil)
					old, err := custcontent.GetPartContent(c.PartId, key)
					So(err, ShouldBeNil)
					So(len(content), ShouldResemble, len(old))
				})
				Convey("Grouped Part Content", func() {
					_, key := getApiKey(allCustContent)
					content, err := GetGroupedPartContent(ids, key)
					So(err, ShouldBeNil)
					old, err := custcontent.GetGroupedPartContent(ids, key)
					So(err, ShouldBeNil)
					So(len(content), ShouldResemble, len(old))
				})
			})
		})
	})

}
