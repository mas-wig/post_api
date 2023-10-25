package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mas-wig/simple-api/controllers"
	db "github.com/mas-wig/simple-api/db/sqlc"
	"github.com/mas-wig/simple-api/middleware"
)

type UserRoutes struct {
	userController controllers.UserHandler
	db             *db.Queries
}

func NewUserRoutes(userController controllers.UserHandler, db *db.Queries) *UserRoutes {
	return &UserRoutes{userController, db}
}

func (u *UserRoutes) UserRouter(rg *gin.RouterGroup) {
	routes := rg.Group("/users")
	routes.GET("/me", middleware.DeserializeUser(u.db), u.userController.MyProfile)
}
