package products

import (
	"github.com/curt-labs/GoAPI/helpers/apicontext"
	"github.com/curt-labs/GoAPI/helpers/database"

	"gopkg.in/mgo.v2"
)

type NoSqlVehicle struct {
	Year  int    `bson:"year" json:"year,omitempty" xml:"year, omitempty"`
	Make  string `bson:"make" json:"make,omitempty" xml:"make, omitempty"`
	Model string `bson:"model" json:"model,omitempty" xml:"model, omitempty"`
	Style string `bson:"style" json:"style,omitempty" xml:"style, omitempty"`
	Parts []Part `bson:"parts" json:"parts,omitempty" xml:"parts, omitempty"`
}

type NoSqlApp struct {
	Year  int    `json:"year,omitempty" xml:"year, omitempty"`
	Make  string `json:"make,omitempty" xml:"make, omitempty"`
	Model string `json:"model,omitempty" xml:"model, omitempty"`
	Style string `json:"style,omitempty" xml:"style, omitempty"`
	Part  int    `json:"part,omitempty" xml:"part, omitempty"`
}

type NoSqlLookup struct {
	Years       []int    `json:"available_years,omitempty" xml:"available_years, omitempty"`
	Makes       []string `json:"available_makes,omitempty" xml:"available_makes, omitempty"`
	Models      []string `json:"available_models,omitempty" xml:"available_models, omitempty"`
	Styles      []string `json:"available_styles,omitempty" xml:"available_styles, omitempty"`
	PartNumbers []int    `json:"partNumbers,omitempty" xml:"partNumbers, omitempty"`
	NoSqlVehicle
}

func FindVehicles(v NoSqlVehicle, collection string, dtx *apicontext.DataContext) (l NoSqlLookup, err error) {
	var apps []NoSqlApp

	session, err := mgo.DialWithInfo(database.MongoConnectionString())
	if err != nil {
		return l, err
	}
	defer session.Close()
	c := session.DB("ariesimport").C(collection)

	queryMap := make(map[string]interface{})

	if v.Year > 0 {
		queryMap["year"] = v.Year
	}

	if v.Make != "" {
		queryMap["make"] = v.Make
	}

	if v.Model != "" {
		queryMap["model"] = v.Model
	}

	if v.Style != "" {
		queryMap["style"] = v.Style
	}

	c.Find(queryMap).All(&apps)

	//ugly maps, dude
	yearMap := make(map[int]int)
	makeMap := make(map[string]int)
	modelMap := make(map[string]int)
	styleMap := make(map[string]int)
	partMap := make(map[int]int)

	for _, app := range apps {
		yearMap[app.Year] = 0
		makeMap[app.Make] = 0
		modelMap[app.Model] = 0
		styleMap[app.Style] = 0
		partMap[app.Part] = 0
	}

	//add to lookup response
	if len(yearMap) == 1 {
		for year, _ := range yearMap {
			l.NoSqlVehicle.Year = year
		}
	} else {
		for year, _ := range yearMap {
			l.Years = append(l.Years, year)
		}
	}

	if len(makeMap) == 1 {
		for ma, _ := range makeMap {
			l.NoSqlVehicle.Make = ma
		}
	} else {
		for ma, _ := range makeMap {
			l.Makes = append(l.Makes, ma)
		}
	}

	if len(modelMap) == 1 {
		for mo, _ := range modelMap {
			l.NoSqlVehicle.Model = mo
		}
	} else {
		for mo, _ := range modelMap {
			l.Models = append(l.Models, mo)
		}
	}

	if len(styleMap) <= 1 {
		for s, _ := range styleMap {
			l.NoSqlVehicle.Style = s
		}
	} else {
		for s, _ := range styleMap {
			l.Styles = append(l.Styles, s)
		}
	}

	// for p, _ := range partMap {
	// 	l.PartNumbers = append(l.PartNumbers, p)
	// }

	//add parts
	if len(yearMap) <= 1 && len(makeMap) <= 1 && len(modelMap) <= 1 && len(styleMap) <= 1 {
		for p, _ := range partMap {
			part := Part{ID: p}
			l.Parts = append(l.Parts, part)
		}
	}

	//get part details
	for i, _ := range l.Parts {
		err = l.Parts[i].Get(dtx)
		if err != nil {
			continue
		}
	}

	return l, err
}
