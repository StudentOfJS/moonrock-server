package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/asdine/storm"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/studentofjs/moonrock-server/models"
)

var e = "test@test.com.au"
var eth = "0xCaE9eFE97895EF43e72791a10254d6abDdb17Ae9"
var f = "Rod"
var l = "Lewis"
var p = "djfiejij5453io5tgrtegdJ"
var rc = "F0001234-0451-4000-B000-000000000000"

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
	r.POST("/tgenews", TGENewsletterHandler)
	return r
}

func removeTestUser() error {
	db, err := storm.Open("my.db")
	defer db.Close()
	if err != nil {
		return errors.New("Database failed to open")
	}
	var user models.User
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
	params.Add("username", e)
	return params.Encode()
}

func getForgotPasswordPOSTPayload() string {
	params := url.Values{}
	params.Add("username", e)
	return params.Encode()
}

func getConfirmAccountPOSTPayload(r string) string {
	params := url.Values{}
	params.Add("resetcode", r)
	return params.Encode()
}

func getLoginPOSTPayload() string {
	params := url.Values{}
	params.Add("password", p)
	params.Add("username", e)
	params.Add("scope", "write:user read:user delete:user")
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

func TestConfirmAccount(t *testing.T) {
	uw := registerUser()
	assert.Equal(t, 200, uw.Code)
	user := models.User{}
	if err := json.NewDecoder(uw.Body).Decode(&user); err != nil {
		t.Fatalf("decoding failed")
	}

	r := router()
	w := httptest.NewRecorder()
	p := getConfirmAccountPOSTPayload(user.ResetCode.String())
	req, _ := http.NewRequest("POST", "/confirm", strings.NewReader(p))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(p)))
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	removeTestUser()
}
