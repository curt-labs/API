package part

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/curt-labs/GoAPI/helpers/database"
	"github.com/curt-labs/GoAPI/helpers/redis"
	_ "github.com/go-sql-driver/mysql"
)

type Attribute struct {
	Key, Value string
}

var (
	partAttrStmt = `select field, value from PartAttribute where partID = ?`
)

func (p *Part) GetAttributes() (err error) {
	redis_key := fmt.Sprintf("part:%d:attributes", p.PartId)

	data, err := redis.Get(redis_key)
	if err == nil && len(data) > 0 {
		if err = json.Unmarshal(data, &p.Attributes); err != nil {
			return nil
		}
	}

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return
	}
	defer db.Close()

	qry, err := db.Prepare(partAttrStmt)
	if err != nil {
		return
	}
	defer qry.Close()

	rows, err := qry.Query(p.PartId)
	if err != nil || rows == nil {
		return err
	}

	var attrs []Attribute
	for rows.Next() {
		var attr Attribute
		if err := rows.Scan(&attr.Key, &attr.Value); err == nil {
			attrs = append(attrs, attr)
		}
	}
	p.Attributes = attrs

	go redis.Setex(redis_key, p.Attributes, redis.CacheTimeout)

	return
}