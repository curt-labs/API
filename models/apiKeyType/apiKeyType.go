package apiKeyType

import (
	"database/sql"
	"time"

	"github.com/curt-labs/API/helpers/database"
	_ "github.com/go-sql-driver/mysql"
)

const (
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

const (
	timeFormat = "2006-01-02 03:04:05"
)

// Get populates the ApiKeyType fields with values from the database if a matching ID was found
func (a *ApiKeyType) Get(db *sql.DB) error {
	id := a.ID

	err := db.QueryRow(getApiKeyType, id).
		Scan(&a.ID, &a.Type, &a.DateAdded)
	switch {
	case err == sql.ErrNoRows:
		return err
	case err != nil:
		return err
	default:
		return nil
	}
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
		var a ApiKeyType
		err := res.Scan(&a.ID, &a.Type, &a.DateAdded)
		if err != nil {
			return as, err
		}
		as = append(as, a)
	}
	defer res.Close()
	return as, err
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
