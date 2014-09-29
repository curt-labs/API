package apifilter

import (
	"database/sql"
	"fmt"
	"github.com/curt-labs/GoAPI/helpers/database"
	"github.com/curt-labs/GoAPI/helpers/sortutil"
	"github.com/curt-labs/GoAPI/models/products"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"strconv"
	"strings"
	"time"
)

var (
	ExcludedCategoryAttributes = []string{"UPC", "Weight"}
	GetCategoryAttributes      = `
		select pa.field, pa.value, group_concat(distinct pa.partID) as parts from PartAttribute as pa
		join Part as p on pa.partID = p.partID
		join CatPart as cp on p.partID = cp.partID
		where
		(p.status = 800 || p.status = 900) &&
		cp.catID = ?
		group by pa.field, pa.value
		order by pa.field, pa.value`
	GetCategoryGroup = `select bottom_category_ids(?) as cats`
)

func CategoryFilter(cat products.ExtendedCategory, specs []interface{}) ([]Options, error) {

	var filtered FilteredOptions

	attrChan := make(chan error)

	go func() {
		if results, err := filtered.categoryGroupAttributes(cat); err == nil {
			filtered = append(filtered, results...)
		} else {
			log.Println(err)
		}

		attrChan <- nil
	}()

	select {
	case err := <-attrChan:
		if err != nil {
			log.Printf("filter error: %s\n", err.Error())
		}
	case <-time.After(5 * time.Second):
		log.Println("filter generation timed out")
	}

	sortutil.AscByField(filtered, "Key")

	return filtered, nil

}

func (filtered FilteredOptions) categoryGroupAttributes(cat products.ExtendedCategory) (FilteredOptions, error) {

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return FilteredOptions{}, err
	}
	defer db.Close()

	idQry, err := db.Prepare(GetCategoryGroup)
	if err != nil {
		return FilteredOptions{}, err
	}
	defer idQry.Close()

	idRows, err := idQry.Query(cat.CategoryId)
	if err != nil {
		return FilteredOptions{}, err
	}
	defer idRows.Close()

	ids := make([]int, 0)
	ids = append(ids, cat.CategoryId) // don' forget the given category
	for idRows.Next() {
		var strIds *string
		if err := idRows.Scan(&strIds); err == nil && strIds != nil {
			strCats := strings.Split(*strIds, ",")
			for _, strCat := range strCats {
				if cID, err := strconv.Atoi(strCat); err == nil {
					ids = append(ids, cID)
				}
			}
		}
	}

	ch := make(chan error)
	for _, id := range ids {
		go func(catID int) {
			if results, err := categoryAttributes(catID); err == nil {
				filtered = append(filtered, results...)
			}
			ch <- nil
		}(id)
	}

	for _, _ = range ids {
		<-ch
	}
	close(ch)

	return filtered, nil
}

func categoryAttributes(catID int) (FilteredOptions, error) {

	mapped := make(map[string]Options, 0)

	var exs string
	for i, ex := range ExcludedCategoryAttributes {
		if i == 0 {
			exs = ex
		} else {
			exs = fmt.Sprintf("%s,%s", exs, ex)
		}
	}

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return FilteredOptions{}, err
	}
	defer db.Close()

	stmt, err := db.Prepare(GetCategoryAttributes)
	if err != nil {
		return FilteredOptions{}, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(catID)
	if err != nil {
		return FilteredOptions{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var key *string
		var val *string
		var parts *string
		err = rows.Scan(&key, &val, &parts)
		if err != nil || key == nil || val == nil || parts == nil || strings.Contains(exs, *key) {
			// We're including the parts nil check here
			// so we don't display attributes when there
			// are no parts, although in theory, if that were
			// the case, the attributes wouldn't be here
			// in the first place.
			continue
		}

		var fOption Options
		var ok bool
		if fOption, ok = mapped[*key]; !ok {
			fOption = Options{
				Key:     *key,
				Options: make([]Option, 0),
			}
		}

		opt := Option{
			Value:    *val,
			Selected: false,
			Products: make([]int, 0),
		}
		strParts := strings.Split(*parts, ",")
		for _, strPart := range strParts {
			if p, err := strconv.Atoi(strPart); err == nil {
				opt.Products = append(opt.Products, p)
			}
		}

		fOption.Options = append(fOption.Options, opt)
		mapped[fOption.Key] = fOption
	}

	f := make(FilteredOptions, 0)
	for _, o := range mapped {
		f = append(f, o)
	}

	return f, nil
}
