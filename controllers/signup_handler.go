package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mas-wig/simple-api/config"
	db "github.com/mas-wig/simple-api/db/sqlc"
	"github.com/mas-wig/simple-api/types"
	"github.com/mas-wig/simple-api/utils"
)

type AuthController struct {
	db  *db.Queries
	ctx context.Context
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

	userResponse := types.FilteredUserResponse(user)
	c.JSON(http.StatusCreated, gin.H{"status": "success", "data": gin.H{"user": userResponse}})
}

func (a *AuthController) SignInUser(c *gin.Context) {
	var credential *types.SignInInput
	if err := c.ShouldBindUri(&credential); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request", "message": err.Error()})
		return
	}

	user, err := a.db.GetUserByEmail(a.ctx, credential.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request", "message": err.Error()})
		return
	}

	if err := utils.ComparePassword(user.Password, credential.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request", "message": err.Error()})
		return
	}
	config, _ := config.LoadConfig()

	accessToken, err := utils.CreateToken(time.Duration(config.AccessTokenMaxAge), user, config.AccessTokenPrivateKey)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request", "message": err.Error()})
		return
	}
	refreshToken, err := utils.CreateToken(time.Duration(config.RefreshTokenMaxAge), user, config.RefreshTokenPrivateKey)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request", "message": err.Error()})
		return
	}

	c.SetCookie("access_token", accessToken, config.AccessTokenMaxAge*60, "/", config.Origin, false, true)
	c.SetCookie("refresh_token", refreshToken, config.AccessTokenMaxAge*60, "/", config.Origin, false, true)
	c.SetCookie("logged_in", "true", config.AccessTokenMaxAge*60, "/", config.Origin, false, true)

	c.JSON(http.StatusOK, gin.H{"status": "success", "access token": accessToken})
}

// TODO: implement this shit tomorrow!!
func (a *AuthController) RefreshAccessToken(c *gin.Context) {

}
