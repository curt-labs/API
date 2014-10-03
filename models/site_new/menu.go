package site_new

import (
	"database/sql"
	"github.com/curt-labs/GoAPI/helpers/database"
	// "github.com/curt-labs/GoAPI/helpers/redis"
	_ "github.com/go-sql-driver/mysql"
)

type Menu struct {
	Id                    int
	Name                  string
	IsPrimary             bool
	Active                bool
	DisplayName           string
	RequireAuthentication bool
	ShowOnSitemap         bool
	Sort                  int
	Contents              Contents
}
type Menus []Menu

const (
	menuFields            = "m.menuID, m.menu_name, m.isPrimary, m.active, m.display_name, m.requireAuthentication, m.showOnSiteMap, m.sort"                                                                                           //menu AS m
	menuSiteContentFields = "msc.menuSort, msc.menuTitle, msc.menuLink, msc.parentID, msc.linkTarget"                                                                                                                                  //omits join ids  as msc
	siteContentFields     = "s.contentID, s.content_type, s.page_title, s.createdDate, s.lastModified, s.meta_title, s.meta_description, s.keywords, s.isPrimary, s.published, s.active, s.slug, s.requireAuthentication, s.canonical" //as s

)

var (
	getMenu         = ` SELECT ` + menuFields + ` FROM Menu AS m WHERE menuID = ? `
	getAllMenus     = ` SELECT ` + menuFields + ` FROM Menu AS m`
	getMenuContents = `SELECT ` + siteContentFields + `, ` + menuSiteContentFields + `  from menu_sitecontent as msc JOIN SiteContent AS s ON s.contentID = msc.ContentID  WHERE msc.menuID = ?`
)

//Get menu by Id
func (m *Menu) Get() (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	stmt, err := db.Prepare(getMenu)
	if err != nil {
		return err
	}
	defer stmt.Close()

	var display *string

	err = stmt.QueryRow(m.Id).Scan(
		&m.Id,
		&m.Name,
		&m.IsPrimary,
		&m.Active,
		&display,
		&m.RequireAuthentication,
		&m.ShowOnSitemap,
		&m.Sort,
	)
	if err != nil {
		return err
	}
	if display != nil {
		m.DisplayName = *display
	}
	return err
}

//Get all menus
func GetAllMenus() (ms Menus, err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return ms, err
	}
	defer db.Close()
	stmt, err := db.Prepare(getAllMenus)
	if err != nil {
		return ms, err
	}
	defer stmt.Close()

	var display *string
	var m Menu

	res, err := stmt.Query()
	if err != nil {
		return ms, err
	}

	for res.Next() {
		res.Scan(
			&m.Id,
			&m.Name,
			&m.IsPrimary,
			&m.Active,
			&display,
			&m.RequireAuthentication,
			&m.ShowOnSitemap,
			&m.Sort,
		)
		if err != nil {
			return ms, err
		}
		if display != nil {
			m.DisplayName = *display
		}
		ms = append(ms, m)
	}
	return ms, err
}

//Get a menu's contents, including latest revision
func (m *Menu) GetContents() (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	stmt, err := db.Prepare(getMenuContents)
	if err != nil {
		return err
	}
	defer stmt.Close()
	res, err := stmt.Query(m.Id)
	var cType, title, mTitle, mDesc, slug, canon, menTitle, mLink *string

	var parent *int
	var c Content
	for res.Next() {
		err = res.Scan(
			&c.Id,
			&cType,
			&title,
			&c.CreatedDate,
			&c.LastModified,
			&mTitle,
			&c.MetaDescription,
			&c.Keywords,
			&c.IsPrimary,
			&c.Published,
			&c.Active,
			&slug,
			&c.RequireAuthentication,
			&canon,
			&c.MenuSort,
			&menTitle,
			&mLink,
			&parent,
			&c.LinkTarget,
		)
		if err != sql.ErrNoRows {
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
			if menTitle != nil {
				c.MenuTitle = *mTitle
			}
			if mLink != nil {
				c.MenuLink = *mLink
			}
			if parent != nil {
				c.ParentId = *parent
			}
			err = c.GetLatestRevision()
			if err != sql.ErrNoRows {
				if err != nil {
					return err
				}
			}
		}
		m.Contents = append(m.Contents, c)

	}
	return err
}
