package salesrep

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/curt-labs/GoAPI/helpers/database"
	_ "github.com/go-sql-driver/mysql"
)

var (
	getAllSalesReps = `select salesRepID, name, code, (select count(cust_id) from Customer where Customer.salesRepID = SalesRepresentative.salesRepID) AS customercount from SalesRepresentative`
	getSalesRep     = `select salesRepID, name, code, (select count(cust_id) from Customer where Customer.salesRepID = SalesRepresentative.salesRepID) AS customercount from SalesRepresentative where salesRepID = ?`
	updateSalesRep  = `update SalesRepresentative set name = ?, code = ? where salesRepID = ?`
	addSalesRep     = `insert into SalesRepresentative (name,code) values (?,?)`
	deleteSalesRep  = `delete from SalesRepresentative where salesRepID = ?`
)

type SalesReps []SalesRep
type SalesRep struct {
	ID            int
	Name          string
	Code          string
	CustomerCount int
}

func GetAllSalesReps() (reps SalesReps, err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return
	}
	defer db.Close()

	stmt, err := db.Prepare(getAllSalesReps)
	if err != nil {
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return
	}

	for rows.Next() {
		var rep SalesRep
		if err = rows.Scan(&rep.ID, &rep.Name, &rep.Code, &rep.CustomerCount); err == nil {
			reps = append(reps, rep)
		}
	}
	defer rows.Close()

	return
}

func (r *SalesRep) Get() error {
	if r.ID == 0 {
		return errors.New("Invalid SalesRep ID")
	}

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(getSalesRep)
	if err != nil {
		return err
	}
	defer stmt.Close()

	var rep SalesRep
	row := stmt.QueryRow(r.ID)
	if err = row.Scan(&rep.ID, &rep.Name, &rep.Code, &rep.CustomerCount); err != nil {
		return err
	}

	r.ID = rep.ID
	r.Name = rep.Name
	r.Code = rep.Code
	r.CustomerCount = rep.CustomerCount

	return nil
}

func (rep *SalesRep) Add() error {
	if len(strings.TrimSpace(rep.Name)) == 0 {
		return errors.New("SalesRep must have a name")
	}

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(addSalesRep)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(rep.Name, rep.Code)
	if err != nil {
		return err
	}

	if id, err := res.LastInsertId(); err != nil {
		return err
	} else {
		rep.ID = int(id)
	}

	return nil
}

func (rep *SalesRep) Update() error {
	if rep.ID == 0 {
		return errors.New("Invalid SalesRep ID")
	}

	if len(strings.TrimSpace(rep.Name)) == 0 {
		return errors.New("SalesRep must have a name")
	}

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(updateSalesRep)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err := stmt.Exec(rep.Name, rep.Code, rep.ID); err != nil {
		return err
	}

	return nil
}

func (rep *SalesRep) Delete() error {
	if rep.ID == 0 {
		return errors.New("Invalid SalesRep ID")
	}

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(deleteSalesRep)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err = stmt.Exec(rep.ID); err != nil {
		return err
	}

	return nil
}
