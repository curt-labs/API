package customer_new

import (
	"database/sql"
	"github.com/curt-labs/GoAPI/helpers/database"
	// "github.com/curt-labs/GoAPI/models/customer"
	_ "github.com/go-sql-driver/mysql"
	. "github.com/smartystreets/goconvey/convey"
	// "math/rand"
	// "strings"
	"testing"

	// "time"
)

func getRandomKey() string {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return ""
	}
	defer db.Close()

	stmt, err := db.Prepare("SELECT api_key FROM ApiKey WHERE type_id = 'EA181F86-3F74-4AD6-8884-829B4558B99D' ORDER BY RAND() LIMIT 1")
	// stmt, err := db.Prepare("SELECT api_key FROM ApiKey WHERE type_id = (SELECT id FROM ApiKeyType WHERE Type = 'Authentication') ORDER BY RAND() LIMIT 1")
	if err != nil {
		return ""
	}
	defer stmt.Close()
	var key string
	err = stmt.QueryRow().Scan(&key)
	if err != nil {
		return ""
	}
	return key
}

func getRandomKeyNonAuthentication() string {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return ""
	}
	defer db.Close()

	stmt, err := db.Prepare("SELECT api_key FROM ApiKey WHERE type_id = '209A05AD-7D42-4C88-B5FA-FEEACDD19AC2' ORDER BY RAND() LIMIT 1")
	// stmt, err := db.Prepare("SELECT api_key FROM ApiKey WHERE type_id = (SELECT id FROM ApiKeyType WHERE Type = 'Authentication') ORDER BY RAND() LIMIT 1")
	if err != nil {
		return ""
	}
	defer stmt.Close()
	var key string
	err = stmt.QueryRow().Scan(&key)
	if err != nil {
		return ""
	}
	return key
}
func updateApiTime(apiKey string) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return
	}
	defer db.Close()

	stmt, err := db.Prepare("UPDATE ApiKey SET date_added = NOW() WHERE api_key = ?")
	if err != nil {
		return
	}
	_, _ = stmt.Exec(apiKey)
	return
}

func randomUserId() (user CustomerUser) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return user
	}
	defer db.Close()

	stmt, err := db.Prepare("SELECT id FROM CustomerUser WHERE active = 1 AND NotCustomer = 0 AND isSudo  = 0 ORDER BY RAND() LIMIT 1")
	if err != nil {
		return user
	}
	defer stmt.Close()

	err = stmt.QueryRow().Scan(&user.Id)
	if err != nil {
		return user
	}

	return user
}

func TestCustomerUser(t *testing.T) {
	Convey("Testing User Registration/ChangePass/Auth ", t, func() {
		Convey("Testing Register()", func() {
			var cu CustomerUser
			var err error
			cu.Email = "bob@bob.com"
			cu.Password = "test"
			cu.OldCustomerID = 888
			cu.Active = true
			cu.Location.Id = 1
			cu.Sudo = true
			cu.CustomerID = 1
			cu.Current = false
			err = cu.Create()
			So(cu, ShouldNotBeNil)
			So(err, ShouldBeNil)
			Convey("BindAPIAccess", func() {
				err = cu.BindApiAccess()
				So(err, ShouldBeNil)
				So(len(cu.Keys), ShouldEqual, 3)
			})
			// Convey("BindLocation", func() {
			// 	err = cu.BindLocation()
			// 	So(err, ShouldBeNil)
			// 	So(cu.Location, ShouldNotBeNil)
			// })
			Convey("Get Location", func() {
				err = cu.GetLocation()
				So(err, ShouldBeNil)
			})
			Convey("Update CustomerUser", func() {
				cu.Name = "Peanut"
				cu.Email = "tim@bob.com"
				cu.Active = false
				cu.Location.Id = 2
				cu.Sudo = false
				cu.Current = true
				err = cu.UpdateCustomerUser()
				So(err, ShouldBeNil)
			})
			Convey("Changing Password", func() {
				So(cu.Id, ShouldNotBeNil)
				oldPass := "test"
				newPass := "jerk"
				str, err := cu.ChangePass(oldPass, newPass)
				So(err, ShouldBeNil)
				So(str, ShouldEqual, "success")
				Convey("Now, Authenticate", func() {
					password := "jerk"
					cust, err := cu.UserAuthentication(password)
					So(err, ShouldBeNil)
					So(cust, ShouldNotBeNil)
					Convey("Reset Password", func() {
						newPass, err := cu.ResetPass(cu.Id)
						So(err, ShouldBeNil)
						So(newPass, ShouldNotEqual, password)

						Convey("Deleting CustomerUser", func() { //Watch - seems to delete; is it true?
							t.Log("cuid", cu.Id)
							err = cu.Delete()
							So(err, ShouldBeNil)
						})

					})

				})
			})
			Convey("Delete CustUsers by CustomerID", func() {
				t.Log(cu.OldCustomerID)
				err = DeleteCustomerUsersByCustomerID(cu.OldCustomerID)
				So(err, ShouldBeNil)
			})
		})
		key := getRandomKey()
		Convey("UserAutByKey", func() {
			t.Log(key)
			cust, err := UserAuthenticationByKey(key)
			So(err, ShouldNotBeNil)
			//update timestamp
			updateApiTime(key)
			cust, err = UserAuthenticationByKey(key)
			t.Log("Cust", cust)
			So(err, ShouldBeNil)
			So(cust, ShouldNotBeNil)

		})

	})

	//meddler calls this auth function
	Convey("Test GetCustomerUserFromKey", t, func() {
		key := getRandomKeyNonAuthentication()
		t.Log("KEY", key)
		u, err := GetCustomerUserFromKey(key)
		So(u, ShouldNotBeNil)
		t.Log(u)
		So(err, ShouldBeNil)
	})

}
