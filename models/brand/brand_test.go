package brand

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestBrands(t *testing.T) {
	var err error
	b := setupDummyBrand()

	Convey("Testing GetAll", t, func() {
		brands, err := GetAllBrands()
		So(err, ShouldBeNil)
		So(len(brands), ShouldBeGreaterThanOrEqualTo, 0)
	})

	Convey("Testing Brands - CRUD", t, func() {
		Convey("Testing Create", func() {
			err = b.Create()
			So(err, ShouldBeNil)
			So(b.ID, ShouldNotEqual, 0)

			Convey("Testing Read/Get", func() {
				err = b.Get()
				So(err, ShouldBeNil)
				So(b.ID, ShouldBeGreaterThan, 0)
				So(b.Name, ShouldEqual, "TESTER")

				Convey("Testing Update", func() {
					b.Name = "TESTING"
					err = b.Update()
					So(err, ShouldBeNil)
					So(b.Name, ShouldEqual, "TESTING")

					Convey("Testing Delete", func() {
						err = b.Delete()
						So(err, ShouldBeNil)
					})
				})
			})
		})
		Convey("Testing Get - Bad ID", func() {
			br := Brand{}
			err = br.Get()
			So(err, ShouldNotBeNil)
		})
	})
}

func BenchmarkGetAllBrands(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetAllBrands()
	}
}

func BenchmarkGetBrand(b *testing.B) {
	br := setupDummyBrand()
	br.Create()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		br.Get()
	}
	b.StopTimer()
	br.Delete()
}

func BenchmarkCreateBrand(b *testing.B) {
	br := setupDummyBrand()
	for i := 0; i < b.N; i++ {
		br.Create()
		b.StopTimer()
		br.Delete()
		b.StartTimer()
	}
}

func BenchmarkUpdateBrand(b *testing.B) {
	br := setupDummyBrand()
	br.Create()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		br.Name = "TESTING"
		br.Code = "TEST"
		br.Update()
	}
	b.StopTimer()
	br.Delete()
}

func BenchmarkDeleteBrand(b *testing.B) {
	br := setupDummyBrand()
	br.Create()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		br.Delete()
	}
	b.StopTimer()
	br.Delete()
}

func setupDummyBrand() *Brand {
	return &Brand{
		Name: "TESTER",
		Code: "TESTER",
	}
}
