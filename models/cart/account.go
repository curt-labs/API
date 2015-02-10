package cart

import (
	"fmt"
	"github.com/curt-labs/GoAPI/helpers/database"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func AuthenticateAccount(token string) (Customer, error) {

	sess, err := mgo.DialWithInfo(database.MongoConnectionString())
	if err != nil {
		return Customer{}, err
	}
	defer sess.Close()

	var cust Customer
	qs := bson.M{
		"token": token,
		"$or": []bson.M{
			bson.M{"customer.addresses.deleted": false},
			bson.M{"customer.addresses.deleted": bson.M{
				"$exists": false,
			}},
		},
	}
	err = sess.DB("CurtCart").C("customer").Find(qs).One(&cust)
	if err != nil || !cust.Id.Valid() {
		return Customer{}, fmt.Errorf("error: %s", "failed to authenticate using JWT")
	}

	return cust, nil
}

func IdentifierFromToken(t string) (bson.ObjectId, error) {
	sess, err := mgo.DialWithInfo(database.MongoConnectionString())
	if err != nil {
		return "", err
	}
	defer sess.Close()

	var cust Customer
	qs := bson.M{
		"token": t,
		"$or": []bson.M{
			bson.M{"customer.addresses.deleted": false},
			bson.M{"customer.addresses.deleted": bson.M{
				"$exists": false,
			}},
		},
	}
	err = sess.DB("CurtCart").C("customer").Find(qs).One(&cust)
	if err != nil || !cust.Id.Valid() {
		return "", fmt.Errorf("error: %s", "failed to identify using JWT")
	}

	return cust.Id, nil
}
