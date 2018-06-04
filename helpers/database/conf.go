package database

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"crypto/tls"
	"database/sql"

	"net"

	"github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/mysql"
	"github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/proxy"
	_ "github.com/go-sql-driver/mysql"
	goauth "golang.org/x/oauth2/google"
	"gopkg.in/mgo.v2"
)

type Scanner interface {
	Scan(...interface{}) error
}

var (
	EmptyDb = flag.String("clean", "", "bind empty database with structure defined")

	ProductCollectionName  = "products"
	CategoryCollectionName = "categories"
	CustomerCollectionName = "customer"
	ProductDatabase        = "product_data"
	CategoryDatabase       = "category_data"
	AriesDatabase          = "aries"

	MongoDatabase        string
	ProductMongoDatabase string
	AriesMongoDatabase   string

	MongoSession         *mgo.Session
	ProductMongoSession  *mgo.Session
	CategoryMongoSession *mgo.Session
	AriesMongoSession    *mgo.Session

	DB     *sql.DB
	VcdbDB *sql.DB
	Driver = "mysql"
)

func Init() error {
	var err error
	if DB == nil {
		if os.Getenv("DATABASE_INSTANCE") == "" {
			DB, err = sql.Open(Driver, ConnectionString())
		} else {
			client, err := clientFromCredentials()
			if err != nil {
				return err
			}

			proxy.Init(client, nil, nil)

			cfg := mysql.Cfg(os.Getenv("DATABASE_INSTANCE"), os.Getenv("DATABASE_USERNAME"), os.Getenv("DATABASE_PASSWORD"))
			cfg.DBName = os.Getenv("CURT_DEV_NAME")
			cfg.ParseTime = true
			cfg.AllowNativePasswords = true
			DB, err = mysql.DialCfg(cfg)
		}
		if err != nil {
			return err
		}
	}

	if VcdbDB == nil {
		if os.Getenv("DATABASE_INSTANCE") == "" {
			VcdbDB, err = sql.Open(Driver, VcdbConnectionString())
		} else {
			client, err := clientFromCredentials()
			if err != nil {
				return err
			}

			proxy.Init(client, nil, nil)

			cfg := mysql.Cfg(os.Getenv("DATABASE_INSTANCE"), os.Getenv("DATABASE_USERNAME"), os.Getenv("DATABASE_PASSWORD"))
			cfg.DBName = os.Getenv("VCDB_NAME")
			cfg.ParseTime = true
			cfg.AllowNativePasswords = true
			VcdbDB, err = mysql.DialCfg(cfg)
		}

		if err != nil {
			return err
		}
	}

	return InitMongo()
}

func ConnectionString() string {
	if addr := os.Getenv("DATABASE_HOST"); addr != "" {
		proto := os.Getenv("DATABASE_PROTOCOL")
		user := os.Getenv("DATABASE_USERNAME")
		pass := os.Getenv("DATABASE_PASSWORD")
		db := os.Getenv("CURT_DEV_NAME")

		return fmt.Sprintf("%s:%s@%s(%s)/%s?parseTime=true&loc=%s", user, pass, proto, addr, db, "America%2FChicago")
	}

	if EmptyDb != nil && *EmptyDb != "" {
		return "root:@tcp(127.0.0.1:3306)/CurtDev_Empty?parseTime=true&loc=America%2FChicago"
	}
	return "root:@tcp(127.0.0.1:3306)/CurtData?parseTime=true&loc=America%2FChicago"
}

func VcdbConnectionString() string {
	if addr := os.Getenv("DATABASE_HOST"); addr != "" {
		proto := os.Getenv("DATABASE_PROTOCOL")
		user := os.Getenv("DATABASE_USERNAME")
		pass := os.Getenv("DATABASE_PASSWORD")
		db := os.Getenv("VCDB_NAME")
		return fmt.Sprintf("%s:%s@%s(%s)/%s?parseTime=true&loc=%s", user, pass, proto, addr, db, "America%2FChicago")
	}

	return "root:@tcp(127.0.0.1:3306)/vcdb?parseTime=true&loc=America%2FChicago"
}

func VintelligencePass() string {
	if vinPin := os.Getenv("VIN_PIN"); vinPin != "" {
		return fmt.Sprintf("%s", vinPin)
	}
	return "curtman:Oct2013!"
}

func MongoConnectionString() *mgo.DialInfo {
	info := mgo.DialInfo{
		Addrs:    []string{"127.0.0.1:27017"},
		Database: "CurtCart",
		Source:   "admin",
		Timeout:  time.Second * 10,
		FailFast: true,
	}

	if addr := os.Getenv("MONGO_URL"); addr != "" {
		info.Addrs = strings.Split(addr, ",")
		info.Username = os.Getenv("MONGO_CART_USERNAME")
		info.Password = os.Getenv("MONGO_CART_PASSWORD")
		info.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
			return tls.Dial("tcp", addr.String(), &tls.Config{
				InsecureSkipVerify: true,
			})
		}
		db := os.Getenv("MONGO_CART_DATABASE")
		if db != "" {
			info.Database = db
		}
		source := os.Getenv("MONGO_AUTH_DATABASE")
		if source != "" {
			info.Source = source
		}
	}

	return &info
}

func MongoPartConnectionString() *mgo.DialInfo {
	info := MongoConnectionString()
	info.Database = ProductDatabase
	return info
}

func AriesMongoConnectionString() *mgo.DialInfo {
	info := mgo.DialInfo{
		Addrs:    []string{"127.0.0.1:27017"},
		Database: "aries",
		Source:   "admin",
		Timeout:  time.Second * 10,
		FailFast: true,
	}

	if addr := os.Getenv("MONGO_URL"); addr != "" {
		info.Addrs = strings.Split(addr, ",")
		info.Username = os.Getenv("MONGO_ARIES_USERNAME")
		info.Password = os.Getenv("MONGO_ARIES_PASSWORD")
		info.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
			return tls.Dial("tcp", addr.String(), &tls.Config{
				InsecureSkipVerify: true,
			})
		}
		db := os.Getenv("MONGO_ARIES_DATABASE")
		if db != "" {
			info.Database = db
		}
		source := os.Getenv("MONGO_AUTH_DATABASE")
		if source != "" {
			info.Source = source
		}
	}

	return &info
}

func AriesConnectionString() string {
	return ConnectionString()
}

func GetCleanDBFlag() string {
	return *EmptyDb
}

func InitMongo() error {
	var err error
	if MongoSession == nil {
		connectionString := MongoConnectionString()
		MongoSession, err = mgo.DialWithInfo(connectionString)
		if err != nil {
			return err
		}
		MongoDatabase = connectionString.Database
	}
	if ProductMongoSession == nil {
		connectionString := MongoPartConnectionString()
		ProductMongoSession, err = mgo.DialWithInfo(connectionString)
		if err != nil {
			return err
		}
		ProductMongoDatabase = connectionString.Database
	}
	if CategoryMongoSession == nil {
		connectionString := MongoPartConnectionString()
		CategoryMongoSession, err = mgo.DialWithInfo(connectionString)
		if err != nil {
			return err
		}
		ProductMongoDatabase = connectionString.Database
	}
	if AriesMongoSession == nil {
		connectionString := AriesMongoConnectionString()
		AriesMongoSession, err = mgo.DialWithInfo(connectionString)
		if err == nil {
			AriesMongoDatabase = connectionString.Database
		}
	}
	return err
}

func clientFromCredentials() (*http.Client, error) {

	const SQLScope = "https://www.googleapis.com/auth/sqlservice.admin"

	ctx := context.Background()

	var client *http.Client

	cfg, err := goauth.JWTConfigFromJSON([]byte(os.Getenv("DATABASE_TOKEN")), SQLScope)
	if err != nil {
		return nil, fmt.Errorf("invalid json file: %v", err)
	}

	client = cfg.Client(ctx)

	return client, nil

}
