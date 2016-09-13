package site

import (
	"database/sql"

	"github.com/curt-labs/API/helpers/apicontext"
	"github.com/curt-labs/API/helpers/database"
	// "github.com/curt-labs/API/helpers/redis"
	_ "github.com/go-sql-driver/mysql"
)

type Website struct {
	ID          int      `json:"id,omitempty" xml:"id,omitempty"`
	Url         string   `json:"url,omitempty" xml:"url,omitempty"`
	Description string   `json:"description,omitempty" xml:"description,omitempty"`
	Menus       Menus    `json:"menus,omitempty" xml:"menus,omitempty"`
	Contents    Contents `json:"contents,omitempty" xml:contents,omitempty"`
	BrandIDs    []int    `json:"brandId,omitempty" xml:brandId,omitempty"`
}
type Websites []Website

var (
	getSite         = `SELECT ID, url, description FROM Website WHERE ID = ?`
	getAllSites     = `SELECT ID, url, description FROM Website `
	createSite      = `INSERT INTO Website (url, description) VALUES (?,?)`
	updateSite      = `UPDATE Website SET url = ?, description = ? WHERE ID = ?`
	deleteSite      = `DELETE FROM Website WHERE ID = ?`
	joinToBrand     = `insert into WebsiteToBrand (WebsiteID, brandID) values (?,?)`
	deleteBrandJoin = `delete from WebsiteToBrand where WebsiteID = ? and brandID = ?`
	getBrands       = `select brandID from WebsiteToBrand where WebsiteID = ?`
)

func (w *Website) Get() (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	stmt, err := db.Prepare(getSite)
	if err != nil {
		return err
	}
	defer stmt.Close()
	var url, desc *string
	err = stmt.QueryRow(w.ID).Scan(
		&w.ID,
		&url,
		&desc,
	)
	if err != nil {
		return err
	}
	if url != nil {
		w.Url = *url
	}
	if desc != nil {
		w.Description = *desc
	}
	//get brands
	stmt, err = db.Prepare(getBrands)
	if err != nil {
		return err
	}
	defer stmt.Close()
	res, err := stmt.Query(w.ID)
	if err != nil {
		return err
	}
	var brandId int
	for res.Next() {
		err = res.Scan(&brandId)
		if err != nil {
			return err
		}
		w.BrandIDs = append(w.BrandIDs, brandId)
	}
	return err
}

func (w *Website) GetDetails(dtx *apicontext.DataContext) (err error) {
	err = w.Get()
	if err != nil {
		return err
	}

	menus, err := GetAllMenus(dtx)
	menuMap := menus.ToMap()

	for _, menu := range menuMap {

		if menu.WebsiteId == w.ID {
			err = menu.GetContents()
			w.Menus = append(w.Menus, menu)
		}
	}

	return err
}

func GetAllWebsites() (ws Websites, err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return ws, err
	}
	defer db.Close()
	stmt, err := db.Prepare(getAllSites)
	if err != nil {
		return ws, err
	}
	defer stmt.Close()
	res, err := stmt.Query()
	if err != nil {
		return ws, err
	}
	var w Website
	var url, desc *string
	for res.Next() {
		err = res.Scan(
			&w.ID,
			&url,
			&desc,
		)
		if err != nil {
			return ws, err
		}
		if url != nil {
			w.Url = *url
		}
		if desc != nil {
			w.Description = *desc
		}
		ws = append(ws, w)
	}
	defer res.Close()
	return ws, err
}

func (w *Website) Create() (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()
	stmt, err := tx.Prepare(createSite)
	if err != nil {
		return err
	}
	defer stmt.Close()
	res, err := stmt.Exec(w.Url, w.Description)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	id, err := res.LastInsertId()
	w.ID = int(id)
	if err != nil {
		return err
	}
	err = w.JoinToBrand()
	if err != nil {
		return err
	}
	return err
}

func (w *Website) Update() error {
	var err error
	for _, brandId := range w.BrandIDs {
		err = w.DeleteBrandJoin(brandId)
		if err != nil {
			return err
		}
	}
	err = w.JoinToBrand()
	if err != nil {
		return err
	}
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()
	stmt, err := tx.Prepare(updateSite)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(w.Url, w.Description, w.ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return err
}

func (w *Website) Delete() (err error) {
	for _, brandId := range w.BrandIDs {
		err = w.DeleteBrandJoin(brandId)
		if err != nil {
			return err
		}
	}

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()
	stmt, err := tx.Prepare(deleteSite)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(w.ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	return err
}

func (w *Website) JoinToBrand() error {
	var err error
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(joinToBrand)
	if err != nil {
		return err
	}
	defer stmt.Close()
	for _, brandId := range w.BrandIDs {
		_, err = stmt.Exec(w.ID, brandId)
		if err != nil {
			return err
		}
	}
	return err
}

func (w *Website) DeleteBrandJoin(brandId int) error {
	var err error
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(deleteBrandJoin)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(w.ID, brandId)
	if err != nil {
		return err
	}
	return err
}

//mapping
func (c Contents) ToMap() map[int]Content {
	theMap := make(map[int]Content)
	for _, v := range c {
		theMap[v.Id] = v
	}
	return theMap
}

func (m Menus) ToMap() map[int]Menu {
	theMap := make(map[int]Menu)
	for _, v := range m {
		theMap[v.Id] = v
	}
	return theMap
}
