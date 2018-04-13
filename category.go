package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
)

type (
	Category struct {
		ID     uint   `json:"id" gorm:"primary_key"`
		Name   string `json:"name" gorm:"size:255; unique; not null" binding:"required"`
		UserID uint   `json:"user_id" gorm:"not null" binding:"required"`
	}

	Categories []Category
)

func (Category) TableName() string {
	return "category"
}

func (c *Category) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("ID", 0)
	return nil
}

func fetchAllCategories(c *gin.Context) {
	var categories Categories

	userId := ParseUserIdFromToken(c)
	db.Where(
		&Category{
			UserID: userId,
		},
	).Find(&categories)

	c.JSON(
		http.StatusOK,
		gin.H{
			"status": http.StatusOK,
			"data":   categories,
		},
	)

}

func addCategory(c *gin.Context) {

	var category Category

	if err = c.ShouldBindJSON(&category); err != nil {
		PanicOnErr(err)
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"status":  http.StatusBadRequest,
				"message": "the data is incorrect: " + err.Error(),
			},
		)
		return
	}

	if err = db.Create(&category).Error; err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"status":  http.StatusBadRequest,
				"message": err,
			},
		)
		return
	}

	c.JSON(
		http.StatusCreated,
		gin.H{
			"status": http.StatusCreated,
			"data":   category,
		},
	)
}

func patchCategory(c *gin.Context) {
	var category Category

	if err = c.ShouldBindJSON(&category); err != nil {
		PanicOnErr(err)
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"status":  http.StatusBadRequest,
				"message": "the data is incorrect #categoty",
			},
		)
		return
	}

	abortIfCategoryNotExists(c)

	if db.Model(&category).Update("name", category.Name).RecordNotFound() {
		c.JSON(
			http.StatusNotFound,
			gin.H{
				"status":  http.StatusNotFound,
				"message": "no category found!",
			},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"status": http.StatusOK,
			"data":   category,
		},
	)

}

func deleteCategory(c *gin.Context) {
	var category Category

	if err = c.BindQuery(&category); err != nil {
		PanicOnErr(err)
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"status":  http.StatusBadRequest,
				"message": "the data is incorrect #categoty",
			},
		)
		return
	}

	abortIfCategoryNotExists(c)

	if db.Delete(&category).RecordNotFound() {
		c.JSON(
			http.StatusNotFound,
			gin.H{
				"status":  http.StatusNotFound,
				"message": "no category found!",
			},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"status": http.StatusOK,
			"data":   category,
		},
	)

}

func abortIfCategoryNotExists(c *gin.Context) {
	var category Category
	userId := ParseUserIdFromToken(c)
	existingCategory := db.Where(
		&Category{
			UserID: userId,
		},
	).First(&category)
	if existingCategory == nil {
		c.JSON(
			http.StatusNotFound,
			gin.H{
				"status":  http.StatusNotFound,
				"message": "Category does not exist.",
			},
		)
		c.Abort()
	}
}
