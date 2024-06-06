package main

import (
	"fmt"
	"log"
	"net/http"
	"teachers-awards/global"
	"teachers-awards/router"
	"time"
)

//	@title			Swagger teachers-awards API
//	@version		1.0
//	@description	This is a teachers-awards server.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger_controller.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host	localhost:7850

//	@securityDefinitions.basic	BasicAuth

//	@externalDocs.description	OpenAPI
//	@externalDocs.url			https://swagger.io/resources/open-api/
func main() {
	err := global.InitConfig()
	if err != nil {
		log.Fatalf("读取配置失败，err:%v", err)
	}

	global.InitGlobalVal()

	//链路跟踪
	closer := global.InitTract()
	if closer != nil {
		defer closer.Close()
	}

	routers := router.Routers()
	server := &http.Server{
		Addr:           fmt.Sprintf(":%v", global.ServeCfg.Default.AppPort),
		Handler:        routers,
		ReadTimeout:    60 * time.Second,
		WriteTimeout:   60 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	if err := server.ListenAndServe(); err != nil {
		fmt.Println(err)
	}
}
