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
	db := database.Db

	row, res, err := db.QueryFirst(partCategoryStmt, part.PartId)
	if database.MysqlError(err) {
		return err
	}

	id := res.Map("catID")
	parent := res.Map("parentID")
	sort := res.Map("sort")
	date := res.Map("dateAdded")
	title := res.Map("catTitle")
	sDesc := res.Map("shortDesc")
	lDesc := res.Map("longDesc")
	img := res.Map("image")
	isLife := res.Map("isLifestyle")
	vSpecific := res.Map("vehicleSpecific")
	cCode := res.Map("code")
	font := res.Map("font")

	da, _ := time.Parse("2006-01-02 15:04:01", row.Str(date))
	imgUrl, _ := url.Parse(row.Str(img))

	colorCode := row.Str(cCode)
	rgbCode := "rgb(" + colorCode[0:3] + "," + colorCode[3:6] + "," + colorCode[6:9] + ")"

	initCat := Category{
		CategoryId:      row.Int(id),
		ParentId:        row.Int(parent),
		Sort:            row.Int(sort),
		DateAdded:       da,
		Title:           row.Str(title),
		ShortDesc:       row.Str(sDesc),
		LongDesc:        row.Str(lDesc),
		FontCode:        "#" + row.Str(font),
		Image:           imgUrl,
		IsLifestyle:     row.ForceBool(isLife),
		VehicleSpecific: row.ForceBool(vSpecific),
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
			row, res, err = db.QueryFirst(parentCategoryStmt, parent)

			id := res.Map("catID")
			parentID := res.Map("parentID")
			sort := res.Map("sort")
			date := res.Map("dateAdded")
			title := res.Map("catTitle")
			sDesc := res.Map("shortDesc")
			lDesc := res.Map("longDesc")
			img := res.Map("image")
			isLife := res.Map("isLifestyle")
			vSpecific := res.Map("vehicleSpecific")
			cCode := res.Map("code")
			font := res.Map("font")

			da, _ := time.Parse("2006-01-02 15:04:01", row.Str(date))
			imgUrl, _ := url.Parse(row.Str(img))

			colorCode := row.Str(cCode)
			rgbCode := "rgb(" + colorCode[0:3] + "," + colorCode[3:6] + "," + colorCode[6:9] + ")"

			subCat := Category{
				CategoryId:      row.Int(id),
				ParentId:        row.Int(parentID),
				Sort:            row.Int(sort),
				DateAdded:       da,
				Title:           row.Str(title),
				ShortDesc:       row.Str(sDesc),
				LongDesc:        row.Str(lDesc),
				FontCode:        "#" + row.Str(font),
				Image:           imgUrl,
				IsLifestyle:     row.ForceBool(isLife),
				VehicleSpecific: row.ForceBool(vSpecific),
				ColorCode:       rgbCode,
			}

			cats = append(cats, subCat)
			parent = subCat.ParentId
		}
	}

	part.Categories = cats
	return nil
}
