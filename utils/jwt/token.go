package jwt

import (
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/totoval/framework/helpers/cache"
	"gopkg.in/dgrijalva/jwt-go.v3"
	"time"
)

const ExpiredTime time.Duration = 4 * time.Hour //@todo move to configration
const RefreshExpiredTime time.Duration = 10 * time.Minute
const MaxRefreshTimes uint = 1

const REFRESH_TOKEN_CACHE_KEY = "TOTOVAL_REFRESH_TOKEN_CACHE_KEY_%s"

func refreshTokenCacheKey(tokenMd5 string) string {
	return fmt.Sprintf(REFRESH_TOKEN_CACHE_KEY, tokenMd5)
}

type refreshToken struct {
	Name string
	RefreshTimes uint
}

var (
	TokenExpired     error  = errors.New("Token is expired")
	TokenNotValidYet error  = errors.New("Token not active yet")
	TokenMalformed   error  = errors.New("That's not even a token")
	TokenInvalid     error  = errors.New("Couldn't handle this token")
	TokenNoSet       error  = errors.New("Token is not set")
)
type UserClaims struct {
	ID   string    `json:"id"`
	Name string `json:"name"`
	//Email string `json:"email"`
	Revoked bool `json:"revoked"`
	jwt.StandardClaims
}
func (c *UserClaims) Revoke() {
	c.Revoked = true
}
func (c UserClaims) Valid() error {
	vErr := new(jwt.ValidationError)
	now := jwt.TimeFunc().Unix()

	// The claims below are optional, by default, so if they are set to the
	// default value in Go, let's not fail the verification for them.
	if c.VerifyExpiresAt(now, false) == false {
		delta := time.Unix(now, 0).Sub(time.Unix(c.ExpiresAt, 0))
		vErr.Inner = fmt.Errorf("token is expired by %v", delta)
		vErr.Errors |= jwt.ValidationErrorExpired
	}

	if c.VerifyIssuedAt(now, false) == false {
		vErr.Inner = fmt.Errorf("Token used before issued")
		vErr.Errors |= jwt.ValidationErrorIssuedAt
	}

	if c.VerifyNotBefore(now, false) == false {
		vErr.Inner = fmt.Errorf("token is not valid yet")
		vErr.Errors |= jwt.ValidationErrorNotValidYet
	}


	if vErr.Errors == 0 {
		return nil
	}

	return vErr
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
	jwt.TimeFunc = time.Now
	now := time.Now()
	claims := UserClaims{
		id,
		name,
		false,
		jwt.StandardClaims {
			IssuedAt: now.Unix(),
			NotBefore: now.Unix(),
			ExpiresAt: now.Add(ExpiredTime).Unix(),
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
				spew.Dump(ve)
				return nil, TokenInvalid
			}
		}
	}

	if claims, ok := token.Claims.(*UserClaims); ok && token.Valid && !claims.Revoked {
		return claims, nil
	}
	return nil, TokenInvalid
}
func (j *JWT) RefreshToken(tokenString string) (string, error) {
	if tokenString == "" {
		return "", TokenNoSet
	}

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
		//claims.StandardClaims.ExpiresAt = time.Now().Add(ExpiredTime).Unix()
		return j.CreateToken(claims.ID, claims.Name)
	}
	return "", TokenInvalid
}

func (j *JWT) tokenMd5(tokenString string) string {
	md5Slice := md5.Sum([]byte(tokenString))
	return string(md5Slice[:])
}
func (j *JWT) checkTokenRefreshTimesValid(tokenString string) bool {
	tokenMd5 := j.tokenMd5(tokenString)
	if cache.Has(tokenMd5) {
		if cache.Get(tokenMd5).(int) >= int(MaxRefreshTimes) {
			return false
		}
	}
	return true
}
func (j *JWT) recordTokenRefreshTimes(tokenString string) {
	tokenMd5 := j.tokenMd5(tokenString)
	var increment int = 1
	if cache.Has(refreshTokenCacheKey(tokenMd5)) {
		cache.Increment(refreshTokenCacheKey(tokenMd5), increment)
	}else{
		cache.Add(refreshTokenCacheKey(tokenMd5), increment, time.Now().Add(ExpiredTime).Add(RefreshExpiredTime))
	}
}

func (j *JWT) RefreshTokenUnverified(tokenString string) (string, error) {
	if tokenString == "" {
		return "", TokenNoSet
	}

	// check this token has been refreshed times
	if !j.checkTokenRefreshTimesValid(tokenString) {
		return "", TokenInvalid
	}

	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, &UserClaims{})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*UserClaims); ok {
		//claims.StandardClaims.ExpiresAt = time.Now().Add(ExpiredTime).Unix()
		// after refresh expired time, then cannot do auto refresh
		if !time.Now().After(time.Unix(claims.ExpiresAt, 0).Add(RefreshExpiredTime)) {
			newToken, err := j.CreateToken(claims.ID, claims.Name)
			if err != nil{
				return "", err
			}

			// record the token refresh times
			j.recordTokenRefreshTimes(tokenString)

			return newToken, nil
		}
		return "", TokenExpired
	}
	return "", TokenInvalid
}

func (j *JWT) RevokeToken(tokenString string) error {
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, &UserClaims{})
	if err != nil {
		return err
	}
	if claims, ok := token.Claims.(*UserClaims); ok {
		claims.Revoke()
		return nil
	}
	return TokenInvalid
}