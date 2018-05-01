package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
	"strconv"
	"time"
)

type (
	Word struct {
		ID          uint       `json:"id" gorm:"primary_key"`
		Word        string     `json:"word" gorm:"size:255; not null"`
		Translation string     `json:"translation" gorm:"size:255; not null"`
		DueDate     *time.Time `json:"due_date"`
		UserID      uint       `gorm:"index"`
		CategoryID  uint       `gorm:"index"`
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
	var words Words
	var all bool
	var categoryId uint
	var db_tmp *gorm.DB
	if _, ok := c.GetQuery("all"); ok {
		all = true
	}
	if category_id, ok := c.GetQuery("category_id"); ok {
		if id, err := strconv.Atoi(category_id); err != nil {
			categoryId = uint(id)
		}
	}
	if !all {
		//db_tmp = db.Where("DueDate < ", time.Now)
		db_tmp = db
	} else {
		db_tmp = db
	}
	db_tmp.Model(&User{ID: getUserId(c)}).
		Model(&Category{ID: categoryId}).
		Find(&words)

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
	if err := c.ShouldBindJSON(&word); err != nil {
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
	word.CategoryID = 0
	if value, ok := c.GetQuery("category_id"); ok {
		if categoryId, err := strconv.Atoi(value); err == nil {
			category_id := uint(categoryId)
			if !db.Model(&User{ID: getUserId(c)}).
				First(&Category{}, category_id).
				RecordNotFound() {
				word.CategoryID = category_id
			}
		}
	}
	if err := db.
		Model(&User{ID: getUserId(c)}).
		Association("Words").
		Append(&word).
		Error; err != nil {
		printLog(err)
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
	if err := c.ShouldBindJSON(&word); err != nil {
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
	if db.Model(&User{ID: getUserId(c)}).
		Model(&word).
		Updates(map[string]interface{}{"word": word.Word, "translation": word.Translation}).
		RecordNotFound() {
		c.JSON(
			http.StatusNotFound,
			gin.H{
				"status":  http.StatusNotFound,
				"message": "no word found!",
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
	if err := c.ShouldBindJSON(&word); err != nil {
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
	if db.
		Model(&User{ID: getUserId(c)}).
		Delete(&word).
		RecordNotFound() {
		c.JSON(
			http.StatusNotFound,
			gin.H{
				"status":  http.StatusNotFound,
				"message": "no word found!",
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
