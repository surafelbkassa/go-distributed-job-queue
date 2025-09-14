package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// TODO: Wire up usecase and repository for actual registration logic
func RegisterHandler(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	// Here you would call your usecase to register the user
	// For now, just return success for demonstration
	c.JSON(http.StatusCreated, gin.H{"message": "user registered"})
}

func R() {

}
