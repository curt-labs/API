package contact

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/curt-labs/GoAPI/helpers/database"
	_ "github.com/go-sql-driver/mysql"
)

var (
	getAllContactTypesStmt = `select contactTypeID, name from ContactType`
	getContactTypeStmt     = `select contactTypeID, name from ContactType where contactTypeID = ?`
	addContactTypeStmt     = `insert into ContactType(name) values (?)`
	updateContactTypeStmt  = `update ContactType set name = ? where contactTypeID = ?`
	deleteContactTypeStmt  = `delete from ContactType where contactTypeID = ?`
)

type ContactTypes []ContactType
type ContactType struct {
	ID   int
	Name string
}

func GetAllContactTypes() (types ContactTypes, err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return
	}
	defer db.Close()

	stmt, err := db.Prepare(getAllContactTypesStmt)
	if err != nil {
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return
	}

	for rows.Next() {
		var ct ContactType
		err = rows.Scan(
			&ct.ID,
			&ct.Name,
		)
		if err != nil {
			return
		}
		types = append(types, ct)
	}
	defer rows.Close()

	return
}

func (ct *ContactType) Get() error {
	if ct.ID > 0 {
		db, err := sql.Open("mysql", database.ConnectionString())
		if err != nil {
			return err
		}
		defer db.Close()

		stmt, err := db.Prepare(getContactTypeStmt)
		if err != nil {
			return err
		}
		defer stmt.Close()

		err = stmt.QueryRow(ct.ID).Scan(
			&ct.ID,
			&ct.Name,
		)
		return err
	}
	return errors.New("Invalid ContactType ID")
}

func (ct *ContactType) Add() error {
	if strings.TrimSpace(ct.Name) == "" {
		return errors.New("Invalid contact name.")
	}

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(addContactTypeStmt)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(ct.Name)
	if err != nil {
		return err
	}

	if id, err := res.LastInsertId(); err != nil {
		return err
	} else {
		ct.ID = int(id)
	}

	return nil
}

func (ct *ContactType) Update() error {
	if ct.ID == 0 {
		return errors.New("Invalid ContactType ID")
	}
	if strings.TrimSpace(ct.Name) == "" {
		return errors.New("Invalid contact name.")
	}

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(updateContactTypeStmt)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(ct.Name, ct.ID)

	return err
}

func (ct *ContactType) Delete() error {
	if ct.ID > 0 {
		db, err := sql.Open("mysql", database.ConnectionString())
		if err != nil {
			return err
		}
		defer db.Close()

		stmt, err := db.Prepare(deleteContactTypeStmt)
		if err != nil {
			return err
		}
		defer stmt.Close()

		_, err = stmt.Exec(ct.ID)

		return err
	}
	return errors.New("Invalid ContactType ID")
}
