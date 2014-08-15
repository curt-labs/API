package aces

import (
	"github.com/curt-labs/GoAPI/helpers/api"
	. "github.com/smartystreets/goconvey/convey"
	"math/rand"
	"testing"
)

func TestGetSubmodels(t *testing.T) {
	var l Lookup
	Convey("Testing GetSubmodels()", t, func() {

		Convey("without year/make/model", func() {
			err := l.GetSubmodels()
			So(err, ShouldEqual, nil)
			So(l.Submodels, ShouldNotEqual, nil)
			So(len(l.Submodels), ShouldEqual, 0)
		})

		Convey("with bogus data", func() {
			l.Vehicle.Base.Year = 1
			l.Vehicle.Base.Make = "KD"
			l.Vehicle.Base.Model = "123"
			err := l.GetSubmodels()
			So(err, ShouldEqual, nil)
			So(l.Submodels, ShouldNotEqual, nil)
			So(len(l.Submodels), ShouldEqual, 0)
		})

		Convey("with year", func() {
			err := l.GetYears()
			So(err, ShouldEqual, nil)
			if len(l.Years) == 0 {
				return
			}
			l.Vehicle.Base.Year = l.Years[api_helpers.RandGenerator(len(l.Years))]

			err = l.GetSubmodels()
			So(err, ShouldEqual, nil)
			So(l.Submodels, ShouldNotEqual, nil)
			So(len(l.Submodels), ShouldEqual, 0)
		})

		Convey("with year/make", func() {
			err := l.GetYears()
			So(err, ShouldEqual, nil)
			if len(l.Years) == 0 {
				return
			}
			l.Vehicle.Base.Year = l.Years[api_helpers.RandGenerator(len(l.Years))]
			err = l.GetMakes()
			So(err, ShouldEqual, nil)
			if len(l.Makes) == 0 {
				return
			}
			l.Vehicle.Base.Make = l.Makes[api_helpers.RandGenerator(len(l.Makes))]

			err = l.GetSubmodels()
			So(err, ShouldEqual, nil)
			So(l.Submodels, ShouldNotEqual, nil)
			So(len(l.Submodels), ShouldEqual, 0)
		})

		Convey("with year/make/model", func() {

			err := l.GetYears()
			So(err, ShouldEqual, nil)
			if len(l.Years) == 0 {
				return
			}
			l.Vehicle.Base.Year = l.Years[api_helpers.RandGenerator(len(l.Years)-1)]

			err = l.GetMakes()
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
			if (len(l.Submodels)) > 0 {
				So(len(l.Submodels), ShouldNotEqual, 0)
				idx := 0
				if len(l.Submodels) > 1 {
					idx = rand.Intn(len(l.Submodels) - 1)
				}

				So(l.Submodels[idx], ShouldNotEqual, "")
			}
		})
	})

}
