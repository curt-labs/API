package products

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"

	"github.com/curt-labs/API/helpers/apicontext"
	"github.com/curt-labs/API/helpers/database"
	"github.com/curt-labs/API/helpers/sortutil"
	"github.com/curt-labs/API/models/contact"
	_ "github.com/go-sql-driver/mysql"
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

	partMatcherStmt = `select p.partID, cat.name, ca.value from vcdb_Vehicle as v
		left join VehicleConfigAttribute as vca on v.ConfigID = vca.VehicleConfigID
		left join ConfigAttribute as ca on vca.AttributeID = ca.ID
		left join ConfigAttributeType as cat on ca.ConfigAttributeTypeID = cat.ID
		left join Submodel as sm on v.SubmodelID = sm.ID
		join vcdb_VehiclePart as vp on v.ID = vp.VehicleID
		join Part as p on vp.PartNumber = p.partID
		join BaseVehicle as bv on v.BaseVehicleID = bv.ID
		join vcdb_Make as ma on bv.MakeID = ma.ID
		join vcdb_Model as mo on bv.ModelID = mo.ID
		where bv.YearID = ? && ma.MakeName = ? && mo.ModelName = ? &&
		(sm.SubmodelName = ? || sm.ID is null) && p.brandID in (?) &&
		(p.status = 800 || p.status = 900)
		group by p.partID, cat.name, ca.value
		order by p.partID, cat.name, ca.value`
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
		count = 50
	}

	err := database.Init()
	if err != nil {
		ch <- nil
		return
	}

	stmt, err := database.DB.Prepare(partMatcherStmt)
	if err != nil {
		ch <- nil
		return
	}
	defer stmt.Close()

	brands := make([]string, 0)
	for _, b := range apicontext.AllBrandsArray {
		brands = append(brands, strconv.Itoa(b))
	}

	rows, err := stmt.Query(l.Vehicle.Base.Year, l.Vehicle.Base.Make, l.Vehicle.Base.Model, l.Vehicle.Submodel, strings.Join(brands, ","))
	if err != nil || rows == nil {
		ch <- nil
		return
	}
	defer rows.Close()

	maps := make(map[int][]ConfigurationOption, 0)
	parts := make([]int, 0)

	// Compile configuration map from database results
	for rows.Next() {
		var part int
		var config_type *string
		var config_val *string
		if err := rows.Scan(&part, &config_type, &config_val); err != nil {
			continue
		}
		if part == 0 {
			continue
		}
		if config_type == nil || config_val == nil || *config_type == "" || *config_val == "" {
			parts = append(parts, part)
			continue
		}

		opt := ConfigurationOption{
			Type:    strings.TrimSpace(*config_type),
			Options: make([]string, 0),
		}

		if maps[part] == nil {
			maps[part] = make([]ConfigurationOption, 0)
			opt.Options = append(opt.Options, *config_val)
			maps[part] = append(maps[part], opt)
		} else {
			for i, conf := range maps[part] {
				if strings.ToLower(strings.TrimSpace(conf.Type)) == strings.ToLower(strings.TrimSpace(*config_type)) {
					maps[part][i].Options = append(maps[part][i].Options, *config_val)
					continue
				}
			}
		}
	}

	// index the qualified configurations
	confIndex := make(map[string]string, 0)
	for _, c := range l.Vehicle.Configurations {
		confIndex[strings.ToLower(strings.TrimSpace(c.Key))] = strings.TrimSpace(c.Value)
	}

	// run a comparison of the database results
	// against the qualified configurations,
	// storing part numbers that are fully matched.
	for part, configs := range maps {
		qualified := 0
		for _, config := range configs {
			if val := confIndex[strings.ToLower(strings.TrimSpace(config.Type))]; val != "" {

				for _, opt := range config.Options {
					if strings.ToLower(strings.TrimSpace(val)) == strings.ToLower(strings.TrimSpace(opt)) {
						qualified = qualified + 1
					}
				}
			}
		}

		if qualified == len(configs) {
			parts = append(parts, part)
		}
	}

	sort.Ints(parts)

	partCount := len(parts)
	pagedParts := parts
	if partCount > count {
		start := 0
		if page > 1 {
			start = count * (page - 1)
		}
		end := start + count
		if len(parts) <= end {
			pagedParts = parts[start:]
		} else {
			pagedParts = parts[start : start+count]
		}
	}

	l.Parts = make([]Part, 0)
	perChan := make(chan int)
	for i, p := range pagedParts {
		go func(j int, prt Part) {
			if err := prt.Get(dtx); err == nil && prt.ShortDesc != "" {
				l.Parts = append(l.Parts, prt)
			}
			perChan <- 1
		}(i, Part{ID: p})
	}

	for _, _ = range pagedParts {
		<-perChan
	}

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

	ch <- nil
	return
}

func (v *Vehicle) GetVcdbID() (int, error) {
	err := database.Init()
	if err != nil {
		return 0, err
	}

	var row *sql.Row
	if v.Submodel != "" {
		stmt, err := database.VcdbDB.Prepare(getVcdbVehicleIDWithSubmodel)
		if err != nil {
			return 0, err
		}
		defer stmt.Close()

		row = stmt.QueryRow(v.Base.Year, v.Base.Make, v.Base.Model, v.Submodel)
		if row == nil {
			return 0, err
		}
	} else {
		stmt, err := database.VcdbDB.Prepare(getVcdbVehicleID)
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

	err := database.Init()
	if err != nil {
		return err
	}

	stmt, err := database.DB.Prepare(insertStmt)
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
	// var cat Category
	// cat.CategoryID = i.Category
	// cat.GetCategory(dtx.APIKey, 1, 1, true, nil, nil, dtx)

	// Start to build email body
	body := fmt.Sprintf("Name: %s\n", i.Name)
	body = fmt.Sprintf("%sEmail: %s\n", body, i.Email)
	body = fmt.Sprintf("%sPhone: %s\n", body, i.Phone)
	// body = fmt.Sprintf("%sCategory: %s\n", body, cat.Title)

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
