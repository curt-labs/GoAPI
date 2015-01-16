package customer_ctlr

import (
	"github.com/curt-labs/GoAPI/helpers/apicontext"
	"github.com/curt-labs/GoAPI/helpers/encoding"
	apierr "github.com/curt-labs/GoAPI/helpers/error"
	"github.com/curt-labs/GoAPI/models/customer"
	"github.com/go-martini/martini"

	"encoding/json"
	"fmt"
	"io/ioutil"
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

	//default brand for setting key
	defaultBrandArray := []int{1}

	err := user.AuthenticateUser(defaultBrandArray)
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
	var err error
	dtx := &apicontext.DataContext{APIKey: key}
	dtx.BrandArray, err = dtx.GetBrandsFromKey()

	cust, err := customer.AuthenticateAndGetCustomer(key, dtx)
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

func ChangePassword(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, dtx *apicontext.DataContext) string {
	email := r.FormValue("email")
	oldPass := r.FormValue("oldPass")
	newPass := r.FormValue("newPass")
	var user customer.CustomerUser
	user.Email = email

	err := user.ChangePass(oldPass, newPass, dtx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
	}
	return encoding.Must(enc.Encode("Success"))
}

func GenerateApiKey(w http.ResponseWriter, r *http.Request, params martini.Params, enc encoding.Encoder, dtx *apicontext.DataContext) string {
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

	generated, err := user.GenerateAPIKey(generateType, dtx.BrandArray)
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

	if email == "" || pass == "" {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return "Email and password are required."
	}
	brands := strings.TrimSpace(r.FormValue("brands"))
	brandStrArr := strings.Split(brands, ",")
	//default to CURT (1)
	if len(brandStrArr) <= 1 {
		brandStrArr = []string{"1"}
	}
	var brandArr []int
	for _, brand := range brandStrArr {
		brint, err := strconv.Atoi(brand)
		if err != nil {
			apierr.GenerateError("Error parsing brand IDs", err, w, r)
		}
		brandArr = append(brandArr, brint)
	}
	//default to CURT (1)
	if len(brandArr) == 0 {
		brandArr = append(brandArr, 1)
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
	err = user.Create(brandArr)
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
