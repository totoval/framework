package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/totoval/framework/config"
	"github.com/totoval/framework/utils/jwt"
	"net/http"
	"strings"
	"unicode/utf8"
)

const (
	CLAIM_KEY = "CLAIM"
	TOKEN_KEY = "TOKEN"
)
type TokenRevokeError struct {}
func (e TokenRevokeError) Error() string{
	return "token revoke failed"
}
type UserNotLoginError struct{}
func (e UserNotLoginError) Error() string {
	return "user not login"
}


func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.DefaultQuery("token", "")
		if token == "" {
			token = c.Request.Header.Get("Authorization")
			if s := strings.Split(token, " "); len(s) == 2 {
				token = s[1]
			}
		}

		// set token
		c.Set(TOKEN_KEY, token)

		j := jwt.NewJWT(config.GetString("auth.sign_key"))
		claims, err := j.ParseToken(token)
		if err != nil {
			if err == jwt.TokenExpired {
				if token, err = j.RefreshTokenUnverified(token); err == nil {
					if claims, err := j.ParseToken(token); err == nil {
						c.Set(CLAIM_KEY, claims)
						c.Header("Authorization", "Bear "+token)
						//c.JSON(http.StatusOK, gin.H{"data": gin.H{"token": token}})
						return
					}
				}
			}
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.Set(CLAIM_KEY, claims)
	}
}

func AuthClaimsID(c *gin.Context) (ID uint, isAbort bool) {
	claims, exist := c.Get(CLAIM_KEY)
	if !exist {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, UserNotLoginError{})
		return 0, true
	}

	r, _ := utf8.DecodeRune([]byte(claims.(*jwt.UserClaims).ID))
	return uint(r), false
}


func Revoke(c *gin.Context) error {
	j := jwt.NewJWT(config.GetString("auth.sign_key"))
	if tokenString, exist := c.Get(TOKEN_KEY); exist {
		if token, ok := tokenString.(string); ok {
			if err := j.RevokeToken(token); err == nil{
				c.Header("Authorization", "")
				return nil
			}
		}
	}
	return TokenRevokeError{}
}
