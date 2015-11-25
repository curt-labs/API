package products

import (
	"github.com/curt-labs/API/helpers/database"
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

	err = session.DB(database.AriesMongoConnectionString().Database).C(collection).Find(bson.M{}).Sort("-year", "make", "model", "style").All(&apps)
	return apps, err
}

func (n *NoSqlVehicle) Update(collection string) error {
	session, err := mgo.DialWithInfo(database.AriesMongoConnectionString())
	if err != nil {
		return err
	}
	defer session.Close()

	update := make(map[string]interface{})
	if n.Year != "" {
		update["year"] = n.Year
	}
	if n.Make != "" {
		update["make"] = n.Make
	}
	if n.Model != "" {
		update["model"] = n.Model
	}
	if n.Style != "" {
		update["style"] = n.Style
	}
	if n.Make != "" {
		update["make"] = n.Make
	}
	if len(n.PartIdentifiers) > 0 {
		update["parts"] = n.PartIdentifiers
	}
	return session.DB(database.AriesMongoConnectionString().Database).C(collection).UpdateId(n.ID, update)
}

func (n *NoSqlVehicle) Delete(collection string) error {
	session, err := mgo.DialWithInfo(database.AriesMongoConnectionString())
	if err != nil {
		return err
	}
	defer session.Close()
	return session.DB(database.AriesMongoConnectionString().Database).C(collection).RemoveId(n.ID)
}

func (n *NoSqlVehicle) Create(collection string) error {
	session, err := mgo.DialWithInfo(database.AriesMongoConnectionString())
	if err != nil {
		return err
	}
	defer session.Close()
	return session.DB(database.AriesMongoConnectionString().Database).C(collection).Insert(n)
}
