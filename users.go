package main

import (
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
	"time"
)

const (
	secretKey        = "SOME_SECRET_KEY"
	token_time_delay = 60 * time.Minute
)

type (
	User struct {
		ID       uint   `json:"id" gorm:"primary_key"`
		Email    string `json:"email" gorm:"size:100; unique; not null" form:"email"`
		Password string `json:"password" gorm:"size:255; not null" form:"password"`
		Name     string `json:"name" gorm:"size:255"`
	}

	Users []User

	userClaims struct {
		UserId uint `json:"id"`
		jwt.StandardClaims
	}
)

func (User) TableName() string {
	return "user"
}

func generateFromPassword(password string) (string, error) {
	pw, err := bcrypt.GenerateFromPassword([]byte(password), 0)
	return string(pw), err
}

func (u *User) BeforeCreate(scope *gorm.Scope) error {
	u.Name = strings.ToLower(u.Name)
	scope.SetColumn("ID", 0)
	if pw, err := generateFromPassword(u.Password); err == nil {
		scope.SetColumn("Password", pw)
	}
	return nil
}

func fetchAllUser(c *gin.Context) {
	var users Users

	db.Find(&users)

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
			http.StatusBadRequest,
			gin.H{
				"status":  http.StatusBadRequest,
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
			http.StatusBadRequest,
			gin.H{
				"status":  http.StatusBadRequest,
				"message": "the password is incorrect!",
				"err":     err,
			},
		)
		return
	}

	token, err := createJwtToken(user.ID)
	PanicOnErr(err)
	c.Header("Authorization", token)

	c.JSON(
		http.StatusOK,
		gin.H{
			"status": http.StatusOK,
			"data":   user,
		},
	)

}

func registerUser(c *gin.Context) {

	var user User

	if err = c.ShouldBindJSON(&user); err != nil {
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

	token, err := createJwtToken(user.ID)
	PanicOnErr(err)
	c.Header("Authorization", token)

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

func loginOutUser(c *gin.Context) {

	c.Header("Authorization", "")

	user := User{}
	if userId, exists := c.Get("userId"); exists {
		user.ID = userId.(uint)
	}

	if db.First(&user, user.ID).RecordNotFound() {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"status":  http.StatusBadRequest,
				"message": "no user found!",
			},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"status": http.StatusOK,
			"data":   user,
		},
	)

}

func createJwtToken(id uint) (string, error) {

	claims := userClaims{
		id,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(token_time_delay).Unix(),
		},
	}

	rawToken := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return rawToken.SignedString([]byte(secretKey))

}

func authorized(c *gin.Context) {

	tokenString := c.GetHeader("Authorization")

	token, _ := jwt.ParseWithClaims(
		tokenString,
		&userClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		},
	)

	if !token.Valid {
		c.JSON(
			http.StatusUnauthorized,
			gin.H{
				"status":  http.StatusUnauthorized,
				"message": "token is invalid",
			},
		)
		c.Abort()
	} else {
		claims := token.Claims.(*userClaims)
		c.Set("userId", claims.UserId)
	}

}
