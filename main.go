package main

import (
	"github.com/gin-gonic/gin"
	"os"
)

func main() {
	router := getRouter()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	router.Run(":" + port)
}

func getRouter() *gin.Engine {
	router := gin.Default()
	v1 := router.Group("/api/v1")

	usersGroup := v1.Group("/users")
	{
		usersGroup.GET("/", authorized, checkAdmin, fetchAllUser)
		usersGroup.POST("/register", registerUser)
		usersGroup.GET("/login", logInUser)
		usersGroup.GET("/logout", authorized, logOutUser)
	}

	categoryGroup := v1.Group("/categories")
	{
		categoryGroup.GET("/", authorized, fetchAllCategories)
		categoryGroup.POST("/", authorized, addCategory)
		categoryGroup.PATCH("/", authorized, patchCategory)
		categoryGroup.DELETE("/", authorized, deleteCategory)
	}

	wordGroup := v1.Group("/words")
	{
		wordGroup.GET("/", authorized, fetchWords)
		wordGroup.POST("/", authorized, addWord)
		wordGroup.PATCH("/", authorized, patchWord)
		wordGroup.DELETE("/", authorized, deleteWord)
	}

	return router
}
