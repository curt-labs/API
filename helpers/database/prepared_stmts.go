package database

import (
	//"../mymysql/autorc"
	"../mymysql/mysql"
	_ "../mymysql/thrsafe"
	"errors"
	"expvar"
)

// Category Raw Queries
var (

	// Get the category that a part is tied to, by PartId
	PartCategoryStmt = `select c.catID, c.parentID, c.sort, c.dateAdded,
				c.catTitle, c.shortDesc, c.longDesc,
				c.image, c.isLifestyle, c.vehicleSpecific,
				cc.code, cc.font from Categories as c
				join CatPart as cp on c.catID = cp.catID
				left join ColorCode as cc on c.codeID = cc.codeID
				where cp.partID = ?
				order by c.sort
				limit 1`

	PartAllCategoryStmt = `select c.catID, c.dateAdded, c.parentID, c.catTitle, c.shortDesc, 
					c.longDesc,c.sort, c.image, c.isLifestyle, c.vehicleSpecific,
					cc.font, cc.code
					from Categories as c
					join CatPart as cp on c.catID = cp.catID
					join ColorCode as cc on c.codeID = cc.codeID
					where cp.partID = ?
					order by c.catID`

	// Get a category by catID
	ParentCategoryStmt = `select c.catID, c.parentID, c.sort, c.dateAdded,
					c.catTitle, c.shortDesc, c.longDesc,
					c.image, c.isLifestyle, c.vehicleSpecific,
					cc.code, cc.font from Categories as c
					left join ColorCode as cc on c.codeID = cc.codeID
					where c.catID = ?
					order by c.sort
					limit 1`

	// Get the top-tier categories i.e Hitches, Electrical
	TopCategoriesStmt = `select c.catID, c.parentID, c.sort, c.dateAdded,
					c.catTitle, c.shortDesc, c.longDesc,
					c.image, c.isLifestyle, c.vehicleSpecific,
					cc.code, cc.font from Categories as c
					left join ColorCode as cc on c.codeID = cc.codeID
					where c.parentID IS NULL or c.parentID = 0
					and isLifestyle = 0
					order by c.sort`

	SubCategoriesStmt = `select c.catID, c.parentID, c.sort, c.dateAdded,
					c.catTitle, c.shortDesc, c.longDesc,
					c.image, c.isLifestyle, c.vehicleSpecific,
					cc.code, cc.font from Categories as c
					left join ColorCode as cc on c.codeID = cc.codeID
					where c.parentID = ?
					and isLifestyle = 0
					order by c.sort`

	CategoryByNameStmt = `select c.catID, c.parentID, c.sort, c.dateAdded,
					c.catTitle, c.shortDesc, c.longDesc,
					c.image, c.isLifestyle, c.vehicleSpecific,
					cc.code, cc.font from Categories as c
					left join ColorCode as cc on c.codeID = cc.codeID
					where c.catTitle = ?
					order by c.sort`

	CategoryByIdStmt = `select c.catID, c.parentID, c.sort, c.dateAdded,
					c.catTitle, c.shortDesc, c.longDesc,
					c.image, c.isLifestyle, c.vehicleSpecific,
					cc.code, cc.font from Categories as c
					left join ColorCode as cc on c.codeID = cc.codeID
					where c.catID = ?
					order by c.sort`

	CategoryPartBasicStmt = `select cp.partID
					from CatPart as cp
					where cp.catID = ?
					order by cp.partID
					limit ?,?`

	CategoryContentStmt = `select ct.type, c.text from ContentBridge cb
					join Content as c on cb.contentID = c.contentID
					left join ContentType as ct on c.cTypeID = ct.cTypeID
					where cb.catID = ?`
)

// Category Prepared Statements
var (
	Statements map[string]mysql.Stmt
)

// Prepare all MySQL statements
func PrepareAll() error {

	Statements = make(map[string]mysql.Stmt, 0)

	if !Db.IsConnected() {
		Db.Connect()
	}

	partCategoryPrepared, err := Db.Prepare(PartCategoryStmt)
	if err != nil {
		return err
	}
	Statements["PartCategoryStmt"] = partCategoryPrepared

	partAllCategoryPrepared, err := Db.Prepare(PartAllCategoryStmt)
	if err != nil {
		return err
	}
	Statements["PartAllCategoryStmt"] = partAllCategoryPrepared

	parentCategoryPrepared, err := Db.Prepare(ParentCategoryStmt)
	if err != nil {
		return err
	}
	Statements["ParentCategoryStmt"] = parentCategoryPrepared

	topCategoriesPrepared, err := Db.Prepare(TopCategoriesStmt)
	if err != nil {
		return err
	}
	Statements["TopCategoriesStmt"] = topCategoriesPrepared

	subCategoriesPrepared, err := Db.Prepare(SubCategoriesStmt)
	if err != nil {
		return err
	}
	Statements["SubCategoriesStmt"] = subCategoriesPrepared

	categoryByNamePrepared, err := Db.Prepare(CategoryByNameStmt)
	if err != nil {
		return err
	}
	Statements["CategoryByNameStmt"] = categoryByNamePrepared

	categoryByIdPrepared, err := Db.Prepare(CategoryByIdStmt)
	if err != nil {
		return err
	}
	Statements["CategoryByIdStmt"] = categoryByIdPrepared

	categoryPartBasicPrepared, err := Db.Prepare(CategoryPartBasicStmt)
	if err != nil {
		return err
	}
	Statements["CategoryPartBasicStmt"] = categoryPartBasicPrepared

	categoryContentPrepared, err := Db.Prepare(CategoryContentStmt)
	if err != nil {
		return err
	}
	Statements["CategoryContentStmt"] = categoryContentPrepared

	return nil
}

func GetStatement(key string) (stmt mysql.Stmt, err error) {
	stmt, ok := Statements[key]
	if !ok {
		qry := expvar.Get(key)
		if qry == nil {
			err = errors.New("Invalid query reference")
		} else {
			stmt, err = Db.Prepare(qry.String())
		}
	}
	return

}
