package models

import (
	"../helpers/database"
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

	// Get a category by catID
	parentCategoryStmt = `select c.catID, c.parentID, c.sort, c.dateAdded,
					c.catTitle, c.shortDesc, c.longDesc,
					c.image, c.isLifestyle, c.vehicleSpecific,
					cc.code, cc.font from Categories as c
					left join ColorCode as cc on c.codeID = cc.codeID
					where c.catID = %d
					order by c.sort
					limit 1`
)

type Category struct {
	CategoryId, ParentId, Sort   int
	DateAdded                    time.Time
	Title, ShortDesc, LongDesc   string
	ColorCode, FontCode          string
	Image                        *url.URL
	IsLifestyle, VehicleSpecific bool
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
	initCat := Category{
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
	var cats []Category
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
			subCat := Category{
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
