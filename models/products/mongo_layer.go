package products

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

func (c *Category) FromMongo() error {

	session, err := mgo.DialWithInfo(database.MongoPartConnectionString())
	if err != nil {
		return err
	}
	defer session.Close()

	return session.DB(database.ProductDatabase).C(database.CategoryCollectionName).Find(bson.M{"id": c.ID}).One(&c)
}

func GetCategoryParts(catId int) ([]Part, error) {
	var parts []Part

	session, err := mgo.DialWithInfo(database.MongoPartConnectionString())
	if err != nil {
		return parts, err
	}
	defer session.Close()
	query := bson.M{"categories": bson.M{"$elemMatch": bson.M{"id": catId}}}
	err = session.DB(database.ProductDatabase).C(database.ProductCollectionName).Find(query).All(&parts)
	return parts, err
}
