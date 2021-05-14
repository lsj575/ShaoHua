package router

import (
	"api/user-web/api"
	"api/user-web/middlewares"
	"github.com/gin-gonic/gin"
)

func InitUserRouter(router *gin.RouterGroup) {
	UserRouter := router.Group("user")
	{
		UserRouter.GET("/:id", middlewares.JWTAuth(), api.GetUserById)
		UserRouter.GET("/email/:email", middlewares.JWTAuth(), api.GetUserByEmail)
		UserRouter.GET("/current", middlewares.JWTAuth(), api.GetCurrentUser)
		UserRouter.GET("/check", api.CheckUserExistByEmail)
	}
}
