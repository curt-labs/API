package cartIntegration

import (
	"github.com/curt-labs/GoAPI/helpers/apicontext"
	"github.com/curt-labs/GoAPI/helpers/database"
	_ "github.com/go-sql-driver/mysql"

	"database/sql"
	"strconv"
	"time"
)

type CustomerPrice struct {
	ID             int        `json:"id,omitempty" xml:"id,omitempty"`
	CustID         int        `json:"custId,omitempty" xml:"custId,omitempty"`
	PartID         int        `json:"partId,omitempty" xml:"partId,omitempty"`
	CustomerPartID int        `json:"customerPartId,omitempty" xml:"customerPartId,omitempty"`
	Price          float64    `json:"price,omitempty" xml:"price,omitempty"`
	IsSale         int        `json:"isSale,omitempty" xml:"isSale,omitempty"`
	SaleStart      *time.Time `json:"saleStart,omitempty" xml:"saleStart,omitempty"`
	SaleEnd        *time.Time `json:"saleEnd,omitempty" xml:"saleEnd,omitempty"`
}

type Price struct {
	PartID int     `json:"partId,omitempty" xml:"partId,omitempty"`
	Type   string  `json:"type,omitempty" xml:"type,omitempty"`
	Price  float64 `json:"price,omitempty" xml:"price,omitempty"`
}

var (
	getPricing = `SELECT cp.cust_price_id, cp.cust_id, cp.partID, cp.price, cp.isSale, cp.sale_start, cp.sale_end FROM CustomerPricing as cp
		JOIN CustomerUser as cu on cu.cust_id = cp.cust_id
		JOIN ApiKey as a on a.user_id = cu.id
		JOIN Part as p on p.partID = cp.partID
		WHERE p.brandID = ?
		AND a.api_key = ?`

	getPricingPaged = `SELECT cp.cust_price_id, cp.cust_id, cp.partID, cp.price, cp.isSale, cp.sale_start, cp.sale_end FROM CustomerPricing as cp
		JOIN CustomerUser as cu on cu.cust_id = cp.cust_id
		JOIN ApiKey as a on a.user_id = cu.id
		JOIN Part as p on p.partID = cp.partID
		WHERE p.brandID = ?
		AND a.api_key = ?
		order by cp.partID
		LIMIT ?,?`

	getPricingCount = `SELECT COUNT(cp.cust_price_id) FROM CustomerPricing as cp
		JOIN CustomerUser as cu on cu.cust_id = cp.cust_id
		JOIN ApiKey as a on a.user_id = cu.id
		JOIN Part as p on p.partID = cp.partID
		WHERE p.brandID = ?
		AND a.api_key = ?`

	getPricingByPart = `SELECT pr.partID, pr.priceType, pr.price FROM Price as pr
		JOIN Part as p ON pr.partID = p.partID
		WHERE p.status != 999 && p.brandID = ? && p.partID = ?
		ORDER BY pr.priceType`
	getAllPricing = `SELECT pr.partID, pr.priceType, pr.price FROM Price as pr
		JOIN Part as p ON pr.partID = p.partID
		WHERE p.status != 999 && p.brandID = ?
		ORDER by p.partID, pr.priceType`
	getAllMAPPricing = `SELECT pr.partID, pr.priceType, pr.price FROM Price as pr
		JOIN Part as p ON pr.partID = p.partID
		WHERE p.status != 999 && p.brandID = ? && pr.priceType = 'Map'
		ORDER by p.partID, pr.priceType`

	updateCustomerPrice = `UPDATE CustomerPricing SET price = ?, isSale = ?, sale_start = ?, sale_end = ? WHERE cust_id = ? AND partID = ?`
	insertCustomerPrice = `INSERT INTO CustomerPricing(cust_id, partID, price, isSale, sale_start, sale_end) VALUES(?, ?, ?, ?, ?, ?)`
	deleteCustomerPrice = `delete from CustomerPricing where cust_price_id = ?`

	getCustomerCartIntegrations = `select c.referenceID, c.partID, c.custPartID, c.custID from CartIntegration as c
		join CustomerUser as cu on cu.cust_id = c.custID
		join ApiKey as a on a.user_id = cu.id
		join Part as p on p.partID = c.partID
		where a.api_key = ?
		and p.brandID = ?
		order by p.partID`
	insertCartIntegration = `INSERT INTO CartIntegration(partID, custPartID, custID) VALUES (?, ?, ?)`
	updateCartIntegration = `UPDATE CartIntegration SET custPartID = ? WHERE partID = ? AND custID = ?`
	getAllPriceTypes      = `SELECT DISTINCT priceType from Price`
)

//Get all of a single customer's prices
func GetCustomerPrices(dtx *apicontext.DataContext) ([]CustomerPrice, error) {
	var cps []CustomerPrice
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return cps, err
	}
	defer db.Close()
	stmt, err := db.Prepare(getPricing)
	if err != nil {
		return cps, err
	}
	defer stmt.Close()
	res, err := stmt.Query(dtx.BrandID, dtx.APIKey)
	if err != nil {
		return cps, err
	}

	//customer cart integration
	cartIntegrations, err := GetCustomerCartIntegrations(dtx)
	if err != nil {
		return cps, err
	}
	cartIntegrationMap := make(map[string]int) //partID+:+custID -to- custPartID
	for _, ci := range cartIntegrations {
		cartIntegrationMap[strconv.Itoa(ci.PartID)+":"+strconv.Itoa(ci.CustID)] = ci.CustomerPartID
	}

	for res.Next() {
		c, err := Scan(res)
		if err != nil {
			return cps, err
		}
		c.CustomerPartID = cartIntegrationMap[strconv.Itoa(c.PartID)+":"+strconv.Itoa(c.CustID)]
		cps = append(cps, c)
	}
	return cps, err
}

//Get a customers prices - paged/limited
func GetPricingPaged(page int, count int, dtx *apicontext.DataContext) ([]CustomerPrice, error) {
	var cps []CustomerPrice
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return cps, err
	}
	defer db.Close()
	stmt, err := db.Prepare(getPricingPaged)
	if err != nil {
		return cps, err
	}
	defer stmt.Close()
	res, err := stmt.Query(dtx.BrandID, dtx.APIKey, (page-1)*count, count)
	if err != nil {
		return cps, err
	}

	//customer cart integration
	cartIntegrations, err := GetCustomerCartIntegrations(dtx)
	if err != nil {
		return cps, err
	}
	cartIntegrationMap := make(map[string]int) //partID+:+custID -to- custPartID
	for _, ci := range cartIntegrations {
		cartIntegrationMap[strconv.Itoa(ci.PartID)+":"+strconv.Itoa(ci.CustID)] = ci.CustomerPartID
	}

	for res.Next() {
		c, err := Scan(res)
		if err != nil {
			return cps, err
		}
		c.CustomerPartID = cartIntegrationMap[strconv.Itoa(c.PartID)+":"+strconv.Itoa(c.CustID)]
		cps = append(cps, c)
	}
	return cps, err
}

//Returns the number of prices that a customer has
func GetPricingCount(dtx *apicontext.DataContext) (int, error) {
	var count int
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return count, err
	}
	defer db.Close()
	stmt, err := db.Prepare(getPricingCount)
	if err != nil {
		return count, err
	}
	defer stmt.Close()
	err = stmt.QueryRow(dtx.BrandID, dtx.APIKey).Scan(&count)
	if err != nil {
		return count, err
	}
	return count, err
}

//Returns Price for a part
func GetPartPricesByPartID(partID int, dtx *apicontext.DataContext) ([]Price, error) {
	var ps []Price
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return ps, err
	}
	defer db.Close()
	stmt, err := db.Prepare(getPricingByPart)
	if err != nil {
		return ps, err
	}
	defer stmt.Close()
	res, err := stmt.Query(dtx.BrandID, partID)
	if err != nil {
		return ps, err
	}
	for res.Next() {
		p, err := ScanPrice(res)
		if err != nil {
			return ps, err
		}
		ps = append(ps, p)
	}
	return ps, err
}

//Returns all Prices
func GetPartPrices(dtx *apicontext.DataContext) ([]Price, error) {
	var ps []Price
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return ps, err
	}
	defer db.Close()
	stmt, err := db.Prepare(getAllPricing)
	if err != nil {
		return ps, err
	}
	defer stmt.Close()
	res, err := stmt.Query(dtx.BrandID)
	if err != nil {
		return ps, err
	}
	for res.Next() {
		p, err := ScanPrice(res)
		if err != nil {
			return ps, err
		}
		ps = append(ps, p)
	}
	return ps, err
}

//Returns Map Price for every part
func GetMAPPartPrices(dtx *apicontext.DataContext) ([]Price, error) {
	var ps []Price
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return ps, err
	}
	defer db.Close()
	stmt, err := db.Prepare(getAllMAPPricing)
	if err != nil {
		return ps, err
	}
	defer stmt.Close()
	res, err := stmt.Query(dtx.BrandID)
	if err != nil {
		return ps, err
	}
	for res.Next() {
		p, err := ScanPrice(res)
		if err != nil {
			return ps, err
		}
		ps = append(ps, p)
	}
	return ps, err
}

//CRUD
func (c *CustomerPrice) Update() error {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	stmt, err := db.Prepare(updateCustomerPrice)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(c.Price, c.IsSale, c.SaleStart, c.SaleEnd, c.CustID, c.PartID)
	return err
}

func (c *CustomerPrice) Create() error {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	stmt, err := db.Prepare(insertCustomerPrice)
	if err != nil {
		return err
	}
	defer stmt.Close()
	res, err := stmt.Exec(c.CustID, c.PartID, c.Price, c.IsSale, c.SaleStart, c.SaleEnd)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	c.ID = int(id)
	return nil
}

func (c *CustomerPrice) Delete() error {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	stmt, err := db.Prepare(deleteCustomerPrice)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(c.ID)
	return err
}

//CartIntegration
func GetCustomerCartIntegrations(dtx *apicontext.DataContext) ([]CustomerPrice, error) {
	var cps []CustomerPrice
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return cps, err
	}
	defer db.Close()
	stmt, err := db.Prepare(getCustomerCartIntegrations)
	if err != nil {
		return cps, err
	}
	defer stmt.Close()
	res, err := stmt.Query(dtx.APIKey, dtx.BrandID)
	if err != nil {
		return cps, err
	}
	for res.Next() {
		c, err := ScanCartIntegration(res)
		if err != nil {
			return cps, err
		}
		cps = append(cps, c)
	}
	return cps, err
}

func (cp *CustomerPrice) UpdateCartIntegration() error {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	stmt, err := db.Prepare(updateCartIntegration)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(cp.CustomerPartID, cp.PartID, cp.CustID)
	return err
}

func (cp *CustomerPrice) InsertCartIntegration() error {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	stmt, err := db.Prepare(insertCartIntegration)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(cp.CustomerPartID, cp.PartID, cp.CustID)
	return err
}

func GetAllPriceTypes() ([]string, error) {
	var types []string
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return types, err
	}
	defer db.Close()
	stmt, err := db.Prepare(getAllPriceTypes)
	if err != nil {
		return types, err
	}
	defer stmt.Close()
	res, err := stmt.Query()
	if err != nil {
		return types, err
	}
	var s string
	for res.Next() {
		err = res.Scan(&s)
		if err != nil {
			return types, err
		}
		types = append(types, s)
	}
	return types, err
}

//Utility
func Scan(rows database.Scanner) (CustomerPrice, error) {
	var c CustomerPrice
	var p *float64
	err := rows.Scan(
		&c.ID,
		&c.CustID,
		&c.PartID,
		&p,
		&c.IsSale,
		&c.SaleStart,
		&c.SaleEnd,
	)
	if err != nil {
		return c, err
	}
	if p != nil {
		c.Price = *p
	}
	return c, err
}

func ScanPrice(rows database.Scanner) (Price, error) {
	var p Price
	err := rows.Scan(
		&p.PartID,
		&p.Type,
		&p.Price,
	)
	return p, err
}

func ScanCartIntegration(rows database.Scanner) (CustomerPrice, error) {
	var c CustomerPrice
	var throwaway int
	err := rows.Scan(
		&throwaway,
		&c.PartID,
		&c.CustomerPartID,
		&c.CustID,
	)
	return c, err
}
