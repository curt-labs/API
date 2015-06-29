package products

// import (
// 	"encoding/csv"
// 	"fmt"
// 	"log"
// 	"mime/multipart"
// 	"strings"

// 	"github.com/curt-labs/GoAPI/helpers/database"
// 	"gopkg.in/mgo.v2"
// )

// func GetCollections() ([]string, error) {
// 	session, err := mgo.DialWithInfo(database.AriesMongoConnectionString())
// 	if err != nil {
// 		return
// 	}
// 	defer session.Close()

// 	return session.DB(database.AriesMongoConnectionString().Database).CollectionNames()

// }
