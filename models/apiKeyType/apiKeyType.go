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
func (akt *ApiKeyType) Get(db *sql.DB) error {
	id := akt.ID

	err := db.QueryRow(getApiKeyType, id).
		Scan(&akt.ID, &akt.Type, &akt.DateAdded)
	switch {
	case err == sql.ErrNoRows:
		return err
	case err != nil:
		return err
	default:
		return nil
	}
}

// GetAllApiKeyTypes fetches all API key types
func GetAllApiKeyTypes(db *sql.DB) ([]ApiKeyType, error) {
	rows, err := db.Query(getAllApiKeyTypes)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	akts := make([]ApiKeyType, 0)
	for rows.Next() {
		var akt ApiKeyType
		if err := rows.Scan(&akt.ID, &akt.Type, &akt.DateAdded); err != nil {
			// stop fetching keys on first error
			return akts, err
		}
		akts = append(akts, akt)
	}

	rerr := rows.Close()
	if rerr != nil {
		return akts, rerr
	}

	if err := rows.Err(); err != nil {
		return akts, err
	}

	return akts, nil
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
