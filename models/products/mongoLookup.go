package products

import (
	"database/sql"

	"github.com/curt-labs/API/helpers/apicontext"
	"github.com/curt-labs/API/helpers/database"

	"sort"
	"strings"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	AriesDb = "aries"
)

var (
	initMap     sync.Once
	finishes    = make(map[string]string, 0)
	colors      = make(map[string]string, 0)
	partMap     = make(map[int]BasicPart, 0)
	partMapStmt = `select p.status, p.dateAdded, p.dateModified, p.shortDesc, p.oldPartNumber, p.partID, p.priceCode, pc.class, p.brandID, c.catTitle, (
						select group_concat(pa.value) from PartAttribute as pa
						where pa.partID = p.partID && pa.field = 'Finish'
					) as finish,
					(
						select group_concat(pa.value) from PartAttribute as pa
						where pa.partID = p.partID && pa.field = 'Color'
					) as color,
					(
						select group_concat(pa.value) from PartAttribute as pa
						where pa.partID = p.partID && pa.field = 'Location'
					) as location,
					con.text as installSheet
					from Part as p
					left join Class as pc on p.classID = pc.classID
					left join CatPart as cp on p.partID = cp.partID
					left join Categories as c on cp.catID = c.catID
					left join ContentBridge as cb on p.partID = cb.partID
					left join Content as con on cb.contentID = con.contentID && con.cTypeID = 5
					where p.brandID = 3 && p.status in (800,900)`
)

type NoSqlVehicle struct {
	ID              bson.ObjectId `bson:"_id" json:"_id" xml:"_id"`
	Year            string        `bson:"year" json:"year,omitempty" xml:"year, omitempty"`
	Make            string        `bson:"make" json:"make,omitempty" xml:"make, omitempty"`
	Model           string        `bson:"model" json:"model,omitempty" xml:"model, omitempty"`
	Style           string        `bson:"style" json:"style,omitempty" xml:"style, omitempty"`
	Parts           []BasicPart   `bson:"-" json:"parts,omitempty" xml:"parts, omitempty"`
	PartIdentifiers []int         `bson:"parts" json:"parts_ids" xml:"-"`
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
	Finish         string    `json:"finish"`
	Color          string    `json:"color"`
	Category       string    `json:"category"`
	Location       string    `json:"location"`
	InstallSheet   string    `json:"install_sheet"`
}

type Result struct {
	Applications []NoSqlVehicle `json:"applications"`
	Finishes     []string       `json:"finishes"`
	Colors       []string       `json:"colors"`
}

func initMaps() {
	buildPartMap()
}

func GetAriesVehicleCollections(session *mgo.Session) ([]string, error) {

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

	if err = database.Init(); err != nil {
		return
	}

	session := database.AriesMongoSession.Copy()
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

func FindApplications(collection string, skip, limit int) (Result, error) {
	initMap.Do(initMaps)

	if limit == 0 || limit > 100 {
		limit = 100
	}

	res := Result{
		Applications: make([]NoSqlVehicle, 0),
		Finishes:     make([]string, 0),
		Colors:       make([]string, 0),
	}

	var apps []NoSqlVehicle
	var err error

	session, err := mgo.DialWithInfo(database.AriesMongoConnectionString())
	if err != nil {
		return res, err
	}
	defer session.Close()

	c := session.DB(AriesDb).C(collection)

	err = c.Find(nil).Sort("make", "model", "style", "-year").Skip(skip).Limit(limit).All(&apps)

	existingFinishes := make(map[string]string, 0)
	existingColors := make(map[string]string, 0)
	for _, app := range apps {
		for _, p := range app.PartIdentifiers {
			if part, ok := partMap[p]; ok {
				app.Parts = append(app.Parts, part)

				_, ok := existingFinishes[part.Finish]
				if part.Finish != "" && !ok {
					res.Finishes = append(res.Finishes, part.Finish)
					existingFinishes[part.Finish] = part.Finish
				}
				_, ok = existingColors[part.Color]
				if part.Color != "" && !ok {
					res.Colors = append(res.Colors, part.Color)
					existingColors[part.Color] = part.Color
				}
			}
		}
		if len(app.Parts) > 0 {
			res.Applications = append(res.Applications, app)
		}
	}

	return res, err
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
	defer rows.Close()

	for rows.Next() {
		var p BasicPart
		var priceCode, cat, class, finish, color, location, install *string
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
			&cat,
			&finish,
			&color,
			&location,
			&install,
		)
		if err != nil {
			continue
		}
		if install != nil {
			p.InstallSheet = *install
		}
		if priceCode != nil {
			p.PriceCode = *priceCode
		}
		if class != nil {
			p.Class = *class
		}
		if cat != nil {
			p.Category = *cat
		}
		if finish != nil {
			p.Finish = *finish
			if _, ok := finishes[p.Finish]; !ok {
				finishes[p.Finish] = p.Finish
			}
		}
		if color != nil {
			p.Color = *color
			if _, ok := colors[p.Color]; !ok {
				colors[p.Color] = p.Color
			}
		}
		if location != nil {
			p.Location = *location
		}

		partMap[p.ID] = p
	}

	return nil
}

func FindVehiclesWithParts(v NoSqlVehicle, collection string, dtx *apicontext.DataContext, sess *mgo.Session) (l NoSqlLookup, err error) {

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
	}
	//
	c := sess.DB(AriesDb).C(collection)
	queryMap := make(map[string]interface{})

	ids := make([]int, 0)
	queryMap["year"] = strings.ToLower(v.Year)
	queryMap["make"] = strings.ToLower(v.Make)
	queryMap["model"] = strings.ToLower(v.Model)
	if v.Style != "" {
		queryMap["style"] = strings.ToLower(v.Style)
	}

	c.Find(queryMap).Distinct("parts", &ids)

	l.Parts, err = GetMany(ids, getBrandsFromDTX(dtx), sess)

	return l, err
}

//from each category:
//if no v.style:
//query base vehicle
//get parts & available_styles
//else:
//query base+style
//get parts

func FindVehiclesFromAllCategories(v NoSqlVehicle, dtx *apicontext.DataContext, sess *mgo.Session) (map[string]NoSqlLookup, error) {
	var l NoSqlLookup
	lookupMap := make(map[string]NoSqlLookup)

	//Get all collections
	cols, err := GetAriesVehicleCollections(sess)
	if err != nil {
		return lookupMap, err
	}

	//from each category
	for _, col := range cols {

		c := sess.DB(AriesDb).C(col)
		queryMap := make(map[string]interface{})
		//query base vehicle
		queryMap["year"] = strings.ToLower(v.Year)
		queryMap["make"] = strings.ToLower(v.Make)
		queryMap["model"] = strings.ToLower(v.Model)
		if (v.Style) != "" {
			queryMap["style"] = strings.ToLower(v.Style)
		} else {
			_, l.Styles, err = GetApps(v, col)
			if err != nil {
				continue
			}
		}

		var ids []int
		err = c.Find(queryMap).Distinct("parts", &ids)
		if err != nil {
			continue
		}
		//add parts
		var partsArray []Part
		l.Parts, err = GetMany(ids, getBrandsFromDTX(dtx), sess)
		if err != nil {
			continue
		}

		partsArray = append(partsArray, l.Parts...)
		if len(partsArray) > 0 {
			var tmp = lookupMap[col]
			tmp.Parts = partsArray
			tmp.Styles = l.Styles
			lookupMap[col] = tmp
			partsArray = nil
		}
	}
	return lookupMap, err
}

func FindPartsFromOneCategory(v NoSqlVehicle, collection string, dtx *apicontext.DataContext, sess *mgo.Session) (map[string]NoSqlLookup, error) {
	var l NoSqlLookup
	var err error
	lookupMap := make(map[string]NoSqlLookup)

	c := sess.DB(AriesDb).C(collection)
	queryMap := make(map[string]interface{})
	//query base vehicle
	queryMap["year"] = strings.ToLower(v.Year)
	queryMap["make"] = strings.ToLower(v.Make)
	queryMap["model"] = strings.ToLower(v.Model)
	if (v.Style) != "" {
		queryMap["style"] = strings.ToLower(v.Style)
	} else {
		_, l.Styles, err = GetApps(v, collection)
		if err != nil {
			return lookupMap, err
		}
	}

	var ids []int
	c.Find(queryMap).Distinct("parts", &ids)
	//add parts

	l.Parts, err = GetMany(ids, getBrandsFromDTX(dtx), sess)
	if err != nil {
		return lookupMap, err
	}

	if len(l.Parts) > 0 {
		var tmp = lookupMap[collection]
		tmp.Parts = l.Parts
		tmp.Styles = l.Styles
		lookupMap[collection] = tmp
	}
	return lookupMap, err
}
