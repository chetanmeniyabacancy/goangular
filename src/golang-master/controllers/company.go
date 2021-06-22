package controllers

import (
	"database/sql"
	"fmt"
	"net/http"
	"golang-master/models"
	"github.com/jmoiron/sqlx"
	"encoding/json"
	"github.com/gorilla/mux"
	"strconv"
	"gopkg.in/go-playground/validator.v9"
	"golang-master/validation"
	"golang-master/lang"
	"golang-master/generallib"
)

// BaseHandler will hold everything that controller needs
type BaseHandler struct {
	db *sql.DB
}

// BaseHandler will hold everything that controller needs
type BaseHandlerSqlx struct {
	db *sqlx.DB
}

type GetCompanies struct {
	Status int64 `json:"status"`
    Message string `json:"message"`
	Data *models.Companies `json:"data"`
}

type GetCompaniesDataTables struct {
	Status int64 `json:"status"`
    Message string `json:"message"`
	Data *models.DataTablesCompanies `json:"data"`
}

type GetCompany struct {
	Status int64 `json:"status"`
    Message string `json:"message"`
	Data *models.Company `json:"data"`
}

type CustomValidationMessages struct {
	messages map[string]string
}

// NewBaseHandler returns a new BaseHandler
func NewBaseHandler(db *sql.DB) *BaseHandler {
	return &BaseHandler{
		db: db,
	}
}

// NewBaseHandler returns a new BaseHandler
func NewBaseHandlerSqlx(db *sqlx.DB) *BaseHandlerSqlx {
	return &BaseHandlerSqlx{
		db: db,
	}
}
// HelloWorld returns Hello, World
func (h *BaseHandler) GetCompanies(w http.ResponseWriter, r *http.Request) {
	companies := models.GetCompanies(h.db);
	
	if err := h.db.Ping(); err != nil {
		fmt.Println("DB Error")
	}

	for _, elem := range *companies {
		w.Write([]byte(elem.Name))
    }

}

// GetCompanies returns companies list
func (h *BaseHandlerSqlx) GetCompaniesSqlxDataTables(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("content-type", "application/json")

	decoder := json.NewDecoder(r.Body)
	var reqcompany *models.DataTablesRequest
	err := decoder.Decode(&reqcompany)
	
	if err != nil {
		json.NewEncoder(w).Encode(ErrHandler(lang.Get("invalid_requuest")))
		return
	}
	
	response := GetCompaniesDataTables{}
	data := models.GetCompaniesSqlxDataTables(h.db,reqcompany);
	
	// for _, elem := range *companies {
	// 	w.Write([]byte(elem.Name))
    // }
	response.Status = 1;
	response.Message = lang.Get("success");
	response.Data = data;

	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(response)
}


// GetCompanies returns companies list
func (h *BaseHandlerSqlx) GetCompaniesSqlx(w http.ResponseWriter, r *http.Request) {
	response := GetCompanies{}
	companies := models.GetCompaniesSqlx(h.db);
	
	// for _, elem := range *companies {
	// 	w.Write([]byte(elem.Name))
    // }
	response.Status = 1;
	response.Message = lang.Get("success");
	response.Data = companies;

	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// PostCompaniesSqlx create companies
func (h *BaseHandlerSqlx) PostCompanySqlx(w http.ResponseWriter, r *http.Request) {
	fmt.Println("sadad")
	w.Header().Set("content-type", "application/json")
	response := GetCompany{}

	decoder := json.NewDecoder(r.Body)
	var reqcompany *models.ReqCompany
	err := decoder.Decode(&reqcompany)
	fmt.Println(err)

	if err != nil {
		json.NewEncoder(w).Encode(ErrHandler(lang.Get("invalid_requuest")))
		return
	}
	
	v := validator.New()
	v = validation.Custom(v)

	err = v.Struct(reqcompany)

	if err != nil {
		resp := validation.ToErrResponse(err)
		response := validation.FinalErrResponse{}
		response.Status = 0;
		response.Message = lang.Get("errors");
		response.Data = resp;
		json.NewEncoder(w).Encode(response)
		return
	}
	
	company,errmessage := models.PostCompanySqlx(h.db,reqcompany);
	if errmessage != "" {
		json.NewEncoder(w).Encode(ErrHandler(errmessage))
		return
	}
	
	generallib.Measure();		
	generallib.GoChannleExample();
	go generallib.SendMail();
	response.Status = 1;
	response.Message = lang.Get("insert_success");
	response.Data = company;
	json.NewEncoder(w).Encode(response)
}

// GetCompany returns company 
func (h *BaseHandlerSqlx) GetCompany(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	response := GetCompany{}
	
	company,errmessage := models.GetCompany(h.db,vars["id"]);
	
	if errmessage != "" {
		json.NewEncoder(w).Encode(ErrHandler(errmessage))
		return
	}

	response.Status = 1;
	response.Message = lang.Get("success");
	response.Data = company;

	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(response)
}

//EditCompanies edit companies
func (h *BaseHandlerSqlx) EditCompany(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()  
	
	w.Header().Set("content-type", "application/json")
	vars := mux.Vars(r)
	response := GetCompany{}
	id,err:=strconv.ParseInt(vars["id"], 10, 64);
	if err != nil {
		json.NewEncoder(w).Encode(ErrHandler(lang.Get("invalid_requuest")))
		return
	}

	var reqcompany models.ReqCompany
	reqcompany.Status, err = strconv.ParseInt(r.FormValue("status"), 10, 64);
	reqcompany.Name = r.FormValue("name");

	if err != nil {
		json.NewEncoder(w).Encode(ErrHandler(lang.Get("invalid_requuest")))
		return
	}

	company,errmessage := models.EditCompany(h.db,&reqcompany,id);
	if errmessage != "" {
		json.NewEncoder(w).Encode(ErrHandler(errmessage))
		return
	}

	response.Status = 1;
	response.Message = lang.Get("update_success");
	response.Data = company;
	json.NewEncoder(w).Encode(response)
}

// DeleteCompany delete company 
func (h *BaseHandlerSqlx) DeleteCompany(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	
	errmessage := models.DeleteCompany(h.db,vars["id"]);
	
	if errmessage != "" {
		json.NewEncoder(w).Encode(ErrHandler(errmessage))
		return
	}
	
	successresponse := CommonSuccess{}
	successresponse.Status = 1;
	successresponse.Message = lang.Get("delete_success");

	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(successresponse)
}