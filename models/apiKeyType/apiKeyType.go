package apiKeyType

import (
	"database/sql"
	"github.com/curt-labs/GoAPI/helpers/database"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

var (
	getApiKeyType     = "SELECT id, type, date_added FROM ApiKeyType WHERE id = ? "
	getAllApiKeyTypes = "SELECT id, type, date_added FROM ApiKeyType "
	getKeyByDateType  = "SELECT id FROM ApiKeyType WHERE type = ?  AND date_added = ?"
	createApiKeyType  = "INSERT INTO ApiKeyType (id, type, date_added) VALUES (UUID(),?,?)"
	deleteApiKeyType  = "DELETE FROM ApiKeyType WHERE id = ? "
)

type ApiKeyType struct {
	ID        string
	Type      string
	DateAdded time.Time
	time.Time
}

type Scanner interface {
	Scan(...interface{}) error
}

func (a *ApiKeyType) Get() (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return
	}
	defer db.Close()
	stmt, err := db.Prepare(getApiKeyType)
	if err != nil {
		return
	}
	defer stmt.Close()
	res := stmt.QueryRow(a.ID)
	a, err = ScanKey(res)

	return
}

func GetAllApiKeyTypes() (as []ApiKeyType, err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return
	}
	defer db.Close()
	stmt, err := db.Prepare(getAllApiKeyTypes)
	if err != nil {
		return
	}
	defer stmt.Close()
	res, err := stmt.Query() //returns *sql.Rows

	for res.Next() {
		a, err := ScanKey(res)
		if err != nil {
			return as, err
		}
		as = append(as, *a)
	}
	return as, err
}

func ScanKey(s Scanner) (*ApiKeyType, error) {
	a := &ApiKeyType{}
	err := s.Scan(&a.ID, &a.Type, &a.DateAdded)
	return a, err
}

func (a *ApiKeyType) Create() (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return
	}
	defer db.Close()
	stmt, err := db.Prepare(createApiKeyType)
	if err != nil {
		return
	}
	defer stmt.Close()
	a.DateAdded = time.Now()
	_, err = stmt.Exec(a.Type, a.DateAdded)
	if err != nil {
		return
	}

	stmt, err = db.Prepare(getKeyByDateType)
	if err != nil {
		return err
	}

	defer stmt.Close()
	err = stmt.QueryRow(a.Type, a.DateAdded).Scan(&a.ID)
	if err != nil {
		return err
	}
	return
}

func (a *ApiKeyType) Delete() (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	stmt, err := db.Prepare(deleteApiKeyType)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(a.ID)
	if err != nil {
		return err
	}
	return
}
