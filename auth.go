package main

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/asdine/storm"
	"github.com/gin-gonic/gin"
	"github.com/maxzerbini/oauth"
)

// Token is the struct of the users login token
type Token struct {
	Credential     string // this field will not be indexed
	RefreshTokenID string // this field will not be indexed
	TokenID        string `storm:"id"` // primary key
	TokenType      string // this field will not be indexed
}

// RegisterAPI registers api endpoints with the auth middleware
func RegisterAPI(router *gin.Engine) {
	s := oauth.NewOAuthBearerServer(
		SecretKey,
		time.Hour*120,
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
	authorized.PUT("/family", TokenSaleUpdatesHandler)

}

// UserVerifier provides user credentials verifier
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
	if scope != "write:subscription" {
		err := errors.New("invalid")
		return err
	}
	err := LoginCheck(clientID, clientSecret)
	return err
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
	// Start boltDB
	db, err := storm.Open("my.db")
	if err != nil {
		log.Fatal(err)
	}

	token := Token{
		Credential:     credential,
		RefreshTokenID: refreshTokenId,
		TokenID:        tokenID,
		TokenType:      tokenType,
	}

	if err := db.Save(&token); err == storm.ErrAlreadyExists {
		err := db.Update(&token{ID: 10, Name: "Jack", Age: 45})

	}
	defer db.Close()

	return nil
}

// AddProperties provides additional information to the token response
func (*UserVerifier) AddProperties(credential, tokenID, tokenType string, scope string) (map[string]string, error) {
	props := make(map[string]string)
	switch scope {
	case "write:subscription":
		props["access_type"] = "client-only"
		props["permission"] = "write"
	case "write:registration":
		props["access"] = "client-only"
		props["permission"] = "write"
	case "write:user read:user delete:user":
		props["access_type"] = "auth-only"
		props["permission"] = "read write delete"
	default:
		props["access_type"] = "read-only"
		props["permission"] = "read"
	}

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
