package customer_new

import (
	"database/sql"
	"github.com/curt-labs/GoAPI/helpers/database"
	"github.com/curt-labs/GoAPI/helpers/redis"
	_ "github.com/go-sql-driver/mysql"
	// "log"
	"strconv"
)

var (
	//New
	getCustomer    = `select ` + customerFields + ` from Customer as c where c.cust_id = ? `
	getDealerTypes = `select dt.dealer_type, ` + dealerTypeFields + ` from DealerTypes as dt`
	getDealerTiers = `select dtr.ID, ` + dealerTierFields + ` from DealerTiers as dtr`
	getMapIcons    = `select mi.ID, mi.tier, mi.dealer_type, ` + mapIconFields + ` from MapIcons as mi`
	getMapixCodes  = ` select mpx.mCodeID, ` + mapixCodeFields + ` from MapixCode as mpx`
	getSalesReps   = ` select sr.salesRepID, ` + salesRepFields + ` from salesRepresentative as sr`
)

func DealerTypeMap() (map[int]DealerType, error) {
	typeMap := make(map[int]DealerType)
	var err error
	dTypes, err := GetDealerTypes()
	if err != nil {
		return typeMap, err
	}
	for _, dType := range dTypes {
		typeMap[dType.Id] = dType
		//set redis
		redis_key := "dealerType:" + strconv.Itoa(dType.Id)
		err = redis.Set(redis_key, dType)
	}
	return typeMap, err
}

func GetDealerTypes() ([]DealerType, error) {
	var dType DealerType
	var dTypes []DealerType
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return dTypes, err
	}
	defer db.Close()
	stmt, err := db.Prepare(getDealerTypes)
	if err != nil {
		return dTypes, err
	}
	defer stmt.Close()
	res, err := stmt.Query()
	if err != nil {
		return dTypes, err
	}
	for res.Next() {
		err = res.Scan(
			&dType.Id,
			&dType.Type,
			&dType.Online,
			&dType.Show,
			&dType.Label,
		)
		if err != nil {
			return dTypes, err
		}
		dTypes = append(dTypes, dType)
	}
	defer res.Close()
	return dTypes, err
}

func DealerTierMap() (map[int]DealerTier, error) {
	tierMap := make(map[int]DealerTier)
	var err error
	dTiers, err := GetDealerTiers()
	if err != nil {
		return tierMap, err
	}
	for _, dTier := range dTiers {
		tierMap[dTier.Id] = dTier
		//set redis
		redis_key := "dealerTier:" + strconv.Itoa(dTier.Id)
		err = redis.Set(redis_key, dTier)
	}
	return tierMap, err
}

func GetDealerTiers() ([]DealerTier, error) {
	var dTier DealerTier
	var dTiers []DealerTier
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return dTiers, err
	}
	defer db.Close()
	stmt, err := db.Prepare(getDealerTiers)
	if err != nil {
		return dTiers, err
	}
	defer stmt.Close()
	res, err := stmt.Query()
	if err != nil {
		return dTiers, err
	}
	for res.Next() {
		err = res.Scan(
			&dTier.Id,
			&dTier.Tier,
			&dTier.Sort,
		)
		if err != nil {
			return dTiers, err
		}

		dTiers = append(dTiers, dTier)
	}
	defer res.Close()
	return dTiers, err
}

func GetMapIcons() ([]MapIcon, error) {
	var mi MapIcon
	var mis []MapIcon
	var err error
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return mis, err
	}
	defer db.Close()
	stmt, err := db.Prepare(getMapIcons)
	if err != nil {
		return mis, err
	}
	defer stmt.Close()
	res, err := stmt.Query()
	for res.Next() {
		err = res.Scan(
			&mi.Id,
			&mi.TierId,
			&mi.DealerTypeId,
			&mi.MapIcon,
			&mi.MapIconShadow,
		)
		if err != nil {
			return mis, err
		}
		mis = append(mis, mi)
	}
	defer res.Close()
	return mis, err
}

func MapixMap() (map[int]MapixCode, error) {
	mapixMap := make(map[int]MapixCode)
	mcs, err := GetMapixCodes()
	if err != nil {
		return mapixMap, err
	}
	for _, mc := range mcs {
		mapixMap[mc.ID] = mc
		//set redis
		redis_key := "mapixCode:" + strconv.Itoa(mc.ID)
		err = redis.Set(redis_key, mc)
	}
	return mapixMap, err
}

func GetMapixCodes() ([]MapixCode, error) {
	var mc MapixCode
	var mcs []MapixCode
	var err error
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return mcs, err
	}
	defer db.Close()
	stmt, err := db.Prepare(getMapixCodes)
	if err != nil {
		return mcs, err
	}
	defer stmt.Close()
	res, err := stmt.Query()
	if err != nil {
		return mcs, err
	}
	for res.Next() {
		err = res.Scan(
			&mc.ID,
			&mc.Code,
			&mc.Description,
		)
		if err != nil {
			return mcs, err
		}
		mcs = append(mcs, mc)
	}
	defer res.Close()
	return mcs, err
}

func SalesRepMap() (map[int]SalesRepresentative, error) {
	repMap := make(map[int]SalesRepresentative)
	reps, err := GetSalesReps()
	if err != nil {
		return repMap, err
	}
	for _, rep := range reps {
		repMap[rep.ID] = rep
		//set redis
		redis_key := "salesRep:" + strconv.Itoa(rep.ID)
		err = redis.Set(redis_key, rep)
	}
	return repMap, err
}

func GetSalesReps() ([]SalesRepresentative, error) {
	var sr SalesRepresentative
	var srs []SalesRepresentative
	var err error
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return srs, err
	}
	defer db.Close()
	stmt, err := db.Prepare(getSalesReps)
	if err != nil {
		return srs, err
	}
	defer stmt.Close()
	res, err := stmt.Query()
	for res.Next() {
		err = res.Scan(
			&sr.ID,
			&sr.Name,
			&sr.Code,
		)
		if err != nil {
			return srs, err
		}
		srs = append(srs, sr)
	}
	defer res.Close()
	return srs, err
}
