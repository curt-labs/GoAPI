package customer_ctlr

import (
	. "../../models"
	"../../plate"
	"net/http"
)

func UserAuthentication(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	pass := r.FormValue("password")

	user := CustomerUser{
		Email: email,
	}
	cust, err := user.UserAuthentication(pass)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	plate.ServeFormatted(w, r, cust)
}
