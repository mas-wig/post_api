package controllers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/mas-wig/simple-api/db/sqlc"
	"github.com/mas-wig/simple-api/types"
)

type UserHandler struct {
	db  *db.Queries
	ctx context.Context
}

func NewUserHandler(db *db.Queries, ctx context.Context) *UserHandler {
	return &UserHandler{db: db, ctx: ctx}
}

func (u *UserHandler) MyProfile(c *gin.Context) {
	currentUser := c.MustGet("current_user").(*types.SignInInput)
	c.JSON(http.StatusOK, gin.H{"status": "success", "user": currentUser})
}
