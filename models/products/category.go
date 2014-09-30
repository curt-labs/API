package products

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/curt-labs/GoAPI/helpers/database"
	"github.com/curt-labs/GoAPI/helpers/redis"
	"github.com/curt-labs/GoAPI/helpers/sortutil"
	"github.com/curt-labs/GoAPI/models/customer/content"
	"github.com/curt-labs/GoAPI/models/vehicle"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var (
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
		where c.ParentID IS NULL or c.ParentID = 0
		and isLifestyle = 0
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
	CategoryByNameStmt = `
		select c.catID, c.ParentID, c.sort, c.dateAdded,
		c.catTitle, c.shortDesc, c.longDesc,
		c.image, c.icon, c.isLifestyle, c.vehicleSpecific,
		c.vehicleRequired,
		cc.code, cc.font from Categories as c
		left join ColorCode as cc on c.codeID = cc.codeID
		where c.catTitle = ?
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
	CategoryPartBasicStmt = `
		select cp.partID
		from CatPart as cp
		where cp.catID = ?
		order by cp.partID
		limit ?,?`
	SubIDStmt = `
		select c.catID, group_concat(p.partID) as parts from Categories as c
		left join CatPart as cp on c.catID = cp.catID
		left join Part as p on cp.partID = p.partID
		where c.ParentID = ? && (p.status = null || (p.status = 800 || p.status = 900))`
	CategoryContentStmt = `
		select ct.type, c.text from ContentBridge cb
		join Content as c on cb.contentID = c.contentID
		left join ContentType as ct on c.cTypeID = ct.cTypeID
		where cb.catID = ?`
)

type Category struct {
	ID              int                     `json:"id" xml:"id"`
	ParentID        int                     `json:"parent_id" xml:"parent_id"`
	Sort            int                     `json:"sort" xml:"sort"`
	DateAdded       time.Time               `json:"date_added" xml:"date_added"`
	Title           string                  `json:"title" xml:"title"`
	ShortDesc       string                  `json:"short_description" xml:"short_description"`
	LongDesc        string                  `json:"long_description" xml:"long_description"`
	ColorCode       string                  `json:"color_code" xml:"color_code"`
	FontCode        string                  `json:"font_code" xml:"font_code"`
	Image           *url.URL                `json:"image" xml:"image"`
	Icon            *url.URL                `json:"icon" xml:"icon"`
	IsLifestyle     bool                    `json:"lifestyle" xml:"lifestyle"`
	VehicleSpecific bool                    `json:"vehicle_specific" xml:"vehicle_specific"`
	VehicleRequired bool                    `json:"vehicle_required" xml:"vehicle_required"`
	Content         []Content               `json:"content,omitempty" xml:"content,omitempty"`
	SubCategories   []Category              `json:"sub_categories,omitempty" xml:"sub_categories,omitempty"`
	ProductListing  PaginatedProductListing `json:"product_listing,omitempty" xml:"product_listing,omitempty"`
	Filter          interface{}             `json:"filter,omitempty" xml:"filter,omitempty"`
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

		// con, err := initCat.GetContent()
		// if err == nil {
		// 	initCat.Content = con
		// }

		// if subCats, err := initCat.GetSubCategories(); err == nil {
		// 	initCat.SubCategories = subCats
		// }

		cats = append(cats, initCat)
	}

	ch <- cats
}

func PopulateCategory(row *sql.Row, ch chan Category) {
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
		log.Println(err)
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

	con, err := initCat.GetContent()
	if err == nil {
		initCat.Content = con
	}

	if subCats, err := initCat.GetSubCategories(); err == nil {
		initCat.SubCategories = subCats
	}

	ch <- initCat
}

// TopTierCategories
// Description: Returns the top tier categories
// Returns: []Category, error
func TopTierCategories(key string) (cats []Category, err error) {

	redis_key := "category:top"

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
	catRows, err := qry.Query()
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
				err := c.GetCategory(key)
				if err == nil {
					cats = append(cats, c)
				}
				ch <- err
			}(cat)
			iter++
		}
	}

	for i := 0; i < iter; i++ {
		<-ch
	}

	sortutil.AscByField(cats, "Sort")

	go redis.Setex(redis_key, cats, 86400)

	return
}

func GetCategoryByTitle(cat_title string) (cat Category, err error) {

	redis_key := "category:title:" + cat_title

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

	qry, err := db.Prepare(CategoryByNameStmt)
	if err != nil {
		return
	}
	defer qry.Close()

	// Execute SQL Query against title
	catRow := qry.QueryRow(cat_title)
	if catRow == nil { // Error occurred while executing query
		return
	}

	ch := make(chan Category)
	go PopulateCategory(catRow, ch)
	cat = <-ch

	go redis.Setex(redis_key, cat, 86400)

	return
}

func GetCategoryById(cat_id int) (cat Category, err error) {

	redis_key := "category:id:" + strconv.Itoa(cat_id)

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
	go PopulateCategory(catRow, ch)
	cat = <-ch

	go redis.Setex(redis_key, cat, 86400)

	return
}

func (c *Category) GetSubCategories() (cats []Category, err error) {

	if c.ID == 0 {
		return
	}

	redis_key := "category:" + strconv.Itoa(c.ID) + ":subs"

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

	go redis.Setex(redis_key, cats, 86400)

	return
}

func (c *Category) GetCategory(key string) error {

	redis_key := "category:" + strconv.Itoa(c.ID)

	// First lets try to access the category:top endpoint in Redis
	cat_bytes, err := redis.Get(redis_key)
	if len(cat_bytes) > 0 {
		err = json.Unmarshal(cat_bytes, &c)
		if err == nil {
			content, err := custcontent.GetCategoryContent(c.ID, key)
			for _, con := range content {
				strArr := strings.Split(con.ContentType.Type, ":")
				cType := con.ContentType.Type
				if len(strArr) > 1 {
					cType = strArr[1]
				}
				c.Content = append(c.Content, Content{
					Key:   cType,
					Value: con.Text,
				})
			}
			return err
		}
	}

	cat, catErr := GetCategoryById(c.ID)
	if catErr != nil {
		return catErr
	}

	c.ID = cat.ID
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

	content, err := custcontent.GetCategoryContent(c.ID, key)
	for _, con := range content {
		strArr := strings.Split(con.ContentType.Type, ":")
		cType := con.ContentType.Type
		if len(strArr) > 1 {
			cType = strArr[1]
		}
		c.Content = append(c.Content, Content{
			Key:   cType,
			Value: con.Text,
		})
	}

	return nil
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
		if err := conRows.Scan(&con.Key, &con.Value); err == nil {
			content = append(content, con)
		}
	}

	return
}

func (c *Category) GetParts(page int, count int, v vehicle.Vehicle) error {
	redis_key := fmt.Sprintf("category:%d:parts:%d:%d", c.ID, page, count)

	parts := make([]Part, 0)
	// First lets try to access the category:top endpoint in Redis
	part_bytes, err := redis.Get(redis_key)
	if len(part_bytes) > 0 {
		err = json.Unmarshal(part_bytes, &parts)
		if err == nil {
			return nil
		}
	}

	return nil
}
