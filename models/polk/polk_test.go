package polk

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestPolk(t *testing.T) {

	// Convey("Test Polk Csv", t, func() {
	// 	file := "/Users/macuser/Desktop/Polk/SampleCurt.csv"
	// 	cs, err := CaptureCsv(file, 1)
	// 	So(err, ShouldBeNil)
	// 	So(len(cs), ShouldBeGreaterThan, 0)
	// })

	Convey("Test Polk Csv", t, func() {
		file := "/Users/macuser/Desktop/Polk/AriesTestData.csv"
		// file := "/Users/macuser/Desktop/Polk/sampleCurt.csv"
		// file := "/Users/macuser/Desktop/Polk/Aries_Offroad_Coverage_US_201410.csv"
		err := RunDiff(file, 1, true)
		So(err, ShouldBeNil)
	})

}
