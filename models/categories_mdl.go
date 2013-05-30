package models

import (
	"../helpers/database"
	"../helpers/mymysql/mysql"
	"../helpers/redis"
	"encoding/json"
	"errors"
	"log"
	"net/url"
	"strconv"
	"time"
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

// PartBreacrumbs
//
// Description: Builds out Category breadcrumb array for the current part object.
//
// Inherited: part Part
// Returns: error
func (part *Part) PartBreadcrumbs() error {

	if part.PartId == 0 {
		return errors.New("Invalid Part Number")
	}

	qry, err := database.GetStatement("PartCategoryStmt")
	if err != nil {
		return err
	}

	parentQuery, err := database.GetStatement("ParentCategoryStmt")
	if err != nil {
		return err
	}

	// Execute SQL Query against current PartId
	catRow, catRes, err := qry.ExecFirst(part.PartId)
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
			catRow, catRes, err = parentQuery.ExecFirst(parent)
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

// func (lookup *Lookup) PartBreadcrumbs() error {

// 	var ids []string
// 	for _, p := range lookup.Parts {
// 		ids = append(ids, strconv.Itoa(p.PartId))
// 	}

// 	rows, res, err := database.Db.Query(partCategoryStmt_Grouped, strings.Join(ids, ","))
// 	if database.MysqlError(err) {
// 		return err
// 	} else if len(rows) == 0 {
// 		return nil
// 	}

// }

func (part *Part) GetPartCategories() (cats []ExtendedCategory, err error) {

	if part.PartId == 0 {
		return
	}

	qry, err := database.GetStatement("PartAllCategoryStmt")
	if err != nil {
		return
	}

	// Execute SQL Query against current PartId
	catRows, catRes, err := qry.Exec(part.PartId)
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

// TopTierCategories
// Description: Returns the top tier categories
// Returns: []Category, error
func TopTierCategories() (cats []Category, err error) {

	// First lets try to access the category:top endpoint in Redis
	cat_bytes, err := redis.RedisClient.Get("category:top")
	if err == nil && len(cat_bytes) > 0 {
		err = json.Unmarshal(cat_bytes, &cats)
		if err == nil {
			return
		}
	}

	qry, err := database.GetStatement("TopCategoriesStmt")
	if err != nil {
		return
	}

	// Execute SQL Query against current PartId
	catRows, catRes, err := qry.Exec()
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

	if cat_bytes, err := json.Marshal(cats); err == nil {
		err = redis.RedisClient.Setex("category:top", 86400, cat_bytes)
		log.Println("Top Tier Category Submission: %s", err)
	}

	return
}

func GetByTitle(cat_title string) (cat Category, err error) {

	qry, err := database.GetStatement("CategoryByNameStmt")
	if err != nil {
		return
	}

	// Execute SQL Query against title
	catRow, catRes, err := qry.ExecFirst(cat_title)
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

	qry, err := database.GetStatement("CategoryByIdStmt")
	if err != nil {
		return
	}

	// Execute SQL Query against title
	catRow, catRes, err := qry.ExecFirst(cat_id)
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

	if c.CategoryId == 0 {
		return
	}

	// First lets try to access the category:top endpoint in Redis
	cat_bytes, err := redis.RedisClient.Get("category:" + strconv.Itoa(c.CategoryId) + ":subs")
	if err == nil && len(cat_bytes) > 0 {
		err = json.Unmarshal(cat_bytes, &cats)
		if err == nil {
			return
		}
	}

	qry, err := database.GetStatement("SubCategoriesStmt")
	if err != nil {
		return
	}

	// Execute SQL Query against current PartId
	catRows, catRes, err := qry.Exec(c.CategoryId)
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

	if cat_bytes, err = json.Marshal(cats); err == nil {
		cat_key := "category:" + strconv.Itoa(c.CategoryId) + ":subs"
		redis.RedisClient.Set(cat_key, cat_bytes)
		redis.RedisClient.Expire(cat_key, 86400)
	}

	return
}

func (c *Category) GetCategoryParts(key string, page int, count int) (parts []Part, err error) {

	if c.CategoryId == 0 {
		return
	}

	if page > 0 {
		page = count * page
	}

	qry, err := database.GetStatement("CategoryPartBasicStmt")
	if err != nil {
		return
	}

	rows, _, err := qry.Exec(c.CategoryId, page, count)
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

	// First lets try to access the category:top endpoint in Redis
	cat_bytes, err := redis.RedisClient.Get("category:" + strconv.Itoa(c.CategoryId))
	if err == nil && len(cat_bytes) > 0 {
		err = json.Unmarshal(cat_bytes, &extended)
		if err == nil {
			return
		}
	}

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
	} else if extended.CategoryId == 0 {
		return extended, errors.New("Invalid Category")
	}

	if cat_bytes, err := json.Marshal(extended); err == nil {
		cat_key := "category:" + strconv.Itoa(c.CategoryId)
		redis.RedisClient.Set(cat_key, cat_bytes)
		redis.RedisClient.Expire(cat_key, 86400)
	}

	return
}

func (c *Category) GetContent() (content []Content, err error) {

	if c.CategoryId == 0 {
		return
	}

	qry, err := database.GetStatement("CategoryContentStmt")
	if err != nil {
		return
	}

	// Execute SQL Query against current CategoryId
	conRows, res, err := qry.Exec(c.CategoryId)
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
