package products

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/curt-labs/GoAPI/helpers/database"
	"github.com/curt-labs/GoAPI/helpers/redis"
	"github.com/curt-labs/GoAPI/models/customer/content"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var (
	PartCategoryStmt = `
		select c.catID, c.parentID, c.sort, c.dateAdded,
		c.catTitle, c.shortDesc, c.longDesc,
		c.image, c.icon, c.isLifestyle, c.vehicleSpecific,
		cc.code, cc.font from Categories as c
		join CatPart as cp on c.catID = cp.catID
		left join ColorCode as cc on c.codeID = cc.codeID
		where cp.partID = ?
		order by c.sort
		limit 1`
	PartAllCategoryStmt = `
		select c.catID, c.parentID, c.sort, c.dateAdded,
		c.catTitle, c.shortDesc, c.longDesc,
		c.image, c.icon, c.isLifestyle, c.vehicleSpecific,
		cc.code, cc.font
		from Categories as c
		join CatPart as cp on c.catID = cp.catID
		join ColorCode as cc on c.codeID = cc.codeID
		where cp.partID = ?
		order by c.catID`
	ParentCategoryStmt = `
		select c.catID, c.parentID, c.sort, c.dateAdded,
		c.catTitle, c.shortDesc, c.longDesc,
		c.image, c.icon, c.isLifestyle, c.vehicleSpecific,
		cc.code, cc.font from Categories as c
		left join ColorCode as cc on c.codeID = cc.codeID
		where c.catID = ?
		order by c.sort
		limit 1`
	TopCategoriesStmt = `
		select c.catID, c.parentID, c.sort, c.dateAdded,
		c.catTitle, c.shortDesc, c.longDesc,
		c.image, c.icon, c.isLifestyle, c.vehicleSpecific,
		cc.code, cc.font from Categories as c
		left join ColorCode as cc on c.codeID = cc.codeID
		where c.parentID IS NULL or c.parentID = 0
		and isLifestyle = 0
		order by c.sort`
	SubCategoriesStmt = `
		select c.catID, c.parentID, c.sort, c.dateAdded,
		c.catTitle, c.shortDesc, c.longDesc,
		c.image, c.icon, c.isLifestyle, c.vehicleSpecific,
		cc.code, cc.font from Categories as c
		left join ColorCode as cc on c.codeID = cc.codeID
		where c.parentID = ?
		and isLifestyle = 0
		order by c.sort`
	CategoryByNameStmt = `
		select c.catID, c.parentID, c.sort, c.dateAdded,
		c.catTitle, c.shortDesc, c.longDesc,
		c.image, c.icon, c.isLifestyle, c.vehicleSpecific,
		cc.code, cc.font from Categories as c
		left join ColorCode as cc on c.codeID = cc.codeID
		where c.catTitle = ?
		order by c.sort`
	CategoryByIdStmt = `
		select c.catID, c.parentID, c.sort, c.dateAdded,
		c.catTitle, c.shortDesc, c.longDesc,
		c.image, c.icon, c.isLifestyle, c.vehicleSpecific,
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
	SubCategoryIdStmt = `
		select c.catID, group_concat(p.partID) as parts from Categories as c
		left join CatPart as cp on c.catID = cp.catID
		left join Part as p on cp.partID = p.partID
		where c.parentID = ? && (p.status = null || (p.status = 800 || p.status = 900))`
	CategoryContentStmt = `
		select ct.type, c.text from ContentBridge cb
		join Content as c on cb.contentID = c.contentID
		left join ContentType as ct on c.cTypeID = ct.cTypeID
		where cb.catID = ?`
)

type Category struct {
	CategoryId, ParentId, Sort   int
	DateAdded                    time.Time
	Title, ShortDesc, LongDesc   string
	ColorCode, FontCode          string
	Image, Icon                  *url.URL
	IsLifestyle, VehicleSpecific bool
	Content                      []Content
}

type ExtendedCategory struct {

	// Replicate of the Category struct
	CategoryId, ParentId, Sort   int
	DateAdded                    time.Time
	Title, ShortDesc, LongDesc   string
	ColorCode, FontCode          string
	Image, Icon                  *url.URL
	IsLifestyle, VehicleSpecific bool

	// Extension for more detail
	SubCategories []Category
	Content       []Content
	Parts         []Part
}

func PopulateExtendedCategoryMulti(rows *sql.Rows, ch chan []ExtendedCategory) {
	cats := make([]ExtendedCategory, 0)
	if rows == nil {
		ch <- cats
		return
	}

	for rows.Next() {
		var initCat ExtendedCategory
		var catImg *string
		var catIcon *string
		var colorCode *string
		var fontCode *string
		err := rows.Scan(
			&initCat.CategoryId,
			&initCat.ParentId,
			&initCat.Sort,
			&initCat.DateAdded,
			&initCat.Title,
			&initCat.ShortDesc,
			&initCat.LongDesc,
			&catImg,
			&catIcon,
			&initCat.IsLifestyle,
			&initCat.VehicleSpecific,
			&colorCode,
			&fontCode)
		if err != nil {
			log.Println(err)
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
		cats = append(cats, initCat)
	}

	ch <- cats
}

func PopulateExtendedCategory(row *sql.Row, ch chan ExtendedCategory) {
	if row == nil {
		ch <- ExtendedCategory{}
		return
	}

	var initCat ExtendedCategory
	var catImg *string
	var catIcon *string
	var colorCode *string
	var fontCode *string
	err := row.Scan(
		&initCat.CategoryId,
		&initCat.ParentId,
		&initCat.Sort,
		&initCat.DateAdded,
		&initCat.Title,
		&initCat.ShortDesc,
		&initCat.LongDesc,
		&catImg,
		&catIcon,
		&initCat.IsLifestyle,
		&initCat.VehicleSpecific,
		&colorCode,
		&fontCode)
	if err != nil {
		log.Println(err)
		ch <- ExtendedCategory{}
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

	ch <- initCat
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
			&initCat.CategoryId,
			&initCat.ParentId,
			&initCat.Sort,
			&initCat.DateAdded,
			&initCat.Title,
			&initCat.ShortDesc,
			&initCat.LongDesc,
			&catImg,
			&catIcon,
			&initCat.IsLifestyle,
			&initCat.VehicleSpecific,
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

		con, err := initCat.GetContent()
		if err == nil {
			initCat.Content = con
		}

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
		&initCat.CategoryId,
		&initCat.ParentId,
		&initCat.Sort,
		&initCat.DateAdded,
		&initCat.Title,
		&initCat.ShortDesc,
		&initCat.LongDesc,
		&catImg,
		&catIcon,
		&initCat.IsLifestyle,
		&initCat.VehicleSpecific,
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

	ch <- initCat
}

// TopTierCategories
// Description: Returns the top tier categories
// Returns: []Category, error
func TopTierCategories() (cats []Category, err error) {

	redis_key := "category:top"

	// First lets try to access the category:top endpoint in Redis
	data, err := redis.Get(redis_key)
	if len(data) > 0 && err != nil {
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

	ch := make(chan []Category, 0)
	go PopulateCategoryMulti(catRows, ch)
	cats = <-ch

	go redis.Setex(redis_key, cats, 86400)

	return
}

func GetCategoryByTitle(cat_title string) (cat Category, err error) {

	redis_key := "category:title:" + cat_title

	// Attempt to get the category from Redis
	data, err := redis.Get(redis_key)
	if len(data) > 0 && err != nil {
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
	if len(data) > 0 && err != nil {
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

func (c *Category) SubCategories() (cats []Category, err error) {

	if c.CategoryId == 0 {
		return
	}

	redis_key := "category:" + strconv.Itoa(c.CategoryId) + ":subs"

	// First lets try to access the category:top endpoint in Redis
	data, err := redis.Get(redis_key)
	if len(data) > 0 && err != nil {
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
	catRows, err := qry.Query(c.CategoryId)
	if err != nil || catRows == nil { // Error occurred while executing query
		return
	}

	ch := make(chan []Category, 0)
	go PopulateCategoryMulti(catRows, ch)
	cats = <-ch

	go redis.Setex(redis_key, cats, 86400)

	return
}

func (c Category) GetCategory(key string) (extended ExtendedCategory, err error) {

	redis_key := "gopapi:category:" + strconv.Itoa(c.CategoryId)

	// First lets try to access the category:top endpoint in Redis
	cat_bytes, err := redis.Get(redis_key)
	if len(cat_bytes) > 0 {
		err = json.Unmarshal(cat_bytes, &extended)
		if err == nil {
			content, err := custcontent.GetCategoryContent(extended.CategoryId, key)
			for _, con := range content {
				strArr := strings.Split(con.ContentType.Type, ":")
				cType := con.ContentType.Type
				if len(strArr) > 1 {
					cType = strArr[1]
				}
				extended.Content = append(extended.Content, Content{
					Key:   cType,
					Value: con.Text,
				})
			}
			return extended, err
		}
	}

	var errs []error
	catChan := make(chan int)
	subChan := make(chan int)
	conChan := make(chan int)

	// Build out generalized category properties
	go func() {
		cat, catErr := GetCategoryById(c.CategoryId)

		if catErr != nil {
			errs = append(errs, catErr)
		} else {
			extended.CategoryId = cat.CategoryId
			extended.ColorCode = cat.ColorCode
			extended.DateAdded = cat.DateAdded
			extended.FontCode = cat.FontCode
			extended.Image = cat.Image
			extended.Icon = cat.Icon
			extended.IsLifestyle = cat.IsLifestyle
			extended.LongDesc = cat.LongDesc
			extended.ParentId = cat.ParentId
			extended.ShortDesc = cat.ShortDesc
			extended.Sort = cat.Sort
			extended.Title = cat.Title
			extended.VehicleSpecific = cat.VehicleSpecific
		}

		catChan <- 1
	}()

	go func() {
		subs, subErr := c.SubCategories()
		extended.SubCategories = subs
		if subErr != nil {
			errs = append(errs, subErr)
		}
		subChan <- 1
	}()

	go func() {
		cons, conErr := c.GetContent()
		if conErr != nil {
			errs = append(errs, conErr)
		} else {
			extended.Content = cons
		}
		conChan <- 1
	}()

	<-catChan
	<-subChan
	<-conChan

	if len(errs) > 1 {
		err = errs[0]
	} else if extended.CategoryId == 0 {
		return extended, errors.New("Invalid Category")
	}

	go redis.Setex(redis_key, extended, 86400)

	content, err := custcontent.GetCategoryContent(extended.CategoryId, key)
	for _, con := range content {
		strArr := strings.Split(con.ContentType.Type, ":")
		cType := con.ContentType.Type
		if len(strArr) > 1 {
			cType = strArr[1]
		}
		extended.Content = append(extended.Content, Content{
			Key:   cType,
			Value: con.Text,
		})
	}

	return
}

func (c *Category) GetContent() (content []Content, err error) {

	if c.CategoryId == 0 {
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

	// Execute SQL Query against current CategoryId
	conRows, err := qry.Query(c.CategoryId)
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
