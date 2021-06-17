package main

import (
	"net/http"
	"log"
	"golang-master/config"
	"golang-master/controllers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"path"
	"os"
)

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
        log.Fatal(err)
    }
	r := mux.NewRouter()
	db := config.ConnectDB()
	dbsqlx := config.ConnectDBSqlx()
	h := controllers.NewBaseHandler(db)
	hsqlx := controllers.NewBaseHandlerSqlx(dbsqlx)
	
	user := r.PathPrefix("/admin").Subrouter()
	user.HandleFunc("/backend", admin).Methods("GET")
	user.HandleFunc("/login", hsqlx.Login).Methods("POST")
	user.HandleFunc("/logout", hsqlx.Logout).Methods("GET")

	company := r.PathPrefix("/admin/company").Subrouter()
	company.HandleFunc("/listfordatatables", hsqlx.GetCompaniesSqlxDataTables).Methods("POST")
	company.HandleFunc("/list", hsqlx.GetCompaniesSqlx).Methods("GET")
	company.HandleFunc("/", hsqlx.PostCompanySqlx).Methods("POST")
	company.HandleFunc("/", hsqlx.GetCompany).Methods("GET")
	company.HandleFunc("/{id}", hsqlx.EditCompany).Methods("PUT")
	company.HandleFunc("/{id}", hsqlx.DeleteCompany).Methods("DELETE")
	company.Use(hsqlx.IsAuthorized)

	r.HandleFunc("/", h.GetCompanies)
	// r.HandleFunc("/sqlx", hsqlx.GetCompaniesSqlx)
	
	c := cors.New(cors.Options{
        AllowedOrigins: []string{os.Getenv("ALLOWED_ORIGINS")},
        AllowCredentials: true,
		AllowedHeaders: []string{"*"},
		AllowedMethods:   []string{"GET", "DELETE", "POST", "PUT"},
    })

    s := c.Handler(r)
    http.ListenAndServe(":5000", s)
}

// serves index file
func admin(w http.ResponseWriter, r*http.Request) {
    p := path.Dir("backend/src/index.html")
    // set header
    w.Header().Set("Content-type", "text/html")
    http.ServeFile(w, r, p)
}