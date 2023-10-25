package controllers

import (
	"context"
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	db "github.com/mas-wig/simple-api/db/sqlc"
	"github.com/mas-wig/simple-api/types"
)

type PostHandler struct {
	db  *db.Queries
	ctx context.Context
}

func NewPostController(db *db.Queries, ctx context.Context) *PostHandler {
	return &PostHandler{db: db, ctx: ctx}
}

func (p *PostHandler) CreatePostHandler(c *gin.Context) {
	var payload *types.CreatePost
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": err.Error()})
		return
	}
	now := time.Now()
	args := &db.CreatePostParams{
		Title:     payload.Title,
		Content:   payload.Content,
		Category:  payload.Category,
		Image:     payload.Image,
		CreatedAt: now,
		UpdatedAt: now,
	}
	postSave, err := p.db.CreatePost(p.ctx, *args)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"status": "error bad gateway", "message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"status": "succcess", "message": postSave})
}

func (p *PostHandler) UpdatePostHandler(c *gin.Context) {
	var payload *types.UpdatePost
	postID := c.Param("id")
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": err.Error()})
		return
	}
	now := time.Now()

	args := &db.UpdatePostParams{
		Title:     sql.NullString{String: payload.Title, Valid: payload.Title != ""},
		Category:  sql.NullString{String: payload.Category, Valid: payload.Category != ""},
		Content:   sql.NullString{String: payload.Content, Valid: payload.Content != ""},
		Image:     sql.NullString{String: payload.Image, Valid: payload.Image != ""},
		UpdatedAt: sql.NullTime{Time: now, Valid: true},
		ID:        uuid.MustParse(postID),
	}

	updateUser, err := p.db.UpdatePost(p.ctx, *args)
	if err != nil {
		if err == sql.ErrNoRows {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"status": "status not found", "message": err.Error()})
			return
		}
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"status": "bad gateway", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success update data", "message": updateUser})
}

func (p *PostHandler) GetPostByID(c *gin.Context) {
	postID := c.Param("id")
	post, err := p.db.GetPostById(p.ctx, uuid.MustParse(postID))
	if err != nil {
		if err == sql.ErrNoRows {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"status": "status not found", "message": err.Error()})
			return
		}
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"status": "bad gateway", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success update data", "message": post})
}

func (p *PostHandler) GetAllPosts(c *gin.Context) {
	var (
		page  = c.DefaultQuery("page", "1")
		limit = c.DefaultQuery("limit", "10")
	)

	intPage, _ := strconv.Atoi(page)
	intLimit, _ := strconv.Atoi(limit)
	offset := (intPage - 1) * intLimit
	args := &db.ListPostsParams{Limit: int32(intLimit), Offset: int32(offset)}
	allPots, err := p.db.ListPosts(p.ctx, *args)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"status": "bad gateway", "message": err.Error()})
		return
	}

	if allPots != nil {
		allPots = []db.Post{}
	}
	c.JSON(http.StatusOK, gin.H{"status": "success update data", "message": allPots})
}

func (p *PostHandler) DeleteUserByID(c *gin.Context) {
	postID := c.Param("id")
	_, err := p.db.GetPostById(p.ctx, uuid.MustParse(postID))
	if err != nil {
		if err == sql.ErrNoRows {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"status": "not found", "message": "no post with exist id"})
			return
		}
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"status": "bad gateway", "message": err.Error()})
		return
	}
	err = p.db.DeletePost(p.ctx, uuid.MustParse(postID))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": "bad gateway", "message": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, gin.H{"status": "success"})
}
