package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mas-wig/simple-api/config"
	db "github.com/mas-wig/simple-api/db/sqlc"
	"github.com/mas-wig/simple-api/utils"
)

func DeserializeUser(db *db.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			accessToken         string
			cookie, err         = c.Cookie("access_token")
			authorizationHeader = c.Request.Header.Get("Authorization")
			field               = strings.Fields(authorizationHeader)
			config, _           = config.LoadConfig()
		)
		if len(field) != 0 && field[0] == "Bearer" {
			accessToken = field[1]
		} else if err != nil {
			accessToken = cookie
		}

		if accessToken == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
			return
		}

		sub, err := utils.ValidateToken(accessToken, config.AccessTokenPublicKey)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "unauthorized", "message": err.Error()})
			return
		}

		user, err := db.GetUserById(context.TODO(), uuid.MustParse(fmt.Sprint(sub)))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "unauthorized", "message": "token not exits"})
			return
		}
		c.Set("current_user", user)
		c.Next()
	}
}
