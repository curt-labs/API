package apiKeyType

import (
	"database/sql"
	"github.com/curt-labs/GoAPI/helpers/database"
	_ "github.com/go-sql-driver/mysql"
	"log"
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
	err = stmt.QueryRow(a.ID).Scan(&a.ID, &a.Type, &a.DateAdded)
	if err != nil {
		return
	}
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
	var a ApiKeyType
	for res.Next() {
		err = res.Scan(&a.ID, &a.Type, &a.DateAdded)
		as = append(as, a)
	}
	return as, err
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
	log.Print("A", a.ID)
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
