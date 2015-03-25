package products

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/curt-labs/GoAPI/helpers/apicontext"
	"github.com/curt-labs/GoAPI/helpers/database"
	"github.com/curt-labs/GoAPI/helpers/redis"
	"github.com/curt-labs/GoAPI/helpers/sortutil"
	"github.com/curt-labs/GoAPI/models/customer/content"
	_ "github.com/go-sql-driver/mysql"

	"net/url"
	"strings"
	"time"
)

var (
	CategoriesByBrandStmt = `
		select c.catID, c.ParentID, c.sort, c.dateAdded,
		c.catTitle, c.shortDesc, c.longDesc,
		c.image, c.icon, c.isLifestyle, c.vehicleSpecific,
		c.vehicleRequired, cc.code, cc.font 
		from Categories as c
		left join ColorCode as cc on c.codeID = cc.codeID
		where c.BrandID = ?`
	PartCategoryStmt = `
		select c.catID, c.ParentID, c.sort, c.dateAdded,
		c.catTitle, c.shortDesc, c.longDesc,
		c.image, c.icon, c.isLifestyle, c.vehicleSpecific,
		c.vehicleRequired,
		cc.code, cc.font from Categories as c
		join CatPart as cp on c.catID = cp.catID
		left join ColorCode as cc on c.codeID = cc.codeID
		where cp.partID = ?
		order by c.sort
		limit 1`
	PartAllCategoryStmt = `
		select c.catID, c.ParentID, c.sort, c.dateAdded,
		c.catTitle, c.shortDesc, c.longDesc,
		c.image, c.icon, c.isLifestyle, c.vehicleSpecific,
		c.vehicleRequired,
		cc.code, cc.font
		from Categories as c
		join CatPart as cp on c.catID = cp.catID
		join ColorCode as cc on c.codeID = cc.codeID
		where cp.partID = ?
		order by c.catID`
	ParentCategoryStmt = `
		select c.catID, c.ParentID, c.sort, c.dateAdded,
		c.catTitle, c.shortDesc, c.longDesc,
		c.image, c.icon, c.isLifestyle, c.vehicleSpecific,
		c.vehicleRequired,
		cc.code, cc.font from Categories as c
		left join ColorCode as cc on c.codeID = cc.codeID
		where c.catID = ?
		order by c.sort
		limit 1`
	TopCategoriesStmt = `
		select c.catID from Categories as c
		join ApiKeyToBrand as akb on akb.brandID = c.brandID
		join ApiKey as ak on ak.id = akb.keyID
		where c.ParentID IS NULL or c.ParentID = 0
		and isLifestyle = 0
		and (ak.api_key = ? && (c.BrandID = ? or 0 = ?))
		order by c.sort`
	SubCategoriesStmt = `
		select c.catID, c.ParentID, c.sort, c.dateAdded,
		c.catTitle, c.shortDesc, c.longDesc,
		c.image, c.icon, c.isLifestyle, c.vehicleSpecific,
		c.vehicleRequired,
		cc.code, cc.font from Categories as c
		left join ColorCode as cc on c.codeID = cc.codeID
		where c.ParentID = ?
		and isLifestyle = 0
		order by c.sort`
	CategoryByTitleStmt = `
		select c.catID, c.ParentID, c.sort, c.dateAdded,
		c.catTitle, c.shortDesc, c.longDesc,
		c.image, c.icon, c.isLifestyle, c.vehicleSpecific,
		c.vehicleRequired,
		cc.code, cc.font from Categories as c
		left join ColorCode as cc on c.codeID = cc.codeID
		join ApiKeyToBrand as akb on akb.brandID = c.brandID
		join ApiKey as ak on ak.id = akb.keyID
		where c.catTitle = ?
		and (ak.api_key = ? && (c.BrandID = ? or 0=?))
		order by c.sort`
	CategoryByIdStmt = `
		select c.catID, c.ParentID, c.sort, c.dateAdded,
		c.catTitle, c.shortDesc, c.longDesc,
		c.image, c.icon, c.isLifestyle, c.vehicleSpecific,
		c.vehicleRequired,
		cc.code, cc.font from Categories as c
		left join ColorCode as cc on c.codeID = cc.codeID
		where c.catID = ?
		order by c.sort`
	CategoryPartsStmt = `
		select p.partID, 1 from Part as p
		join CatPart as cp on p.partID = cp.partID
		where (p.status = 800 || p.status = 900) && FIND_IN_SET(cp.catID,bottom_category_ids(?))
		order by p.partID
		limit ?,?`
	CategoryPartsFilteredStmt = `select p.partID, (
			select count(pa.pAttrID) from PartAttribute as pa
			where pa.partID = cp.partID && FIND_IN_SET(REPLACE(pa.field, ',','|'),?) &&
			FIND_IN_SET(REPLACE(pa.value, ',','|'),?)
		) as cnt from Part as p
		join CatPart as cp on p.partID = cp.partID
		where (p.status = 800 || p.status = 900) && FIND_IN_SET(cp.catID,bottom_category_ids(?))
		having cnt >= ?
		order by p.partID
		limit ?,?`
	CategoryPartCountStmt = `
		select count(p.partID) as count from Part as p
		join CatPart as cp on p.partID = cp.partID
		where (p.status = 800 || p.status = 900) && FIND_IN_SET(cp.catID,bottom_category_ids(?))
		order by p.partID`
	CategoryFilteredPartCountStmt = `select p.partID, (
			select count(pa.pAttrID) from PartAttribute as pa
			where pa.partID = cp.partID && FIND_IN_SET(pa.field,?) &&
			FIND_IN_SET(pa.value,?)
		) as cnt from Part as p
		join CatPart as cp on p.partID = cp.partID
		where (p.status = 800 || p.status = 900) && FIND_IN_SET(cp.catID,bottom_category_ids(?))
		having cnt >= ?
		order by p.partID`
	SubIDStmt = `
		select c.catID, group_concat(p.partID) as parts from Categories as c
		left join CatPart as cp on c.catID = cp.catID
		left join Part as p on cp.partID = p.partID
		where c.ParentID = ? && (p.status = null || (p.status = 800 || p.status = 900))`
	CategoryContentStmt = `
		select ct.cTypeID, ct.type, c.text from ContentBridge cb
		join Content as c on cb.contentID = c.contentID
		left join ContentType as ct on c.cTypeID = ct.cTypeID
		where cb.catID = ?`
	createCategory = `insert into Categories (dateAdded, parentID, catTitle, shortDesc, longDesc, image, isLifestyle, codeId, sort, vehicleSpecific, vehicleRequired, metaTitle, metaDesc, metaKeywords, icon, path, brandID)
						values(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`
	deleteCategory = `delete from Categories where catID = ? `
)

const (
	DefaultPageCount = 10
)

type Category struct {
	ID              int                      `json:"id" xml:"id,attr"`
	ParentID        int                      `json:"parent_id" xml:"parent_id,attr"`
	Sort            int                      `json:"sort" xml:"sort,attr"`
	DateAdded       time.Time                `json:"date_added" xml:"date_added,attr"`
	Title           string                   `json:"title" xml:"title,attr"`
	ShortDesc       string                   `json:"short_description" xml:"short_description"`
	LongDesc        string                   `json:"long_description" xml:"long_description"`
	ColorCode       string                   `json:"color_code" xml:"color_code,attr"`
	FontCode        string                   `json:"font_code" xml:"font_code,attr"`
	Image           *url.URL                 `json:"image" xml:"image"`
	Icon            *url.URL                 `json:"icon" xml:"icon"`
	IsLifestyle     bool                     `json:"lifestyle" xml:"lifestyle,attr"`
	VehicleSpecific bool                     `json:"vehicle_specific" xml:"vehicle_specific,attr"`
	VehicleRequired bool                     `json:"vehicle_required" xml:"vehicle_required,attr"`
	Content         []Content                `json:"content,omitempty" xml:"content,omitempty"`
	SubCategories   []Category               `json:"sub_categories,omitempty" xml:"sub_categories,omitempty"`
	ProductListing  *PaginatedProductListing `json:"product_listing,omitempty" xml:"product_listing,omitempty"`
	Filter          interface{}              `json:"filter,omitempty" xml:"filter,omitempty"`
	MetaTitle       string                   `json:"metaTitle,omitempty" xml:"v,omitempty"`
	MetaDescription string                   `json:"metaDescription,omitempty" xml:"metaDescription,omitempty"`
	MetaKeywords    string                   `json:"metaKeywords,omitempty" xml:"metaKeywords,omitempty"`
	BrandID         int                      `json:"categoryId,omitempty" xml:"categoryId,omitempty"`
}

func PopulateCategoryMulti(rows *sql.Rows, ch chan []Category) {
	cats := make([]Category, 0)
	if rows == nil {
		ch <- cats
		return
	}

	for rows.Next() {
		var initCat Category
		var catImg *string
		var catIcon *string
		var colorCode *string
		var fontCode *string
		err := rows.Scan(
			&initCat.ID,
			&initCat.ParentID,
			&initCat.Sort,
			&initCat.DateAdded,
			&initCat.Title,
			&initCat.ShortDesc,
			&initCat.LongDesc,
			&catImg,
			&catIcon,
			&initCat.IsLifestyle,
			&initCat.VehicleSpecific,
			&initCat.VehicleRequired,
			&colorCode,
			&fontCode)
		if err != nil {
			ch <- cats
			return
		}

		// Attempt to parse out the image Url
		if catImg != nil {
			initCat.Image, _ = url.Parse(*catImg)
		}
		if catIcon != nil {
			initCat.Icon, _ = url.Parse(*catIcon)
		}

		// Build out RGB value for color coding
		if colorCode != nil && fontCode != nil && len(*colorCode) == 9 {
			cc := fmt.Sprintf("%s", *colorCode)
			initCat.ColorCode = fmt.Sprintf("rgb(%s,%s,%s)", cc[0:3], cc[3:6], cc[6:9])
			initCat.FontCode = fmt.Sprintf("#%s", *fontCode)
		}

		initCat.ProductListing = nil

		cats = append(cats, initCat)
	}
	defer rows.Close()

	ch <- cats
}

func PopulateCategory(row *sql.Row, ch chan Category, dtx *apicontext.DataContext) {
	if row == nil {
		ch <- Category{}
		return
	}

	var initCat Category
	var catImg *string
	var catIcon *string
	var colorCode *string
	var fontCode *string
	err := row.Scan(
		&initCat.ID,
		&initCat.ParentID,
		&initCat.Sort,
		&initCat.DateAdded,
		&initCat.Title,
		&initCat.ShortDesc,
		&initCat.LongDesc,
		&catImg,
		&catIcon,
		&initCat.IsLifestyle,
		&initCat.VehicleSpecific,
		&initCat.VehicleRequired,
		&colorCode,
		&fontCode)
	if err != nil {
		ch <- Category{}
		return
	}

	// Attempt to parse out the image Url
	if catImg != nil {
		initCat.Image, _ = url.Parse(*catImg)
	}
	if catIcon != nil {
		initCat.Icon, _ = url.Parse(*catIcon)
	}

	// Build out RGB value for color coding
	if colorCode != nil && fontCode != nil && len(*colorCode) == 9 {
		cc := fmt.Sprintf("%s", *colorCode)
		initCat.ColorCode = fmt.Sprintf("rgb(%s,%s,%s)", cc[0:3], cc[3:6], cc[6:9])
		initCat.FontCode = fmt.Sprintf("#%s", *fontCode)
	}

	conChan := make(chan []Content)
	go func() {
		con, _ := initCat.GetContent()
		conChan <- con
	}()

	if subCats, err := initCat.GetSubCategories(dtx); err == nil {
		initCat.SubCategories = subCats
	}

	select {
	case con := <-conChan:
		initCat.Content = con
	}

	ch <- initCat
}

// TopTierCategories
// Description: Returns the top tier categories
// Returns: []Category, error
func TopTierCategories(dtx *apicontext.DataContext) (cats []Category, err error) {
	cats = make([]Category, 0)

	redis_key := fmt.Sprintf("category:top:%s", dtx.BrandString)
	// First lets try to access the category:top endpoint in Redis
	data, err := redis.Get(redis_key)
	if len(data) > 0 && err == nil {
		err = json.Unmarshal(data, &cats)
		if err == nil {
			return
		}
	}

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return
	}
	defer db.Close()

	qry, err := db.Prepare(TopCategoriesStmt)
	if err != nil {
		return
	}
	defer qry.Close()

	// Execute SQL Query against current PartId
	catRows, err := qry.Query(dtx.APIKey, dtx.BrandID, dtx.BrandID)
	if err != nil || catRows == nil { // Error occurred while executing query
		return
	}

	var iter int
	ch := make(chan error)
	for catRows.Next() {
		var cat Category
		err := catRows.Scan(&cat.ID)
		if err == nil {
			go func(c Category) {
				err := c.GetCategory(dtx.APIKey, 0, 0, true, nil, nil, dtx)
				if err == nil {
					cats = append(cats, c)
				}
				ch <- err
			}(cat)
			iter++
		}
	}
	defer catRows.Close()

	for i := 0; i < iter; i++ {
		<-ch
	}

	sortutil.AscByField(cats, "Sort")
	if dtx.BrandString != "" {
		go redis.Setex(redis_key, cats, 86400)
	}
	return
}

func GetCategoryByTitle(cat_title string, dtx *apicontext.DataContext) (cat Category, err error) {
	redis_key := fmt.Sprintf("category:title:%s:%s", cat_title, dtx.BrandString)
	// Attempt to get the category from Redis
	data, err := redis.Get(redis_key)
	if len(data) > 0 && err == nil {
		err = json.Unmarshal(data, &cat)
		if err == nil {
			return cat, err
		}
	}

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return cat, err
	}
	defer db.Close()

	qry, err := db.Prepare(CategoryByTitleStmt)
	if err != nil {
		return cat, err
	}
	defer qry.Close()

	// Execute SQL Query against title
	catRow := qry.QueryRow(cat_title, dtx.APIKey, dtx.BrandID, dtx.BrandID)
	if catRow == nil { // Error occurred while executing query
		return cat, err
	}

	ch := make(chan Category)
	go PopulateCategory(catRow, ch, dtx)
	cat = <-ch
	if dtx.BrandString != "" {
		go redis.Setex(redis_key, cat, 86400)
	}
	return cat, err
}

func GetCategoryById(cat_id int, dtx *apicontext.DataContext) (cat Category, err error) {

	redis_key := fmt.Sprintf("category:id:%d", cat_id)

	// Attempt to get the category from Redis
	data, err := redis.Get(redis_key)
	if len(data) > 0 && err == nil {
		err = json.Unmarshal(data, &cat)
		if err == nil {
			return
		}
	}

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return
	}
	defer db.Close()

	qry, err := db.Prepare(CategoryByIdStmt)
	if err != nil {
		return
	}
	defer qry.Close()

	// Execute SQL Query against title
	catRow := qry.QueryRow(cat_id)
	if catRow == nil { // Error occurred while executing query
		return
	}

	ch := make(chan Category)
	go PopulateCategory(catRow, ch, dtx)
	cat = <-ch

	go redis.Setex(redis_key, cat, 86400)

	return
}

func (c *Category) GetSubCategories(dtx *apicontext.DataContext) (cats []Category, err error) {
	cats = make([]Category, 0)

	if c.ID == 0 {
		return
	}

	redis_key := fmt.Sprintf("category:%d:subs:%s", c.ID, dtx.BrandString)
	// First lets try to access the category:top endpoint in Redis
	data, err := redis.Get(redis_key)
	if len(data) > 0 && err == nil {
		err = json.Unmarshal(data, &cats)
		if err == nil {
			return cats, err
		}
	}

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return
	}
	defer db.Close()

	qry, err := db.Prepare(SubCategoriesStmt)
	if err != nil {
		return
	}
	defer qry.Close()

	// Execute SQL Query against current PartId
	catRows, err := qry.Query(c.ID)
	if err != nil || catRows == nil { // Error occurred while executing query
		return
	}

	ch := make(chan []Category, 0)
	go PopulateCategoryMulti(catRows, ch)
	cats = <-ch
	if dtx.BrandString != "" {
		go redis.Setex(redis_key, cats, 86400)
	}
	return
}

func (c *Category) GetCategory(key string, page int, count int, ignoreParts bool, v *Vehicle, specs *map[string][]string, dtx *apicontext.DataContext) error {
	var err error
	if c.ID == 0 {
		return fmt.Errorf("error: %s", "invalid category reference")
	}

	if v != nil && v.Base.Year == 0 {
		v = nil
	}

	redis_key := fmt.Sprintf("category:%d:%d", c.BrandID, c.ID)

	// First lets try to access the category:top endpoint in Redis
	cat_bytes, err := redis.Get(redis_key)
	if len(cat_bytes) > 0 {
		err = json.Unmarshal(cat_bytes, &c)
	}

	if err != nil || c.ShortDesc == "" {
		cat, catErr := GetCategoryById(c.ID, dtx)
		if catErr != nil {
			return catErr
		}

		c.ID = cat.ID
		c.BrandID = cat.BrandID
		c.ColorCode = cat.ColorCode
		c.DateAdded = cat.DateAdded
		c.FontCode = cat.FontCode
		c.Image = cat.Image
		c.Icon = cat.Icon
		c.IsLifestyle = cat.IsLifestyle
		c.LongDesc = cat.LongDesc
		c.ParentID = cat.ParentID
		c.ShortDesc = cat.ShortDesc
		c.Sort = cat.Sort
		c.Title = cat.Title
		c.VehicleSpecific = cat.VehicleSpecific
		c.VehicleRequired = cat.VehicleRequired
		c.Content = cat.Content
		c.SubCategories = cat.SubCategories
		c.ProductListing = cat.ProductListing
		c.Filter = cat.Filter

		go redis.Setex(redis_key, c, 86400)
	}

	partChan := make(chan *PaginatedProductListing)
	if !ignoreParts {
		go func() {
			c.GetParts(key, page, count, v, specs, dtx)
			partChan <- c.ProductListing
		}()
	} else {
		close(partChan)
	}

	content, err := custcontent.GetCategoryContent(c.ID, key)
	for _, con := range content {
		strArr := strings.Split(con.ContentType.Type, ":")
		cType := con.ContentType.Type
		if len(strArr) > 1 {
			cType = strArr[1]
		}
		var co Content
		co.ContentType.Type = cType
		co.Text = con.Text
		c.Content = append(c.Content, co)
	}

	if !ignoreParts {
		c.ProductListing = <-partChan
		close(partChan)
	}
	return err
}

func (c *Category) GetContent() (content []Content, err error) {
	if c.ID == 0 {
		return
	}

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return
	}
	defer db.Close()

	qry, err := db.Prepare(CategoryContentStmt)
	if err != nil {
		return
	}
	defer qry.Close()

	// Execute SQL Query against current ID
	conRows, err := qry.Query(c.ID)
	if err != nil || conRows == nil {
		return
	}

	for conRows.Next() {
		var con Content
		if err := conRows.Scan(&con.ContentType.Id, &con.ContentType.Type, &con.Text); err == nil {
			content = append(content, con)
		}
	}
	defer conRows.Close()

	return
}

func (c *Category) GetParts(key string, page int, count int, v *Vehicle, specs *map[string][]string, dtx *apicontext.DataContext) error {
	var err error
	c.ProductListing = &PaginatedProductListing{}
	if c.ID == 0 {
		return fmt.Errorf("error: %s %d", "invalid category reference", c.ID)
	}

	if count == 0 {
		count = DefaultPageCount
	}
	queryPage := page
	if page == 1 {
		queryPage = count
	} else if page > 1 {
		queryPage = count * (page - 1)
	}

	if v != nil {
		vehicleChan := make(chan []Part)
		l := Lookup{
			Vehicle:     *v,
			CustomerKey: key,
		}
		go l.LoadParts(vehicleChan, page, count, dtx)

		parts := <-vehicleChan

		for _, p := range parts {
			for _, partCat := range p.Categories {
				if partCat.ID == c.ID {
					c.ProductListing.Parts = append(c.ProductListing.Parts, p)
					break
				}
			}
		}

		c.ProductListing.ReturnedCount = len(c.ProductListing.Parts)
		c.ProductListing.PerPage = c.ProductListing.ReturnedCount
		c.ProductListing.Page = 1
		c.ProductListing.TotalItems = c.ProductListing.ReturnedCount
		c.ProductListing.TotalPages = 1

		return nil
	}

	parts := make([]Part, 0)

	redis_key := fmt.Sprintf("category:%d:%d:parts:%d:%d", c.BrandID, c.ID, queryPage, count)

	// First lets try to access the category:top endpoint in Redis
	part_bytes, err := redis.Get(redis_key)
	if len(part_bytes) > 0 {
		err = json.Unmarshal(part_bytes, &parts)
	}

	if err != nil || len(parts) == 0 || specs != nil {
		db, err := sql.Open("mysql", database.ConnectionString())
		if err != nil {
			return err
		}
		defer db.Close()

		var rows *sql.Rows

		if specs != nil && len(*specs) > 0 {
			keys := make([]string, 0)
			values := make([]string, 0)
			for k, vals := range *specs {
				keys = append(keys, strings.Replace(k, ",", "|", -1))
				for _, val := range vals {
					values = append(values, strings.Replace(val, ",", "|", -1))
				}
			}

			stmt, err := db.Prepare(CategoryPartsFilteredStmt)
			if err != nil {
				return err
			}
			defer stmt.Close()

			rows, err = stmt.Query(strings.Join(keys, ","), strings.Join(values, ","), c.ID, len(keys), queryPage, count)
		} else {
			stmt, err := db.Prepare(CategoryPartsStmt)
			if err != nil {
				return err
			}
			defer stmt.Close()

			rows, err = stmt.Query(c.ID, queryPage, count)
		}

		if err != nil || rows == nil {
			return err
		}

		ch := make(chan error)
		count := 0

		for rows.Next() {
			var ct *int
			var id *int
			if err := rows.Scan(&id, &ct); err == nil && id != nil {
				go func(i int) {
					p := Part{
						ID: i,
					}
					err := p.Get(dtx)
					if err == nil {
						c.ProductListing.Parts = append(c.ProductListing.Parts, p)
					}
					ch <- err
				}(*id)
				count++

			}
		}
		defer rows.Close()

		for i := 0; i < count; i++ {
			<-ch
		}
	}

	sortutil.AscByField(c.ProductListing.Parts, "ID")

	c.ProductListing.ReturnedCount = len(c.ProductListing.Parts)
	c.ProductListing.PerPage = count
	if page == 0 {
		page = 1
	}
	c.ProductListing.Page = page
	c.GetPartCount(key, v, specs)
	c.ProductListing.TotalPages = c.ProductListing.TotalItems / count

	go redis.Setex(redis_key, parts, redis.CacheTimeout)

	return nil
}

func (c *Category) GetPartCount(key string, v *Vehicle, specs *map[string][]string) error {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	var total *int
	if specs != nil && len(*specs) > 0 {
		keys := make([]string, 0)
		values := make([]string, 0)
		for k, vals := range *specs {
			keys = append(keys, k)
			values = append(values, vals...)
		}

		stmt, err := db.Prepare(CategoryFilteredPartCountStmt)
		if err != nil {
			return err
		}
		defer stmt.Close()

		rows, err := stmt.Query(strings.Join(keys, ","), strings.Join(values, ","), c.ID, len(keys))
		if err != nil {
			return err
		}

		counter := 0
		for rows.Next() {
			counter++
		}
		defer rows.Close()
		total = &counter
	} else {
		stmt, err := db.Prepare(CategoryPartCountStmt)
		if err != nil {
			return err
		}
		defer stmt.Close()

		row := stmt.QueryRow(c.ID)
		if row == nil {
			return fmt.Errorf("error: %s", "failed to retrieve part count")
		}

		if err := row.Scan(&total); err != nil || total == nil {
			return err
		}
	}

	c.ProductListing.TotalItems = *total

	return nil
}

func (c *Category) Create() (err error) {
	go redis.Delete("category:top*")
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	stmt, err := db.Prepare(createCategory)
	if err != nil {
		return err
	}
	defer stmt.Close()
	res, err := stmt.Exec(
		c.DateAdded,
		c.ParentID,
		c.Title,
		c.ShortDesc,
		c.LongDesc,
		c.Image,
		c.IsLifestyle,
		0, // this should be selected from the ColorCode table
		c.Sort,
		c.VehicleSpecific,
		c.VehicleRequired,
		c.MetaTitle,
		c.MetaDescription,
		c.MetaKeywords,
		c.Icon,
		c.Image,
		c.BrandID,
	)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	c.ID = int(id)

	redis.Setex(fmt.Sprintf("category:%d:%d", c.BrandID, c.ID), c, redis.CacheTimeout)

	return err
}

func (c *Category) Delete(dtx *apicontext.DataContext) (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(deleteCategory)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(c.ID)
	if err == nil {
		go redis.Delete(fmt.Sprintf("category:%d:%d", c.BrandID, c.ID))
		go redis.Delete("category:title:" + dtx.BrandString)
		go redis.Delete("category:top:" + dtx.BrandString)
	}

	return err
}
