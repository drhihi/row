package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/garyburd/redigo/redis"
	"os"
)

var (
	db *gorm.DB
	cr redis.Conn
)

func init() {

	connStr := os.Getenv("DATABASE_URL")
	db, err = gorm.Open("postgres", connStr)

	PanicOnErr(err)

	db.AutoMigrate(&User{})

}

func joinCR(){
	cr, err = redis.DialURL(os.Getenv("REDIS_URL"))
	PanicOnErr(err)
}

func closeCR() {
	cr.Close()
}