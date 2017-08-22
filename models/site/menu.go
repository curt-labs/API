package site

import (
	"database/sql"

	"github.com/curt-labs/API/helpers/apicontext"
	"github.com/curt-labs/API/helpers/database"

	// "github.com/curt-labs/API/helpers/redis"
	_ "github.com/go-sql-driver/mysql"
)

type Menu struct {
	Id                    int      `json:"id,omitempty" xml:"id,omitempty"`
	Name                  string   `json:"name,omitempty" xml:"name,omitempty"`
	IsPrimary             bool     `json:"isPrimary,omitempty" xml:"isPrimary,omitempty"`
	Active                bool     `json:"active,omitempty" xml:"active,omitempty"`
	DisplayName           string   `json:"displayName,omitempty" xml:"displayName,omitempty"`
	RequireAuthentication bool     `json:"requireAuthentication,omitempty" xml:"requireAuthentication,omitempty"`
	ShowOnSitemap         bool     `json:"showOnSitemap,omitempty" xml:showOnSitemap,omitempty"`
	Sort                  int      `json:"sort,omitempty" xml:"sort,omitempty"`
	WebsiteId             int      `json:"websiteId,omitempty" xml:"websiteId,omitempty"`
	Contents              Contents `json:"contents,omitempty" xml:"contents,omitempty"`
}
type Menus []Menu

const (
	menuFields            = "m.menuID, m.menu_name, m.isPrimary, m.active, m.display_name, m.requireAuthentication, m.showOnSiteMap, m.sort, m.websiteID"                                                                                           //menu AS m
	menuSiteContentFields = "msc.menuSort, msc.menuTitle, msc.menuLink, msc.parentID, msc.linkTarget"                                                                                                                                               //omits join ids  as msc
	siteContentFields     = "s.contentID, s.content_type, s.page_title, s.createdDate, s.lastModified, s.meta_title, s.meta_description, s.keywords, s.isPrimary, s.published, s.active, s.slug, s.requireAuthentication, s.canonical, s.contentID" //as s

)

var (
	getMenu = ` SELECT ` + menuFields + ` FROM Menu AS m
								Join WebsiteToBrand as wub on wub.WebsiteID = m.websiteID
								Join ApiKeyToBrand as akb on akb.brandID = wub.brandID
								Join ApiKey as ak on akb.keyID = ak.id
	 							WHERE menuID = ? && (ak.api_key = ? && (wub.brandID = ? OR 0=?))`
	getAllMenus = ` SELECT ` + menuFields + ` FROM Menu AS m
								Join WebsiteToBrand as wub on wub.WebsiteID = m.websiteID
								Join ApiKeyToBrand as akb on akb.brandID = wub.brandID
								Join ApiKey as ak on akb.keyID = ak.id
	 							WHERE (ak.api_key = ? && (wub.brandID = ? OR 0=?))`
	getMenuContents = `SELECT ` + siteContentFields + `, ` + menuSiteContentFields + `  from Menu_SiteContent as msc JOIN SiteContent AS s ON s.contentID = msc.ContentID  WHERE msc.menuID = ?`
	getMenuByName   = ` SELECT ` + menuFields + ` FROM Menu AS m
		JOIN WebsiteToBrand AS wub ON wub.WebsiteID = m.websiteID
		WHERE menu_name = ? && (wub.brandID = ? OR 0 = ?)`
	//operations
	createMenu                    = `INSERT INTO Menu (menu_name, isPrimary, active, display_name, requireAuthentication, showOnSiteMap, sort, websiteID) VALUES(?,?,?,?,?,?,?,?)`
	updateMenu                    = `UPDATE Menu SET menu_name = ?, isPrimary = ?, active = ?, display_name = ?, requireAuthentication = ?, showOnSiteMap = ?, sort = ?, websiteID = ? WHERE menuID = ?`
	deleteMenu                    = `DELETE FROM Menu WHERE menuID = ?`
	deleteMenuSiteContentByMenuId = `DELETE FROM Menu_SiteContent WHERE menuID = ?` //used when deleting menu
	createMenuContentJoin         = `INSERT INTO Menu_SiteContent (menuID, contentID, menuSort, menuTitle, menuLink, parentID, linkTarget) VALUES(?,?,?,?,?,?,?)`
	deleteMenuSiteContentJoin     = `DELETE FROM Menu_SiteContent WHERE menuID = ? AND contentID = ?`
)

//Fetch menu by Id
func (m *Menu) Get(dtx *apicontext.DataContext) (err error) {
	err = database.Init()
	if err != nil {
		return err
	}

	stmt, err := database.DB.Prepare(getMenu)
	if err != nil {
		return err
	}
	defer stmt.Close()

	var display *string

	err = stmt.QueryRow(m.Id, dtx.APIKey, dtx.BrandID, dtx.BrandID).Scan(
		&m.Id,
		&m.Name,
		&m.IsPrimary,
		&m.Active,
		&display,
		&m.RequireAuthentication,
		&m.ShowOnSitemap,
		&m.Sort,
		&m.WebsiteId,
	)
	if err != nil {
		return err
	}
	if display != nil {
		m.DisplayName = *display
	}
	return err
}

//Fetch up a menu by name
func (m *Menu) GetByName(dtx *apicontext.DataContext) (err error) {
	err = database.Init()
	if err != nil {
		return err
	}

	stmt, err := database.DB.Prepare(getMenuByName)
	if err != nil {
		return err
	}
	defer stmt.Close()

	var display *string

	err = stmt.QueryRow(m.Name, dtx.BrandID, dtx.BrandID).Scan(
		&m.Id,
		&m.Name,
		&m.IsPrimary,
		&m.Active,
		&display,
		&m.RequireAuthentication,
		&m.ShowOnSitemap,
		&m.Sort,
		&m.WebsiteId,
	)
	if err != nil {
		return err
	}
	if display != nil {
		m.DisplayName = *display
	}
	return err
}

//Fetch all menus
func GetAllMenus(dtx *apicontext.DataContext) (ms Menus, err error) {
	err = database.Init()
	if err != nil {
		return ms, err
	}

	stmt, err := database.DB.Prepare(getAllMenus)
	if err != nil {
		return ms, err
	}
	defer stmt.Close()

	var display *string
	var m Menu

	res, err := stmt.Query(dtx.APIKey, dtx.BrandID, dtx.BrandID)
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
			&m.WebsiteId,
		)
		if err != nil {
			return ms, err
		}
		if display != nil {
			m.DisplayName = *display
		}
		ms = append(ms, m)
	}
	defer res.Close()
	return ms, err
}

//Fetch a menu's contents, including latest revision
func (m *Menu) GetContents() (err error) {
	err = database.Init()
	if err != nil {
		return err
	}

	stmt, err := database.DB.Prepare(getMenuContents)
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
			&m.WebsiteId,
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
			if err == sql.ErrNoRows {
				err = nil
			}
		}
		m.Contents = append(m.Contents, c)
	}
	defer res.Close()
	return err
}

//creating a menu
func (m *Menu) Create() (err error) {
	err = database.Init()
	if err != nil {
		return err
	}

	tx, err := database.DB.Begin()
	stmt, err := tx.Prepare(createMenu)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(
		m.Name,
		m.IsPrimary,
		m.Active,
		m.DisplayName,
		m.RequireAuthentication,
		m.ShowOnSitemap,
		m.Sort,
		m.WebsiteId,
	)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	m.Id = int(id)
	return err
}

//updating a menu
func (m *Menu) Update() (err error) {
	err = database.Init()
	if err != nil {
		return err
	}

	tx, err := database.DB.Begin()
	stmt, err := tx.Prepare(updateMenu)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(m.Name,
		m.IsPrimary,
		m.Active,
		m.DisplayName,
		m.RequireAuthentication,
		m.ShowOnSitemap,
		m.Sort,
		m.WebsiteId,
		m.Id,
	)
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

//deleting a menu, takes a content_sitecontent join with
func (m *Menu) Delete() (err error) {
	err = database.Init()
	if err != nil {
		return err
	}

	tx, err := database.DB.Begin()

	//delete menu content join
	stmt, err := tx.Prepare(deleteMenuSiteContentByMenuId)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(m.Id)
	if err != nil {
		tx.Rollback()
		return err
	}

	//delete menu
	stmt, err = tx.Prepare(deleteMenu)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(m.Id)
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

//there needs to exist a menu object with id > 0 for there be a FK relation
func (m *Menu) JoinToContent(c Content) (err error) {
	err = database.Init()
	if err != nil {
		return err
	}

	tx, err := database.DB.Begin()

	stmt, err := tx.Prepare(createMenuContentJoin)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(m.Id, c.Id, c.MenuSort, c.MenuTitle, c.MenuLink, c.ParentId, c.LinkTarget)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	return err
}

//For deleting a join
func (m *Menu) DeleteMenuContentJoin(c Content) (err error) {
	err = database.Init()
	if err != nil {
		return err
	}

	tx, err := database.DB.Begin()

	stmt, err := tx.Prepare(deleteMenuSiteContentJoin)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(m.Id, c.Id)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	return err
}
