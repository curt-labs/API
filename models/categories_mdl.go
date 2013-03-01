package models

import (
	"../helpers/database"
	"net/url"
	"time"
)

var (
	partCategoryStmt = `select c.catID, c.parentID, c.sort, c.dateAdded,
				c.catTitle, c.shortDesc, c.longDesc,
				c.image, c.isLifestyle, c.vehicleSpecific,
				cc.code, cc.font from Categories as c
				join CatPart as cp on c.catID = cp.catID
				left join ColorCode as cc on c.codeID = cc.codeID
				where cp.partID = %d
				order by c.sort
				limit 1`

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

func (part *Part) PartBreadcrumbs() error {
	//db := database.Db

	catRow, catRes, err := database.Db.QueryFirst(partCategoryStmt, part.PartId)
	if database.MysqlError(err) {
		return err
	}

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

	da, _ := time.Parse("2006-01-02 15:04:01", catRow.Str(date))
	imgUrl, _ := url.Parse(catRow.Str(img))

	colorCode := catRow.Str(cCode)
	rgbCode := ""
	if len(colorCode) == 9 {
		rgbCode = "rgb(" + colorCode[0:3] + "," + colorCode[3:6] + "," + colorCode[6:9] + ")"
	}

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

	var cats []Category
	cats = append(cats, initCat)

	if initCat.ParentId > 0 {
		parent := initCat.ParentId
		for {
			if parent == 0 {
				break
			}
			catRow, catRes, err = db.QueryFirst(parentCategoryStmt, parent)

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

			da, _ := time.Parse("2006-01-02 15:04:01", catRow.Str(date))
			imgUrl, _ := url.Parse(catRow.Str(img))

			colorCode := catRow.Str(cCode)
			rgbCode := ""
			if len(colorCode) == 9 {
				rgbCode = "rgb(" + colorCode[0:3] + "," + colorCode[3:6] + "," + colorCode[6:9] + ")"
			}

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

			cats = append(cats, subCat)
			parent = subCat.ParentId
		}
	}

	part.Categories = cats
	return nil
}
