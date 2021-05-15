package apiKeyType

import (
	"database/sql"
	"log"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/bmizerany/assert"
	. "github.com/smartystreets/goconvey/convey"
)

func NewMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	return db, mock
}

func TestApiKeyType_Get(t *testing.T) {
	var getTests = []struct {
		name   string
		in     *ApiKeyType
		rows   *sqlmock.Rows
		outErr error
		outAkt *ApiKeyType
	}{
		{
			name:   "no API key type found",
			in:     &ApiKeyType{ID: "00000000-0000-0000-0000-000000000000"},
			outErr: sql.ErrNoRows,
			outAkt: &ApiKeyType{ID: "00000000-0000-0000-0000-000000000000"},
		},
		{
			name: "API key type found",
			in:   &ApiKeyType{ID: "00000000-0000-0000-0000-000000000000"},
			rows: sqlmock.NewRows([]string{"id", "type", "date_added"}).
				AddRow("99900000-0000-0000-0000-000000000000", "TestKey", time.Date(2020, 4, 25, 12, 14, 00, 00, time.UTC)),
			outErr: nil,
			outAkt: &ApiKeyType{ID: "99900000-0000-0000-0000-000000000000", Type: "TestKey", DateAdded: time.Date(2020, 4, 25, 12, 14, 00, 00, time.UTC)},
		},
	}

	for _, tt := range getTests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := NewMock()
			defer db.Close()

			query := getApiKeyType

			var rows *sqlmock.Rows
			if tt.rows != nil {
				rows = tt.rows
			} else {
				rows = sqlmock.NewRows([]string{"id", "type", "date_added"})
			}

			mock.ExpectQuery(query).WithArgs(tt.in.ID).WillReturnRows(rows)

			err := tt.in.Get(db)
			assert.Equal(t, tt.outErr, err)
			assert.Equal(t, tt.outAkt, tt.in)
		})
	}
}

func TestApiKeyType(t *testing.T) {
	Convey("Test Create AppGuide", t, func() {
		var err error
		var akt ApiKeyType

		//create
		akt.Type = "testType"

		err = akt.Create()
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
		//akt.Get()
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
