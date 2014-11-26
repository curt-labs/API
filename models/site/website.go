package site

import (
	"database/sql"
	"github.com/curt-labs/GoAPI/helpers/database"
	// "github.com/curt-labs/GoAPI/helpers/redis"
	_ "github.com/go-sql-driver/mysql"
)

type Website struct {
	ID          int      `json:"id,omitempty" xml:"id,omitempty"`
	Url         string   `json:"url,omitempty" xml:"url,omitempty"`
	Description string   `json:"description,omitempty" xml:"description,omitempty"`
	Menus       Menus    `json:"menus,omitempty" xml:"menus,omitempty"`
	Contents    Contents `json:"contents,omitempty" xml:contents,omitempty"`
}
type Websites []Website

var (
	getSite     = `SELECT ID, url, description FROM Website WHERE ID = ?`
	getAllSites = `SELECT ID, url, description FROM Website `
	createSite  = `INSERT INTO Website (url, description) VALUES (?,?)`
	updateSite  = `UPDATE Website SET url = ?, description = ? WHERE ID = ?`
	deleteSite  = `DELETE FROM Website WHERE ID = ?`
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
	return err
}

func (w *Website) GetDetails() (err error) {
	err = w.Get()
	if err != nil {
		return err
	}

	menus, err := GetAllMenus()
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
	return err
}

func (w *Website) Update() (err error) {
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
	tx.Commit()
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
