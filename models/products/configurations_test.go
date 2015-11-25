package products

import (
	"github.com/curt-labs/API/helpers/api"
	"github.com/curt-labs/API/helpers/apicontextmock"
	. "github.com/smartystreets/goconvey/convey"
	"math/rand"
	"testing"
)

func TestGetConfigurations(t *testing.T) {
	var l Lookup
	l.Brands = append(l.Brands, 1)
	MockedDTX, err := apicontextmock.Mock()
	if err != nil {
		return
	}

	Convey("Testing GetConfigurations()", t, func() {

		Convey("without year/make/model", func() {
			err := l.GetConfigurations()
			So(err, ShouldEqual, nil)
			So(l.Configurations, ShouldNotEqual, nil)
			So(len(l.Configurations), ShouldEqual, 0)
		})

		Convey("with year", func() {
			err := l.GetYears(MockedDTX)
			So(err, ShouldEqual, nil)
			if len(l.Years) == 0 {
				return
			}
			l.Vehicle.Base.Year = l.Years[api_helpers.RandGenerator(len(l.Years)-1)]

			err = l.GetConfigurations()
			So(err, ShouldEqual, nil)
			So(l.Configurations, ShouldNotEqual, nil)
			So(len(l.Configurations), ShouldEqual, 0)
		})

		Convey("with year/make", func() {
			err := l.GetYears(MockedDTX)
			So(err, ShouldEqual, nil)
			if len(l.Years) == 0 {
				return
			}
			l.Vehicle.Base.Year = l.Years[api_helpers.RandGenerator(len(l.Years)-1)]
			err = l.GetMakes(MockedDTX)
			So(err, ShouldEqual, nil)
			if len(l.Makes) == 0 {
				return
			}
			l.Vehicle.Base.Make = l.Makes[api_helpers.RandGenerator(len(l.Makes)-1)]

			err = l.GetConfigurations()
			So(err, ShouldEqual, nil)
			So(l.Configurations, ShouldNotEqual, nil)
			So(len(l.Configurations), ShouldEqual, 0)
		})

		Convey("with year/make/model", func() {

			err := l.GetYears(MockedDTX)
			So(err, ShouldEqual, nil)
			if len(l.Years) == 0 {
				return
			}
			l.Vehicle.Base.Year = l.Years[api_helpers.RandGenerator(len(l.Years)-1)]

			err = l.GetMakes(MockedDTX)
			So(err, ShouldEqual, nil)
			if len(l.Makes) == 0 {
				return
			}
			l.Vehicle.Base.Make = l.Makes[api_helpers.RandGenerator(len(l.Makes)-1)]

			err = l.GetModels()
			So(err, ShouldEqual, nil)
			if len(l.Models) == 0 {
				return
			}
			l.Vehicle.Base.Model = l.Models[api_helpers.RandGenerator(len(l.Models)-1)]

			err = l.GetConfigurations()
			So(err, ShouldEqual, nil)
			So(l.Configurations, ShouldNotEqual, nil)
			if (len(l.Configurations)) > 0 {
				So(len(l.Configurations), ShouldNotEqual, 0)
				idx := 0
				if len(l.Configurations) > 1 {
					idx = rand.Intn(len(l.Submodels) - 1)
				}

				So(l.Configurations[idx], ShouldNotEqual, nil)
				So(l.Configurations[idx].Type, ShouldNotEqual, nil)
				So(l.Configurations[idx].Type, ShouldNotEqual, "")
				So(l.Configurations[idx].Options, ShouldNotEqual, nil)
				So(l.Configurations[idx].Options, ShouldNotBeEmpty)
			}
		})

		Convey("with year/make/model/submodel", func() {

			err := l.GetYears(MockedDTX)
			So(err, ShouldEqual, nil)
			if len(l.Years) == 0 {
				return
			}
			l.Vehicle.Base.Year = l.Years[api_helpers.RandGenerator(len(l.Years)-1)]

			err = l.GetMakes(MockedDTX)
			So(err, ShouldEqual, nil)
			if len(l.Makes) == 0 {
				return
			}
			l.Vehicle.Base.Make = l.Makes[api_helpers.RandGenerator(len(l.Makes)-1)]

			err = l.GetModels()
			So(err, ShouldEqual, nil)
			if len(l.Models) == 0 {
				return
			}
			l.Vehicle.Base.Model = l.Models[api_helpers.RandGenerator(len(l.Models)-1)]

			err = l.GetSubmodels()
			So(err, ShouldEqual, nil)
			So(l.Submodels, ShouldNotEqual, nil)
			if (len(l.Submodels)) == 0 {
				return
			}
			l.Vehicle.Submodel = l.Submodels[api_helpers.RandGenerator(len(l.Submodels)-1)]

			err = l.GetConfigurations()
			So(err, ShouldEqual, nil)
			So(l.Configurations, ShouldNotEqual, nil)
			//TODO fix - runtime error

			if len(l.Configurations) > 0 {
				So(len(l.Configurations), ShouldNotEqual, 0)
				idx := 0
				if len(l.Configurations) > 1 {
					idx = rand.Intn(len(l.Submodels) - 1)
				}
				t.Log(idx)

				So(l.Configurations[idx], ShouldNotEqual, nil)
				So(l.Configurations[idx].Type, ShouldNotEqual, nil)
				So(l.Configurations[idx].Type, ShouldNotEqual, "")
				So(l.Configurations[idx].Options, ShouldNotEqual, nil)
				So(l.Configurations[idx].Options, ShouldNotBeEmpty)
			}
		})

		Convey("with bogus data", func() {
			l.Vehicle.Base.Year = 1
			l.Vehicle.Base.Make = "KD"
			l.Vehicle.Base.Model = "123"
			l.Vehicle.Submodel = "LKJ"
			err := l.GetConfigurations()
			So(err, ShouldEqual, nil)
			So(l.Configurations, ShouldNotEqual, nil)
			So(len(l.Configurations), ShouldEqual, 0)
		})
	})

	Convey("Test getDefinedConfigurations()", t, func() {

		configs, err := l.Vehicle.getDefinedConfigurations(MockedDTX.APIKey)
		So(err, ShouldEqual, nil)
		So(configs, ShouldNotEqual, nil)

		Convey("with year/make/model/submodel", func() {

			err := l.GetYears(MockedDTX)
			So(err, ShouldEqual, nil)
			if len(l.Years) == 0 {
				return
			}
			l.Vehicle.Base.Year = l.Years[api_helpers.RandGenerator(len(l.Years)-1)]

			err = l.GetMakes(MockedDTX)
			So(err, ShouldEqual, nil)
			if len(l.Makes) == 0 {
				return
			}
			l.Vehicle.Base.Make = l.Makes[api_helpers.RandGenerator(len(l.Makes)-1)]

			err = l.GetModels()
			So(err, ShouldEqual, nil)
			if len(l.Models) == 0 {
				return
			}
			l.Vehicle.Base.Model = l.Models[api_helpers.RandGenerator(len(l.Models)-1)]

			err = l.GetSubmodels()
			So(err, ShouldEqual, nil)
			So(l.Submodels, ShouldNotEqual, nil)
			if (len(l.Submodels)) == 0 {
				return
			}
			l.Vehicle.Submodel = l.Submodels[api_helpers.RandGenerator(len(l.Submodels)-1)]

			configs, err := l.Vehicle.getDefinedConfigurations(MockedDTX.APIKey)
			So(err, ShouldEqual, nil)
			So(configs, ShouldNotEqual, nil)
		})
	})
	_ = apicontextmock.DeMock(MockedDTX)

}
