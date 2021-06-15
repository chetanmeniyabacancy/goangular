package controllers

import (
	"golang-master/models"
	"net/http"
	"encoding/json"
	"gopkg.in/go-playground/validator.v9"
	"golang-master/validation"
	"golang-master/lang"
	"github.com/gorilla/sessions"
	"os"
	"fmt"
)

type LoginSuccess struct {
	Status int64 `json:"status"`
    Message string `json:"message"`
	Data *models.User `json:"data"`
}

var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))

//Login
func (h *BaseHandlerSqlx) Login(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "login")
	r.ParseForm()  
	
	w.Header().Set("content-type", "application/json")
	response := LoginSuccess{}

	var reqlogin models.ReqLogin
	reqlogin.Email = r.FormValue("email");
	reqlogin.Password = r.FormValue("password");

	v := validator.New()
	v = validation.Custom(v)

	err := v.Struct(reqlogin)

	if err != nil {
		resp := validation.ToErrResponse(err)
		response := validation.FinalErrResponse{}
		response.Status = 0;
		response.Message = lang.Get("errors");
		response.Data = resp;
		json.NewEncoder(w).Encode(response)
		return
	}
	
	user,errmessage := models.Login(h.db,&reqlogin);
	if errmessage != "" {
		json.NewEncoder(w).Encode(ErrHandler(errmessage))
		return
	}

	session.Values["authenticated"] = "1"
    session.Save(r, w)
	session, _ = store.Get(r, "login")
	fmt.Println("authenticated")
	fmt.Println(session.Values["authenticated"])

	response.Status = 1;
	response.Message = lang.Get("update_success");
	response.Data = user;
	json.NewEncoder(w).Encode(response)
	
}

func (h *BaseHandlerSqlx) Secret(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "login")
		// Check if user is authenticated
		fmt.Println("authenticated2")
		fmt.Println(session.Values["authenticated"])
		if session.Values["authenticated"] == "1" {
			next.ServeHTTP(w, r)
		} else {
			next.ServeHTTP(w, r)
			// http.Error(w, "Not authorized", 401)
			// return
		}
	})
}

func (h *BaseHandlerSqlx) Logout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
    session, _ := store.Get(r, "login")
	response := CommonSuccess{}
    // Revoke users authentication
    session.Values["authenticated"] = "0"
    session.Save(r, w)
	response.Status = 1;
	response.Message = lang.Get("logout_success");
	json.NewEncoder(w).Encode(response)
}