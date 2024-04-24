package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"teachers-awards/common/jwt"
	"teachers-awards/global"
	"time"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		//从请求头获取token
		token := c.Request.Header.Get(global.AuthToken)
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, jwt.TokenInvalid.Error())
			return
		}
		// parseToken 解析token包含的信息
		claims, err := global.Jwt.ParseToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, err.Error())
			return
		}

		userInfo := global.GetUserInfoFromClaims(claims)
		ctx := global.GetContext(c)
		ctx = global.SetUserInfo(ctx, userInfo)
		global.SetContext(c, ctx)

		//单点登录结合redis，无需使用注释
		//tokenKey := global.JwtKey + userInfo.UserId
		//redisToken, err := global.RedisClient.Get(ctx, tokenKey).Result()
		//if err != nil {
		//	c.AbortWithStatusJSON(http.StatusUnauthorized, err.Error())
		//	return
		//}
		//if token != redisToken {
		//	c.AbortWithStatusJSON(http.StatusUnauthorized, errorz.Code(errorz.RESP_LOGIN_EXIST).Error())
		//	return
		//}

		// RefreshTime时间区内刷新token，否则无需刷新（减少服务压力）
		unix := time.Now().Unix()
		if claims.ExpiresAt-unix < global.Jwt.RefreshTime {
			claims.ExpiresAt = unix + global.Jwt.ExpiresTime
			token, _ = global.Jwt.CreateToken(*claims)
			//global.RedisClient.Set(ctx, tokenKey, newToken, global.Jwt.ExpiresTime)
		}
		c.Header(global.AuthToken, token)
		c.Header(global.ExpiresAt, fmt.Sprintf("%v", claims.ExpiresAt))
		c.Next()
	}
}
