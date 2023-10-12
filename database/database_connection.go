package database

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB() {
	var error error
	dbDriver := "mysql"
	dbUser := "username"
	dbPass := "password"
	dbName := "go_sql_user"
	DB, error = sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if error != nil {
		log.Fatal("Sql open connection failed", error)
	}
}
