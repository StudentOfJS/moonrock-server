package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/maxzerbini/oauth"
)

// RegisterAPI registers api endpoints with the auth middleware
func RegisterAPI(router *gin.Engine) {
	authorized := router.Group("/")
	// use the Bearer Athentication middleware
	authorized.Use(oauth.Authorize("U#cXHY_w4Vg$FCaJ7-jtjr##xMrmgydy", nil))
	authorized.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	authorized.PUT("/tgenews", TokenSaleUpdatesHandler)

	s := oauth.NewOAuthBearerServer(
		"U#cXHY_w4Vg$FCaJ7-jtjr##xMrmgydy",
		time.Second*120,
		&TestUserVerifier{},
		nil)
	router.POST("/token", s.UserCredentials)
	router.POST("/auth", s.ClientCredentials)

}

// TestUserVerifier provides user credentials verifier for testing.
type TestUserVerifier struct {
	username string
	password string
}

// ValidateUser validates username and password returning an error if the user credentials are wrong
func (*TestUserVerifier) ValidateUser(username, password, scope string, req *http.Request) error {
	if username == "user01" && password == "12345" {
		return nil
	}
	return errors.New("Wrong user")
}

// ValidateClient validates clientId and secret returning an error if the client credentials are wrong
func (*TestUserVerifier) ValidateClient(clientID, clientSecret, scope string, req *http.Request) error {
	if clientID == "abcdef" && clientSecret == "12345" {
		return nil
	}
	return errors.New("Wrong client")
}

// AddClaims provides additional claims to the token
func (*TestUserVerifier) AddClaims(credential, tokenID, tokenType, scope string) (map[string]string, error) {
	claims := make(map[string]string)
	claims["customerId"] = "1001"
	claims["customerData"] = `{"OrderDate":"2016-12-14","OrderId":"9999"}`
	return claims, nil
}

// StoreTokenID saves the token Id generated for the user
func (*TestUserVerifier) StoreTokenID(credential, tokenID, refreshTokenID, tokenType string) error {
	return nil
}

// AddProperties provides additional information to the token response
func (*TestUserVerifier) AddProperties(credential, tokenID, tokenType string, scope string) (map[string]string, error) {
	props := make(map[string]string)
	props["customerName"] = "Gopher"
	return props, nil
}

// ValidateTokenID validates token Id
func (*TestUserVerifier) ValidateTokenID(credential, tokenID, refreshTokenID, tokenType string) error {
	return nil
}

// ValidateCode validates token Id
func (*TestUserVerifier) ValidateCode(clientID, clientSecret, code, redirectURI string, req *http.Request) (string, error) {
	return "", nil
}
