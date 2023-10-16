package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mas-wig/simple-api/config"
	"github.com/mas-wig/simple-api/controllers"

	"github.com/mas-wig/simple-api/routes"

	_ "github.com/lib/pq"
)

var (
	server         *gin.Engine
	db             *dbConn.Queries
	AuthController controllers.AuthController
	AuthRoutes     routes.AuthRoutes
)

func init() {
	config, _ := config.LoadConfig()
	conn, err := sql.Open(config.PostgreDriver, config.PostgreURI)
	if err != nil {
		log.Fatal("tidak bisa melakukan koneksi kedalam database")
	}
	db = dbConn.New(conn)
	fmt.Println("PostgreSQL connected successfully...")

	AuthController = *controllers.NewAuthController(db)
	AuthRoutes = *routes.NewAuthRoutes(AuthController)
	server = gin.Default()
}

func main() {
	config, _ := config.LoadConfig()
	router := server.Group("/api")

	AuthRoutes.AuthRouteUser(router)

	router.GET("/healtchecker", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "siap onfire"})
	})
	log.Fatal(server.Run(":" + config.PORT))
}
