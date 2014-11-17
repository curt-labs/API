package cartIntegration

import (
	"database/sql"
	"github.com/curt-labs/GoAPI/helpers/database"
	_ "github.com/go-sql-driver/mysql"
)

type CartIntegration struct {
	ID         int `json:"id,omitempty" xml:"id,omitempty"`
	PartID     int `json:"partId,omitempty" xml:"partId,omitempty"`
	CustPartID int `json:"custPartId,omitempty" xml:"custPartId,omitempty"`
	CustID     int `json:"custId,omitempty" xml:"custId,omitempty"`
}

var (
	getAllCI       = `select referenceID, partID, custPartID, custID from CartIntegration`
	getCIsByPartID = `select referenceID, partID, custPartID, custID from CartIntegration where partID = ?`
	getCIsByCustID = `select referenceID, partID, custPartID, custID from CartIntegration where custID = (select cust_id from Customer where customerID =  ?)`
	getCI          = `select referenceID, partID, custPartID, custID from CartIntegration where custID = ? && partID = ?`
	insertCI       = `insert into CartIntegration (partID, custPartID, custID) values (?,?,?)`
	updateCI       = `update CartIntegration set partID = ?, custPartID = ?, custID = ? where referenceID = ?`
	deleteCI       = `delete from CartIntegration where referenceID = ?`
)

func GetAllCartIntegrations() (cis []CartIntegration, err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return cis, err
	}
	defer db.Close()
	stmt, err := db.Prepare(getAllCI)
	if err != nil {
		return cis, err
	}
	defer stmt.Close()
	res, err := stmt.Query()
	var c CartIntegration
	for res.Next() {
		err = res.Scan(&c.ID, &c.PartID, &c.CustPartID, &c.CustID)
		if err != nil {
			return cis, err
		}
		cis = append(cis, c)
	}
	defer res.Close()
	return cis, err
}

func GetCartIntegrationsByPart(ci CartIntegration) (cis []CartIntegration, err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return cis, err
	}
	defer db.Close()
	stmt, err := db.Prepare(getCIsByPartID)
	if err != nil {
		return cis, err
	}
	defer stmt.Close()
	res, err := stmt.Query(ci.PartID)
	var c CartIntegration
	for res.Next() {
		err = res.Scan(&c.ID, &c.PartID, &c.CustPartID, &c.CustID)
		if err != nil {
			return cis, err
		}
		cis = append(cis, c)
	}
	defer res.Close()
	return cis, err
}

func GetCartIntegrationsByCustomer(ci CartIntegration) (cis []CartIntegration, err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return cis, err
	}
	defer db.Close()
	stmt, err := db.Prepare(getCIsByCustID)
	if err != nil {
		return cis, err
	}
	defer stmt.Close()
	res, err := stmt.Query(ci.CustID)
	var c CartIntegration
	for res.Next() {
		err = res.Scan(&c.ID, &c.PartID, &c.CustPartID, &c.CustID)
		if err != nil {
			return cis, err
		}
		cis = append(cis, c)
	}
	defer res.Close()
	return cis, err
}

func (c *CartIntegration) Get() (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	stmt, err := db.Prepare(getCI)
	if err != nil {
		return err
	}
	defer stmt.Close()
	err = stmt.QueryRow(c.CustID, c.PartID).Scan(&c.ID, &c.PartID, &c.CustPartID, &c.CustID)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	return nil
}

func (c *CartIntegration) Create() (err error) {
	if err := c.Get(); err == nil && c.ID > 0 {
		return c.Update()
	}

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(insertCI)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(c.PartID, c.CustPartID, c.CustID)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	c.ID = int(id)
	return nil
}

func (c *CartIntegration) Update() (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	stmt, err := db.Prepare(updateCI)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(c.PartID, c.CustPartID, c.CustID, c.ID)
	if err != nil {
		return err
	}
	return err
}

func (c *CartIntegration) Delete() (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	stmt, err := db.Prepare(deleteCI)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(c.ID)
	if err != nil {
		return err
	}
	return err
}
