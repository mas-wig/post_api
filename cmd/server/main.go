package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/mas-wig/simple-api/config"
	"github.com/mas-wig/simple-api/controllers"
	"github.com/mas-wig/simple-api/routes"

	_ "github.com/lib/pq"
	dbConn "github.com/mas-wig/simple-api/db/sqlc"
)

var (
	server         *gin.Engine
	db             *dbConn.Queries
	AuthController controllers.AuthController
	AuthRoutes     routes.AuthRoutes
	UserController controllers.UserHandler
	UserRoutes     routes.UserRoutes
	PostController controllers.PostHandler
	PostRoutes     routes.PostRoutes
	ctx            context.Context
)

func init() {
	config, _ := config.LoadConfig()
	conn, err := sql.Open(config.PostgreDriver, config.PostgreURI)
	if err != nil {
		log.Fatal("tidak bisa melakukan koneksi kedalam database")
	}
	db = dbConn.New(conn)

	fmt.Println("PostgreSQL connected successfully...")

	ctx = context.TODO()

	AuthController = *controllers.NewAuthController(db, ctx)
	AuthRoutes = *routes.NewAuthRoutes(AuthController, db)
	UserController = *controllers.NewUserHandler(db, ctx)
	UserRoutes = *routes.NewUserRoutes(UserController, db)
	PostController = *controllers.NewPostController(db, ctx)
	PostRoutes = *routes.NewPostRoutes(PostController)
	server = gin.Default()
}

func main() {
	config, _ := config.LoadConfig()
	server.Use(
		cors.New(cors.Config{
			AllowAllOrigins: false,
			AllowOrigins:    []string{config.Origin},
			AllowHeaders:    []string{"Origin"},
			AllowMethods:    []string{"GET", "POST"},
		}),
	)
	router := server.Group("/api")
	AuthRoutes.AuthRouteUser(router)
	UserRoutes.UserRouter(router)
	PostRoutes.PostRouter(router)

	router.GET("/healtchecker", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "siap onfire"})
	})
	log.Fatal(server.Run(":" + config.PORT))
}
