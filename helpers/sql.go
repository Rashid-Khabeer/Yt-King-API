package helpers

import (
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func InitiateMySql() {
	instance, err := sql.Open("mysql", "root:@tcp(localhost:3306)/subsubbo_ytboost?parseTime=true")
	if err != nil {
		panic(err)
	}
	instance.SetConnMaxLifetime(time.Minute * 3)
	db = instance
	// db.SetMaxOpenConns(10)
	// db.SetMaxIdleConns(10)
}

func GetDB() *sql.DB {
	return db
}
