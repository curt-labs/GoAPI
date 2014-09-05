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

func TestCustomerUser(t *testing.T) {
	Convey("Testing User Registration/ChangePass/Auth ", t, func() {
		Convey("Testing Register()", func() {
			var cu CustomerUser
			cu.Email = "bob@bob.com"
			pass := "test"
			customerID := 888
			isActive := true
			locationID := 1
			isSudo := true
			cust_ID := 1
			notCustomer := false
			custUser, err := cu.Register(pass, customerID, isActive, locationID, isSudo, cust_ID, notCustomer)
			So(custUser, ShouldNotBeNil)
			So(err, ShouldBeNil)
			Convey("BindAPIAccess", func() {
				err = cu.BindApiAccess()
				So(err, ShouldBeNil)
				So(len(cu.Keys), ShouldEqual, 3)
			})
			Convey("BindLocation", func() {
				err = cu.BindLocation()
				So(err, ShouldBeNil)
				So(cu.Location, ShouldNotBeNil)
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
				str, err := cu.ChangePass(oldPass, newPass, customerID)
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
				t.Log(customerID)
				err = DeleteCustomerUsersByCustomerID(customerID)
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
}