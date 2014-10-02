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
	"log"
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
	PartClass         string            `json:"part_class" xml:"part_class,attr"`
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
}

type CustomerPart struct {
	Price         float64 `json:"price" xml:"price,attr"`
	CartReference int     `json:"cart_reference" xml:"cart_reference,attr"`
}

type Content struct {
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
			log.Printf("Attribute Error: %s", attrErr.Error())
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

	p.ShortDesc = fmt.Sprintf("CURT %s %d", p.ShortDesc, p.ID)

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

	qry, err := database.Db.Prepare(relatedPartStmt)
	if err != nil {
		return err
	}

	rows, _, err := qry.Exec(p.ID)
	if err != nil {
		return err
	} else if rows == nil {
		return errors.New("No related found for part: " + string(p.ID))
	}

	var related []int
	for _, row := range rows {
		related = append(related, row.Int(0))
	}
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

	qry, err := database.Db.Prepare(partContentStmt)
	if err != nil {
		return err
	}

	rows, _, err := qry.Exec(p.ID)
	if err != nil {
		return err
	} else if rows == nil {
		return errors.New("No content found for part: " + string(p.ID))
	}

	var content []Content
	for _, row := range rows {
		con := Content{
			Key:   row.Str(0),
			Value: row.Str(1),
		}

		if strings.Contains(strings.ToLower(con.Key), "install") {
			//sheetUrl, _ := url.Parse(con.Value)
			p.InstallSheet, _ = url.Parse(api_helpers.API_DOMAIN + "/part/" + strconv.Itoa(p.ID) + ".pdf")
		} else {
			content = append(content, con)
		}
	}
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

	qry, err := database.Db.Prepare(partInstallSheetStmt)
	if err != nil {
		return
	}

	row, _, err := qry.ExecFirst(p.ID)
	if err != nil || row == nil {
		return
	}

	data, err = rest.GetPDF(row.Str(0), r)

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

// 	db, err := sql.Open("mysql", database.ConnectionString())
// 	if err != nil {
// 		return
// 	}
// 	defer db.Close()

// 	subQry, err := db.Prepare(SubIDStmt)
// 	if err != nil {
// 		return
// 	}
// 	defer subQry.Close()

// 	// Execute against current Category Id
// 	// to retrieve all category Ids that are children.
// 	rows, err := subQry.Query(tree.ID)
// 	if err != nil {
// 		return
// 	}

// 	chans := make(chan int, 0)
// 	var rowCount int
// 	for rows.Next() {
// 		var catID int
// 		if err := rows.Scan(&catID); err != nil {
// 			continue
// 		}

// 		go func(catID int) {

// 			// Need to parse out string array into ints and populate
// 			cat := Category{
// 				ID: catID,
// 			}
// 			tree.SubCategories = append(tree.SubCategories, cat.ID)

// 			subRows, err := subQry.Query(cat.ID)
// 			if err == nil && subRows != nil {
// 				for subRows.Next() {
// 					var subID int
// 					if err := subRows.Scan(&subID); err == nil {
// 						tree.SubCategories = append(tree.SubCategories, subID)
// 					}

// 					// subTree.CategoryTreeBuilder()
// 					// tree.SubCategories = append(tree.SubCategories, subTree.SubCategories...)
// 				}
// 			}
// 			chans <- 1
// 		}(catID)
// 		rowCount++
// 	}

// 	for i := 0; i < rowCount; i++ {
// 		<-chans
// 	}

// 	return
// }
