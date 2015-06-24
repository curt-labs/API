package products

import (
	"database/sql"
	"github.com/curt-labs/GoAPI/helpers/apicontext"
	"github.com/curt-labs/GoAPI/helpers/database"
	"sort"
	"strings"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/mgo.v2"
)

const (
	AriesDb = "aries"
)

var (
	initMap     sync.Once
	partMap     = make(map[int]BasicPart, 0)
	partMapStmt = `select p.status, p.dateAdded, p.dateModified, p.shortDesc, p.oldPartNumber, p.partID, p.priceCode, pc.class, p.brandID
                from Part as p
                left join Class as pc on p.classID = pc.classID
                where p.brandID = 3 && p.status in (800,900)`
)

type NoSqlVehicle struct {
	Year            string      `bson:"year" json:"year,omitempty" xml:"year, omitempty"`
	Make            string      `bson:"make" json:"make,omitempty" xml:"make, omitempty"`
	Model           string      `bson:"model" json:"model,omitempty" xml:"model, omitempty"`
	Style           string      `bson:"style" json:"style,omitempty" xml:"style, omitempty"`
	Parts           []BasicPart `bson:"-" json:"parts,omitempty" xml:"parts, omitempty"`
	PartIdentifiers []int       `bson:"parts" json:"parts_ids" xml:"-"`
}

type NoSqlApp struct {
	Year  int    `json:"year,omitempty" xml:"year,omitempty"`
	Make  string `json:"make,omitempty" xml:"make,omitempty"`
	Model string `json:"model,omitempty" xml:"model,omitempty"`
	Style string `json:"style,omitempty" xml:"style,omitempty"`
	Part  int    `json:"part,omitempty" xml:"part,omitempty"`
}

type NoSqlLookup struct {
	Years  []string `json:"available_years,omitempty" xml:"available_years, omitempty"`
	Makes  []string `json:"available_makes,omitempty" xml:"available_makes, omitempty"`
	Models []string `json:"available_models,omitempty" xml:"available_models, omitempty"`
	Styles []string `json:"available_styles,omitempty" xml:"available_styles, omitempty"`
	Parts  []Part   `json:"parts,omitempty" xml:"parts, omitempty"`
	NoSqlVehicle
}

type BasicPart struct {
	ID             int       `json:"id" xml:"id,attr"`
	BrandID        int       `json:"brandId" xml:"brandId,attr"`
	Status         int       `json:"status" xml:"status,attr"`
	PriceCode      string    `json:"price_code" xml:"price_code,attr"`
	Class          string    `json:"class" xml:"class,attr"`
	DateModified   time.Time `json:"date_modified" xml:"date_modified,attr"`
	DateAdded      time.Time `json:"date_added" xml:"date_added,attr"`
	ShortDesc      string    `json:"short_description" xml:"short_description,attr"`
	Featured       bool      `json:"featured,omitempty" xml:"featured,omitempty"`
	AcesPartTypeID int       `json:"acesPartTypeId,omitempty" xml:"acesPartTypeId,omitempty"`
	OldPartNumber  string    `json:"oldPartNumber,omitempty" xml:"oldPartNumber,omitempty"`
	UPC            string    `json:"upc,omitempty" xml:"upc,omitempty"`
}

func initMaps() {
	buildPartMap()
}

func GetAriesVehicleCollections() ([]string, error) {
	session, err := mgo.DialWithInfo(database.AriesMongoConnectionString())
	if err != nil {
		return []string{}, err
	}
	defer session.Close()

	cols, err := session.DB(AriesDb).CollectionNames()
	if err != nil {
		return []string{}, err
	}

	validCols := make([]string, 0)
	for _, col := range cols {
		if !strings.Contains(col, "system") {
			validCols = append(validCols, col)
		}
	}

	return validCols, nil
}

func GetApps(v NoSqlVehicle, collection string) (stage string, vals []string, err error) {

	if v.Year != "" && v.Make != "" && v.Model != "" && v.Style != "" {
		return
	}

	session, err := mgo.DialWithInfo(database.AriesMongoConnectionString())
	if err != nil {
		return
	}
	defer session.Close()

	c := session.DB(AriesDb).C(collection)

	queryMap := make(map[string]interface{})

	if v.Year != "" {
		queryMap["year"] = strings.ToLower(v.Year)
	} else {
		c.Find(queryMap).Distinct("year", &vals)
		sort.Sort(sort.Reverse(sort.StringSlice(vals)))
		stage = "year"
		return
	}

	if v.Make != "" {
		queryMap["make"] = strings.ToLower(v.Make)
	} else {
		c.Find(queryMap).Sort("make").Distinct("make", &vals)
		sort.Strings(vals)
		stage = "make"
		return
	}

	if v.Model != "" {
		queryMap["model"] = strings.ToLower(v.Model)
	} else {
		c.Find(queryMap).Sort("model").Distinct("model", &vals)
		sort.Strings(vals)
		stage = "model"
		return
	}

	c.Find(queryMap).Distinct("style", &vals)
	if len(vals) == 1 && vals[0] == "" {
		vals = []string{}
	}

	sort.Strings(vals)
	stage = "style"

	return
}

func FindVehicles(v NoSqlVehicle, collection string, dtx *apicontext.DataContext) (l NoSqlLookup, err error) {

	l = NoSqlLookup{}

	stage, vals, err := GetApps(v, collection)
	if err != nil {
		return
	}

	if stage != "" {
		switch stage {
		case "year":
			l.Years = vals
		case "make":
			l.Makes = vals
		case "model":
			l.Models = vals
		case "style":
			l.Styles = vals
		}

		if stage != "style" || len(l.Styles) > 0 {
			return
		}
	}

	session, err := mgo.DialWithInfo(database.AriesMongoConnectionString())
	if err != nil {
		return
	}
	defer session.Close()

	c := session.DB(AriesDb).C(collection)
	queryMap := make(map[string]interface{})

	ids := make([]int, 0)
	queryMap["year"] = strings.ToLower(v.Year)
	queryMap["make"] = strings.ToLower(v.Make)
	queryMap["model"] = strings.ToLower(v.Model)
	queryMap["style"] = strings.ToLower(v.Style)

	c.Find(queryMap).Distinct("parts", &ids)

	//add parts
	for _, id := range ids {
		p := Part{ID: id}
		if err := p.Get(dtx); err != nil {
			continue
		}
		l.Parts = append(l.Parts, p)
	}

	return l, err
}

func FindApplications(collection string) ([]NoSqlVehicle, error) {
	initMap.Do(initMaps)

	var apps []NoSqlVehicle
	var err error

	session, err := mgo.DialWithInfo(database.AriesMongoConnectionString())
	if err != nil {
		return apps, err
	}
	defer session.Close()

	c := session.DB(AriesDb).C(collection)

	err = c.Find(nil).Sort("make", "model", "style", "-year").All(&apps)

	fulfilled := make([]NoSqlVehicle, 0)
	for _, app := range apps {
		for _, p := range app.PartIdentifiers {
			if part, ok := partMap[p]; ok {
				app.Parts = append(app.Parts, part)
			}
		}
		if len(app.Parts) > 0 {
			fulfilled = append(fulfilled, app)
		}
	}

	return fulfilled, err
}

func buildPartMap() error {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	qry, err := db.Prepare(partMapStmt)
	if err != nil {
		return err
	}
	defer qry.Close()

	rows, err := qry.Query()
	if err != nil || rows == nil {
		return err
	}

	for rows.Next() {
		var p BasicPart
		var priceCode, class *string
		err = rows.Scan(
			&p.Status,
			&p.DateAdded,
			&p.DateModified,
			&p.ShortDesc,
			&p.OldPartNumber,
			&p.ID,
			&priceCode,
			&class,
			&p.BrandID,
		)
		if err != nil {
			continue
		}
		if priceCode != nil {
			p.PriceCode = *priceCode
		}
		if class != nil {
			p.Class = *class
		}

		partMap[p.ID] = p
	}
	rows.Close()

	return nil
}
