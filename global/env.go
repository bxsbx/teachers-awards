package global

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
)

const (
	DEV  = "dev"
	TEST = "test"
	PROD = "prod"
)

func GetEnvMode() string {
	//DREAMENV 环境标识
	env := os.Getenv("DREAMENV")
	fmt.Println("env:", env)
	switch env {
	case PROD:
		return gin.DebugMode
	case TEST:
		return gin.TestMode
	default:
		return gin.DebugMode
	}
}
