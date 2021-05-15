package apiKeyType

import (
	"database/sql"
	"errors"
	"fmt"
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

func TestApiKeyType_GetAllApiKeyTypes(t *testing.T) {
	var getTests = []struct {
		name     string
		rows     *sqlmock.Rows
		queryErr error
		outErr   error
		outAkt   []ApiKeyType
	}{
		{
			name:   "no API key types found",
			rows:   sqlmock.NewRows([]string{"id", "type", "date_added"}),
			outErr: nil,
			outAkt: []ApiKeyType{},
		},
		{
			name:     "query failed",
			rows:     sqlmock.NewRows([]string{"id", "type", "date_added"}),
			queryErr: errors.New("query error"),
			outErr:   errors.New("query error"),
			outAkt:   nil,
		},
		{
			name: "API key type found",
			//in:   &ApiKeyType{ID: "00000000-0000-0000-0000-000000000000"},
			rows: sqlmock.NewRows([]string{"id", "type", "date_added"}).
				AddRow("99900000-0000-0000-0000-000000000000", "TestKey", time.Date(2020, 4, 25, 12, 14, 00, 00, time.UTC)).
				AddRow("99800000-0000-0000-0000-000000000000", "TestKey2", time.Date(2020, 4, 26, 12, 14, 00, 00, time.UTC)).
				AddRow("99700000-0000-0000-0000-000000000000", "TestKey3", time.Date(2020, 4, 27, 12, 14, 00, 00, time.UTC)),
			outErr: nil,
			outAkt: []ApiKeyType{
				{ID: "99900000-0000-0000-0000-000000000000", Type: "TestKey", DateAdded: time.Date(2020, 4, 25, 12, 14, 00, 00, time.UTC)},
				{ID: "99800000-0000-0000-0000-000000000000", Type: "TestKey2", DateAdded: time.Date(2020, 4, 26, 12, 14, 00, 00, time.UTC)},
				{ID: "99700000-0000-0000-0000-000000000000", Type: "TestKey3", DateAdded: time.Date(2020, 4, 27, 12, 14, 00, 00, time.UTC)},
			},
		},
		{
			name: "scan error API key type found",
			//in:   &ApiKeyType{ID: "00000000-0000-0000-0000-000000000000"},
			rows: sqlmock.NewRows([]string{"id", "type", "date_added"}).
				AddRow("99900000-0000-0000-0000-000000000000", "TestKey", time.Date(2020, 4, 25, 12, 14, 00, 00, time.UTC)).
				AddRow("99800000-0000-0000-0000-000000000000", "TestKey2", "asfd"). // causes Scan error
				AddRow("99700000-0000-0000-0000-000000000000", "TestKey3", time.Date(2020, 4, 27, 12, 14, 00, 00, time.UTC)),
			outErr: fmt.Errorf(`sql: Scan error on column index %d, name %q: %w`, 2, "date_added", errors.New("unsupported Scan, storing driver.Value type string into type *time.Time")),
			outAkt: []ApiKeyType{
				{ID: "99900000-0000-0000-0000-000000000000", Type: "TestKey", DateAdded: time.Date(2020, 4, 25, 12, 14, 00, 00, time.UTC)},
			},
		},
		{
			name: "row error API key type found",
			//in:   &ApiKeyType{ID: "00000000-0000-0000-0000-000000000000"},
			rows: sqlmock.NewRows([]string{"id", "type", "date_added"}).
				AddRow("99900000-0000-0000-0000-000000000000", "TestKey", time.Date(2020, 4, 25, 12, 14, 00, 00, time.UTC)).
				AddRow("99800000-0000-0000-0000-000000000000", "TestKey2", time.Date(2020, 4, 26, 12, 14, 00, 00, time.UTC)).
				AddRow("99700000-0000-0000-0000-000000000000", "TestKey3", time.Date(2020, 4, 27, 12, 14, 00, 00, time.UTC)).
				RowError(1, errors.New("scan error")),
			outErr: errors.New("scan error"),
			outAkt: []ApiKeyType{
				{ID: "99900000-0000-0000-0000-000000000000", Type: "TestKey", DateAdded: time.Date(2020, 4, 25, 12, 14, 00, 00, time.UTC)},
			},
		},
	}

	for _, tt := range getTests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := NewMock()
			defer db.Close()

			query := getAllApiKeyTypes

			var rows *sqlmock.Rows
			if tt.rows != nil {
				rows = tt.rows
			} else {
				rows = sqlmock.NewRows([]string{"id", "type", "date_added"})
			}

			mock.ExpectQuery(query).WillReturnRows(rows).WillReturnError(tt.queryErr)

			akts, err := GetAllApiKeyTypes(db)

			assert.Equal(t, tt.outErr, err)
			assert.Equal(t, tt.outAkt, akts)
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

		//as, err := GetAllApiKeyTypes()
		//So(err, ShouldBeNil)
		//So(len(as), ShouldBeGreaterThan, 0)

		//delete
		err = akt.Delete()
		So(err, ShouldBeNil)
	})
}

func BenchmarkGetAllApiKeyTypes(b *testing.B) {
	//for i := 0; i < b.N; i++ {
	//GetAllApiKeyTypes()
	//}
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
