package products

import (
	"database/sql"
	"github.com/curt-labs/GoAPI/helpers/database"
	"github.com/curt-labs/GoAPI/models/customer_new/content"
	// "github.com/curt-labs/GoAPI/helpers/redis"
	_ "github.com/go-sql-driver/mysql"
)

type Content struct {
	ID          int                     `json:"id,omitempty" xml:"id,omitempty"`
	Text        string                  `json:"text,omitempty" xml:"text,omitempty"`
	ContentType custcontent.ContentType `json:"contentType,omitempty" xml:"contentType,omitempty"`
	UserID      string                  `json:"userId,omitempty" xml:"userId,omitempty"`
	Deleted     bool                    `json:"deleted,omitempty" xml:"deleted,omitempty"`
}

var (
	createContent = `insert into Content (text, cTypeID, userID, deleted) values (?,?,?,?)`
	deleteContent = `delete from Content where contentID = ? `
)

func (c *Content) Create() error {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	stmt, err := db.Prepare(createContent)
	if err != nil {
		return err
	}
	defer stmt.Close()
	res, err := stmt.Exec(c.Text, c.ContentType.Id, c.UserID, c.Deleted)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	c.ID = int(id)
	return err
}

func (c *Content) Delete() error {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	stmt, err := db.Prepare(deleteContent)
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
