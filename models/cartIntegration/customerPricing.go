package cartIntegration

import (
	"github.com/curt-labs/API/helpers/database"
	_ "github.com/go-sql-driver/mysql"

	"database/sql"
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
	ListPrice      Price      `json:"listPrice,omitempty" xml:"listPrice,omitempty"`
}

type Price struct {
	PartID int     `json:"partId,omitempty" xml:"partId,omitempty"`
	Type   string  `json:"type,omitempty" xml:"type,omitempty"`
	Price  float64 `json:"price,omitempty" xml:"price,omitempty"`
}

var (
	getPricing = `SELECT distinct cp.cust_price_id, cp.cust_id, p.partID, ci.custPartID, cp.price, cp.isSale, cp.sale_start, cp.sale_end, pr.priceType, pr.price FROM Part p
		LEFT JOIN CustomerPricing cp ON cp.partID = p.partID AND cp.cust_id = ?
		LEFT JOIN CartIntegration ci ON ci.partID = p.partID AND ci.custID = ?
		left join Price pr on pr.partID = p.partID and pr.priceType = 'list'
		WHERE (p.status = 800 OR p.status = 900) && p.brandID = ?
		ORDER BY p.partID`
	getPricingPaged = `SELECT distinct cp.cust_price_id, cp.cust_id, p.partID, ci.custPartID, cp.price, cp.isSale, cp.sale_start, cp.sale_end, pr.priceType, pr.price FROM Part p
		LEFT JOIN CustomerPricing cp ON cp.partID = p.partID AND cp.cust_id = ?
		LEFT JOIN CartIntegration ci ON ci.partID = p.partID AND ci.custID = ?
		left join Price pr on pr.partID = p.partID and pr.priceType = 'list'
		WHERE (p.status = 800 OR p.status = 900) && p.brandID = ?
		ORDER BY p.partID
		LIMIT ?, ?`
	getPricingCount = `SELECT count(distinct cp.cust_price_id) FROM Part p
		LEFT JOIN CustomerPricing cp ON cp.partID = p.partID AND cp.cust_id = ?
		LEFT JOIN CartIntegration ci ON ci.partID = p.partID AND ci.custID = ?
		WHERE (p.status = 800 OR p.status = 900) && p.brandID = ?
		ORDER BY p.partID`
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
	updateCustomerPrice         = `UPDATE CustomerPricing SET price = ?, isSale = ?, sale_start = ?, sale_end = ? WHERE cust_id = ? AND partID = ?`
	insertCustomerPrice         = `INSERT INTO CustomerPricing(cust_id, partID, price, isSale, sale_start, sale_end) VALUES(?, ?, ?, ?, ?, ?)`
	deleteCustomerPrice         = `delete from CustomerPricing where cust_price_id = ?`
	getCustomerCartIntegrations = `select c.referenceID, c.partID, c.custPartID, c.custID from CartIntegration as c
		join CustomerUser as cu on cu.cust_id = c.custID
		join ApiKey as a on a.user_id = cu.id
		join Part as p on p.partID = c.partID
		where a.api_key = ?
		and p.brandID = ?
		order by p.partID`
	insertCartIntegration = `INSERT INTO CartIntegration(partID, custPartID, custID) VALUES (?, ?, ?)`
	deleteCartIntegration = `delete from CartIntegration where partID = ? and custPartID = ? and custID = ?`
	updateCartIntegration = `UPDATE CartIntegration SET custPartID = ? WHERE partID = ? AND custID = ?`
	getAllPriceTypes      = `SELECT DISTINCT priceType from Price`
)

var (
	Brand_ID    int
	Customer_ID int
)

func initDB() (*sql.DB, error) {
	connStr := database.ConnectionString()
	db, err := sql.Open("mysql", connStr)
	if err != nil {
		return nil, err
	}
	return db, err
}

//Get all of a single customer's prices
func GetCustomerPrices() ([]CustomerPrice, error) {
	var cps []CustomerPrice
	db, err := initDB()
	if err != nil {
		return cps, err
	}
	defer db.Close()
	stmt, err := db.Prepare(getPricing)
	if err != nil {
		return cps, err
	}
	defer stmt.Close()
	res, err := stmt.Query(Customer_ID, Customer_ID, Brand_ID)
	if err != nil {
		return cps, err
	}

	for res.Next() {
		c, err := Scan(res)
		if err != nil {
			return cps, err
		}
		cps = append(cps, c)
	}
	return cps, err
}

//Get a customers prices - paged/limited
func GetPricingPaged(page int, count int) ([]CustomerPrice, error) {
	var cps []CustomerPrice
	db, err := initDB()
	if err != nil {
		return cps, err
	}
	defer db.Close()

	stmt, err := db.Prepare(getPricingPaged)
	if err != nil {
		return cps, err
	}
	defer stmt.Close()
	res, err := stmt.Query(Customer_ID, Customer_ID, Brand_ID, (page-1)*count, count)
	if err != nil {
		return cps, err
	}

	for res.Next() {
		c, err := Scan(res)
		if err != nil {
			return cps, err
		}
		cps = append(cps, c)
	}
	return cps, err
}

//Returns the number of prices that a customer has
func GetPricingCount() (int, error) {
	var count int
	db, err := initDB()
	if err != nil {
		return count, err
	}
	defer db.Close()
	stmt, err := db.Prepare(getPricingCount)
	if err != nil {
		return count, err
	}
	defer stmt.Close()
	err = stmt.QueryRow(Customer_ID, Customer_ID, Brand_ID).Scan(&count)
	if err != nil {
		return count, err
	}
	return count, err
}

//Returns Price for a part
func GetPartPricesByPartID(partID int) ([]Price, error) {
	var ps []Price
	db, err := initDB()
	if err != nil {
		return ps, err
	}
	defer db.Close()
	stmt, err := db.Prepare(getPricingByPart)
	if err != nil {
		return ps, err
	}
	defer stmt.Close()
	res, err := stmt.Query(Brand_ID, partID)
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
func GetPartPrices() ([]Price, error) {
	var ps []Price
	db, err := initDB()
	if err != nil {
		return ps, err
	}
	defer db.Close()
	stmt, err := db.Prepare(getAllPricing)
	if err != nil {
		return ps, err
	}
	defer stmt.Close()
	res, err := stmt.Query(Brand_ID)
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
func GetMAPPartPrices() ([]Price, error) {
	var ps []Price
	db, err := initDB()
	if err != nil {
		return ps, err
	}
	defer db.Close()
	stmt, err := db.Prepare(getAllMAPPricing)
	if err != nil {
		return ps, err
	}
	defer stmt.Close()
	res, err := stmt.Query(Brand_ID)
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
	db, err := initDB()
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
	db, err := initDB()
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
	db, err := initDB()
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
func GetCustomerCartIntegrations(key string) ([]CustomerPrice, error) {
	var cps []CustomerPrice
	db, err := initDB()
	if err != nil {
		return cps, err
	}
	defer db.Close()
	stmt, err := db.Prepare(getCustomerCartIntegrations)
	if err != nil {
		return cps, err
	}
	defer stmt.Close()
	res, err := stmt.Query(key, Brand_ID)
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
	db, err := initDB()
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
	db, err := initDB()
	if err != nil {
		return err
	}
	defer db.Close()
	stmt, err := db.Prepare(insertCartIntegration)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(cp.PartID, cp.CustomerPartID, cp.CustID)
	return err
}

func (cp *CustomerPrice) DeleteCartIntegration() error {
	db, err := initDB()
	if err != nil {
		return err
	}
	defer db.Close()
	stmt, err := db.Prepare(deleteCartIntegration)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(cp.PartID, cp.CustomerPartID, cp.CustID)
	return err
}

func GetAllPriceTypes() ([]string, error) {
	var types []string
	db, err := initDB()
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
	var p, lp *float64
	var custPartId, id, custId, isSale *int
	var ss, se *time.Time
	var pt *string

	err := rows.Scan(
		&id,
		&custId,
		&c.PartID,
		&custPartId,
		&p,
		&isSale,
		&ss,
		&se,
		&pt,
		&lp,
	)
	if err != nil {
		return c, err
	}

	if id != nil {
		c.ID = *id
	}
	if custId != nil {
		c.CustID = *custId
	}
	if p != nil {
		c.Price = *p
	}
	if custPartId != nil {
		c.CustomerPartID = *custPartId
	}
	if isSale != nil {
		c.IsSale = *isSale
	}
	if ss != nil {
		c.SaleStart = ss
	}
	if se != nil {
		c.SaleEnd = se
	}
	if pt != nil {
		c.ListPrice.Type = *pt
	}
	if lp != nil {
		c.ListPrice.Price = *lp
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
