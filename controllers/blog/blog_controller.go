package blog_controller

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/curt-labs/API/helpers/apicontext"
	"github.com/curt-labs/API/helpers/encoding"
	"github.com/curt-labs/API/helpers/error"
	"github.com/curt-labs/API/helpers/sortutil"
	"github.com/curt-labs/API/models/blog"
	"github.com/go-martini/martini"
)

const (
	timeFormat = "2006-01-02 15:04:05"
)

func GetAll(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder, dtx *apicontext.DataContext) string {
	blogs, err := blog_model.GetAll(dtx)
	if err != nil {
		apierror.GenerateError("Trouble getting all blogs", err, rw, r)
	}
	sort := r.FormValue("sort")
	direction := r.FormValue("direction")
	if sort != "" {
		if strings.ContainsAny(direction, "esc") {
			sortutil.DescByField(blogs, sort)
		} else {
			sortutil.AscByField(blogs, sort)
		}

	}
	return encoding.Must(enc.Encode(blogs))
}

func GetAllCategories(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder, dtx *apicontext.DataContext) string {
	cats, err := blog_model.GetAllCategories(dtx)
	if err != nil {
		apierror.GenerateError("Trouble getting blog categories", err, rw, r)
	}
	sort := r.FormValue("sort")
	direction := r.FormValue("direction")
	if sort != "" {
		if strings.ContainsAny(direction, "esc") {
			sortutil.DescByField(cats, sort)
		} else {
			sortutil.AscByField(cats, sort)
		}

	}
	return encoding.Must(enc.Encode(cats))
}

func GetBlog(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params, dtx *apicontext.DataContext) string {
	var b blog_model.Blog
	var err error
	idStr := r.FormValue("id")
	if idStr != "" {
		b.ID, err = strconv.Atoi(idStr)
		if err != nil {
			apierror.GenerateError("Trouble getting blog ID", err, rw, r)
		}
	} else {
		b.ID, err = strconv.Atoi(params["id"])
		if err != nil {
			apierror.GenerateError("Trouble getting blog ID", err, rw, r)
		}
	}
	err = b.Get(dtx)
	if err != nil {
		apierror.GenerateError("Trouble getting blog", err, rw, r)
	}
	return encoding.Must(enc.Encode(b))
}

func CreateBlog(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder, dtx *apicontext.DataContext) string {
	var b blog_model.Blog
	var err error

	b.Title = r.FormValue("title")
	b.Slug = r.FormValue("slug")
	b.Text = r.FormValue("text")
	b.PublishedDate, err = time.Parse(timeFormat, r.FormValue("publishedDate"))
	b.UserID, err = strconv.Atoi(r.FormValue("userID"))
	b.MetaTitle = r.FormValue("metaTitle")
	b.MetaDescription = r.FormValue("metaDescription")
	b.Keywords = r.FormValue("keywords")
	b.Active, err = strconv.ParseBool(r.FormValue("active"))
	categoryIDs := r.Form["categoryID"]
	for _, v := range categoryIDs {
		var bc blog_model.BlogCategory
		bc.Category.ID, err = strconv.Atoi(v)
		b.BlogCategories = append(b.BlogCategories, bc)
	}

	err = b.Create(dtx)
	if err != nil {
		apierror.GenerateError("Trouble creating blog", err, rw, r)
	}
	return encoding.Must(enc.Encode(b))
}
func GetBlogCategory(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params, dtx *apicontext.DataContext) string {
	var c blog_model.Category
	var err error
	idStr := r.FormValue("id")
	if idStr != "" {
		c.ID, err = strconv.Atoi(idStr)
		if err != nil {
			apierror.GenerateError("Trouble getting blog category ID", err, rw, r)
		}
	} else {
		c.ID, err = strconv.Atoi(params["id"])
		if err != nil {
			apierror.GenerateError("Trouble getting blog category ID", err, rw, r)
		}
	}
	err = c.Get(dtx)
	if err != nil {
		apierror.GenerateError("Trouble getting blog category", err, rw, r)
	}
	return encoding.Must(enc.Encode(c))
}
func CreateBlogCategory(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder, dtx *apicontext.DataContext) string {
	var c blog_model.Category
	var err error

	c.Name = r.FormValue("name")
	c.Slug = r.FormValue("slug")
	c.Active, err = strconv.ParseBool(r.FormValue("active"))

	err = c.Create(dtx)
	if err != nil {
		apierror.GenerateError("Trouble creating blog category", err, rw, r)
	}
	return encoding.Must(enc.Encode(c))
}

func DeleteBlogCategory(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params, dtx *apicontext.DataContext) string {
	var c blog_model.Category
	var err error

	idStr := r.FormValue("id")
	if idStr != "" {
		c.ID, err = strconv.Atoi(idStr)
		if err != nil {
			apierror.GenerateError("Trouble getting blog category ID", err, rw, r)
		}
	} else {
		c.ID, err = strconv.Atoi(params["id"])
		if err != nil {
			apierror.GenerateError("Trouble getting blog category ID", err, rw, r)
		}
	}

	err = c.Delete(dtx)
	if err != nil {
		apierror.GenerateError("Trouble deleting blog category", err, rw, r)
	}
	return encoding.Must(enc.Encode(c))
}

func UpdateBlog(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params, dtx *apicontext.DataContext) string {
	var b blog_model.Blog
	var err error

	idStr := r.FormValue("id")
	if idStr != "" {
		b.ID, err = strconv.Atoi(idStr)
		if err != nil {
			apierror.GenerateError("Trouble getting blog ID", err, rw, r)
		}
	} else {
		b.ID, err = strconv.Atoi(params["id"])
		if err != nil {
			apierror.GenerateError("Trouble getting blog ID", err, rw, r)
		}
	}
	b.Get(dtx)

	var tempBC []blog_model.BlogCategory

	title := r.FormValue("title")
	slug := r.FormValue("slug")
	text := r.FormValue("text")
	publishedDate := r.FormValue("publishedDate")
	userID := r.FormValue("userID")
	metaTitle := r.FormValue("metaTitle")
	metaDescription := r.FormValue("metaDescription")
	keywords := r.FormValue("keywords")
	active := r.FormValue("active")
	categoryIDs := r.Form["categoryID"]
	for _, v := range categoryIDs {
		var bc blog_model.BlogCategory
		bc.Category.ID, err = strconv.Atoi(v)
		tempBC = append(tempBC, bc)
	}

	if err != nil {
		apierror.GenerateError("Trouble getting blog", err, rw, r)
		return err.Error()
	}
	if title != "" {
		b.Title = title
	}
	if slug != "" {
		b.Slug = slug
	}
	if text != "" {
		b.Text = text
	}
	if publishedDate != "" {
		b.PublishedDate, err = time.Parse(timeFormat, publishedDate)
	}
	if userID != "" {
		b.UserID, err = strconv.Atoi(userID)
	}
	if metaTitle != "" {
		b.MetaTitle = metaTitle
	}
	if metaDescription != "" {
		b.MetaDescription = metaDescription
	}
	if keywords != "" {
		b.Keywords = keywords
	}
	if active != "" {
		b.Active, err = strconv.ParseBool(active)
	}
	if categoryIDs != nil {
		b.BlogCategories = tempBC
	}

	err = b.Update(dtx)
	if err != nil {
		apierror.GenerateError("Trouble updating blog", err, rw, r)
	}
	return encoding.Must(enc.Encode(b))
}

func DeleteBlog(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params, dtx *apicontext.DataContext) string {
	var b blog_model.Blog
	var err error
	idStr := r.FormValue("id")
	if idStr != "" {
		b.ID, err = strconv.Atoi(idStr)
		if err != nil {
			apierror.GenerateError("Trouble getting blog ID", err, rw, r)
		}
	} else {
		b.ID, err = strconv.Atoi(params["id"])
		if err != nil {
			apierror.GenerateError("Trouble getting blog ID", err, rw, r)
		}
	}
	err = b.Delete(dtx)
	if err != nil {
		apierror.GenerateError("Trouble deleting blog", err, rw, r)
	}
	return encoding.Must(enc.Encode(b))
}

func Search(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder, dtx *apicontext.DataContext) string {
	var err error

	title := r.FormValue("title")
	slug := r.FormValue("slug")
	text := r.FormValue("text")
	createdDate := r.FormValue("createdDate")
	publishedDate := r.FormValue("publishedDate")
	lastModified := r.FormValue("lastModified")
	userID := r.FormValue("userID")
	metaTitle := r.FormValue("metaTitle")
	metaDescription := r.FormValue("metaDescription")
	keywords := r.FormValue("keywords")
	active := r.FormValue("active")

	page := r.FormValue("page")
	results := r.FormValue("results")

	l, err := blog_model.Search(title, slug, text, publishedDate, createdDate, lastModified, userID, metaTitle, metaDescription, keywords, active, page, results, dtx)
	if err != nil {
		apierror.GenerateError("Trouble searching for blog", err, rw, r)
	}

	return encoding.Must(enc.Encode(l))
}
