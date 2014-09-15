package products

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/curt-labs/GoAPI/helpers/database"
	"github.com/curt-labs/GoAPI/helpers/redis"
	_ "github.com/go-sql-driver/mysql"
)

type Price struct {
	Type     string
	Price    float64
	Enforced bool
}

var (
	partPriceStmt = `
		select priceType, price, enforced from Price
		where partID = ?
		order by priceType`
)

func (p *Part) GetPricing() error {
	redis_key := fmt.Sprintf("part:%d:pricing", p.PartId)

	data, err := redis.Get(redis_key)
	if err == nil && len(data) > 0 {
		if err = json.Unmarshal(data, &p.Pricing); err != nil {
			return nil
		}
	}

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	qry, err := db.Prepare(partPriceStmt)
	if err != nil {
		return err
	}
	defer qry.Close()

	rows, err := qry.Query(p.PartId)
	if err != nil || rows == nil {
		return err
	}

	for rows.Next() {
		var pr Price
		err = rows.Scan(&pr.Type, &pr.Price, &pr.Enforced)
		if err == nil {
			p.Pricing = append(p.Pricing, pr)
		}
	}

	go redis.Setex(redis_key, p.Pricing, 86400)

	return nil
}
