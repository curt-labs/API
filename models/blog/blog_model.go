package blog_model

import (
	"database/sql"
	"encoding/json"
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
	getAllBlogs              = "SELECT b.blogPostID, b.post_title ,b.slug ,b.post_text ,b.publishedDate, b.createdDate, b.lastModified, b.userID, b.meta_title, b.meta_description, b.keywords, b.active FROM  BlogPosts AS b "
	getAllCategories         = "SELECT b.blogCategoryID, b.name, b.slug, b.active FROM BlogCategories AS b"
	stmtGetAllBlogCategories = "SELECT bc.postCategoryID, bc.blogPostID, bc.blogCategoryID, b.blogCategoryID, b.name, b.slug, b.active FROM BlogPost_BlogCategory AS bc LEFT JOIN blogCategories AS b ON b.blogCategoryID = bc.blogCategoryID"
	getBlog                  = "SELECT b.blogPostID, b.post_title ,b.slug, COALESCE(b.post_text,'') ,COALESCE(b.publishedDate,''), COALESCE(b.createdDate,''), COALESCE(b.lastModified,''), b.userID, COALESCE(b.meta_title,''), COALESCE(b.meta_description,''), COALESCE(b.keywords,''), b.active FROM  BlogPosts AS b WHERE b.blogPostID = ?"
	create                   = `INSERT INTO BlogPosts (post_title ,slug ,post_text, createdDate, publishedDate, lastModified, userID, meta_title, meta_description, keywords, active)
								VALUES (?,?,?,?,?,?,?,?,?,?,?)`
	getCategory     = "SELECT blogCategoryID, name, slug,active FROM BlogCategories WHERE blogCategoryID = ?"
	createCategory  = "INSERT INTO BlogCategories (name,slug,active) VALUES (?,?,?)"
	deleteCategory  = "DELETE FROM BlogCategories WHERE blogCategoryID = ?"
	createCatBridge = `INSERT INTO BlogPost_BlogCategory (blogPostID, blogCategoryID) VALUES (?,?)`
	deleteCatBridge = `DELETE FROM BlogPost_BlogCategory WHERE blogPostID = ?`
	update          = `UPDATE BlogPosts SET post_title = ? ,slug = ? ,post_text = ?, publishedDate= ?, lastModified = ?,userID = ?, meta_title = ?, meta_description = ?, keywords = ?, active = ? WHERE blogPostID = ?`
	deleteBlog      = "DELETE FROM BlogPosts WHERE blogPostID = ?"
	search          = `SELECT b.blogPostID, b.post_title ,b.slug ,b.post_text ,b.publishedDate, b.createdDate, b.lastModified, b.userID, b.meta_title, b.meta_description, b.keywords, b.active FROM  BlogPosts AS b
						WHERE b.post_title LIKE ? AND b.slug LIKE ? AND b.post_text LIKE ? AND b.publishedDate LIKE ? AND b.createdDate LIKE ? AND b.lastModified LIKE ? AND b.userID LIKE ? AND b.meta_title LIKE ? AND b.meta_description LIKE ? AND b.keywords LIKE ? AND b.active LIKE ?`
)

const (
	timeFormat = "2006-01-02 15:04:05"
)

func GetAll() (Blogs, error) {
	var bs Blogs
	var err error

	redis_key := "blogs"
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
	bcs, err := getAllBlogCategories()
	bcMap := bcs.ToMap()

	res, err := stmt.Query()
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
	go redis.Setex(redis_key, bs, 86400)
	return bs, err
}

func getAllBlogCategories() (BlogCategories, error) {
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
	res, err := stmt.Query()
	for res.Next() {
		var temp BlogCategory
		res.Scan(&temp.ID, &temp.BlogPostID, &temp.BlogCategoryID, &temp.Category.ID, &temp.Category.Name, &temp.Category.Slug, &temp.Category.Active)
		bcs = append(bcs, temp)
	}
	return bcs, err

}

func GetAllCategories() (Categories, error) {
	var cs Categories
	var err error
	redis_key := "blogs:categories"
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
	res, err := stmt.Query()
	for res.Next() {
		var c Category
		res.Scan(&c.ID, &c.Name, &c.Slug, &c.Active)
		cs = append(cs, c)
	}
	go redis.Setex(redis_key, cs, 86400)

	return cs, err
}
func (c *Category) Get() error {
	var err error
	redis_key := "blogs:category" + strconv.Itoa(c.ID)
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
	err = stmt.QueryRow(c.ID).Scan(&c.ID, &c.Name, &c.Slug, &c.Active)

	go redis.Setex(redis_key, c, 86400)
	return err
}

func (b *Blog) Get() error {
	var err error

	redis_key := "blogs:" + strconv.Itoa(b.ID)
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
	bcs, err := getAllBlogCategories()
	bcMap := bcs.ToMap()
	var p, c, l *string

	res, err := stmt.Query(b.ID)
	for res.Next() {

		res.Scan(&b.ID, &b.Title, &b.Slug, &b.Text, &p, &c, &l, &b.UserID, &b.MetaTitle, &b.MetaDescription, &b.Keywords, &b.Active)
		if err != nil {
			return err
		}
		if p != nil {
			b.PublishedDate, err = time.Parse(timeFormat, *p)
		}
		if c != nil {
			b.CreatedDate, err = time.Parse(timeFormat, *c)
		}
		if l != nil {
			b.LastModified, err = time.Parse(timeFormat, *l)

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
	go redis.Setex(redis_key, b, 86400)
	return err
}

func (b *Blog) Create() error {
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
	res, err := stmt.Exec(b.Title, b.Slug, b.Text, b.CreatedDate, b.PublishedDate, b.LastModified, b.UserID, b.MetaTitle, b.MetaDescription, b.Keywords, b.Active)
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
	return nil
}

func (c *Category) Create() error {
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

	res, err := stmt.Exec(c.Name, c.Slug, c.Active)
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
	tx.Commit()
	return nil
}

func (c *Category) Delete() error {
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
	return nil
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

//implement redis caching - redis get - else db - redis.Setex (see message)
func (b *Blog) Update() error {
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

func (b *Blog) Delete() error {
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
	return nil
}

func Search(title, slug, text, publishedDate, createdDate, lastModified, userID, metaTitle, metaDescription, keywords, active, pageStr, resultsStr string) (pagination.Objects, error) {
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

	res, err := stmt.Query("%"+title+"%", "%"+slug+"%", "%"+text+"%", "%"+publishedDate+"%", "%"+createdDate+"%", "%"+lastModified+"%", "%"+userID+"%", "%"+metaTitle+"%", "%"+metaDescription+"%", "%"+keywords+"%", "%"+active+"%")
	for res.Next() {
		var n Blog
		res.Scan(&n.ID, &n.Title, &n.Slug, &n.Text, &n.PublishedDate, &n.CreatedDate, &n.LastModified, &n.UserID, &n.MetaTitle, &n.MetaDescription, &n.Keywords, &n.Active)
		fs = append(fs, n)
	}
	l = pagination.Paginate(pageStr, resultsStr, fs)
	return l, err
}
