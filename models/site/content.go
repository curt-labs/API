package site

import (
	"github.com/curt-labs/API/helpers/apicontext"
	"github.com/curt-labs/API/helpers/database"
	_ "github.com/go-sql-driver/mysql"

	"database/sql"
	"time"
)

type Content struct {
	Id                    int              `json:"id,omitempty" xml:"id,omitempty"`
	Type                  string           `json:"type,omitempty" xml:"type,omitempty"`
	Title                 string           `json:"title,omitempty" xml:"title,omitempty"`
	CreatedDate           time.Time        `json:"createdDate,omitempty" xml:"createdDate,omitempty"`
	LastModified          time.Time        `json:"lastModified,omitempty" xml:"lastModified,omitempty"`
	MetaTitle             string           `json:"metaTitle,omitempty" xml:"metaTitle,omitempty"`
	MetaDescription       string           `json:"metaDescription,omitempty" xml:"metaDescription,omitempty"`
	Keywords              string           `json:"keywords,omitempty" xml:"keywords,omitempty"`
	IsPrimary             bool             `json:"isPrimary,omitempty" xml:"isPrimary,omitempty"`
	Published             bool             `json:"published,omitempty" xml:"published,omitempty"`
	Active                bool             `json:"active,omitempty" xml:"active,omitempty"`
	Slug                  string           `json:"slug,omitempty" xml:"slug,omitempty"`
	RequireAuthentication bool             `json:"requireAuthentication,omitempty" xml:"requireAuthentication,omitempty"`
	Canonical             string           `json:"canonical,omitempty" xml:"canonical,omitempty"`
	WebsiteId             int              `json:"websiteId,omitempty" xml:"websiteId,omitempty"`
	ContentRevisions      ContentRevisions `json:"contentRevisions,omitempty" xml:"contentRevisions,omitempty"`
	MenuSort              int              `json:"menuSort,omitempty" xml:"menuSort,omitempty"`
	MenuTitle             string           `json:"menuTitle,omitempty" xml:"menuTitle,omitempty"`
	MenuLink              string           `json:"menuLink,omitempty" xml:"menuLink,omitempty"`
	ParentId              int              `json:"parentId,omitempty" xml:"parentId,omitempty"`
	LinkTarget            bool             `json:"linkTarget,omitempty" xml:"v,omitempty"`
}

type Contents []Content

type ContentRevision struct {
	Id          int       `json:"id,omitempty" xml:"id,omitempty"`
	ContentId   int       `json:"contentId,omitempty" xml:"contentId,omitempty"`
	Text        string    `json:"text,omitempty" xml:"text,omitempty"`
	CreatedDate time.Time `json:"createdDate,omitempty" xml:"createdDate,omitempty"`
	Active      bool      `json:"active,omitempty" xml:"active,omitempty"`
}
type ContentRevisions []ContentRevision

const (
	siteContentColumns = "s.contentID, s.content_type, s.page_title, s.createdDate, s.lastModified, s.meta_title, s.meta_description, s.keywords, s.isPrimary, s.published, s.active, s.slug, s.requireAuthentication, s.canonical, s.websiteID" //as s
)

var (
	getLatestRevision = `SELECT revisionID, content_text, createdOn, active FROM SiteContentRevision AS scr WHERE scr.contentID = ? ORDER BY createdOn DESC LIMIT 1`
	getContent        = `SELECT ` + siteContentColumns + ` FROM SiteContent AS s
								Join WebsiteToBrand as wub on wub.WebsiteID = s.websiteID
								Join ApiKeyToBrand as akb on akb.brandID = wub.brandID
								Join ApiKey as ak on akb.keyID = ak.id
								where s.contentID = ? && (ak.api_key = ? && (wub.brandID = ? OR 0=?))`
	getAllContent = `SELECT ` + siteContentColumns + ` FROM SiteContent AS s  
								Join WebsiteToBrand as wub on wub.WebsiteID = s.websiteID
								Join ApiKeyToBrand as akb on akb.brandID = wub.brandID
								Join ApiKey as ak on akb.keyID = ak.id
								where s.active = true && s.published = true && (ak.api_key = ? && (wub.brandID = ? OR 0=?)  && (wub.WebsiteID = ? OR 0=?))`
	getContentRevisions    = `SELECT revisionID, content_text, createdOn, active FROM SiteContentRevision AS scr WHERE scr.contentID = ? `
	getAllContentRevisions = `SELECT revisionID, content_text, createdOn, active FROM SiteContentRevision AS scr `
	getContentRevision     = `SELECT revisionID, content_text, createdOn, active FROM SiteContentRevision AS scr WHERE revisionID = ?`
	getContentBySlug       = `SELECT ` + siteContentColumns + ` FROM SiteContent AS s 
								Join WebsiteToBrand as wub on wub.WebsiteID = s.websiteID
								Join ApiKeyToBrand as akb on akb.brandID = wub.brandID
								Join ApiKey as ak on akb.keyID = ak.id
								WHERE s.slug = ? && (ak.api_key = ? && (wub.brandID = ? OR 0=?))`

	//operations
	createRevision = `INSERT INTO SiteContentRevision (contentID, content_text, createdOn, active) VALUES (?,?,?,?)`
	createContent  = `INSERT INTO SiteContent
						(content_type, page_title, createdDate, meta_title, meta_description, keywords, isPrimary, published, active, slug, requireAuthentication, canonical, websiteID)
						VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?)`
	updateRevision = `UPDATE SiteContentRevision SET contentID = ?, content_text = ?, active = ? WHERE revisionID = ?`
	updateContent  = `UPDATE SiteContent SET
					content_type = ?, page_title = ?,  meta_title = ?, meta_description = ?, keywords = ?, isPrimary = ?, published = ?, active = ?, slug = ?, requireAuthentication = ?, canonical  = ?, websiteID = ?
					WHERE contentID = ?`

	deleteRevision                   = `DELETE FROM SiteContentRevision WHERE revisionID = ?`
	deleteContent                    = `DELETE FROM SiteContent WHERE contentID = ?`
	deleteRevisionbyContentID        = `DELETE FROM SiteContentRevision WHERE contentID = ?`
	deleteMenuSiteContentByContentId = `DELETE FROM Menu_SiteContent WHERE contentID = ?`
)

//Fetch content by id
func (c *Content) Get(dtx *apicontext.DataContext) (err error) {
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
	err = stmt.QueryRow(c.Id, dtx.APIKey, dtx.BrandID, dtx.BrandID).Scan(
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
		&c.WebsiteId,
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
	//get latest revision
	err = c.GetLatestRevision()
	return err
}

//Fetch content by slug
func (c *Content) GetBySlug(dtx *apicontext.DataContext) (err error) {
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
	err = stmt.QueryRow(c.Slug, dtx.APIKey, dtx.BrandID, dtx.BrandID).Scan(
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
		&c.WebsiteId,
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
	//get latest revision
	err = c.GetLatestRevision()
	if err == sql.ErrNoRows {
		err = nil
	}
	return err
}

//Fetch a great many contents
func GetAllContents(dtx *apicontext.DataContext, siteID int) (cs Contents, err error) {
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

	var cType, title, mTitle, mDesc, slug, canon, keywords *string
	var c Content
	res, err := stmt.Query(dtx.APIKey, dtx.BrandID, dtx.BrandID, siteID, siteID)
	for res.Next() {
		err = res.Scan(
			&c.Id,
			&cType,
			&title,
			&c.CreatedDate,
			&c.LastModified,
			&mTitle,
			&mDesc,
			&keywords,
			&c.IsPrimary,
			&c.Published,
			&c.Active,
			&slug,
			&c.RequireAuthentication,
			&canon,
			&c.WebsiteId,
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
		if keywords != nil {
			c.Keywords = *keywords
		}
		if slug != nil {
			c.Slug = *slug
		}
		if canon != nil {
			c.Canonical = *canon
		}
		err = c.GetLatestRevision()
		if err != nil {
			return cs, err
		}

		cs = append(cs, c)
	}
	defer res.Close()
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

	c.ContentRevisions = nil //refresh contentrevision array
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
	defer res.Close()
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
	defer res.Close()
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
		c.WebsiteId,
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
	//create content revisions
	for _, cr := range c.ContentRevisions {
		cr.ContentId = c.Id
		err = cr.Create()
		if err != nil {
			return err
		}
	}
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
		c.WebsiteId,
		c.Id,
	)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	//create/update content revisions
	for _, cr := range c.ContentRevisions {
		cr.ContentId = c.Id
		if cr.Id > 0 {
			err = cr.Update()
		} else {
			err = cr.Create()
		}
		if err != nil {
			return err
		}
	}
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
