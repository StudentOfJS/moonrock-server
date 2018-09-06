package handlers

import (
	"strconv"

	"github.com/asdine/storm"
	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
	"github.com/studentofjs/moonrock-server/models"
)

// ContributionAddressHandler uses an ID to find user and updates their contribution address
func ContributionAddressHandler(c *gin.Context) {
	e := c.PostForm("ethereum")
	idStr := c.PostForm("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.String(401, "unauthenticated")
		return
	}

	if !models.UpdateContributionAddress(id, e) {
		c.String(400, "update failed")
		return
	}

	c.JSON(200, gin.H{
		"status":   "updated",
		"ethereum": e,
	})
	return
}

// ConfirmAccountHandler checks a resetCode against the DB and returns an error string or
func ConfirmAccountHandler(c *gin.Context) {
	resetcode := c.PostForm("resetcode")
	code := models.ConfirmAccount(resetcode)
	c.JSON(code.serverCode, gin.H{"status": code.response})
}

// ForgotPasswordHandler sends a reset email with unique password reset link
func ForgotPasswordHandler(c *gin.Context) {
	username := c.PostForm("username")
	code := models.ForgotPassword(username)
	c.JSON(code.serverCode, gin.H{"status": code.response})
}

// GetContributionAddressHandler returns the saved address of the user
func GetContributionAddressHandler(c *gin.Context) {
	i := c.PostForm("id")
	eth, code := models.GetContributionAddress(i)
	if eth != nil {
		c.JSON(200, gin.H{
			"status":   "ok",
			"ethereum": eth,
		})
	} else {
		c.JSON(code.serverCode, gin.H{"status": code.response})
	}
}

// RegisterHandler validates the user signup form and saves to db
func RegisterHandler(c *gin.Context) {
	address := c.PostForm("address")
	country := c.PostForm("country")
	ethereum := c.PostForm("ethereum")
	firstname := c.PostForm("firstname")
	lastname := c.PostForm("lastname")
	password := c.PostForm("password")
	username := c.PostForm("username")
	code := models.Register(address, country, ethereum, firstname, lastname, password, username)
	if code.serverCode == 200 {
		c.JSON(200, gin.H{
			"status":    "updated",
			"address":   address,
			"country":   country,
			"ethereum":  ethereum,
			"firstName": firstname,
			"lastName":  lastname,
		})
	} else {
		c.JSON(code.serverCode, gin.H{"status": code.response})
	}
}

// ResetPasswordHandler handles the reset code checking and password change
func ResetPasswordHandler(c *gin.Context) {
	password := c.PostForm("password")
	resetcode := c.PostForm("resetcode")
	username := c.PostForm("username")
	code := ResetPassword(p, r, u string)
	c.JSON(code.serverCode, gin.H{"status": code.response})
}

// TGENewsletterHandler - signs up from PUT request with email to newsletter
func TGENewsletterHandler(c *gin.Context) {
	email := c.PostForm("email")
	code := TGENewsletter(email)
	c.JSON(code.serverCode, gin.H{"status": code.response})
}

// UpdateUserDetailsHandler updates user details supplied to API
func UpdateUserDetailsHandler(c *gin.Context) {
	// @todo check for values
	a := c.PostForm("address")
	c := c.PostForm("country")
	f := c.PostForm("firstname")
	i := c.PostForm("id")
	l := c.PostForm("lastname")
	code := UpdateUserDetails(a, c, f, i, l)
	if code.serverCode == 200 {
		c.JSON(200, gin.H{
			"status":    "updated",
			"address":   address,
			"country":   country,
			"firstName": firstname,
			"lastName":  lastname,
		})
	} else {
		c.JSON(code.serverCode, gin.H{"status": code.response})
	}
}
