package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/asdine/storm"
	"github.com/gin-gonic/gin"
	"github.com/maxzerbini/oauth"
	"github.com/studentofjs/moonrock-server/database"
	"github.com/studentofjs/moonrock-server/handlers"
	"github.com/studentofjs/moonrock-server/models"
	"github.com/studentofjs/moonrock-server/secrets"
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
		secrets.SecretKey,
		time.Hour*120,
		&UserVerifier{},
		nil)
	router.POST("/token", s.UserCredentials)
	router.POST("/auth", s.ClientCredentials)

	authorized := router.Group("/u")
	// use the Bearer Athentication middleware
	authorized.Use(oauth.Authorize(secrets.SecretKey, nil))
	// update user's contribution address
	authorized.PUT("/address", handlers.ContributionAddressHandler)
	// get user's contribution address
	authorized.GET("/address", handlers.GetContributionAddressHandler)
	// update user details
	authorized.PUT("/update", handlers.UpdateUserDetailsHandler)

	router.PUT("/confirm", handlers.ConfirmAccountHandler)          // confirm user account
	router.POST("/register", handlers.RegisterHandler)              // register user account
	router.PUT("/reset_password", handlers.ResetPasswordHandler)    // reset password action
	router.POST("/forgot_password", handlers.ForgotPasswordHandler) // forgot password process
	router.POST("/tgenews", handlers.TGENewsletterHandler)          // signup to token sale news
}

// UserVerifier provides user credentials verifier
type UserVerifier struct {
}

// ValidateUser validates username and password returning an error if the user credentials are wrong
func (*UserVerifier) ValidateUser(username, password, scope string, req *http.Request) error {
	err := models.LoginCheck(username, password)
	return err
}

// ValidateClient validates clientId and secret returning an error if the client credentials are wrong
func (*UserVerifier) ValidateClient(clientID, clientSecret, scope string, req *http.Request) error {
	if scope != "write:subscription" {
		err := errors.New("invalid")
		return err
	}
	err := models.LoginCheck(clientID, clientSecret)
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
	db, err := database.OpenProdDB("./database/")
	if err != nil {
		return errors.New("database error")
	}
	defer db.Close()

	token := Token{
		Credential:     credential,
		RefreshTokenID: refreshTokenId,
		TokenID:        tokenID,
		TokenType:      tokenType,
	}

	if err := db.Save(&token); err == storm.ErrAlreadyExists {
		if err := db.Update(&token); err != nil {
			return err
		}

	}
	return nil
}

// AddProperties provides additional information to the token response
func (*UserVerifier) AddProperties(credential, tokenID, tokenType string, scope string) (map[string]string, error) {
	db, err := database.OpenProdDB("./database/")
	if err != nil {

		return nil, errors.New("database error")
	}
	defer db.Close()

	var user models.User

	if err := db.One("Username", credential, &user); err != nil {
		return nil, errors.New("invalid login")
	}

	props := make(map[string]string)
	switch scope {
	case "write:subscription":
		props["access_type"] = "client"
		props["permission"] = "write"
	case "write:registration":
		props["access"] = "client"
		props["permission"] = "write"
	case "write:user read:user delete:user":
		props["access_type"] = "auth"
		props["permission"] = "read write delete"
		props["id"] = string(user.ID)
		props["firstname"] = user.FirstName
		props["lastname"] = user.LastName
		props["country"] = user.CountryCode
	default:
		props["access_type"] = "read-only"
		props["permission"] = "read"
	}

	return props, nil
}

// ValidateTokenId validates token Id
func (*UserVerifier) ValidateTokenId(credential, tokenId, refreshTokenID, tokenType string) error {
	db, err := database.OpenProdDB("./database/")
	if err != nil {
		return errors.New("database error")
	}
	defer db.Close()
	var token Token
	if err = db.One("TokenID", tokenId, &token); err != nil {
		return err
	}
	if credential != token.Credential || refreshTokenID != token.RefreshTokenID || tokenType != tokenType {
		return errors.New("invalid token")
	}
	return nil
}

// ValidateCode validates token Id - currently only partially implemented so will skip
func (*UserVerifier) ValidateCode(clientID, clientSecret, code, redirectURI string, req *http.Request) (string, error) {
	fmt.Println(redirectURI)
	return "", nil
}
