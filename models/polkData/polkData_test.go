package polkData

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestPolk(t *testing.T) {

	// Convey("Test Polk Csv", t, func() {
	// 	file := "/Users/macuser/Desktop/Polk/AriesTestData.csv"
	// 	// file := "/Users/macuser/Desktop/Polk/Aries_Offroad_Coverage_US_201410.csv"
	// 	cs, partsNeededFile, missingBaseVehicles, missingSubmodels, err := CaptureCsv(file, 1, false, false)
	// 	So(err, ShouldBeNil)
	// 	So(len(cs), ShouldBeGreaterThan, 0)
	// 	So(partsNeededFile, ShouldHaveSameTypeAs, &[]CsvDatum{})
	// 	So(missingBaseVehicles, ShouldHaveSameTypeAs, &[]CsvDatum{})
	// 	So(missingSubmodels, ShouldHaveSameTypeAs, &[]CsvDatum{})

	// })

	Convey("Test Polk Diff", t, func() {
		// file := "/Users/macuser/Desktop/Polk/AriesTestData.csv"
		// file := "/Users/macuser/Desktop/Polk/CurtTestData.csv"
		file := "/Users/macuser/Desktop/Polk/Aries_Offroad_Coverage_US_201410.csv"
		// file := "/Users/macuser/Desktop/Polk/Trailer_Hitches_Coverage_US_201410.csv"
		err := Run(file, 1, true, false)
		So(err, ShouldBeNil)
	})

	// Convey("Test maps", t, func() {
	// 	configMap, _ := GetConfigMap()
	// 	So(configMap, ShouldNotBeNil)
	// 	t.Log((configMap))
	// })

}
