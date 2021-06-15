package main

import (
	"net/http"
	"log"
	"golang-master/config"
	"golang-master/controllers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
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
	user.HandleFunc("/login", hsqlx.Login).Methods("POST")
	user.HandleFunc("/logout", hsqlx.Logout).Methods("GET")

	company := r.PathPrefix("/admin/company").Subrouter()
	company.HandleFunc("/listfordatatables", hsqlx.GetCompaniesSqlxDataTables).Methods("POST")
	company.HandleFunc("/list", hsqlx.GetCompaniesSqlx).Methods("GET")
	company.HandleFunc("/", hsqlx.PostCompanySqlx).Methods("POST")
	company.HandleFunc("/", hsqlx.GetCompany).Methods("GET")
	company.HandleFunc("/{id}", hsqlx.EditCompany).Methods("PUT")
	company.HandleFunc("/{id}", hsqlx.DeleteCompany).Methods("DELETE")
	company.Use(hsqlx.Secret)

	r.HandleFunc("/", h.GetCompanies)
	// r.HandleFunc("/sqlx", hsqlx.GetCompaniesSqlx)
	
	c := cors.New(cors.Options{
        AllowedOrigins: []string{"*"},
        AllowCredentials: true,
		AllowedHeaders: []string{"*"},
		AllowedMethods:   []string{"GET", "DELETE", "POST", "PUT"},
    })

    s := c.Handler(r)
    http.ListenAndServe(":5000", s)
}

// Middleware function, which will be called for each request
func AdminValidationMiddleware(next http.Handler) http.Handler {
	    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// token := r.Header.Get("X-Session-Token")

			// if user, found := amw.tokenUsers[token]; found {
			// 	// We found the token in our map
			// 	log.Printf("Authenticated user %s\n", user)
			// 	// Pass down the request to the next middleware (or final handler)
				next.ServeHTTP(w, r)
			// } else {
				// Write an error and stop the handler chain
				// http.Error(w, "Not authorized", 401)
				// return
			// }
		})
}