package cartIntegration

import (
	"database/sql"
	"fmt"
	"github.com/curt-labs/GoAPI/helpers/database"
	"github.com/curt-labs/GoAPI/models/products"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
	"time"
)

const (
	NULL_DATE_STRING  = "0000-00-00"
	MYSQL_DATE_FORMAT = "YYYY-MM-DD"
)

type PricePoint struct {
	CartIntegration CartIntegration `json:"cartIntegration,omitempty" xml:"id,omitempty"`
	Price           float64         `json:"price,omitempty" xml:"price,omitempty"`
	IsSale          int             `json:"isSale,omitempty" xml:"isSale,omitempty"`
	Sale_start      time.Time       `json:"saleStart,omitempty" xml:"saleStart,omitempty"`
	Sale_end        time.Time       `json:"saleEnd,omitempty" xml:"saleEnd,omitempty"`
	PriceString     string          `json:"priceString,omitempty" xml:"priceString,omitempty"`
	IsSaleString    string          `json:"isSaleString,omitempty" xml:"isSaleString,omitempty"`
	SaleStartString string          `json:"saleStartString,omitempty" xml:"saleStartString,omitempty"`
	SaleEndString   string          `json:"saleEndString,omitempty" xml:"saleEndString,omitempty"`
}

type PriceMatrix struct {
	PartID int
	Prices []products.Part
}

var (
	getCustomerPricing = `SELECT p.partID, ci.custPartID, cp.price, cp.isSale, cp.sale_start, cp.sale_end FROM Part p
							LEFT JOIN CustomerPricing cp ON cp.partID = p.partID AND cp.cust_id = (select cust_id from Customer where customerID = ?)
							LEFT JOIN CartIntegration ci ON ci.partID = p.partID AND ci.custID = (select cust_id from Customer where customerID = ?)
							WHERE p.status = 800 OR p.status = 900
							ORDER BY p.partID`
	getLimitedCustomerPricing = `SELECT p.partID, ci.custPartID, cp.price, cp.isSale, cp.sale_start, cp.sale_end FROM Part p
							LEFT JOIN CustomerPricing cp ON cp.partID = p.partID AND cp.cust_id = (select cust_id from Customer where customerID = ?)
							LEFT JOIN CartIntegration ci ON ci.partID = p.partID AND ci.custID = (select cust_id from Customer where customerID = ?)
							WHERE p.status = 800 OR p.status = 900
							ORDER BY p.partID
							limit ?,?`
	getPricingCountForCustomer = `SELECT count(p.partID) FROM Part p
							LEFT JOIN CartIntegration ci ON ci.partID = p.partID AND ci.custID = (select cust_id from Customer where customerID = ?)
							WHERE p.status = 800 OR p.status = 900
							ORDER BY p.partID`

	getIntegration = `SELECT custPartID FROM CartIntegration WHERE custID = ? AND partID = ? limit 1`
)

//all pricepoints
func GetPricesByCustomerID(custID int) (priceslist []PricePoint, err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return priceslist, err
	}
	defer db.Close()
	stmt, err := db.Prepare(getCustomerPricing)
	if err != nil {
		return priceslist, err
	}
	defer stmt.Close()
	res, err := stmt.Query(custID, custID)
	if err != nil {
		return priceslist, err
	}
	var p PricePoint
	var cpid, isSale *int
	var price *float64
	var ss, se *time.Time
	for res.Next() {
		err = res.Scan(
			&p.CartIntegration.PartID,
			&cpid,
			&price,
			&isSale,
			&ss,
			&se,
		)

		if err != nil {
			return priceslist, err
		}
		if cpid != nil {
			p.CartIntegration.CustPartID = *cpid
		}
		if price != nil {
			p.Price = *price
		}
		if isSale != nil {
			p.IsSale = *isSale
		}
		if ss != nil {
			p.Sale_start = *ss
		}
		if se != nil {
			p.Sale_end = *se
		}

		err = p.toString()
		priceslist = append(priceslist, p)
	}
	return priceslist, err
}

//get paged pricePoints
func GetPricesByCustomerIDPaged(custID int, page, count int) (priceslist []PricePoint, err error) {
	page = ((page - 1) * count)
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return priceslist, err
	}
	defer db.Close()
	stmt, err := db.Prepare(getLimitedCustomerPricing)
	if err != nil {
		return priceslist, err
	}
	defer stmt.Close()
	res, err := stmt.Query(custID, custID, page, count)
	var p PricePoint
	var cpid, isSale *int
	var price *float64
	var ss, se *time.Time
	for res.Next() {
		err = res.Scan(
			&p.CartIntegration.PartID,
			&cpid,
			&price,
			&isSale,
			&ss,
			&se,
		)
		if err != nil {
			return priceslist, err
		}
		if cpid != nil {
			p.CartIntegration.CustPartID = *cpid
		}
		if price != nil {
			p.Price = *price
		}
		if isSale != nil {
			p.IsSale = *isSale
		}
		if ss != nil {
			p.Sale_start = *ss
		}
		if se != nil {
			p.Sale_end = *se
		}

		err = p.toString()
		priceslist = append(priceslist, p)
	}
	defer res.Close()
	return priceslist, err
}

func (p *PricePoint) toString() (err error) {
	p.PriceString = "$" + strconv.FormatFloat(p.Price, 'f', 2, 64)
	if p.IsSale == 1 {
		p.IsSaleString = "Yes"
		p.SaleStartString = p.Sale_start.String()
		p.SaleEndString = p.Sale_end.String()
	} else {
		p.IsSaleString = "No"
	}
	return nil
}

func GetPricingCount(custID int) (count int, err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return count, err
	}
	defer db.Close()
	stmt, err := db.Prepare(getPricingCountForCustomer)
	if err != nil {
		return count, err
	}
	defer stmt.Close()
	err = stmt.QueryRow(custID).Scan(&count)
	if err != nil {
		return count, err
	}
	return count, err
}

//what? Cartintegration has some of these types of array-find functions. Why? Who knows?
func GetPriceMatrix(partID int, matrices []PriceMatrix) (pm PriceMatrix) {
	for _, matrix := range matrices {
		if matrix.PartID == partID {
			pm = matrix
			return
		}
	}
	return
}

//Useful to CRUD customer prices - customer/price
func (p *PricePoint) GetCustPriceID() (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(getIntegration)
	if err != nil {
		return err
	}
	defer stmt.Close()

	row := stmt.QueryRow(p.CartIntegration.CustID, p.CartIntegration.PartID)
	if row == nil {
		return fmt.Errorf("%s", "failed to retrieve integration")
	}

	return row.Scan(&p.CartIntegration.CustPartID)
}
