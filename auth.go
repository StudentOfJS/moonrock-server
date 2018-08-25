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
	s := oauth.NewOAuthBearerServer(
		SecretKey,
		time.Second*120,
		&UserVerifier{},
		nil)
	router.POST("/token", s.UserCredentials)
	router.POST("/auth", s.ClientCredentials)

	authorized := router.Group("/")
	// use the Bearer Athentication middleware
	authorized.Use(oauth.Authorize(SecretKey, nil))
	authorized.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	authorized.PUT("/tgenews", TokenSaleUpdatesHandler)
	authorized.PUT("/family/", TokenSaleUpdatesHandler)
	authorized.PUT("/tgenews", TokenSaleUpdatesHandler)
	authorized.PUT("/tgenews", TokenSaleUpdatesHandler)

}

// UserVerifier provides user credentials verifier for testing.
type UserVerifier struct {
}

// ValidateUser validates username and password returning an error if the user credentials are wrong
func (*UserVerifier) ValidateUser(username, password, scope string, req *http.Request) error {
	if err := LoginCheck(username, password); err == nil {
		return nil
	}
	return errors.New("invalid login")
}

// ValidateClient validates clientId and secret returning an error if the client credentials are wrong
func (*UserVerifier) ValidateClient(clientID, clientSecret, scope string, req *http.Request) error {
	if clientID == "abcdef" && clientSecret == "12345" {
		return nil
	}
	return errors.New("Wrong client")
}

// AddClaims provides additional claims to the token
func (*UserVerifier) AddClaims(credential, tokenID, tokenType, scope string) (map[string]string, error) {
	claims := make(map[string]string)
	claims["customerId"] = "1001"
	claims["customerData"] = `{"OrderDate":"2016-12-14","OrderId":"9999"}`
	return claims, nil
}

// StoreTokenId saves the token Id generated for the user
func (*UserVerifier) StoreTokenId(credential, tokenID, refreshTokenId, tokenType string) error {
	return nil
}

// AddProperties provides additional information to the token response
func (*UserVerifier) AddProperties(credential, tokenID, tokenType string, scope string) (map[string]string, error) {
	props := make(map[string]string)
	props["customerName"] = "Gopher"
	return props, nil
}

// ValidateTokenId validates token Id
func (*UserVerifier) ValidateTokenId(credential, tokenId, refreshTokenID, tokenType string) error {
	return nil
}

// ValidateCode validates token Id
func (*UserVerifier) ValidateCode(clientID, clientSecret, code, redirectURI string, req *http.Request) (string, error) {
	return "", nil
}
