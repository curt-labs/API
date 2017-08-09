package contact

import (
	"errors"

	"github.com/curt-labs/API/helpers/database"
	"github.com/curt-labs/API/helpers/email"
	_ "github.com/go-sql-driver/mysql"
)

var (
	getAllContactReceiversStmt              = `select contactReceiverID, first_name, last_name, email from ContactReceiver`
	getContactReceiverStmt                  = `select contactReceiverID, first_name, last_name, email from ContactReceiver where contactReceiverID = ?`
	addContactReceiverStmt                  = `insert into ContactReceiver(first_name, last_name, email) values (?,?,?)`
	updateContactReceiverStmt               = `update ContactReceiver set first_name = ?, last_name = ?, email = ? where contactReceiverID = ?`
	deleteContactReceiverStmt               = `delete from ContactReceiver where contactReceiverID = ?`
	createReceiverContactTypeJoin           = `insert into ContactReceiver_ContactType (ContactReceiverID, ContactTypeID) values (?,?)`
	deleteReceiverContactTypeJoin           = `delete from ContactReceiver_ContactType where ContactReceiverID = ? and  ContactTypeID = ?`
	deleteReceiverContactTypeJoinByReceiver = `delete from ContactReceiver_ContactType where ContactReceiverID = ?`
	getContactTypesByReceiver               = `select crct.contactTypeID, ct.name, ct.showOnWebsite, ct.brandID from ContactReceiver_ContactType as crct
												left join ContactType as ct on crct.ContactTypeID = ct.contactTypeID where crct.contactReceiverID = ?`
)

type ContactReceivers []ContactReceiver
type ContactReceiver struct {
	ID           int
	FirstName    string
	LastName     string
	Email        string
	ContactTypes ContactTypes
}

func GetAllContactReceivers() (ContactReceivers, error) {
	err := database.Init()
	if err != nil {
		return ContactReceivers{}, err
	}

	stmt, err := database.DB.Prepare(getAllContactReceiversStmt)
	if err != nil {
		return ContactReceivers{}, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return ContactReceivers{}, err
	}

	var receivers ContactReceivers

	for rows.Next() {
		var cr ContactReceiver
		err = rows.Scan(
			&cr.ID,
			&cr.FirstName,
			&cr.LastName,
			&cr.Email,
		)
		if err != nil {
			return receivers, err
		}
		err = cr.GetContactTypes()
		if err != nil {
			return receivers, err
		}
		receivers = append(receivers, cr)
	}
	defer rows.Close()

	return receivers, nil
}

func (cr *ContactReceiver) Get() error {
	if cr.ID > 0 {
		err := database.Init()
		if err != nil {
			return err
		}

		stmt, err := database.DB.Prepare(getContactReceiverStmt)
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
		err = cr.GetContactTypes()

		return err
	}
	return errors.New("Invalid ContactReceiver ID")
}

func (cr *ContactReceiver) Add() error {
	if !email.IsEmail(cr.Email) {
		return errors.New("Empty or invalid email address.")
	}

	err := database.Init()
	if err != nil {
		return err
	}

	stmt, err := database.DB.Prepare(addContactReceiverStmt)
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
	//add contact types
	if len(cr.ContactTypes) > 0 {
		for _, ct := range cr.ContactTypes {
			err = cr.CreateTypeJoin(ct)
			if err != nil {
				return err
			}
		}
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

	err := database.Init()
	if err != nil {
		return err
	}

	stmt, err := database.DB.Prepare(updateContactReceiverStmt)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(cr.FirstName, cr.LastName, cr.Email, cr.ID)

	//update type joins
	if len(cr.ContactTypes) > 0 {
		err = cr.DeleteTypeJoinByReceiver()
		if err != nil {
			return err
		}
		for _, ct := range cr.ContactTypes {
			err = cr.CreateTypeJoin(ct)
			if err != nil {
				return err
			}
		}
	}

	return err
}

func (cr *ContactReceiver) Delete() error {
	if cr.ID > 0 {
		err := database.Init()
		if err != nil {
			return err
		}

		stmt, err := database.DB.Prepare(deleteContactReceiverStmt)
		if err != nil {
			return err
		}
		defer stmt.Close()

		_, err = stmt.Exec(cr.ID)

		//delete receiver-type join
		err = cr.DeleteTypeJoinByReceiver()

		return err
	}
	return errors.New("Invalid ContactReceiver ID")
}

//get a contact receiver's contact types
func (cr *ContactReceiver) GetContactTypes() error {
	err := database.Init()
	if err != nil {
		return err
	}

	stmt, err := database.DB.Prepare(getContactTypesByReceiver)
	if err != nil {
		return err
	}
	defer stmt.Close()
	var ct ContactType
	res, err := stmt.Query(cr.ID)
	for res.Next() {
		err = res.Scan(&ct.ID, &ct.Name, &ct.ShowOnWebsite, &ct.BrandID)
		if err != nil {
			return err
		}
		cr.ContactTypes = append(cr.ContactTypes, ct)
	}
	defer res.Close()
	return err
}

func (cr *ContactReceiver) CreateTypeJoin(ct ContactType) error {
	err := database.Init()
	if err != nil {
		return err
	}

	stmt, err := database.DB.Prepare(createReceiverContactTypeJoin)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(cr.ID, ct.ID)
	if err != nil {
		return err
	}
	return nil
}

//delete joins for a receiver-type pair
func (cr *ContactReceiver) DeleteTypeJoin(ct ContactType) error {
	err := database.Init()
	if err != nil {
		return err
	}

	stmt, err := database.DB.Prepare(deleteReceiverContactTypeJoin)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(cr.ID, ct.ID)
	if err != nil {
		return err
	}
	return nil
}

//delete all type-receiver joins for a receiver
func (cr *ContactReceiver) DeleteTypeJoinByReceiver() error {
	err := database.Init()
	if err != nil {
		return err
	}

	stmt, err := database.DB.Prepare(deleteReceiverContactTypeJoinByReceiver)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(cr.ID)
	if err != nil {
		return err
	}
	return nil
}
