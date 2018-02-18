package main

import (
	"github.com/gin-gonic/gin"
	"os"
)

func main() {

	router := gin.Default()
	v1 := router.Group("/api/v1")
	{
		v1.GET("/", fetchAllUser)
		v1.GET("/login", loginUser)
		v1.POST("/register", registerUser)
	}

	port := os.Getenv("PORT")

	if port == "" {
		port = "8000"
	}

	router.Run(":" + port)

}
