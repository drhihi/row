package main

import (
	"bytes"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

const (
	urlVersion      = "/api/v1"
	countUsers      = 2
	countUsersFail  = 2
	countCategories = 2
	countWords      = 3
	withFalse       = false
)

func getTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := getRouter()
	r.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "Not Found")
	})
	return r
}

func testHTTPResponse(t *testing.T, r *gin.Engine, req *http.Request, f func(w *httptest.ResponseRecorder) bool) {
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if !f(w) {
		t.Fail()
	}
}

func getUserIdFromToken(tokenString string) (userId uint) {
	token, _ := jwt.ParseWithClaims(
		tokenString,
		&userClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		},
	)
	if token != nil && token.Valid {
		claims := token.Claims.(*userClaims)
		userId = claims.UserId
	}
	return
}

func getTestUsers(fail bool) Users {

	var users Users
	for i := 1; i < countUsers+1; i++ {
		s := strconv.Itoa(i)
		name := fmt.Sprint(s, "@", s)
		user := User{
			Email:     fmt.Sprint(name, ".com"),
			Name:      name,
			Password:  fmt.Sprint(s, s, s),
			RoleAdmin: i%2 == 0,
		}
		users = append(users, user)
	}

	if withFalse && fail {
		for i := 1; i < countUsersFail+1; i++ {
			s := strconv.Itoa(i)
			name := fmt.Sprint(s, "@", s)
			user := User{
				Email:     fmt.Sprint(name, ".com"),
				Name:      name,
				Password:  "",
				RoleAdmin: i%2 == 0,
			}
			users = append(users, user)
		}
	}

	return users
}

func deleteUsers() {
	users := getTestUsers(false)
	for _, user := range users {
		db.Delete(user)
	}
}

func logInUserTest(t *testing.T, r *gin.Engine, user *User) string {
	t.Log("\n", user)
	urlPath := urlVersion + "/users/login"
	req, err := http.NewRequest("GET", urlPath, nil)
	req.Header.Set("Content-Type", "application/json")
	q := req.URL.Query()
	q.Set("email", user.Email)
	q.Set("password", user.Password)
	req.URL.RawQuery = q.Encode()
	if err != nil {
		t.Fatal(err)
	}
	testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		valid := w.Code == http.StatusOK
		if !valid {
			t.Log(w.Body)
		}
		token := w.Header().Get("Authorization")
		t.Log(token)
		req.Header.Set("Authorization", token)
		return withFalse || valid

	},
	)
	token := req.Header.Get("Authorization")
	return token
}

func TestUsersRegister(t *testing.T) {
	//deleteUsers()
	urlPath := urlVersion + "/users/register"
	r := getTestRouter()
	users := getTestUsers(false)
	for _, user := range users {
		b, err := json.Marshal(&user)
		if err != nil {
			t.Fatal(err)
		}
		req, err := http.NewRequest("POST", urlPath, bytes.NewBuffer(b))
		req.Header.Set("Content-Type", "application/json")
		if err != nil {
			t.Fatal(err)
		}
		testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
			valid := w.Code == http.StatusCreated
			if !valid {
				t.Log(w.Body)
			}
			return withFalse || valid
		},
		)
	}
}

func TestUsersLogin(t *testing.T) {
	r := getTestRouter()
	users := getTestUsers(true)
	for _, user := range users {
		logInUserTest(t, r, &user)
	}
}

func getTestCategories(email string) Categories {
	var categories Categories
	for i := 1; i < countCategories+1; i++ {
		s := strconv.Itoa(i)
		name := fmt.Sprintf("cat #%s (user: %s)", s, email)
		category := Category{
			Name: name,
		}
		categories = append(categories, category)
	}
	return categories
}

func TestCategoriesAdd(t *testing.T) {
	urlPath := urlVersion + "/categories/"
	r := getTestRouter()
	users := getTestUsers(true)
	for _, user := range users {
		token := logInUserTest(t, r, &user)
		categories := getTestCategories(user.Email)
		for _, category := range categories {
			b, err := json.Marshal(&category)
			if err != nil {
				t.Fatal(err)
			}
			req, err := http.NewRequest("POST", urlPath, bytes.NewBuffer(b))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", token)
			if err != nil {
				t.Fatal(err)
			}
			testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
				valid := w.Code == http.StatusCreated
				t.Log(w.Body)
				return withFalse || valid
			},
			)
		}
	}
}

func TestCategoriesPatch(t *testing.T) {
	urlPath := urlVersion + "/categories/"
	r := getTestRouter()
	users := getTestUsers(true)
	for _, user := range users {
		token := logInUserTest(t, r, &user)
		var categories Categories
		db.Model(&User{ID: getUserIdFromToken(token)}).Related(&categories)
		for _, category := range categories {
			if category.ID%2 == 0 {
				continue
			}
			category.Name = fmt.Sprint(category.Name, " *patch")
			category.UserID = 0
			b, err := json.Marshal(&category)
			if err != nil {
				t.Fatal(err)
			}
			req, err := http.NewRequest("PATCH", urlPath, bytes.NewBuffer(b))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", token)
			if err != nil {
				t.Fatal(err)
			}
			testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
				valid := w.Code == http.StatusOK
				t.Log(w.Body)
				return withFalse || valid
			},
			)
		}
	}
}

func TestCategoriesGet(t *testing.T) {
	urlPath := urlVersion + "/categories/"
	r := getTestRouter()
	users := getTestUsers(true)
	for _, user := range users {
		token := logInUserTest(t, r, &user)
		req, err := http.NewRequest("GET", urlPath, nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", token)
		if err != nil {
			t.Fatal(err)
		}
		testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
			valid := w.Code == http.StatusOK
			t.Log(w.Body)
			return withFalse || valid
		},
		)
	}
}

func TestCategoriesDelete(t *testing.T) {
	urlPath := urlVersion + "/categories/"
	r := getTestRouter()
	users := getTestUsers(true)
	for _, user := range users {
		token := logInUserTest(t, r, &user)
		var categories Categories
		db.Model(&User{ID: getUserIdFromToken(token)}).Related(&categories)
		for _, category := range categories {
			category.Name = ""
			category.UserID = 0
			b, err := json.Marshal(&category)
			if err != nil {
				t.Fatal(err)
			}
			req, err := http.NewRequest("DELETE", urlPath, bytes.NewBuffer(b))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", token)
			if err != nil {
				t.Fatal(err)
			}
			testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
				valid := w.Code == http.StatusOK
				t.Log(w.Body)
				return withFalse || valid
			},
			)
		}
	}
}

func getTestWords() Words {
	var words Words
	for i := 1; i < countWords+1; i++ {
		s := strconv.Itoa(i)
		name := fmt.Sprintf("word #%s", s)
		transl := fmt.Sprint(name, " *transl")
		word := Word{
			Word:        name,
			Translation: transl,
		}
		words = append(words, word)
	}
	return words
}

func TestWordsAdd(t *testing.T) {
	urlPath := urlVersion + "/words/"
	r := getTestRouter()
	words := getTestWords()
	users := getTestUsers(true)
	for _, user := range users {
		token := logInUserTest(t, r, &user)
		userId := getUserIdFromToken(token)
		categories := getTestCategories(user.Email)
		categories = append(categories, Category{Name: "Fail cat"})
		for _, category := range categories {
			var currentCategory Category
			currentCategory.Name = category.Name
			db.Model(&User{ID: userId}).
				First(&currentCategory, &currentCategory)
			for _, words := range words {
				b, err := json.Marshal(&words)
				if err != nil {
					t.Fatal(err)
				}
				req, err := http.NewRequest("POST", urlPath, bytes.NewBuffer(b))
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authorization", token)
				if err != nil {
					t.Fatal(err)
				}
				if currentCategory.ID != 0 {
					q := req.URL.Query()
					q.Set("category_id", strconv.Itoa(int(currentCategory.ID)))
					req.URL.RawQuery = q.Encode()
				}

				testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
					valid := w.Code == http.StatusCreated
					t.Log(w.Body)
					return withFalse || valid
				},
				)
			}
		}
	}
}

func TestWordsPatch(t *testing.T) {
	urlPath := urlVersion + "/words/"
	r := getTestRouter()
	users := getTestUsers(true)
	for _, user := range users {
		token := logInUserTest(t, r, &user)
		var words Words
		db.Model(&User{ID: getUserIdFromToken(token)}).Related(&words)
		for _, word := range words {
			if word.ID%2 == 0 {
				continue
			}
			word.Word = fmt.Sprint(word.Word, " *patch")
			word.Translation = fmt.Sprint(word.Translation, " *patch")
			word.UserID = 0
			b, err := json.Marshal(&word)
			if err != nil {
				t.Fatal(err)
			}
			req, err := http.NewRequest("PATCH", urlPath, bytes.NewBuffer(b))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", token)
			if err != nil {
				t.Fatal(err)
			}
			testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
				valid := w.Code == http.StatusOK
				t.Log(w.Body)
				return withFalse || valid
			},
			)
		}
	}
}

func TestWordsDelete(t *testing.T) {
	urlPath := urlVersion + "/words/"
	r := getTestRouter()
	users := getTestUsers(true)
	for _, user := range users {
		token := logInUserTest(t, r, &user)
		var words Words
		db.Model(&User{ID: getUserIdFromToken(token)}).Related(&words)
		for _, word := range words {
			word.UserID = 0
			word.CategoryID = 0
			b, err := json.Marshal(&word)
			if err != nil {
				t.Fatal(err)
			}
			req, err := http.NewRequest("DELETE", urlPath, bytes.NewBuffer(b))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", token)
			if err != nil {
				t.Fatal(err)
			}
			testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
				valid := w.Code == http.StatusOK
				t.Log(w.Body)
				return withFalse || valid
			},
			)
		}
	}
}

func TestWordsGetAll(t *testing.T) {
	urlPath := urlVersion + "/words/?all"
	r := getTestRouter()
	users := getTestUsers(true)
	for _, user := range users {
		token := logInUserTest(t, r, &user)
		req, err := http.NewRequest("GET", urlPath, nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", token)
		if err != nil {
			t.Fatal(err)
		}
		testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
			valid := w.Code == http.StatusOK
			t.Log(w.Body)
			return withFalse || valid
		},
		)
	}
}

func TestWordsGetToCategory(t *testing.T) {
	urlPath := urlVersion + "/words/"
	r := getTestRouter()
	users := getTestUsers(true)
	for _, user := range users {
		token := logInUserTest(t, r, &user)
		userId := getUserIdFromToken(token)
		categories := getTestCategories(user.Email)
		categories = append(categories, Category{Name: "Fail cat"})
		for _, category := range categories {
			var currentCategory Category
			currentCategory.Name = category.Name
			db.Model(&User{ID: userId}).
				First(&currentCategory, &currentCategory)
			t.Log(currentCategory)
			req, err := http.NewRequest("GET", urlPath, nil)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", token)
			if err != nil {
				t.Fatal(err)
			}
			if currentCategory.ID != 0 {
				q := req.URL.Query()
				q.Set("category_id", strconv.Itoa(int(currentCategory.ID)))
				req.URL.RawQuery = q.Encode()
			}

			testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
				valid := w.Code == http.StatusOK
				t.Log(w.Body)
				return withFalse || valid
			},
			)
		}
	}
}
