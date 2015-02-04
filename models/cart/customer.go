package cart

import (
	"code.google.com/p/go.crypto/bcrypt"
	"fmt"
	"github.com/curt-labs/GoAPI/helpers/database"
	jwt "github.com/dgrijalva/jwt-go"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

const (
	jwtSigningKey = "goapi_curt"
)

type Customer struct {
	Id               bson.ObjectId     `json:"id,omitempty" xml:"id" bson:"_id"`
	ShopId           bson.ObjectId     `json:"-" xml:"-" bson:"shop_id"`
	AcceptsMarketing bool              `json:"accepts_marketing" xml:"accepts_marketing,attr" bson:"accepts_marketing"`
	Addresses        []CustomerAddress `json:"addresses" xml:"addresses>addres" bson:"addresses"`
	DefaultAddress   *CustomerAddress  `json:"default_address" xml:"default_address" bson:"default_address"`
	CreatedAt        time.Time         `json:"created_at" xml:"created_at,attr" bson:"created_at"`
	Email            string            `json:"email" xml:"email,attr" bson:"email"`
	Password         string            `json:"password,omitempty" xml:"password,attr,omitempty" bson:"password"`
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
	Token            string            `json:"token" xml:"token,attr" bson:"token"`
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
	err = c.Find(qs).Skip(page * limit).Limit(limit).All(&custs)
	if err != nil {
		return []Customer{}, err
	}

	for i, _ := range custs {
		custs[i].Password = ""
	}

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

	err = c.Find(qs).Skip(page * limit).Limit(limit).All(&custs)
	if err != nil {
		return []Customer{}, nil
	}

	for i, _ := range custs {
		custs[i].Password = ""
	}

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
		"shop_id":  shopId,
		"password": 0,
	}

	err = sess.DB("CurtCart").C("customer").Find(qs).All(&custs)
	if err != nil {
		return []Customer{}, err
	}

	for i, _ := range custs {
		custs[i].Password = ""
	}

	return custs, err
}

// Login a customer.
func (c *Customer) Login(ref string) error {
	pass := c.Password
	c.Password = ""

	sess, err := mgo.DialWithInfo(database.MongoConnectionString())
	if err != nil {
		return err
	}
	defer sess.Close()

	col := sess.DB("CurtCart").C("customer")

	var custs []Customer
	err = col.Find(bson.M{"email": c.Email, "shop_id": c.ShopId}).All(&custs)
	if err != nil {
		return err
	}

	if custs == nil || len(custs) == 0 {
		return fmt.Errorf("error: %s", "no account for this email address")
	}

	for _, cust := range custs {
		if err := bcrypt.CompareHashAndPassword([]byte(cust.Password), []byte(pass)); err != nil {
			continue
		}

		c.Id = cust.Id
		c.ShopId = cust.ShopId
		c.AcceptsMarketing = cust.AcceptsMarketing
		c.Addresses = cust.Addresses
		c.DefaultAddress = cust.DefaultAddress
		c.CreatedAt = cust.CreatedAt
		c.Email = cust.Email
		c.FirstName = cust.FirstName
		c.LastName = cust.LastName
		c.MetaFields = cust.MetaFields
		c.LastOrderId = cust.LastOrderId
		c.LastOrderName = cust.LastOrderName
		c.OrdersCount = cust.OrdersCount
		c.Note = cust.Note
		c.State = cust.State
		c.Tags = cust.Tags
		c.UpdatedAt = cust.UpdatedAt
		c.VerifiedEmail = cust.VerifiedEmail
		c.Orders = cust.Orders
		c.generateToken(ref)

		return nil
	}

	return fmt.Errorf("error: %s", "credentials do not match")
}

// Get a customer.
func (c *Customer) Get() error {
	sess, err := mgo.DialWithInfo(database.MongoConnectionString())
	if err != nil {
		return err
	}
	defer sess.Close()

	col := sess.DB("CurtCart").C("customer")

	err = col.Find(bson.M{"_id": c.Id, "shop_id": c.ShopId}).One(&c)
	if err != nil {
		return err
	}
	c.Password = ""

	return nil
}

// Get a customer by email.
func (c *Customer) GetByEmail() error {
	sess, err := mgo.DialWithInfo(database.MongoConnectionString())
	if err != nil {
		return err
	}
	defer sess.Close()

	col := sess.DB("CurtCart").C("customer")

	err = col.Find(bson.M{"email": c.Email, "shop_id": c.ShopId}).One(&c)
	if err != nil {
		return err
	}
	c.Password = ""

	return nil
}

// Add new customer.
func (c *Customer) Insert(ref string) error {
	if c.Email == "" {
		c.Password = ""
		return fmt.Errorf("error: %s", "invalid email address")
	}
	if c.Password == "" {
		return fmt.Errorf("error: %s", "invalid password")
	}
	if c.FirstName == "" {
		c.Password = ""
		return fmt.Errorf("error: %s", "invalid first name")
	}
	if c.LastName == "" {
		c.Password = ""
		return fmt.Errorf("error: %s", "invalid last name")
	}
	if c.Id.Hex() == "" {
		c.Id = bson.NewObjectId()
	}
	if c.CreatedAt.IsZero() {
		c.CreatedAt = time.Now()
	}
	c.UpdatedAt = time.Now()

	cryptic, err := bcrypt.GenerateFromPassword([]byte(c.Password), bcrypt.DefaultCost)
	if err != nil {
		c.Password = ""
		return fmt.Errorf("error: %s", err.Error())
	}
	pass := string(cryptic)
	c.generateToken(ref)
	c.Password = pass

	sess, err := mgo.DialWithInfo(database.MongoConnectionString())
	if err != nil {
		c.Password = ""
		return err
	}
	defer sess.Close()

	col := sess.DB("CurtCart").C("customer")

	_, err = col.UpsertId(c.Id, c)
	if err != nil {
		c.Password = ""
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
	c.Password = ""

	return nil
}

// Update a customer.
// Updates updated_at, accepts_marketing, addresses, default_address,
// email, first_name, last_name, meta_fields, note, state, tags.
func (c *Customer) Update() error {
	if c.Id.Hex() == "" {
		c.Password = ""
		return fmt.Errorf("error: %s", "cannot update a customer that doesn't exist")
	}
	if c.Email == "" {
		c.Password = ""
		return fmt.Errorf("error: %s", "invalid email address")
	}
	if c.FirstName == "" {
		c.Password = ""
		return fmt.Errorf("error: %s", "invalid first anem")
	}
	if c.LastName == "" {
		c.Password = ""
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

	c.Password = ""
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

func (c *Customer) generateToken(referer string) error {
	c.Password = ""
	var err error

	// create token
	token := jwt.New(jwt.SigningMethodHS256)

	// assign claims
	token.Claims["iss"] = "carter.curtmfg.com"
	token.Claims["sub"] = referer
	token.Claims["exp"] = time.Now().Add(time.Hour * 1).Unix()
	token.Claims["iat"] = time.Now().Unix()

	c.Token, err = token.SignedString([]byte(jwtSigningKey))
	if err != nil {
		return err
	}

	sess, err := mgo.DialWithInfo(database.MongoConnectionString())
	if err != nil {
		return err
	}
	defer sess.Close()

	var change = mgo.Change{
		ReturnNew: true,
		Update: bson.M{
			"$set": bson.M{
				"token": c.Token,
			},
		},
	}

	_, err = sess.DB("CurtCart").C("customer").Find(bson.M{"_id": c.Id, "shop_id": c.ShopId}).Apply(change, c)

	c.Password = ""
	return err
}
