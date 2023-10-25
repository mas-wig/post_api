package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mas-wig/simple-api/config"
	db "github.com/mas-wig/simple-api/db/sqlc"
	"github.com/mas-wig/simple-api/types"
	"github.com/mas-wig/simple-api/utils"
)

type AuthController struct {
	db  *db.Queries
	ctx context.Context
}

func NewAuthController(db *db.Queries, ctx context.Context) *AuthController {
	return &AuthController{db: db, ctx: ctx}
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
	if err := c.ShouldBindJSON(&credential); err != nil {
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

	c.Request.Header.Add("X-XSS-Protection", "1; mode=block")
	c.SetCookie("access_token", accessToken, config.AccessTokenMaxAge*60, "/", config.Origin, false, true)
	c.SetCookie("refresh_token", refreshToken, config.AccessTokenMaxAge*60, "/", config.Origin, false, true)
	c.SetCookie("logged_in", "true", config.AccessTokenMaxAge*60, "/", config.Origin, false, true)
	c.JSON(http.StatusOK, gin.H{"status": "success", "access token": accessToken})
}

func (a *AuthController) RefreshAccessToken(c *gin.Context) {
	cookie, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "status forbidden", "message": err.Error()})
		return
	}
	config, _ := config.LoadConfig()
	sub, err := utils.ValidateToken(cookie, config.RefreshTokenPublicKey)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "status forbidden", "message": err.Error()})
		return
	}
	user, err := a.db.GetUserById(a.ctx, uuid.MustParse(fmt.Sprint(sub)))
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "status forbidden", "message": err.Error()})
		return
	}
	accessToken, err := utils.CreateToken(config.AccessTokenExpiredIn, user.ID, config.AccessTokenPrivateKey)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "status forbidden", "message": err.Error()})
		return
	}
	c.SetCookie("access_token", accessToken, config.AccessTokenMaxAge, "/", config.Origin, false, true)
	c.SetCookie("logged_in", "true", config.AccessTokenMaxAge, "/", config.Origin, false, true)
	c.JSON(http.StatusOK, gin.H{"access token": accessToken})
}

func (ac *AuthController) LogoutUser(ctx *gin.Context) {
	config, _ := config.LoadConfig()
	ctx.SetCookie("access_token", "", -1, "/", config.Origin, false, true)
	ctx.SetCookie("refresh_token", "", -1, "/", config.Origin, false, true)
	ctx.SetCookie("logged_in", "", -1, "/", config.Origin, false, true)
	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}
