package products

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/curt-labs/GoAPI/helpers/api"
	"github.com/curt-labs/GoAPI/helpers/database"
	"github.com/curt-labs/GoAPI/helpers/redis"
	"github.com/curt-labs/GoAPI/helpers/rest"
	"github.com/curt-labs/GoAPI/helpers/sortutil"
	"github.com/curt-labs/GoAPI/models/customer"
	"github.com/curt-labs/GoAPI/models/customer/content"
	"github.com/curt-labs/GoAPI/models/vehicle"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var (
	GetPaginatedPartNumbers = `select distinct p.partID
                               from Part as p
                               where p.status = 800 || p.status = 900
                               order by p.partID
                               limit ?,?`
	GetFeaturedParts = `select distinct p.partID
                        from Part as p
                        where (p.status = 800 || p.status = 900) && p.featured = 1
                        order by p.dateAdded desc
                        limit 0, ?`
	GetLatestParts = `select distinct p.partID
                      from Part as p
                      where p.status = 800 || p.status = 900
                      order by p.dateAdded desc
                      limit 0,?`
	SubCategoryIDStmt = `select distinct cp.partID
                         from CatPart as cp
                         join Part as p on cp.partID = p.partID
                         where cp.catID IN(%s) and (p.status = 800 || p.status = 900)
                         order by cp.partID
                         limit %d, %d`
	basicsStmt = `select p.status, p.dateAdded, p.dateModified, p.shortDesc, p.partID, p.priceCode, pc.class
                from Part as p
                left join Class as pc on p.classID = pc.classID
                where p.partID = ? && p.status in (800,900) limit 1`
	relatedPartStmt = `select distinct relatedID from RelatedPart
                where partID = ?
                order by relatedID`
	partContentStmt = `select ct.type, con.text
                from Content as con
                join ContentBridge as cb on con.contentID = cb.contentID
                join ContentType as ct on con.cTypeID = ct.cTypeID
                where cb.partID = ? && LOWER(ct.type) != 'appguide'
                order by ct.type`
	partInstallSheetStmt = `select c.text from ContentBridge as cb
                    join Content as c on cb.contentID = c.contentID
                    join ContentType as ct on c.cTypeID = ct.cTypeID
                    where partID = ? && ct.type = 'InstallationSheet'
                    limit 1`

	getPartByOldpartNumber = `select partID, status, dateModified, dateAdded, shortDesc, priceCode, classID, featured, ACESPartTypeID from Part where oldPartNumber = ?`
	//create
	createPart = `INSERT INTO Part (partID, status, dateAdded, shortDesc, priceCode, classID, featured, ACESPartTypeID)
                    VALUES(?,?,?,?,?,?,?, ?)`
	createPartAttributeJoin = `INSERT INTO PartAttribute (partID, value, field, sort) VALUES (?,?,?,?)`
	createVehiclePartJoin   = `INSERT INTO VehiclePart (vehicleID, partID, drilling, exposed, installTime) VALUES (?,?,?,?,?)`
	createContentBridge     = `INSERT INTO ContentBridge (catID, partID, contentID) VALUES (?,?,?)`
	createRelatedPart       = `INSERT INTO RelatedPart (partID, relatedID) VALUES (?,?)`
	createPartCategoryJoin  = `INSERT INTO CatPart (catID, partID) VALUES (?,?)`

	//delete
	deletePart               = `DELETE FROM Part WHERE partID  = ?`
	deletePartAttributeJoins = `DELETE FROM PartAttribute WHERE partID = ?`
	deleteVehiclePartJoins   = `DELETE FROM VehiclePart WHERE partID = ?`
	deleteContentBridgeJoins = `DELETE FROM ContentBridge WHERE partID = ?`
	deleteRelatedParts       = `DELETE FROM RelatedPart WHERE partID = ?`
	deletePartCategoryJoins  = `DELETE FROM CatPart WHERE partID = ?`

	//update
	updatePart = `UPDATE Part SET status = ?, shortDesc = ?, priceCode = ?, classID = ?, featured = ?, ACESPartTypeID = ? WHERE partID = ?`
)

type Part struct {
	ID                int               `json:"id" xml:"id,attr"`
	Status            int               `json:"status" xml:"status,attr"`
	PriceCode         int               `json:"price_code" xml:"price_code,attr"`
	RelatedCount      int               `json:"related_count" xml:"related_count,attr"`
	AverageReview     float64           `json:"average_review" xml:"average_review,attr"`
	DateModified      time.Time         `json:"date_modified" xml:"date_modified,attr"`
	DateAdded         time.Time         `json:"date_added" xml:"date_added,attr"`
	ShortDesc         string            `json:"short_description" xml:"short_description,attr"`
	PartClass         string            `json:"part_class" xml:"part_class,attr"` //sloppy - delete me in favor of child object "Class"
	InstallSheet      *url.URL          `json:"install_sheet" xml:"install_sheet"`
	Attributes        []Attribute       `json:"attributes" xml:"attributes"`
	VehicleAttributes []string          `json:"vehicle_atttributes" xml:"vehicle_attributes"`
	Vehicles          []vehicle.Vehicle `json:"vehicles,omitempty" xml:"vehicles,omitempty"`
	Content           []Content         `json:"content" xml:"content"`
	Pricing           []Price           `json:"pricing" xml:"pricing"`
	Reviews           []Review          `json:"reviews" xml:"reviews"`
	Images            []Image           `json:"images" xml:"images"`
	Related           []int             `json:"related" xml:"related"`
	Categories        []Category        `json:"categories" xml:"categories"`
	Videos            []PartVideo       `json:"videos" xml:"videos"`
	Packages          []Package         `json:"packages" xml:"packages"`
	Customer          CustomerPart      `json:"customer,omitempty" xml:"customer,omitempty"`
	Class             Class             `json:"class,omitempty" xml:"class,omitempty"`
	Featured          bool              `json:"featured,omitempty" xml:"featured,omitempty"`
	AcesPartTypeID    int               `json:"acesPartTypeId,omitempty" xml:"acesPartTypeId,omitempty"`
	Installations     Installations     `json:"installation,omitempty" xml:"installation,omitempty"`
	Inventory         PartInventory     `json:"inventory,omitempty" xml:"inventory,omitempty"`
	OldPartNumber     string            `json:"oldPartNumber,omitempty" xml:"oldPartNumber,omitempty"`
}

type CustomerPart struct {
	Price         float64 `json:"price" xml:"price,attr"`
	CartReference int     `json:"cart_reference" xml:"cart_reference,attr"`
}

type Content struct {
	ID    int    `json:"id,omitempty" xml:"id,omitempty"`
	Key   string `json:"key" xml:"key,attr"`
	Value string `json:"value" xml:",chardata"`
}

type PaginatedProductListing struct {
	Parts         []Part `json:"parts" xml:"parts"`
	TotalItems    int    `json:"total_items" xml:"total_items,attr"`
	ReturnedCount int    `json:"returned_count" xml:"returned_count,attr"`
	Page          int    `json:"page" xml:"page,attr"`
	PerPage       int    `json:"per_page" xml:"per_page,attr"`
	TotalPages    int    `json:"total_pages" xml:"total_pages,attr"`
}

type Class struct {
	ID    int    `json:"id,omitempty" xml:"id,omitempty"`
	Name  string `json:"name,omitempty" xml:"name,omitempty"`
	Image string `json:"image,omitempty" xml:"image,omitempty"`
}

type Installation struct { //aka VehiclePart Table
	ID          int             `json:"id,omitempty" xml:"id,omitempty"`
	Vehicle     vehicle.Vehicle `json:"vehicle,omitempty" xml:"vehicle,omitempty"`
	Part        Part            `json:"part,omitempty" xml:"part,omitempty"`
	Drilling    string          `json:"drilling,omitempty" xml:"v,omitempty"`
	Exposed     string          `json:"exposed,omitempty" xml:"exposed,omitempty"`
	InstallTime int             `json:"installTime,omitempty" xml:"installTime,omitempty"`
}

type Installations []Installation

func (p *Part) FromDatabase() error {

	var errs []string

	attrChan := make(chan int)
	priceChan := make(chan int)
	reviewChan := make(chan int)
	imageChan := make(chan int)
	videoChan := make(chan int)
	relatedChan := make(chan int)
	packageChan := make(chan int)
	categoryChan := make(chan int)
	contentChan := make(chan int)

	go func() {
		attrErr := p.GetAttributes()
		if attrErr != nil {
			errs = append(errs, attrErr.Error())
		}
		attrChan <- 1
	}()

	go func() {
		priceErr := p.GetPricing()
		if priceErr != nil {
			errs = append(errs, priceErr.Error())
		}
		priceChan <- 1
	}()

	go func() {
		reviewErr := p.GetActiveApprovedReviews()
		if reviewErr != nil {
			errs = append(errs, reviewErr.Error())
		}
		reviewChan <- 1
	}()

	go func() {
		imgErr := p.GetImages()
		if imgErr != nil {
			errs = append(errs, imgErr.Error())
		}
		imageChan <- 1
	}()

	go func() {
		vidErr := p.GetVideos()
		if vidErr != nil {
			errs = append(errs, vidErr.Error())
		}
		videoChan <- 1
	}()

	go func() {
		relErr := p.GetRelated()
		if relErr != nil {
			errs = append(errs, relErr.Error())
		}
		relatedChan <- 1
	}()

	go func() {
		pkgErr := p.GetPartPackaging()
		if pkgErr != nil {
			errs = append(errs, pkgErr.Error())
		}
		packageChan <- 1
	}()

	go func() {
		catErr := p.PartBreadcrumbs()
		if catErr != nil {
			errs = append(errs, catErr.Error())
		}
		categoryChan <- 1
	}()

	go func() {
		conErr := p.GetContent()
		if conErr != nil {
			errs = append(errs, conErr.Error())
		}
		contentChan <- 1
	}()

	if basicErr := p.Basics(); basicErr != nil {
		errs = append(errs, basicErr.Error())
	}

	<-attrChan
	<-priceChan
	<-reviewChan
	<-imageChan
	<-videoChan
	<-relatedChan
	<-packageChan
	<-categoryChan
	<-contentChan

	go func(tmp Part) {
		redis.Setex("part:"+strconv.Itoa(tmp.ID), tmp, redis.CacheTimeout)
	}(*p)

	return nil
}

func (p *Part) Get(key string) error {

	customerChan := make(chan int)

	var err error

	go func(api_key string) {
		err = p.BindCustomer(api_key)

		p.GetInventory(api_key, "")
		customerChan <- 1
	}(key)

	redis_key := fmt.Sprintf("part:%d", p.ID)

	part_bytes, err := redis.Get(redis_key)
	if len(part_bytes) > 0 {
		json.Unmarshal(part_bytes, &p)
	}

	if p.Status == 0 {
		if err := p.FromDatabase(); err != nil {
			customerChan <- 1
			close(customerChan)
			return err
		}
	}

	<-customerChan
	close(customerChan)

	return err
}

func All(key string, page, count int) ([]Part, error) {

	parts := make([]Part, 0)

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return parts, err
	}
	defer db.Close()

	stmt, err := db.Prepare(GetPaginatedPartNumbers)
	if err != nil {
		return parts, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(page, count)
	if err != nil {
		return parts, err
	}

	iter := 0
	partChan := make(chan int)
	for rows.Next() {
		var partID int
		if err = rows.Scan(&partID); err != nil {
			return parts, err
		}

		go func(id int) {
			p := Part{ID: id}
			p.Get(key)
			parts = append(parts, p)
			partChan <- 1
		}(partID)
		iter++
	}

	for i := 0; i < iter; i++ {
		<-partChan
	}

	sortutil.AscByField(parts, "ID")

	return parts, nil
}

func Featured(key string, count int) ([]Part, error) {
	parts := make([]Part, 0)

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return parts, err
	}
	defer db.Close()

	stmt, err := db.Prepare(GetFeaturedParts)
	if err != nil {
		return parts, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(count)
	if err != nil {
		return parts, err
	}

	iter := 0
	partChan := make(chan int)
	for rows.Next() {
		var partID int
		if err = rows.Scan(&partID); err != nil {
			return parts, err
		}

		go func(id int) {
			p := Part{ID: id}
			p.Get(key)
			parts = append(parts, p)
			partChan <- 1
		}(partID)
		iter++
	}

	for i := 0; i < iter; i++ {
		<-partChan
	}

	sortutil.DescByField(parts, "DateAdded")

	return parts, nil
}

func Latest(key string, count int) ([]Part, error) {
	parts := make([]Part, 0)

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return parts, err
	}
	defer db.Close()

	stmt, err := db.Prepare(GetLatestParts)
	if err != nil {
		return parts, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(count)
	if err != nil {
		return parts, err
	}

	iter := 0
	partChan := make(chan int)
	for rows.Next() {
		var partID int
		if err = rows.Scan(&partID); err != nil {
			return parts, err
		}

		go func(id int) {
			p := Part{ID: id}
			p.Get(key)
			parts = append(parts, p)
			partChan <- 1
		}(partID)
		iter++
	}

	for i := 0; i < iter; i++ {
		<-partChan
	}

	sortutil.DescByField(parts, "DateAdded")

	return parts, nil
}

func (p *Part) GetWithVehicle(vehicle *vehicle.Vehicle, api_key string) error {

	var errs []string

	superChan := make(chan int)
	noteChan := make(chan int)
	go func(key string) {
		p.Get(key)
		superChan <- 1
	}(api_key)
	go func() {
		notes, nErr := vehicle.GetNotes(p.ID)
		if nErr != nil && notes != nil {
			errs = append(errs, nErr.Error())
			p.VehicleAttributes = []string{}
		} else {
			p.VehicleAttributes = notes
		}
		noteChan <- 1
	}()

	<-superChan
	<-noteChan

	if len(errs) > 0 {
		return errors.New("Error: " + strings.Join(errs, ", "))
	}
	return nil
}

func (p *Part) GetById(id int, key string) {
	p.ID = id

	p.Get(key)
}

func (p *Part) Basics() error {
	redis_key := fmt.Sprintf("part:%d:basics", p.ID)

	data, err := redis.Get(redis_key)
	if err == nil && len(data) > 0 {
		if err = json.Unmarshal(data, &p); err != nil {
			return nil
		}
	}

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	qry, err := db.Prepare(basicsStmt)
	if err != nil {
		return err
	}
	defer qry.Close()

	row := qry.QueryRow(p.ID)
	if row == nil {
		return errors.New("No Part Found for:" + string(p.ID))
	}

	row.Scan(
		&p.Status,
		&p.DateAdded,
		&p.DateModified,
		&p.ShortDesc,
		&p.ID,
		&p.PriceCode,
		&p.PartClass)

	if !strings.Contains(p.ShortDesc, "CURT") {
		p.ShortDesc = fmt.Sprintf("CURT %s %d", p.ShortDesc, p.ID)
	}

	go func(tmp Part) {
		redis.Setex(redis_key, tmp, redis.CacheTimeout)
	}(*p)

	return nil
}

func (p *Part) GetRelated() error {
	redis_key := fmt.Sprintf("part:%d:related", p.ID)

	data, err := redis.Get(redis_key)
	if err == nil && len(data) > 0 {
		if err = json.Unmarshal(data, &p.Related); err != nil {
			return nil
		}
	}

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(relatedPartStmt)
	if err != nil {
		return err
	}
	defer stmt.Close()

	rows, err := stmt.Query(p.ID)
	if err != nil {
		return err
	}

	var related []int
	var relatedID int
	for rows.Next() {
		err = rows.Scan(&relatedID)
		if err != nil {
			return err
		}
		related = append(related, relatedID)
	}
	defer rows.Close()

	p.Related = related
	p.RelatedCount = len(related)

	go func(rel []int) {
		redis.Setex(redis_key, rel, redis.CacheTimeout)
	}(p.Related)

	return nil
}

func (p *Part) GetContent() error {
	redis_key := fmt.Sprintf("part:%d:content", p.ID)

	data, err := redis.Get(redis_key)
	if err == nil && len(data) > 0 {
		if err = json.Unmarshal(data, &p.Content); err != nil {
			return nil
		}
	}

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(partContentStmt)
	if err != nil {
		return err
	}
	defer stmt.Close()

	rows, err := stmt.Query(p.ID)
	if err != nil {
		return err
	}

	var content []Content
	for rows.Next() {
		var con Content
		err = rows.Scan(
			&con.Key,
			&con.Value,
		)
		if err != nil {
			return err
		}

		if strings.Contains(strings.ToLower(con.Key), "install") {
			//sheetUrl, _ := url.Parse(con.Value)
			p.InstallSheet, _ = url.Parse(api_helpers.API_DOMAIN + "/part/" + strconv.Itoa(p.ID) + ".pdf")
		} else {
			content = append(content, con)
		}
	}
	defer rows.Close()

	p.Content = content

	go redis.Setex(redis_key, p.Content, redis.CacheTimeout)

	return nil
}

func (p *Part) BindCustomer(key string) error {
	var price float64
	var ref int

	priceChan := make(chan int)
	refChan := make(chan int)
	contentChan := make(chan int)

	go func() {
		price, _ = customer.GetCustomerPrice(key, p.ID)
		priceChan <- 1
	}()

	go func() {
		ref, _ = customer.GetCustomerCartReference(key, p.ID)
		refChan <- 1
	}()

	go func() {
		content, _ := custcontent.GetPartContent(p.ID, key)
		for _, con := range content {

			strArr := strings.Split(con.ContentType.Type, ":")
			cType := con.ContentType.Type
			if len(strArr) > 1 {
				cType = strArr[1]
			}
			p.Content = append(p.Content, Content{
				Key:   cType,
				Value: con.Text,
			})
		}
		contentChan <- 1
	}()

	<-priceChan
	<-refChan
	<-contentChan

	cust := CustomerPart{
		Price:         price,
		CartReference: ref,
	}
	p.Customer = cust
	return nil
}

func (p *Part) GetInstallSheet(r *http.Request) (data []byte, err error) {
	redis_key := fmt.Sprintf("part:%d:installsheet", p.ID)

	data, err = redis.Get(redis_key)
	if err == nil && len(data) > 0 {
		return data, nil
	}

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return
	}
	defer db.Close()

	stmt, err := db.Prepare(partInstallSheetStmt)
	if err != nil {
		return
	}
	defer stmt.Close()

	var text string
	err = stmt.QueryRow(p.ID).Scan(
		&text,
	)
	if err != nil {
		return
	}

	data, err = rest.GetPDF(text, r)

	go func(dt []byte) {
		redis.Setex(redis_key, dt, redis.CacheTimeout)
	}(data)

	return
}

// PartBreacrumbs
//
// Description: Builds out Category breadcrumb array for the current part object.
//
// Inherited: part Part
// Returns: error
func (p *Part) PartBreadcrumbs() error {

	redis_key := fmt.Sprintf("part:%d:breadcrumbs", p.ID)

	data, err := redis.Get(redis_key)
	if err == nil && len(data) > 0 {
		if err = json.Unmarshal(data, &p.Categories); err == nil {
			return nil
		}
	}

	if p.ID == 0 {
		return errors.New("Invalid Part Number")
	}

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	qry, err := db.Prepare(PartCategoryStmt)
	if err != nil {
		return err
	}
	defer qry.Close()

	parentQuery, err := db.Prepare(ParentCategoryStmt)
	if err != nil {
		return err
	}
	defer parentQuery.Close()

	// Execute SQL Query against current ID
	catRow := qry.QueryRow(p.ID)
	if catRow == nil {
		return errors.New("No part found for " + string(p.ID))
	}

	ch := make(chan Category)
	go PopulateCategory(catRow, ch)
	initCat := <-ch

	// Instantiate our array with the initial category
	var cats []Category
	cats = append(cats, initCat)

	if initCat.ParentID > 0 { // Not top level category

		// Loop through the categories retrieving parents until we
		// hit the top-tier category
		parent := initCat.ParentID
		for {
			if parent == 0 {
				break
			}

			// Execute out SQL query to retrieve a category by ParentID
			catRow = parentQuery.QueryRow(parent)
			if catRow == nil {
				break
			}

			ch := make(chan Category)
			go PopulateCategory(catRow, ch)

			// Append new Category onto array
			subCat := <-ch
			cats = append(cats, subCat)
			parent = subCat.ParentID
		}
	}

	// Apply breadcrumbs to our part object and return
	p.Categories = cats

	go func(cats []Category) {
		redis.Setex(redis_key, cats, redis.CacheTimeout)
	}(p.Categories)

	return nil
}

func (p *Part) GetPartCategories(key string) (cats []Category, err error) {

	redis_key := fmt.Sprintf("part:%d:categories", p.ID)

	data, err := redis.Get(redis_key)
	if err == nil && len(data) > 0 {
		if err := json.Unmarshal(data, &cats); err == nil {
			return cats, nil
		}
	}

	if p.ID == 0 {
		return
	}

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return
	}
	defer db.Close()

	qry, err := db.Prepare(PartAllCategoryStmt)
	if err != nil {
		return
	}
	defer qry.Close()

	// Execute SQL Query against current ID
	catRows, err := qry.Query(p.ID)
	if err != nil || catRows == nil { // Error occurred while executing query
		return
	}

	ch := make(chan []Category, 0)
	go PopulateCategoryMulti(catRows, ch)
	cats = <-ch

	for _, cat := range cats {

		contentChan := make(chan int)
		subChan := make(chan int)
		customerChan := make(chan int)

		c := Category{
			ID: cat.ID,
		}

		go func() {
			content, contentErr := c.GetContent()
			if contentErr == nil {
				cat.Content = content
			}
			contentChan <- 1
		}()

		go func() {
			subs, subErr := c.GetSubCategories()
			if subErr == nil {
				cat.SubCategories = subs
			}
			subChan <- 1
		}()

		go func() {
			content, _ := custcontent.GetCategoryContent(cat.ID, key)
			for _, con := range content {
				strArr := strings.Split(con.ContentType.Type, ":")
				cType := con.ContentType.Type
				if len(strArr) > 1 {
					cType = strArr[1]
				}
				cat.Content = append(cat.Content, Content{
					Key:   cType,
					Value: con.Text,
				})
			}
			customerChan <- 1
		}()

		<-contentChan
		<-subChan
		<-customerChan

		cats = append(cats, cat)
	}

	go func(cts []Category) {
		redis.Setex(redis_key, cts, redis.CacheTimeout)
	}(cats)

	return
}

// func (tree *CategoryTree) CategoryTreeBuilder() {

//  db, err := sql.Open("mysql", database.ConnectionString())
//  if err != nil {
//      return
//  }
//  defer db.Close()

//  subQry, err := db.Prepare(SubIDStmt)
//  if err != nil {
//      return
//  }
//  defer subQry.Close()

//  // Execute against current Category Id
//  // to retrieve all category Ids that are children.
//  rows, err := subQry.Query(tree.ID)
//  if err != nil {
//      return
//  }

//  chans := make(chan int, 0)
//  var rowCount int
//  for rows.Next() {
//      var catID int
//      if err := rows.Scan(&catID); err != nil {
//          continue
//      }

//      go func(catID int) {

//          // Need to parse out string array into ints and populate
//          cat := Category{
//              ID: catID,
//          }
//          tree.SubCategories = append(tree.SubCategories, cat.ID)

//          subRows, err := subQry.Query(cat.ID)
//          if err == nil && subRows != nil {
//              for subRows.Next() {
//                  var subID int
//                  if err := subRows.Scan(&subID); err == nil {
//                      tree.SubCategories = append(tree.SubCategories, subID)
//                  }

//                  // subTree.CategoryTreeBuilder()
//                  // tree.SubCategories = append(tree.SubCategories, subTree.SubCategories...)
//              }
//          }
//          chans <- 1
//      }(catID)
//      rowCount++
//  }

//  for i := 0; i < rowCount; i++ {
//      <-chans
//  }

//  return
// }

func (p *Part) GetPartByOldPartNumber() (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	stmt, err := db.Prepare(getPartByOldpartNumber)
	if err != nil {
		return err
	}
	defer stmt.Close()
	err = stmt.QueryRow(p.OldPartNumber).Scan(
		&p.ID,
		&p.Status,
		&p.DateModified,
		&p.DateAdded,
		&p.OldPartNumber,
		&p.PriceCode,
		&p.Class.ID,
		&p.Featured,
		&p.AcesPartTypeID,
	)
	if err != sql.ErrNoRows {
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Part) Create() (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()
	stmt, err := tx.Prepare(createPart)
	if err != nil {
		return err
	}
	defer stmt.Close()
	p.DateAdded = time.Now()
	_, err = stmt.Exec(
		p.ID,
		p.Status,
		p.DateAdded,
		p.ShortDesc,
		p.PriceCode,
		p.Class.ID,
		p.Featured,
		p.AcesPartTypeID,
	)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()

	pajChan := make(chan int)
	diChan := make(chan int)
	dcbChan := make(chan int)
	priceChan := make(chan int)
	revChan := make(chan int)
	imageChan := make(chan int)
	relatedChan := make(chan int)
	pcjChan := make(chan int)
	videoChan := make(chan int)
	packChan := make(chan int)

	go func() (err error) {
		for _, attribute := range p.Attributes {
			err = p.CreatePartAttributeJoin(attribute)
			if err != nil {
				pajChan <- 1
				return err
			}
		}
		pajChan <- 1
		return err
	}()
	go func() (err error) {
		for _, installation := range p.Installations {
			err = p.CreateInstallation(installation)
			if err != nil {
				diChan <- 1
				return err
			}
		}
		diChan <- 1
		return err
	}()
	go func() (err error) {
		for _, content := range p.Content {
			err = p.CreateContentBridge(p.Categories, content)
			if err != nil {
				dcbChan <- 1
				return err
			}
		}
		dcbChan <- 1
		return err
	}()
	go func() (err error) {
		for _, price := range p.Pricing {
			price.PartId = p.ID
			err = price.Create()
			if err != nil {
				priceChan <- 1
				return err
			}
		}
		priceChan <- 1
		return err
	}()
	go func() (err error) {
		for _, review := range p.Reviews {
			review.PartID = p.ID
			err = review.Create()
			if err != nil {
				revChan <- 1
				return err
			}
		}
		revChan <- 1
		return err
	}()
	go func() (err error) {
		for _, image := range p.Images {
			image.PartID = p.ID
			err = image.Create()
			if err != nil {
				imageChan <- 1
				return err
			}
		}
		imageChan <- 1
		return err
	}()
	go func() (err error) {
		for _, related := range p.Related {
			err = p.CreateRelatedPart(related)
			if err != nil {
				relatedChan <- 1
				return err
			}
		}
		relatedChan <- 1
		return err
	}()
	go func() (err error) {
		for _, category := range p.Categories {
			err = p.CreatePartCategoryJoin(category)
			if err != nil {
				pcjChan <- 1
				return err
			}
		}
		pcjChan <- 1
		return err
	}()
	go func() (err error) {
		for _, video := range p.Videos {
			video.PartID = p.ID
			err = video.CreatePartVideo()
			if err != nil {
				videoChan <- 1
				return err
			}
		}
		videoChan <- 1
		return err
	}()
	go func() (err error) {
		for _, pack := range p.Packages {
			pack.PartID = p.ID
			err = pack.Create()
			if err != nil {
				packChan <- 1
				return err
			}
		}
		packChan <- 1
		return err
	}()

	<-pajChan
	<-diChan
	<-dcbChan
	<-priceChan
	<-revChan
	<-imageChan
	<-relatedChan
	<-pcjChan
	<-videoChan
	<-packChan

	return err
}

func (p *Part) Delete() (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	stmt, err := db.Prepare(deletePart)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(p.ID)
	if err != nil {
		return err
	}

	pajChan := make(chan int)
	diChan := make(chan int)
	dcbChan := make(chan int)
	priceChan := make(chan int)
	revChan := make(chan int)
	imageChan := make(chan int)
	relatedChan := make(chan int)
	pcjChan := make(chan int)
	videoChan := make(chan int)
	packChan := make(chan int)

	go func() (err error) {
		err = p.DeletePartAttributeJoins()
		if err != nil {
			pajChan <- 1
			return err
		}
		pajChan <- 1
		return err
	}()
	go func() (err error) {
		err = p.DeleteInstallations()
		if err != nil {
			diChan <- 1
			return err
		}
		diChan <- 1
		return err
	}()
	go func() (err error) {
		err = p.DeleteContentBridges()
		if err != nil {
			dcbChan <- 1
			return err
		}
		dcbChan <- 1
		return err
	}()
	go func() (err error) {
		var price Price
		price.PartId = p.ID
		err = price.DeleteByPart()
		if err != nil {
			priceChan <- 1
			return err
		}
		priceChan <- 1
		return err
	}()
	go func() (err error) {
		var review Review
		review.PartID = p.ID
		err = review.Delete()
		if err != nil {
			revChan <- 1
			return err
		}
		revChan <- 1
		return err
	}()
	go func() (err error) {
		var image Image
		image.PartID = p.ID
		err = image.DeleteByPart()
		if err != nil {
			imageChan <- 1
			return err
		}
		imageChan <- 1
		return err
	}()
	go func() (err error) {
		err = p.DeleteRelatedParts()
		if err != nil {
			relatedChan <- 1
			return err
		}
		relatedChan <- 1
		return err
	}()
	go func() (err error) {
		err = p.DeletePartCategoryJoins()
		if err != nil {
			pcjChan <- 1
			return err
		}
		pcjChan <- 1
		return err
	}()
	go func() (err error) {
		var v PartVideo
		v.PartID = p.ID
		err = v.DeleteByPart()
		if err != nil {
			videoChan <- 1
			return err
		}
		videoChan <- 1
		return err
	}()
	go func() (err error) {
		var pack Package
		pack.PartID = p.ID
		err = pack.DeleteByPart()
		if err != nil {
			packChan <- 1
			return err
		}
		packChan <- 1
		return err
	}()

	<-pajChan
	<-diChan
	<-dcbChan
	<-priceChan
	<-revChan
	<-imageChan
	<-relatedChan
	<-pcjChan
	<-videoChan
	<-packChan

	return nil
}

func (p *Part) Update() (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	stmt, err := db.Prepare(updatePart)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(p.Status, p.ShortDesc, p.PriceCode, p.Class.ID, p.Featured, p.AcesPartTypeID, p.ID)
	if err != nil {
		return err
	}
	//Refresh joins
	pajChan := make(chan int)
	diChan := make(chan int)
	dcbChan := make(chan int)
	priceChan := make(chan int)
	revChan := make(chan int)
	imageChan := make(chan int)
	relatedChan := make(chan int)
	pcjChan := make(chan int)
	videoChan := make(chan int)
	packChan := make(chan int)
	pajChanC := make(chan int)
	diChanC := make(chan int)
	dcbChanC := make(chan int)
	priceChanC := make(chan int)
	revChanC := make(chan int)
	imageChanC := make(chan int)
	relatedChanC := make(chan int)
	pcjChanC := make(chan int)
	videoChanC := make(chan int)
	packChanC := make(chan int)

	go func() (err error) {
		err = p.DeletePartAttributeJoins()
		if err != nil {
			pajChan <- 1
			return err
		}
		pajChan <- 1
		return err
	}()
	go func() (err error) {
		err = p.DeleteInstallations()
		if err != nil {
			diChan <- 1
			return err
		}
		diChan <- 1
		return err
	}()
	go func() (err error) {
		err = p.DeleteContentBridges()
		if err != nil {
			dcbChan <- 1
			return err
		}
		dcbChan <- 1
		return err
	}()
	go func() (err error) {
		var price Price
		price.PartId = p.ID
		err = price.DeleteByPart()
		if err != nil {
			priceChan <- 1
			return err
		}
		priceChan <- 1
		return err
	}()
	go func() (err error) {
		var review Review
		review.PartID = p.ID
		err = review.Delete()
		if err != nil {
			revChan <- 1
			return err
		}
		revChan <- 1
		return err
	}()
	go func() (err error) {
		var image Image
		image.PartID = p.ID
		err = image.DeleteByPart()
		if err != nil {
			imageChan <- 1
			return err
		}
		imageChan <- 1
		return err
	}()
	go func() (err error) {
		err = p.DeleteRelatedParts()
		if err != nil {
			relatedChan <- 1
			return err
		}
		relatedChan <- 1
		return err
	}()
	go func() (err error) {
		err = p.DeletePartCategoryJoins()
		if err != nil {
			pcjChan <- 1
			return err
		}
		pcjChan <- 1
		return err
	}()
	go func() (err error) {
		var v PartVideo
		v.PartID = p.ID
		err = v.DeleteByPart()
		if err != nil {
			videoChan <- 1
			return err
		}
		videoChan <- 1
		return err
	}()
	go func() (err error) {
		var pack Package
		pack.PartID = p.ID
		err = pack.DeleteByPart()
		if err != nil {
			packChan <- 1
			return err
		}
		packChan <- 1
		return err
	}()

	go func() (err error) {
		for _, attribute := range p.Attributes {
			err = p.CreatePartAttributeJoin(attribute)
			if err != nil {
				pajChanC <- 1
				return err
			}
		}
		pajChanC <- 1
		return err
	}()
	go func() (err error) {
		for _, installation := range p.Installations {
			err = p.CreateInstallation(installation)
			if err != nil {
				diChanC <- 1
				return err
			}
		}
		diChanC <- 1
		return err
	}()
	go func() (err error) {
		for _, content := range p.Content {
			err = p.CreateContentBridge(p.Categories, content)
			if err != nil {
				dcbChanC <- 1
				return err
			}
		}
		dcbChanC <- 1
		return err
	}()
	go func() (err error) {
		for _, price := range p.Pricing {
			price.PartId = p.ID
			err = price.Create()
			if err != nil {
				priceChanC <- 1
				return err
			}
		}
		priceChanC <- 1
		return err
	}()
	go func() (err error) {
		for _, review := range p.Reviews {
			review.PartID = p.ID
			err = review.Create()
			if err != nil {
				revChanC <- 1
				return err
			}
		}
		revChanC <- 1
		return err
	}()
	go func() (err error) {
		for _, image := range p.Images {
			image.PartID = p.ID
			err = image.Create()
			if err != nil {
				imageChanC <- 1
				return err
			}
		}
		imageChanC <- 1
		return err
	}()
	go func() (err error) {
		for _, related := range p.Related {
			err = p.CreateRelatedPart(related)
			if err != nil {
				relatedChanC <- 1
				return err
			}
		}
		relatedChanC <- 1
		return err
	}()
	go func() (err error) {
		for _, category := range p.Categories {
			err = p.CreatePartCategoryJoin(category)
			if err != nil {
				pcjChanC <- 1
				return err
			}
		}
		pcjChanC <- 1
		return err
	}()
	go func() (err error) {
		for _, video := range p.Videos {
			video.PartID = p.ID
			err = video.CreatePartVideo()
			if err != nil {
				videoChanC <- 1
				return err
			}
		}
		videoChanC <- 1
		return err
	}()
	go func() (err error) {
		for _, pack := range p.Packages {
			pack.PartID = p.ID
			err = pack.Create()
			if err != nil {
				packChanC <- 1
				return err
			}
		}
		packChanC <- 1
		return err
	}()

	<-pajChan
	<-diChan
	<-dcbChan
	<-priceChan
	<-revChan
	<-imageChan
	<-relatedChan
	<-pcjChan
	<-videoChan
	<-packChan
	<-pajChanC
	<-diChanC
	<-dcbChanC
	<-priceChanC
	<-revChanC
	<-imageChanC
	<-relatedChanC
	<-pcjChanC
	<-videoChanC
	<-packChanC

	return err
}

//Join Creators
func (p *Part) CreatePartAttributeJoin(a Attribute) (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	stmt, err := db.Prepare(createPartAttributeJoin)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(p.ID, a.Value, a.Key, a.Sort)
	if err != nil {
		return err
	}
	return nil
}

//Creates "VehiclePart" Join, which also contains installation fields
func (p *Part) CreateInstallation(i Installation) (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	stmt, err := db.Prepare(createVehiclePartJoin)
	if err != nil {
		return err
	}
	defer stmt.Close()
	res, err := stmt.Exec(i.Vehicle.ID, p.ID, i.Drilling, i.Exposed, i.InstallTime)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	i.ID = int(id)
	return nil
}

func (p *Part) CreateContentBridge(cats []Category, c Content) (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()
	stmt, err := tx.Prepare(createContentBridge)
	if err != nil {
		return err
	}
	defer stmt.Close()
	for _, cat := range cats {
		_, err = stmt.Exec(cat.ID, p.ID, c.ID)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	err = tx.Commit()
	return err
}

func (p *Part) CreateRelatedPart(relatedID int) (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	stmt, err := db.Prepare(createRelatedPart)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(p.ID, relatedID)
	if err != nil {
		return err
	}
	return nil
}

func (p *Part) CreatePartCategoryJoin(c Category) (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	stmt, err := db.Prepare(createPartCategoryJoin)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(c.ID, p.ID)
	if err != nil {
		return err
	}
	return nil
}

//delete Joins
func (p *Part) DeletePartAttributeJoins() (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	stmt, err := db.Prepare(deletePartAttributeJoins)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(p.ID)
	if err != nil {
		return err
	}
	return nil
}
func (p *Part) DeleteInstallations() (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	stmt, err := db.Prepare(deleteVehiclePartJoins)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(p.ID)
	if err != nil {
		return err
	}
	return nil
}
func (p *Part) DeleteContentBridges() (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	stmt, err := db.Prepare(deleteContentBridgeJoins)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(p.ID)
	if err != nil {
		return err
	}
	return nil
}
func (p *Part) DeleteRelatedParts() (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	stmt, err := db.Prepare(deleteRelatedParts)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(p.ID)
	if err != nil {
		return err
	}
	return nil
}
func (p *Part) DeletePartCategoryJoins() (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	stmt, err := db.Prepare(deletePartCategoryJoins)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(p.ID)
	if err != nil {
		return err
	}
	return nil
}
