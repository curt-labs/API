package products

import (
	"github.com/curt-labs/GoAPI/helpers/database"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func GetAllCollectionApplications(collection string) ([]NoSqlVehicle, error) {
	var apps []NoSqlVehicle
	session, err := mgo.DialWithInfo(database.AriesMongoConnectionString())
	if err != nil {
		return apps, err
	}
	defer session.Close()

	err = session.DB(database.AriesMongoConnectionString().Database).C(collection).Find(bson.M{}).All(&apps)
	return apps, err
}
