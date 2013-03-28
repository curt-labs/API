package models

import (
	"../helpers/database"
	"errors"
	"github.com/ziutek/mymysql/mysql"
	"net/url"
	"time"
)

var (

	// Get the category that a part is tied to, by PartId
	partCategoryStmt = `select c.catID, c.parentID, c.sort, c.dateAdded,
				c.catTitle, c.shortDesc, c.longDesc,
				c.image, c.isLifestyle, c.vehicleSpecific,
				cc.code, cc.font from Categories as c
				join CatPart as cp on c.catID = cp.catID
				left join ColorCode as cc on c.codeID = cc.codeID
				where cp.partID = %d
				order by c.sort
				limit 1`

	partAllCategoryStmt = `select c.catID, c.dateAdded, c.parentID, c.catTitle, c.shortDesc, 
					c.longDesc,c.sort, c.image, c.isLifestyle, c.vehicleSpecific,
					cc.font, cc.code
					from Categories as c
					join CatPart as cp on c.catID = cp.catID
					join ColorCode as cc on c.codeID = cc.codeID
					where cp.partID = %d
					order by c.catID`

	// Get a category by catID
	parentCategoryStmt = `select c.catID, c.parentID, c.sort, c.dateAdded,
					c.catTitle, c.shortDesc, c.longDesc,
					c.image, c.isLifestyle, c.vehicleSpecific,
					cc.code, cc.font from Categories as c
					left join ColorCode as cc on c.codeID = cc.codeID
					where c.catID = %d
					order by c.sort
					limit 1`

	// Get the top-tier categories i.e Hitches, Electrical
	topCategoriesStmt = `select c.catID, c.parentID, c.sort, c.dateAdded,
					c.catTitle, c.shortDesc, c.longDesc,
					c.image, c.isLifestyle, c.vehicleSpecific,
					cc.code, cc.font from Categories as c
					left join ColorCode as cc on c.codeID = cc.codeID
					where c.parentID IS NULL or c.parentID = 0
					and isLifestyle = 0
					order by c.sort`

	subCategoriesStmt = `select c.catID, c.parentID, c.sort, c.dateAdded,
					c.catTitle, c.shortDesc, c.longDesc,
					c.image, c.isLifestyle, c.vehicleSpecific,
					cc.code, cc.font from Categories as c
					left join ColorCode as cc on c.codeID = cc.codeID
					where c.parentID = %d
					and isLifestyle = 0
					order by c.sort`

	categoryByNameStmt = `select c.catID, c.parentID, c.sort, c.dateAdded,
					c.catTitle, c.shortDesc, c.longDesc,
					c.image, c.isLifestyle, c.vehicleSpecific,
					cc.code, cc.font from Categories as c
					left join ColorCode as cc on c.codeID = cc.codeID
					where c.catTitle = '%s'
					order by c.sort`

	categoryByIdStmt = `select c.catID, c.parentID, c.sort, c.dateAdded,
					c.catTitle, c.shortDesc, c.longDesc,
					c.image, c.isLifestyle, c.vehicleSpecific,
					cc.code, cc.font from Categories as c
					left join ColorCode as cc on c.codeID = cc.codeID
					where c.catID = %d
					order by c.sort`

	categoryPartBasicStmt = `select cp.partID
					from CatPart as cp
					where cp.catID = %d 
					order by cp.partID
					limit %d,%d`

	categoryContentStmt = `select ct.type, c.text from ContentBridge cb
					join Content as c on cb.contentID = c.contentID
					left join ContentType as ct on c.cTypeID = ct.cTypeID
					where cb.catID = %d`
)

type Category struct {
	CategoryId, ParentId, Sort   int
	DateAdded                    time.Time
	Title, ShortDesc, LongDesc   string
	ColorCode, FontCode          string
	Image                        *url.URL
	IsLifestyle, VehicleSpecific bool
}

type ExtendedCategory struct {

	// Replicate of the Category struct
	CategoryId, ParentId, Sort   int
	DateAdded                    time.Time
	Title, ShortDesc, LongDesc   string
	ColorCode, FontCode          string
	Image                        *url.URL
	IsLifestyle, VehicleSpecific bool

	// Extension for more detail
	SubCategories []Category
	Content       []Content
}

/*
* PartBreacrumbs
*
* Description: Builds out Category breadcrumb array for the current
*  		part object.
*
*  Inherited: part Part
*
* Retruns: error
 */
func (part *Part) PartBreadcrumbs() error {

	// Execute SQL Query against current PartId
	catRow, catRes, err := database.Db.QueryFirst(partCategoryStmt, part.PartId)
	if database.MysqlError(err) { // Error occurred while executing query
		return err
	} else if catRow == nil {
		return errors.New("No part found for " + string(part.PartId))
	}

	// Map the different columns to variables
	id := catRes.Map("catID")
	parent := catRes.Map("parentID")
	sort := catRes.Map("sort")
	date := catRes.Map("dateAdded")
	title := catRes.Map("catTitle")
	sDesc := catRes.Map("shortDesc")
	lDesc := catRes.Map("longDesc")
	img := catRes.Map("image")
	isLife := catRes.Map("isLifestyle")
	vSpecific := catRes.Map("vehicleSpecific")
	cCode := catRes.Map("code")
	font := catRes.Map("font")

	// Attempt to parse out the dataAdded field
	da, _ := time.Parse("2006-01-02 15:04:01", catRow.Str(date))

	// Attempt to parse out the image Url
	imgUrl, _ := url.Parse(catRow.Str(img))

	// Build out RGB value for color coding
	colorCode := catRow.Str(cCode)
	rgbCode := ""
	if len(colorCode) == 9 {
		rgbCode = "rgb(" + colorCode[0:3] + "," + colorCode[3:6] + "," + colorCode[6:9] + ")"
	}

	// Populate our lowest level Category
	initCat := ExtendedCategory{
		CategoryId:      catRow.Int(id),
		ParentId:        catRow.Int(parent),
		Sort:            catRow.Int(sort),
		DateAdded:       da,
		Title:           catRow.Str(title),
		ShortDesc:       catRow.Str(sDesc),
		LongDesc:        catRow.Str(lDesc),
		FontCode:        "#" + catRow.Str(font),
		Image:           imgUrl,
		IsLifestyle:     catRow.ForceBool(isLife),
		VehicleSpecific: catRow.ForceBool(vSpecific),
		ColorCode:       rgbCode,
	}

	// Instantiate our array with the initial category
	var cats []ExtendedCategory
	cats = append(cats, initCat)

	if initCat.ParentId > 0 { // Not top level category

		// Loop through the categories retrieving parents until we
		// hit the top-tier category
		parent := initCat.ParentId
		for {
			if parent == 0 {
				break
			}

			// Execute out SQL query to retrieve a category by ParentId
			catRow, catRes, err = db.QueryFirst(parentCategoryStmt, parent)
			if database.MysqlError(err) {
				return err
			}

			// Map the columns
			id := catRes.Map("catID")
			parentID := catRes.Map("parentID")
			sort := catRes.Map("sort")
			date := catRes.Map("dateAdded")
			title := catRes.Map("catTitle")
			sDesc := catRes.Map("shortDesc")
			lDesc := catRes.Map("longDesc")
			img := catRes.Map("image")
			isLife := catRes.Map("isLifestyle")
			vSpecific := catRes.Map("vehicleSpecific")
			cCode := catRes.Map("code")
			font := catRes.Map("font")

			// Parse the dateAdded field
			da, _ := time.Parse("2006-01-02 15:04:01", catRow.Str(date))

			// Parse the image Url
			imgUrl, _ := url.Parse(catRow.Str(img))

			// Build out RGB for color coding
			colorCode := catRow.Str(cCode)
			rgbCode := ""
			if len(colorCode) == 9 {
				rgbCode = "rgb(" + colorCode[0:3] + "," + colorCode[3:6] + "," + colorCode[6:9] + ")"
			}

			// Create Category object
			subCat := ExtendedCategory{
				CategoryId:      catRow.Int(id),
				ParentId:        catRow.Int(parentID),
				Sort:            catRow.Int(sort),
				DateAdded:       da,
				Title:           catRow.Str(title),
				ShortDesc:       catRow.Str(sDesc),
				LongDesc:        catRow.Str(lDesc),
				FontCode:        "#" + catRow.Str(font),
				Image:           imgUrl,
				IsLifestyle:     catRow.ForceBool(isLife),
				VehicleSpecific: catRow.ForceBool(vSpecific),
				ColorCode:       rgbCode,
			}

			// Append new Category onto array
			cats = append(cats, subCat)
			parent = subCat.ParentId
		}
	}

	// Apply breadcrumbs to our part object and return
	part.Categories = cats
	return nil
}

func (part *Part) GetPartCategories() (cats []ExtendedCategory, err error) {
	// Execute SQL Query against current PartId
	catRows, catRes, err := database.Db.Query(partAllCategoryStmt, part.PartId)
	if database.MysqlError(err) || catRows == nil { // Error occurred while executing query
		return
	}

	// Map the different columns to variables
	id := catRes.Map("catID")
	parent := catRes.Map("parentID")
	sort := catRes.Map("sort")
	date := catRes.Map("dateAdded")
	title := catRes.Map("catTitle")
	sDesc := catRes.Map("shortDesc")
	lDesc := catRes.Map("longDesc")
	img := catRes.Map("image")
	isLife := catRes.Map("isLifestyle")
	vSpecific := catRes.Map("vehicleSpecific")
	cCode := catRes.Map("code")
	font := catRes.Map("font")

	for _, catRow := range catRows {

		// Attempt to parse out the dataAdded field
		da, _ := time.Parse("2006-01-02 15:04:01", catRow.Str(date))

		// Attempt to parse out the image Url
		imgUrl, _ := url.Parse(catRow.Str(img))

		// Build out RGB value for color coding
		colorCode := catRow.Str(cCode)
		rgbCode := ""
		if len(colorCode) == 9 {
			rgbCode = "rgb(" + colorCode[0:3] + "," + colorCode[3:6] + "," + colorCode[6:9] + ")"
		}

		// Populate our lowest level Category
		cat := ExtendedCategory{
			CategoryId:      catRow.Int(id),
			ParentId:        catRow.Int(parent),
			Sort:            catRow.Int(sort),
			DateAdded:       da,
			Title:           catRow.Str(title),
			ShortDesc:       catRow.Str(sDesc),
			LongDesc:        catRow.Str(lDesc),
			FontCode:        "#" + catRow.Str(font),
			Image:           imgUrl,
			IsLifestyle:     catRow.ForceBool(isLife),
			VehicleSpecific: catRow.ForceBool(vSpecific),
			ColorCode:       rgbCode,
		}

		contentChan := make(chan int)
		subChan := make(chan int)

		c := Category{
			CategoryId: cat.CategoryId,
		}

		go func() {
			content, contentErr := c.GetContent()
			if contentErr == nil {
				cat.Content = content
			}
			contentChan <- 1
		}()

		go func() {
			subs, subErr := c.SubCategories()
			if subErr == nil {
				cat.SubCategories = subs
			}
			subChan <- 1
		}()

		<-contentChan
		<-subChan

		cats = append(cats, cat)
	}

	return
}

func TopTierCategories() (cats []Category, err error) {

	// Execute SQL Query against current PartId
	catRows, catRes, err := database.Db.Query(topCategoriesStmt)
	if database.MysqlError(err) || catRows == nil { // Error occurred while executing query
		return
	}

	// Map the different columns to variables
	id := catRes.Map("catID")
	parent := catRes.Map("parentID")
	sort := catRes.Map("sort")
	date := catRes.Map("dateAdded")
	title := catRes.Map("catTitle")
	sDesc := catRes.Map("shortDesc")
	lDesc := catRes.Map("longDesc")
	img := catRes.Map("image")
	isLife := catRes.Map("isLifestyle")
	vSpecific := catRes.Map("vehicleSpecific")
	cCode := catRes.Map("code")
	font := catRes.Map("font")

	for _, catRow := range catRows {
		// Attempt to parse out the dataAdded field
		da, _ := time.Parse("2006-01-02 15:04:01", catRow.Str(date))

		// Attempt to parse out the image Url
		imgUrl, _ := url.Parse(catRow.Str(img))

		// Build out RGB value for color coding
		colorCode := catRow.Str(cCode)
		rgbCode := ""
		if len(colorCode) == 9 {
			rgbCode = "rgb(" + colorCode[0:3] + "," + colorCode[3:6] + "," + colorCode[6:9] + ")"
		}

		// Populate our lowest level Category
		cat := Category{
			CategoryId:      catRow.Int(id),
			ParentId:        catRow.Int(parent),
			Sort:            catRow.Int(sort),
			DateAdded:       da,
			Title:           catRow.Str(title),
			ShortDesc:       catRow.Str(sDesc),
			LongDesc:        catRow.Str(lDesc),
			FontCode:        "#" + catRow.Str(font),
			Image:           imgUrl,
			IsLifestyle:     catRow.ForceBool(isLife),
			VehicleSpecific: catRow.ForceBool(vSpecific),
			ColorCode:       rgbCode,
		}
		cats = append(cats, cat)
	}

	return
}

func GetByTitle(cat_title string) (cat Category, err error) {

	// Execute SQL Query against title
	catRow, catRes, err := database.Db.QueryFirst(categoryByNameStmt, cat_title)
	if database.MysqlError(err) || catRow == nil { // Error occurred while executing query
		return
	}

	// Map the different columns to variables
	id := catRes.Map("catID")
	parent := catRes.Map("parentID")
	sort := catRes.Map("sort")
	date := catRes.Map("dateAdded")
	title := catRes.Map("catTitle")
	sDesc := catRes.Map("shortDesc")
	lDesc := catRes.Map("longDesc")
	img := catRes.Map("image")
	isLife := catRes.Map("isLifestyle")
	vSpecific := catRes.Map("vehicleSpecific")
	cCode := catRes.Map("code")
	font := catRes.Map("font")

	// Attempt to parse out the dataAdded field
	da, _ := time.Parse("2006-01-02 15:04:01", catRow.Str(date))

	// Attempt to parse out the image Url
	imgUrl, _ := url.Parse(catRow.Str(img))

	// Build out RGB value for color coding
	colorCode := catRow.Str(cCode)
	rgbCode := ""
	if len(colorCode) == 9 {
		rgbCode = "rgb(" + colorCode[0:3] + "," + colorCode[3:6] + "," + colorCode[6:9] + ")"
	}

	// Populate our lowest level Category
	cat = Category{
		CategoryId:      catRow.Int(id),
		ParentId:        catRow.Int(parent),
		Sort:            catRow.Int(sort),
		DateAdded:       da,
		Title:           catRow.Str(title),
		ShortDesc:       catRow.Str(sDesc),
		LongDesc:        catRow.Str(lDesc),
		FontCode:        "#" + catRow.Str(font),
		Image:           imgUrl,
		IsLifestyle:     catRow.ForceBool(isLife),
		VehicleSpecific: catRow.ForceBool(vSpecific),
		ColorCode:       rgbCode,
	}

	return
}

func GetById(cat_id int) (cat Category, err error) {

	// Execute SQL Query against title
	catRow, catRes, err := database.Db.QueryFirst(categoryByIdStmt, cat_id)
	if database.MysqlError(err) || catRow == nil { // Error occurred while executing query
		return
	}

	// Map the different columns to variables
	id := catRes.Map("catID")
	parent := catRes.Map("parentID")
	sort := catRes.Map("sort")
	date := catRes.Map("dateAdded")
	title := catRes.Map("catTitle")
	sDesc := catRes.Map("shortDesc")
	lDesc := catRes.Map("longDesc")
	img := catRes.Map("image")
	isLife := catRes.Map("isLifestyle")
	vSpecific := catRes.Map("vehicleSpecific")
	cCode := catRes.Map("code")
	font := catRes.Map("font")

	// Attempt to parse out the dataAdded field
	da, _ := time.Parse("2006-01-02 15:04:01", catRow.Str(date))

	// Attempt to parse out the image Url
	imgUrl, _ := url.Parse(catRow.Str(img))

	// Build out RGB value for color coding
	colorCode := catRow.Str(cCode)
	rgbCode := ""
	if len(colorCode) == 9 {
		rgbCode = "rgb(" + colorCode[0:3] + "," + colorCode[3:6] + "," + colorCode[6:9] + ")"
	}

	// Populate our lowest level Category
	cat = Category{
		CategoryId:      catRow.Int(id),
		ParentId:        catRow.Int(parent),
		Sort:            catRow.Int(sort),
		DateAdded:       da,
		Title:           catRow.Str(title),
		ShortDesc:       catRow.Str(sDesc),
		LongDesc:        catRow.Str(lDesc),
		FontCode:        "#" + catRow.Str(font),
		Image:           imgUrl,
		IsLifestyle:     catRow.ForceBool(isLife),
		VehicleSpecific: catRow.ForceBool(vSpecific),
		ColorCode:       rgbCode,
	}

	return
}

func (c *Category) SubCategories() (cats []Category, err error) {

	// Execute SQL Query against current PartId
	catRows, catRes, err := database.Db.Query(subCategoriesStmt, c.CategoryId)
	if database.MysqlError(err) || catRows == nil { // Error occurred while executing query
		return
	}

	// Map the different columns to variables
	id := catRes.Map("catID")
	parent := catRes.Map("parentID")
	sort := catRes.Map("sort")
	date := catRes.Map("dateAdded")
	title := catRes.Map("catTitle")
	sDesc := catRes.Map("shortDesc")
	lDesc := catRes.Map("longDesc")
	img := catRes.Map("image")
	isLife := catRes.Map("isLifestyle")
	vSpecific := catRes.Map("vehicleSpecific")
	cCode := catRes.Map("code")
	font := catRes.Map("font")

	for _, catRow := range catRows {
		// Attempt to parse out the dataAdded field
		da, _ := time.Parse("2006-01-02 15:04:01", catRow.Str(date))

		// Attempt to parse out the image Url
		imgUrl, _ := url.Parse(catRow.Str(img))

		// Build out RGB value for color coding
		colorCode := catRow.Str(cCode)
		rgbCode := ""
		if len(colorCode) == 9 {
			rgbCode = "rgb(" + colorCode[0:3] + "," + colorCode[3:6] + "," + colorCode[6:9] + ")"
		}

		// Populate our lowest level Category
		cat := Category{
			CategoryId:      catRow.Int(id),
			ParentId:        catRow.Int(parent),
			Sort:            catRow.Int(sort),
			DateAdded:       da,
			Title:           catRow.Str(title),
			ShortDesc:       catRow.Str(sDesc),
			LongDesc:        catRow.Str(lDesc),
			FontCode:        "#" + catRow.Str(font),
			Image:           imgUrl,
			IsLifestyle:     catRow.ForceBool(isLife),
			VehicleSpecific: catRow.ForceBool(vSpecific),
			ColorCode:       rgbCode,
		}
		cats = append(cats, cat)
	}

	return
}

func (c *Category) GetCategoryParts(key string, page int, count int) (parts []Part, err error) {

	if page > 0 {
		page = count * page
	}

	rows, _, err := database.Db.Query(categoryPartBasicStmt, c.CategoryId, page, count)
	if database.MysqlError(err) || rows == nil {
		return
	}

	chans := make(chan int, len(rows))

	for _, r := range rows {
		go func(row mysql.Row) {
			if len(row) == 1 {
				p := Part{
					PartId: row.Int(0),
				}
				p.Get(key)
				parts = append(parts, p)
				chans <- 1
			} else {
				chans <- 1
			}
		}(r)

	}

	for i := 0; i < len(rows); i++ {
		<-chans
	}
	return
}

func (c Category) GetCategory(key string) (extended ExtendedCategory, err error) {

	var errs []error
	catChan := make(chan int)
	subChan := make(chan int)
	conChan := make(chan int)
	// partChan := make(chan int)

	// Build out generalized category properties
	go func() {
		cat, catErr := GetById(c.CategoryId)

		if catErr != nil {
			errs = append(errs, catErr)
		} else {
			extended.CategoryId = cat.CategoryId
			extended.ColorCode = cat.ColorCode
			extended.DateAdded = cat.DateAdded
			extended.FontCode = cat.FontCode
			extended.Image = cat.Image
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

	// go func() {
	// 	parts, partErr := c.GetCategoryParts(key)
	// 	paged := PagedParts{
	// 		Parts: parts,
	// 	}
	// 	extended.Parts = paged
	// 	if partErr != nil {
	// 		errs = append(errs, partErr)
	// 	}
	// 	partChan <- 1
	// }()

	<-catChan
	<-subChan
	<-conChan
	// <-partChan

	if len(errs) > 1 {
		err = errs[0]
	}

	return
}

func (c *Category) GetContent() (content []Content, err error) {

	// Execute SQL Query against current CategoryId
	conRows, res, err := database.Db.Query(categoryContentStmt, c.CategoryId)
	if database.MysqlError(err) || conRows == nil {
		return
	}

	typ := res.Map("type")
	text := res.Map("text")

	for _, conRow := range conRows {
		con := Content{
			Key:   conRow.Str(typ),
			Value: conRow.Str(text),
		}
		content = append(content, con)
	}
	return
}
