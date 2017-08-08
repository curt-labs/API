package customer

import (
	"github.com/curt-labs/API/helpers/database"
	_ "github.com/go-sql-driver/mysql"
)

var (
	createCustomerBrand     = `insert into CustomerToBrand (cust_id,brandID) values(?,?)`
	deleteCustomerBrand     = `delete from CustomerToBrand where cust_id = ? and brandID = ?`
	deleteAllCustomerBrands = `delete from CustomerToBrand where cust_id = ?`
)

func (c *Customer) CreateCustomerBrand(brandID int) error {
	err := database.Init()
	if err != nil {
		return err
	}

	stmt, err := database.DB.Prepare(createCustomerBrand)
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
	err := database.Init()
	if err != nil {
		return err
	}

	stmt, err := database.DB.Prepare(deleteCustomerBrand)
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

func (c *Customer) DeleteAllCustomerBrands() error {
	err := database.Init()
	if err != nil {
		return err
	}

	stmt, err := database.DB.Prepare(deleteAllCustomerBrands)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(c.Id)
	if err != nil {
		return err
	}
	return err
}
