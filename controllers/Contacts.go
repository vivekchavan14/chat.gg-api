package controllers

import (
	"net/http"
	"server/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ContactsController struct {
	DB *gorm.DB
}

func (cc *ContactsController) CheckRegisteredContacts(c *gin.Context) {
	var req struct {
		PhoneNumbers []string `json:"phone_numbers"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(req.PhoneNumbers) == 0 {
		c.JSON(http.StatusOK, []models.User{})
		return
	}

	var registeredUsers []models.User
	if err := cc.DB.Where("phone_number IN ?", req.PhoneNumbers).Find(&registeredUsers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query registered contacts"})
		return
	}

	// Prepare response by excluding sensitive information
	var contacts []struct {
		Username    string `json:"username"`
		PhoneNumber string `json:"phone_number"`
	}

	for _, user := range registeredUsers {
		contacts = append(contacts, struct {
			Username    string `json:"username"`
			PhoneNumber string `json:"phone_number"`
		}{
			Username:    user.Username,
			PhoneNumber: user.PhoneNumber,
		})
	}

	c.JSON(http.StatusOK, contacts)
}
