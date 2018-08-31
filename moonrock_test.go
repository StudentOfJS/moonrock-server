package main

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/asdine/storm"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var e = "test@test.com.au"
var eth = "0xCaE9eFE97895EF43e72791a10254d6abDdb17Ae9"
var f = "Rod"
var l = "Lewis"
var p = "djfiejij5453io5tgrtegdJ"
var rc = "F0001234-0451-4000-B000-000000000000"

/* ------------------- utils ------------------------- */
func TestHashPassword(t *testing.T) {
	hash, err := HashPassword(p)
	if err != nil {
		t.Errorf("Expected a hash but recieved an error: %s", err)
	} else {
		hashType := reflect.TypeOf(hash).String()
		if hashType != "[]uint8" {
			t.Error("Expected a []uint8, but recieved: " + hashType)
		}
	}
}

func TestLoginValid(t *testing.T) {
	if err := LoginValid(e, p); err != nil {
		t.Errorf("Provided valid username and password got: %s", err)
	}
	if err := LoginValid(e, "x"); err == nil {
		t.Error("Provided invalid password, expected login check to fail, but it passed")
	}
	if err := LoginValid("username_not_valid", p); err == nil {
		t.Error("Provided invalid username, expected login check to fail but it passed")
	}
}

func TestUserValid(t *testing.T) {
	if err := UserValid(eth, f, l); err != nil {
		t.Errorf("Provided valid user details, but recieved: %s", err)
	}
	if err := UserValid("not_valid", f, l); err == nil {
		t.Error("Provided invalid eth address, but check passed")
	}
	if err := UserValid(eth, "12132", l); err == nil {
		t.Error("Provided invalid name, but check passed")
	}
}

func TestEmailValid(t *testing.T) {
	if err := EmailValid(e); err != nil {
		t.Errorf("Provided valid email, but recieved: %s", err)
	}
}

func TestCreateUUID(t *testing.T) {
	id, err := CreateUUID(rc)
	if err != nil {
		t.Errorf("Provided valid string, expected id, but recieved: %s", err)
	} else {
		if reflect.TypeOf(id).String() != "uuid.UUID" {
			t.Error("Expected type of id to be uuid.UUID")
		}
	}
	if _, err := CreateUUID("not_valid_string"); err == nil {
		t.Error("Provided invalid string, expected error, but recieved nil")
	}
}

func TestLoginCheck(t *testing.T) {
	if err := LoginCheck(e, "invalid_login"); err == nil {
		t.Error("Provided invalid login details, expected error, but recieved nil")
	}
	if err := LoginCheck(TestUser, TestPass); err.Error() != "confirm email" {
		t.Errorf("Provided valid login details, but recieved: %s", err)
	}
}

/* ------------------- Database ------------------------- */

func TestDB(t *testing.T) {
	// Start boltDB
	db, err := storm.Open("my.db")
	defer db.Close()
	if err != nil {
		t.Errorf("Database failed to open: %v", err)
		return
	}
}

/* ------------------- API ------------------------- */

func router() *gin.Engine {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
	r.POST("/confirm", ConfirmAccountHandler)
	r.POST("/register", RegisterHandler)
	r.POST("/reset_password", ForgotPasswordHandler)
	r.POST("/forgot_password", ForgotPasswordHandler)
	r.POST("/tgenews", TokenSaleUpdatesHandler)
	return r
}

func removeTestUser() error {
	db, err := storm.Open("my.db")
	defer db.Close()
	if err != nil {
		return errors.New("Database failed to open")
	}
	var user User
	err = db.One("EthereumAddress", eth, &user)
	if err != nil {
		return errors.New("test user not in db")
	}
	if err := db.DeleteStruct(&user); err != nil {
		return errors.New("can't delete test user")
	}
	return nil
}

func registerUser() *httptest.ResponseRecorder {
	r := router()
	w := httptest.NewRecorder()
	p := getRegisterPOSTPayload()
	req, _ := http.NewRequest("POST", "/register", strings.NewReader(p))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(p)))
	r.ServeHTTP(w, req)
	return w
}

func getTgeNewsSignupPOSTPayload() string {
	params := url.Values{}
	params.Add("email", e)
	return params.Encode()
}

func getRegisterPOSTPayload() string {
	params := url.Values{}
	params.Add("address", "test address in the test village")
	params.Add("country", "AU")
	params.Add("ethereum", eth)
	params.Add("firstname", f)
	params.Add("firstname", l)
	params.Add("password", p)
	params.Add("resetcode", rc)
	params.Add("username", e)
	return params.Encode()
}

func getForgotPasswordPOSTPayload() string {
	params := url.Values{}
	params.Add("username", e)
	return params.Encode()
}

func TestPingRoute(t *testing.T) {
	r := router()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "pong", w.Body.String())
}

func TestSignupForTokenSaleNews(t *testing.T) {
	r := router()
	w := httptest.NewRecorder()
	p := getTgeNewsSignupPOSTPayload()
	req, _ := http.NewRequest("POST", "/tgenews", strings.NewReader(p))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(p)))

	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

func TestRegisterUser(t *testing.T) {
	w := registerUser()
	assert.Equal(t, 200, w.Code)
	removeTestUser()
}

func TestForgotPassword(t *testing.T) {
	uw := registerUser()
	assert.Equal(t, 200, uw.Code)
	r := router()
	w := httptest.NewRecorder()
	p := getForgotPasswordPOSTPayload()
	req, _ := http.NewRequest("POST", "/forgot_password", strings.NewReader(p))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(p)))

	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	removeTestUser()
}
