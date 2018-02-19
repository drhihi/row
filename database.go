package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"os"
)

var (
	db *gorm.DB
)

func init() {

	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = "user=postgres dbname=vocaburise sslmode=disable"
	}
	db, err = gorm.Open("postgres", connStr)

	PanicOnErr(err)

	db.AutoMigrate(&User{})

}
