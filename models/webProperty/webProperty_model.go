package webProperty_model

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/curt-labs/API/helpers/apicontext"
	"github.com/curt-labs/API/helpers/database"
	"github.com/curt-labs/API/helpers/pagination"
	"github.com/curt-labs/API/helpers/redis"
	_ "github.com/go-sql-driver/mysql"
)

// Web Property is any online site, or presence that sells or markets our products.
type WebProperty struct {
	ID                      int                     `json:"id,omitempty" xml:"id,omitempty"`
	Name                    string                  `json:"name,omitempty" xml:"name,omitempty"`
	CustID                  int                     `json:"custId,omitempty" xml:"custId,omitempty"`
	BadgeID                 string                  `json:"badgeId,omitempty" xml:"badgeId,omitempty"`
	Url                     string                  `json:"url,omitempty" xml:"url,omitempty"`
	IsEnabled               bool                    `json:"isEnabled,omitempty" xml:"isEnabled,omitempty"`
	SellerID                string                  `json:"sellerId,omitempty" xml:"v,omitempty"`
	WebPropertyNotes        WebPropertyNotes        `json:"webPropertyNotes,omitempty" xml:"webPropertyNotes,omitempty"`
	WebPropertyType         WebPropertyType         `json:"webPropertyTypes,omitempty" xml:"webPropertyTypes,omitempty"`
	WebPropertyRequirements WebPropertyRequirements `json:"webPropertyRequirements,omitempty" xml:"webPropertyRequirements,omitempty"`
	IsFinalApproved         bool                    `json:"isFinalApproved,omitempty" xml:"isFinalApproved,omitempty"`
	IsEnabledDate           *time.Time              `json:"isEnabledDate,omitempty" xml:"isEnabledDate,omitempty"`
	IsDenied                bool                    `json:"isDenied,omitempty" xml:"isDenied,omitempty"`
	RequestedDate           *time.Time              `json:"requestedDate,omitempty" xml:"requestedDate,omitempty"`
	AddedDate               *time.Time              `json:"addedDate,omitempty" xml:"addedDate,omitempty"`
}

// WebProperties is just an easier type to work with than using an array of WebProperty's.
type WebProperties []WebProperty

// The Type of a WebProperty. Examples are: Website, Ebay Store, Amazon Store
type WebPropertyType struct {
	ID     int    `json:"id,omitempty" xml:"id,omitempty"`
	TypeID int    `json:"typeId,omitempty" xml:"typeId,omitempty"`
	Type   string `json:"type,omitempty" xml:"type,omitempty"`
}

// WebPropertiesTypes is just an easier type to work with than using an array of WebPropertyType's.
type WebPropertyTypes []WebPropertyType

// WebPropertyNote is just notes about the web property that should be taken into consideration when marketing our products.
type WebPropertyNote struct {
	ID        int        `json:"id,omitempty" xml:"id,omitempty"`
	WebPropID int        `json:"webPropId,omitempty" xml:"webPropId,omitempty"`
	Text      string     `json:"text,omitempty" xml:"text,omitempty"`
	DateAdded *time.Time `json:"dateAdded,omitempty" xml:"dateAdded,omitempty"`
}

// WebPropertiesNotes is just an easier type to work with than using an array of WebPropertyNote's.
type WebPropertyNotes []WebPropertyNote

// WebPropertyRequirement is a requirement for your web property to pass before being approved for an Authorized Dealer badge.
type WebPropertyRequirement struct {
	ID            int    `json:"id,omitempty" xml:"id,omitempty"`
	ReqType       string `json:"reqType,omitempty" xml:"reqType,omitempty"`
	Requirement   string `json:"requirement,omitempty" xml:"requirement,omitempty"`
	RequirementID int    `json:"requirementId,omitempty" xml:"requirementId,omitempty"`
	Compliance    bool   `json:"compliance,omitempty" xml:"compliance,omitempty"`
	WebPropID     int    `json:"webPropId,omitempty" xml:"webPropId,omitempty"`
}

// WebPropertiesRequirements is just an easier type to work with than using an array of WebPropertyRequirement's.
type WebPropertyRequirements []WebPropertyRequirement

var (
	getAllWebProperties = `SELECT w.id, w.name, w.cust_ID, w.badgeID, w.url, w.isEnabled,w.sellerID, w.typeID , w.isFinalApproved, w.isEnabledDate, w.isDenied, w.requestedDate, w.addedDate
				FROM WebProperties as w
				join CustomerToBrand as ctb on ctb.cust_id = w.cust_id
				join ApiKeyToBrand as atb on atb.brandID = ctb.brandID
				join ApiKey as a on a.id = atb.keyID
				where a.api_key = ? && (ctb.brandID = ? or 0 = ?)`
	getWebProperty             = "SELECT id, name, cust_ID, badgeID, url, isEnabled,sellerID, typeID , isFinalApproved, isEnabledDate, isDenied, requestedDate, addedDate FROM WebProperties WHERE id = ?"
	getWebPropertiesByCustomer = "SELECT id, name, cust_ID, badgeID, url, isEnabled,sellerID, typeID , isFinalApproved, isEnabledDate, isDenied, requestedDate, addedDate FROM WebProperties WHERE cust_ID = ?"
	getAllWebPropertyTypes     = `SELECT DISTINCT
		wt.id,
		wt.typeID,
		wt.type
		FROM WebPropertyTypes  AS wt
		JOIN WebProperties   AS w   ON w.typeID    = wt.id
		JOIN CustomerToBrand AS ctb ON ctb.cust_id = w.cust_id
		WHERE ctb.brandID = ?`
	getAllWebPropertyNotes = `SELECT wn.id, wn.webPropID, wn.text, wn.dateAdded
		FROM WebPropNotes AS wn
		JOIN WebProperties AS w ON w.id = wn.webPropID
		JOIN CustomerToBrand AS ctb ON ctb.cust_id = w.cust_id
		WHERE ctb.brandID = ?
		ORDER BY wn.id`
	getAllWebPropertyRequirements = `SELECT DISTINCT wprc.ID, wpr.ID, wpr.ReqType, wpr.Requirement, wprc.Compliance, wprc.WebPropertiesID
		FROM WebPropRequirementCheck AS wprc
		LEFT JOIN WebPropRequirements AS wpr ON wpr.ID = wprc.WebPropRequirementsID
		join WebProperties as w on w.ID = wprc.WebPropertiesID
		join CustomerToBrand as ctb on ctb.cust_id = w.cust_id
		join ApiKeyToBrand as atb on atb.brandID = ctb.brandID
		join ApiKey as a on a.id = atb.keyID
		where a.api_key = ? && (ctb.brandID = ? or 0 = ?)`
	create                                  = "INSERT INTO WebProperties (name, cust_ID, badgeID, url, isEnabled,sellerID, typeID , isFinalApproved, isEnabledDate, isDenied, requestedDate, addedDate) VALUES (?,?,UUID(),?,?,?,?,?,?,?,?,?)"
	deleteWebProp                           = "DELETE FROM WebProperties WHERE id = ?"
	createNote                              = "INSERT INTO WebPropNotes (webPropID, text, dateAdded) VALUES (?,?,?)"
	updateNote                              = "UPDATE WebPropNotes SET webPropID = ?, text = ?, dateAdded = ? WHERE id =?"
	deleteNote                              = "DELETE FROM WebPropNotes WHERE id = ?"
	deletePropertyNotes                     = "DELETE FROM WebPropNotes WHERE WebPropID = ?"
	createRequirementsBridge                = "INSERT INTO WebPropRequirementCheck (WebPropertiesID, Compliance, WebPropRequirementsID) VALUES (?,?,?)"
	deleteRequirementsBridge                = "DELETE FROM WebPropRequirementCheck WHERE id = ?"
	deleteRequirementsBridgeByRequirementId = "DELETE FROM WebPropRequirementCheck WHERE WebPropRequirementsID = ?"
	deletePropertyRequirementsBridges       = "DELETE FROM WebPropRequirementCheck WHERE WebPropertiesID = ?"
	update                                  = "UPDATE WebProperties SET name = ?, cust_ID = ?,url = ?, isEnabled = ?,sellerID = ?, typeID = ?, isFinalApproved = ?, isEnabledDate = ?, isDenied = ?, requestedDate = ? WHERE id = ?"
	search                                  = `SELECT id, name, cust_ID, badgeID, url, isEnabled,sellerID, typeID , isFinalApproved, isEnabledDate, isDenied, requestedDate, addedDate FROM WebProperties
									 WHERE  name LIKE ? AND cust_ID LIKE ? AND url LIKE ? AND isEnabled LIKE ? AND sellerID LIKE ? AND typeID  LIKE ? AND isFinalApproved LIKE ? AND isEnabledDate LIKE ? AND
									 isDenied LIKE ? AND requestedDate LIKE ? AND addedDate LIKE ? `
	createRequirement = "INSERT INTO WebPropRequirements (ReqType, Requirement) VALUES (?,?)"
	updateRequirement = "UPDATE WebPropRequirements SET ReqType = ?, Requirement = ? WHERE ID = ?"
	deleteRequirement = "DELETE FROM WebPropRequirements WHERE ID = ?"
	getNote           = "SELECT id, webPropID, text, dateAdded FROM WebPropNotes WHERE id = ?"
	getRequirement    = "SELECT ID, ReqType, Requirement FROM WebPropRequirements WHERE ID = ?"
	getType           = "SELECT id, typeID, type FROM WebPropertyTypes WHERE id = ?"
	createType        = "INSERT INTO WebPropertyTypes (typeID, type) VALUES (?,?)"
	updateType        = "UPDATE WebPropertyTypes SET typeID = ?, type=? WHERE id = ?"
	deleteType        = "DELETE FROM WebPropertyTypes WHERE id = ?"
)

const (
	timeFormat = "2006-01-02 15:04:05"
)

// Gets a specific Web Property
func (w *WebProperty) Get(dtx *apicontext.DataContext) error {
	var err error

	redis_key := "webproperty:" + strconv.Itoa(w.ID)
	data, err := redis.Get(redis_key)
	if err == nil && len(data) > 0 {
		err = json.Unmarshal(data, &w)
		if err == nil {
			return err
		}
	}

	err = database.Init()
	if err != nil {
		return err
	}

	stmt, err := database.DB.Prepare(getWebProperty)
	if err != nil {
		return err
	}
	defer stmt.Close()

	webPropTypes, err := GetAllWebPropertyTypes(dtx)
	webPropNotes, err := GetAllWebPropertyNotes(dtx)
	WebPropertyRequirements, err := GetAllWebPropertyRequirements(dtx)
	if err != nil {
		return err
	}

	typesMap := webPropTypes.ToMap()
	notesMap := webPropNotes.ToMap()
	requirementsMap := WebPropertyRequirements.ToMap()

	var url, sid *string
	var tid *int
	err = stmt.QueryRow(w.ID).Scan(
		&w.ID,
		&w.Name,
		&w.CustID,
		&w.BadgeID,
		&url,
		&w.IsEnabled,
		&sid,
		&tid,
		&w.IsFinalApproved,
		&w.IsEnabledDate,
		&w.IsDenied,
		&w.RequestedDate,
		&w.AddedDate,
	)
	if err != nil {
		return err
	}

	if url != nil {
		w.Url = *url
	}
	if sid != nil {
		w.SellerID = *sid
	}
	if tid != nil {
		w.WebPropertyType.ID = *tid
	}

	typeChan := make(chan int)
	notesChan := make(chan int)
	requirementsChan := make(chan int)
	go func() error {
		if _, ok := typesMap[w.WebPropertyType.ID]; ok {
			w.WebPropertyType = typesMap[w.WebPropertyType.ID]
		}
		typeChan <- 1
		return nil
	}()
	go func() error {
		if _, ok := notesMap[w.ID]; ok {
			w.WebPropertyNotes = append(w.WebPropertyNotes, notesMap[w.ID])
		}
		notesChan <- 1
		return nil
	}()
	go func() error {
		if _, ok := requirementsMap[w.ID]; ok {
			w.WebPropertyRequirements = append(w.WebPropertyRequirements, requirementsMap[w.ID])
		}
		requirementsChan <- 1
		return nil
	}()

	<-typeChan
	<-notesChan
	<-requirementsChan

	go redis.Setex(redis_key, w, 86400)
	return err
}

// Gets all of the web properties associated to a specific customer.
func GetByCustomer(CustID int, dtx *apicontext.DataContext) (ws WebProperties, err error) {
	redis_key := "webpropertyByCustomer:" + strconv.Itoa(CustID)
	data, err := redis.Get(redis_key)
	if err == nil && len(data) > 0 {
		err = json.Unmarshal(data, &ws)
		if err == nil {
			return
		}
	}

	err = database.Init()
	if err != nil {
		return
	}

	stmt, err := database.DB.Prepare(getWebPropertiesByCustomer)
	if err != nil {
		return
	}
	defer stmt.Close()

	webPropTypes, err := GetAllWebPropertyTypes(dtx)
	webPropNotes, err := GetAllWebPropertyNotes(dtx)
	WebPropertyRequirements, err := GetAllWebPropertyRequirements(dtx)
	if err != nil {
		return
	}
	typesMap := webPropTypes.ToMap()
	notesMap := webPropNotes.ToMap()
	requirementsMap := WebPropertyRequirements.ToMap()

	res, err := stmt.Query(CustID)
	var w WebProperty
	var url, sid *string
	var tid *int
	for res.Next() {
		err = res.Scan(
			&w.ID,
			&w.Name,
			&w.CustID,
			&w.BadgeID,
			&url,
			&w.IsEnabled,
			&sid,
			&tid,
			&w.IsFinalApproved,
			&w.IsEnabledDate,
			&w.IsDenied,
			&w.RequestedDate,
			&w.AddedDate,
		)
		if err != nil {
			return
		}
		if url != nil {
			w.Url = *url
		}
		if sid != nil {
			w.SellerID = *sid
		}
		if tid != nil {
			w.WebPropertyType.ID = *tid
		}

		typeChan := make(chan int)
		notesChan := make(chan int)
		requirementsChan := make(chan int)
		go func() error {
			if _, ok := typesMap[w.WebPropertyType.ID]; ok {
				w.WebPropertyType = typesMap[w.WebPropertyType.ID]
			}
			typeChan <- 1
			return nil
		}()
		go func() error {
			if _, ok := notesMap[w.ID]; ok {
				w.WebPropertyNotes = append(w.WebPropertyNotes, notesMap[w.ID])
			}
			notesChan <- 1
			return nil
		}()
		go func() error {
			if _, ok := requirementsMap[w.ID]; ok {
				w.WebPropertyRequirements = append(w.WebPropertyRequirements, requirementsMap[w.ID])
			}
			requirementsChan <- 1
			return nil
		}()

		<-typeChan
		<-notesChan
		<-requirementsChan
		ws = append(ws, w)
	}
	defer res.Close()
	go redis.Setex(redis_key, ws, 86400)
	return
}

// Gets All Web Properties
func GetAll(dtx *apicontext.DataContext) (WebProperties, error) {
	var ws WebProperties
	var err error

	redis_key := "webproperties:" + dtx.BrandString
	data, err := redis.Get(redis_key)
	if err == nil && len(data) > 0 {
		err = json.Unmarshal(data, &ws)
		if err == nil {
			return ws, err
		}
	}

	err = database.Init()

	if err != nil {
		return ws, err
	}

	stmt, err := database.DB.Prepare(getAllWebProperties)
	if err != nil {
		return ws, err
	}
	defer stmt.Close()

	webPropTypes, err := GetAllWebPropertyTypes(dtx)
	webPropNotes, err := GetAllWebPropertyNotes(dtx)
	WebPropertyRequirements, err := GetAllWebPropertyRequirements(dtx)
	if err != nil {
		return ws, err
	}

	typesMap := webPropTypes.ToMap()
	notesMap := webPropNotes.ToMap()
	requirementsMap := WebPropertyRequirements.ToMap()

	var url, sid *string
	var tid *int

	res, err := stmt.Query(dtx.APIKey, dtx.BrandID, dtx.BrandID)
	for res.Next() {
		var w WebProperty
		res.Scan(
			&w.ID,
			&w.Name,
			&w.CustID,
			&w.BadgeID,
			&url,
			&w.IsEnabled,
			&sid,
			&tid,
			&w.IsFinalApproved,
			&w.IsEnabledDate,
			&w.IsDenied,
			&w.RequestedDate,
			&w.AddedDate,
		)
		if err != nil {
			return ws, err
		}
		if url != nil {
			w.Url = *url
		}
		if sid != nil {
			w.SellerID = *sid
		}
		if tid != nil {
			w.WebPropertyType.ID = *tid
		}

		typeChan := make(chan int)
		notesChan := make(chan int)
		requirementsChan := make(chan int)
		go func() error {
			if _, ok := typesMap[w.WebPropertyType.ID]; ok {
				w.WebPropertyType = typesMap[w.WebPropertyType.ID]
			}
			typeChan <- 1
			return nil
		}()
		go func() error {
			if _, ok := notesMap[w.ID]; ok {
				w.WebPropertyNotes = append(w.WebPropertyNotes, notesMap[w.ID])
			}
			notesChan <- 1
			return nil
		}()
		go func() error {
			if _, ok := requirementsMap[w.ID]; ok {
				w.WebPropertyRequirements = append(w.WebPropertyRequirements, requirementsMap[w.ID])
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

// Creates a Web Property
func (w *WebProperty) Create(dtx *apicontext.DataContext) (err error) {
	go redis.Delete("webproperties:" + dtx.BrandString)
	err = database.Init()
	if err != nil {
		return err
	}

	tx, err := database.DB.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(create)
	if err != nil {
		return err
	}
	defer stmt.Close()

	add := time.Now()
	w.AddedDate = &add
	res, err := stmt.Exec(
		w.Name,
		w.CustID,
		w.Url,
		w.IsEnabled,
		w.SellerID,
		w.WebPropertyType.ID,
		w.IsFinalApproved,
		w.IsEnabledDate,
		w.IsDenied,
		w.RequestedDate,
		w.AddedDate,
	)
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

	//create/update web properties check
	for _, req := range w.WebPropertyRequirements {
		req.WebPropID = w.ID

		err = req.CreateJoin()
		if err != nil {
			return err
		}
	}
	//create/updated notes
	for _, note := range w.WebPropertyNotes {
		note.WebPropID = w.ID

		err = note.Create(dtx)
		if err != nil {
			return err
		}
	}

	return nil
}

// Updates a Web Property
func (w *WebProperty) Update(dtx *apicontext.DataContext) (err error) {
	go redis.Delete("webproperties:" + dtx.BrandString)

	err = database.Init()
	if err != nil {
		return err
	}

	tx, err := database.DB.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(update)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(w.Name, w.CustID, w.Url, w.IsEnabled, w.SellerID, w.WebPropertyType.ID, w.IsFinalApproved, w.IsEnabledDate, w.IsDenied, w.RequestedDate, w.ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	//create/update web properties check
	for _, req := range w.WebPropertyRequirements {
		req.WebPropID = w.ID
		err = req.DeleteJoin()
		if err != nil {
			return err
		}
		err = req.CreateJoin()
		if err != nil {
			return err
		}
	}
	//create/updated notes
	for _, note := range w.WebPropertyNotes {
		note.WebPropID = w.ID
		if note.ID > 0 {
			err = note.Update(dtx)
			if err != nil {
				return err
			}
		} else {
			err = note.Create(dtx)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Deletes a Web Property and any associations.
func (w *WebProperty) Delete(dtx *apicontext.DataContext) error {
	go redis.Delete("webproperties:" + dtx.BrandString)
	err := database.Init()
	if err != nil {
		return err
	}

	notesChan := make(chan int)
	requirementsChan := make(chan int)

	go func() {
		err = w.DeleteNotesByPropId(dtx)
		notesChan <- 1
	}()
	go func() {
		err = w.DeleteJoinByPropId()
		requirementsChan <- 1
	}()
	<-notesChan
	<-requirementsChan

	tx, err := database.DB.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(deleteWebProp)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(w.ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	return nil
}

// Gets All the available WebPropertyTypes
func GetAllWebPropertyTypes(dtx *apicontext.DataContext) (WebPropertyTypes, error) {
	var ws WebPropertyTypes
	var err error

	redis_key := "webpropertytypes:" + dtx.BrandString
	data, err := redis.Get(redis_key)
	if err == nil && len(data) > 0 {
		err = json.Unmarshal(data, &ws)
		if err == nil {
			return ws, err
		}
	}
	err = database.Init()
	if err != nil {
		return ws, err
	}

	stmt, err := database.DB.Prepare(getAllWebPropertyTypes)
	if err != nil {
		return ws, err
	}
	defer stmt.Close()

	res, err := stmt.Query(dtx.APIKey, dtx.BrandID, dtx.BrandID)
	for res.Next() {
		var w WebPropertyType
		res.Scan(&w.ID, &w.TypeID, &w.Type)
		ws = append(ws, w)
	}
	go redis.Setex(redis_key, ws, 86400)
	return ws, err
}

// gets all of the web property notes - rarely called.
func GetAllWebPropertyNotes(dtx *apicontext.DataContext) (WebPropertyNotes, error) {
	var ws WebPropertyNotes
	var err error

	redis_key := "webpropertynotes:" + dtx.BrandString
	data, err := redis.Get(redis_key)
	if err == nil && len(data) > 0 {
		err = json.Unmarshal(data, &ws)
		if err == nil {
			return ws, err
		}
	}
	err = database.Init()
	if err != nil {
		return ws, err
	}

	stmt, err := database.DB.Prepare(getAllWebPropertyNotes)
	if err != nil {
		return ws, err
	}
	defer stmt.Close()

	res, err := stmt.Query(dtx.APIKey, dtx.BrandID, dtx.BrandID)
	for res.Next() {
		var w WebPropertyNote
		res.Scan(&w.ID, &w.WebPropID, &w.Text, &w.DateAdded)
		ws = append(ws, w)
	}
	go redis.Setex(redis_key, ws, 86400)
	return ws, err
}

// Gets All Web Property Requirements
func GetAllWebPropertyRequirements(dtx *apicontext.DataContext) (WebPropertyRequirements, error) {
	var ws WebPropertyRequirements
	var err error

	redis_key := "webpropertyrequirements:" + dtx.BrandString
	data, err := redis.Get(redis_key)
	if err == nil && len(data) > 0 {
		err = json.Unmarshal(data, &ws)
		if err == nil {
			return ws, err
		}
	}
	err = database.Init()
	if err != nil {
		return ws, err
	}

	stmt, err := database.DB.Prepare(getAllWebPropertyRequirements)
	if err != nil {
		return ws, err
	}
	defer stmt.Close()

	res, err := stmt.Query(dtx.APIKey, dtx.BrandID, dtx.BrandID)
	for res.Next() {
		var w WebPropertyRequirement
		var reqType, req *string
		var comp *bool
		var wpid *int
		err = res.Scan(
			&w.ID,
			&w.RequirementID,
			&reqType,
			&req,
			&comp,
			&wpid,
		)
		if err != nil {
			return ws, err
		}
		if reqType != nil {
			w.ReqType = *reqType
		}
		if req != nil {
			w.Requirement = *req
		}
		if comp != nil {
			w.Compliance = *comp
		}
		if wpid != nil {
			w.WebPropID = *wpid
		}
		ws = append(ws, w)
	}
	go redis.Setex(redis_key, ws, 86400)
	return ws, err
}

// gets a specific Web Property Note
func (n *WebPropertyNote) Get() error {
	redis_key := "webpropertynote:" + strconv.Itoa(n.ID)
	data, err := redis.Get(redis_key)
	if err == nil && len(data) > 0 {
		err = json.Unmarshal(data, &n)
		return err
	}
	err = database.Init()
	if err != nil {
		return err
	}

	stmt, err := database.DB.Prepare(getNote)
	if err != nil {
		return err
	}
	defer stmt.Close()
	err = stmt.QueryRow(n.ID).Scan(&n.ID, &n.WebPropID, &n.Text, &n.DateAdded)
	if err != nil {
		return err
	}

	return nil
}

// Creates a new note for a web property.
func (n *WebPropertyNote) Create(dtx *apicontext.DataContext) error {
	go redis.Delete("webpropertynotes:" + dtx.BrandString)
	err := database.Init()
	if err != nil {
		return err
	}

	tx, err := database.DB.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(createNote)
	if err != nil {
		return err
	}
	defer stmt.Close()
	da := time.Now()
	n.DateAdded = &da
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

// Updates a note on a web property
func (n *WebPropertyNote) Update(dtx *apicontext.DataContext) error {
	go redis.Delete("webpropertynotes:" + dtx.BrandString)
	err := database.Init()
	if err != nil {
		return err
	}

	tx, err := database.DB.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(updateNote)
	if err != nil {
		return err
	}
	defer stmt.Close()
	da := time.Now()
	n.DateAdded = &da
	_, err = stmt.Exec(n.WebPropID, n.Text, n.DateAdded, n.ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	return nil
}

// Deletes a Web Property Note
func (n *WebPropertyNote) Delete(dtx *apicontext.DataContext) error {
	go redis.Delete("webpropertynotes:" + dtx.BrandString)
	err := database.Init()
	if err != nil {
		return err
	}

	tx, err := database.DB.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(deleteNote)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(n.ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	return nil
}

// Deletes all of the notes of a specific Web Property
func (n *WebProperty) DeleteNotesByPropId(dtx *apicontext.DataContext) error {
	go redis.Delete("webpropertynotes:" + dtx.BrandString)
	err := database.Init()
	if err != nil {
		return err
	}

	tx, err := database.DB.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(deletePropertyNotes)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(n.ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	return nil
}

// Makes an association between a Web Property and a Requirement
func (r *WebPropertyRequirement) CreateJoin() error {
	err := database.Init()
	if err != nil {
		return err
	}

	tx, err := database.DB.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(createRequirementsBridge)
	if err != nil {
		return err
	}
	defer stmt.Close()
	res, err := stmt.Exec(r.WebPropID, r.Compliance, r.RequirementID)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	id, err := res.LastInsertId()
	r.ID = int(id)
	return nil
}

// removes an association between a Web Property and a Requirement
func (r *WebPropertyRequirement) DeleteJoin() error {
	err := database.Init()
	if err != nil {
		return err
	}

	tx, err := database.DB.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(deleteRequirementsBridge)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(r.ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	return nil
}

// removes a web property requirement association by its own requirementID
func (r *WebPropertyRequirement) DeleteJoinByRequirementId() error {
	err := database.Init()
	if err != nil {
		return err
	}

	tx, err := database.DB.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(deleteRequirementsBridgeByRequirementId)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(r.RequirementID)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	return nil
}

// removes all associations between a Web Property and it's Requirements
func (r *WebProperty) DeleteJoinByPropId() error {
	err := database.Init()
	if err != nil {
		return err
	}

	tx, err := database.DB.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(deletePropertyRequirementsBridges)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(r.ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	return nil
}

// Gets a specific WebPropertyRequirement
func (r *WebPropertyRequirement) Get() error {
	redis_key := "webpropertyrequirement:" + strconv.Itoa(r.RequirementID)
	data, err := redis.Get(redis_key)
	if err == nil && len(data) > 0 {
		err = json.Unmarshal(data, &r)
		return err
	}
	err = database.Init()
	if err != nil {
		return err
	}

	stmt, err := database.DB.Prepare(getRequirement)
	if err != nil {
		return err
	}
	defer stmt.Close()
	var req, reqType *string
	err = stmt.QueryRow(r.RequirementID).Scan(
		&r.ID,
		&reqType,
		&req,
	)
	if err != nil {
		return err
	}
	if reqType != nil {
		r.ReqType = *reqType
	}
	if req != nil {
		r.Requirement = *req
	}

	return nil
}

// Creates a WebPropertyRequirement
func (r *WebPropertyRequirement) Create(dtx *apicontext.DataContext) error {
	go redis.Delete("webpropertyrequirements:" + dtx.BrandString)
	err := database.Init()
	if err != nil {
		return err
	}

	tx, err := database.DB.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(createRequirement)
	if err != nil {
		return err
	}
	defer stmt.Close()
	res, err := stmt.Exec(r.ReqType, r.Requirement)
	if err != nil {
		tx.Rollback()
		return err
	}
	id, err := res.LastInsertId()

	r.RequirementID = int(id)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	return nil
}

// Updates a WebPropertyRequirement
func (r *WebPropertyRequirement) Update(dtx *apicontext.DataContext) error {
	go redis.Delete("webpropertyrequirements:" + dtx.BrandString)
	err := database.Init()
	if err != nil {
		return err
	}

	tx, err := database.DB.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(updateRequirement)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(r.ReqType, r.Requirement, r.RequirementID)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	return nil
}

// Deletes a Web Property Requirement
func (r *WebPropertyRequirement) Delete(dtx *apicontext.DataContext) error {
	go redis.Delete("webpropertyrequirements:" + dtx.BrandString)
	var err error
	err = r.Get()
	if err != nil {
		return err
	}
	err = r.DeleteJoinByRequirementId()
	if err != nil {
		return err
	}
	err = database.Init()
	if err != nil {
		return err
	}

	tx, err := database.DB.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(deleteRequirement)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(r.RequirementID)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	return nil
}

// Gets a WebPropertyType
func (t *WebPropertyType) Get() error {
	redis_key := "webpropertytype:" + strconv.Itoa(t.ID)
	data, err := redis.Get(redis_key)
	if err == nil && len(data) > 0 {
		err = json.Unmarshal(data, &t)
		return err
	}
	err = database.Init()
	if err != nil {
		return err
	}

	stmt, err := database.DB.Prepare(getType)
	if err != nil {
		return err
	}
	defer stmt.Close()
	err = stmt.QueryRow(t.ID).Scan(&t.ID, &t.TypeID, &t.Type)
	if err != nil {
		return err
	}

	return nil
}

// Updates a WebPropertyType
func (t *WebPropertyType) Update(dtx *apicontext.DataContext) error {
	go redis.Delete("webpropertytypes:" + dtx.BrandString)
	err := database.Init()
	if err != nil {
		return err
	}

	tx, err := database.DB.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(updateType)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(t.TypeID, t.Type, t.ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	return nil
}

// creates a WebPropertyType
func (t *WebPropertyType) Create(dtx *apicontext.DataContext) error {
	go redis.Delete("webpropertytypes:" + dtx.BrandString)
	err := database.Init()
	if err != nil {
		return err
	}

	tx, err := database.DB.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(createType)
	if err != nil {
		return err
	}
	defer stmt.Close()
	res, err := stmt.Exec(t.TypeID, t.Type)
	if err != nil {
		tx.Rollback()
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	t.ID = int(id)
	tx.Commit()
	return nil
}

// Deletes a WebPropertyType
func (t *WebPropertyType) Delete(dtx *apicontext.DataContext) error {
	go redis.Delete("webpropertytypes:" + dtx.BrandString)
	err := database.Init()
	if err != nil {
		return err
	}

	tx, err := database.DB.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(deleteType)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(t.ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	return nil
}

// Searches for a web property given all the web properties properties as search parameters.
func Search(name, custID, badgeID, url, isEnabled, sellerID, webPropertyTypeID, isFinalApproved, isEnabledDate, isDenied, requestedDate, typeID, pageStr, resultsStr string) (pagination.Objects, error) {
	var err error
	var l pagination.Objects
	var fs []interface{}

	err = database.Init()
	if err != nil {
		return l, err
	}

	stmt, err := database.DB.Prepare(search)
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
