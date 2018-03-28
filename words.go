package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
)

type (
	Word struct {
		ID         uint     `json:"id" gorm:"primary_key"`
		Name       string   `json:"name" gorm:"size:255; unique; not null"`
		NameEng    string   `json:"name_eng" gorm:"size:255; unique; not null"`
		Category   Category `json:"category" gorm:"auto_preload"`
		CategoryID uint
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

	if id_category, ok := c.Get("id_category"); ok && id_category != nil {
		id, _ := id_category.(uint)
		db.Where(
			&Word{
				Category: Category{
					ID: id,
				},
			},
		).Find(&words)
	} else {
		db.Find(&words)
	}

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

	if db.Model(&word).Updates(map[string]interface{}{"name": word.Name, "name_eng": word.NameEng}).RecordNotFound() {
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
