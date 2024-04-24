package swagger

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"strings"
	"teachers-awards/global"
	"teachers-awards/swagger/docs"
)

func Swagger(Router *gin.Engine) {
	host := global.ServeCfg.TeachersAwards.Domain
	docs.SwaggerInfo.Host = strings.TrimLeft(host, "https://")
	//http://127.0.0.1:7850/swagger/index.html访问
	Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
