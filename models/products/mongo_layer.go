package products

import (
	"github.com/curt-labs/GoAPI/helpers/database"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func GetCategoryTree() ([]Category, error) {
	var cats []Category

	session, err := mgo.DialWithInfo(database.MongoCategoryConnectionString())
	if err != nil {
		return cats, err
	}
	defer session.Close()
	query := bson.M{"parent_id": 0}
	err = session.DB(database.CategoryDatabase).C(database.CategoryCollectionName).Find(query).All(&cats)
	return cats, err
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
