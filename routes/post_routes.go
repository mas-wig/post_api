package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mas-wig/simple-api/controllers"
)

type PostRoutes struct {
	postController controllers.PostHandler
}

func NewPostRoutes(pc controllers.PostHandler) *PostRoutes {
	return &PostRoutes{postController: pc}
}

func (pc *PostRoutes) PostRouter(rg *gin.RouterGroup) {
	router := rg.Group("/post")
	router.POST("/", pc.postController.CreatePostHandler)
	router.GET("/", pc.postController.GetAllPosts)
	router.PATCH("/:id", pc.postController.UpdatePostHandler)
	router.GET("/:id", pc.postController.GetPostByID)
	router.DELETE("/:id", pc.postController.DeleteUserByID)
}
