package main

import (
	"github.com/gin-gonic/gin"
)

func main() {

	//joinCR()
	//defer closeCR()

	router := gin.Default()
	v1 := router.Group("/api/v1")
	{
		v1.GET("/", fetchAllUser)
		v1.GET("/login", loginUser)
		v1.POST("/register", registerUser)
	}

	router.Run()

}
