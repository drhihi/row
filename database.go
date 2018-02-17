package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var (
	db *gorm.DB
)

func init() {

	//connStr := "postgresql://postgres:PostGreSQL@localhost/vocaburise?sslmode=verify-full"
	connStr := "user=postgres dbname=vocaburise sslmode=disable"
	db, err = gorm.Open("postgres", connStr)

	PanicOnErr(err)

	db.AutoMigrate(&User{})

}
