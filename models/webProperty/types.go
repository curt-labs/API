package webProperty_model

import (
	"database/sql"
	"encoding/json"

	"github.com/curt-labs/API/helpers/apicontext"
	"github.com/curt-labs/API/helpers/database"
	"github.com/curt-labs/API/helpers/redis"
	_ "github.com/go-sql-driver/mysql"
)

// The Type of a WebProperty. Examples are: Website, Ebay Store, Amazon Store
type WebPropertyType2 struct {
	ID     int    `json:"id,omitempty" xml:"id,omitempty"`
	TypeID int    `json:"typeId,omitempty" xml:"typeId,omitempty"`
	Type   string `json:"type,omitempty" xml:"type,omitempty"`
}

// WebPropertiesTypes is just an easier type to work with than using an array of WebPropertyType's.
type WebPropertyTypes2 []WebPropertyType

var getAllWebPropertyTypes2     = `SELECT DISTINCT wt.id, wt.typeID, wt.type
		FROM WebPropertyTypes as wt
		join WebProperties as w on w.typeID = wt.id
		join CustomerToBrand as ctb on ctb.cust_id = w.cust_id
		join ApiKeyToBrand as atb on atb.brandID = ctb.brandID
		join ApiKey as a on a.id = atb.keyID
		where a.api_key = ? && (ctb.brandID = ? or 0 = ?)`

func GetAllWebPropertyTypesByBrand(apiKey string, brandId int) (ws WebPropertyTypes, err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return ws, err
	}
	defer db.Close()

	stmt, err := db.Prepare(getAllWebPropertyTypes2)
	if err != nil {
		return ws, err
	}
	defer stmt.Close()

	res, err := stmt.Query(apiKey, brandId, brandId)
	for res.Next() {
		var w WebPropertyType
		res.Scan(&w.ID, &w.TypeID, &w.Type)
		ws = append(ws, w)
	}

	return ws, err
}

// Gets All the available WebPropertyTypes
func GetAllWebPropertyTypes2(dtx *apicontext.DataContext) (ws WebPropertyTypes, err error) {
	redis_key := "webpropertytypes:" + dtx.BrandString
	data, err := redis.Get(redis_key)
	if err == nil && len(data) > 0 {
		err = json.Unmarshal(data, &ws)
		if err == nil {
			return ws, err
		}
	}

	ws, err = GetAllWebPropertyTypesByBrand(dtx.APIKey, dtx.BrandID)

	go redis.Setex(redis_key, ws, 86400) // Expire in 1 Day

	return ws, err
}
