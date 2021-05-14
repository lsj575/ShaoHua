package initialize

import (
	"api/user-web/global"
	"api/user-web/middlewares"
	"api/user-web/router"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Routers() *gin.Engine {
	totalRouter := gin.Default()
	// 负责给consul做健康检查
	totalRouter.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, global.JsonSuccess("success", []int{}))
	})

	totalRouter.Use(middlewares.Cors())
	apiGroup := totalRouter.Group("v1")
	router.InitBaseRouter(apiGroup)
	router.InitUserRouter(apiGroup)
	router.InitLoginRouter(apiGroup)

	return totalRouter
}
