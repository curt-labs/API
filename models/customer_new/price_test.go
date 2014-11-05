package customer_new

// import (
// 	"database/sql"
// 	"github.com/curt-labs/GoAPI/helpers/database"
// 	_ "github.com/go-sql-driver/mysql"
// 	. "github.com/smartystreets/goconvey/convey"
// 	"math/rand"
// 	"testing"
// 	"time"
// )

// func getCustWithSale() (custId int) {
// 	db, err := sql.Open("mysql", database.ConnectionString())
// 	if err != nil {
// 		return 0
// 	}
// 	defer db.Close()
// 	var id int
// 	stmt, err := db.Prepare("select cust_id from CustomerPricing where sale_start > 2013 ORDER BY RAND() LIMIT 1")
// 	if err != nil {
// 		return 0
// 	}
// 	defer stmt.Close()
// 	err = stmt.QueryRow().Scan(&id)
// 	//need customerID, not cust_id
// 	stmt, err = db.Prepare("SELECT customerID FROM Customer WHERE cust_id = ?")
// 	if err != nil {
// 		return 0
// 	}
// 	defer stmt.Close()
// 	err = stmt.QueryRow(id).Scan(&custId)

// 	return custId
// }

// func TestCustomerPriceModel(t *testing.T) {
// 	Convey("Testing Price - Gets", t, func() {
// 		Convey("Testing GetAllPrices()", func() {
// 			ps, err := GetAllPrices()
// 			So(len(ps), ShouldBeGreaterThan, 200000)
// 			So(err, ShouldBeNil)
// 		})
// 		Convey("Gets random CustomerPrice", func() {
// 			ps, err := GetAllPrices()
// 			So(err, ShouldBeNil)
// 			if len(ps) > 0 {
// 				x := rand.Intn(len(ps))
// 				p := ps[x]

// 				Convey("Testing Get()", func() {
// 					err := p.Get()
// 					So(p.Price, ShouldHaveSameTypeAs, 0.00)
// 					So(p, ShouldNotBeNil)
// 					So(err, ShouldBeNil)
// 				})

// 				Convey("Testing GetPricesByCustomer()", func() {
// 					var c Customer
// 					c.Id = p.CustID
// 					custPrices, err := c.GetPricesByCustomer()
// 					So(custPrices, ShouldNotBeNil)
// 					So(err, ShouldBeNil)
// 				})
// 				Convey("Testing GetPricesByPart()", func() {
// 					partID := p.PartID
// 					prices, err := GetPricesByPart(partID)
// 					So(len(prices), ShouldNotBeNil)
// 					So(err, ShouldBeNil)
// 				})
// 				Convey("Testing GetPricesBySaleRange", func() {
// 					var s time.Time
// 					var e time.Time
// 					id := getCustWithSale()
// 					c := Customer{Id: id}
// 					var err error
// 					s, err = time.Parse(inputTimeFormat, "2006-01-02 15:04:05")
// 					e, err = time.Parse(inputTimeFormat, "2016-01-02 15:04:05")
// 					prices, err := c.GetPricesBySaleRange(s, e)
// 					t.Log(c)
// 					if err != sql.ErrNoRows {
// 						So(err, ShouldBeNil)
// 						So(len(prices), ShouldBeGreaterThan, 0)
// 						So(prices, ShouldNotBeNil)
// 					}
// 				})

// 				Convey("Testing Price -  CUD", func() {
// 					Convey("Testing Create() Update() Delete() Price", func() {
// 						var pr Price
// 						var err error
// 						pr.CustID = p.CustID
// 						pr.SaleEnd, err = time.Parse(inputTimeFormat, "02/12/2006")
// 						pr.IsSale = 1
// 						pr.Price = 666.00
// 						err = pr.Create()
// 						So(err, ShouldBeNil)
// 						pr.SaleStart, err = time.Parse(inputTimeFormat, "01/02/2007")
// 						err = pr.Update()
// 						So(err, ShouldBeNil)
// 						err = pr.Get()
// 						So(err, ShouldBeNil)
// 						t, err := time.Parse(inputTimeFormat, "02/12/2006")
// 						So(pr.SaleStart, ShouldHaveSameTypeAs, t)
// 						err = pr.Delete()
// 						So(err, ShouldBeNil)
// 					})

// 				})
// 			}
// 		})
// 	})
// }
