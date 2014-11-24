package brand

import (
	"database/sql"
	"errors"

	"github.com/curt-labs/GoAPI/helpers/database"
	_ "github.com/go-sql-driver/mysql"
)

var (
	getAllBrandsStmt = `select ID, name, code from Brand`
	getBrandStmt     = `select ID, name, code from Brand where ID = ?`
	insertBrandStmt  = `insert into Brand(name, code) values (?,?)`
	updateBrandStmt  = `update Brand set name = ?, code = ? where ID = ?`
	deleteBrandStmt  = `delete from Brand where ID = ?`
)

type Brands []Brand
type Brand struct {
	ID   int    `json:"id" xml:"id,attr"`
	Name string `json:"name" xml:"name,attr"`
	Code string `json:"code" xml:"code,attr"`
}

func GetAllBrands() (brands Brands, err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return
	}
	defer db.Close()

	stmt, err := db.Prepare(getAllBrandsStmt)
	if err != nil {
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return
	}

	var b Brand
	for rows.Next() {
		b = Brand{}
		if err = rows.Scan(&b.ID, &b.Name, &b.Code); err != nil {
			return
		}
		brands = append(brands, b)
	}
	defer rows.Close()

	return
}

func (b *Brand) Get() error {
	if b.ID > 0 {
		db, err := sql.Open("mysql", database.ConnectionString())
		if err != nil {
			return err
		}
		defer db.Close()

		stmt, err := db.Prepare(getBrandStmt)
		if err != nil {
			return err
		}
		defer stmt.Close()

		err = stmt.QueryRow(b.ID).Scan(&b.ID, &b.Name, &b.Code)

		if err == sql.ErrNoRows {
			return errors.New("Invalid Brand ID")
		}

		return err
	}
	return errors.New("Invalid Brand ID")
}

func (b *Brand) Create() error {
	if b.Name == "" {
		return errors.New("Brand must have a name.")
	}
	if b.Code == "" {
		return errors.New("Brand must have a code.")
	}
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(insertBrandStmt)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(b.Name, b.Code)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	b.ID = int(id)
	return err
}

func (b *Brand) Update() error {
	if b.Name == "" {
		return errors.New("Brand must have a name.")
	}
	if b.Code == "" {
		return errors.New("Brand must have a code.")
	}
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(updateBrandStmt)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(b.Name, b.Code, b.ID)
	return err
}

func (b *Brand) Delete() error {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(deleteBrandStmt)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(b.ID)
	return err
}
