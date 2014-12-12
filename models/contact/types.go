package contact

import (
	"database/sql"
	"errors"
	"github.com/curt-labs/GoAPI/helpers/apicontext"
	"github.com/curt-labs/GoAPI/helpers/database"
	_ "github.com/go-sql-driver/mysql"
	"strings"
)

var (
	getAllContactTypesStmt = `select contactTypeID, name, showOnWebsite from ContactType
		join apiKeyToBrand as akb on akb.brandID = ContactType.brandID
		join apiKey as ak on ak.id = akb.keyID
		where ak.api_key = ? && (ContactType.BrandID = ? or 0 = ?)`
	getContactTypeStmt    = `select contactTypeID, name, showOnWebsite from ContactType where contactTypeID = ?`
	addContactTypeStmt    = `insert into ContactType(name,showOnWebsite, brandID) values (?,?,?)`
	updateContactTypeStmt = `update ContactType set name = ?, showOnWebsite = ?, brandID = ? where contactTypeID = ?`
	deleteContactTypeStmt = `delete from ContactType where contactTypeID = ?`
	getReceiverByType     = `select cr.contactReceiverID, cr.first_name, cr.last_name, cr.email from ContactReceiver_ContactType as crct 
								left join ContactReceiver as cr on crct.contactReceiverID = cr.contactReceiverID 
								where crct.contactTypeID = ?`
)

type ContactTypes []ContactType
type ContactType struct {
	ID            int    `json:"id" xml:"id"`
	Name          string `json:"name" xml: "name"`
	ShowOnWebsite bool   `json:"show" xml:"show"`
	BrandID       int    `json:"brandId" xml:"brandId"`
}

func GetAllContactTypes(dtx *apicontext.DataContext) (types ContactTypes, err error) {
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

	rows, err := stmt.Query(dtx.APIKey, dtx.BrandID, dtx.BrandID)
	if err != nil {
		return
	}
	for rows.Next() {
		var ct ContactType
		err = rows.Scan(
			&ct.ID,
			&ct.Name,
			&ct.ShowOnWebsite,
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
			&ct.ShowOnWebsite,
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

	res, err := stmt.Exec(ct.Name, ct.ShowOnWebsite, ct.BrandID)
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

func (ct *ContactType) GetReceivers() (crs ContactReceivers, err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return crs, err
	}
	defer db.Close()
	stmt, err := db.Prepare(getReceiverByType)
	if err != nil {
		return crs, err
	}
	defer stmt.Close()
	var cr ContactReceiver
	res, err := stmt.Query(ct.ID)
	if err != nil {
		return crs, err
	}
	for res.Next() {
		err = res.Scan(
			&cr.ID,
			&cr.FirstName,
			&cr.LastName,
			&cr.Email,
		)
		if err != nil {
			return crs, err
		}
		crs = append(crs, cr)
	}
	defer res.Close()
	return crs, err
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

	_, err = stmt.Exec(ct.Name, ct.ShowOnWebsite, ct.BrandID, ct.ID)

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
