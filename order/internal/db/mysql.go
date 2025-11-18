package db

import (
	"database/sql"
	"fmt"

	"github.com/example/order-service/config"
	_ "github.com/go-sql-driver/mysql"
)

func NewMySQLConnection() *sql.DB {
	cfg := config.LoadConfig()
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		cfg.DBUser, cfg.DBPass, cfg.DBHost, cfg.DBPort, cfg.DBName,
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
