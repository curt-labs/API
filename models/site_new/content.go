package site_new

import (
	"database/sql"
	"github.com/curt-labs/GoAPI/helpers/database"
	// "github.com/curt-labs/GoAPI/helpers/redis"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

type Content struct {
	Id                    int
	Type                  string
	Title                 string
	CreatedDate           time.Time
	LastModified          time.Time
	MetaTitle             string
	MetaDescription       string
	Keywords              string
	IsPrimary             bool
	Published             bool
	Active                bool
	Slug                  string
	RequireAuthentication bool
	Canonical             string
	ContentRevisions      ContentRevisions
	MenuSort              int
	MenuTitle             string
	MenuLink              string
	ParentId              int
	LinkTarget            bool
}

type Contents []Content

type ContentRevision struct {
	Id          int
	ContentId   int
	Text        string
	CreatedDate time.Time
	Active      bool
}
type ContentRevisions []ContentRevision

const (
	siteContentColumns = "s.contentID, s.content_type, s.page_title, s.createdDate, s.lastModified, s.meta_title, s.meta_description, s.keywords, s.isPrimary, s.published, s.active, s.slug, s.requireAuthentication, s.canonical" //as s
)

var (
	getLatestRevision      = `SELECT revisionID, content_text, createdOn, active FROM SiteContentRevision AS scr WHERE scr.contentID = ? ORDER BY createdOn DESC LIMIT 1`
	getContent             = `SELECT ` + siteContentColumns + ` FROM SiteContent AS s WHERE s.contentID = ? `
	getAllContent          = `SELECT ` + siteContentColumns + ` FROM SiteContent AS s  `
	getContentRevisions    = `SELECT revisionID, content_text, createdOn, active FROM SiteContentRevision AS scr WHERE scr.contentID = ? `
	getAllContentRevisions = `SELECT revisionID, content_text, createdOn, active FROM SiteContentRevision AS scr `
	getContentRevision     = `SELECT revisionID, content_text, createdOn, active FROM SiteContentRevision AS scr WHERE revisionID = ?`
	getContentBySlug       = `SELECT ` + siteContentColumns + ` FROM SiteContent AS s WHERE s.slug = ? `

	//operations
	createRevision = `INSERT INTO SiteContentRevision (contentID, content_text, createdOn, active) VALUES (?,?,?,?)`
	createContent  = `INSERT INTO SiteContent 
						(content_type, page_title, createdDate, meta_title, meta_description, keywords, isPrimary, published, active, slug, requireAuthentication, canonical)
						VALUES (?,?,?,?,?,?,?,?,?,?,?,?)`
	updateRevision = `UPDATE SiteContentRevision SET contentID = ?, content_text = ?, active = ? WHERE revisionID = ?`
	updateContent  = `UPDATE SiteContent SET 
					content_type = ?, page_title = ?,  meta_title = ?, meta_description = ?, keywords = ?, isPrimary = ?, published = ?, active = ?, slug = ?, requireAuthentication = ?, canonical  = ?
					WHERE contentID = ?`

	deleteRevision                   = `DELETE FROM SiteContentRevision WHERE revisionID = ?`
	deleteContent                    = `DELETE FROM SiteContent WHERE contentID = ?`
	deleteRevisionbyContentID        = `DELETE FROM SiteContentRevision WHERE contentID = ?`
	deleteMenuSiteContentByContentId = `DELETE FROM Menu_SiteContent WHERE contentID = ?`
)

//Fetch content by id
func (c *Content) Get() (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	stmt, err := db.Prepare(getContent)
	if err != nil {
		return err
	}
	defer stmt.Close()

	var cType, title, mTitle, mDesc, slug, canon *string
	err = stmt.QueryRow(c.Id).Scan(
		&c.Id,
		&cType,
		&title,
		&c.CreatedDate,
		&c.LastModified,
		&mTitle,
		&mDesc,
		&c.Keywords,
		&c.IsPrimary,
		&c.Published,
		&c.Active,
		&slug,
		&c.RequireAuthentication,
		&canon,
	)
	if err != nil {
		return err
	}

	if cType != nil {
		c.Type = *cType
	}
	if title != nil {
		c.Title = *title
	}
	if mTitle != nil {
		c.MetaTitle = *mTitle
	}
	if mDesc != nil {
		c.MetaDescription = *mDesc
	}
	if slug != nil {
		c.Slug = *slug
	}
	if canon != nil {
		c.Canonical = *canon
	}
	return err
}

//Fetch content by slug
func (c *Content) GetbySlug() (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	stmt, err := db.Prepare(getContentBySlug)
	if err != nil {
		return err
	}
	defer stmt.Close()

	var cType, title, mTitle, mDesc, slug, canon *string
	err = stmt.QueryRow(c.Slug).Scan(
		&c.Id,
		&cType,
		&title,
		&c.CreatedDate,
		&c.LastModified,
		&mTitle,
		&mDesc,
		&c.Keywords,
		&c.IsPrimary,
		&c.Published,
		&c.Active,
		&slug,
		&c.RequireAuthentication,
		&canon,
	)
	if err != nil {
		return err
	}

	if cType != nil {
		c.Type = *cType
	}
	if title != nil {
		c.Title = *title
	}
	if mTitle != nil {
		c.MetaTitle = *mTitle
	}
	if mDesc != nil {
		c.MetaDescription = *mDesc
	}
	if slug != nil {
		c.Slug = *slug
	}
	if canon != nil {
		c.Canonical = *canon
	}
	return err
}

//Fetch a great many contents
func GetAllContents() (cs Contents, err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return cs, err
	}
	defer db.Close()
	stmt, err := db.Prepare(getAllContent)
	if err != nil {
		return cs, err
	}
	defer stmt.Close()

	var cType, title, mTitle, mDesc, slug, canon *string
	var c Content
	res, err := stmt.Query()
	for res.Next() {
		err = res.Scan(
			&c.Id,
			&cType,
			&title,
			&c.CreatedDate,
			&c.LastModified,
			&mTitle,
			&mDesc,
			&c.Keywords,
			&c.IsPrimary,
			&c.Published,
			&c.Active,
			&slug,
			&c.RequireAuthentication,
			&canon,
		)
		if err != nil {
			return cs, err
		}

		if cType != nil {
			c.Type = *cType
		}
		if title != nil {
			c.Title = *title
		}
		if mTitle != nil {
			c.MetaTitle = *mTitle
		}
		if mDesc != nil {
			c.MetaDescription = *mDesc
		}
		if slug != nil {
			c.Slug = *slug
		}
		if canon != nil {
			c.Canonical = *canon
		}
		cs = append(cs, c)
	}
	return cs, err
}

//Fetch a content's most recent revision
func (c *Content) GetLatestRevision() (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	stmt, err := db.Prepare(getLatestRevision)
	if err != nil {
		return err
	}
	defer stmt.Close()

	var rev ContentRevision
	err = stmt.QueryRow(c.Id).Scan(
		&rev.Id,
		&rev.Text,
		&rev.CreatedDate,
		&rev.Active,
	)
	c.ContentRevisions = append(c.ContentRevisions, rev)
	return err
}

//Fetch all of thine content's revisions
func (c *Content) GetContentRevisions() (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	stmt, err := db.Prepare(getContentRevisions)
	if err != nil {
		return err
	}
	defer stmt.Close()

	var rev ContentRevision
	res, err := stmt.Query(c.Id)
	for res.Next() {
		err = res.Scan(
			&rev.Id,
			&rev.Text,
			&rev.CreatedDate,
			&rev.Active,
		)
		if err != nil {
			return err
		}
		c.ContentRevisions = append(c.ContentRevisions, rev)
	}
	return err
}

//Fetch a single revision by Id
func (rev *ContentRevision) Get() (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	stmt, err := db.Prepare(getContentRevision)
	if err != nil {
		return err
	}
	defer stmt.Close()

	err = stmt.QueryRow(rev.Id).Scan(
		&rev.Id,
		&rev.Text,
		&rev.CreatedDate,
		&rev.Active,
	)
	if err != nil {
		return err
	}
	return err
}

//Fetch a great many revisions
func GetAllContentRevisions() (cr ContentRevisions, err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return cr, err
	}
	defer db.Close()
	stmt, err := db.Prepare(getAllContentRevisions)
	if err != nil {
		return cr, err
	}
	defer stmt.Close()

	var rev ContentRevision
	res, err := stmt.Query()
	for res.Next() {
		err = res.Scan(
			&rev.Id,
			&rev.Text,
			&rev.CreatedDate,
			&rev.Active,
		)
		if err != nil {
			return cr, err
		}
		cr = append(cr, rev)
	}
	return cr, err
}

//creatin' content
func (c *Content) Create() (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()
	stmt, err := tx.Prepare(createContent)
	if err != nil {
		return err
	}
	defer stmt.Close()
	c.CreatedDate = time.Now()
	res, err := stmt.Exec(
		c.Type,
		c.Title,
		c.CreatedDate,
		c.MetaTitle,
		c.MetaDescription,
		c.Keywords,
		c.IsPrimary,
		c.Published,
		c.Active,
		c.Slug,
		c.RequireAuthentication,
		c.Canonical,
	)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	c.Id = int(id)
	return err
}

//updatin' content
func (c *Content) Update() (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()
	stmt, err := tx.Prepare(updateContent)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		c.Type,
		c.Title,
		c.MetaTitle,
		c.MetaDescription,
		c.Keywords,
		c.IsPrimary,
		c.Published,
		c.Active,
		c.Slug,
		c.RequireAuthentication,
		c.Canonical,
		c.Id,
	)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	return err
}

//deletin' content, brings joined revisions and menu join with
func (c *Content) Delete() (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()

	//adios revisions
	stmt, err := tx.Prepare(deleteRevisionbyContentID)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(c.Id)
	if err != nil {
		tx.Rollback()
		return err
	}

	//adios menu join
	stmt, err = tx.Prepare(deleteMenuSiteContentByContentId)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(c.Id)
	if err != nil {
		tx.Rollback()
		return err
	}

	//adios content
	stmt, err = tx.Prepare(deleteContent)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(c.Id)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	return err
}

//creatin' a revision, requires content to exist
func (rev *ContentRevision) Create() (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()
	stmt, err := tx.Prepare(createRevision)
	if err != nil {
		return err
	}
	defer stmt.Close()
	rev.CreatedDate = time.Now()
	res, err := stmt.Exec(rev.ContentId, rev.Text, rev.CreatedDate, rev.Active)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	rev.Id = int(id)
	return err
}

//updatin' a revision, requires content to exisi
func (rev *ContentRevision) Update() (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()
	stmt, err := tx.Prepare(updateRevision)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(rev.ContentId, rev.Text, rev.Active, rev.Id)
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

//deletin' a revision
func (rev *ContentRevision) Delete() (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()
	stmt, err := tx.Prepare(deleteRevision)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(rev.Id)
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
