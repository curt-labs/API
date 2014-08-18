package database

import (
	"errors"
	"expvar"
	"github.com/ziutek/mymysql/autorc"
	_ "github.com/ziutek/mymysql/thrsafe"
	"log"
)

var (
	Statements = make(map[string]*autorc.Stmt, 0)
)

func PrepareAuthentication(authChan chan int) {
	UnPreparedStatements := make(map[string]string, 0)

	UnPreparedStatements["AuthStmt"] = `select id from ApiKey where api_key = ?`
	UnPreparedStatements["PrivateAuthStmt"] = `select ak.id from ApiKey as ak
				join ApiKeyType as akt on ak.type_id = akt.id
				where akt.type = 'PRIVATE'
				&& api_key = ?`

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

	authChan <- 1
}

func PrepareGeoLocation(prepChan chan int) {

	UnPreparedStatements := make(map[string]string, 0)

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
											c.postal_code, s.stateID, s.state, s.abbr as state_abbr, cty.countryID, cty.name as country_name, cty.abbr as country_abbr,
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
	UnPreparedStatements["PolygonStmt"] = `select s.stateID, s.state, s.abbr,
											(
												select COUNT(cl.locationID) from CustomerLocations as cl
												join Customer as c on cl.cust_id = c.cust_id
												join DealerTypes as dt on c.dealer_type = dt.dealer_type
												where dt.online = 0 && cl.stateID = s.stateID
											) as count
											from States as s
											where (
												select COUNT(cl.locationID) from CustomerLocations as cl
												join Customer as c on cl.cust_id = c.cust_id
												join DealerTypes as dt on c.dealer_type = dt.dealer_type
												where dt.online = 0 && cl.stateID = s.stateID
											) > 0
											order by s.state`
	UnPreparedStatements["MapPolygonCoordinatesForStateStmt"] = `select mp.ID,mpc.latitude, mpc.longitude
																	from MapPolygonCoordinates as mpc
																	join MapPolygon as mp on mpc.MapPolygonID = mp.ID
																	where mp.stateID = ?`
	UnPreparedStatements["WhereToBuyDealersStmt"] = `select c.customerID, c.name, c.email, c.address, c.address2, c.city, c.phone, c.fax, c.contact_person,
														c.latitude, c.longitude, c.searchURL, c.logo, c.website,
														c.postal_code, s.stateID, s.state, s.abbr as state_abbr, cty.countryID, cty.name as country_name, cty.abbr as country_abbr,
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

	UnPreparedStatements["LocalDealersStmt"] = `select cl.locationID, c.customerID, cl.name, c.email, cl.address, cl.city, cl.phone, cl.fax, cl.contact_person,
												cl.latitude, cl.longitude, c.searchURL, c.logo, c.website,
												cl.postalCode, s.stateID, s.state, s.abbr as state_abbr, cty.countryID, cty.name as country_name, cty.abbr as country_abbr,
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
												where dt.online = 0 && c.isDummy = 0 && dt.show = 1 && dtr.ID = mi.tier &&
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

	prepChan <- 1
}

func PrepareCustomer(prepChan chan int) {

	UnPreparedStatements := make(map[string]string, 0)

	UnPreparedStatements["CustomerStmt"] = `select c.customerID, c.name, c.email, c.address, c.address2, c.city, c.phone, c.fax, c.contact_person,
												c.latitude, c.longitude, c.searchURL, c.logo, c.website,
												c.postal_code, s.stateID, s.state, s.abbr as state_abbr, cty.countryID, cty.name as country_name, cty.abbr as country_abbr,
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
														s.stateID, s.state, s.abbr as state_abbr, cty.countryID, cty.name as cty_name, cty.abbr as cty_abbr
														from CustomerLocations as cl
														left join States as s on cl.stateID = s.stateID
														left join Country as cty on s.countryID = cty.countryID
														where cl.cust_id = ?`
	UnPreparedStatements["CustomerLocationByIdStmt"] = `select cls.*, dt.dealer_type as typeID, dt.type as dealerType, dt.online as typeOnline, dt.show as typeShow, dt.label as typeLabel,
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
													where cl.locationID = ? limit 1`

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

	UnPreparedStatements["CustomerIDFromKeyStmt"] = `select c.customerID from Customer as c
														join CustomerUser as cu on c.cust_id = cu.cust_ID
														join ApiKey as ak on cu.id = ak.user_id
														where ak.api_key = ?
														limit 1`

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

	prepChan <- 1
}

func PrepareCustomerUser(prepChan chan int) {

	UnPreparedStatements := make(map[string]string, 0)

	UnPreparedStatements["CustomerUserStmt"] = `select cu.* from CustomerUser as cu
													join Customer as c on cu.cust_ID = c.cust_id
													where c.customerID = '?'
													&& cu.active = 1`

	UnPreparedStatements["UserCustomerStmt"] = `select c.customerID, c.name, c.email, c.address, c.address2, c.city, c.phone, c.fax, c.contact_person,
												c.latitude, c.longitude, c.searchURL, c.logo, c.website,
												c.postal_code, s.stateID, s.state, s.abbr as state_abbr, cty.countryID, cty.name as country_name, cty.abbr as country_abbr,
												dt.dealer_type as typeID, dt.type as dealerType, dt.online as typeOnline, dt.show as typeShow, dt.label as typeLabel,
												dtr.ID as tierID, dtr.tier as tier, dtr.sort as tierSort,
												mi.ID as iconID, mi.mapicon, mi.mapiconshadow,
												mpx.code as mapix_code, mpx.description as mapic_desc,
												sr.name as rep_name, sr.code as rep_code, c.parentID
												from Customer as c
												join CustomerUser as cu on c.cust_id = cu.cust_ID
												left join States as s on c.stateID = s.stateID
												left join Country as cty on s.countryID = cty.countryID
												left join DealerTypes as dt on c.dealer_type = dt.dealer_type
												left join MapIcons as mi on dt.dealer_type = mi.dealer_type
												left join DealerTiers dtr on c.tier = dtr.ID
												left join MapixCode as mpx on c.mCodeID = mpx.mCodeID
												left join SalesRepresentative as sr on c.salesRepID = sr.salesRepID
												where cu.id = ?`
	UnPreparedStatements["CustomerUserAuthStmt"] = `select password, id, name, email, date_added, active, isSudo, passwordConverted from CustomerUser
													where email = ?
													&& active = 1
													limit 1`
	UnPreparedStatements["UpdateCustomerUserPassStmt"] = `update CustomerUser set password = ?, passwordConverted = 1
															where id = ? && active = 1`
	UnPreparedStatements["CustomerUserKeyAuthStmt"] = `select cu.* from CustomerUser as cu
														join ApiKey as ak on cu.id = ak.user_id
														join ApiKeyType as akt on ak.type_id = akt.id
														where UPPER(akt.type) = ?
														&& ak.api_key = ?
														&& cu.active = 1 && ak.date_added >= ?`
	UnPreparedStatements["CustomerUserKeysStmt"] = `select ak.api_key, akt.type, ak.date_added from ApiKey as ak
													join ApiKeyType as akt on ak.type_id = akt.id
													where user_id = ? && UPPER(akt.type) NOT IN (?)`
	UnPreparedStatements["UserAuthenticationKeyStmt"] = `select ak.api_key, ak.type_id, akt.type from ApiKey as ak
															join ApiKeyType as akt on ak.type_id = akt.id
															where UPPER(akt.type) = ?
															&& ak.user_id = ?`

	// This statement will run the trigger on the
	// ApiKey table to regenerate the api_key column
	// for the updated record
	UnPreparedStatements["ResetUserAuthenticationStmt"] = `update ApiKey as ak
															set ak.date_added = ?
															where ak.type_id = ?
															&& ak.user_id = ?`

	// This statement will renew the timer on the
	// authentication API key for the given user.
	// The disabling of the trigger is to turn off the
	// key regeneration trigger for this table
	UnPreparedStatements["EnableTriggerStmt"] = `SET @disable_trigger = 0`
	UnPreparedStatements["DisableTriggerStmt"] = `SET @disable_trigger = 1`
	UnPreparedStatements["RenewUserAuthenticationStmt"] = `update ApiKey as ak
															join ApiKeyType as akt on ak.type_id = akt.id
															set ak.date_added = ?
															where UPPER(akt.type) = ? && ak.user_id = ?`

	UnPreparedStatements["UserLocationStmt"] = `select cl.locationID, cl.name, cl.email, cl.address, cl.city,
												cl.postalCode, cl.phone, cl.fax, cl.latitude, cl.longitude,
												cl.cust_id, cl.contact_person, cl.isprimary, cl.ShippingDefault,
												s.stateID, s.state, s.abbr as state_abbr, cty.countryID, cty.name as cty_name, cty.abbr as cty_abbr
												from CustomerLocations as cl
												left join States as s on cl.stateID = s.stateID
												left join Country as cty on s.countryID = cty.countryID
												join CustomerUser as cu on cl.locationID = cu.locationID
												where cu.id = ?`

	UnPreparedStatements["CustomerUserFromKeyStmt"] = `select cu.* from CustomerUser as cu
														join ApiKey as ak on cu.id = ak.user_id
														join ApiKeyType as akt on ak.type_id = akt.id
														where akt.type != ? && ak.api_key = ?
														limit 1`

	UnPreparedStatements["CustomerUserFromId"] = `select cu.* from CustomerUser as cu
														join ApiKey as ak on cu.id = ak.user_id
														join ApiKeyType as akt on ak.type_id = akt.id
														where cu.id = ?
														limit 1`

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

	prepChan <- 1
}

func PrepareSearch(prepChan chan int) {

	UnPreparedStatements := make(map[string]string, 0)

	UnPreparedStatements["SearchPartAttributes"] = `select partID, value, field, sort,
														(
															match(field,value) against (? in natural language mode)
														) as score from PartAttribute
														where match(field,value) against (? in natural language mode)
														order by score desc
														limit ?,?`

	UnPreparedStatements["SearchPart"] = `select *, (
												match(shortDesc) against (? in natural language mode)
											) as score from Part
											where match(shortDesc)
											against (? in natural language mode)
											|| partID = ?
											order by score desc
											limit ?,?`

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

	prepChan <- 1
}

func PrepareCMS(prepChan chan int) {
	UnPreparedStatements := make(map[string]string, 0)

	UnPreparedStatements["AllCustomerContent"] = `select cc.id, cc.text,cc.added,cc.modified,cc.deleted,
													ct.type,ct.allowHTML,
													ccb.partID, ccb.catID
													from CustomerContent as cc
													left join CustomerContentBridge as ccb on cc.id = ccb.contentID
													join ContentType as ct on cc.typeID = ct.cTypeID
													join Customer as c on cc.custID = c.cust_id
													join CustomerUser as cu on c.cust_id = cu.cust_ID
													join ApiKey as ak on cu.id = ak.user_id
													where api_key = ?
													group by cc.id`

	UnPreparedStatements["CustomerContent"] = `select cc.id, cc.text,cc.added,cc.modified,cc.deleted,
													ct.type,ct.allowHTML,ccb.partID,ccb.catID
													from CustomerContent as cc
													join CustomerContentBridge as ccb on cc.id = ccb.contentID
													join ContentType as ct on cc.typeID = ct.cTypeID
													join Customer as c on cc.custID = c.cust_id
													join CustomerUser as cu on c.cust_id = cu.cust_ID
													join ApiKey as ak on cu.id = ak.user_id
													where api_key = ? and cc.id = ?
													limit 1`

	UnPreparedStatements["CustomerContentRevisions"] = `select ccr.old_text, ccr.new_text, ccr.date, ccr.changeType,
														ct1.type as newType, ct1.allowHTML as newAllowHtml,
														ct2.type as oldType, ct2.allowHTML as oldAllowHtml,
														ccr.userID as userId
														from CustomerContent_Revisions ccr
														left join ContentType ct1 on ccr.new_type = ct1.cTypeId
														left join ContentType ct2 on ccr.old_type = ct2.cTypeId
														join CustomerContent cc on ccr.contentID = cc.id
														join Customer as c on cc.custID = c.cust_id
														join CustomerUser as cu on c.cust_id = cu.cust_ID
														join ApiKey as ak on cu.id = ak.user_id
														where ak.api_key = ? and ccr.contentID = ?
														order by ccr.date`

	UnPreparedStatements["AllCustomerPartContent"] = `select cc.id, cc.text,cc.added,cc.modified,cc.deleted,
													ct.type,ct.allowHTML,ccb.partID
													from CustomerContent as cc
													join CustomerContentBridge as ccb on cc.id = ccb.contentID
													join ContentType as ct on cc.typeID = ct.cTypeID
													join Customer as c on cc.custID = c.cust_id
													join CustomerUser as cu on c.cust_id = cu.cust_ID
													join ApiKey as ak on cu.id = ak.user_id
													where api_key = ? and ccb.partID > 0
													group by ccb.partID, cc.id
													order by ccb.partID`

	UnPreparedStatements["CustomerPartContent"] = `select cc.id, cc.text,cc.added,cc.modified,cc.deleted,
													ct.type,ct.allowHTML,ccb.partID
													from CustomerContent as cc
													join CustomerContentBridge as ccb on cc.id = ccb.contentID
													join ContentType as ct on cc.typeID = ct.cTypeID
													join Customer as c on cc.custID = c.cust_id
													join CustomerUser as cu on c.cust_id = cu.cust_ID
													join ApiKey as ak on cu.id = ak.user_id
													where api_key = ? and ccb.partID = ?
													group by cc.id`

	UnPreparedStatements["AllCustomerCategoryContent"] = `select cc.id, cc.text,cc.added,cc.modified,cc.deleted,
													ct.type,ct.allowHTML,ccb.catID
													from CustomerContent as cc
													join CustomerContentBridge as ccb on cc.id = ccb.contentID
													join ContentType as ct on cc.typeID = ct.cTypeID
													join Customer as c on cc.custID = c.cust_id
													join CustomerUser as cu on c.cust_id = cu.cust_ID
													join ApiKey as ak on cu.id = ak.user_id
													where api_key = ? and ccb.catID > 0
													group by ccb.catID, cc.id
													order by ccb.catID`

	UnPreparedStatements["CustomerCategoryContent"] = `select cc.id, cc.text,cc.added,cc.modified,cc.deleted,
													ct.type,ct.allowHTML,ccb.catID
													from CustomerContent as cc
													join CustomerContentBridge as ccb on cc.id = ccb.contentID
													join ContentType as ct on cc.typeID = ct.cTypeID
													join Customer as c on cc.custID = c.cust_id
													join CustomerUser as cu on c.cust_id = cu.cust_ID
													join ApiKey as ak on cu.id = ak.user_id
													where api_key = ? and ccb.catID = ?
													group by cc.id`

	UnPreparedStatements["InsertCustomerContent"] = `insert into CustomerContent (
														text, custID, modified, userID, typeID, deleted
													)
													select ?, c.cust_id, now(), cu.id, ?, 0
													from Customer as c
													join CustomerUser as cu on c.cust_id = cu.cust_ID
													join ApiKey as ak on cu.id = ak.user_id
													where ak.api_key = ?`

	UnPreparedStatements["UpdateCustomerContent"] = `update CustomerContent as cc
														join Customer as c on cc.custID = c.cust_id
														join CustomerUser as cu on c.cust_id = cu.cust_ID
														join ApiKey as ak on cu.id = ak.user_id
														set cc.text = ?, cc.modified = now(),
														cc.userID = cu.id, cc.typeID = ?, cc.deleted = ?
														where ak.api_key = ? and cc.id = ?`

	UnPreparedStatements["CheckExistingCustomerContentBridge"] = `select count(id) from CustomerContentBridge
																	where partID = ? and catID = ? and contentID = ?`

	UnPreparedStatements["CreateCustomerContentBridge"] = `insert into CustomerContentBridge
															(partID, catID, contentID)
															values (?,?,?)`

	UnPreparedStatements["DeleteCustomerContentBridge"] = `delete from CustomerContentBridge
															where contentID in(
																select cc.id from CustomerContent as cc
																join Customer as c on cc.custID = c.cust_id
																join CustomerUser as cu on c.cust_id = cu.cust_ID
																join ApiKey ak on cu.id = ak.user_id
																where api_key = ? and contentID = ?
															) and partID = ? and catID = ?`

	UnPreparedStatements["MarkCustomerContentDeleted"] = `update CustomerContent as cc
															join Customer as c on cc.custID = c.cust_id
															join CustomerUser as cu on c.cust_id = cu.cust_ID
															join ApiKey as ak on cu.id = ak.user_id
															set cc.deleted = 1, cc.modified = now(),
															cc.userID = cu.id where ak.api_key = ?
															and cc.id = ?`

	UnPreparedStatements["GetContentTypeId"] = `select cTypeID, type, allowHTML from ContentType where type = ? limit 1`

	UnPreparedStatements["GetAllContentTypes"] = `select type, allowHTML from ContentType order by type`

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

	prepChan <- 1
}

func PrepareACES(acesChan chan int) {
	UnPreparedStatements := make(map[string]string, 0)

	UnPreparedStatements["GetAcesTypes"] = `select ID, name from AcesType order by ID`
	UnPreparedStatements["BigFatACESQuery"] = `select vp.ID, p.ACESPartTypeID, p.shortDesc, vp.PartNumber, bv.AAIABaseVehicleID, s.AAIASubmodelID, p.classID,
												(
												select group_concat(trim(ca1.vcdbID) order by cat1.sort) from ConfigAttribute as ca1
												join ConfigAttributeType as cat1 on ca1.ConfigAttributeTypeID = cat1.ID
												join VehicleConfigAttribute as vca1 on ca1.ID = vca1.AttributeID
												where vca1.VehicleConfigID = v.ConfigID && ca1.vcdbID > 0 && ca1.value != ""
												order by cat1.sort
												) as configIDs,
												(
												select group_concat(at.name order by cat2.sort) from AcesType as at
												join ConfigAttributeType as cat2 on at.ID = cat2.AcesTypeID
												join ConfigAttribute as ca2 on cat2.ID = ca2.ConfigAttributeTypeID
												join VehicleConfigAttribute as vca2 on ca2.ID = vca2.AttributeID
												join vcdb_Vehicle as v2 on vca2.VehicleConfigID = v2.ConfigID
												where v2.ID = v.ID && cat2.AcesTypeID > 0
												order by cat2.sort
												) as configNames,
												(
												select group_concat(ca3.value order by cat3.sort) from ConfigAttribute as ca3
												join ConfigAttributeType as cat3 on ca3.ConfigAttributeTypeID = cat3.ID
												join VehicleConfigAttribute as vca3 on ca3.ID = vca3.AttributeID
												join vcdb_Vehicle as v3 on vca3.VehicleConfigID = v3.ConfigID
												where v3.ID = v.ID && ca3.vcdbID = 0
												order by cat3.sort
												) as notes,
												(
												select group_concat(n.note) from Note as n
												where n.vehiclePartID = vp.ID
												) as part_notes
												from BaseVehicle as bv
												join vcdb_Make as ma on bv.MakeID = ma.ID
												join vcdb_Model as mo on bv.ModelID = mo.ID
												join vcdb_Vehicle as v on bv.ID = v.BaseVehicleID
												join vcdb_VehiclePart as vp on v.ID = vp.VehicleID
												join Part as p on vp.PartNumber = p.partID
												left join Submodel as s on v.SubModelID = s.ID
												order by vp.ID`

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

	acesChan <- 1
}

// Prepare all MySQL statements
func PrepareAll() error {

	// Bind the connection string from
	// environment variables
	BindDatabase()

	authChan := make(chan int)
	acesChan := make(chan int)
	geoChan := make(chan int)
	custChan := make(chan int)
	userChan := make(chan int)
	searchChan := make(chan int)
	cmsChan := make(chan int)

	go PrepareAuthentication(authChan)
	go PrepareACES(acesChan)
	go PrepareGeoLocation(geoChan)
	go PrepareCustomer(custChan)
	go PrepareCustomerUser(userChan)
	go PrepareSearch(searchChan)
	go PrepareCMS(cmsChan)

	log.Print("Executing Prepared Statements...")

	<-authChan
	log.Println("Authentication Statements Completed.......[OK]")
	<-acesChan
	log.Println("ACES Statements Completed.................[OK]")
	<-geoChan
	log.Println("GeoLocation Statements Completed..........[OK]")
	<-custChan
	log.Println("Customer Statements Completed.............[OK]")
	<-userChan
	log.Println("CustomerUser Statements Completed.........[OK]")
	<-searchChan
	log.Println("Search Statements Completed...............[OK]")
	<-cmsChan
	log.Println("CMS Statements Completed..................[OK]")

	log.Println("Finished.")

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
