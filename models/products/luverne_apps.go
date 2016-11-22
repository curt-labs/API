package products

import (
	"strings"

	"github.com/curt-labs/API/helpers/database"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type LuverneResult struct {
	Applications []NoSqlLuverneVehicle `json:"applications"`
	Finishes     []string              `json:"finishes"`
}

type NoSqlLuverneVehicle struct {
	ID      bson.ObjectId  `bson:"_id" json:"_id" xml:"_id"`
	Year    string         `bson:"year" json:"year,omitempty" xml:"year, omitempty"`
	Make    string         `bson:"make" json:"make,omitempty" xml:"make, omitempty"`
	Model   string         `bson:"model" json:"model,omitempty" xml:"model, omitempty"`
	Body    string         `bson:"body" json:"body,omitempty" xml:"body, omitempty"`
	Box     string         `bson:"boxLength" json:"boxLength,omitempty" xml:"boxLength, omitempty"`
	Cab     string         `bson:"cabLength" json:"cabLength,omitempty" xml:"cabLength, omitempty"`
	Fuel    string         `bson:"fuelType" json:"fuelType,omitempty" xml:"fuelType, omitempty"`
	Wheel   string         `bson:"wheelType" json:"wheelType,omitempty" xml:"wheelType, omitempty"`
	Parts   []*CompactPart `bson:"products" json:"parts,omitempty" xml:"parts, omitempty"`
	MinYear string         `bson:"min_year" json:"min_year" xml:"minYear"`
	MaxYear string         `bson:"max_year" json:"max_year" xml:"maxYear"`
}

type CompactPart struct {
	PartNumber   string `bson:"part_number" json:"part_number" xml:"part_number"`
	InstallSheet string `bson:"install_sheet" json:"install_sheet" xml:"install_sheet"`
	ShortDesc    string `bson:"short_description" json:"short_description" xml:"short_description"`
	Finish       string `bson:"finish" json:"finish" xml:"finish"`
	Attributes   []Atr  `bson:"attributes" json:"attributes" xml:"attributes"`
}

type Atr struct {
	Name  string `bson:"name" json:"name" xml:"name"`
	Value string `bson:"value" json:"value" xml:"value"`
}

func FindAppsLuverne(catID, skip, limit int) (LuverneResult, error) {
	res := LuverneResult{
		Applications: make([]NoSqlLuverneVehicle, 0),
		Finishes:     make([]string, 0),
	}

	if limit == 0 || limit > 100 {
		limit = 100
	}

	var apps []NoSqlLuverneVehicle
	var err error

	session, err := mgo.DialWithInfo(database.AriesMongoConnectionString())
	if err != nil {
		return res, err
	}
	defer session.Close()

	c := session.DB(database.ProductMongoDatabase).C(database.ProductCollectionName)

	pipe := c.Pipe([]bson.D{
		bson.D{{"$match", bson.M{"brand.id": 4, "categories.id": catID}}},
		bson.D{{"$unwind", "$luverne_applications"}},
		bson.D{
			{
				"$group", bson.M{
					"_id": bson.M{
						"make":      "$luverne_applications.make",
						"model":     "$luverne_applications.model",
						"body":      "$luverne_applications.body",
						"boxLength": "$luverne_applications.boxLength",
						"cabLength": "$luverne_applications.cabLength",
						"fuelType":  "$luverne_applications.fuelType",
						"wheelType": "$luverne_applications.wheelType",
					},
					"min_year": bson.M{"$min": "$luverne_applications.year"},
					"max_year": bson.M{"$max": "$luverne_applications.year"},
					"products": bson.M{
						"$addToSet": bson.M{
							"part_number":       "$part_number",
							"install_sheet":     "$install_sheet",
							"short_description": "$short_description",
							"attributes":        "$attributes",
						},
					},
				},
			},
		},
		bson.D{
			{
				"$project", bson.M{
					"make":      bson.M{"$toUpper": "$_id.make"},
					"model":     bson.M{"$toUpper": "$_id.model"},
					"body":      bson.M{"$toUpper": "$_id.body"},
					"boxLength": bson.M{"$toUpper": "$_id.boxLength"},
					"cabLength": bson.M{"$toUpper": "$_id.cabLength"},
					"fuelType":  bson.M{"$toUpper": "$_id.fuelType"},
					"wheelType": bson.M{"$toUpper": "$_id.wheelType"},
					"products":  1,
					"min_year":  1,
					"max_year":  1,
					"_id":       0,
				},
			},
		},
		bson.D{
			{
				"$group", bson.M{
					"_id": bson.M{
						"min_year":  "$min_year",
						"max_year":  "$max_year",
						"make":      "$make",
						"model":     "$model",
						"body":      "$body",
						"boxLength": "$boxLength",
						"cabLength": "$cabLength",
						"fuelType":  "$fuelType",
						"wheelType": "$wheelType",
					},
					"products":  bson.M{"$first": "$products"},
					"make":      bson.M{"$first": "$make"},
					"model":     bson.M{"$first": "$model"},
					"body":      bson.M{"$first": "$body"},
					"boxLength": bson.M{"$first": "$boxLength"},
					"cabLength": bson.M{"$first": "$cabLength"},
					"fuelType":  bson.M{"$first": "$fuelType"},
					"wheelType": bson.M{"$first": "$wheelType"},
					"min_year":  bson.M{"$min": "$min_year"},
					"max_year":  bson.M{"$max": "$max_year"},
				},
			},
		},
		bson.D{
			{
				"$sort", bson.D{
					{"_id.make", 1},
					{"_id.model", 1},
					{"_id.body", 1},
					{"_id.boxLength", 1},
					{"_id.cabLength", 1},
					{"_id.fuelType", 1},
					{"_id.wheelType", 1},
				},
			},
		},
		bson.D{{"$skip", skip}},
		bson.D{{"$limit", limit}},
	})
	err = pipe.All(&apps)
	if err != nil {
		return res, err
	}
	res.Applications = apps

	// find all available finishes
	var finishesMap = make(map[string]bool, 0)
	for _, app := range apps {
		for _, p := range app.Parts {
			for _, atr := range p.Attributes {
				if strings.ToLower(atr.Name) == "finish" {
					finishesMap[atr.Value] = true
					p.Finish = atr.Value
				}
			}
		}
	}
	// loop over the keys of the map and assign them as the values for the res.finishes array
	for k, _ := range finishesMap {
		res.Finishes = append(res.Finishes, k)
	}

	return res, nil
}
