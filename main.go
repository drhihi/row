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
		v1.GET("/login", logInUser)
		v1.POST("/register", registerUser)
		v1.GET("/logout", authorized, logOutUser)
	}

	categoryGroup := v1.Group("/category")
	{
		categoryGroup.GET("/", fetchAllCategory)
		categoryGroup.POST("/", checkAdmin, addCategory)
		categoryGroup.PATCH("/", checkAdmin, patchCategory)
		categoryGroup.DELETE("/", checkAdmin, deleteCategory)
	}

	wordGroup := v1.Group("/word")
	{
		wordGroup.GET("/", fetchWords)
		wordGroup.POST("/", checkAdmin, addWord)
		wordGroup.PATCH("/", checkAdmin, patchWord)
		wordGroup.DELETE("/", checkAdmin, deleteWord)
	}

	port := os.Getenv("PORT")

	if port == "" {
		port = "8000"
	}

	router.Run(":" + port)

}
