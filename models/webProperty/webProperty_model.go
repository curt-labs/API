package webProperty_model

import (
	"database/sql"
	"encoding/json"
	"github.com/curt-labs/GoAPI/helpers/database"
	"github.com/curt-labs/GoAPI/helpers/pagination"
	"github.com/curt-labs/GoAPI/helpers/redis"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"time"
)

type WebProperty struct {
	ID                      int `json:"id,omitempty" xml:"id,omitempty"`
	Name                    string
	CustID                  int
	BadgeID                 string
	Url                     string
	IsEnabled               bool
	SellerID                string
	WebPropertyNotes        WebPropertyNotes
	WebPropertyType         WebPropertyType
	WebPropertyRequirements WebPropertyRequirements
	IsFinalApproved         bool
	IsEnabledDate           time.Time
	IsDenied                bool
	RequestedDate           time.Time
	AddedDate               time.Time
}

type WebProperties []WebProperty

type WebPropertyType struct {
	ID     int
	TypeID int
	Type   string
}
type WebPropertyTypes []WebPropertyType

type WebPropertyNote struct {
	ID        int
	WebPropID int
	Text      string
	DateAdded time.Time
}

type WebPropertyNotes []WebPropertyNote

type WebPropertyRequirement struct {
	ID            int
	ReqType       string
	Requirement   string
	RequirementID int
	Compliance    bool
	WebPropID     int
}

type WebPropertyRequirements []WebPropertyRequirement

var (
	getAllWebProperties           = "SELECT id, name, cust_ID, badgeID, url, isEnabled,sellerID, typeID , isFinalApproved, isEnabledDate, isDenied, requestedDate, addedDate FROM WebProperties"
	getWebProperty                = "SELECT id, name, cust_ID, badgeID, url, isEnabled,sellerID, typeID , isFinalApproved, isEnabledDate, isDenied, requestedDate, addedDate FROM WebProperties WHERE id = ?"
	getAllWebPropertyTypes        = "SELECT id, typeID, type FROM WebPropertyTypes"
	getAllWebPropertyNotes        = "SELECT id, webPropID, text, dateAdded FROM WebPropNotes"
	getAllWebPropertyRequirements = "SELECT wprc.ID, wpr.ID, wpr.ReqType, wpr.Requirement, wprc.Compliance, wprc.WebPropertiesID FROM WebPropRequirementCheck AS wprc LEFT JOIN WebPropRequirements AS wpr ON wpr.ID = wprc.WebPropRequirementsID"
	create                        = "INSERT INTO WebProperties (name, cust_ID, badgeID, url, isEnabled,sellerID, typeID , isFinalApproved, isEnabledDate, isDenied, requestedDate, addedDate) VALUES (?,?,?,?,?,?,?,?,?,?,?,?)"
	deleteWebProp                 = "DELETE FROM WebProperties WHERE id = ?"
	createNote                    = "INSERT INTO WebPropNotes (webPropID, text, dateAdded) VALUES (?,?,?)"
	updateNote                    = "UPDATE WebPropNotes SET webPropID = ?, text = ?, dateAdded = ? WHERE id =?"
	deleteNote                    = "DELETE FROM WebPropNotes WHERE id = ?"
	createRequirementsBridge      = "INSERT INTO WebPropRequirementCheck (WebPropertiesID, Compliance, WebPropRequirementsID) VALUES (?,?,?)"
	deleteRequirementsBridge      = "DELETE FROM WebPropRequirementCheck WHERE id = ?"
	updateRequirementsBridge      = "UPDATE WebPropRequirementCheck SET WebPropertiesID = ?, Compliance = ?, WebPropRequirementsID = ? WHERE ID = ?"
	update                        = "UPDATE WebProperties SET name = ?, cust_ID = ?, badgeID = ?, url = ?, isEnabled = ?,sellerID = ?, typeID = ?, isFinalApproved = ?, isEnabledDate = ?, isDenied = ?, requestedDate = ? WHERE id = ?"
	search                        = `SELECT id, name, cust_ID, badgeID, url, isEnabled,sellerID, typeID , isFinalApproved, isEnabledDate, isDenied, requestedDate, addedDate FROM WebProperties
									 WHERE  name LIKE ? AND cust_ID LIKE ? AND url LIKE ? AND isEnabled LIKE ? AND sellerID LIKE ? AND typeID  LIKE ? AND isFinalApproved LIKE ? AND isEnabledDate LIKE ? AND
									 isDenied LIKE ? AND requestedDate LIKE ? AND addedDate LIKE ? `
	createRequirement    = "INSERT INTO WebPropRequirements (ReqType, Requirement) VALUES (?,?)"
	updateRequirement    = "UPDATE WebPropRequirements SET ReqType = ?, Requirement = ? WHERE ID = ?"
	deleteRequirement    = "DELETE FROM WebPropRequirements WHERE ID = ?"
	getNote              = "SELECT id, webPropID, text, dateAdded FROM WebPropNotes WHERE id = ?"
	getRequirement       = "SELECT ID, ReqType, Requirement FROM WebPropRequirements WHERE ID = ?"
	getRequirementBridge = "SELECT ID, WebPropertiesID, Compliance, WebPropRequirementsID FROM WebPropRequirementCheck WHERE ID = ?"
)

const (
	timeFormat = "2006-01-02 15:04:05"
)

func (w *WebProperty) Get() error {
	var ws WebProperties
	var err error

	redis_key := "goacespi:webproperties"
	data, err := redis.Get(redis_key)
	if err == nil && len(data) > 0 {
		err = json.Unmarshal(data, &ws)
		return err
	}

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(getWebProperty)
	if err != nil {
		return err
	}
	defer stmt.Close()

	webPropTypes, err := GetAllWebPropertyTypes()
	webPropNotes, err := GetAllWebPropertyNotes()
	WebPropertyRequirements, err := GetAllWebPropertyRequirements()
	if err != nil {
		return err
	}

	typesMap := webPropTypes.ToMap()
	notesMap := webPropNotes.ToMap()
	requirementsMap := WebPropertyRequirements.ToMap()

	res, err := stmt.Query(w.ID)
	for res.Next() {
		res.Scan(&w.ID, &w.Name, &w.CustID, &w.BadgeID, &w.Url, &w.IsEnabled, &w.SellerID, &w.WebPropertyType.ID, &w.IsFinalApproved, &w.IsEnabledDate, &w.IsDenied, &w.RequestedDate, &w.AddedDate)

		typeChan := make(chan int)
		notesChan := make(chan int)
		requirementsChan := make(chan int)
		go func() error {
			for _, val := range typesMap {
				if val.TypeID == w.WebPropertyType.ID {
					w.WebPropertyType = val
				}
			}
			typeChan <- 1
			return nil
		}()
		go func() error {
			for _, val := range notesMap {
				if val.WebPropID == w.ID {
					w.WebPropertyNotes = append(w.WebPropertyNotes, val)
				}
			}
			notesChan <- 1
			return nil
		}()
		go func() error {
			for _, val := range requirementsMap {
				if val.WebPropID == w.ID {
					w.WebPropertyRequirements = append(w.WebPropertyRequirements, val)
				}
			}
			requirementsChan <- 1
			return nil
		}()

		<-typeChan
		<-notesChan
		<-requirementsChan

	}
	go redis.Setex(redis_key, w, 86400)
	return err
}

func GetAll() (WebProperties, error) {
	var ws WebProperties
	var err error

	redis_key := "goacespi:webproperties"
	data, err := redis.Get(redis_key)
	if err == nil && len(data) > 0 {
		err = json.Unmarshal(data, &ws)
		return ws, err
	}

	db, err := sql.Open("mysql", database.ConnectionString())

	if err != nil {
		return ws, err
	}
	defer db.Close()

	stmt, err := db.Prepare(getAllWebProperties)
	if err != nil {
		return ws, err
	}
	defer stmt.Close()

	webPropTypes, err := GetAllWebPropertyTypes()
	webPropNotes, err := GetAllWebPropertyNotes()
	WebPropertyRequirements, err := GetAllWebPropertyRequirements()
	if err != nil {
		return ws, err
	}

	typesMap := webPropTypes.ToMap()
	notesMap := webPropNotes.ToMap()
	requirementsMap := WebPropertyRequirements.ToMap()

	res, err := stmt.Query()
	for res.Next() {
		var w WebProperty
		res.Scan(&w.ID, &w.Name, &w.CustID, &w.BadgeID, &w.Url, &w.IsEnabled, &w.SellerID, &w.WebPropertyType.ID, &w.IsFinalApproved, &w.IsEnabledDate, &w.IsDenied, &w.RequestedDate, &w.AddedDate)

		typeChan := make(chan int)
		notesChan := make(chan int)
		requirementsChan := make(chan int)
		go func() error {

			for _, val := range typesMap {
				if val.TypeID == w.WebPropertyType.ID {
					w.WebPropertyType = val
				}
			}
			typeChan <- 1
			return nil
		}()
		go func() error {

			for _, val := range notesMap {
				if val.WebPropID == w.ID {
					w.WebPropertyNotes = append(w.WebPropertyNotes, val)
				}
			}
			notesChan <- 1
			return nil
		}()
		go func() error {

			for _, val := range requirementsMap {
				if val.WebPropID == w.ID {
					w.WebPropertyRequirements = append(w.WebPropertyRequirements, val)
				}
			}
			requirementsChan <- 1
			return nil
		}()

		<-typeChan
		<-notesChan
		<-requirementsChan

		ws = append(ws, w)
	}
	go redis.Setex(redis_key, ws, 86400)
	return ws, err
}
func (w *WebProperty) Create() error {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(create)
	w.AddedDate = time.Now()
	res, err := stmt.Exec(w.Name, w.CustID, w.BadgeID, w.Url, w.IsEnabled, w.SellerID, w.WebPropertyType.ID, w.IsFinalApproved, w.IsEnabledDate, w.IsDenied, w.RequestedDate, w.AddedDate)
	if err != nil {
		tx.Rollback()
		return err
	}
	id, err := res.LastInsertId()
	w.ID = int(id)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	//create notes
	for _, note := range w.WebPropertyNotes {
		note.WebPropID = w.ID
		err = note.Create()
		if err != nil {
			return err
		}
	}
	//create web properties check
	for _, req := range w.WebPropertyRequirements {
		req.WebPropID = w.ID
		err = req.CreateJoin()
		if err != nil {
			return err
		}
	}
	return nil
}

func (w *WebProperty) Update() error {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(update)
	_, err = stmt.Exec(w.Name, w.CustID, w.BadgeID, w.Url, w.IsEnabled, w.SellerID, w.WebPropertyType.ID, w.IsFinalApproved, w.IsEnabledDate, w.IsDenied, w.RequestedDate, w.ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	return nil
}
func (w *WebProperty) Delete() error {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	webPropNotes, err := GetAllWebPropertyNotes()
	notesMap := webPropNotes.ToMap()

	WebPropertyRequirements, err := GetAllWebPropertyRequirements()
	requirementsMap := WebPropertyRequirements.ToMap()

	notesChan := make(chan int)
	requirementsChan := make(chan int)

	go func() {
		for _, val := range notesMap {
			if val.WebPropID == w.ID {
				val.Delete()
			}
		}
		notesChan <- 1
	}()
	go func() {
		for _, val := range requirementsMap {
			if val.WebPropID == w.ID {
				val.DeleteJoin()
			}
		}
		requirementsChan <- 1
	}()
	<-notesChan
	<-requirementsChan

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(deleteWebProp)
	_, err = stmt.Exec(w.ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	return nil
}

func GetAllWebPropertyTypes() (WebPropertyTypes, error) {
	var ws WebPropertyTypes
	var err error

	redis_key := "goacespi:webpropertytypes"
	data, err := redis.Get(redis_key)
	if err == nil && len(data) > 0 {
		err = json.Unmarshal(data, &ws)
		return ws, err
	}
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return ws, err
	}
	defer db.Close()

	stmt, err := db.Prepare(getAllWebPropertyTypes)
	if err != nil {
		return ws, err
	}
	defer stmt.Close()

	res, err := stmt.Query()
	for res.Next() {
		var w WebPropertyType
		res.Scan(&w.ID, &w.TypeID, &w.Type)
		ws = append(ws, w)
	}
	go redis.Setex(redis_key, ws, 86400)
	return ws, err
}

func GetAllWebPropertyNotes() (WebPropertyNotes, error) {
	var ws WebPropertyNotes
	var err error

	redis_key := "goacespi:webpropertynotes"
	data, err := redis.Get(redis_key)
	if err == nil && len(data) > 0 {
		err = json.Unmarshal(data, &ws)
		return ws, err
	}
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return ws, err
	}
	defer db.Close()

	stmt, err := db.Prepare(getAllWebPropertyNotes)
	if err != nil {
		return ws, err
	}
	defer stmt.Close()

	res, err := stmt.Query()
	for res.Next() {
		var w WebPropertyNote
		res.Scan(&w.ID, &w.WebPropID, &w.Text, &w.DateAdded)
		ws = append(ws, w)
	}
	go redis.Setex(redis_key, ws, 86400)
	return ws, err
}
func GetAllWebPropertyRequirements() (WebPropertyRequirements, error) {
	var ws WebPropertyRequirements
	var err error

	redis_key := "goacespi:webpropertyrequirements"
	data, err := redis.Get(redis_key)
	if err == nil && len(data) > 0 {
		err = json.Unmarshal(data, &ws)
		return ws, err
	}
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return ws, err
	}
	defer db.Close()

	stmt, err := db.Prepare(getAllWebPropertyRequirements)
	if err != nil {
		return ws, err
	}
	defer stmt.Close()

	res, err := stmt.Query()
	for res.Next() {
		var w WebPropertyRequirement
		res.Scan(&w.ID, &w.RequirementID, &w.ReqType, &w.Requirement, &w.Compliance, &w.WebPropID)
		ws = append(ws, w)
	}
	go redis.Setex(redis_key, ws, 86400)
	return ws, err
}

func (n *WebPropertyNote) Get() error {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(getNote)
	err = stmt.QueryRow(n.ID).Scan(&n.ID, &n.WebPropID, &n.Text, &n.DateAdded)
	if err != nil {
		return err
	}

	return nil
}

func (n *WebPropertyNote) Create() error {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(createNote)
	n.DateAdded = time.Now()
	res, err := stmt.Exec(n.WebPropID, n.Text, n.DateAdded)
	if err != nil {
		tx.Rollback()
		return err
	}
	id, err := res.LastInsertId()
	n.ID = int(id)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	return nil
}

func (n *WebPropertyNote) Update() error {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(updateNote)
	n.DateAdded = time.Now()
	_, err = stmt.Exec(n.WebPropID, n.Text, n.DateAdded, n.ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	return nil
}

func (n *WebPropertyNote) Delete() error {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(deleteNote)
	_, err = stmt.Exec(n.ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	return nil
}

func (r *WebPropertyRequirement) GetJoin() error {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(getRequirementBridge)
	err = stmt.QueryRow(r.ID).Scan(&r.ID, &r.WebPropID, &r.Compliance, &r.RequirementID)
	if err != nil {
		return err
	}
	err = r.Get()
	return nil
}

func (r *WebPropertyRequirement) CreateJoin() error {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(createRequirementsBridge)
	_, err = stmt.Exec(r.WebPropID, r.Compliance, r.RequirementID)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	return nil
}

func (r *WebPropertyRequirement) UpdateJoin() error {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(updateRequirementsBridge)
	_, err = stmt.Exec(r.WebPropID, r.Compliance, r.RequirementID, r.ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	r.Get()

	return nil
}

func (r *WebPropertyRequirement) DeleteJoin() error {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(deleteRequirementsBridge)
	_, err = stmt.Exec(r.ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	return nil
}

func (r *WebPropertyRequirement) Get() error {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(getRequirement)
	err = stmt.QueryRow(r.RequirementID).Scan(&r.ID, &r.ReqType, &r.Requirement)
	if err != nil {
		return err
	}

	return nil
}

func (r *WebPropertyRequirement) Create() error {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(createRequirement)
	res, err := stmt.Exec(r.ReqType, r.Requirement)
	if err != nil {
		tx.Rollback()
		return err
	}
	id, err := res.LastInsertId()
	r.ID = int(id)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	return nil
}

func (r *WebPropertyRequirement) Update() error {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	log.Print("ID: ", r.RequirementID, " TYPE:", r.ReqType, " REQ:", r.Requirement)
	stmt, err := tx.Prepare(updateRequirement)
	_, err = stmt.Exec(r.ReqType, r.Requirement, r.RequirementID)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	return nil
}

func (r *WebPropertyRequirement) Delete() error {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(deleteRequirement)
	_, err = stmt.Exec(r.RequirementID)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	return nil
}
func Search(name, custID, badgeID, url, isEnabled, sellerID, webPropertyTypeID, isFinalApproved, isEnabledDate, isDenied, requestedDate, typeID, pageStr, resultsStr string) (pagination.Objects, error) {
	var err error
	var l pagination.Objects
	var fs []interface{}

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return l, err
	}
	defer db.Close()

	stmt, err := db.Prepare(search)
	if err != nil {
		return l, err
	}
	defer stmt.Close()

	res, err := stmt.Query("%"+name+"%", "%"+custID+"%", "%"+url+"%", "%"+isEnabled+"%", "%"+sellerID+"%", "%"+webPropertyTypeID+"%", "%"+isFinalApproved+"%", "%"+isEnabledDate+"%", "%"+isDenied+"%", "%"+requestedDate+"%", "%"+typeID+"%")
	for res.Next() {
		var w WebProperty
		res.Scan(&w.ID, &w.Name, &w.CustID, &w.BadgeID, &w.Url, &w.IsEnabled, &w.SellerID, &w.WebPropertyType.ID, &w.IsFinalApproved, &w.IsEnabledDate, &w.IsDenied, &w.RequestedDate, &w.AddedDate)
		fs = append(fs, w)
	}
	l = pagination.Paginate(pageStr, resultsStr, fs)
	return l, err
}
