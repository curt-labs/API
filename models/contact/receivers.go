package contact

import (
	"database/sql"
	"errors"

	"github.com/curt-labs/GoAPI/helpers/database"
	"github.com/curt-labs/GoAPI/helpers/email"
	_ "github.com/go-sql-driver/mysql"
)

var (
	getAllContactReceiversStmt = `select contactReceiverID, first_name, last_name, email from ContactReceiver`
	getContactReceiverStmt     = `select contactReceiverID, first_name, last_name, email from ContactReceiver where contactReceiverID = ?`
	addContactReceiverStmt     = `insert into ContactReceiver(first_name, last_name, email) values (?,?,?)`
	updateContactReceiverStmt  = `update ContactReceiver set first_name = ?, last_name = ?, email = ? where contactReceiverID = ?`
	deleteContactReceiverStmt  = `delete from ContactReceiver where contactReceiverID = ?`
)

type ContactReceivers []ContactReceiver
type ContactReceiver struct {
	ID        int
	FirstName string
	LastName  string
	Email     string
}

func GetAllContactReceivers() (receivers ContactReceivers, err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return
	}
	defer db.Close()

	stmt, err := db.Prepare(getAllContactReceiversStmt)
	if err != nil {
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return
	}

	for rows.Next() {
		var cr ContactReceiver
		err = rows.Scan(
			&cr.ID,
			&cr.FirstName,
			&cr.LastName,
			&cr.Email,
		)
		if err != nil {
			return
		}
		receivers = append(receivers, cr)
	}
	defer rows.Close()

	return
}

func (cr *ContactReceiver) Get() error {
	if cr.ID > 0 {
		db, err := sql.Open("mysql", database.ConnectionString())
		if err != nil {
			return err
		}
		defer db.Close()

		stmt, err := db.Prepare(getContactReceiverStmt)
		if err != nil {
			return err
		}
		defer stmt.Close()

		err = stmt.QueryRow(cr.ID).Scan(
			&cr.ID,
			&cr.FirstName,
			&cr.LastName,
			&cr.Email,
		)

		return err
	}
	return errors.New("Invalid ContactReceiver ID")
}

func (cr *ContactReceiver) Add() error {
	if !email.IsEmail(cr.Email) {
		return errors.New("Empty or invalid email address.")
	}

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(addContactReceiverStmt)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(cr.FirstName, cr.LastName, cr.Email)
	if err != nil {
		return err
	}

	if id, err := res.LastInsertId(); err != nil {
		return err
	} else {
		cr.ID = int(id)
	}

	return nil
}

func (cr *ContactReceiver) Update() error {
	if cr.ID == 0 {
		return errors.New("Invalid ContactReceiver ID")
	}
	if !email.IsEmail(cr.Email) {
		return errors.New("Empty or invalid email address.")
	}

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(updateContactReceiverStmt)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(cr.FirstName, cr.LastName, cr.Email, cr.ID)

	return err
}

func (cr *ContactReceiver) Delete() error {
	if cr.ID > 0 {
		db, err := sql.Open("mysql", database.ConnectionString())
		if err != nil {
			return err
		}
		defer db.Close()

		stmt, err := db.Prepare(deleteContactReceiverStmt)
		if err != nil {
			return err
		}
		defer stmt.Close()

		_, err = stmt.Exec(cr.ID)

		return err
	}
	return errors.New("Invalid ContactReceiver ID")
}
