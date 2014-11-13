package geography

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestGeography(t *testing.T) {
	Convey("Testing Gets", t, func() {
		Convey("Testing GetAllCountriesAndStates()", func() {
			countrystates, err := GetAllCountriesAndStates()
			So(len(countrystates), ShouldBeGreaterThanOrEqualTo, 0)
			So(err, ShouldBeNil)
		})

		Convey("Testing GetAllCountries", func() {
			countries, err := GetAllCountries()
			So(len(countries), ShouldBeGreaterThanOrEqualTo, 0)
			So(err, ShouldBeNil)
		})

		Convey("Testing GetAllStates", func() {
			states, err := GetAllStates()
			So(len(states), ShouldBeGreaterThanOrEqualTo, 0)
			So(err, ShouldBeNil)
		})
	})
}

func BenchmarkGetAllCountriesAndStates(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetAllCountriesAndStates()
	}
}

func BenchmarkGetAllCountries(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetAllCountries()
	}
}

func BenchmarkGetAllStates(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetAllStates()
	}
}
