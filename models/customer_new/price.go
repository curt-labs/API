package customer_new

import (
	"database/sql"
	"encoding/json"
	"github.com/curt-labs/GoAPI/helpers/database"
	"time"
	// "github.com/curt-labs/goacesapi/helpers/pagination"
	"github.com/curt-labs/GoAPI/helpers/redis"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
)

type Price struct {
	ID        int
	CustID    int
	PartID    int
	Price     float64
	IsSale    int
	SaleStart time.Time
	SaleEnd   time.Time
}

type Prices []Price

type CustomerPrices struct {
	Customer Customer `json:"customer" xml:"customer"`
	Prices   Prices   `json:"prices" xml:"prices"`
}

var (
	getPrice             = "SELECT cust_price_id, cust_id, partID, price, isSale, sale_start, sale_end FROM CustomerPricing WHERE cust_price_id = ?"
	getPrices            = "SELECT cust_price_id, cust_id, partID, price, isSale, sale_start, sale_end FROM CustomerPricing"
	createPrice          = "INSERT INTO CustomerPricing (cust_id, partID, price, isSale, sale_start, sale_end) VALUES (?,?,?,?,?,?)"
	updatePrice          = "UPDATE CustomerPricing SET cust_id = ?, partID = ?, price = ?, isSale = ?, sale_start = ?, sale_end = ? WHERE cust_price_id = ?"
	deletePrice          = "DELETE FROM CustomerPricing WHERE cust_price_id = ?"
	getPricesByCustomer  = "SELECT cust_price_id, cust_id, partID, price, isSale, sale_start, sale_end FROM CustomerPricing WHERE cust_id = (select cust_id from Customer where customerID = ?)"
	getPricesByPart      = "SELECT cust_price_id, cust_id, partID, price, isSale, sale_start, sale_end FROM CustomerPricing WHERE partID = ?"
	getPricesBySaleRange = "SELECT cust_price_id, cust_id, partID, price, isSale, sale_start, sale_end FROM CustomerPricing WHERE sale_start >= ? AND sale_end <= ? AND cust_id = (select cust_id from Customer where customerID = ?)"
)

const (
	timeFormat = "2006-01-02"
)

func (p *Price) Get() error {
	redis_key := "goapi:price:" + strconv.Itoa(p.ID)
	data, err := redis.Get(redis_key)
	if err == nil && len(data) > 0 {
		err = json.Unmarshal(data, &p)
		return err
	}

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(getPrice)
	if err != nil {
		return err
	}
	defer stmt.Close()
	err = stmt.QueryRow(p.ID).Scan(&p.ID, &p.CustID, &p.PartID, &p.Price, &p.IsSale, &p.SaleStart, &p.SaleEnd)
	if err != nil {
		return err
	}

	go redis.Setex(redis_key, p, 86400)
	return nil
}

func GetAllPrices() (Prices, error) {
	var ps Prices
	var err error
	redis_key := "goapi:prices"
	data, err := redis.Get(redis_key)
	if err == nil && len(data) > 0 {
		err = json.Unmarshal(data, &ps)
		return ps, err
	}

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return ps, err
	}
	defer db.Close()

	stmt, err := db.Prepare(getPrices)
	if err != nil {
		return ps, err
	}
	defer stmt.Close()
	res, err := stmt.Query()
	for res.Next() {
		var p Price
		err = res.Scan(&p.ID, &p.CustID, &p.PartID, &p.Price, &p.IsSale, &p.SaleStart, &p.SaleEnd)
		ps = append(ps, p)
	}

	go redis.Setex(redis_key, ps, 86400)
	return ps, nil
}

func (p *Price) Create() error {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(createPrice)
	res, err := stmt.Exec(p.CustID, p.PartID, p.Price, p.IsSale, p.SaleStart, p.SaleEnd)
	if err != nil {
		tx.Rollback()
		return err
	}
	id, err := res.LastInsertId()
	p.ID = int(id)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	redis.Setex("goapi:prices:"+strconv.Itoa(p.ID), p, 86400)
	return nil
}
func (p *Price) Update() error {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(updatePrice)

	_, err = stmt.Exec(p.CustID, p.PartID, p.Price, p.IsSale, p.SaleStart, p.SaleEnd, p.ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	err = redis.Setex("goapi:prices:"+strconv.Itoa(p.ID), p, 86400)
	return nil
}

func (p *Price) Delete() error {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(deletePrice)
	_, err = stmt.Exec(p.ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	err = redis.Delete("goapi:prices:" + strconv.Itoa(p.ID))
	return nil
}

func (c *Customer) GetPricesByCustomer() (CustomerPrices, error) {
	var cps CustomerPrices
	redis_key := "goapi:customers:prices:" + strconv.Itoa(c.Id)
	data, err := redis.Get(redis_key)
	if err == nil && len(data) > 0 {
		err = json.Unmarshal(data, &cps)
		return cps, err
	}

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return cps, err
	}
	defer db.Close()

	stmt, err := db.Prepare(getPricesByCustomer)
	if err != nil {
		return cps, err
	}
	defer stmt.Close()
	res, err := stmt.Query(c.Id)
	for res.Next() {
		var p Price
		res.Scan(&p.ID, &p.CustID, &p.PartID, &p.Price, &p.IsSale, &p.SaleStart, &p.SaleEnd)

		cps.Prices = append(cps.Prices, p)
	}
	go redis.Setex(redis_key, cps, 86400)
	return cps, err
}

func GetPricesByPart(partID int) (Prices, error) {
	var ps Prices
	redis_key := "goapi:prices:part:" + strconv.Itoa(partID)
	data, err := redis.Get(redis_key)
	if err == nil && len(data) > 0 {
		err = json.Unmarshal(data, &ps)
		return ps, err
	}

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return ps, err
	}
	defer db.Close()

	stmt, err := db.Prepare(getPricesByPart)
	if err != nil {
		return ps, err
	}
	defer stmt.Close()
	res, err := stmt.Query(partID)
	for res.Next() {
		var p Price
		res.Scan(&p.ID, &p.CustID, &p.PartID, &p.Price, &p.IsSale, &p.SaleStart, &p.SaleEnd)
		ps = append(ps, p)
	}
	go redis.Setex(redis_key, ps, 86400)
	return ps, nil
}

func (c *Customer) GetPricesBySaleRange(startDate, endDate time.Time) (Prices, error) {
	var err error
	var ps Prices
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return ps, err
	}

	defer db.Close()

	stmt, err := db.Prepare(getPricesBySaleRange)
	if err != nil {
		return ps, err
	}
	defer stmt.Close()
	res, err := stmt.Query(startDate, endDate, c.Id)
	for res.Next() {
		var p Price
		res.Scan(&p.ID, &p.CustID, &p.PartID, &p.Price, &p.IsSale, &p.SaleStart, &p.SaleEnd)
		ps = append(ps, p)
	}

	return ps, nil
}
