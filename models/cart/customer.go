package cart

import (
	"fmt"
	"github.com/curt-labs/GoAPI/helpers/database"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Customer struct {
	Id               bson.ObjectId     `json:"id,omitempty" xml:"id" bson:"_id"`
	ShopId           bson.ObjectId     `json:"-" xml:"-" bson:"shop_id"`
	AcceptsMarketing bool              `json:"accepts_marketing" xml:"accepts_marketing,attr" bson:"accepts_marketing"`
	Addresses        []CustomerAddress `json:"addresses" xml:"addresses>addres" bson:"addresses"`
	DefaultAddress   *CustomerAddress  `json:"default_address" xml:"default_address" bson:"default_address"`
	CreatedAt        time.Time         `json:"created_at" xml:"created_at,attr" bson:"created_at"`
	Email            string            `json:"email" xml:"email,attr" bson:"email"`
	FirstName        string            `json:"first_name" xml:"first_name,attr" bson:"first_name"`
	LastName         string            `json:"last_name" xml:"last_name,attr" bson:"last_name"`
	MetaFields       []MetaField       `json:"meta_fields" xml:"meta_fields>meta_field" bson:"meta_fields"`
	LastOrderId      *bson.ObjectId    `json:"last_order_id" xml:"last_order_id,attr" bson:"last_order_id"`
	LastOrderName    string            `json:"last_order_name,omitempty" xml:"last_order_name,attr,omitempty" bson:"last_order_name"`
	OrdersCount      int               `json:"orders_count" xml:"orders_count" bson:"orders_count"`
	Note             string            `json:"note" xml:"note" bson:"note"`
	State            string            `json:"state" xml:"state,attr" bson:"state"` // Disabled, invited, etc.
	Tags             []string          `json:"tags" xml:"tags>tag" bson:"tags"`
	TotalSpent       float64           `json:"total_spent" xml:"total_spent,attr" bson:"total_spent"`
	UpdatedAt        time.Time         `json:"updated_at" xml:"updated_at,attr" bson:"updated_at"`
	VerifiedEmail    bool              `json:"verified_email" xml:"verified_email,attr" bson:"verified_email"`
	Orders           []Order           `json:"orders" xml:"orders" bson:"orders"`
}

type MetaField struct {
	Key       string `json:"key" xml:"key,attr" bson:"key"`
	Namespace string `json:"namespace" xml:"namespace,attr" bson:"namespace"`
	ValueType string `json:"value_type" xml:"value_type,attr" bson:"value_type"`
	Value     string `json:"value" xml:"value,attr" bson:"value"`
}

// Get all customers since a defined Id.
func CustomersSinceId(shopId bson.ObjectId, since_id bson.ObjectId, page, limit int, created_at_min, created_at_max, updated_at_min, updated_at_max *time.Time) ([]Customer, error) {
	custs := []Customer{}
	sess, err := mgo.DialWithInfo(database.MongoConnectionString())
	if err != nil {
		return custs, err
	}
	defer sess.Close()

	c := sess.DB("CurtCart").C("customer")
	qs := bson.M{
		"shop_id": shopId.String(),
		"_id": bson.M{
			"$gt": since_id.String(),
		},
	}
	if created_at_min != nil || created_at_max != nil {
		createdQs := bson.M{}
		if created_at_min != nil {
			createdQs["&qt"] = created_at_min.String()
		}
		if created_at_max != nil {
			createdQs["&lt"] = created_at_max.String()
		}
		qs["created_at"] = createdQs
	}
	if updated_at_min != nil || updated_at_max != nil {
		updatedQs := bson.M{}
		if updated_at_min != nil {
			updatedQs["&qt"] = updated_at_min.String()
		}
		if updated_at_max != nil {
			updatedQs["&lt"] = updated_at_max.String()
		}
		qs["updated_at"] = updatedQs
	}

	if page == 1 {
		page = 0
	}
	c.Find(qs).Skip(page * limit).Limit(limit)

	return custs, err
}

// Get all customers.
func GetCustomers(id bson.ObjectId, page, limit int, created_at_min, created_at_max, updated_at_min, updated_at_max *time.Time) ([]Customer, error) {
	custs := []Customer{}
	sess, err := mgo.DialWithInfo(database.MongoConnectionString())
	if err != nil {
		return custs, err
	}
	defer sess.Close()

	c := sess.DB("CurtCart").C("customer")
	qs := bson.M{
		"shop_id": id,
	}
	if created_at_min != nil || created_at_max != nil {
		createdQs := bson.M{}
		if created_at_min != nil {
			createdQs["&qt"] = created_at_min.String()
		}
		if created_at_max != nil {
			createdQs["&lt"] = created_at_max.String()
		}
		qs["created_at"] = createdQs
	}
	if updated_at_min != nil || updated_at_max != nil {
		updatedQs := bson.M{}
		if updated_at_min != nil {
			updatedQs["&qt"] = updated_at_min.String()
		}
		if updated_at_max != nil {
			updatedQs["&lt"] = updated_at_max.String()
		}
		qs["updated_at"] = updatedQs
	}

	if page == 1 {
		page = 0
	}
	c.Find(qs).Skip(page * limit).Limit(limit).All(&custs)

	return custs, err
}

func CustomerCount(shopId bson.ObjectId) (int, error) {
	if shopId.Hex() == "" {
		return 0, fmt.Errorf("error: %s", "invalid shop reference")
	}

	sess, err := mgo.DialWithInfo(database.MongoConnectionString())
	if err != nil {
		return 0, err
	}
	defer sess.Close()

	return sess.DB("CurtCart").C("customer").Find(bson.M{"shop_id": shopId}).Count()
}

func SearchCustomers(query string, shopId bson.ObjectId) ([]Customer, error) {
	var custs []Customer
	if query == "" {
		return custs, fmt.Errorf("error: %s", "invalid query")
	}

	sess, err := mgo.DialWithInfo(database.MongoConnectionString())
	if err != nil {
		return custs, err
	}
	defer sess.Close()

	qs := bson.M{
		"$text": bson.M{
			"$search": query,
		},
		"shop_id": shopId,
	}

	err = sess.DB("CurtCart").C("customer").Find(qs).All(&custs)

	return custs, err
}

// Get a customer.
func (c *Customer) Get() error {
	sess, err := mgo.DialWithInfo(database.MongoConnectionString())
	if err != nil {
		return err
	}
	defer sess.Close()

	col := sess.DB("CurtCart").C("customer")

	return col.Find(bson.M{"_id": c.Id, "shop_id": c.ShopId}).One(&c)
}

// Add new customer.
func (c *Customer) Insert() error {
	if c.Email == "" {
		return fmt.Errorf("error: %s", "invalid email address")
	}
	if c.FirstName == "" {
		return fmt.Errorf("error: %s", "invalid first anem")
	}
	if c.LastName == "" {
		return fmt.Errorf("error: %s", "invalid last name")
	}
	if c.Id.Hex() == "" {
		c.Id = bson.NewObjectId()
	}
	if c.CreatedAt.IsZero() {
		c.CreatedAt = time.Now()
	}
	c.UpdatedAt = time.Now()

	sess, err := mgo.DialWithInfo(database.MongoConnectionString())
	if err != nil {
		return err
	}
	defer sess.Close()

	col := sess.DB("CurtCart").C("customer")

	_, err = col.UpsertId(c.Id, c)
	if err != nil {
		return err
	}

	// index the document
	idx := mgo.Index{
		Key:        []string{"email", "first_name", "last_name", "meta_fields", "note", "state"},
		Background: true,
		Sparse:     false,
		DropDups:   true,
	}
	col.EnsureIndex(idx)

	return nil
}

// Update a customer.
// Updates updated_at, accepts_marketing, addresses, default_address,
// email, first_name, last_name, meta_fields, note, state, tags.
func (c *Customer) Update() error {
	if c.Id.Hex() == "" {
		return fmt.Errorf("error: %s", "cannot update a customer that doesn't exist")
	}
	if c.Email == "" {
		return fmt.Errorf("error: %s", "invalid email address")
	}
	if c.FirstName == "" {
		return fmt.Errorf("error: %s", "invalid first anem")
	}
	if c.LastName == "" {
		return fmt.Errorf("error: %s", "invalid last name")
	}

	c.UpdatedAt = time.Now()

	sess, err := mgo.DialWithInfo(database.MongoConnectionString())
	if err != nil {
		return err
	}
	defer sess.Close()

	var change = mgo.Change{
		ReturnNew: true,
		Update: bson.M{
			"$set": bson.M{
				"accepts_marketing": c.AcceptsMarketing,
				"addresses":         c.Addresses,
				"default_address":   c.DefaultAddress,
				"email":             c.Email,
				"first_name":        c.FirstName,
				"last_name":         c.LastName,
				"meta_fields":       c.MetaFields,
				"note":              c.Note,
				"state":             c.State,
				"tags":              c.Tags,
			},
		},
	}

	_, err = sess.DB("CurtCart").C("customer").Find(bson.M{"_id": c.Id, "shop_id": c.ShopId}).Apply(change, c)

	return err
}

// Delete a customer.
// A customer can't be deleted if they have existing orders
func (c *Customer) Delete() error {
	if c.Id.Hex() == "" {
		return fmt.Errorf("error: %s", "invalid customer reference")
	}

	sess, err := mgo.DialWithInfo(database.MongoConnectionString())
	if err != nil {
		return err
	}
	defer sess.Close()

	if err := c.Get(); err != nil {
		return err
	}

	if c.Orders != nil && len(c.Orders) > 0 {
		return fmt.Errorf("error: %s", "can't remove a customer that has order information")
	}

	return sess.DB("CurtCart").C("customer").RemoveId(c.Id)
}
