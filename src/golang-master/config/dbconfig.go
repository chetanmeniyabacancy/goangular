package config

import ( 
	"database/sql"
	"os"
	"golang-master/constants"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

)

// ConnectDB opens a connection to the database
func ConnectDB() *sql.DB {
	db, err := sql.Open(os.Getenv("DBTYPE"), os.Getenv("DBUSERNAME")+":"+os.Getenv("DBPASSWORD")+"@/"+os.Getenv("DBNAME"))
	if err != nil {
		panic(err.Error())
	}
	return db
}

// ConnectDB opens a connection to the database
func ConnectDBSqlx() *sqlx.DB {
	db, err := sqlx.Open(constants.DBTYPE, constants.DBUSERNAME+":"+constants.DBPASSWORD+"@/"+constants.DBNAME)
	if err != nil {
		panic(err.Error())
	}
	
	return db
}