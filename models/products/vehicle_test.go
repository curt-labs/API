package products

import (
	"github.com/curt-labs/GoAPI/helpers/api"
	"github.com/curt-labs/GoAPI/helpers/apicontext"
	"github.com/curt-labs/GoAPI/helpers/apicontextmock"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestLoadParts(t *testing.T) {
	var l Lookup
	var err error
	l.Brands = append(l.Brands, 1)

	MockedDTX := &apicontext.DataContext{}
	if MockedDTX, err = apicontextmock.Mock(); err != nil {
		return
	}
	Convey("Testing LoadParts()", t, func() {

		Convey("without year/make/model", func() {
			ch := make(chan []Part)
			go l.LoadParts(ch, MockedDTX)
			parts := <-ch
			So(parts, ShouldNotEqual, nil)
			So(len(parts), ShouldEqual, 0)
		})

		Convey("with bogus data", func() {
			l.Vehicle.Base.Year = 1
			l.Vehicle.Base.Make = "KD"
			l.Vehicle.Base.Model = "123"
			l.Vehicle.Submodel = "LKJ"
			ch := make(chan []Part)
			go l.LoadParts(ch, MockedDTX)
			parts := <-ch
			So(parts, ShouldNotEqual, nil)
			So(len(parts), ShouldEqual, 0)
		})

		Convey("with year", func() {
			err := l.GetYears(MockedDTX)
			So(err, ShouldEqual, nil)
			if len(l.Years) == 0 {
				return
			}
			l.Vehicle.Base.Year = l.Years[api_helpers.RandGenerator(len(l.Years)-1)]

			ch := make(chan []Part)
			go l.LoadParts(ch, MockedDTX)
			parts := <-ch
			So(parts, ShouldNotEqual, nil)
			So(len(parts), ShouldEqual, 0)
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

			ch := make(chan []Part)
			go l.LoadParts(ch, MockedDTX)
			parts := <-ch
			So(parts, ShouldNotEqual, nil)
			So(len(parts), ShouldEqual, 0)
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

			ch := make(chan []Part)
			go l.LoadParts(ch, MockedDTX)
			parts := <-ch
			So(parts, ShouldNotEqual, nil)
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

			ch := make(chan []Part)
			go l.LoadParts(ch, MockedDTX)
			parts := <-ch
			So(parts, ShouldNotEqual, nil)
		})
	})
	_ = apicontextmock.DeMock(MockedDTX)
}

func TestGetVcdbID(t *testing.T) {
	var l Lookup
	var err error
	l.Brands = append(l.Brands, 1)
	MockedDTX := &apicontext.DataContext{}
	if MockedDTX, err = apicontextmock.Mock(); err != nil {
		return
	}
	Convey("Testing GetVcdbID()", t, func() {

		Convey("without year/make/model", func() {
			id, err := l.Vehicle.GetVcdbID()
			So(err, ShouldNotEqual, nil)
			So(id, ShouldEqual, 0)
		})

		Convey("with year", func() {
			err := l.GetYears(MockedDTX)
			So(err, ShouldEqual, nil)
			if len(l.Years) == 0 {
				return
			}
			l.Vehicle.Base.Year = l.Years[api_helpers.RandGenerator(len(l.Years)-1)]

			id, err := l.Vehicle.GetVcdbID()
			So(err, ShouldNotEqual, nil)
			So(id, ShouldEqual, 0)
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

			id, err := l.Vehicle.GetVcdbID()
			So(err, ShouldNotEqual, nil)
			So(id, ShouldEqual, 0)
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

			id, err := l.Vehicle.GetVcdbID()
			if err != nil {
				So(id, ShouldEqual, 0)
			} else {
				So(id, ShouldNotEqual, 0)
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

			id, err := l.Vehicle.GetVcdbID()
			So(err, ShouldEqual, nil)
			So(id, ShouldNotEqual, 0)
		})

		Convey("with bogus data", func() {
			l.Vehicle.Base.Year = 1
			l.Vehicle.Base.Make = "KD"
			l.Vehicle.Base.Model = "123"
			l.Vehicle.Submodel = "LKJ"
			id, err := l.Vehicle.GetVcdbID()
			So(err, ShouldNotEqual, nil)
			So(id, ShouldEqual, 0)
		})
	})
	_ = apicontextmock.DeMock(MockedDTX)
}
