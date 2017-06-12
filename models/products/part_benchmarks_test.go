package products

import (
	"testing"
	"github.com/curt-labs/API/mocks"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"os"
	"net/http"
)

// You should already have an existing MongoDB server running with the test data (dump of production databases)
// After the test database server has been set up, set the DOCKER_MACHINE_NAME
//	# (Linux) Setting up the benchmarking environment
// 	export DOCKER_MACHINE_NAME=goapi_test_data
//	docker run --name $DOCKER_MACHINE_NAME mongo:3.2
//	# Backup production database
//	export MONGO_BACKUP=mongo_backup_$(date +%Y%m%d)
//	mongodump --ssl \
// 	          --host=$PRODUCTION_HOST \
//	          --username=$PRODUCTION_USER
// 	          --password=$PRODUCTION_PASS
// 	          --authenticationDatabase=$PRODUCTION_AUTH_DB
// 	          --out=$MONGO_BACKUP
//	mongorestore --host=localhost:27017 --authenticationDatabase=$PRODCUTION_AUTH_DB --dir=$MONGO_BACKUP


func BenchmarkGetMany(b *testing.B) {
	mongo, err := mocks.NewDockertestMongo()
	if err != nil {
		b.Fatal(err)
	}
	defer mongo.Close()

	http.status

	// verify localhost
	c := mongo.Session.DB("product_data").C("products")
	qry := bson.M{"id": 121151}
	mongo.Pool.

	var parts []Part
	err = c.Find(qry).All(&parts)
	fmt.Printf("env = %s\n", os.Getenv("DOCKER_MACHINE_NAME"))
	fmt.Println(mongo.Pool.Client)
	fmt.Println(err)
	fmt.Println(parts)
	// end localhost

	b.Run("Testing", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			GetMany([]int{}, []int{}, mongo.Session)
		}
	})
}