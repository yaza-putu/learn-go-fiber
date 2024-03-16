package database

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

func Conn() *sql.DB {
	dsn := "root:Temp123!@tcp(localhost:3306)/crud_fiber"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	return db
}
