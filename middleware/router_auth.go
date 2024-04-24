package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"teachers-awards/common/errorz"
	"teachers-awards/global"
)

// 对路由进行简单的权限控制（仅针对增删改操作）
func RouterAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.String()
		if roles, ok := global.RouterAuthMap[path]; ok {
			userInfo := global.GetUserInfo(global.GetContext(c))
			flag := false
			for _, userRole := range userInfo.UserRoles {
				for _, role := range roles {
					if userRole == role {
						flag = true
						break
					}
				}
			}
			if flag {
				c.Next()
			} else {
				c.AbortWithStatusJSON(http.StatusForbidden, errorz.Code(errorz.NO_PERMISSIONS_ACCESS).Error())
			}
		} else { // 不需要权限的路由直接通过
			c.Next()
		}
	}
}
