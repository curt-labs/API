package mongoData

import (
	"github.com/curt-labs/GoAPI/helpers/apicontext"
	"github.com/curt-labs/GoAPI/helpers/database"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func GetCategoryTree(dtx *apicontext.DataContext) ([]Category, error) {
	var cats []Category

	session, err := mgo.DialWithInfo(database.MongoPartConnectionString())
	if err != nil {
		return cats, err
	}
	defer session.Close()
	query := bson.M{"parent_id": 0, "is_lifestyle": false, "brand.id": bson.M{"$in": dtx.BrandArray}}
	err = session.DB(database.ProductDatabase).C(database.CategoryCollectionName).Find(query).Sort("title").All(&cats)
	return cats, err
}

func (c *Category) FromMongo(page, count int) error {

	session, err := mgo.DialWithInfo(database.MongoPartConnectionString())
	if err != nil {
		return err
	}
	defer session.Close()

	err = session.DB(database.ProductDatabase).C(database.CategoryCollectionName).Find(bson.M{"id": c.CategoryID}).One(&c)
	if err != nil {
		return err
	}

	if page < 1 {
		page = 1
	}
	if count < 0 {
		count = 1
	} else if count > 50 {
		count = 50
	}

	var skip int
	if page > 1 {
		skip = page * count
	}

	c.ProductListing = &PaginatedProductListing{
		Page:    page,
		PerPage: count,
		Parts:   []Product{},
	}

	c.ProductListing.TotalItems, err = session.DB(database.ProductDatabase).C(database.ProductCollectionName).Find(bson.M{"id": bson.M{"$in": c.ProductIdentifiers}}).Count()
	if err != nil {
		c.ProductListing.TotalItems = 1
	}

	err = session.DB(database.ProductDatabase).C(database.ProductCollectionName).Find(bson.M{"id": bson.M{"$in": c.ProductIdentifiers}}).Sort("id").Skip(skip).Limit(count).All(&c.ProductListing.Parts)
	if err != nil {
		return err
	}

	c.ProductListing.ReturnedCount = len(c.ProductListing.Parts)
	c.ProductListing.TotalPages = c.ProductListing.TotalItems / c.ProductListing.PerPage

	return nil
}

func GetCategoryParts(catId int) ([]Product, error) {
	var parts []Product

	session, err := mgo.DialWithInfo(database.MongoPartConnectionString())
	if err != nil {
		return parts, err
	}
	defer session.Close()
	query := bson.M{"categories": bson.M{"$elemMatch": bson.M{"id": catId}}}
	err = session.DB(database.ProductDatabase).C(database.ProductCollectionName).Find(query).All(&parts)
	return parts, err
}
