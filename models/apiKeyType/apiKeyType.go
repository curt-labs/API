package apiKeyType

import (
	"time"

	"github.com/curt-labs/API/helpers/database"
	_ "github.com/go-sql-driver/mysql"
)

var (
	getApiKeyType     = "SELECT id, type, date_added FROM ApiKeyType WHERE id = ? "
	getAllApiKeyTypes = "SELECT id, type, date_added FROM ApiKeyType "
	getKeyByDateType  = "SELECT id FROM ApiKeyType WHERE type = ?  AND date_added = ?"
	createApiKeyType  = "INSERT INTO ApiKeyType (id, type, date_added) VALUES (UUID(),?,?)"
	deleteApiKeyType  = "DELETE FROM ApiKeyType WHERE id = ? "
)

type ApiKeyType struct {
	ID        string    `json:"_id" xml:"id"`
	Type      string    `json:"type" xml:"type"`
	DateAdded time.Time `json:"dateAdded" xml:"dateAdded"`
}

type Scanner interface {
	Scan(...interface{}) error
}

const (
	timeFormat = "2006-01-02 03:04:05"
)

func (a *ApiKeyType) Get() error {
	err := database.Init()
	if err != nil {
		return err
	}

	stmt, err := database.DB.Prepare(getApiKeyType)
	if err != nil {
		return err
	}
	defer stmt.Close()
	res := stmt.QueryRow(a.ID)
	a, err = ScanKey(res)

	return nil
}

func GetAllApiKeyTypes() ([]ApiKeyType, error) {
	err := database.Init()
	if err != nil {
		return nil, err
	}

	stmt, err := database.DB.Prepare(getAllApiKeyTypes)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	res, err := stmt.Query() //returns *sql.Rows
	if err != nil {
		return nil, err
	}

	var as []ApiKeyType

	for res.Next() {
		a, err := ScanKey(res)
		if err != nil {
			return as, err
		}
		as = append(as, *a)
	}
	defer res.Close()
	return as, err
}

func ScanKey(s Scanner) (*ApiKeyType, error) {
	a := &ApiKeyType{}
	err := s.Scan(&a.ID, &a.Type, &a.DateAdded)
	return a, err
}

func (a *ApiKeyType) Create() error {
	err := database.Init()
	if err != nil {
		return err
	}

	stmt, err := database.DB.Prepare(createApiKeyType)
	if err != nil {
		return err
	}
	defer stmt.Close()
	added := time.Now().Format(timeFormat)
	_, err = stmt.Exec(a.Type, added)
	if err != nil {
		return err
	}

	stmt, err = database.DB.Prepare(getKeyByDateType)
	if err != nil {
		return err
	}

	defer stmt.Close()

	err = stmt.QueryRow(a.Type, added).Scan(&a.ID)
	if err != nil {
		return err
	}
	return nil
}

func (a *ApiKeyType) Delete() error {
	err := database.Init()
	if err != nil {
		return err
	}

	stmt, err := database.DB.Prepare(deleteApiKeyType)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(a.ID)
	if err != nil {
		return err
	}
	return nil
}
