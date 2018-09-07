package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/studentofjs/moonrock-server/models"
)

// TGENewsletterHandler - signs up from PUT request with email to newsletter
func TGENewsletterHandler(c *gin.Context) {
	email := c.PostForm("email")
	code := models.TGENewsletter(email)
	c.JSON(code.ServerCode, gin.H{"status": code.Response})
}
