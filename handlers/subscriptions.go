package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/studentofjs/moonrock-server/models"
)

type email struct {
	Email string `json:"email" binding:"required"`
}

// TGENewsletterHandler - signs up from PUT request with email to newsletter
func TGENewsletterHandler(c *gin.Context) {
	var json email
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	code := models.TGENewsletter(string(json.Email))
	c.JSON(code.ServerCode, gin.H{"status": code.Response})
}
