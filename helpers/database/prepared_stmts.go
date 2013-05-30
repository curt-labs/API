package database

import (
	"../mymysql/autorc"
	_ "../mymysql/thrsafe"
	"errors"
	"expvar"
	"log"
)

var (
	Statements = make(map[string]*autorc.Stmt, 0)
)

// Prepare all MySQL statements
func PrepareAll() error {

	UnPreparedStatements := make(map[string]string, 0)

	// Get the category that a part is tied to, by PartId
	UnPreparedStatements["PartCategoryStmt"] = `select c.catID, c.parentID, c.sort, c.dateAdded,
				c.catTitle, c.shortDesc, c.longDesc,
				c.image, c.isLifestyle, c.vehicleSpecific,
				cc.code, cc.font from Categories as c
				join CatPart as cp on c.catID = cp.catID
				left join ColorCode as cc on c.codeID = cc.codeID
				where cp.partID = ?
				order by c.sort
				limit 1`

	UnPreparedStatements["PartAllCategoryStmt"] = `select c.catID, c.dateAdded, c.parentID, c.catTitle, c.shortDesc, 
					c.longDesc,c.sort, c.image, c.isLifestyle, c.vehicleSpecific,
					cc.font, cc.code
					from Categories as c
					join CatPart as cp on c.catID = cp.catID
					join ColorCode as cc on c.codeID = cc.codeID
					where cp.partID = ?
					order by c.catID`

	// Get a category by catID
	UnPreparedStatements["ParentCategoryStmt"] = `select c.catID, c.parentID, c.sort, c.dateAdded,
					c.catTitle, c.shortDesc, c.longDesc,
					c.image, c.isLifestyle, c.vehicleSpecific,
					cc.code, cc.font from Categories as c
					left join ColorCode as cc on c.codeID = cc.codeID
					where c.catID = ?
					order by c.sort
					limit 1`

	// Get the top-tier categories i.e Hitches, Electrical
	UnPreparedStatements["TopCategoriesStmt"] = `select c.catID, c.parentID, c.sort, c.dateAdded,
					c.catTitle, c.shortDesc, c.longDesc,
					c.image, c.isLifestyle, c.vehicleSpecific,
					cc.code, cc.font from Categories as c
					left join ColorCode as cc on c.codeID = cc.codeID
					where c.parentID IS NULL or c.parentID = 0
					and isLifestyle = 0
					order by c.sort`

	UnPreparedStatements["SubCategoriesStmt"] = `select c.catID, c.parentID, c.sort, c.dateAdded,
					c.catTitle, c.shortDesc, c.longDesc,
					c.image, c.isLifestyle, c.vehicleSpecific,
					cc.code, cc.font from Categories as c
					left join ColorCode as cc on c.codeID = cc.codeID
					where c.parentID = ?
					and isLifestyle = 0
					order by c.sort`

	UnPreparedStatements["CategoryByNameStmt"] = `select c.catID, c.parentID, c.sort, c.dateAdded,
					c.catTitle, c.shortDesc, c.longDesc,
					c.image, c.isLifestyle, c.vehicleSpecific,
					cc.code, cc.font from Categories as c
					left join ColorCode as cc on c.codeID = cc.codeID
					where c.catTitle = ?
					order by c.sort`

	UnPreparedStatements["CategoryByIdStmt"] = `select c.catID, c.parentID, c.sort, c.dateAdded,
					c.catTitle, c.shortDesc, c.longDesc,
					c.image, c.isLifestyle, c.vehicleSpecific,
					cc.code, cc.font from Categories as c
					left join ColorCode as cc on c.codeID = cc.codeID
					where c.catID = ?
					order by c.sort`

	UnPreparedStatements["CategoryPartBasicStmt"] = `select cp.partID
					from CatPart as cp
					where cp.catID = ?
					order by cp.partID
					limit ?,?`

	UnPreparedStatements["CategoryContentStmt"] = `select ct.type, c.text from ContentBridge cb
					join Content as c on cb.contentID = c.contentID
					left join ContentType as ct on c.cTypeID = ct.cTypeID
					where cb.catID = ?`

	UnPreparedStatements["SearchDealerLocations"] = `select cls.*, dt.dealer_type as typeID, dt.type as dealerType, dt.online as typeOnline, dt.show as typeShow, dt.label as typeLabel,
														dtr.ID as tierID, dtr.tier as tier, dtr.sort as tierSort,
														cl.locationID, cl.name, cl.address,cl.city,
														cl.postalCode, cl.email, cl.phone,cl.fax,
														cl.latitude, cl.longitude, cl.cust_id, cl.isPrimary, cl.ShippingDefault, cl.contact_person,
														c.showWebsite, c.website, c.eLocalURL
														from CustomerLocations as cl
														join States as cls on cl.stateID = cls.stateID
														join Customer as c on cl.cust_id = c.cust_id
														join DealerTypes as dt on c.dealer_type = dt.dealer_type
														join DealerTiers as dtr on c.tier = dtr.ID
														where (dt.dealer_type = 2 or dt.dealer_type = 3) and c.isDummy = false
														and dt.show = true and (lower(cl.name) like ? || lower(c.name) like ?)`
	UnPreparedStatements["SearchDealerLocationsByType"] = `select cls.*, dt.dealer_type as typeID, dt.type as dealerType, dt.online as typeOnline, dt.show as typeShow, dt.label as typeLabel,
															dtr.ID as tierID, dtr.tier as tier, dtr.sort as tierSort,
															cl.locationID, cl.name, cl.address,cl.city,
															cl.postalCode, cl.email, cl.phone,cl.fax,
															cl.latitude, cl.longitude, cl.cust_id, cl.isPrimary, cl.ShippingDefault, cl.contact_person,
															c.showWebsite, c.website, c.eLocalURL
															from CustomerLocations as cl
															join States as cls on cl.stateID = cls.stateID
															join Customer as c on cl.cust_id = c.cust_id
															join DealerTypes as dt on c.dealer_type = dt.dealer_type
															join DealerTiers as dtr on c.tier = dtr.ID
															where dt.online = false and c.isDummy = false
															and dt.show = true and (lower(cl.name) like ? || lower(c.name) like ?)`
	UnPreparedStatements["SearchDealerLocationsByLatLng"] = `select cls.*, dt.dealer_type as typeID, dt.type as dealerType, dt.online as typeOnline, dt.show as typeShow, dt.label as typeLabel,
															dtr.ID as tierID, dtr.tier as tier, dtr.sort as tierSort,
															cl.locationID, cl.name, cl.address,cl.city,
															cl.postalCode, cl.email, cl.phone,cl.fax,
															cl.latitude, cl.longitude, cl.cust_id, cl.isPrimary, cl.ShippingDefault, cl.contact_person,
															c.showWebsite, c.website, c.eLocalURL
															from CustomerLocations as cl
															join States as cls on cl.stateID = cls.stateID
															join Customer as c on cl.cust_id = c.cust_id
															join DealerTypes as dt on c.dealer_type = dt.dealer_type
															join DealerTiers as dtr on c.tier = dtr.ID
															where dt.online = false and c.isDummy = false
															and dt.show = true and
															( ? * (
																2 * ATAN2(
																	SQRT((SIN(((cl.latitude - ?) * (PI() / 180)) / 2) * SIN(((cl.latitude - ?) * (PI() / 180)) / 2)) + ((SIN(((cl.longitude - ?) * (PI() / 180)) / 2)) * (SIN(((cl.longitude - ?) * (PI() / 180)) / 2))) * COS(? * (PI() / 180)) * COS(cl.latitude * (PI() / 180))),
																	SQRT(1 - ((SIN(((cl.latitude - ?) * (PI() / 180)) / 2) * SIN(((cl.latitude - ?) * (PI() / 180)) / 2)) + ((SIN(((cl.longitude - ?) * (PI() / 180)) / 2)) * (SIN(((cl.longitude - ?) * (PI() / 180)) / 2))) * COS(? * (PI() / 180)) * COS(cl.latitude * (PI() / 180))))
																)
															) < 100.0)`

	UnPreparedStatements["GetLocalDealerTiers"] = `select distinct dtr.* from DealerTiers as dtr
													join Customer as c on dtr.ID = c.tier
													join DealerTypes as dt on c.dealer_type = dt.dealer_type
													where dt.online = false and dt.show = true
													order by dtr.sort`
	UnPreparedStatements["GetLocalDealerTypes"] = `select m.ID as iconId, m.mapicon, m.mapiconshadow,
													dtr.ID as tierID, dtr.tier as tier, dtr.sort as tierSort,
													dt.dealer_type as dealerTypeId, dt.type as dealerType, dt.online, dt.show, dt.label
													from MapIcons as m
													join DealerTypes as dt on m.dealer_type = dt.dealer_type
													join DealerTiers as dtr on m.tier = dtr.ID
													where dt.show = true
													order by dtr.sort desc`
	UnPreparedStatements["GetEtailers"] = `select c.customerID, c.name, c.email, c.address, c.address2, c.city, c.phone, c.fax, c.contact_person,
											c.latitude, c.longitude, c.searchURL, c.logo, c.website,
											c.postal_code, s.state, s.abbr as state_abbr, cty.name as country_name, cty.abbr as country_abbr,
											dt.dealer_type as typeID, dt.type as dealerType, dt.online as typeOnline, dt.show as typeShow, dt.label as typeLabel,
											dtr.ID as tierID, dtr.tier as tier, dtr.sort as tierSort,
											mpx.code as mapix_code, mpx.description as mapic_desc,
											sr.name as rep_name, sr.code as rep_code, c.parentID
											from Customer as c
											left join States as s on c.stateID = s.stateID
											left join Country as cty on s.countryID = cty.countryID
											left join DealerTypes as dt on c.dealer_type = dt.dealer_type
											left join DealerTiers dtr on c.tier = dtr.ID
											left join MapixCode as mpx on c.mCodeID = mpx.mCodeID
											left join SalesRepresentative as sr on c.salesRepID = sr.salesRepID
											where dt.online = true && c.isDummy = false
											order by dtr.sort`
	UnPreparedStatements["PolygonStmt"] = `select s.stateID, s.state, s.abbr,(
												select COUNT(cl.locationID) from CustomerLocations as cl
												join Customer as c on cl.cust_id = c.cust_id
												join DealerTypes as dt on c.dealer_type = dt.dealer_type
												where dt.online = 0 && cl.stateID = s.stateID
											) as count, 
											(select group_concat(mpc.latitude)
											from MapPolygonCoordinates as mpc
											join MapPolygon as mp on mpc.MapPolygonID = mp.ID
											where mp.stateID = s.stateID
											order by mpc.ID) as latitudes,
											(select group_concat(mpc.longitude)
											from MapPolygonCoordinates as mpc
											join MapPolygon as mp on mpc.MapPolygonID = mp.ID
											where mp.stateID = s.stateID
											order by mpc.ID) as longitudes
											from States as s
											where (
												select COUNT(cl.locationID) from CustomerLocations as cl
												join Customer as c on cl.cust_id = c.cust_id
												join DealerTypes as dt on c.dealer_type = dt.dealer_type
												where dt.online = 0 && cl.stateID = s.stateID
											) > 0
											order by s.state`
	UnPreparedStatements["MapPolygonCoordinatesForStateStmt"] = `select mpc.latitude, mpc.longitude
															from MapPolygonCoordinates as mpc
															join MapPolygon as mp on mpc.MapPolygonID = mp.ID
															where mp.stateID = ?`
	UnPreparedStatements["WhereToBuyDealersStmt"] = `select c.customerID, c.name, c.email, c.address, c.address2, c.city, c.phone, c.fax, c.contact_person,
														c.latitude, c.longitude, c.searchURL, c.logo, c.website,
														c.postal_code, s.state, s.abbr as state_abbr, cty.name as country_name, cty.abbr as country_abbr,
														dt.dealer_type as typeID, dt.type as dealerType, dt.online as typeOnline, dt.show as typeShow, dt.label as typeLabel,
														dtr.ID as tierID, dtr.tier as tier, dtr.sort as tierSort,
														mi.ID as iconID, mi.mapicon, mi.mapiconshadow,
														mpx.code as mapix_code, mpx.description as mapic_desc,
														sr.name as rep_name, sr.code as rep_code, c.parentID
														from Customer as c
														join DealerTypes as dt on c.dealer_type = dt.dealer_type
														join DealerTiers dtr on c.tier = dtr.ID
														left join MapIcons as mi on dt.dealer_type = mi.dealer_type
														left join States as s on c.stateID = s.stateID
														left join Country as cty on s.countryID = cty.countryID
														left join MapixCode as mpx on c.mCodeID = mpx.mCodeID
														left join SalesRepresentative as sr on c.salesRepID = sr.salesRepID
														where dt.dealer_type = 1 and dtr.ID = 4 and c.isDummy = false and length(c.searchURL) > 1`
	UnPreparedStatements["CustomerStmt"] = `select c.customerID, c.name, c.email, c.address, c.address2, c.city, c.phone, c.fax, c.contact_person,
												c.latitude, c.longitude, c.searchURL, c.logo, c.website,
												c.postal_code, s.state, s.abbr as state_abbr, cty.name as country_name, cty.abbr as country_abbr,
												dt.dealer_type as typeID, dt.type as dealerType, dt.online as typeOnline, dt.show as typeShow, dt.label as typeLabel,
												dtr.ID as tierID, dtr.tier as tier, dtr.sort as tierSort,
												mi.ID as iconID, mi.mapicon, mi.mapiconshadow,
												mpx.code as mapix_code, mpx.description as mapic_desc,
												sr.name as rep_name, sr.code as rep_code, c.parentID
												from Customer as c
												left join States as s on c.stateID = s.stateID
												left join Country as cty on s.countryID = cty.countryID
												left join DealerTypes as dt on c.dealer_type = dt.dealer_type
												left join MapIcons as mi on dt.dealer_type = mi.dealer_type
												left join DealerTiers as dtr on c.tier = dtr.ID
												left join MapixCode as mpx on c.mCodeID = mpx.mCodeID
												left join SalesRepresentative as sr on c.salesRepID = sr.salesRepID
												where c.customerID = ?`
	UnPreparedStatements["CustomerLocationStmt"] = `select cl.locationID, cl.name, cl.email, cl.address, cl.city,
														cl.postalCode, cl.phone, cl.fax, cl.latitude, cl.longitude,
														cl.cust_id, cl.contact_person, cl.isprimary, cl.ShippingDefault,
														s.state, s.abbr as state_abbr, cty.name as cty_name, cty.abbr as cty_abbr
														from CustomerLocations as cl
														left join States as s on cl.stateID = s.stateID
														left join Country as cty on s.countryID = cty.countryID
														where cl.cust_id = ?`
	UnPreparedStatements["CustomerUserStmt"] = `select cu.* from CustomerUser as cu
												join Customer as c on cu.cust_ID = c.cust_id
												where c.customerID = '?'
												&& cu.active = 1`
	UnPreparedStatements["CustomerPriceStmt"] = `select distinct cp.price from ApiKey as ak
													join CustomerUser cu on ak.user_id = cu.id
													join Customer c on cu.cust_ID = c.cust_id
													join CustomerPricing cp on c.customerID = cp.cust_id
													where api_key = ?
													and cp.partID = ?`
	UnPreparedStatements["CustomerPartStmt"] = `select distinct ci.custPartID from ApiKey as ak
												join CustomerUser cu on ak.user_id = cu.id
												join Customer c on cu.cust_ID = c.cust_id
												join CartIntegration ci on c.customerID = ci.custID
												where ak.api_key = ?
												and ci.partID = ?`
	UnPreparedStatements["LocalDealersStmt"] = `select cl.locationID, c.customerID, cl.name, c.email, cl.address, cl.city, cl.phone, cl.fax, cl.contact_person,
												cl.latitude, cl.longitude, c.searchURL, c.logo, c.website,
												cl.postalCode, s.state, s.abbr as state_abbr, cty.name as country_name, cty.abbr as country_abbr,
												dt.dealer_type as typeID, dt.type as dealerType, dt.online as typeOnline, dt.show as typeShow, dt.label as typeLabel,
												dtr.ID as tierID, dtr.tier as tier, dtr.sort as tierSort,
												mi.ID as iconID, mi.mapicon, mi.mapiconshadow,
												mpx.code as mapix_code, mpx.description as mapic_desc,
												sr.name as rep_name, sr.code as rep_code, c.parentID
												from CustomerLocations as cl
												join Customer as c on cl.cust_id = c.cust_id
												join DealerTypes as dt on c.dealer_type = dt.dealer_type
												left join MapIcons as mi on dt.dealer_type = mi.dealer_type
												join DealerTiers as dtr on c.tier = dtr.ID
												left join States as s on cl.stateID = s.stateID
												left join Country as cty on s.countryID = cty.countryID
												left join MapixCode as mpx on c.mCodeID = mpx.mCodeID
												left join SalesRepresentative as sr on c.salesRepID = sr.salesRepID
												where dt.online = 0 && c.isDummy = 0 && dt.show = 1 &&
												( ? * (
													2 * ATAN2(
														SQRT((SIN(((cl.latitude - ?) * (PI() / 180)) / 2) * SIN(((cl.latitude - ?) * (PI() / 180)) / 2)) + ((SIN(((cl.longitude - ?) * (PI() / 180)) / 2)) * (SIN(((cl.longitude - ?) * (PI() / 180)) / 2))) * COS(? * (PI() / 180)) * COS(cl.latitude * (PI() / 180))),
														SQRT(1 - ((SIN(((cl.latitude - ?) * (PI() / 180)) / 2) * SIN(((cl.latitude - ?) * (PI() / 180)) / 2)) + ((SIN(((cl.longitude - ?) * (PI() / 180)) / 2)) * (SIN(((cl.longitude - ?) * (PI() / 180)) / 2))) * COS(? * (PI() / 180)) * COS(cl.latitude * (PI() / 180))))
													)
												) < ?)
												&& (
													(cl.latitude >= ? && cl.latitude <= ?) 
													&&
													(cl.longitude >= ? && cl.longitude <= ?) 
													||
													(cl.longitude >= ? && cl.longitude <= ?)
												)
												group by cl.locationID
												order by dtr.sort desc`

	if !Db.Raw.IsConnected() {
		Db.Raw.Connect()
	}

	c := make(chan int)

	for stmtname, stmtsql := range UnPreparedStatements {
		go PrepareStatement(stmtname, stmtsql, c)
	}

	for _, _ = range UnPreparedStatements {
		<-c
	}

	return nil
}

func PrepareStatement(name string, sql string, ch chan int) {
	stmt, err := Db.Prepare(sql)
	if err == nil {
		Statements[name] = stmt
	} else {
		log.Fatal(err)
	}
	ch <- 1
}

func GetStatement(key string) (stmt *autorc.Stmt, err error) {
	stmt, ok := Statements[key]
	if !ok {
		qry := expvar.Get(key)
		if qry == nil {
			err = errors.New("Invalid query reference")
		}
	}
	return

}
