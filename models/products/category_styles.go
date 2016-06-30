package products

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/curt-labs/API/helpers/database"
	"github.com/curt-labs/API/helpers/redis"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// AllPlaceholder Delete style value
const AllPlaceholder = "All"

// LookupContext Holds required configuration settings and resources.
type LookupContext struct {
	Session  *mgo.Session
	Statuses []int
	Brands   []int
}

// CategoryVehicleBase The lowest level vehicle properties used to query
// a CategoryVehicle.
type CategoryVehicleBase struct {
	Year  string `json:"year" xml:"year,attr"`
	Make  string `json:"make" xml:"make,attr"`
	Model string `json:"model" xml:"model,attr"`
}

// CategoryVehicle Represents the requested `Base` vehicle and all matching
// LookupCategory types that are fitments of the `Base`.
type CategoryVehicle struct {
	Base CategoryVehicleBase `json:"base" xml:"BaseVehicle"`

	Years      []string         `json:"availableYears,omitempty" xml:"availableYears,omitempty"`
	Makes      []string         `json:"availableMakes,omitempty" xml:"availableMakes,omitempty"`
	Models     []string         `json:"availableModels,omitempty" xml:"availableModels,omitempty"`
	Categories []LookupCategory `json:"lookup_category" xml:"StyleOptions"`
}

// LookupCategory Represents a specific category of `StyleOption` fitments.
type LookupCategory struct {
	Category     Category      `json:"category" xml:"category"`
	StyleOptions []StyleOption `json:"style_options" xml:"StyleOptions"`
}

// StyleOption Matches a slice of `Part` that have equal fitments to the matched
// `Style`.
type StyleOption struct {
	Style    string           `json:"style" xml:"Style"`
	Fitments []FitmentProduct `json:"products" xml:"Products"`
}

// FitmentProduct Defines the matching product along with any application specific
// attributes.
type FitmentProduct struct {
	Attributes []FitmentAttribute `json:"attributes" xml:"Attributes"`
	Product    Part               `json:"product" xml:"Product"`
}

// FitmentAttribute A name value for a note of a fitment application.
type FitmentAttribute struct {
	Key   string `json:"key" xml:"Key"`
	Value string `json:"value" xml:"Value"`
}

// Query Returns a `CategoryVehicle` that holds matching information for the
// queried `CategoryVehicleBase` attributes.
func Query(ctx *LookupContext, args ...string) (*CategoryVehicle, error) {

	var redisKey string
	var category string
	var vehicle CategoryVehicle

	for i, arg := range args {
		if i == 0 {
			redisKey = arg
		} else {
			redisKey = fmt.Sprintf("%s:%s", redisKey, arg)
		}
	}
	data, err := redis.Get(redisKey)
	log.Printf("REDIS_GET_ERROR: %+v\n", err)
	if err == nil && len(data) > 0 {
		err = json.Unmarshal(data, &vehicle)
		if err == nil {
			log.Printf("CATEGORY_STYLES :: %s\n", redisKey)
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
		category = args[3]
	}

	if vehicle.Base.Year == "" {
		vehicle.Years, err = getYears(ctx)
	} else if vehicle.Base.Year != "" && vehicle.Base.Make == "" {
		vehicle.Makes, err = getMakes(ctx, vehicle.Base.Year)
	} else if vehicle.Base.Year != "" && vehicle.Base.Make != "" && vehicle.Base.Model == "" {
		vehicle.Models, err = getModels(ctx, vehicle.Base.Year, vehicle.Base.Make)
	} else if vehicle.Base.Year != "" && vehicle.Base.Make != "" && vehicle.Base.Model != "" {
		vehicle.Categories, err = getStyles(ctx, vehicle.Base.Year, vehicle.Base.Make, vehicle.Base.Model, category)
	}

	log.Printf("SETTING_CATEGORY_STYLES :: %s\n", redisKey)
	redis.Setex(redisKey, vehicle, 60*60*24)

	return &vehicle, err
}

func getYears(ctx *LookupContext) ([]string, error) {
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
		"brand.id": bson.M{
			"$in": ctx.Brands,
		},
	}

	var res []string
	err := c.Find(qry).Select(bson.M{
		"vehicle_applications.year": 1,
		"_id": -1,
	}).Distinct("vehicle_applications.year", &res)

	if err != nil {
		return nil, err
	}

	sort.Sort(sort.Reverse(sort.StringSlice(res)))

	return res, err
}

func getMakes(ctx *LookupContext, year string) ([]string, error) {
	if ctx == nil {
		return nil, fmt.Errorf("missing context")
	} else if ctx.Session == nil {
		return nil, fmt.Errorf("invalid mongodb connection")
	}

	c := ctx.Session.DB(database.ProductMongoDatabase).C(database.ProductCollectionName)

	type Apps struct {
		Apps []VehicleApplication `bson:"vehicle_applications"`
	}

	var apps []Apps
	qry := bson.M{
		"vehicle_applications": bson.M{
			"$elemMatch": bson.M{
				"year": year,
			},
		},
		"status": bson.M{
			"$in": ctx.Statuses,
		},
		"brand.id": bson.M{
			"$in": ctx.Brands,
		},
	}
	err := c.Find(qry).Select(bson.M{"vehicle_applications.make": 1, "_id": 0}).All(&apps)
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

func getModels(ctx *LookupContext, year, vehicleMake string) ([]string, error) {
	if ctx == nil {
		return nil, fmt.Errorf("missing context")
	} else if ctx.Session == nil {
		return nil, fmt.Errorf("invalid mongodb connection")
	}

	c := ctx.Session.DB(database.ProductMongoDatabase).C(database.ProductCollectionName)

	type Apps struct {
		Apps []VehicleApplication `bson:"vehicle_applications"`
	}

	var apps []Apps
	err := c.Find(bson.M{
		"vehicle_applications": bson.M{
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
		"brand.id": bson.M{
			"$in": ctx.Brands,
		},
	}).Select(bson.M{"vehicle_applications.$.model": 1, "_id": 0}).All(&apps)
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

func getStyles(ctx *LookupContext, year, vehicleMake, model, category string) ([]LookupCategory, error) {
	if ctx == nil {
		return nil, fmt.Errorf("missing context")
	} else if ctx.Session == nil {
		return nil, fmt.Errorf("invalid mongodb connection")
	}

	c := ctx.Session.DB(database.ProductMongoDatabase).C(database.ProductCollectionName)

	type Apps struct {
		Apps    []VehicleApplication `bson:"vehicle_applications"`
		PartNum string               `bson:"part_number"`
	}

	var parts []Part
	qry := bson.M{
		"vehicle_applications": bson.M{
			"$elemMatch": bson.M{
				"year": year,
				"make": bson.RegEx{
					Pattern: "^" + vehicleMake + "$",
					Options: "i",
				},
				"model": bson.RegEx{
					Pattern: "^" + model + "$",
					Options: "i",
				},
			},
		},
		"status": bson.M{
			"$in": []int{800, 900},
		},
		"brand.id": bson.M{
			"$in": ctx.Brands,
		},
	}

	if category != "" {
		qry["categories.title"] = category
	}
	err := c.Find(qry).All(&parts)
	if err != nil || len(parts) == 0 {
		return nil, err
	}

	return generateCategoryStyles(parts, year, vehicleMake, model), nil
}

func generateCategoryStyles(parts []Part, year, vehicleMake, model string) []LookupCategory {
	lc := make(map[string]LookupCategory, 0)
	y := year
	ma := strings.ToLower(vehicleMake)
	mod := strings.ToLower(model)
	for _, p := range parts {
		if len(p.Categories) == 0 {
			continue
		}

		for _, va := range p.Vehicles {
			if va.Year != y || strings.ToLower(va.Make) != ma || strings.ToLower(va.Model) != mod {
				continue
			}

			lc = mapPartToCategoryStyles(p, lc, va.Style)
		}
	}

	var cats []LookupCategory
	for _, l := range lc {
		cats = append(cats, l)
	}

	return cats
}

// AddPart Creates a record of the provided part under the referenced style.
func (lc *LookupCategory) AddPart(style string, p Part) {
	if strings.TrimSpace(style) == "" {
		style = AllPlaceholder
	}

	for i, options := range lc.StyleOptions {
		if strings.TrimSpace(options.Style) == "" {
			options.Style = AllPlaceholder
			lc.StyleOptions[i].Style = AllPlaceholder
		}
		if strings.Compare(
			strings.ToLower(options.Style),
			strings.ToLower(style),
		) == 0 {
			lc.StyleOptions[i].Fitments = append(lc.StyleOptions[i].Fitments,
				FitmentProduct{
					Product: p,
				},
			)
			return
		}
	}

	lc.StyleOptions = append(lc.StyleOptions, StyleOption{
		Style: style,
		Fitments: []FitmentProduct{
			FitmentProduct{
				Product: p,
			},
		},
	})
}

func mapPartToCategoryStyles(p Part, lookupCats map[string]LookupCategory, style string) map[string]LookupCategory {
	childCat, err := getChildCategory(p.Categories)
	if err != nil || childCat.Identifier.String() == "" {
		return lookupCats
	}

	lc, ok := lookupCats[childCat.Identifier.String()]
	if !ok {
		childCat.PartIDs = nil
		childCat.Children = nil
		childCat.ProductListing = nil
		lc = LookupCategory{
			Category: childCat,
		}
	}

	// we're going to clear out the category information here, since
	// the products are already being grouped into their respective
	// categories at the higher level. (saves on da bits)
	p.Categories = nil

	// add the part to the appropriate style
	lc.AddPart(style, p)

	// currentStyle := StyleOption{
	// 	Style: strings.Title(style),
	// 	Fitments: []FitmentProduct{
	// 		FitmentProduct{
	// 			Product: p,
	// 		},
	// 	},
	// }
	//
	// lc.StyleOptions = append(lc.StyleOptions, currentStyle)
	lookupCats[childCat.Identifier.String()] = lc

	return lookupCats
}

func getChildCategory(cats []Category) (Category, error) {
	for _, cat := range cats {
		if len(cat.Children) == 0 {
			return cat, nil
		}
	}

	return Category{}, fmt.Errorf("failed to locate child")
}