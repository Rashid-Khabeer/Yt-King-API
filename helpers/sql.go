package helpers

import (
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func InitiateMySql() {
	/// Local Mysql
	// instance, err := sql.Open("mysql", "admin:root@tcp(localhost:3307)/ytboost?parseTime=true")
	/// Local DB
	// instance, err := sql.Open("mysql", "rashid_yt5:Softstack1234$#@!@tcp(198.211.97.172:3306)/subsubbo_ytboost?parseTime=true")
	/// Live DB
	instance, err := sql.Open("mysql", "rashid_yt:Subsubpass1234$#@!@tcp(localhost:3306)/subsubbo_ytboost?parseTime=true")
	if err != nil {
		panic(err)
	}
	instance.SetConnMaxLifetime(time.Minute * 3)
	db = instance
	// db.SetMaxOpenConns(10)
}

func GetDB() *sql.DB {
	return db
}
