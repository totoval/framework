package jwt

import (
	"errors"
	"github.com/totoval/framework/helpers"
	"gopkg.in/dgrijalva/jwt-go.v3"
	"time"
)

const ExpiredTime time.Duration = 1 * time.Hour

var (
	TokenExpired     error  = errors.New("Token is expired")
	TokenNotValidYet error  = errors.New("Token not active yet")
	TokenMalformed   error  = errors.New("That's not even a token")
	TokenInvalid     error  = errors.New("Couldn't handle this token:")
	TokenNoSet       error  = errors.New("Token is not set")
)
type UserClaims struct {
	ID   string    `json:"id"`
	Name string `json:"name"`
	//Email string `json:"email"`
	jwt.StandardClaims
}

type JWT struct {
	SigningKey []byte
}

func NewJWT(signKey string) *JWT {
	return &JWT{
		[]byte(signKey),
	}
}
func (j *JWT) CreateToken(id string, name string) (string, error) {
	now := time.Now()
	claims := UserClaims{
		id,
		name,
		jwt.StandardClaims {
			ExpiresAt: now.Add(ExpiredTime).Unix(),
			Issuer: "totoval",
			IssuedAt: now.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}
func (j *JWT) ParseToken(tokenString string) (*UserClaims, error) {
	if tokenString == "" {
		return nil, TokenNoSet
	}

	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				helpers.Dump(ve)
				return nil, TokenInvalid
			}
		}
	}
	if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, TokenInvalid
}
func (j *JWT) RefreshToken(tokenString string) (string, error) {
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		jwt.TimeFunc = time.Now
		//claims.StandardClaims.ExpiresAt = time.Now().Add(ExpiredTime).Unix()
		return j.CreateToken(claims.ID, claims.Name)
	}
	return "", TokenInvalid
}

//@todo RevokeToken