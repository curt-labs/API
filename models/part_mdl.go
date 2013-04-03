package models

import (
	"../helpers/database"
	"../helpers/redis"
	"../helpers/rest"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var (
	basicsStmt = `select p.status, p.dateAdded, p.dateModified, p.shortDesc, p.partID, p.priceCode, pc.class
				from Part as p
				left join Class as pc on p.classID = pc.classID
				where p.partID = ? && p.status in (800,900) limit 1`

	basicsStmt_Grouped = `select p.status, p.dateAdded, p.dateModified, p.shortDesc, p.partID, p.priceCode, pc.class
				from Part as p
				left join Class as pc on p.classID = pc.classID
				where p.partID IN (%s) && p.status in (800,900)`

	partAttrStmt = `select field, value from PartAttribute where partID = ?`

	partAttrStmt_Grouped = `select partID,field, value from PartAttribute where partID IN (%s)`

	partPriceStmt = `select priceType, price, enforced from Price where partID = ?`

	partPriceStmt_Grouped = `select partID, priceType, price, enforced from Price where partID IN (%s)`

	relatedPartStmt = `select distinct relatedID from RelatedPart
				where partID = ?
				order by relatedID`

	relatedPartStmt_Grouped = `select distinct relatedID, partID from RelatedPart
				where partID IN (%s)
				order by relatedID`

	partContentStmt = `select ct.type, con.text
				from Content as con
				join ContentBridge as cb on con.contentID = cb.contentID
				join ContentType as ct on con.cTypeID = ct.cTypeID
				where cb.partID = ?
				order by ct.type`

	partContentStmt_Grouped = `select cb.partID, ct.type, con.text
				from Content as con
				join ContentBridge as cb on con.contentID = cb.contentID
				join ContentType as ct on con.cTypeID = ct.cTypeID
				where cb.partID IN (%s)
				order by ct.type`

	partInstallSheetStmt = `select c.text from ContentBridge as cb
					join Content as c on cb.contentID = c.contentID
					join ContentType as ct on c.cTypeID = ct.cTypeID
					where partID = ? && ct.type = 'InstallationSheet'
					limit 1`

	partInstallSheetStmt_Grouped = `select partID, c.text from ContentBridge as cb
					join Content as c on cb.contentID = c.contentID
					join ContentType as ct on c.cTypeID = ct.cTypeID
					where partID IN (?) && ct.type = 'InstallationSheet'`
)

type Part struct {
	PartId, Status, PriceCode, RelatedCount int
	AverageReview                           float64
	DateModified, DateAdded                 time.Time
	ShortDesc, PartClass                    string
	InstallSheet                            *url.URL
	Attributes                              []Attribute
	VehicleAttributes                       []string
	Content                                 []Content
	Pricing                                 []Pricing
	Reviews                                 []Review
	Images                                  []Image
	Related                                 []int
	Categories                              []ExtendedCategory
	Videos                                  []Video
	Packages                                []Package
	Customer                                CustomerPart
}

type PagedParts struct {
	Parts  []Part
	Paging []Paging
}

type Paging struct {
	CurrentIndex int
	PageCount    int
}

type CustomerPart struct {
	Price         float64
	CartReference int
}

type Attribute struct {
	Key, Value string
}

type Content struct {
	Key, Value string
}

type Pricing struct {
	Type     string
	Price    float64
	Enforced bool
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

	if len(errs) > 0 {
		return errors.New("Error: " + strings.Join(errs, ", "))
	}

	if part_bytes, err := json.Marshal(p); err == nil {
		part_key := "part:" + strconv.Itoa(p.PartId)
		redis.RedisClient.Set(part_key, part_bytes)
		redis.RedisClient.Expire(part_key, 86400)
	}

	return nil
}

func (p *Part) FromCache() error {

	part_bytes, err := redis.RedisClient.Get("part:" + strconv.Itoa(p.PartId))
	if err != nil {
		return err
	} else if len(part_bytes) == 0 {
		return errors.New("Part does not exist in cache")
	}

	err = json.Unmarshal(part_bytes, &p)

	return err
}

func (p *Part) Get(key string) error {

	partChan := make(chan int)
	customerChan := make(chan int)

	var err error

	go func() {
		if err = p.FromCache(); err != nil {
			err = p.FromDatabase()
		}
		partChan <- 1
	}()

	go func(api_key string) {
		err = p.BindCustomer(api_key)
		customerChan <- 1
	}(key)

	<-partChan
	<-customerChan

	return err
}

func (p *Part) GetWithVehicle(vehicle *Vehicle, api_key string) error {

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

func (p *Part) GetAttributes() (err error) {
	qry, err := database.Db.Prepare(partAttrStmt)
	if err != nil {
		return
	}

	rows, _, err := qry.Exec(p.PartId)
	if database.MysqlError(err) {
		return err
	} else if rows == nil {
		return
	}

	var attrs []Attribute
	for _, row := range rows {
		attr := Attribute{
			Key:   row.Str(0),
			Value: row.Str(1),
		}
		attrs = append(attrs, attr)
	}
	p.Attributes = attrs

	return
}

func (p *Part) Basics() error {
	qry, err := database.Db.Prepare(basicsStmt)
	if err != nil {
		return err
	}

	row, res, err := qry.ExecFirst(p.PartId)
	if database.MysqlError(err) {
		return err
	} else if row == nil {
		return errors.New("No Part Found for:" + string(p.PartId))
	}
	status := res.Map("status")
	dateAdded := res.Map("dateAdded")
	dateModified := res.Map("dateModified")
	shortDesc := res.Map("shortDesc")
	partID := res.Map("partID")
	priceCode := res.Map("priceCode")
	class := res.Map("class")

	p.PartId = row.Int(partID)
	date_add, _ := time.Parse("2006-01-02 15:04:15", row.Str(dateAdded))
	p.DateAdded = date_add

	date_mod, _ := time.Parse("2006-01-02 15:04:15", row.Str(dateModified))
	p.DateModified = date_mod

	p.ShortDesc = row.Str(shortDesc)
	p.PriceCode = row.Int(priceCode)
	p.PartClass = row.Str(class)
	p.Status = row.Int(status)

	return nil
}

func (p *Part) GetPricing() error {
	qry, err := database.Db.Prepare(partPriceStmt)
	if err != nil {
		return err
	}

	rows, res, err := qry.Exec(p.PartId)
	if database.MysqlError(err) {
		return err
	} else if rows == nil {
		return errors.New("No pricing found for part: " + string(p.PartId))
	}

	typ := res.Map("priceType")
	price := res.Map("price")
	enforced := res.Map("enforced")

	var prices []Pricing
	for _, row := range rows {
		pr := Pricing{
			row.Str(typ),
			row.Float(price),
			row.ForceBool(enforced),
		}

		if pr.Type == "Map" {
			pr.Enforced = true
		}
		prices = append(prices, pr)
	}

	p.Pricing = prices

	return nil
}

func (p *Part) GetRelated() error {
	qry, err := database.Db.Prepare(relatedPartStmt)
	if err != nil {
		return err
	}

	rows, _, err := qry.Exec(p.PartId)
	if database.MysqlError(err) {
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
	return nil
}

func (p *Part) GetContent() error {
	qry, err := database.Db.Prepare(partContentStmt)
	if err != nil {
		return err
	}

	rows, _, err := qry.Exec(p.PartId)
	if database.MysqlError(err) {
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
			sheetUrl, _ := url.Parse(con.Value)
			p.InstallSheet = sheetUrl
		} else {
			content = append(content, con)
		}
	}
	p.Content = content
	return nil
}

func (p *Part) BindCustomer(key string) error {
	price, err := GetCustomerPrice(key, p.PartId)
	if err != nil {
		return err
	}

	ref, err := GetCustomerCartReference(key, p.PartId)
	if err != nil {
		return err
	}

	cust := CustomerPart{
		Price:         price,
		CartReference: ref,
	}
	p.Customer = cust
	return nil
}

func (p *Part) GetInstallSheet(r *http.Request) (data []byte, err error) {
	qry, err := database.Db.Prepare(partInstallSheetStmt)
	if err != nil {
		return
	}

	row, _, err := qry.ExecFirst(p.PartId)
	if database.MysqlError(err) || row == nil {
		return
	}

	data, err = rest.GetPDF(row.Str(0), r)

	return
}

/*** Grouped Queries ****/

func GetByGroup(existing map[int]Part, key string) (parts map[int]Part, err error) {

	partChan := make(chan int)
	customerChan := make(chan int)

	genParts := make(map[int]Part, 0)
	custParts := make(map[int]CustomerPart, 0)
	go func() {
		//if err = p.FromCache(); err != nil {
		genParts, err = FromDatabaseByGroup(existing)
		//}
		partChan <- 1
	}()

	go func(api_key string) {
		custParts, err = BindCustomerByGroup(existing, api_key)
		customerChan <- 1
	}(key)

	<-partChan
	<-customerChan

	parts = genParts
	for k, _ := range parts {
		tmp := parts[k]
		tmp.Customer = custParts[k]
		parts[k] = tmp
	}

	return
}

func GetWithVehicleByGroup(existing map[int]Part, vehicle *Vehicle, api_key string) (parts map[int]Part, err error) {

	var errs []string
	parts = make(map[int]Part, len(existing))

	superChan := make(chan int)
	noteChan := make(chan int)

	// Get the Part data
	go func(key string) {
		getParts, getErr := GetByGroup(existing, key)
		if getErr != nil {
			errs = append(errs, getErr.Error())
		} else {
			for k, v := range getParts {
				parts[k] = v
			}
		}
		superChan <- 1
	}(api_key)

	// Get Notes for each part with this vehicle information
	go func() {
		// if nErr := vehicle.GetNotesByGroup(parts); nErr != nil {
		// 	errs = append(errs, nErr.Error())
		// }
		noteChan <- 1
	}()

	<-superChan
	<-noteChan

	if len(errs) > 0 {
		err = errors.New("Error: " + strings.Join(errs, ", "))
	}
	return
}

func FromDatabaseByGroup(existing map[int]Part) (parts map[int]Part, err error) {

	parts = make(map[int]Part, len(existing))
	var basicParts map[int]Part
	var attrParts map[int]Part
	var priceParts map[int]Part
	var relatedParts map[int]Part
	var contentParts map[int]Part

	for k, _ := range existing {
		parts[k] = Part{PartId: k}
	}

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

	go func(ps map[int]Part) {
		partArr, basicErr := BasicsByGroup(ps)
		if basicErr != nil {
			errs = append(errs, basicErr.Error())
		} else {
			basicParts = partArr
		}
		basicChan <- 1
	}(existing)

	go func(ps map[int]Part) {
		partArr, attrErr := GetAttributesByGroup(ps)
		if attrErr != nil {
			errs = append(errs, attrErr.Error())
		} else {
			attrParts = partArr
		}
		attrChan <- 1
	}(existing)

	go func(ps map[int]Part) {
		partArr, priceErr := GetPricingByGroup(ps)
		if priceErr != nil {
			errs = append(errs, priceErr.Error())
		} else {
			priceParts = partArr
		}
		priceChan <- 1
	}(existing)

	go func(ps map[int]Part) {
		// reviewErr := GetReviewsByGroup(ps)
		// if reviewErr != nil {
		// 	errs = append(errs, reviewErr.Error())
		// }
		reviewChan <- 1
	}(existing)

	go func(ps map[int]Part) {
		// imgErr := GetImagesByGroup(ps)
		// if imgErr != nil {
		// 	errs = append(errs, imgErr.Error())
		// }
		imageChan <- 1
	}(existing)

	go func(ps map[int]Part) {
		// vidErr := GetVideosByGroup(ps)
		// if vidErr != nil {
		// 	errs = append(errs, vidErr.Error())
		// }
		videoChan <- 1
	}(existing)

	go func(ps map[int]Part) {
		partArr, relErr := GetRelatedByGroup(ps)
		if relErr != nil {
			errs = append(errs, relErr.Error())
		} else {
			relatedParts = partArr
		}
		relatedChan <- 1
	}(existing)

	go func(ps map[int]Part) {
		// pkgErr := GetPartPackagingByGroup(ps)
		// if pkgErr != nil {
		// 	errs = append(errs, pkgErr.Error())
		// }
		packageChan <- 1
	}(existing)

	go func(ps map[int]Part) {
		// catErr := PartBreadcrumbsByGroup(ps)
		// if catErr != nil {
		// 	errs = append(errs, catErr.Error())
		// }
		categoryChan <- 1
	}(existing)

	go func(ps map[int]Part) {
		partArr, conErr := GetContentByGroup(ps)
		if conErr != nil {
			errs = append(errs, conErr.Error())
		} else {
			contentParts = partArr
		}
		contentChan <- 1
	}(existing)

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

	if len(errs) > 0 {
		err = errors.New("Error: " + strings.Join(errs, ", "))
		return
	}

	for k, v := range basicParts {
		tmp := parts[k]
		tmp.PartId = v.PartId
		tmp.DateAdded = v.DateAdded
		tmp.DateModified = v.DateModified
		tmp.ShortDesc = v.ShortDesc
		tmp.PriceCode = v.PriceCode
		tmp.PartClass = v.PartClass
		tmp.Status = v.Status
		parts[k] = tmp
	}

	for k, v := range attrParts {
		tmp := parts[k]
		tmp.Attributes = v.Attributes
		parts[k] = tmp
	}

	for k, v := range priceParts {
		tmp := parts[k]
		tmp.Pricing = v.Pricing
		parts[k] = tmp
	}

	for k, v := range relatedParts {
		tmp := parts[k]
		tmp.Related = v.Related
		tmp.RelatedCount = len(v.Related)
		parts[k] = tmp
	}

	for k, v := range contentParts {
		tmp := parts[k]
		tmp.Content = v.Content
		tmp.InstallSheet = v.InstallSheet
		parts[k] = tmp
	}

	// if part_bytes, err := json.Marshal(p); err == nil {
	// 	part_key := "part:" + strconv.Itoa(p.PartId)
	// 	redis.RedisClient.Set(part_key, part_bytes)
	// 	redis.RedisClient.Expire(part_key, 86400)
	// }

	return
}

// func FromCacheByGroup(ids []int) (parts []Part, err error) {

// 	part_bytes, err := redis.RedisClient.Get("part:" + strconv.Itoa(p.PartId))
// 	if err != nil {
// 		return err
// 	} else if len(part_bytes) == 0 {
// 		return errors.New("Part does not exist in cache")
// 	}

// 	err = json.Unmarshal(part_bytes, &p)

// 	return err
// }

func BasicsByGroup(existing map[int]Part) (parts map[int]Part, err error) {

	parts = make(map[int]Part, len(existing))

	var ids []string
	for k, _ := range existing {
		parts[k] = Part{PartId: k}
		ids = append(ids, strconv.Itoa(k))
	}

	rows, res, err := database.Db.Query(basicsStmt_Grouped, strings.Join(ids, ","))
	if database.MysqlError(err) {
		return
	} else if len(rows) == 0 {
		err = errors.New("No Parts Found for:" + strings.Join(ids, ","))
		return
	}

	status := res.Map("status")
	dateAdded := res.Map("dateAdded")
	dateModified := res.Map("dateModified")
	shortDesc := res.Map("shortDesc")
	partID := res.Map("partID")
	priceCode := res.Map("priceCode")
	class := res.Map("class")

	for _, row := range rows {
		date_add, _ := time.Parse("2006-01-02 15:04:15", row.Str(dateAdded))
		date_mod, _ := time.Parse("2006-01-02 15:04:15", row.Str(dateModified))

		pId := row.Int(partID)
		part := Part{
			PartId:       pId,
			DateAdded:    date_add,
			DateModified: date_mod,
			ShortDesc:    row.Str(shortDesc),
			PriceCode:    row.Int(priceCode),
			PartClass:    row.Str(class),
			Status:       row.Int(status),
		}
		parts[pId] = part
	}

	return
}

func GetAttributesByGroup(existing map[int]Part) (parts map[int]Part, err error) {

	parts = make(map[int]Part, len(existing))
	var ids []string
	for k, _ := range existing {
		parts[k] = Part{PartId: k}
		ids = append(ids, strconv.Itoa(k))
	}

	rows, _, err := database.Db.Query(partAttrStmt_Grouped, strings.Join(ids, ","))
	if database.MysqlError(err) || len(rows) == 0 {
		return
	}

	for _, row := range rows {
		pId := row.Int(0)
		attr := Attribute{
			Key:   row.Str(1),
			Value: row.Str(2),
		}
		tmp := parts[pId]
		tmp.Attributes = append(tmp.Attributes, attr)

		parts[pId] = tmp
	}

	return
}

func GetPricingByGroup(existing map[int]Part) (parts map[int]Part, err error) {

	parts = make(map[int]Part, len(existing))
	var ids []string
	for k, _ := range existing {
		parts[k] = Part{PartId: k}
		ids = append(ids, strconv.Itoa(k))
	}

	rows, res, err := database.Db.Query(partPriceStmt_Grouped, strings.Join(ids, ","))
	if database.MysqlError(err) {
		return
	} else if len(rows) == 0 {
		err = errors.New("No pricing found")
		return
	}

	partID := res.Map("partID")
	typ := res.Map("priceType")
	price := res.Map("price")
	enforced := res.Map("enforced")

	for _, row := range rows {
		pId := row.Int(partID)

		pr := Pricing{
			row.Str(typ),
			row.Float(price),
			row.ForceBool(enforced),
		}

		if pr.Type == "Map" {
			pr.Enforced = true
		}

		tmp := parts[pId]
		tmp.Pricing = append(tmp.Pricing, pr)
		parts[pId] = tmp
	}

	return
}

func GetRelatedByGroup(existing map[int]Part) (parts map[int]Part, err error) {

	parts = make(map[int]Part, len(existing))
	var ids []string
	for k, _ := range existing {
		parts[k] = Part{PartId: k}
		ids = append(ids, strconv.Itoa(k))
	}

	rows, res, err := database.Db.Query(relatedPartStmt_Grouped, strings.Join(ids, ","))
	if database.MysqlError(err) {
		return
	} else if len(rows) == 0 {
		err = errors.New("No related found")
		return
	}

	relatedID := res.Map("relatedID")
	partID := res.Map("partID")

	for _, row := range rows {
		pId := row.Int(partID)

		tmp := parts[pId]
		tmp.Related = append(tmp.Related, row.Int(relatedID))
		parts[pId] = tmp
	}

	return
}

func GetContentByGroup(existing map[int]Part) (parts map[int]Part, err error) {

	parts = make(map[int]Part, len(existing))
	var ids []string
	for k, _ := range existing {
		parts[k] = Part{PartId: k}
		ids = append(ids, strconv.Itoa(k))
	}

	rows, _, err := database.Db.Query(partContentStmt_Grouped, strings.Join(ids, ","))
	if database.MysqlError(err) {
		return
	} else if len(rows) == 0 {
		err = errors.New("No content found")
		return
	}

	for _, row := range rows {
		pId := row.Int(0)
		con := Content{
			Key:   row.Str(1),
			Value: row.Str(2),
		}

		tmp := parts[pId]
		if strings.Contains(strings.ToLower(con.Key), "install") {
			sheetUrl, _ := url.Parse(con.Value)
			tmp.InstallSheet = sheetUrl
		} else {
			tmp.Content = append(tmp.Content, con)
		}
		parts[pId] = tmp
	}

	return
}

func BindCustomerByGroup(existing map[int]Part, key string) (custs map[int]CustomerPart, err error) {

	custs = make(map[int]CustomerPart, len(existing))
	var ids []string
	for k, _ := range existing {
		custs[k] = CustomerPart{}
		ids = append(ids, strconv.Itoa(k))
	}

	prices, err := GetCustomerPriceByGroup(key, existing)
	if err != nil {
		return
	}

	refs, err := GetCustomerCartReferenceByGroup(key, existing)
	if err != nil {
		return
	}

	for k, v := range custs {
		if _, ok := prices[k]; ok {
			v.Price = prices[k]
		}

		if _, ok := refs[k]; ok {
			v.CartReference = refs[k]
		}

		custs[k] = v
	}
	return
}
