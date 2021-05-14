package router

import (
	"api/user-web/api"
	"github.com/gin-gonic/gin"
)

func InitLoginRouter(router *gin.RouterGroup) {
	router.POST("/login", api.PasswordLogin)
	router.POST("/non_password_login", api.NonPasswordLogin)
	router.POST("/register", api.Register)
}
