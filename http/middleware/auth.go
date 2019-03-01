package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/totoval/framework/config"
	"github.com/totoval/framework/utils/jwt"
	"net/http"
	"strings"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.DefaultQuery("token", "")
		if token == "" {
			token = c.Request.Header.Get("Authorization")
			if s := strings.Split(token, " "); len(s) == 2 {
				token = s[1]
			}
		}

		j := jwt.NewJWT(config.GetString("auth.sign_key"))
		claims, err := j.ParseToken(token)
		if err != nil {
			if err == jwt.TokenExpired {
				if token, err = j.RefreshToken(token); err == nil {
					c.Header("Authorization", "Bear "+token)
					c.JSON(http.StatusOK, gin.H{"error": 0, "message": "refresh token", "token": token})
					return
				}
			}
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": 1, "message": err.Error()})
			return
		}
		c.Set("claims", claims)
	}
}
