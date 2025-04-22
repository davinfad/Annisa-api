package database

import (
	"database/sql"
	"log"
)

func InitDb() (*sql.DB, error) {
	dsn := "root:@tcp(127.0.0.1:3306)/annisa-api?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal("DB Ping Error:", err)
		return nil, err
	}
	return db, nil
}
