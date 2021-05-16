package apiKeyType

import (
	"database/sql"
	"time"

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

func (akt *ApiKeyType) Create(db *sql.DB) error {
	keyType := akt.Type

	// We need to manually set this instead of using Database created field because we need it for lookup later.
	// Keeping the code as close as possible for now, we could make this a pure function if we set ApiKeyType.DateAdded
	// and used that in the SQL insert instead of time.Now(), which is hard to mock and adds an implicit dependency for
	// this function.
	added := time.Now().Format(timeFormat)

	result, err := db.Exec(createApiKeyType, keyType, added)
	if err != nil {
		return err
	}

	_, err = result.RowsAffected()
	if err != nil {
		return err
	}

	// get the ID (UUID) for new ApiKeyType
	//
	// In the future we can avoid this query if we either add a real primary key column to the ApiKeyType table
	// and the get the LastInsertID(), or if we generate the UUID ourselves instead of relying on the MySQL version.
	// When we remove this query we can set the timestamp on the MySQL server side instead of needing to pass it in.
	err = db.QueryRow(getKeyByDateType, keyType, added).
		Scan(&akt.ID)
	switch {
	case err == sql.ErrNoRows:
		return err
	case err != nil:
		return err
	default:
		return nil
	}
}

func (akt *ApiKeyType) Delete(db *sql.DB) error {
	id := akt.ID

	result, err := db.Exec(deleteApiKeyType, id)
	if err != nil {
		return err
	}

	_, err = result.RowsAffected()
	if err != nil {
		return err
	}

	return nil
}
