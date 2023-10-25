package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mas-wig/simple-api/controllers"
	db "github.com/mas-wig/simple-api/db/sqlc"
	"github.com/mas-wig/simple-api/middleware"
)

type AuthRoutes struct {
	handler controllers.AuthController
	db      *db.Queries
}

func NewAuthRoutes(handler controllers.AuthController, db *db.Queries) *AuthRoutes {
	return &AuthRoutes{handler: handler, db: db}
}

func (a *AuthRoutes) AuthRouteUser(rg *gin.RouterGroup) {
	routes := rg.Group("/auth")
	routes.POST("/register", a.handler.SignUpUser)
	routes.POST("/login", a.handler.SignInUser)
	routes.GET("/refresh", a.handler.RefreshAccessToken)
	routes.GET("/logut", middleware.DeserializeUser(a.db), a.handler.LogoutUser)
}
