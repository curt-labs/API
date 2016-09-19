package products

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/curt-labs/API/helpers/database"
	"github.com/curt-labs/API/helpers/redis"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// LuverneLookupContext Holds required configuration settings and resources.
type LuverneLookupContext struct {
	Session  *mgo.Session
	Statuses []int
}

// LuverneVehicleBase The lowest level vehicle properties used to query
// a LuverneCategoryVehicle.
type LuverneVehicleBase struct {
	Year  string `json:"year" xml:"year,attr"`
	Make  string `json:"make" xml:"make,attr"`
	Model string `json:"model" xml:"model,attr"`
}

// LuverneCategoryVehicle Represents the requested `Base` vehicle and all matching
// LookupCategory types that are fitments of the `Base`.
type LuverneCategoryVehicle struct {
	Base CategoryVehicleBase `json:"base" xml:"BaseVehicle"`

	Years      []string         `json:"availableYears,omitempty" xml:"availableYears,omitempty"`
	Makes      []string         `json:"availableMakes,omitempty" xml:"availableMakes,omitempty"`
	Models     []string         `json:"availableModels,omitempty" xml:"availableModels,omitempty"`
	Categories []LookupCategory `json:"lookup_category" xml:"StyleOptions"`
	Products   []Part           `json:"products" xml:"products"`
}

// LuverneLookupCategory Represents a specific category of `StyleOption` fitments.
type LuverneLookupCategory struct {
	Category     Category      `json:"category" xml:"category"`
	StyleOptions []StyleOption `json:"style_options" xml:"StyleOptions"`
}

// // LuverneStyleOption Matches a slice of `Part` that have equal fitments to the matched
// // `Style`.
// type LuverneStyleOption struct {
// 	Style          string           `json:"style" xml:"Style"`
// 	FitmentNumbers []FitmentMapping `json:"fitments" xml:"fitments"`
// }
//
// // FitmentMapping defines the product and any associated attributes that
// // relevant to the fitment on a specific application (install time, drilling, etc)
// type FitmentMapping struct {
// 	Attributes []FitmentAttribute `json:"attributes" xml:"Attributes"`
// 	Number     string             `json:"product_identifier" xml:"product_identifier"`
// }
//
// // FitmentAttribute A name value for a note of a fitment application.
// type FitmentAttribute struct {
// 	Key   string `json:"key" xml:"Key"`
// 	Value string `json:"value" xml:"Value"`
// }

// Query Returns a `CategoryVehicle` that holds matching information for the
// queried `CategoryVehicleBase` attributes.
func LuverneQuery(ctx *LuverneLookupContext, args ...string) (*CategoryVehicle, error) {

	var redisKey string
	// var category string
	var vehicle CategoryVehicle

	for i, arg := range args {
		if i == 0 {
			redisKey = fmt.Sprintf("luverne:%s", arg)
		} else {
			redisKey = fmt.Sprintf("%s:%s", redisKey, arg)
		}
	}
	data, err := redis.Get(redisKey)
	if err == nil && len(data) > 0 {
		err = json.Unmarshal(data, &vehicle)
		if err == nil {
			return &vehicle, nil
		}
	}

	switch len(args) {
	case 1:
		vehicle.Base.Year = args[0]
	case 2:
		vehicle.Base.Year = args[0]
		vehicle.Base.Make = args[1]
	case 3:
		vehicle.Base.Year = args[0]
		vehicle.Base.Make = args[1]
		vehicle.Base.Model = args[2]
	case 4:
		vehicle.Base.Year = args[0]
		vehicle.Base.Make = args[1]
		vehicle.Base.Model = args[2]
		// category = args[3]
	}

	if vehicle.Base.Year == "" {
		vehicle.Years, err = getLuverneYears(ctx)
	} else if vehicle.Base.Year != "" && vehicle.Base.Make == "" {
		vehicle.Makes, err = getLuverneMakes(ctx, vehicle.Base.Year)
	} else if vehicle.Base.Year != "" && vehicle.Base.Make != "" && vehicle.Base.Model == "" {
		vehicle.Models, err = getLuverneModels(ctx, vehicle.Base.Year, vehicle.Base.Make)
	}

	redis.Setex(redisKey, vehicle, 60*60*24)

	return &vehicle, err
}

func getLuverneYears(ctx *LuverneLookupContext) ([]string, error) {
	if ctx == nil {
		return nil, fmt.Errorf("missing context")
	} else if ctx.Session == nil {
		return nil, fmt.Errorf("invalid mongodb connection")
	}

	c := ctx.Session.DB(database.ProductMongoDatabase).C(database.ProductCollectionName)

	qry := bson.M{
		"status": bson.M{
			"$in": ctx.Statuses,
		},
		"brand.id": 4,
	}

	var res []string
	err := c.Find(qry).Select(bson.M{
		"luverne_applications.year": 1,
		"_id": -1,
	}).Distinct("luverne_applications.year", &res)

	if err != nil {
		return nil, err
	}

	sort.Sort(sort.Reverse(sort.StringSlice(res)))

	return res, err
}

func getLuverneMakes(ctx *LuverneLookupContext, year string) ([]string, error) {
	if ctx == nil {
		return nil, fmt.Errorf("missing context")
	} else if ctx.Session == nil {
		return nil, fmt.Errorf("invalid mongodb connection")
	}

	c := ctx.Session.DB(database.ProductMongoDatabase).C(database.ProductCollectionName)

	type Apps struct {
		Apps []VehicleApplication `bson:"luverne_applications"`
	}

	var apps []Apps
	qry := bson.M{
		"luverne_applications": bson.M{
			"$elemMatch": bson.M{
				"year": year,
			},
		},
		"status": bson.M{
			"$in": ctx.Statuses,
		},
		"brand.id": 4,
	}
	err := c.Find(qry).Select(bson.M{"luverne_applications.make": 1, "_id": 0}).All(&apps)
	if err != nil {
		return nil, err
	}

	var makes []string

	existing := make(map[string]string, 0)
	for _, app := range apps {
		for _, a := range app.Apps {
			a.Make = strings.Title(a.Make)
			if _, ok := existing[a.Make]; !ok {
				makes = append(makes, a.Make)
				existing[a.Make] = a.Make
			}
		}
	}

	sort.Strings(makes)

	return makes, err
}

func getLuverneModels(ctx *LuverneLookupContext, year, vehicleMake string) ([]string, error) {
	if ctx == nil {
		return nil, fmt.Errorf("missing context")
	} else if ctx.Session == nil {
		return nil, fmt.Errorf("invalid mongodb connection")
	}

	c := ctx.Session.DB(database.ProductMongoDatabase).C(database.ProductCollectionName)

	type Apps struct {
		Apps []VehicleApplication `bson:"luverne_applications"`
	}

	var apps []Apps
	err := c.Find(bson.M{
		"luverne_applications": bson.M{
			"$elemMatch": bson.M{
				"year": year,
				"make": bson.RegEx{
					Pattern: "^" + vehicleMake + "$",
					Options: "i",
				},
			},
		},
		"status": bson.M{
			"$in": ctx.Statuses,
		},
		"brand.id": 4,
	}).Select(bson.M{"luverne_applications.$.model": 1, "_id": 0}).All(&apps)
	if err != nil {
		return nil, err
	}

	var models []string

	existing := make(map[string]string, 0)
	for _, app := range apps {
		for _, a := range app.Apps {
			a.Model = strings.Title(a.Model)
			if _, ok := existing[a.Model]; !ok {
				models = append(models, a.Model)
				existing[a.Model] = a.Model
			}
		}
	}

	sort.Strings(models)

	return models, err
}
