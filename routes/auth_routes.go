package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mas-wig/simple-api/controllers"
)

type AuthRoutes struct {
	handler controllers.AuthController
}

func NewAuthRoutes(handler controllers.AuthController) *AuthRoutes {
	return &AuthRoutes{handler: handler}
}

func (a *AuthRoutes) AuthRouteUser(rg *gin.RouterGroup) {
	routes := rg.Group("/auth")
	routes.POST("/register", a.handler.SignUpUser)
}
