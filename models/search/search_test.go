package search

import (
	"github.com/curt-labs/API/helpers/apicontextmock"
	. "github.com/smartystreets/goconvey/convey"
	"os"
	"testing"
)

func TestDsl(t *testing.T) {
	dtx, err := apicontextmock.Mock()
	if err != nil {
		t.Log(err)
		return
	}
	ip := os.Getenv("ELASTICSEARCH_IP")
	port := os.Getenv("ELASTICSEARCH_PORT")
	user := os.Getenv("ELASTICSEARCH_USER")
	pass := os.Getenv("ELASTICSEARCH_PASS")
	Convey("Testing Search Dsl", t, func() {

		Convey("empty query", func() {
			res, err := Dsl("", 0, 0, dtx)
			So(err, ShouldNotBeNil)
			So(res, ShouldBeNil)
		})
		Convey("query of `hitch` but bad connections", func() {
			os.Setenv("ELASTICSEARCH_IP", "")
			os.Setenv("ELASTICSEARCH_PORT", "")
			os.Setenv("ELASTICSEARCH_USER", "")
			os.Setenv("ELASTICSEARCH_PASS", "")
			res, err := Dsl("hitch", 0, 0, dtx)
			So(err, ShouldNotBeNil)
			So(res, ShouldBeNil)
			os.Setenv("ELASTICSEARCH_IP", ip)
			os.Setenv("ELASTICSEARCH_PORT", port)
			os.Setenv("ELASTICSEARCH_USER", user)
			os.Setenv("ELASTICSEARCH_PASS", pass)
		})
		Convey("query of `hitch` with no brand", func() {
			dtx.BrandArray = []int{}
			res, err := Dsl("hitch", 1, 0, dtx)
			So(err, ShouldNotBeNil)
			So(res, ShouldBeNil)
		})
		Convey("query of `hitch` with aries", func() {
			dtx.BrandArray = []int{3}
			res, err := Dsl("hitch", 0, 0, dtx)
			So(err, ShouldBeNil)
			So(res, ShouldNotBeNil)
		})
		Convey("query of `hitch` with curt", func() {
			dtx.BrandArray = []int{1}
			res, err := Dsl("hitch", 0, 0, dtx)
			So(err, ShouldBeNil)
			So(res, ShouldNotBeNil)
		})
		Convey("query of `hitch`", func() {
			dtx.BrandArray = []int{1, 3}
			res, err := Dsl("hitch", 0, 0, dtx)
			So(err, ShouldBeNil)
			So(res, ShouldNotBeNil)
		})
		Convey("query of `hitch` with 1 result", func() {
			dtx.BrandArray = []int{1, 3}
			res, err := Dsl("hitch", 0, 1, dtx)
			So(err, ShouldBeNil)
			So(res, ShouldNotBeNil)
			So(res.Hits.Len(), ShouldEqual, 1)
		})

	})
}
