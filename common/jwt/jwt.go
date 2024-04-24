package jwt

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
)

// 自定义jwt 所需Claims结构
type CustomClaims struct {
	DataJson []byte
	jwt.StandardClaims
}

type Jwt struct {
	CustomSecret string
	ExpiresTime  int64 // token有效时间区 (有效期一周)
	RefreshTime  int64 // token刷新时间区
}

var (
	TokenExpired     = errors.New("token已过期")
	TokenNotValidYet = errors.New("token未激活")
	TokenMalformed   = errors.New("非法token")
	TokenInvalid     = errors.New("无效token")
)

// 创建一个token
func (j *Jwt) CreateToken(claims CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.CustomSecret))
}

// 解析 token
func (j *Jwt) ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (i interface{}, e error) {
		return []byte(j.CustomSecret), nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}
		}
	}
	if token != nil {
		if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
			return claims, nil
		}
	}
	return nil, TokenInvalid
}
