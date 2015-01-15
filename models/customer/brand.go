package customer

import (
	"github.com/curt-labs/GoAPI/helpers/database"
	_ "github.com/go-sql-driver/mysql"

	"database/sql"
)

var (
	createCustomerBrand = `insert into CustomerToBrand (cust_id,brandID) values(?,?)`
	deleteCustomerBrand = `delete from CustomerToBrand where cust_id = ? and brandID = ?`
)

func (c *Customer) CreateCustomerBrand(brandID int) error {
	var err error
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(createCustomerBrand)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(c.Id, brandID)
	if err != nil {
		return err
	}
	return err
}

func (c *Customer) DeleteCustomerBrand(brandID int) error {
	var err error
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(deleteCustomerBrand)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(c.Id, brandID)
	if err != nil {
		return err
	}
	return err
}
