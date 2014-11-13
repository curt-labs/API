package apiKeyType

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestApiKeyType(t *testing.T) {
	Convey("Test Create AppGuide", t, func() {
		var err error
		var akt ApiKeyType

		//create
		akt.Type = "testType"

		err = akt.Create()
		So(err, ShouldBeNil)

		//get
		err = akt.Get()
		So(err, ShouldBeNil)

		as, err := GetAllApiKeyTypes()
		So(err, ShouldBeNil)
		So(len(as), ShouldBeGreaterThan, 0)

		//delete
		err = akt.Delete()
		So(err, ShouldBeNil)
	})
}

func BenchmarkGetAllApiKeyTypes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetAllApiKeyTypes()
	}
}

func BenchmarkGetApiKeyType(b *testing.B) {
	akt := ApiKeyType{Type: "TESTER"}
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		akt.Create()
		b.StartTimer()
		akt.Get()
		b.StopTimer()
		akt.Delete()
	}
}

func BenchmarkCreateApiKeyType(b *testing.B) {
	akt := ApiKeyType{Type: "TESTER"}
	for i := 0; i < b.N; i++ {
		b.StartTimer()
		akt.Create()
		b.StopTimer()
		akt.Delete()
	}
}

func BenchmarkDeleteApiKeyType(b *testing.B) {
	akt := ApiKeyType{Type: "TESTER"}
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		akt.Create()
		b.StartTimer()
		akt.Delete()
	}
}
