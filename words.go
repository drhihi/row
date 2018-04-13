package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
	"time"
	"strconv"
)

type (
	Word struct {
		ID          uint       `json:"id" gorm:"primary_key"`
		Word        string     `json:"word" gorm:"size:255; unique; not null" binding:"required"`
		Translation string     `json:"translation" gorm:"size:255; unique; not null" binding:"required"`
		CategoryID  uint       `json:"category_id" gorm:"not null" binding:"required"`
		UserID      uint       `json:"user_id" gorm:"not null" binding:"required"`
		DueDate     *time.Time `json:"due_date"`
	}

	Words []Word
)

func (Word) TableName() string {
	return "word"
}

func (w *Word) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("ID", 0)
	return nil
}

func fetchWords(c *gin.Context) {

	idCategory := c.Query("id_category")
	id, ok := strconv.Atoi(idCategory)
	if ok != nil || len(idCategory) == 0 {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"status":  http.StatusBadRequest,
				"message": "Incorrect URL - category id not specified.",
			},
		)
		return
	}

	var words Words
	userId := ParseUserIdFromToken(c)

	db.Where(
		&Word{
			CategoryID: uint(id),
			UserID:     userId,
		},
	).Find(&words)

	c.JSON(
		http.StatusOK,
		gin.H{
			"status": http.StatusOK,
			"data":   words,
		},
	)

}

func addWord(c *gin.Context) {

	var word Word

	if err = c.ShouldBindJSON(&word); err != nil {
		PanicOnErr(err)
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"status":  http.StatusBadRequest,
				"message": "the data is incorrect",
			},
		)
		return
	}

	if err = db.Create(&word).Error; err != nil {
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
			"data":   word,
		},
	)
}

func patchWord(c *gin.Context) {
	var word Word

	if err = c.BindQuery(&word); err != nil {
		PanicOnErr(err)
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"status":  http.StatusBadRequest,
				"message": "the data is incorrect #word",
			},
		)
		return
	}

	abortIfWordNotExists(c)

	if db.Model(&word).Updates(map[string]interface{}{"word": word.Word, "translation": word.Translation}).RecordNotFound() {
		c.JSON(
			http.StatusNotFound,
			gin.H{
				"status":  http.StatusNotFound,
				"message": "Update failed",
			},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"status": http.StatusOK,
			"data":   word,
		},
	)

}

func deleteWord(c *gin.Context) {
	var word Word

	if err = c.BindQuery(&word); err != nil {
		PanicOnErr(err)
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"status":  http.StatusBadRequest,
				"message": "the data is incorrect #word",
			},
		)
		return
	}

	abortIfWordNotExists(c)

	if db.Delete(&word).RecordNotFound() {
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
			"data":   word,
		},
	)

}

func abortIfWordNotExists(c *gin.Context) {
	var word Word
	userId := ParseUserIdFromToken(c)
	existingWord := db.Where(
		&Word{
			UserID: userId,
		},
	).First(&word)
	if existingWord == nil {
		c.JSON(
			http.StatusNotFound,
			gin.H{
				"status":  http.StatusNotFound,
				"message": "Word does not exist.",
			},
		)
		c.Abort()
	}
}
