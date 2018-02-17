package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"
)

type (
	User struct {
		ID       uint   `json:"id" gorm:"primary_key"`
		Email    string `json:"email" gorm:"size:100; unique; not null" form:"email"`
		Password string `json:"password" gorm:"size:255; not null" form:"password"`
		Name     string `json:"name" gorm:"size:255"`
	}

	Users []User
)

func (User) TableName() string {
	return "user"
}

func generateFromPassword(password string) (string, error) {
	pw, err := bcrypt.GenerateFromPassword([]byte(password), 0)
	return string(pw), err
}

func (u *User) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("ID", 0)
	if pw, err := generateFromPassword(u.Password); err == nil {
		scope.SetColumn("Password", pw)
	}
	return nil
}

func fetchAllUser(c *gin.Context) {
	var users Users

	db.Find(&users)

	if len(users) <= 0 {
		c.JSON(
			http.StatusNotFound,
			gin.H{
				"status":  http.StatusNotFound,
				"message": "no user found!",
			},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"status": http.StatusOK,
			"data":   users,
		},
	)

}

func loginUser(c *gin.Context) {
	var user User

	if err = c.BindQuery(&user); err != nil {
		PanicOnErr(err)
		c.JSON(
			http.StatusNotFound,
			gin.H{
				"status":  http.StatusNotFound,
				"message": "the data is incorrect #1",
			},
		)
		return
	}

	password := user.Password

	if db.Where("email = ?", user.Email).First(&user).RecordNotFound() {
		c.JSON(
			http.StatusNotFound,
			gin.H{
				"status":  http.StatusNotFound,
				"message": "no user found!",
			},
		)
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		c.JSON(
			http.StatusNonAuthoritativeInfo,
			gin.H{
				"status":  http.StatusNonAuthoritativeInfo,
				"message": "the password is incorrect!",
				"err":     err,
			},
		)
		return
	}

	pw, err := generateFromPassword(strconv.Itoa(int(user.ID)))
	PanicOnErr(err)
	c.SetCookie("auth", pw, 30, "", "", true, false)
	cr.Do("SET", pw, user.ID, "EX", 30)
	result, err := cr.Do("CLIENT LIST")

	c.JSON(
		http.StatusOK,
		gin.H{
			"status": http.StatusOK,
			"data":   user,
			"cr":     result,
		},
	)

}

func registerUser(c *gin.Context) {

	var user User

	if err = c.ShouldBindJSON(&user); err != nil {
		PanicOnErr(err)
		c.JSON(
			http.StatusNotFound,
			gin.H{
				"status":  http.StatusNotFound,
				"message": "the data is incorrect",
			},
		)
		return
	}

	if err = db.Create(&user).Error; err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"status":  http.StatusBadRequest,
				"message": err,
			},
		)
		return
	}

	pw, err := generateFromPassword(strconv.Itoa(int(user.ID)))
	PanicOnErr(err)
	c.SetCookie("auth", pw, 30, "", "", true, false)

	c.JSON(
		http.StatusCreated,
		gin.H{
			"status": http.StatusCreated,
			"data": struct {
				ID    uint   `json:"id"`
				Email string `json:"email"`
			}{
				user.ID,
				user.Email,
			},
		},
	)
}
