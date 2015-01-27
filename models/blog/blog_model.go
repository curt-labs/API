package blog_model

import (
	"database/sql"
	"encoding/json"
	"github.com/curt-labs/GoAPI/helpers/apicontext"
	"github.com/curt-labs/GoAPI/helpers/database"
	"github.com/curt-labs/GoAPI/helpers/pagination"
	"github.com/curt-labs/GoAPI/helpers/redis"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
	"time"
)

type Blog struct {
	ID              int            `json:"id,omitempty" xml:"id,omitempty"`
	Title           string         `json:"title,omitempty" xml:"title,omitempty"`
	Slug            string         `json:"slug,omitempty" xml:"slug,omitempty"`
	Text            string         `json:"text,omitempty" xml:"text,omitempty"`
	PublishedDate   time.Time      `json:"publishedDate,omitempty" xml:"publishedDate,omitempty"`
	CreatedDate     time.Time      `json:"createdDate,omitempty" xml:"createdDate,omitempty"`
	LastModified    time.Time      `json:"lastModified,omitempty" xml:"lastModified,omitempty"`
	UserID          int            `json:"userID,omitempty" xml:"userID,omitempty"`
	MetaTitle       string         `json:"metaTitle,omitempty" xml:"metaTitle,omitempty"`
	MetaDescription string         `json:"metaDescription,omitempty" xml:"metaDescription,omitempty"`
	Keywords        string         `json:"keywords,omitempty" xml:"keywords,omitempty"`
	Active          bool           `json:"active,omitempty" xml:"active,omitempty"`
	BlogCategories  BlogCategories `json:"blogCategories,omitempty" xml:"blogCategories,omitempty"`
}

type Blogs []Blog
type Categories []Category

type Category struct {
	ID     int    `json:"id,omitempty" xml:"id,omitempty"`
	Name   string `json:"name,omitempty" xml:"name,omitempty"`
	Slug   string `json:"slug,omitempty" xml:"slug,omitempty"`
	Active bool   `json:"active,omitempty" xml:"active,omitempty"`
}

type BlogCategory struct {
	ID             int      `json:"id,omitempty" xml:"id,omitempty"`
	BlogPostID     int      `json:"blogPostID,omitempty" xml:"blogPostID,omitempty"`
	BlogCategoryID int      `json:"blogCategoryID,omitempty" xml:"blogCategoryID,omitempty"`
	Category       Category `json:"category,omitempty" xml:"category,omitempty"`
}

type BlogCategories []BlogCategory

var (
	getAllBlogs = `SELECT b.blogPostID, b.post_title ,b.slug ,b.post_text ,b.publishedDate, b.createdDate, b.lastModified, b.userID, b.meta_title, b.meta_description, b.keywords, b.active FROM  BlogPosts AS b 
									Join BlogPost_BlogCategory as bpbc on bpbc.blogPostID = b.blogPostID
									Join BlogCategories as bc on bc.blogCategoryID = bpbc.blogCategoryID
									Join ApiKeyToBrand as akb on akb.brandID = bc.brandID
									Join ApiKey as ak on akb.keyID = ak.id
									where (ak.api_key = ? && (bc.brandID = ? OR 0=?))`
	getAllCategories = `SELECT bc.blogCategoryID, bc.name, bc.slug, bc.active FROM BlogCategories AS bc
									Join ApiKeyToBrand as akb on akb.brandID = bc.brandID
									Join ApiKey as ak on akb.keyID = ak.id
									where (ak.api_key = ? && (bc.brandID = ? OR 0=?))`

	stmtGetAllBlogCategories = `SELECT bpbc.postCategoryID, bpbc.blogPostID, bpbc.blogCategoryID, bc.blogCategoryID, bc.name, bc.slug, bc.active FROM BlogPost_BlogCategory AS bpbc 
									LEFT JOIN blogCategories AS bc ON bc.blogCategoryID = bpbc.blogCategoryID
									Join ApiKeyToBrand as akb on akb.brandID = bc.brandID
									Join ApiKey as ak on akb.keyID = ak.id
									where (ak.api_key = ? && (bc.brandID = ? OR 0=?))`

	getBlog = `SELECT b.blogPostID, b.post_title ,b.slug, COALESCE(b.post_text,'') ,COALESCE(b.publishedDate,''), COALESCE(b.createdDate,''), COALESCE(b.lastModified,''), b.userID, COALESCE(b.meta_title,''), COALESCE(b.meta_description,''), COALESCE(b.keywords,''), b.active 
					FROM BlogPosts AS b 
						Join BlogPost_BlogCategory as bpbc on bpbc.blogPostID = b.blogPostID
						Join BlogCategories as bc on bc.blogCategoryID = bpbc.blogCategoryID
						Join ApiKeyToBrand as akb on akb.brandID = bc.brandID
						Join ApiKey as ak on akb.keyID = ak.id
						where (ak.api_key = ? && (bc.brandID = ? OR 0=?)) && b.blogPostID = ?`
	create      = `INSERT INTO BlogPosts (post_title ,slug ,post_text, createdDate, publishedDate, lastModified, userID, meta_title, meta_description, keywords, active, thumbnail) VALUES (?,?,?,?,?,?,?,?,?,?,?,?)`
	getCategory = `SELECT bc.blogCategoryID, bc.name, bc.slug,bc.active FROM BlogCategories as bc 
									Join ApiKeyToBrand as akb on akb.brandID = bc.brandID
									Join ApiKey as ak on akb.keyID = ak.id
									where (ak.api_key = ? && (bc.brandID = ? OR 0=?)) && bc.blogCategoryID = ?`
	createCategory         = `INSERT INTO BlogCategories (name,slug,active, brandID) VALUES (?,?,?,?)`
	deleteCategory         = `DELETE FROM BlogCategories WHERE blogCategoryID = ?`
	createCatBridge        = `INSERT INTO BlogPost_BlogCategory (blogPostID, blogCategoryID) VALUES (?,?)`
	deleteCatBridge        = `DELETE FROM BlogPost_BlogCategory WHERE blogPostID = ?`
	deleteCatBridgeByCatID = `DELETE FROM BlogPost_BlogCategory WHERE blogCategoryID = ?`
	update                 = `UPDATE BlogPosts SET post_title = ? ,slug = ? ,post_text = ?, publishedDate= ?, lastModified = ?,userID = ?, meta_title = ?, meta_description = ?, keywords = ?, active = ? WHERE blogPostID = ?`
	deleteBlog             = "DELETE FROM BlogPosts WHERE blogPostID = ?"
	search                 = `SELECT b.blogPostID, b.post_title ,b.slug ,b.post_text ,b.publishedDate, b.createdDate, b.lastModified, b.userID, b.meta_title, b.meta_description, b.keywords, b.active FROM BlogPosts AS b
									Join BlogPost_BlogCategory as bpbc on bpbc.blogPostID = b.blogPostID
									Join BlogCategories as bc on bc.blogCategoryID = bpbc.blogCategoryID
									Join ApiKeyToBrand as akb on akb.brandID = bc.brandID
									Join ApiKey as ak on akb.keyID = ak.id
									where (ak.api_key = ? && (bc.brandID = ? OR 0=?)) && b.post_title LIKE ? AND b.slug LIKE ? AND b.post_text LIKE ? AND b.publishedDate LIKE ? AND b.createdDate LIKE ? AND b.lastModified LIKE ? AND b.userID LIKE ? AND b.meta_title LIKE ? AND b.meta_description LIKE ? AND b.keywords LIKE ? AND b.active LIKE ?`
)

const (
	timeFormat = "2006-01-02 15:04:05"
)

func GetAll(dtx *apicontext.DataContext) (Blogs, error) {
	var bs Blogs
	var err error

	redis_key := "blogs:" + dtx.BrandString
	data, err := redis.Get(redis_key)
	if err == nil && len(data) > 0 {
		err = json.Unmarshal(data, &bs)
		return bs, err
	}

	db, err := sql.Open("mysql", database.ConnectionString())

	if err != nil {
		return bs, err
	}
	defer db.Close()

	stmt, err := db.Prepare(getAllBlogs)
	if err != nil {
		return bs, err
	}
	defer stmt.Close()
	bcs, err := getAllBlogCategories(dtx)
	bcMap := bcs.ToMap()

	res, err := stmt.Query(dtx.APIKey, dtx.BrandID, dtx.BrandID)
	for res.Next() {
		var b Blog
		res.Scan(&b.ID, &b.Title, &b.Slug, &b.Text, &b.PublishedDate, &b.CreatedDate, &b.LastModified, &b.UserID, &b.MetaTitle, &b.MetaDescription, &b.Keywords, &b.Active)
		bcChan := make(chan int)

		go func() {
			for _, val := range bcMap {
				if val.BlogPostID == b.ID {
					b.BlogCategories = append(b.BlogCategories, val)
				}
			}
			bcChan <- 1
		}()
		<-bcChan
		bs = append(bs, b)
	}
	defer res.Close()
	go redis.Setex(redis_key, bs, 86400)
	return bs, err
}

func getAllBlogCategories(dtx *apicontext.DataContext) (BlogCategories, error) {
	var bcs BlogCategories
	var err error
	db, err := sql.Open("mysql", database.ConnectionString())

	if err != nil {
		return bcs, err
	}
	defer db.Close()

	stmt, err := db.Prepare(stmtGetAllBlogCategories)
	if err != nil {
		return bcs, err
	}
	defer stmt.Close()
	res, err := stmt.Query(dtx.APIKey, dtx.BrandID, dtx.BrandID)
	for res.Next() {
		var temp BlogCategory
		res.Scan(&temp.ID, &temp.BlogPostID, &temp.BlogCategoryID, &temp.Category.ID, &temp.Category.Name, &temp.Category.Slug, &temp.Category.Active)
		bcs = append(bcs, temp)
	}
	defer res.Close()
	return bcs, err
}

func GetAllCategories(dtx *apicontext.DataContext) (Categories, error) {
	var cs Categories
	var err error
	redis_key := "blogs:categories:" + dtx.BrandString
	data, err := redis.Get(redis_key)
	if err == nil && len(data) > 0 {
		err = json.Unmarshal(data, &cs)
		return cs, err
	}

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return cs, err
	}
	defer db.Close()

	stmt, err := db.Prepare(getAllCategories)
	if err != nil {
		return cs, err
	}
	defer stmt.Close()
	res, err := stmt.Query(dtx.APIKey, dtx.BrandID, dtx.BrandID)
	for res.Next() {
		var c Category
		res.Scan(&c.ID, &c.Name, &c.Slug, &c.Active)
		cs = append(cs, c)
	}
	defer res.Close()
	go redis.Setex(redis_key, cs, 86400)

	return cs, err
}

func (b *Blog) Get(dtx *apicontext.DataContext) error {
	var err error

	redis_key := "blog:" + strconv.Itoa(b.ID) + ":" + dtx.BrandString
	data, err := redis.Get(redis_key)
	if err == nil && len(data) > 0 {
		err = json.Unmarshal(data, &b)
		return err
	}

	db, err := sql.Open("mysql", database.ConnectionString())

	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(getBlog)
	if err != nil {
		return err
	}
	defer stmt.Close()
	bcs, err := getAllBlogCategories(dtx)
	bcMap := bcs.ToMap()
	var p, c, l *string

	res, err := stmt.Query(dtx.APIKey, dtx.BrandID, dtx.BrandID, b.ID)
	for res.Next() {

		res.Scan(&b.ID, &b.Title, &b.Slug, &b.Text, &p, &c, &l, &b.UserID, &b.MetaTitle, &b.MetaDescription, &b.Keywords, &b.Active)
		if err != nil {
			return err
		}
		if b.ID == 0 {
			continue
		}

		if p != nil {
			b.PublishedDate, _ = time.Parse(timeFormat, *p)
		}
		if c != nil {
			b.CreatedDate, _ = time.Parse(timeFormat, *c)
		}
		if l != nil {
			b.LastModified, _ = time.Parse(timeFormat, *l)

		}

		bcChan := make(chan int)

		go func() {
			var tempBlogCat []BlogCategory
			for _, val := range bcMap {
				if val.BlogPostID == b.ID {
					tempBlogCat = append(tempBlogCat, val)
				}
			}
			b.BlogCategories = tempBlogCat
			bcChan <- 1
		}()
		<-bcChan
	}
	defer res.Close()
	go redis.Setex(redis_key, b, 86400)
	return err
}

func (b *Blog) Create(dtx *apicontext.DataContext) error {
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
	b.LastModified = time.Now()
	b.CreatedDate = time.Now()
	res, err := stmt.Exec(b.Title, b.Slug, b.Text, b.CreatedDate, b.PublishedDate, b.LastModified, b.UserID, b.MetaTitle, b.MetaDescription, b.Keywords, b.Active, "")
	if err != nil {
		tx.Rollback()
		return err
	}
	id, err := res.LastInsertId()
	b.ID = int(id)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	b.createCatBridge()
	err = redis.Setex("blog:"+strconv.Itoa(b.ID)+":"+dtx.BrandString, b, 86400)
	return nil
}

func (b *Blog) Delete(dtx *apicontext.DataContext) error {
	var err error
	err = b.deleteCatBridge()
	if err != nil {
		return err
	}
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(deleteBlog)
	_, err = stmt.Exec(b.ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	if err == nil {
		redis.Delete("blog:" + strconv.Itoa(b.ID) + ":" + dtx.BrandString)
	}

	return err
}

func (b *Blog) createCatBridge() error {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	for _, v := range b.BlogCategories {
		tx, err := db.Begin()
		if err != nil {
			return err
		}

		stmt, err := tx.Prepare(createCatBridge)
		_, err = stmt.Exec(b.ID, v.Category.ID)
		if err != nil {
			tx.Rollback()
			return err
		}
		if err != nil {
			tx.Rollback()
			return err
		}
		tx.Commit()
	}
	return nil
}

func (b *Blog) Update(dtx *apicontext.DataContext) error {
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
	b.LastModified = time.Now()
	_, err = stmt.Exec(b.Title, b.Slug, b.Text, b.PublishedDate, b.LastModified, b.UserID, b.MetaTitle, b.MetaDescription, b.Keywords, b.Active, b.ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	bcChan := make(chan int)
	createCatChan := make(chan int)
	go func() {
		err = b.deleteCatBridge()
		bcChan <- 1
	}()
	<-bcChan //need these synchrnous, I guess
	go func() {
		err = b.createCatBridge()
		createCatChan <- 1
	}()
	<-createCatChan
	err = redis.Setex("blog:"+strconv.Itoa(b.ID)+":"+dtx.BrandString, b, 86400)
	return nil
}

func (b *Blog) deleteCatBridge() error {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(deleteCatBridge)
	_, err = stmt.Exec(b.ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	return nil
}

func (c *Category) Get(dtx *apicontext.DataContext) error {
	var err error
	redis_key := "blogs:category:" + strconv.Itoa(c.ID) + ":" + dtx.BrandString
	data, err := redis.Get(redis_key)
	if err == nil && len(data) > 0 {
		err = json.Unmarshal(data, &c)
		return err
	}

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(getCategory)
	if err != nil {
		return err
	}
	defer stmt.Close()
	err = stmt.QueryRow(dtx.APIKey, dtx.BrandID, dtx.BrandID, c.ID).Scan(&c.ID, &c.Name, &c.Slug, &c.Active)

	go redis.Setex(redis_key, c, 86400)
	return err
}

func (c *Category) Create(dtx *apicontext.DataContext) error {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(createCategory)

	res, err := stmt.Exec(c.Name, c.Slug, c.Active, dtx.BrandID)
	if err != nil {
		tx.Rollback()
		return err
	}
	id, err := res.LastInsertId()
	c.ID = int(id)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	redis.Setex("blogs:category:"+strconv.Itoa(c.ID)+":"+dtx.BrandString, c, redis.CacheTimeout)

	return nil
}

func (c *Category) Delete(dtx *apicontext.DataContext) error {
	var err error
	err = c.deleteCatBridgeByCategory()
	if err != nil {
		return err
	}

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(deleteCategory)

	_, err = stmt.Exec(c.ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	if err == nil {
		redis.Delete("blogs:category:" + strconv.Itoa(c.ID) + ":" + dtx.BrandString)
	}

	return err
}

func (c *Category) deleteCatBridgeByCategory() error {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(deleteCatBridgeByCatID)
	_, err = stmt.Exec(c.ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	return nil
}

func Search(title, slug, text, publishedDate, createdDate, lastModified, userID, metaTitle, metaDescription, keywords, active, pageStr, resultsStr string, dtx *apicontext.DataContext) (pagination.Objects, error) {
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

	res, err := stmt.Query(dtx.APIKey, dtx.BrandID, dtx.BrandID, "%"+title+"%", "%"+slug+"%", "%"+text+"%", "%"+publishedDate+"%", "%"+createdDate+"%", "%"+lastModified+"%", "%"+userID+"%", "%"+metaTitle+"%", "%"+metaDescription+"%", "%"+keywords+"%", "%"+active+"%")
	for res.Next() {
		var n Blog
		res.Scan(&n.ID, &n.Title, &n.Slug, &n.Text, &n.PublishedDate, &n.CreatedDate, &n.LastModified, &n.UserID, &n.MetaTitle, &n.MetaDescription, &n.Keywords, &n.Active)
		fs = append(fs, n)
	}
	l = pagination.Paginate(pageStr, resultsStr, fs)
	return l, err
}
