package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
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
	code := models.UpdateContributionAddress(id, e)
	if code.ServerCode == 200 {
		c.JSON(200, gin.H{
			"status":   "updated",
			"ethereum": e,
		})
	} else {
		c.JSON(code.ServerCode, gin.H{"status": code.Response})
	}
}

// ConfirmAccountHandler checks a resetCode against the DB and returns an error string or
func ConfirmAccountHandler(c *gin.Context) {
	resetcode := c.PostForm("resetcode")
	code := models.ConfirmAccount(resetcode)
	c.JSON(code.ServerCode, gin.H{"status": code.Response})
}

// ForgotPasswordHandler sends a reset email with unique password reset link
func ForgotPasswordHandler(c *gin.Context) {
	username := c.PostForm("username")
	code := models.ForgotPassword(username)
	c.JSON(code.ServerCode, gin.H{"status": code.Response})
}

// GetContributionAddressHandler returns the saved address of the user
func GetContributionAddressHandler(c *gin.Context) {
	i := c.PostForm("id")
	eth, code := models.GetContributionAddress(i)
	if eth != "" {
		c.JSON(200, gin.H{
			"status":   "ok",
			"ethereum": eth,
		})
	} else {
		c.JSON(code.ServerCode, gin.H{"status": code.Response})
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
	if code.ServerCode == 200 {
		c.JSON(200, gin.H{
			"status":    "updated",
			"address":   address,
			"country":   country,
			"ethereum":  ethereum,
			"firstName": firstname,
			"lastName":  lastname,
		})
	} else {
		c.JSON(code.ServerCode, gin.H{"status": code.Response})
	}
}

// ResetPasswordHandler handles the reset code checking and password change
func ResetPasswordHandler(c *gin.Context) {
	password := c.PostForm("password")
	resetcode := c.PostForm("resetcode")
	username := c.PostForm("username")
	code := models.ResetPassword(password, resetcode, username)
	c.JSON(code.ServerCode, gin.H{"status": code.Response})
}

// UpdateUserDetailsHandler updates user details supplied to API
func UpdateUserDetailsHandler(c *gin.Context) {
	// @todo check for values
	a := c.PostForm("address")
	cc := c.PostForm("country")
	f := c.PostForm("firstname")
	i := c.PostForm("id")
	l := c.PostForm("lastname")
	code := models.UpdateUserDetails(a, cc, f, i, l)
	if code.ServerCode == 200 {
		c.JSON(200, gin.H{
			"status":    "updated",
			"address":   a,
			"country":   cc,
			"firstName": f,
			"lastName":  l,
		})
	} else {
		c.JSON(code.ServerCode, gin.H{"status": code.Response})
	}
}
