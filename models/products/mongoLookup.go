package products

import (
	"github.com/curt-labs/GoAPI/helpers/apicontext"
	"github.com/curt-labs/GoAPI/helpers/database"
	"sort"
	"strings"

	"gopkg.in/mgo.v2"
)

const (
	AriesDb = "aries"
)

type NoSqlVehicle struct {
	Year  string `bson:"year" json:"year,omitempty" xml:"year, omitempty"`
	Make  string `bson:"make" json:"make,omitempty" xml:"make, omitempty"`
	Model string `bson:"model" json:"model,omitempty" xml:"model, omitempty"`
	Style string `bson:"style" json:"style,omitempty" xml:"style, omitempty"`
	Parts []Part `bson:"parts" json:"parts,omitempty" xml:"parts, omitempty"`
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
	if len(vals) == 1 {
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

	if v.Style == "ANYTHING_YOU_WANT" {
		v.Style = ""
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

	ch := make(chan *Part)
	//add parts
	for _, id := range ids {
		go func(i int) {
			p := &Part{ID: i}
			if err := p.Get(dtx); err != nil {
				ch <- nil
			} else {
				ch <- p
			}
		}(id)
	}

	for _, _ = range ids {
		res := <-ch
		if res != nil {
			l.Parts = append(l.Parts, *res)
		}
	}

	return l, err
}
