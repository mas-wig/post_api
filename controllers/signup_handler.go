package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/mas-wig/simple-api/db/sqlc"
	"github.com/mas-wig/simple-api/utils"
)

type AuthController struct {
	db *db.Queries
}

func NewAuthController(db *db.Queries) *AuthController {
	return &AuthController{db: db}
}

func (a *AuthController) SignUpUser(c *gin.Context) {
	var credential *db.User
	if err := c.ShouldBindJSON(&credential); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	hashedPassword := utils.HashPassword(credential.Password)

	args := &db.CreateUserParams{
		Name:      credential.Name,
		Email:     credential.Email,
		Photo:     credential.Photo,
		Verified:  true,
		Password:  hashedPassword,
		Role:      credential.Role,
		UpdatedAt: time.Now(),
	}
	user, err := a.db.CreateUser(c, *args)

	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "success", "data": gin.H{"user": user}})
}
