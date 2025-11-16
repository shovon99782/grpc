package db

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func NewMySQLConnection() *sql.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		"user", "12345", "127.0.0.1", "3306", "orders",
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}

	if err := db.Ping(); err != nil {
		fmt.Println(err)
		panic(err)
	}

	fmt.Println("✅ Connected to stock_service MySQL")
	return db
}
