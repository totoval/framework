package middleware

import (
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/gin-gonic/gin"

	"github.com/totoval/framework/config"
	"github.com/totoval/framework/model"
	"github.com/totoval/framework/utils/jwt"
)

const (
	CLAIM_KEY = "CLAIM"
	TOKEN_KEY = "TOKEN"
)

type TokenRevokeError struct{}

func (e TokenRevokeError) Error() string {
	return "token revoke failed"
}

type UserNotLoginError struct{}

func (e UserNotLoginError) Error() string {
	return "user not login"
}

type UserNotExistError struct{}

func (e UserNotExistError) Error() string {
	return "user not exist"
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
				if token, _err := j.RefreshTokenUnverified(token); _err == nil {
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

func authClaimID(c *gin.Context) (ID uint, exist bool) {
	claims, exist := c.Get(CLAIM_KEY)
	if !exist {
		return 0, false
	}
	r, _ := utf8.DecodeRune([]byte(claims.(*jwt.UserClaims).ID))
	return uint(r), true
}

func AuthUser(c *gin.Context, userPtr model.IUser) (isAbort bool) {
	userId, exist := authClaimID(c)
	if !exist {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": UserNotLoginError{}.Error()})
		return true
	}
	if err := userPtr.Scan(userId); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": UserNotExistError{}.Error()})
		return true
	}
	return false
}

func AuthUserId(c *gin.Context) (userId uint, isAbort bool) {
	exist := false
	userId, exist = authClaimID(c)
	if !exist {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": UserNotLoginError{}.Error()})
		return 0, true
	}
	return userId, false
}

func Revoke(c *gin.Context) error {
	j := jwt.NewJWT(config.GetString("auth.sign_key"))
	if tokenString, exist := c.Get(TOKEN_KEY); exist {
		if token, ok := tokenString.(string); ok {
			if err := j.RevokeToken(token); err == nil {
				c.Header("Authorization", "")
				return nil
			}
		}
	}
	return TokenRevokeError{}
}
