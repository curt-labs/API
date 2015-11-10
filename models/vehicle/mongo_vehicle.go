package vehicle

import (
	"github.com/curt-labs/GoAPI/helpers/database"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type MgoVehicle struct {
	Identifier bson.ObjectId `bson:"_id" json:"-" xml:"-"`
	Year       string        `bson:"year" json:"year,omitempty" xml:"year,omitempty"`
	Make       string        `bson:"make" json:"make,omitempty" xml:"make,omitempty"`
	Model      string        `bson:"model" json:"model,omitempty" xml:"model,omitempty"`
	Style      string        `bson:"style" json:"style,omitempty" xml:"style,omitempty"`
}

const (
	AriesDb = "aries"
)

func ReverseMongoLookup(partId int) (vehicles []MgoVehicle, err error) {
	session, err := mgo.DialWithInfo(database.AriesMongoConnectionString())
	if err != nil {
		return
	}
	defer session.Close()

	collections, err := session.DB(AriesDb).CollectionNames()
	if err != nil {
		return
	}
	for _, collection := range collections {
		var temps []MgoVehicle
		query := bson.M{
			"parts": partId,
		}
		err = session.DB(AriesDb).C(collection).Find(query).All(&temps)
		if err != nil {
			return
		}
		vehicles = append(vehicles, temps...)
	}
	return
}
