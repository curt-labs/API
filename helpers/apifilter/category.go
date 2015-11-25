package apifilter

import (
	"database/sql"
	"fmt"
	"github.com/curt-labs/API/helpers/database"
	"github.com/curt-labs/API/helpers/sortutil"
	"github.com/curt-labs/API/models/products"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

var (
	GetCategoryAttributes = `
		select pa.field, pa.value, group_concat(distinct pa.partID) as parts from PartAttribute as pa
		join Part as p on pa.partID = p.partID
		join CatPart as cp on p.partID = cp.partID
		where
		(p.status = 800 || p.status = 900) &&
		cp.catID = ? && !FIND_IN_SET(pa.field, ?) && pa.canFilter = 1
		group by pa.field, pa.value
		order by pa.field, pa.value`
	GetCategoryPrices = `
		select distinct pr.price, GROUP_CONCAT(pr.partID) as parts from Price as pr
		join Part as p on pr.partID = p.partID
		join CatPart as cp on p.partID = cp.partID
		where cp.catID = ? && lower(pr.priceType) = 'list' &&
		(p.status = 800 || p.status = 900)
		group by pr.price`
	GetCategoryGroup = `select distinct cp.catID as cats
											from CatPart as cp
											where FIND_IN_SET(cp.catID, bottom_category_ids(?))`
)

func RenderSelections(r *http.Request) map[string][]string {
	data := make(map[string][]string, 0)

	return data
}

func CategoryFilter(cat products.Category, specs *map[string][]string) ([]Options, error) {

	var filtered FilteredOptions

	attrChan := make(chan error)

	go func() {
		if results, err := filtered.categoryGroupAttributes(cat, specs); err == nil {
			filtered = append(filtered, results...)
		}
		attrChan <- nil
	}()

	select {
	case <-attrChan:

	case <-time.After(5 * time.Second):
		return FilteredOptions{}, nil
	}

	sortutil.AscByField(filtered, "Key")

	return filtered, nil
}

func (filtered FilteredOptions) categoryGroupAttributes(cat products.Category, specs *map[string][]string) (FilteredOptions, error) {

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

	idRows, err := idQry.Query(cat.CategoryID)
	if err != nil {
		return FilteredOptions{}, err
	}
	defer idRows.Close()

	ids := make([]int, 0)
	// ids = append(ids, cat.ID) // don't forget the given category
	for idRows.Next() {
		var i *int
		err := idRows.Scan(&i)
		if err == nil && i != nil {
			ids = append(ids, *i)
		}
	}

	excludedAttributeTypes := getExcludedAttributeTypes()
	filterResults := make(map[string]Options, 0)
	attrCh := make(chan error)
	for _, id := range ids {
		go func(catID int) {
			if results, err := categoryAttributes(catID, excludedAttributeTypes); err == nil {
				for key, result := range results {
					filterResults[key] = result
				}
			}
			attrCh <- nil
		}(id)
	}

	priceCh := make(chan error)
	filterResults["Price"] = Options{
		Key:     "Price",
		Options: make([]Option, 0),
	}

	for _, id := range ids {
		go func(catID int) {
			if results, err := categoryPrices(catID); err == nil {
				opts := make([]Option, 0)
				for _, res := range results {
					opts = append(opts, res.Options...)
				}

				fr := filterResults["Price"]
				fr.Options = append(fr.Options, opts...)
				filterResults["Price"] = fr
			}
			priceCh <- nil
		}(id)
	}

	for _, _ = range ids {
		<-attrCh
		<-priceCh
	}
	close(attrCh)
	close(priceCh)

	for key, res := range filterResults {
		indexed := make(map[string]int, 0)
		opts := make(map[string]Option, 0)
		for i, opt := range res.Options {
			if _, ok := indexed[opt.Value]; ok {
				idxOpt := opts[opt.Value]
				idxOpt.Products = append(idxOpt.Products, opt.Products...)
				curOpt := opts[opt.Value]
				curOpt.Products = removeDuplicates(idxOpt.Products)
				opts[opt.Value] = curOpt
			} else {
				res.Options = append(res.Options, opt)
				if specs != nil {
					for k, vals := range *specs {
						if strings.ToLower(key) == strings.ToLower(k) {
							for _, val := range vals {
								if strings.ToLower(opt.Value) == strings.ToLower(val) {
									opt.Selected = true
								}
							}
							break
						}
					}
				}
				opts[opt.Value] = opt
				indexed[opt.Value] = i
			}
		}

		res.Options = make([]Option, 0)
		mapped := make(map[string]string, 0)
		for _, opt := range opts {
			if _, ok := mapped[opt.Value]; !ok {
				sort.Ints(opt.Products)
				res.Options = append(res.Options, opt)
				mapped[opt.Value] = opt.Value
			}
		}
		if len(res.Options) > 1 {
			sortutil.AscByField(res.Options, "Value")
		}
		if len(res.Options) > 0 {
			filtered = append(filtered, res)
		}
	}

	sortutil.AscByField(filtered, "Key")
	return filtered, nil
}

func categoryAttributes(catID int, excludedAttributeTypes []string) (map[string]Options, error) {
	mapped := make(map[string]Options, 0)

	var exs string
	for i, ex := range excludedAttributeTypes {
		if i == 0 {
			exs = ex
		} else {
			exs = fmt.Sprintf("%s,%s", exs, ex)
		}
	}

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return map[string]Options{}, err
	}
	defer db.Close()

	stmt, err := db.Prepare(GetCategoryAttributes)
	if err != nil {
		return map[string]Options{}, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(catID, strings.Join(excludedAttributeTypes, ","))
	if err != nil {
		return map[string]Options{}, err
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

		var opts Options
		var ok bool
		if opts, ok = mapped[*key]; !ok {
			opts = Options{
				Key:     *key,
				Options: make([]Option, 0),
			}
			mapped[*key] = opts
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

		opts.Options = append(opts.Options, opt)
		mapped[opts.Key] = opts
	}

	return mapped, nil
}

func categoryPrices(catID int) (map[string]Options, error) {
	mapped := make(map[string]Options, 0)
	opt := Options{
		Key:     "Price",
		Options: make([]Option, 0),
	}

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return mapped, err
	}
	defer db.Close()

	stmt, err := db.Prepare(GetCategoryPrices)
	if err != nil {
		return mapped, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(catID)
	if err != nil {
		return mapped, err
	}
	defer rows.Close()

	priceMap := make(map[float64][]int, 0)

	lows := make([]int, 0)
	exists := false
	for rows.Next() {
		var price *float64
		var parts *string
		err = rows.Scan(&price, &parts)
		if err != nil || price == nil || parts == nil {
			// We're including the parts nil check here
			// so we don't display attributes when there
			// are no parts, although in theory, if that were
			// the case, the attributes wouldn't be here
			// in the first place.
			continue
		}

		partIDS := make([]int, 0)
		strParts := strings.Split(*parts, ",")
		for _, strPart := range strParts {
			if pID, err := strconv.Atoi(strPart); err == nil {
				partIDS = append(partIDS, pID)
			}
		}
		if _, ok := priceMap[*price]; !ok {
			priceMap[*price] = partIDS
		} else {
			ids := priceMap[*price]
			ids = append(ids, partIDS...)
			priceMap[*price] = ids
		}

		for _, def := range opt.Options {
			val := strings.Replace(def.Value, "$", "", -1)
			segs := strings.Split(val, " - ")
			if len(segs) < 2 {
				continue
			}

			low, err := strconv.ParseFloat(segs[0], 64)
			if err != nil {
				continue
			}
			high, err := strconv.ParseFloat(segs[1], 64)
			if err != nil {
				continue
			}

			if *price >= low && *price <= high {
				exists = true
			}
		}

		if !exists {
			lows = append(lows, (int(*price)/50)*50)
		}
	}

	sort.Ints(lows)
	existing := make(map[string]Option, 0)
	for _, low := range lows {
		val := fmt.Sprintf("$%d - $%d", low, low+50)
		o := Option{
			Value:    val,
			Selected: false,
		}
		for key, pm := range priceMap {
			if key >= float64(low) && key <= float64(low+50) {
				o.Products = append(o.Products, pm...)
			}
		}

		o.Products = removeDuplicates(o.Products)
		if ex, ok := existing[val]; !ok {
			existing[val] = o
		} else {
			ex.Products = append(ex.Products, o.Products...)
			ex.Products = removeDuplicates(ex.Products)
			existing[val] = ex
		}
	}

	for _, ex := range existing {
		opt.Options = append(opt.Options, ex)
	}

	mapped["Price"] = opt
	return mapped, nil
}
