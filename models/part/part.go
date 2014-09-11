package part

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
	"github.com/curt-labs/GoAPI/models/category"
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
	SubCategoryPartIdStmt = `select distinct cp.partID
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
	PartId, Status, PriceCode, RelatedCount int
	AverageReview                           float64
	DateModified, DateAdded                 time.Time
	ShortDesc, PartClass                    string
	InstallSheet                            *url.URL
	Attributes                              []Attribute
	VehicleAttributes                       []string
	Vehicles                                []vehicle.Vehicle `json:",omitempty" xml:",omitempty"`
	Content                                 []Content
	Pricing                                 []Price
	Reviews                                 []Review
	Images                                  []Image
	Related                                 []int
	Categories                              []category.ExtendedCategory
	Videos                                  []PartVideo
	Packages                                []Package
	Customer                                CustomerPart
}

type PagedParts struct {
	Parts  []Part
	Paging []Paging
}

type CategoryTree struct {
	CategoryId    int
	SubCategories []int
	Parts         []Part
}

type Paging struct {
	CurrentIndex int
	PageCount    int
}

type CustomerPart struct {
	Price         float64
	CartReference int
}

type Content struct {
	Key, Value string
}

func (p *Part) FromDatabase() error {

	var errs []string

	basicChan := make(chan int)
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
		basicErr := p.Basics()
		if basicErr != nil {
			errs = append(errs, basicErr.Error())
		}
		basicChan <- 1
	}()

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
		reviewErr := p.GetReviews()
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

	<-basicChan
	<-attrChan
	<-priceChan
	<-reviewChan
	<-imageChan
	<-videoChan
	<-relatedChan
	<-packageChan
	<-categoryChan
	<-contentChan

	go redis.Setex("part:"+strconv.Itoa(p.PartId), p, redis.CacheTimeout)

	return nil
}

func (p *Part) FromCache() error {

	part_bytes, err := redis.Get("part:" + strconv.Itoa(p.PartId))
	if err != nil {
		return err
	} else if len(part_bytes) == 0 {
		return errors.New("Part does not exist in cache")
	}

	return json.Unmarshal(part_bytes, &p)
}

func (p *Part) Get(key string) error {

	// partChan := make(chan int)
	customerChan := make(chan int)

	var err error

	// go func() {
	// 	if err = p.FromCache(); err != nil {
	// 		err = p.FromDatabase()
	// 	}
	// 	partChan <- 1
	// }()

	go func(api_key string) {
		err = p.BindCustomer(api_key)
		customerChan <- 1
	}(key)

	if err := p.FromDatabase(); err != nil {
		return err
	}

	// <-partChan
	<-customerChan

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
			p := Part{PartId: id}
			p.Get(key)
			parts = append(parts, p)
			partChan <- 1
		}(partID)
		iter++
	}

	for i := 0; i < iter; i++ {
		<-partChan
	}

	sortutil.AscByField(parts, "PartId")

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
		notes, nErr := vehicle.GetNotes(p.PartId)
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
	p.PartId = id

	p.Get(key)
}

func (p *Part) Basics() error {
	redis_key := fmt.Sprintf("part:%d:basics", p.PartId)

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

	row := qry.QueryRow(p.PartId)
	if row == nil {
		return errors.New("No Part Found for:" + string(p.PartId))
	}

	row.Scan(
		&p.Status,
		&p.DateAdded,
		&p.DateModified,
		&p.ShortDesc,
		&p.PartId,
		&p.PriceCode,
		&p.PartClass)

	p.ShortDesc = fmt.Sprintf("CURT %s %d", p.ShortDesc, p.PartId)

	go redis.Setex(redis_key, p, redis.CacheTimeout)

	return nil
}

func (p *Part) GetRelated() error {
	redis_key := fmt.Sprintf("part:%d:related", p.PartId)

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

	rows, _, err := qry.Exec(p.PartId)
	if err != nil {
		return err
	} else if rows == nil {
		return errors.New("No related found for part: " + string(p.PartId))
	}

	var related []int
	for _, row := range rows {
		related = append(related, row.Int(0))
	}
	p.Related = related
	p.RelatedCount = len(related)

	go redis.Setex(redis_key, p.Related, redis.CacheTimeout)

	return nil
}

func (p *Part) GetContent() error {
	redis_key := fmt.Sprintf("part:%d:content", p.PartId)

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

	rows, _, err := qry.Exec(p.PartId)
	if err != nil {
		return err
	} else if rows == nil {
		return errors.New("No content found for part: " + string(p.PartId))
	}

	var content []Content
	for _, row := range rows {
		con := Content{
			Key:   row.Str(0),
			Value: row.Str(1),
		}

		if strings.Contains(strings.ToLower(con.Key), "install") {
			//sheetUrl, _ := url.Parse(con.Value)
			p.InstallSheet, _ = url.Parse(api_helpers.API_DOMAIN + "/part/" + strconv.Itoa(p.PartId) + ".pdf")
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
		price, _ = customer.GetCustomerPrice(key, p.PartId)
		priceChan <- 1
	}()

	go func() {
		ref, _ = customer.GetCustomerCartReference(key, p.PartId)
		refChan <- 1
	}()

	go func() {
		content, _ := custcontent.GetPartContent(p.PartId, key)
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
	redis_key := fmt.Sprintf("part:%d:installsheet", p.PartId)

	data, err = redis.Get(redis_key)
	if err == nil && len(data) > 0 {
		return data, nil
	}

	qry, err := database.Db.Prepare(partInstallSheetStmt)
	if err != nil {
		return
	}

	row, _, err := qry.ExecFirst(p.PartId)
	if err != nil || row == nil {
		return
	}

	data, err = rest.GetPDF(row.Str(0), r)

	go redis.Setex(redis_key, data, redis.CacheTimeout)

	return
}

// PartBreacrumbs
//
// Description: Builds out Category breadcrumb array for the current part object.
//
// Inherited: part Part
// Returns: error
func (p *Part) PartBreadcrumbs() error {

	redis_key := fmt.Sprintf("part:%d:breadcrumbs", p.PartId)

	data, err := redis.Get(redis_key)
	if err == nil && len(data) > 0 {
		if err = json.Unmarshal(data, &p.Categories); err == nil {
			return nil
		}
	}

	if p.PartId == 0 {
		return errors.New("Invalid Part Number")
	}

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	qry, err := db.Prepare(category.PartCategoryStmt)
	if err != nil {
		return err
	}
	defer qry.Close()

	parentQuery, err := db.Prepare(category.ParentCategoryStmt)
	if err != nil {
		return err
	}
	defer parentQuery.Close()

	// Execute SQL Query against current PartId
	catRow := qry.QueryRow(p.PartId)
	if catRow == nil {
		return errors.New("No part found for " + string(p.PartId))
	}

	ch := make(chan category.ExtendedCategory)
	go category.PopulateExtendedCategory(catRow, ch)
	initCat := <-ch

	// Instantiate our array with the initial category
	var cats []category.ExtendedCategory
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
			catRow = parentQuery.QueryRow(parent)
			if catRow == nil {
				break
			}

			ch := make(chan category.ExtendedCategory)
			go category.PopulateExtendedCategory(catRow, ch)

			// Append new Category onto array
			subCat := <-ch
			cats = append(cats, subCat)
			parent = subCat.ParentId
		}
	}

	// Apply breadcrumbs to our part object and return
	p.Categories = cats

	go redis.Setex(redis_key, p.Categories, redis.CacheTimeout)

	return nil
}

func (p *Part) GetPartCategories(key string) (cats []category.ExtendedCategory, err error) {

	redis_key := fmt.Sprintf("part:%d:categories", p.PartId)

	data, err := redis.Get(redis_key)
	if err == nil && len(data) > 0 {
		if err := json.Unmarshal(data, &cats); err == nil {
			return cats, nil
		}
	}

	if p.PartId == 0 {
		return
	}

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return
	}
	defer db.Close()

	qry, err := db.Prepare(category.PartAllCategoryStmt)
	if err != nil {
		return
	}
	defer qry.Close()

	// Execute SQL Query against current PartId
	catRows, err := qry.Query(p.PartId)
	if err != nil || catRows == nil { // Error occurred while executing query
		return
	}

	ch := make(chan []category.ExtendedCategory, 0)
	go category.PopulateExtendedCategoryMulti(catRows, ch)
	cats = <-ch

	for _, cat := range cats {

		contentChan := make(chan int)
		subChan := make(chan int)
		customerChan := make(chan int)

		c := category.Category{
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

		go func() {
			content, _ := custcontent.GetCategoryContent(cat.CategoryId, key)
			for _, con := range content {
				strArr := strings.Split(con.ContentType.Type, ":")
				cType := con.ContentType.Type
				if len(strArr) > 1 {
					cType = strArr[1]
				}
				cat.Content = append(cat.Content, category.Content{
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

	go redis.Setex(redis_key, cats, redis.CacheTimeout)

	return
}

func GetCategoryParts(c category.Category, key string, page int, count int) (parts []Part, err error) {

	if c.CategoryId == 0 {
		return
	}

	if page > 0 {
		page = count * page
	}

	redis_key := "category:" + strconv.Itoa(c.CategoryId) + ":tree:" + strconv.Itoa(page) + ":" + strconv.Itoa(count)
	data, err := redis.Get(redis_key)
	if err == nil && len(data) > 0 {
		err = json.Unmarshal(data, &parts)
		if err != nil {
			return
		}

		chans := make(chan int, len(parts))
		for _, part := range parts {
			go func(p Part, k string) {
				p.BindCustomer(k)
				chans <- 1
			}(part, key)
		}

		for i := 0; i < len(parts); i++ {
			<-chans
		}

		return
	}
	log.Println("missed redis")

	tree := CategoryTree{
		CategoryId: c.CategoryId,
	}

	tree.CategoryTreeBuilder()
	catIdStr := strconv.Itoa(tree.CategoryId)
	for _, treeId := range tree.SubCategories {
		catIdStr = catIdStr + "," + strconv.Itoa(treeId)
	}

	rows, _, err := database.Db.Query(SubCategoryPartIdStmt, catIdStr, page, count)
	if err != nil {
		return nil, err
	}

	for _, row := range rows {
		tree.Parts = append(tree.Parts, Part{PartId: row.Int(0)})
	}

	// This will work for populating the
	// parts that match this exact category.
	chans := make(chan int, len(tree.Parts))

	for _, part := range tree.Parts {
		go func(p Part) {
			p.Get(key)
			parts = append(parts, p)
			chans <- 1
		}(part)

	}

	for i := 0; i < len(tree.Parts); i++ {
		<-chans
	}

	go redis.Setex(redis_key, parts, redis.CacheTimeout)

	return
}

func (tree *CategoryTree) CategoryTreeBuilder() {

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return
	}
	defer db.Close()

	subQry, err := db.Prepare(category.SubCategoryIdStmt)
	if err != nil {
		return
	}
	defer subQry.Close()

	// Execute against current Category Id
	// to retrieve all category Ids that are children.
	rows, err := subQry.Query(tree.CategoryId)
	if err != nil {
		return
	}

	chans := make(chan int, 0)
	var rowCount int
	for rows.Next() {
		var catID int
		if err := rows.Scan(&catID); err != nil {
			continue
		}

		go func(catID int) {

			// Need to parse out string array into ints and populate
			cat := category.Category{
				CategoryId: catID,
			}
			tree.SubCategories = append(tree.SubCategories, cat.CategoryId)

			subRows, err := subQry.Query(cat.CategoryId)
			if err == nil && subRows != nil {
				for subRows.Next() {
					var subID int
					if err := subRows.Scan(&subID); err == nil {
						tree.SubCategories = append(tree.SubCategories, subID)
					}

					// subTree.CategoryTreeBuilder()
					// tree.SubCategories = append(tree.SubCategories, subTree.SubCategories...)
				}
			}
			chans <- 1
		}(catID)
		rowCount++
	}

	for i := 0; i < rowCount; i++ {
		<-chans
	}

	return
}
