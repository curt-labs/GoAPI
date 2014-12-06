package customer_ctlr

import (
	"encoding/json"
	"fmt"
	"github.com/curt-labs/GoAPI/helpers/encoding"
	"github.com/curt-labs/GoAPI/models/customer"
	"github.com/go-martini/martini"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

//Post - Form Authentication
func AuthenticateUser(w http.ResponseWriter, r *http.Request, enc encoding.Encoder) string {
	email := r.FormValue("email")
	pass := r.FormValue("password")

	var user customer.CustomerUser
	user.Email = email
	user.Password = pass

	err := user.AuthenticateUser()
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return ""
	}
	err = user.GetLocation()
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return ""
	}

	err = user.GetKeys()
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return ""
	}

	var key string
	if len(user.Keys) != 0 {
		key = user.Keys[0].Key
	}
	cust, err := user.GetCustomer(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return ""
	}

	return encoding.Must(enc.Encode(cust))
}

//Get - Key (in params) Authentication
func KeyedUserAuthentication(w http.ResponseWriter, r *http.Request, enc encoding.Encoder) string {
	qs := r.URL.Query()
	key := qs.Get("key")
	log.Print("key", key)
	cust, err := customer.AuthenticateAndGetCustomer(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return ""
	}

	return encoding.Must(enc.Encode(cust))
}

//Makes user current
// func ResetAuthentication(w http.ResponseWriter, r *http.Request, enc encoding.Encoder) string { //Testing only
// 	var err error
// 	qs := r.URL.Query()
// 	id := qs.Get("id")
// 	var u customer.CustomerUser
// 	u.Id = id
// 	err = u.ResetAuthentication()
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return ""
// 	}
// 	return "Success"
// }

func GetUserById(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	qs := r.URL.Query()
	key := qs.Get("key")

	var err error
	id := params["id"]
	if id == "" {
		id = r.FormValue("id")
		if id == "" {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return ""
		}
	}

	var user customer.CustomerUser
	user.Id = id

	err = user.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return ""
	}
	return encoding.Must(enc.Encode(user))
}

func ResetPassword(w http.ResponseWriter, r *http.Request, enc encoding.Encoder) string {
	email := r.FormValue("email")
	custID := r.FormValue("customerID")
	if email == "" {
		http.Error(w, "no email provided", http.StatusInternalServerError)
		return ""
	}
	if custID == "" {
		http.Error(w, "customerID cannot be blank", http.StatusInternalServerError)
		return ""
	}

	var user customer.CustomerUser
	user.Email = email

	resp, err := user.ResetPass()
	if err != nil || resp == "" {
		http.Error(w, err.Error(), http.StatusUnauthorized)
	}
	return encoding.Must(enc.Encode(resp))
}

func ChangePassword(w http.ResponseWriter, r *http.Request, enc encoding.Encoder) string {
	email := r.FormValue("email")
	oldPass := r.FormValue("oldPass")
	newPass := r.FormValue("newPass")
	log.Print("old pass ", oldPass)
	var user customer.CustomerUser
	user.Email = email

	err := user.ChangePass(oldPass, newPass)
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
	}
	return encoding.Must(enc.Encode("Success"))
}

func GenerateApiKey(w http.ResponseWriter, r *http.Request, params martini.Params, enc encoding.Encoder) string {
	qs := r.URL.Query()
	key := qs.Get("key")
	if key == "" {
		key = r.FormValue("key")
	}

	user, err := customer.GetCustomerUserFromKey(key)
	if err != nil || user.Id == "" {
		http.Error(w, "failed to authenticate API key; you must provide a private key.", http.StatusInternalServerError)
		return ""
	}

	authed := false
	if user.Sudo == false {
		for _, k := range user.Keys {
			if k.Type == customer.PRIVATE_KEY_TYPE && k.Key == key {
				authed = true
				break
			}
		}
	} else {
		authed = true
	}

	if !authed {
		http.Error(w, "you do not have sufficient permissions to perform this operation.", http.StatusInternalServerError)
		return ""
	}

	generateType := params["type"]
	id := params["id"]
	if id == "" {
		http.Error(w, "you must provide a reference to the user whose key should be generated", http.StatusInternalServerError)
		return ""
	}
	if generateType == "" {
		http.Error(w, "you must provide the type of key to be generated", http.StatusInternalServerError)
		return ""
	}

	user.Id = id
	if err := user.Get(key); err != nil {
		http.Error(w, "failed to retrieve the reference user account", http.StatusInternalServerError)
		return ""
	}

	generated, err := user.GenerateAPIKey(generateType)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to generate an API key: %s", err.Error()), http.StatusInternalServerError)
		return ""
	}

	return encoding.Must(enc.Encode(generated))
}

//a/k/a CreateUser
func RegisterUser(w http.ResponseWriter, r *http.Request, enc encoding.Encoder) string {
	name := r.FormValue("name")
	email := r.FormValue("email")
	pass := r.FormValue("pass")
	customerID, err := strconv.Atoi(r.FormValue("customerID"))
	isActive, err := strconv.ParseBool(r.FormValue("isActive"))
	locationID, err := strconv.Atoi(r.FormValue("locationID"))
	isSudo, err := strconv.ParseBool(r.FormValue("isSudo"))
	cust_ID, err := strconv.Atoi(r.FormValue("cust_ID"))
	notCustomer, err := strconv.ParseBool(r.FormValue("notCustomer"))
	log.Print("pass creation ", pass)
	if email == "" || pass == "" {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return "Email and password are required."
	}

	var user customer.CustomerUser
	user.Email = email
	user.Password = pass
	if name != "" {
		user.Name = name
	}
	if customerID != 0 {
		user.OldCustomerID = customerID
	}
	if locationID != 0 {
		user.Location.Id = locationID
	}
	if cust_ID != 0 {
		user.CustomerID = cust_ID
	}
	user.Active = isActive
	user.Sudo = isSudo
	user.Current = notCustomer
	err = user.Create()
	// cu, err := user.Register(pass, customerID, isActive, locationID, isSudo, cust_ID, notCustomer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return ""
	}

	return encoding.Must(enc.Encode(user))
}
func DeleteCustomerUser(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	id := params["id"]
	var err error

	var cu customer.CustomerUser
	cu.Id = id
	err = cu.Delete()
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return ""
	}

	return encoding.Must(enc.Encode(cu))
}
func DeleteCustomerUsersByCustomerID(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	customerID, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return ""
	}

	err = customer.DeleteCustomerUsersByCustomerID(customerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return ""
	}

	return encoding.Must(enc.Encode("Success."))
}

func UpdateCustomerUser(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params) string {
	qs := r.URL.Query()
	key := qs.Get("key")

	var err error
	id := params["id"]
	if id == "" {
		id = r.FormValue("id")
		if id == "" {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return ""
		}
	}

	var cu customer.CustomerUser
	cu.Id = id
	err = cu.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}

	if strings.ToLower(r.Header.Get("Content-Type")) == "application/json" {
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return ""
		}

		if err := json.Unmarshal(data, &cu); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return ""
		}
	} else {
		name := r.FormValue("name")
		email := r.FormValue("email")
		isActive := r.FormValue("isActive")
		locationID := r.FormValue("locationID")
		isSudo := r.FormValue("isSudo")
		notCustomer := r.FormValue("notCustomer")
		if name != "" {
			cu.Name = name
		}
		if email != "" {
			cu.Email = email
		}
		if isActive != "" {
			cu.Active, err = strconv.ParseBool(isActive)
		}
		if locationID != "" {
			cu.Location.Id, err = strconv.Atoi(locationID)
		}
		if isSudo != "" {
			cu.Sudo, err = strconv.ParseBool(isSudo)
		}
		if notCustomer != "" {
			cu.Current, err = strconv.ParseBool(notCustomer)
		}
	}

	err = cu.UpdateCustomerUser()
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return ""
	}

	return encoding.Must(enc.Encode(cu))
}