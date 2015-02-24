package products

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/curt-labs/GoAPI/helpers/apicontext"
	"github.com/curt-labs/GoAPI/helpers/database"
	"github.com/curt-labs/GoAPI/helpers/sortutil"
	"github.com/curt-labs/GoAPI/models/contact"
	_ "github.com/go-sql-driver/mysql"
	"math"
	"strconv"
	"strings"
)

var (
	getVcdbVehicleIDWithSubmodel = `
		select v.VehicleID from Vehicle as v
		join BaseVehicle as bv on v.BaseVehicleID = bv.BaseVehicleID
		join Make as ma on bv.MakeID = ma.MakeID
		join Model as mo on bv.ModelID = mo.ModelID
		join Submodel as s on v.SubmodelID = s.SubmodelID
		where bv.YearID = ? && ma.MakeName = ? && mo.ModelName = ? && s.SubmodelName = ?
		limit 1`
	getVcdbVehicleID = `
		select v.VehicleID from Vehicle as v
		join BaseVehicle as bv on v.BaseVehicleID = bv.BaseVehicleID
		join Make as ma on bv.MakeID = ma.MakeID
		join Model as mo on bv.ModelID = mo.ModelID
		where bv.YearID = ? && ma.MakeName = ? && mo.ModelName = ? && (v.SubmodelID = 0 || v.SubmodelID is null)
		limit 1`
	getVehicleParts = `
		select distinct vp.PartNumber
		from vcdb_Vehicle as v
		join Submodel as s on v.SubModelID = s.ID
		join BaseVehicle as bv on v.BaseVehicleID = bv.ID
		join vcdb_Model as mo on bv.ModelID = mo.ID
		join vcdb_Make as ma on bv.MakeID = ma.ID
		join vcdb_VehiclePart as vp on v.ID = vp.VehicleID
		join Part as p on vp.PartNumber = p.partID
		where (p.status = 800 || p.status = 900) && 
		bv.YearID = ? && ma.MakeName = ? &&
		mo.ModelName = ? && s.SubmodelName = ? &&
		(v.ConfigID = 0 || v.ConfigID is null)
		&& p.brandID in(?)
		order by vp.PartNumber`
	getBaseVehicleParts = `
		select distinct vp.PartNumber
		from vcdb_Vehicle as v
		join BaseVehicle as bv on v.BaseVehicleID = bv.ID
		join vcdb_Model as mo on bv.ModelID = mo.ID
		join vcdb_Make as ma on bv.MakeID = ma.ID
		join vcdb_VehiclePart as vp on v.ID = vp.VehicleID
		join Part as p on vp.PartNumber = p.partID
		where (p.status = 800 || p.status = 900) &&
		bv.YearID = ? && ma.MakeName = ? &&
		mo.ModelName = ? && (v.SubmodelID = 0 || v.SubmodelID is null) &&
		(v.ConfigID = 0 || v.ConfigID is null)
		&& p.brandID in (?)
		order by vp.PartNumber`
)

type Vehicle struct {
	Base           BaseVehicle     `json:"base" xml:"base"`
	Submodel       string          `json:"submodel" xml:"submodel"`
	Configurations []Configuration `json:"configurations" xml:"configurations"`
}

type Lookup struct {
	Years          []int                 `json:"available_years,omitempty" xml:"available_years,omitempty"`
	Makes          []string              `json:"available_makes,omitempty" xml:"available_makes,omitempty"`
	Models         []string              `json:"available_models,omitempty" xml:"available_models,omitempty"`
	Submodels      []string              `json:"available_submodels,omitempty" xml:"available_submodels,omitempty"`
	Configurations []ConfigurationOption `json:"available_configurations,omitempty" xml:"available_configurations,omitempty"`
	Vehicle        Vehicle               `json:"vehicle" xml:"vehicle"`
	Parts          []Part                `json:"parts" xml:"parts"`
	Filter         interface{}           `json:"filter" xml:"filter"`
	Pagination     Pagination            `json:"pagination" xml:"pagination"`
	CustomerKey    string                `json:"-" xml:"-"`
	Brands         []int                 `json:"-" xml:"-"`
}

type Pagination struct {
	TotalItems    int `json:"total_items" xml:"total_items"`
	ReturnedCount int `json:"returned_count" xml:"returned_count"`
	Page          int `json:"page" xml:"page"`
	PerPage       int `json:"per_page" xml:"per_page"`
	TotalPages    int `json:"total_pages" xml:"total_pages"`
}

func (l *Lookup) LoadParts(ch chan []Part, page int, count int, dtx *apicontext.DataContext) {
	if count == 0 {
		count = DefaultPageCount
	}

	parts := make([]Part, 0)

	vehicleChan := make(chan error)
	baseVehicleChan := make(chan error)
	go l.loadVehicleParts(vehicleChan)
	go l.loadBaseVehicleParts(baseVehicleChan)

	if len(l.Vehicle.Configurations) > 0 {
		configs, err := l.Vehicle.getDefinedConfigurations(l.CustomerKey)
		if err != nil || configs == nil {
			ch <- parts
			return
		}

		chosenValArr := make(map[string]string, 0)
		for _, config := range l.Vehicle.Configurations {
			chosenValArr[strings.TrimSpace(strings.ToLower(config.Value))] = strings.TrimSpace(strings.ToLower(config.Value))
		}

		for _, config := range *configs {
			// configValArr := make(map[string]string, 0)
			matches := true
			for _, val := range config {

				v := strings.TrimSpace(strings.ToLower(val.Value))

				if _, ok := chosenValArr[v]; !ok {
					matches = false
				}
			}
			if matches {
				for _, partID := range config[0].Parts {
					p := Part{ID: partID}
					l.Parts = append(l.Parts, p)
				}
			}
		}
	}

	<-vehicleChan
	<-baseVehicleChan
	removeDuplicates(&l.Parts)

	// we need to strip the result set down the paginated
	// version
	partCount := len(l.Parts)
	pagedParts := l.Parts
	if partCount > count {
		start := 0
		if page > 1 {
			start = count * (page - 1)
		}
		end := start + count
		if len(l.Parts) <= end {
			pagedParts = l.Parts[start:]
		} else {
			pagedParts = l.Parts[start : start+count]
		}
	}

	parts = make([]Part, 0)
	perChan := make(chan int)
	for i, p := range pagedParts {
		go func(j int, prt Part) {
			if err := prt.Get(dtx); err == nil && prt.ShortDesc != "" {
				parts = append(parts, prt)
			}
			perChan <- 1
		}(i, p)
	}

	for _, _ = range pagedParts {
		<-perChan
	}
	l.Parts = parts

	sortutil.AscByField(l.Parts, "ID")

	mod := math.Mod(float64(partCount), float64(count))
	totalPages := partCount / count
	if mod > 0 {
		totalPages++
	}
	if page == 0 {
		page = 1
	}

	l.Pagination = Pagination{
		TotalItems:    partCount,
		ReturnedCount: len(l.Parts),
		Page:          page,
		PerPage:       count,
		TotalPages:    totalPages,
	}

	ch <- parts
}

func (l *Lookup) loadVehicleParts(ch chan error) {
	stmtBeginning := `select distinct vp.PartNumber
		from vcdb_Vehicle as v
		join Submodel as s on v.SubModelID = s.ID
		join BaseVehicle as bv on v.BaseVehicleID = bv.ID
		join vcdb_Model as mo on bv.ModelID = mo.ID
		join vcdb_Make as ma on bv.MakeID = ma.ID
		join vcdb_VehiclePart as vp on v.ID = vp.VehicleID
		join Part as p on vp.PartNumber = p.partID
		where (p.status = 800 || p.status = 900) && 
		bv.YearID = ? && ma.MakeName = ? &&
		mo.ModelName = ? && s.SubmodelName = ? &&
		(v.ConfigID = 0 || v.ConfigID is null)`
	stmtEnd := `order by vp.PartNumber`
	brandStmt := " && p.brandID in ("
	for _, b := range l.Brands {
		brandStmt += strconv.Itoa(b) + ","
	}
	brandStmt = strings.TrimRight(brandStmt, ",") + ")"
	wholeStmt := stmtBeginning + brandStmt + stmtEnd

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		ch <- err
		return
	}
	defer db.Close()

	stmt, err := db.Prepare(wholeStmt)
	if err != nil {
		ch <- err
		return
	}

	rows, err := stmt.Query(l.Vehicle.Base.Year, l.Vehicle.Base.Make, l.Vehicle.Base.Model, l.Vehicle.Submodel)
	if err != nil || rows == nil {
		ch <- err
		return
	}

	for rows.Next() {
		var p Part
		if err = rows.Scan(&p.ID); err == nil {
			l.Parts = append(l.Parts, p)
		}
	}
	defer rows.Close()
	ch <- nil
	return
}

func (l *Lookup) loadBaseVehicleParts(ch chan error) {
	stmtBeginning := `select distinct vp.PartNumber
		from vcdb_Vehicle as v
		join BaseVehicle as bv on v.BaseVehicleID = bv.ID
		join vcdb_Model as mo on bv.ModelID = mo.ID
		join vcdb_Make as ma on bv.MakeID = ma.ID
		join vcdb_VehiclePart as vp on v.ID = vp.VehicleID
		join Part as p on vp.PartNumber = p.partID
		where (p.status = 800 || p.status = 900) &&
		bv.YearID = ? && ma.MakeName = ? &&
		mo.ModelName = ? && (v.SubmodelID = 0 || v.SubmodelID is null) &&
		(v.ConfigID = 0 || v.ConfigID is null)`
	stmtEnd := `order by vp.PartNumber`
	brandStmt := " && p.brandID in ("
	for _, b := range l.Brands {
		brandStmt += strconv.Itoa(b) + ","
	}
	brandStmt = strings.TrimRight(brandStmt, ",") + ")"
	wholeStmt := stmtBeginning + brandStmt + stmtEnd

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		ch <- err
		return
	}
	defer db.Close()

	stmt, err := db.Prepare(wholeStmt)
	if err != nil {
		ch <- err
		return
	}

	rows, err := stmt.Query(l.Vehicle.Base.Year, l.Vehicle.Base.Make, l.Vehicle.Base.Model)
	if err != nil || rows == nil {
		ch <- err
		return
	}

	for rows.Next() {
		var p Part
		if err = rows.Scan(&p.ID); err == nil {
			l.Parts = append(l.Parts, p)
		}
	}
	defer rows.Close()

	ch <- nil
	return
}

func removeDuplicates(xs *[]Part) {
	found := make(map[int]bool)
	j := 0
	for i, x := range *xs {
		if !found[x.ID] {
			found[x.ID] = true
			(*xs)[j] = (*xs)[i]
			j++
		}
	}
	*xs = (*xs)[:j]
}

func (v *Vehicle) GetVcdbID() (int, error) {
	db, err := sql.Open("mysql", database.VcdbConnectionString())
	if err != nil {
		return 0, err
	}
	defer db.Close()

	var row *sql.Row
	if v.Submodel != "" {
		stmt, err := db.Prepare(getVcdbVehicleIDWithSubmodel)
		if err != nil {
			return 0, err
		}
		defer stmt.Close()

		row = stmt.QueryRow(v.Base.Year, v.Base.Make, v.Base.Model, v.Submodel)
		if row == nil {
			return 0, err
		}
	} else {
		stmt, err := db.Prepare(getVcdbVehicleID)
		if err != nil {
			return 0, err
		}
		defer stmt.Close()

		row = stmt.QueryRow(v.Base.Year, v.Base.Make, v.Base.Model)
		if row == nil {
			return 0, err
		}
	}

	var id int
	err = row.Scan(&id)

	return id, err
}

func (v *Vehicle) stringify() string {
	str := fmt.Sprintf("%d %s %s", v.Base.Year, v.Base.Make, v.Base.Model)
	if v.Submodel != "" {
		str = fmt.Sprintf("%s %s", str, v.Submodel)
	}
	if len(v.Configurations) > 0 {
		for _, conf := range v.Configurations {
			if conf.Key != "" && conf.Value != "" {
				str = fmt.Sprintf("%s %s:%s", str, conf.Key, conf.Value)
			}
		}
	}

	return str
}

// Vehicle Inquiry
type VehicleInquiry struct {
	Name     string `json:"name" xml:"name,attr"`
	Category int    `json:"category" xml:"category,attr"`
	Phone    string `json:"phone" xml:"phone,attr"`
	Email    string `json:"email" xml:"email,attr"`
	Vehicle  string `json:"vehicle" xml:"vehicle"`
	Message  string `json:"message" xml:"message"`
}

var (
	insertStmt = `insert into VehicleInquiry(name, category, phone, email, vehicle, message, date_added) values(?,?,?,?,?,?, now())`
)

func (i *VehicleInquiry) Push() error {

	if i.Name == "" {
		return fmt.Errorf("%s", "name is required")
	}
	if i.Category == 0 {
		return fmt.Errorf("%s", "category is required")
	}
	if i.Phone == "" && i.Email == "" {
		return fmt.Errorf("%s", "a form of contact is required")
	}
	if i.Vehicle == "" {
		return fmt.Errorf("%s", "the vehicle of inquiry is required")
	}

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(insertStmt)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(i.Name, i.Category, i.Phone, i.Email, i.Vehicle, i.Message)
	return err
}

func (i *VehicleInquiry) SendEmail(dtx *apicontext.DataContext) error {

	cts, err := contact.GetAllContactTypes(dtx)
	if err != nil {
		return err
	}

	var ct contact.ContactType
	for _, t := range cts {
		if t.Name == "Vehicle Inquiry" {
			ct = t
			break
		}
	}

	// Get Category
	var cat Category
	cat.ID = i.Category
	cat.GetCategory(dtx.APIKey, 1, 1, true, nil, nil, dtx)

	// Start to build email body
	body := fmt.Sprintf("Name: %s\n", i.Name)
	body = fmt.Sprintf("%sEmail: %s\n", body, i.Email)
	body = fmt.Sprintf("%sPhone: %s\n", body, i.Phone)
	body = fmt.Sprintf("%sCategory: %s\n", body, cat.Title)

	// Decode vehicle
	var v Vehicle
	if err := json.Unmarshal([]byte(i.Vehicle), &v); err == nil {
		str := v.stringify()
		if str != "" {
			body = fmt.Sprintf("%sVehicle: %s\n", body, str)
		}
	}

	if i.Message != "" {
		body = fmt.Sprintf("%s\nMessage: %s\n", body, i.Message)
	}

	// Send Email
	return contact.SendEmail(ct, "Email from VehicleInquiry Request Form", body) //contact type id, subject, techSupport

}
